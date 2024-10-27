// ui.js
import { loadOntologies, loadElementTypes } from './api.js';

// Éléments du DOM
const ontologySelect = document.getElementById('ontology-select');
const elementTypeSelect = document.getElementById('element-type-select');
const uploadButton = document.getElementById('upload-button');
const uploadModal = document.getElementById('upload-modal');
const closeModal = document.querySelector('.close');

// Initialisation de l'interface utilisateur
export async function initUI() {
    await populateOntologySelect();
    await populateElementTypeSelect();
    await updateOntologySelect();
    setupModalListeners();

    document.addEventListener('ontologyLoaded', updateOntologySelect);

}

// Remplir le sélecteur d'ontologies
async function populateOntologySelect() {
    try {
        const ontologies = await loadOntologies();
        ontologySelect.innerHTML = '<option value="">Toutes les ontologies</option>';
        ontologies.forEach(ontology => {
            const option = document.createElement('option');
            option.value = ontology.id;
            option.textContent = ontology.name;
            ontologySelect.appendChild(option);
        });
    } catch (error) {
        console.error('Erreur lors du chargement des ontologies:', error);
        throw error;
    }
}

// Remplir le sélecteur de types d'éléments
async function populateElementTypeSelect() {
    try {
        const types = await loadElementTypes();
        elementTypeSelect.innerHTML = '<option value="">Tous les types</option>';
        types.forEach(type => {
            const option = document.createElement('option');
            option.value = type;
            option.textContent = type;
            elementTypeSelect.appendChild(option);
        });
    } catch (error) {
        console.error('Erreur lors du chargement des types d\'éléments:', error);
        throw error;
    }
}

// Configuration des écouteurs d'événements pour le modal
function setupModalListeners() {
    uploadButton.addEventListener('click', () => uploadModal.style.display = 'block');
    closeModal.addEventListener('click', () => uploadModal.style.display = 'none');
    window.addEventListener('click', (event) => {
        if (event.target === uploadModal) {
            uploadModal.style.display = 'none';
        }
    });
}

export async function updateOntologySelect() {
    try {
        const ontologies = await loadOntologies();
        ontologySelect.innerHTML = '<option value="">Toutes les ontologies</option>';
        ontologies.forEach(ontology => {
            const option = document.createElement('option');
            option.value = ontology.id;
            option.textContent = ontology.name;
            ontologySelect.appendChild(option);
        });
    } catch (error) {
        console.error('Erreur lors de la mise à jour des ontologies:', error);
        showErrorMessage('Erreur lors de la mise à jour de la liste des ontologies.');
    }
}