# Expression de besoin - Serveur Ontology avec API et interface web

## 1. Objectif général

Développer un logiciel serveur en Go avec une API RESTful et une interface web permettant de :

1. Charger et gérer des fichiers d'ontologie (TSV et JSON) enrichis de leurs positions.
2. Effectuer des recherches dans ces ontologies basées sur des requêtes utilisateur.
3. Afficher pour chaque résultat :
   - L'ontologie concernée
   - Les éléments correspondants avec leurs types et descriptions
   - Le nom du fichier source et ses caractéristiques
   - Le contexte de chaque élément extrait du fichier JSON associé
4. Permettre la sélection d'éléments spécifiques pour afficher plus de détails.

## 2. Fonctionnalités principales

### 2.1 Gestion des fichiers d'ontologie (Priorité Haute)

#### 2.1.1 Importation des fichiers
- Supporter l'importation de fichiers aux formats TSV et JSON
- Détecter automatiquement le format du fichier basé sur son extension
- Valider l'intégrité et la structure du fichier avant l'importation
- Gérer les erreurs d'importation avec des messages détaillés

#### 2.1.2 Parsing des fichiers TSV
- Implémenter un parser TSV robuste qui extrait pour chaque ligne :
  - Nom de l'élément (colonne 1)
  - Type de l'élément (colonne 2)
  - Description de l'élément (colonne 3)
  - Positions (colonne 4) : parser la liste d'entiers séparés par des virgules

#### 2.1.3 Parsing des fichiers JSON
- Implémenter un parser JSON qui extrait pour chaque objet :
  - position : index du mot dans le texte
  - before : tableau de mots précédant l'élément
  - after : tableau de mots suivant l'élément
  - element : nom de l'élément de l'ontologie
  - length : nombre de mots que l'élément occupe

#### 2.1.4 Extraction et stockage des métadonnées
Pour chaque fichier importé, extraire et stocker :
- Nom du fichier
- Chemin complet
- Taille en octets
- Hash SHA256
- Date et heure d'importation
- Format du fichier (TSV ou JSON)
- Nombre d'éléments dans l'ontologie

#### 2.1.5 Gestion des ontologies
- Implémenter des fonctions pour ajouter, supprimer et mettre à jour les ontologies
- Attribuer un identifiant unique à chaque ontologie importée
- Maintenir une liste des ontologies chargées en mémoire

### 2.2 API RESTful (Priorité Haute)

#### 2.2.1 Recherche avancée
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
      - Nom de l'ontologie
      - Nom du fichier source et ses caractéristiques (taille, hash, etc.)
      - Élément trouvé (nom, type, description)
      - Contexte de l'élément extrait du fichier JSON
      - Position de l'élément dans le fichier

#### 2.2.2 Détails d'un élément
- GET /api/v1/elements/{element_id}
  - Description : Obtenir les détails complets d'un élément spécifique
  - Réponse :
    - Toutes les informations de l'élément
    - Contexte étendu extrait du fichier JSON
    - Liens vers d'autres éléments associés dans l'ontologie

### 2.3 Moteur de recherche avancé (Priorité Haute)

- Implémenter une fonction de recherche qui :
  - Accepte un terme de recherche, un ID d'ontologie optionnel et un type d'élément optionnel
  - Utilise un algorithme de correspondance exacte et partielle
  - Supporte la recherche insensible à la casse
  - Recherche dans les noms, types et descriptions des éléments
  - Retourne une liste de correspondances avec toutes les informations nécessaires pour l'affichage détaillé

### 2.4 Interface utilisateur web (Priorité Moyenne)

Développer une interface web en HTML/JavaScript (ou React si le temps le permet) avec :

#### 2.4.1 Page d'accueil
- Champ de recherche principal
- Options de filtrage (par ontologie, par type d'élément)

#### 2.4.2 Page de résultats de recherche
- Liste des résultats montrant :
  - Nom de l'élément
  - Type de l'élément
  - Nom de l'ontologie
  - Extrait court du contexte (mots avant et après)
- Possibilité de cliquer sur un résultat pour voir plus de détails

#### 2.4.3 Page de détails d'un élément
- Affichage complet des informations de l'élément
- Contexte étendu extrait du fichier JSON
- Informations sur le fichier source (nom, taille, hash, etc.)
- Liens vers d'autres éléments associés dans l'ontologie

## 3. Spécifications techniques

### 3.1 Backend (Priorité Haute)

#### 3.1.1 Langage et framework
- Utiliser Go version 1.16 ou supérieure
- Utiliser le package net/http standard de Go pour le serveur HTTP
- Utiliser le package encoding/json pour le traitement JSON

#### 3.1.2 Structure du projet
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
/parser
tsv.go
json.go
/search
engine.go
/storage
memory.go
/web
index.html
script.js
style.css


#### 3.1.3 Gestion de la concurrence
- Utiliser des goroutines pour le traitement parallèle des fichiers
- Implémenter des mécanismes de synchronisation (mutex) pour l'accès aux données partagées

#### 3.1.4 Gestion des erreurs
- Utiliser des erreurs personnalisées avec des messages clairs
- Implémenter un middleware pour la gestion centralisée des erreurs

### 3.2 Stockage des données (Priorité Haute)

#### 3.2.1 Structure de données en mémoire
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
    ID          string
    Name        string
    Filename    string
    Format      string
    Size        int64
    SHA256      string
    ImportedAt  time.Time
    Elements    []*OntologyElement
    Relations   []*Relation
}

