Voici le markdown complet mis à jour pour décrire la structure des fichiers de sortie .tsv et .json :

# Description des fichiers de sortie de l'ontologie

La création de cette ontologie génère deux types de fichiers de sortie : un fichier `.tsv` (Tab-Separated Values) et un fichier `.json` (JavaScript Object Notation). Ces fichiers représentent différents aspects de l'ontologie extraite du texte source.

## 1. Fichier .tsv

Le fichier `.tsv` contient les éléments de base de l'ontologie dans un format tabulaire. Chaque ligne représente un élément distinct et contient les informations suivantes, séparées par des tabulations :

- Nom de l'élément
- Type de l'élément (Concept, Rôle, Organisation, etc.)
- Description de l'élément
- Liste des positions dans le texte où l'élément apparaît

Exemple :
```
Qualification_Juridique    Concept    Processus de détermination du statut légal d'une situation ou d'un individu    1
Agent_Service_Public    Concept    Personne travaillant pour un service public    23,47,189,24,48,62,81,117,120,154,177,190,208,244,443
```

Ce format permet une visualisation rapide des éléments de l'ontologie et de leur distribution dans le texte.

## 2. Fichier .json

Le fichier `.json` offre une représentation plus détaillée et structurée de l'ontologie. Chaque objet dans le tableau JSON correspond à une occurrence d'un élément de l'ontologie dans le texte. Voici la structure de chaque objet :

- `position` : L'index du mot dans le texte où l'élément commence.
- `before` : Un tableau de mots qui précèdent immédiatement l'élément dans le texte.
- `after` : Un tableau de mots qui suivent immédiatement l'élément dans le texte.
- `element` : Le nom de l'élément de l'ontologie identifié à cette position.
- `length` : Le nombre de mots que l'élément occupe dans le texte.

Exemple :
```json
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
```

Cette structure JSON permet une analyse contextuelle détaillée de chaque occurrence des éléments de l'ontologie dans le texte, facilitant diverses tâches de traitement du langage naturel et d'analyse de texte.