
Tâche : Implémenter le moteur de recherche pour les ontologies dans le fichier internal/search/engine.go.

Contexte :
- Le projet est un serveur d'ontologie en Go.
- La structure de données en mémoire (MemoryStorage) est implémentée dans internal/storage/memory.go.
- L'API RESTful est implémentée dans internal/api/handlers.go et internal/api/router.go.
- Les modèles de données sont définis dans internal/models/ontology.go et internal/models/element.go.
- Le système de logging est en place dans internal/logger/logger.go.

Spécifications :

1. Créer une structure SearchEngine avec les champs nécessaires, notamment :
   - storage : un pointeur vers MemoryStorage pour accéder aux données
   - logger : un pointeur vers le logger pour la journalisation

2. Implémenter une fonction NewSearchEngine qui initialise et retourne un nouveau SearchEngine.

3. Implémenter la méthode principale de recherche :
   func (se *SearchEngine) Search(query string, ontologyID string, elementType string, contextSize int) ([]SearchResult, error)
   
   Cette méthode doit :
   - Accepter les paramètres de recherche (query, ontologyID optionnel, elementType optionnel, contextSize)
   - Effectuer une recherche dans les ontologies stockées
   - Retourner une liste de résultats de recherche et une erreur éventuelle

4. Implémenter des fonctions auxiliaires pour la recherche :
   a. Une fonction pour la correspondance exacte
   b. Une fonction pour la correspondance partielle
   c. Une fonction pour la recherche insensible à la casse
   d. Une fonction pour extraire le contexte d'un élément trouvé

5. Créer une structure SearchResult pour représenter les résultats de recherche, contenant :
   - OntologyID : l'ID de l'ontologie
   - ElementName : le nom de l'élément trouvé
   - ElementType : le type de l'élément
   - Description : la description de l'élément
   - Context : le contexte de l'élément (mots avant et après)
   - Position : la position de l'élément dans le fichier source

6. Implémenter une logique de scoring simple pour classer les résultats de recherche par pertinence.

7. Optimiser les performances de recherche :
   - Utiliser des index en mémoire si nécessaire
   - Implémenter une recherche parallèle si possible

8. Gérer les erreurs de manière appropriée et renvoyer des messages d'erreur clairs.

9. Utiliser le logger pour enregistrer les opérations importantes (par exemple, début et fin de recherche, erreurs).

10. Ajouter des commentaires de documentation Go pour chaque fonction exportée.

11. Implémenter des tests unitaires pour le moteur de recherche dans un fichier engine_test.go.

Contraintes :
- Suivre les conventions de nommage et de formatage standard de Go.
- Limiter chaque fonction à un maximum de 50 lignes de code.
- Assurer la thread-safety si nécessaire.
- Optimiser pour la performance, en particulier pour les grandes ontologies.

