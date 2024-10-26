# Expression de besoin révisée - Serveur Ontology avec API et interface web

## 1. Objectif général

Développer un logiciel serveur avec API s'appuyant sur le code existant du client Ontology, capable de gérer plusieurs ontologies et de fournir un moteur de recherche avancé et performant, avec un support robuste pour les requêtes complexes.

## 2. Fonctionnalités principales

### 2.1 Importation et gestion des données
- Importer et gérer plusieurs fichiers d'ontologie aux formats TSV, RDF ou OWL
- Importer automatiquement les fichiers de contexte JSON associés (NOMDUDOCUMENT_context.json) pour chaque ontologie
- Extraire et stocker les métadonnées de chaque fichier :
  - Localisation
  - Taille
  - SHA256
- Permettre l'ajout, la mise à jour et la suppression d'ontologies
- "La création de cette ontologie génère deux types de fichiers de sortie : un fichier .tsv (Tab-Separated Values) et un fichier .json (JavaScript Object Notation). Ces fichiers représentent différents aspects de l'ontologie extraite du texte source.

1. Fichier .tsv :

Le fichier .tsv contient les éléments de base de l'ontologie dans un format tabulaire. Chaque ligne représente un élément distinct et contient les informations suivantes, séparées par des tabulations :

- Nom de l'élément
- Type de l'élément (Concept, Rôle, Organisation, etc.)
- Description de l'élément
- Liste des positions dans le texte où l'élément apparaît

Par exemple :
Qualification_Juridique    Concept    Processus de détermination du statut légal d'une situation ou d'un individu    1
Agent_Service_Public    Concept    Personne travaillant pour un service public    23,47,189,24,48,62,81,117,120,154,177,190,208,244,443

Ce format permet une visualisation rapide des éléments de l'ontologie et de leur distribution dans le texte.

2. Fichier .json :

Le fichier .json offre une représentation plus détaillée et structurée de l'ontologie, incluant à la fois les éléments et leurs relations. Il se compose de deux parties principales :

a) Liste des éléments :
Chaque élément est représenté par un objet JSON contenant :
- "name": Nom de l'élément
- "type": Type de l'élément
- "description": Description de l'élément
- "positions": Liste des positions dans le texte

b) Liste des relations :
Chaque relation est représentée par un objet JSON contenant :
- "source": Élément source de la relation
- "type": Type de relation
- "target": Élément cible de la relation
- "description": Description de la relation

Par exemple :
{
  "elements": [
    {"name": "Qualification_Juridique", "type": "Concept", "description": "Processus de détermination du statut légal d'une situation ou d'un individu", "positions": [1]},
    {"name": "Agent_Service_Public", "type": "Concept", "description": "Personne travaillant pour un service public", "positions": [23,47,189,24,48,62,81,117,120,154,177,190,208,244,443]}
  ],
  "relations": [
    {"source": "Qualification_Juridique", "type": "Détermine", "target": "Statut_Légal", "description": "La qualification juridique détermine le statut légal d'une situation"}
  ]
}

Cette structure JSON permet une représentation complète et flexible de l'ontologie, facilitant son utilisation dans diverses applications d'analyse et de traitement du langage naturel."


### 2.2 Gestion des ontologies
- Importer chaque ontologie avec tous ses éléments et leurs positions
- Stocker les informations de positions en base de données pour chaque ontologie
- Stocker les ontologies complètes en base de données
- Gérer les relations entre les différentes ontologies si nécessaire

### 2.3 Moteur de recherche avancé
- Implémenter une API de recherche capable de chercher dans toutes les ontologies stockées
- Supporter des requêtes simples, composées et complexes (AND, OR, NOT, phrases exactes)
- Gérer les requêtes imbriquées avec parenthèses
- Utiliser un système d'indexation inversée pour optimiser les performances de recherche
- Implémenter un système de scoring avancé pour les correspondances exactes et partielles, y compris pour les requêtes NOT
- Permettre la recherche dans une ontologie spécifique ou dans toutes les ontologies
- Optimiser les performances pour de grandes ontologies
- Implémenter une logique de déduplication des résultats
- Retourner les résultats de recherche incluant :
  - Nom du document/ontologie
  - Position dans le document
  - Extrait de texte (30 mots avant et après la position par défaut)
  - Score de pertinence

