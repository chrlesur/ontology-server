Tâche : Implémenter l'API RESTful pour le serveur d'ontologie dans les fichiers internal/api/handlers.go et internal/api/router.go.

Contexte :
- Le projet est un serveur d'ontologie en Go.
- La structure de données en mémoire (MemoryStorage) est implémentée dans internal/storage/memory.go.
- Les modèles de données sont définis dans internal/models/ontology.go et internal/models/element.go.
- Le système de logging est en place dans internal/logger/logger.go.

Spécifications :

1. Dans internal/api/handlers.go :

   a. Implémenter un handler pour la recherche (SearchHandler) :
      - Méthode HTTP : GET
      - Chemin : /api/v1/search
      - Paramètres de requête : q (terme de recherche), ontology_id (optionnel), type (optionnel), context_size (optionnel, défaut : 5)
      - Réponse : Liste des correspondances au format JSON

   b. Implémenter un handler pour les détails d'un élément (ElementDetailsHandler) :
      - Méthode HTTP : GET
      - Chemin : /api/v1/elements/{element_id}
      - Réponse : Détails complets de l'élément au format JSON

   c. Implémenter un handler pour l'ajout d'une ontologie (AddOntologyHandler) :
      - Méthode HTTP : POST
      - Chemin : /api/v1/ontologies
      - Corps de la requête : Données de l'ontologie au format JSON
      - Réponse : ID de l'ontologie créée

   d. Implémenter un handler pour la suppression d'une ontologie (DeleteOntologyHandler) :
      - Méthode HTTP : DELETE
      - Chemin : /api/v1/ontologies/{ontology_id}
      - Réponse : Statut de succès

2. Dans internal/api/router.go :

   a. Créer une fonction SetupRouter qui configure toutes les routes de l'API :
      - Utiliser le package http standard de Go ou un routeur tiers comme gorilla/mux
      - Définir les routes pour tous les handlers implémentés
      - Ajouter un middleware pour la gestion des erreurs et la journalisation des requêtes

3. Pour chaque handler :
   - Valider les entrées de l'utilisateur
   - Utiliser MemoryStorage pour accéder aux données
   - Gérer les erreurs et renvoyer des réponses HTTP appropriées
   - Utiliser le logger pour enregistrer les actions importantes

4. Implémenter une gestion des erreurs cohérente :
   - Créer des types d'erreur personnalisés si nécessaire
   - Renvoyer des réponses JSON structurées pour les erreurs

5. Ajouter des commentaires de documentation Go pour chaque fonction exportée.

6. Implémenter des tests unitaires pour chaque handler dans un fichier handlers_test.go.

Contraintes :
- Suivre les conventions de nommage et de formatage standard de Go.
- Limiter chaque fonction à un maximum de 50 lignes de code.
- Utiliser des constantes pour les codes de statut HTTP et les messages d'erreur communs.
