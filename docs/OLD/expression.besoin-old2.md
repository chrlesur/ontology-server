
# Expression de besoin - Serveur Ontology avec API et interface web

## 1. Objectif général

Développer un logiciel serveur en Go avec une API RESTful et une interface web permettant de :

1. Charger et gérer des fichiers d'ontologie (TSV, RDF, OWL) et leurs fichiers JSON associés.
2. Effectuer des recherches dans ces ontologies basées sur des requêtes utilisateur.
3. Afficher pour chaque résultat :
   - L'ontologie concernée
   - Les éléments correspondants avec leurs types et descriptions
   - Le contexte de chaque élément extrait du fichier JSON associé
4. Permettre la sélection d'éléments spécifiques pour afficher plus de détails.

## 2. Fonctionnalités principales

### 2.1 Gestion des fichiers d'ontologie (Priorité Haute)

#### 2.1.1 Importation des fichiers
- Supporter l'importation de fichiers aux formats TSV, RDF et OWL
- Détecter automatiquement le format du fichier basé sur son extension et son contenu
- Valider l'intégrité et la structure du fichier avant l'importation
- Gérer les erreurs d'importation avec des messages détaillés

#### 2.1.2 Gestion des fichiers JSON associés
- Rechercher automatiquement le fichier JSON associé (NOMDUDOCUMENT_context.json) dans le même répertoire
- Si trouvé, importer et lier le fichier JSON à l'ontologie correspondante
- Utiliser le fichier JSON pour extraire le contexte des éléments

#### 2.1.3 Extraction et stockage des métadonnées
Pour chaque fichier importé, extraire et stocker :
- Nom du fichier
- Format du fichier (TSV, RDF, OWL)
- Nombre d'éléments dans l'ontologie
- Nombre de relations dans l'ontologie

#### 2.1.4 Gestion des ontologies
- Implémenter des fonctions pour ajouter, supprimer et mettre à jour les ontologies
- Attribuer un identifiant unique à chaque ontologie importée
- Maintenir une liste des ontologies chargées en mémoire

### 2.2 Traitement des fichiers (Priorité Haute)

#### 2.2.1 Parsing des fichiers d'ontologie
- Implémenter des parsers pour les formats TSV, RDF et OWL
- Extraire les éléments et les relations selon la structure définie :
  ```go
  type OntologyElement struct {
      Name        string
      Type        string
      Positions   []int
      Description string 
  }

  type Relation struct {
      Source      string
      Type        string
      Target      string
      Description string
  }

  type Ontology struct {
      Elements  []*OntologyElement
      Relations []*Relation
  }
  ```

#### 2.2.2 Parsing des fichiers JSON de contexte
- Implémenter un parser JSON qui extrait :
  ```json
  {
    "position": int,
    "before": [string],
    "after": [string],
    "element": string,
    "length": int
  }
  ```
- Associer chaque entrée JSON à l'élément correspondant dans l'ontologie

### 2.3 API RESTful (Priorité Haute)

#### 2.3.1 Recherche avancée
- GET /api/v1/search
  - Description : Effectuer une recherche dans les ontologies
  - Paramètres de requête : 
    - q : terme de recherche
    - ontology_id : ID de l'ontologie (optionnel)
    - type : type d'élément à rechercher (optionnel)
    - context_size : nombre de mots avant et après l'élément pour le contexte (défaut : 5)
  - Réponse : 
    - Liste des correspondances avec :
      - ID de l'ontologie
      - Élément trouvé (nom, type, description)
      - Contexte de l'élément extrait du fichier JSON
      - Positions de l'élément dans le fichier

#### 2.3.2 Détails d'un élément
- GET /api/v1/elements/{element_id}
  - Description : Obtenir les détails complets d'un élément spécifique
  - Réponse :
    - Toutes les informations de l'élément
    - Contexte étendu extrait du fichier JSON
    - Relations associées à cet élément

### 2.4 Moteur de recherche avancé (Priorité Haute)

- Implémenter une fonction de recherche qui :
  - Accepte un terme de recherche, un ID d'ontologie optionnel et un type d'élément optionnel
  - Recherche dans les noms, types et descriptions des éléments
  - Retourne une liste de correspondances avec toutes les informations nécessaires pour l'affichage détaillé

### 2.5 Interface utilisateur web (Priorité Moyenne)

Développer une interface web en HTML/JavaScript avec :

