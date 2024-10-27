# Serveur d'Ontologie

## Description
Ce projet est un serveur d'ontologie développé en Go, offrant une API RESTful et une interface utilisateur web pour la gestion, la recherche et la visualisation d'ontologies.
Il charge les ontologies crée avec le client **ontology**.

**VERSION DE RECHERCHE ET DEVELOPPEMENT**

## Fonctionnalités
- Chargement et gestion de fichiers d'ontologie (TSV, JSON, RDF, OWL)
- API RESTful pour l'interaction avec les ontologies
- Moteur de recherche avancé
- Interface utilisateur web responsive et performante
- Visualisation des relations entre éléments d'ontologie

## Prérequis
- Go 1.16 ou supérieur

## Installation

1. Clonez le dépôt :
   ```
   git clone https://github.com/chrlesur/ontology-server.git
   cd ontology-server
   ```

2. Installez les dépendances Go :
   ```
   go mod tidy
   ```

3. Compilez le serveur :
   ```
   go build -o serveur-ontologie cmd/server/main.go
   ```

4. (Optionnel) Pour le développement frontend :
   ```
   cd web
   npm install
   ```

## Configuration

1. Copiez le fichier de configuration exemple :
   ```
   cp config.example.yaml config.yaml
   ```

2. Modifiez `config.yaml` selon vos besoins :
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

## Utilisation

1. Démarrez le serveur :
   ```
   ./serveur-ontologie
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

## Licence
Ce projet est sous licence GNU General Public License v3.0 (GPL-3.0). Voir le fichier [LICENSE](LICENSE) pour plus de détails.

## Contact
Pour toute question ou suggestion, veuillez ouvrir une issue sur le dépôt GitHub du projet.