### 2.4 Interface utilisateur
- Développer une interface web en HTML5/React
- Fournir une interface pour :
  - L'importation et la gestion de multiples ontologies
  - La recherche avancée dans une ou plusieurs ontologies, supportant les requêtes complexes
  - L'affichage des résultats de recherche avec mise en évidence des correspondances
  - La visualisation et la navigation dans les ontologies

## 3. Spécifications techniques

### 3.1 Backend
- Développer un serveur en Go (dernière version stable)
- Utiliser le framework Gin pour l'API RESTful
- Utiliser et adapter le code existant du client Ontology
- Implémenter une API RESTful avec des endpoints pour gérer plusieurs ontologies
- Utiliser les goroutines et les canaux pour le traitement concurrent et l'optimisation des performances
- Assurer la compatibilité avec différents LLM et leurs limites de contexte spécifiques

### 3.2 Stockage et indexation
- Implémenter une interface de stockage abstraite pour permettre différentes implémentations
- Implémenter initialement un stockage en mémoire avec indexation inversée
- Prévoir l'intégration future d'une base de données (par exemple PostgreSQL) pour un stockage persistant
- Optimiser l'indexation pour les grandes ontologies et les requêtes complexes

### 3.3 Moteur de recherche
- Implémenter un système d'indexation inversée efficace
- Développer un parser de requêtes récursif capable de gérer des requêtes complexes et imbriquées
- Implémenter une structure Query pour représenter les requêtes imbriquées
- Implémenter un système de scoring avancé pour améliorer la pertinence des résultats, y compris pour les requêtes NOT
- Optimiser les performances pour les grandes ontologies
- Implémenter un système de cache pour les requêtes fréquentes
- Développer une logique de déduplication des résultats

### 3.4 Frontend
- Développer une interface web en React
- Implémenter des composants pour la gestion et la visualisation de multiples ontologies
- Créer une interface de recherche avancée supportant les requêtes complexes et imbriquées

### 3.6 Gestion des erreurs et validation
- Implémenter un système robuste de gestion des erreurs, notamment pour les requêtes mal formées
- Développer un système de validation des entrées utilisateur
- Utiliser des middlewares pour la gestion globale des erreurs et la validation des requêtes

### 3.7 Tests et qualité du code
- Implémenter des tests unitaires exhaustifs pour tous les packages
- Développer des tests d'intégration pour valider le comportement du système complet
- Implémenter des tests de performance pour valider l'efficacité avec de grandes ontologies
- Ajouter des tests spécifiques pour les requêtes booléennes complexes
- Utiliser des outils de linting et de formatage de code
- Implémenter des benchmarks pour mesurer et optimiser les performances

### 3.8 Documentation
- Fournir des commentaires de documentation conformes aux standards GoDoc
- Maintenir un README à jour avec les informations sur la structure du projet et son utilisation
- Créer un guide d'utilisation détaillé pour le moteur de recherche, incluant des exemples de requêtes complexes

### 3.9 Sécurité
- Implémenter une structure de base pour l'authentification et l'autorisation

### Journalisation

 voici une description du package logger que le LLM devrait mettre en œuvre :

1. Structure principale :
   - Le package doit définir une structure `Logger` qui encapsule toutes les fonctionnalités de journalisation.
   - Utiliser un modèle de singleton pour s'assurer qu'il n'y a qu'une seule instance de `Logger`.

2. Niveaux de journalisation :
   - Définir une énumération `LogLevel` avec quatre niveaux : Debug, Info, Warning, et Error.
   - Implémenter une méthode pour définir et obtenir le niveau de journalisation actuel.