#### 2.5.1 Page d'accueil
- Champ de recherche principal
- Options de filtrage (par ontologie, par type d'élément)

#### 2.5.2 Page de résultats de recherche
- Liste des résultats montrant :
  - Nom de l'élément
  - Type de l'élément
  - Extrait court du contexte
- Possibilité de cliquer sur un résultat pour voir plus de détails

#### 2.5.3 Page de détails d'un élément
- Affichage complet des informations de l'élément
- Contexte étendu extrait du fichier JSON
- Liste des relations associées à cet élément

## 3. Spécifications techniques

### 3.1 Backend (Priorité Haute)

#### 3.1.1 Langage et framework
- Utiliser Go version 1.16 ou supérieure
- Utiliser le package net/http standard de Go pour le serveur HTTP
- Utiliser le package encoding/json pour le traitement JSON

#### 3.1.2 Structure du projet
```
/cmd
  /server
    main.go
/internal
  /api
    handlers.go
    router.go
  /config
    config.go
  /logger
    logger.go
  /models
    ontology.go
    element.go
    relation.go
  /parser
    tsv.go
    rdf.go
    owl.go
    json.go
  /search
    engine.go
  /storage
    memory.go
/web
  index.html
  script.js
  style.css
```

#### 3.1.3 Gestion de la concurrence
- Utiliser des goroutines pour le traitement parallèle des fichiers
- Implémenter des mécanismes de synchronisation (mutex) pour l'accès aux données partagées

#### 3.1.4 Gestion des erreurs
- Utiliser des erreurs personnalisées avec des messages clairs
- Implémenter un middleware pour la gestion centralisée des erreurs

### 3.2 Stockage des données (Priorité Haute)

#### 3.2.1 Structure de données en mémoire
Utiliser les structures définies dans la section 2.2.1 pour représenter les ontologies en mémoire.

### 3.3 Logging (Priorité Moyenne)

Implémenter un package logger avec :
- Niveaux de log : Debug, Info, Warning, Error
- Rotation des fichiers de log quotidienne
- Format de log : `[LEVEL] [TIMESTAMP] [FILE:LINE] Message`

### 3.4 Configuration (Priorité Moyenne)

Utiliser un fichier config.yaml avec les paramètres suivants :
```yaml
server:
  port: 8080
  host: localhost
logging:
  level: info
  directory: ./logs
storage:
  temp_directory: ./temp
```

### 3.5 Tests (Priorité Haute)

- Implémenter des tests unitaires pour chaque package en parallèle avec le développement
- Utiliser le package `testing` de Go
- Viser une couverture de code d'au moins 80% pour chaque package
- Implémenter des benchmarks pour les fonctions critiques (parsing, recherche)
- Créer des mocks et des stubs pour les dépendances externes lors des tests

### 3.6 Documentation (Priorité Basse)

- Utiliser les commentaires de documentation Go pour toutes les fonctions exportées
- Créer un README.md détaillé avec :
  - Description du projet
  - Instructions d'installation
  - Guide de configuration
  - Exemples d'utilisation de l'API
  - Instructions pour exécuter les tests

### 3.7 Performance (Priorité Haute)
- Optimiser les algorithmes de recherche pour gérer efficacement de grandes ontologies
- Implémenter un système de mise en cache des résultats de recherche fréquents
- Utiliser des index en mémoire pour accélérer les recherches

## 4. Contraintes techniques et directives de développement

- Limiter chaque fichier de code source Go à un maximum de 500 lignes
- Limiter chaque fonction à un maximum de 50 lignes de code
- Suivre les conventions de nommage et de formatage standard de Go
- Utiliser `gofmt` pour formater le code
- Utiliser `golint` et `go vet` pour la vérification du code
- Gérer toutes les erreurs de manière appropriée, sans utiliser de `panic`
- Utiliser des interfaces pour les composants principaux pour faciliter les tests et l'extensibilité future

## 5. Plan de développement

1. Configuration du projet et mise en place de la structure de base
2. Implémentation du package logger et de la gestion des configurations
3. Développement des parsers de fichiers (TSV, RDF, OWL, JSON)
4. Implémentation de la structure de données en mémoire pour les ontologies
5. Développement de l'API RESTful avec endpoints de recherche avancée
6. Implémentation du moteur de recherche avancé
7. Création de l'interface utilisateur web
8. Optimisation des performances et mise en place du système de cache
9. Finalisation des tests unitaires et d'intégration
10. Documentation et nettoyage du code
11. Tests finaux, débogage et ajustements

Pour chaque étape du développement, des tests unitaires correspondants doivent être créés et exécutés.