type JSONContext struct {
    Position int      `json:"position"`
    Before   []string `json:"before"`
    After    []string `json:"after"`
    Element  string   `json:"element"`
    Length   int      `json:"length"`
}
3.3 Logging (Priorité Moyenne)
Implémenter un package logger avec :

Niveaux de log : Debug, Info, Warning, Error
Rotation des fichiers de log quotidienne
Format de log : [LEVEL] [TIMESTAMP] [FILE:LINE] Message
3.4 Configuration (Priorité Moyenne)
Utiliser un fichier config.yaml avec les paramètres suivants :

server:
  port: 8080
  host: localhost
logging:
  level: info
  directory: ./logs
storage:
  temp_directory: ./temp
3.5 Tests (Priorité Haute)
Implémenter des tests unitaires pour chaque package en parallèle avec le développement
Utiliser le package testing de Go
Viser une couverture de code d'au moins 90% pour chaque package
Implémenter des benchmarks pour les fonctions critiques (parsing, recherche)
Créer des mocks et des stubs pour les dépendances externes lors des tests
Assurer que tous les packages ont des tests unitaires complets
Structure des fichiers de test :

/internal
  /api
    handlers_test.go
    router_test.go
  /config
    config_test.go
  /logger
    logger_test.go
  /models
    ontology_test.go
    element_test.go
  /parser
    tsv_test.go
    json_test.go
  /search
    engine_test.go
  /storage
    memory_test.go
3.6 Documentation (Priorité Basse)
Utiliser les commentaires de documentation Go pour toutes les fonctions exportées
Créer un README.md détaillé avec :
Description du projet
Instructions d'installation
Guide de configuration
Exemples d'utilisation de l'API
Instructions pour exécuter les tests
3.7 Performance (Priorité Haute)
Optimiser les algorithmes de recherche pour gérer efficacement de grandes ontologies
Implémenter un système de mise en cache des résultats de recherche fréquents
Utiliser des index en mémoire pour accélérer les recherches
4. Contraintes techniques et directives de développement
Limiter chaque fichier de code source Go à un maximum de 500 lignes
Limiter chaque fonction à un maximum de 50 lignes de code
Suivre les conventions de nommage et de formatage standard de Go
Utiliser gofmt pour formater le code
Utiliser golint et go vet pour la vérification du code
Gérer toutes les erreurs de manière appropriée, sans utiliser de panic
Utiliser des interfaces pour les composants principaux pour faciliter les tests et l'extensibilité future
5. Structure des fichiers de sortie de l'ontologie
5.1 Fichier .tsv
Le fichier .tsv contient les éléments de base de l'ontologie dans un format tabulaire. Chaque ligne représente un élément distinct et contient les informations suivantes, séparées par des tabulations :

Nom de l'élément
Type de l'élément (Concept, Rôle, Organisation, etc.)
Description de l'élément
Liste des positions dans le texte où l'élément apparaît
Exemple :

Qualification_Juridique    Concept    Processus de détermination du statut légal d'une situation ou d'un individu    1
Agent_Service_Public    Concept    Personne travaillant pour un service public    23,47,189,24,48,62,81,117,120,154,177,190,208,244,443
5.2 Fichier .json
Le fichier .json offre une représentation plus détaillée et structurée de l'ontologie. Il contient un tableau d'objets JSON, où chaque objet correspond à une occurrence d'un élément de l'ontologie dans le texte. Voici la structure de chaque objet :

position : L'index du mot dans le texte où l'élément commence.
before : Un tableau de mots qui précèdent immédiatement l'élément dans le texte.
after : Un tableau de mots qui suivent immédiatement l'élément dans le texte.
element : Le nom de l'élément de l'ontologie identifié à cette position.
length : Le nombre de mots que l'élément occupe dans le texte.
Exemple :

[
  {
    "position": 1,
    "before": ["La"],
    "after": ["juridique","sert","de","premier","plan.","Dans","le","cas","des","internationales,","le","Conseil","d'État","nous","dit","que","ce","sont","les","agents","du","service","public,","donc","principe","de","neutralité","applicable.","Dans","le"],
    "element": "Qualification_Juridique",
    "length": 1
  },
  {
    "position": 8,
    "before": ["La","qualification","juridique","sert","de","premier","plan.","Dans"],
    "after": ["cas","des","internationales,","le","Conseil","d'État","nous","dit","que","ce","sont","les","agents","du","service","public,","donc","principe","de","neutralité","applicable.","Dans","le","cas","des","simples","licenciés,","il","nous","dit"],
    "element": "Joueuse_Internationale",
    "length": 1
  }
]
Cette structure JSON permet une analyse contextuelle détaillée de chaque occurrence des éléments de l'ontologie dans le texte, facilitant diverses tâches de traitement du langage naturel et d'analyse de texte.
