// main.js
import { initUI } from './ui.js';
import { initSearch } from './search.js';
import { initOntologyLoader } from './ontologyLoader.js';

// Fonction utilitaire pour afficher les messages d'erreur
export function showErrorMessage(message) {
    // TODO: Implémenter une meilleure UI pour les messages d'erreur
    alert(message);
}

// Initialisation de l'application
document.addEventListener('DOMContentLoaded', async () => {
    try {
        await initUI();
        initSearch();
        initOntologyLoader();
    } catch (error) {
        console.error('Erreur lors de l\'initialisation:', error);
        showErrorMessage('Une erreur est survenue lors de l\'initialisation de l\'application.');
    }
});