

Tâche : Implémenter la structure de données en mémoire pour stocker et gérer les ontologies dans le fichier internal/storage/memory.go.

Contexte :
- Le projet est un serveur d'ontologie en Go.
- Les modèles de données sont définis dans internal/models/ontology.go et internal/models/element.go.
- Le système de logging est en place dans internal/logger/logger.go.

Spécifications :

1. Créer une structure MemoryStorage avec les champs suivants :
   - ontologies : une map[string]*models.Ontology pour stocker les ontologies
   - mutex : un sync.RWMutex pour la gestion de la concurrence

2. Implémenter une fonction NewMemoryStorage() qui initialise et retourne un nouveau MemoryStorage.

3. Implémenter les méthodes suivantes pour MemoryStorage :

   a. AddOntology(ontology *models.Ontology) error
      - Ajoute une nouvelle ontologie à la map
      - Utilise le mutex pour assurer la sécurité des opérations concurrentes
      - Retourne une erreur si l'ontologie existe déjà (basée sur l'ID)

   b. GetOntology(id string) (*models.Ontology, error)
      - Récupère une ontologie par son ID
      - Utilise le mutex pour la lecture
      - Retourne une erreur si l'ontologie n'est pas trouvée

   c. UpdateOntology(ontology *models.Ontology) error
      - Met à jour une ontologie existante
      - Utilise le mutex pour assurer la sécurité des opérations concurrentes
      - Retourne une erreur si l'ontologie n'existe pas

   d. DeleteOntology(id string) error
      - Supprime une ontologie par son ID
      - Utilise le mutex pour assurer la sécurité des opérations concurrentes
      - Retourne une erreur si l'ontologie n'existe pas

   e. ListOntologies() []*models.Ontology
      - Retourne une liste de toutes les ontologies stockées
      - Utilise le mutex pour la lecture

4. Utiliser le logger pour enregistrer les opérations importantes (par exemple, ajout ou suppression d'une ontologie).

5. Gérer les erreurs de manière appropriée et renvoyer des messages d'erreur clairs.

6. Assurer la thread-safety en utilisant correctement le mutex pour toutes les opérations.

7. Ajouter des commentaires de documentation Go pour chaque fonction exportée.

Contraintes :
- Suivre les conventions de nommage et de formatage standard de Go.
- Limiter chaque fonction à un maximum de 50 lignes de code.
- Utiliser des interfaces pour faciliter les tests futurs.
