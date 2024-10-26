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