
Tâche : Développer une interface utilisateur web simple en HTML/JavaScript pour interagir avec l'API du serveur d'ontologie.

Contexte :
- Le projet est un serveur d'ontologie en Go avec une API RESTful fonctionnelle.
- Le moteur de recherche est implémenté et accessible via l'API.
- L'interface doit être simple et fonctionnelle, sans nécessiter de framework JavaScript complexe.

Spécifications :

1. Créer les fichiers suivants dans le répertoire web/ :
   - index.html : la page principale de l'application
   - style.css : les styles CSS pour l'interface
   - script.js : le code JavaScript pour la logique côté client

2. Dans index.html :
   a. Créer une structure HTML5 de base
   b. Inclure une section pour la recherche avec :
      - Un champ de recherche principal
      - Des options de filtrage (par ontologie, par type d'élément)
      - Un bouton de recherche
   c. Inclure une section pour afficher les résultats de recherche
   d. Inclure une section pour afficher les détails d'un élément sélectionné
   e. Lier les fichiers style.css et script.js

3. Dans style.css :
   a. Définir des styles de base pour rendre l'interface utilisable et responsive
   b. Utiliser une approche mobile-first pour le design
   c. Définir des styles pour les éléments de formulaire, les résultats de recherche et les détails des éléments

4. Dans script.js :
   a. Implémenter une fonction pour envoyer une requête de recherche à l'API
   b. Implémenter une fonction pour afficher les résultats de recherche
   c. Implémenter une fonction pour récupérer et afficher les détails d'un élément
   d. Ajouter des gestionnaires d'événements pour :
      - La soumission du formulaire de recherche
      - Le clic sur un résultat de recherche pour voir plus de détails
   e. Implémenter une gestion basique des erreurs et afficher des messages à l'utilisateur

5. Fonctionnalités à implémenter :
   a. Recherche en temps réel (mise à jour des résultats pendant la saisie)
   b. Pagination des résultats de recherche
   c. Mise en évidence des termes de recherche dans les résultats
   d. Possibilité de filtrer les résultats par ontologie ou type d'élément

6. Assurer la compatibilité avec les navigateurs modernes (Chrome, Firefox, Safari, Edge)

7. Implémenter des mesures de sécurité de base :
   a. Échapper les entrées utilisateur pour prévenir les attaques XSS
   b. Utiliser HTTPS pour les requêtes API (à configurer côté serveur)

8. Optimiser les performances :
   a. Minimiser le nombre de requêtes API
   b. Utiliser le cache du navigateur quand c'est approprié

9. Ajouter des commentaires explicatifs dans le code JavaScript

Contraintes :
- Utiliser du JavaScript vanilla (pas de frameworks comme React ou Vue.js)
- Assurer que l'interface est accessible et utilisable sur les appareils mobiles et de bureau
- Maintenir une séparation claire entre la structure (HTML), le style (CSS) et la logique (JavaScript)
- Limiter la taille des fichiers pour des temps de chargement rapides

