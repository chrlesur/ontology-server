# Serveur d'Ontologie

## Description
Le Serveur d'Ontologie est une application sophistiquée développée en Go, conçue pour gérer, explorer et visualiser des ontologies complexes. Ce projet s'inscrit dans le domaine de la gestion des connaissances et de l'intelligence artificielle, offrant une plateforme robuste pour travailler avec des structures de données ontologiques.

### Objectif du projet
L'objectif principal de ce serveur est de fournir une interface puissante et flexible pour interagir avec des ontologies, facilitant ainsi la recherche sémantique, l'analyse de relations conceptuelles et la découverte de connaissances dans divers domaines d'application.

### Caractéristiques principales
- **Gestion avancée d'ontologies** : Capable de charger, stocker et manipuler des ontologies complexes créées par le client Ontology.
- **API RESTful complète** : Offre une interface programmatique riche pour l'interaction avec les ontologies, permettant une intégration facile dans d'autres systèmes et applications.
- **Moteur de recherche sémantique** : Implémente des algorithmes de recherche avancés pour explorer efficacement les structures ontologiques.
- **Visualisation interactive** : Propose une interface utilisateur web intuitive pour explorer visuellement les relations entre les éléments d'une ontologie.
- **Haute performance** : Conçu pour gérer efficacement de grandes quantités de données ontologiques avec une faible latence.
- **Extensibilité** : Architecture modulaire permettant l'ajout facile de nouvelles fonctionnalités et l'intégration de différents formats d'ontologie.

### Applications potentielles
Ce serveur d'ontologie peut être utilisé dans divers domaines, notamment :
- Recherche scientifique et académique
- Systèmes d'aide à la décision
- Gestion des connaissances d'entreprise
- Analyse de données biomédicales
- Traitement du langage naturel et IA
- Systèmes d'information géographique

### État du projet
**VERSION DE RECHERCHE ET DÉVELOPPEMENT**
Ce projet est actuellement en phase active de développement et de recherche. Il est conçu pour être un outil puissant pour les chercheurs et les développeurs travaillant avec des ontologies, mais peut contenir des fonctionnalités expérimentales ou en cours de perfectionnement.

## Fonctionnalités
- Chargement et gestion de fichiers d'ontologie créés par le client Ontology
- API RESTful complète pour l'interaction avec les ontologies
- Moteur de recherche avancé avec capacités de requêtes complexes
- Visualisation interactive des relations entre éléments d'ontologie
- Support pour différents formats d'ontologie (TSV, JSON, potentiellement OWL et RDF)
- Système de logging avancé pour le suivi et le débogage
- Interface utilisateur web responsive et intuitive

## Prérequis
- Go 1.16 ou supérieur
- Git

## Installation

### macOS et Linux

1. Ouvrez un terminal.

2. Installez Go si ce n'est pas déjà fait :
   ```
   # macOS avec Homebrew
   brew install go

   # Linux (Ubuntu/Debian)
   sudo apt-get update
   sudo apt-get install golang-go
   ```

3. Clonez le dépôt :
   ```
   git clone https://github.com/chrlesur/ontology-server.git
   cd ontology-server
   ```

4. Installez les dépendances Go :
   ```
   go mod tidy
   ```

5. Compilez le serveur :
   ```
   go build -o serveur-ontologie cmd/server/main.go
   ```

### Windows

1. Ouvrez PowerShell.

2. Installez Go si ce n'est pas déjà fait :
   - Téléchargez l'installateur depuis [le site officiel de Go](https://golang.org/dl/)
   - Exécutez l'installateur et suivez les instructions

3. Installez Git si nécessaire :
   - Téléchargez depuis [le site officiel de Git](https://git-scm.com/download/win)
   - Exécutez l'installateur

4. Clonez le dépôt :
   ```
   git clone https://github.com/chrlesur/ontology-server.git
   cd ontology-server
   ```

5. Installez les dépendances Go :
   ```
   go mod tidy
   ```

6. Compilez le serveur :
   ```
   go build -o serveur-ontologie.exe cmd/server/main.go
   ```

## Configuration

Modifiez `config.yaml` selon vos besoins :
```yaml
server:
  port: 8080
  host: localhost
  mode: debug  # Peut être "debug" ou "release"
logging:
  level: info
  directory: ./logs
storage:
  temp_directory: ./temp
```

## Utilisation

1. Démarrez le serveur :
   ```
   # macOS et Linux
   ./serveur-ontologie

   # Windows
   .\serveur-ontologie.exe
   ```

2. Accédez à l'interface web à l'adresse `http://localhost:8080`

3. Utilisez l'API RESTful :
   - GET `/api/v1/search` : Recherche dans les ontologies
   - GET `/api/v1/elements/{element_id}` : Détails d'un élément
   - POST `/api/v1/ontologies` : Ajout d'une ontologie
   - DELETE `/api/v1/ontologies/{ontology_id}` : Suppression d'une ontologie

## Développement

### Structure du projet
```
.
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── api/
│   ├── config/
│   ├── logger/
│   ├── models/
│   ├── parser/
│   ├── search/
│   └── storage/
├── web/
│   ├── index.html
│   ├── style.css
│   └── script.js
├── config.yaml
└── README.md
```

### Exécution des tests
```
go test ./...
```

## Dépannage

- Si vous rencontrez des problèmes de permission sur macOS ou Linux, assurez-vous que le fichier compilé est exécutable :
  ```
  chmod +x serveur-ontologie
  ```

- Sur Windows, si PowerShell bloque l'exécution de scripts, vous pouvez modifier la politique d'exécution :
  ```
  Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
  ```

- Vérifiez que les ports nécessaires (par défaut 8080) ne sont pas utilisés par d'autres applications.

## Licence
Ce projet est sous licence GNU General Public License v3.0 (GPL-3.0). Voir le fichier [LICENSE](LICENSE) pour plus de détails.

## Contact
Pour toute question ou suggestion, [christophe.lesur@cloud-temple.com].

