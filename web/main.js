// main.js
import { initUI } from './ui.js';
import { initSearch, handleSearch } from './search.js';
import { initOntologyLoader } from './ontologyLoader.js';
import './results.js';

// Fonction utilitaire pour afficher les messages d'erreur
export function showErrorMessage(message) {
    alert(message);
}

// Initialisation de l'application
document.addEventListener('DOMContentLoaded', async () => {
    try {
        await initUI();
        initSearch();
        initOntologyLoader();

        // Ajouter un écouteur pour l'événement personnalisé
        document.addEventListener('performSearch', (event) => {
            const searchQuery = event.detail.query;
            document.getElementById('search-input').value = searchQuery;
            handleSearch();
        });

    } catch (error) {
        console.error('Erreur lors de l\'initialisation:', error);
        showErrorMessage('Une erreur est survenue lors de l\'initialisation de l\'application.');
    }
});