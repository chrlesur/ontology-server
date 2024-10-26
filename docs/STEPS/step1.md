En tant que développeur Go expérimenté, vous devez configurer un nouveau projet pour un serveur Ontology avec API et interface web. Voici les tâches spécifiques à réaliser :

1. Initialisation du projet :
   - Créez un nouveau répertoire nommé 'ontology-server'
   - Initialisez un nouveau module Go avec 'go mod init github.com/votre-nom/ontology-server'

2. Structure du projet :
   Créez la structure de répertoires suivante :
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

3. Configuration de base :
   - Dans /internal/config/config.go, implémentez une structure Config qui peut lire les paramètres suivants depuis un fichier YAML :
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
   - Utilisez la bibliothèque 'gopkg.in/yaml.v2' pour parser le YAML
   - Implémentez une fonction LoadConfig(filename string) (*Config, error) qui charge la configuration depuis un fichier

4. Système de logging :
   - Dans /internal/logger/logger.go, implémentez un logger personnalisé avec les niveaux : Debug, Info, Warning, Error
   - Le logger doit supporter la rotation quotidienne des fichiers de log
   - Utilisez le format de log : [LEVEL] [TIMESTAMP] [FILE:LINE] Message

5. Fichier principal :
   - Dans /cmd/server/main.go, créez la structure de base du programme principal
   - Chargez la configuration
   - Initialisez le logger
   - Préparez le code pour démarrer le serveur HTTP (sans l'implémenter complètement pour le moment)

Assurez-vous que le code est bien commenté, suit les conventions de Go, et gère correctement les erreurs. Fournissez le code pour chaque fichier mentionné ci-dessus.