3. Configuration :
   - Permettre la configuration du répertoire de journalisation via un fichier de configuration externe.
   - Créer automatiquement le répertoire de journalisation s'il n'existe pas.

4. Fichiers de journalisation :
   - Créer un nouveau fichier journal chaque jour, nommé avec la date actuelle.
   - Écrire les messages de journal à la fois dans le fichier et sur la console (stdout).

5. Méthodes de journalisation :
   - Implémenter des méthodes distinctes pour chaque niveau de journalisation : Debug(), Info(), Warning(), et Error().
   - Inclure les informations de fichier et de ligne dans chaque message de journal.

6. Gestion des erreurs :
   - Gérer gracieusement les erreurs lors de la création de fichiers ou de répertoires.

7. Concurrence :
   - Utiliser des mutex pour assurer la sécurité des threads lors de l'accès aux ressources partagées.

8. Fonctionnalités supplémentaires :
   - Implémenter une méthode UpdateProgress() pour afficher la progression sur la console.
   - Prévoir une méthode RotateLogs() pour la rotation des fichiers journaux (à implémenter selon les besoins spécifiques).

10. Fermeture propre :
    - Fournir une méthode Close() pour fermer proprement le fichier journal.

11. Parsing du niveau de log :
    - Implémenter une fonction ParseLevel() pour convertir une chaîne de caractères en LogLevel.

Le LLM devrait être capable de créer un package de journalisation robuste et thread-safe, capable de gérer différents niveaux de journalisation, d'écrire dans des fichiers et sur la console, et d'inclure des informations utiles comme l'emplacement du fichier et le numéro de ligne dans chaque message de journal.

## 4. Contraintes techniques et directives de développement

- Limiter chaque fichier de code source Go à un maximum de 3000 tokens
- Limiter chaque package à un maximum de 10 méthodes exportées
- Limiter chaque méthode à un maximum de 80 lignes de code
- Suivre les meilleures pratiques et les modèles idiomatiques de Go
- Définir toutes les valeurs configurables dans le package 'internal/config'
- Optimiser le code pour la performance, particulièrement pour le traitement de grands documents et ontologies
- tout le code doit êre commenté
- Tout le code doit s'appuyer sur un package logger 

## 5. Évolutivité et maintenance

- Concevoir le système de manière à faciliter l'ajout de nouvelles fonctionnalités
- Prévoir des mécanismes de sauvegarde et de restauration des données
- Implémenter un système de logging détaillé pour le suivi des opérations et le débogage
- Planifier des phases de refactoring régulières pour améliorer la structure du code et sa maintenabilité
- Prévoir des mécanismes pour mettre à jour l'index de recherche lorsque les ontologies sont modifiées
- Effectuer des revues de code régulières pour identifier et résoudre les problèmes potentiels tôt dans le processus de développement

Directives générales (à suivre impérativement pour toutes les étapes du projet) :
2. Assurez-vous qu'aucun fichier de code source ne dépasse 3000 tokens.
3. Limitez chaque package à un maximum de 10 méthodes exportées.
4. Aucune méthode ne doit dépasser 80 lignes de code.
5. Suivez les meilleures pratiques et les modèles idiomatiques de Go.
6. Tous les messages visibles par l'utilisateur doivent être en anglais.
7. Chaque fonction, méthode et type exporté doit avoir un commentaire de documentation conforme aux standards GoDoc.
8. Utilisez le package 'internal/logger' pour toute journalisation. Implémentez les niveaux de log : debug, info, warning, et error.
9. Toutes les valeurs configurables doivent être définies dans le package 'internal/config'.
10. Gérez toutes les erreurs de manière appropriée, en utilisant error wrapping lorsque c'est pertinent.
14. Implémentez des tests unitaires pour chaque nouvelle fonction ou méthode.
15. Veillez à ce que le code soit sécurisé, en particulier lors du traitement des entrées utilisateur.