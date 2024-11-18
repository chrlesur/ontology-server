// web/search.js

import { searchOntologies } from './api.js';
import { displayResults } from './results.js';
import { showErrorMessage } from './main.js';

const searchInput = document.getElementById('search-input');
const searchButton = document.getElementById('search-button');
const ontologySelect = document.getElementById('ontology-select');
const elementTypeSelect = document.getElementById('element-type-select');

// Initialisation de la recherche
export function initSearch() {
    if (!searchButton || !searchInput || !ontologySelect || !elementTypeSelect) {
        console.error('Éléments de recherche manquants');
        return;
    }

    // Gestionnaire de clic sur le bouton de recherche
    searchButton.addEventListener('click', handleSearch);

    // Gestionnaire de la touche Entrée dans le champ de recherche
    searchInput.addEventListener('keypress', (e) => {
        if (e.key === 'Enter') {
            e.preventDefault();
            handleSearch();
        }
    });

    // Recherche automatique lors du changement de filtre
    ontologySelect.addEventListener('change', debounce(handleSearch, 300));
    elementTypeSelect.addEventListener('change', debounce(handleSearch, 300));

    // Recherche automatique lors de la saisie
    searchInput.addEventListener('input', debounce(handleSearch, 300));
}

// Gestion de la recherche
export async function handleSearch(event) {
    if (event) event.preventDefault();

    const query = document.getElementById('search-input').value.trim();
    const fileId = document.getElementById('ontology-select').value;
    const elementType = document.getElementById('element-type-select').value;

    console.log("Search parameters:", { query, fileId, elementType });

    if (!query) {
        console.log("Search aborted: query is empty");
        showErrorMessage('Veuillez entrer un terme de recherche.');
        return;
    }

    const loadingSpinner = document.getElementById('loading-spinner');
    if (loadingSpinner) loadingSpinner.classList.remove('hidden');

    try {
        const results = await searchOntologies(query, fileId, elementType);
        console.log("Search results:", results);
        
        displayResults(results);

    } catch (error) {
        console.error('Erreur lors de la recherche:', error);
        showErrorMessage('Une erreur est survenue lors de la recherche: ' + error.message);
        displayResults([]);
    } finally {
        if (loadingSpinner) loadingSpinner.classList.add('hidden');
    }
}

// Fonction de recherche explicite pour être appelée depuis d'autres modules
export function performSearch(query) {
    if (searchInput) {
        searchInput.value = query;
        handleSearch();
    }
}

// Utilitaire pour le debounce
function debounce(func, wait) {
    let timeout;
    return function executedFunction(...args) {
        const later = () => {
            clearTimeout(timeout);
            func(...args);
        };
        clearTimeout(timeout);
        timeout = setTimeout(later, wait);
    };
}