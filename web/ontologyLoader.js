// ontologyLoader.js
import { uploadOntology } from './api.js';

const uploadForm = document.getElementById('upload-form');

// Initialisation du chargeur d'ontologie
export function initOntologyLoader() {
    uploadForm.addEventListener('submit', handleUpload);
}

// Gestion du chargement de fichiers
async function handleUpload(event) {
    event.preventDefault();
    const formData = new FormData(uploadForm);
    try {
        const response = await uploadOntology(formData);
        console.log('Ontologie chargée avec succès:', response);
        // TODO: Mettre à jour l'interface utilisateur après le chargement réussi
    } catch (error) {
        console.error('Erreur lors du chargement de l\'ontologie:', error);
        showErrorMessage('Une erreur est survenue lors du chargement de l\'ontologie.');
    }
}