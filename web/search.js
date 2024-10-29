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

    // Récupérer les valeurs de recherche
    const query = searchInput.value.trim();
    const sourceFile = ontologySelect.value;
    const elementType = elementTypeSelect.value;

    // Afficher le spinner de chargement
    const loadingSpinner = document.getElementById('loading-spinner');
    if (loadingSpinner) loadingSpinner.classList.remove('hidden');

    try {
        const results = await searchOntologies(query, sourceFile, elementType);
        console.log("Résultats de la recherche:", results);
        
        displayResults(results);

    } catch (error) {
        console.error('Erreur lors de la recherche:', error);
        showErrorMessage('Une erreur est survenue lors de la recherche.');
        // Afficher un résultat vide en cas d'erreur
        displayResults([]);
    } finally {
        // Cacher le spinner
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