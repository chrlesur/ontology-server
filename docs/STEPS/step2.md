En tant que développeur Go expérimenté, vous devez implémenter les parsers de fichiers pour un serveur Ontology. Ces parsers doivent traiter les formats TSV, RDF/OWL, et JSON. Voici les tâches spécifiques à réaliser :

1. Parser TSV (/internal/parser/tsv.go) :
   - Implémentez une fonction ParseTSV(filename string) ([]Element, error)
   - La fonction doit lire le fichier TSV ligne par ligne
   - Gérez les cas spéciaux : lignes vides, champs entre guillemets, caractères d'échappement
   - Pour chaque ligne valide, extrayez :
     - Nom de l'élément (colonne 1)
     - Type de l'élément (colonne 2)
     - Description de l'élément (colonne 3)
     - Positions (colonne 4) : parser la liste d'entiers séparés par des virgules
   - Retournez un slice de structures Element

2. Parser RDF/OWL (/internal/parser/rdf.go et /internal/parser/owl.go) :
   - Utilisez la bibliothèque github.com/knakk/rdf pour le parsing RDF/OWL
   - Implémentez une fonction ParseRDF(filename string) ([]Element, error)
   - Implémentez une fonction ParseOWL(filename string) ([]Element, error)
   - Extrayez les triplets (sujet, prédicat, objet)
   - Si disponible, associez les positions aux éléments
   - Convertissez la structure RDF/OWL en un slice de structures Element

3. Parser JSON (/internal/parser/json.go) :
   - Implémentez une fonction ParseJSON(filename string) (*JSONContext, error)
   - Utilisez encoding/json pour parser le fichier JSON
   - Extrayez :
     - La liste des éléments avec leurs propriétés
     - La liste des relations entre les éléments
     - Le contexte de chaque élément (texte avant et après)
   - Validez la structure du JSON selon le schéma suivant :
     ```go
     type JSONContext struct {
         Elements  map[string]JSONElement
         Relations []JSONRelation
     }

     type JSONElement struct {
         Name        string
         Type        string
         Description string
         Context     string
     }

     type JSONRelation struct {
         Source string
         Type   string
         Target string
     }
     ```

4. Gestion des erreurs :
   - Implémentez des erreurs personnalisées pour chaque type de problème de parsing
   - Assurez-vous que chaque fonction de parsing retourne des erreurs appropriées

5. Tests unitaires :
   - Pour chaque parser, créez un fichier de test correspondant (tsv_test.go, rdf_test.go, owl_test.go, json_test.go)
   - Écrivez des tests unitaires couvrant les cas normaux et les cas d'erreur
   - Incluez des tests de table pour couvrir plusieurs scénarios
   - Visez une couverture de code d'au moins 80% pour chaque parser

Assurez-vous que le code est bien commenté, suit les conventions de Go, et gère correctement les erreurs. Fournissez le code pour chaque fichier mentionné ci-dessus, ainsi que les tests unitaires correspondants.