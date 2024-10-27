// ontologyLoader.js
import { uploadOntology } from './api.js';
import { showErrorMessage } from './main.js';

// Initialisation du chargeur d'ontologie
export function initOntologyLoader() {
    const uploadButton = document.getElementById('upload-button');
    const uploadModal = document.getElementById('upload-modal');
    const closeButton = uploadModal ? uploadModal.querySelector('.close') : null;
    const uploadForm = document.getElementById('upload-form');
    const uploadProgress = document.getElementById('upload-progress');

    if (!uploadButton || !uploadModal || !closeButton || !uploadForm || !uploadProgress) {
        console.error("One or more required elements are missing");
        return;
    }

    uploadButton.addEventListener('click', () => {
        uploadModal.style.display = 'block';
    });

    closeButton.addEventListener('click', () => {
        uploadModal.style.display = 'none';
    });

    uploadForm.addEventListener('submit', async (event) => {
        event.preventDefault();
        
        try {
            uploadForm.classList.add('hidden');
            uploadProgress.classList.remove('hidden');
            
            const formData = new FormData(uploadForm);
            const response = await uploadOntology(formData);
            console.log('Ontologie chargée avec succès:', response);
            
            // Fermer le modal après un court délai
            setTimeout(() => {
                uploadModal.style.display = 'none';
                // TODO: Mettre à jour l'interface utilisateur après le chargement réussi
            }, 1000);
        } catch (error) {
            console.error('Erreur lors du chargement de l\'ontologie:', error);
            showErrorMessage('Une erreur est survenue lors du chargement de l\'ontologie.');
        } finally {
            uploadForm.classList.remove('hidden');
            uploadProgress.classList.add('hidden');
        }
    });
}