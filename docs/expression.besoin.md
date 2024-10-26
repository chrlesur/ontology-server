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

### 3.5 Intégration LLM
- Intégrer le LLM existant pour la composition des prompts de recherche
- Adapter le système pour prendre en compte la recherche dans plusieurs ontologies

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

## 4. Contraintes techniques et directives de développement

- Limiter chaque fichier de code source Go à un maximum de 3000 tokens
- Limiter chaque package à un maximum de 10 méthodes exportées
- Limiter chaque méthode à un maximum de 80 lignes de code
- Suivre les meilleures pratiques et les modèles idiomatiques de Go
- Utiliser le package 'internal/logger' pour toute journalisation
- Définir toutes les valeurs configurables dans le package 'internal/config'
- Utiliser les constantes définies dans le package 'internal/i18n' pour les messages utilisateur
- Optimiser le code pour la performance, particulièrement pour le traitement de grands documents et ontologies
- Préparer le code pour de futurs efforts de localisation

## 5. Évolutivité et maintenance

- Concevoir le système de manière à faciliter l'ajout de nouvelles fonctionnalités
- Prévoir des mécanismes de sauvegarde et de restauration des données
- Implémenter un système de logging détaillé pour le suivi des opérations et le débogage
- Planifier des phases de refactoring régulières pour améliorer la structure du code et sa maintenabilité
- Prévoir des mécanismes pour mettre à jour l'index de recherche lorsque les ontologies sont modifiées
- Effectuer des revues de code régulières pour identifier et résoudre les problèmes potentiels tôt dans le processus de développement