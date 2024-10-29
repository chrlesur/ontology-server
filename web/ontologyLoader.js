// web/ontologyLoader.js
import { uploadOntology } from './api.js';
import { showErrorMessage } from './main.js';

// Initialisation du chargeur d'ontologie
export function initOntologyLoader() {
    const uploadButton = document.getElementById('upload-button');
    const uploadModal = document.getElementById('upload-modal');
    const closeButton = uploadModal?.querySelector('.close');
    const uploadForm = document.getElementById('upload-form');
    const uploadProgress = document.getElementById('upload-progress');
    const errorMessage = document.getElementById('error-message');
    const successMessage = document.getElementById('success-message');

    // Vérifier que tous les éléments nécessaires sont présents
    if (!uploadButton || !uploadModal || !closeButton || !uploadForm || !uploadProgress) {
        console.error("Un ou plusieurs éléments requis sont manquants");
        return;
    }

    // Fonction pour réinitialiser les messages
    const resetMessages = () => {
        if (errorMessage) errorMessage.style.display = 'none';
        if (successMessage) successMessage.style.display = 'none';
    };

    // Ouvrir le modal
    uploadButton.addEventListener('click', () => {
        uploadModal.style.display = 'block';
        resetMessages();
    });

    // Fermer le modal
    closeButton.addEventListener('click', () => {
        uploadModal.style.display = 'none';
        resetMessages();
        uploadForm.reset();
    });

    // Fermer le modal en cliquant à l'extérieur
    window.addEventListener('click', (event) => {
        if (event.target === uploadModal) {
            uploadModal.style.display = 'none';
            resetMessages();
            uploadForm.reset();
        }
    });

    // Gérer la soumission du formulaire
    uploadForm.addEventListener('submit', async (event) => {
        event.preventDefault();
        resetMessages();
        
        try {
            // Récupérer les fichiers
            const ontologyFile = document.getElementById('ontology-file')?.files[0];
            const metadataFile = document.getElementById('metadata-file')?.files[0];

            // Vérifier que les fichiers requis sont présents
            if (!ontologyFile || !metadataFile) {
                throw new Error('Les fichiers d\'ontologie et de métadonnées sont obligatoires');
            }

            // Afficher le spinner de chargement
            uploadForm.classList.add('hidden');
            uploadProgress.classList.remove('hidden');
            
            // Préparer les données
            const formData = new FormData();
            formData.append('ontologyFile', ontologyFile);
            formData.append('metadataFile', metadataFile);
            
            // Ajouter le fichier de contexte s'il est présent
            const contextFile = document.getElementById('context-file')?.files[0];
            if (contextFile) {
                formData.append('contextFile', contextFile);
            }

            // Envoyer les fichiers
            const response = await uploadOntology(formData);
            
            // Afficher le message de succès
            if (successMessage) {
                successMessage.textContent = 'Ontologie chargée avec succès';
                successMessage.style.display = 'block';
            }

            // Réinitialiser et fermer après un délai
            setTimeout(() => {
                uploadModal.style.display = 'none';
                uploadForm.reset();
                resetMessages();
                // Déclencher un événement pour mettre à jour l'interface
                document.dispatchEvent(new CustomEvent('ontologyLoaded'));
            }, 2000);

        } catch (error) {
            console.error('Erreur lors du chargement:', error);
            if (errorMessage) {
                errorMessage.textContent = error.message || 'Une erreur est survenue lors du chargement';
                errorMessage.style.display = 'block';
            }
        } finally {
            // Toujours réafficher le formulaire et cacher le spinner
            uploadForm.classList.remove('hidden');
            uploadProgress.classList.add('hidden');
        }
    });
}

// Fonction pour fermer le modal de manière programmée
export function closeModal() {
    const uploadModal = document.getElementById('upload-modal');
    const uploadForm = document.getElementById('upload-form');
    const errorMessage = document.getElementById('error-message');
    const successMessage = document.getElementById('success-message');

    if (uploadModal) {
        uploadModal.style.display = 'none';
        if (uploadForm) uploadForm.reset();
        if (errorMessage) errorMessage.style.display = 'none';
        if (successMessage) successMessage.style.display = 'none';
    }
}