// search.js
import { searchOntologies } from './api.js';
import { displayResults } from './results.js';
import { showErrorMessage } from './main.js';

const searchInput = document.getElementById('search-input');
const searchButton = document.getElementById('search-button');
const ontologySelect = document.getElementById('ontology-select');
const elementTypeSelect = document.getElementById('element-type-select');

// Initialisation de la recherche
export function initSearch() {
    searchButton.addEventListener('click', handleSearch);
    searchInput.addEventListener('input', debounce(handleSearch, 300));
}

// Gestion de la recherche
async function handleSearch(event) {
    if (event) event.preventDefault();
    const query = searchInput.value;
    const ontologyId = ontologySelect.value;
    const elementType = elementTypeSelect.value;

    try {
        const results = await searchOntologies(query, ontologyId, elementType);
        displayResults(results);
    } catch (error) {
        console.error('Erreur lors de la recherche:', error);
        showErrorMessage('Une erreur est survenue lors de la recherche.');
    }
}

// Fonction utilitaire pour le debounce
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