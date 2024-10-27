// ontologyLoader.js
import { uploadOntology } from './api.js';
import { showErrorMessage } from './main.js';

const uploadForm = document.getElementById('upload-form');
const uploadModal = document.getElementById('upload-modal');
const uploadProgress = document.getElementById('upload-progress');

// Initialisation du chargeur d'ontologie
export function initOntologyLoader() {
    uploadForm.addEventListener('submit', handleUpload);
    document.querySelector('.close').addEventListener('click', closeModal);
}

// Gestion du chargement de fichiers
async function handleUpload(event) {
    event.preventDefault();
    const formData = new FormData(uploadForm);
    
    try {
        uploadForm.classList.add('hidden');
        uploadProgress.classList.remove('hidden');
        
        const response = await uploadOntology(formData);
        console.log('Ontologie chargée avec succès:', response);
        
        // Émettre un événement personnalisé pour signaler le chargement réussi
        document.dispatchEvent(new CustomEvent('ontologyLoaded'));
        
        // Fermer le modal après un court délai
        setTimeout(() => {
            closeModal();
        }, 1000);
    } catch (error) {
        console.error('Erreur lors du chargement de l\'ontologie:', error);
        showErrorMessage('Une erreur est survenue lors du chargement de l\'ontologie.');
    } finally {
        uploadForm.classList.remove('hidden');
        uploadProgress.classList.add('hidden');
    }
}

function closeModal() {
    uploadModal.style.display = 'none';
    uploadForm.reset();
}

// Exporter la fonction closeModal pour l'utiliser ailleurs si nécessaire
export { closeModal };