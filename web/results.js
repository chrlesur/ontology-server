// web/results.js

import { 
    getElementDetails, 
    getElementRelations,
    loadOntologies 
} from './api.js';
import { createRelationsGraph } from './graph.js';
import { showErrorMessage } from './main.js';
import { performSearch } from './search.js'; // Ajout de l'import manquant

export function displayResults(results) {
    console.log("Received results:", results);
    const resultsList = document.getElementById('results-list');
    resultsList.innerHTML = '';

    if (!results || results.length === 0) {
        resultsList.innerHTML = '<div class="empty-state">Aucun résultat trouvé</div>';
        return;
    }
    
    results.forEach(result => {
        const resultItem = document.createElement('div');
        resultItem.className = 'result-item';
        
        // Extraire les informations spécifiques au fichier
        const metadata = result.SourceMetadata;
        let fileID = result.FileID;  // Ceci devrait être le FileID spécifique à ce résultat
        let fileInfo = null;

        console.log("Processing result:", result);

        if (metadata && metadata.files) {
            if (fileID && metadata.files[fileID]) {
                fileInfo = metadata.files[fileID];
            } else {
                console.warn("FileID not found in metadata, using first available file");
                const fileIds = Object.keys(metadata.files);
                if (fileIds.length > 0) {
                    fileID = fileIds[0];
                    fileInfo = metadata.files[fileID];
                }
            }
        }

        console.log("Using FileID:", fileID, "FileInfo:", fileInfo);

        // Stocker les informations dans dataset
        resultItem.dataset.fileId = fileID || '';
        resultItem.dataset.sourceFile = fileInfo ? fileInfo.source_file : '';
        resultItem.dataset.directory = fileInfo ? fileInfo.directory : '';
        resultItem.dataset.fileDate = fileInfo ? fileInfo.file_date : '';
        resultItem.dataset.sha256Hash = fileInfo ? fileInfo.sha256_hash : '';
        resultItem.dataset.ontologyFile = metadata ? metadata.ontology_file : '';
        resultItem.dataset.processingDate = metadata ? metadata.processing_date : '';
        
        resultItem.innerHTML = `
            <div class="result-item-content">
                <h3>${escapeHtml(result.ElementName)}</h3>
                <p>${escapeHtml(result.Description || '')}</p>
            </div>
            <div class="result-item-meta">
                <div class="file-name">${escapeHtml(fileInfo ? fileInfo.source_file : 'Fichier inconnu')} (${fileID || 'ID inconnu'})</div>
            </div>
        `;

        // Gestion de la fenêtre flottante (tooltip)
        resultItem.addEventListener('mouseenter', (e) => {
            showMetadataTooltip(e.currentTarget);
        });

        resultItem.addEventListener('mousemove', (e) => {
            updateTooltipPosition(e);
        });

        resultItem.addEventListener('mouseleave', () => {
            hideMetadataTooltip();
        });

        // Gestion de la sélection de l'élément
        resultItem.addEventListener('click', () => {
            document.querySelectorAll('.result-item').forEach(item => 
                item.classList.remove('selected'));
            resultItem.classList.add('selected');
            showElementDetails(result.ElementName);
        });

        resultsList.appendChild(resultItem);
    });
}

function escapeHtml(unsafe) {
    if (!unsafe) return '';
    return unsafe
        .replace(/&/g, "&amp;")
        .replace(/</g, "&lt;")
        .replace(/>/g, "&gt;")
        .replace(/"/g, "&quot;")
        .replace(/'/g, "&#039;");
}

function openSourceFile(metadata, fileID) {
    if (!metadata || !fileID) {
        console.error('Métadonnées ou FileID manquants pour ouvrir le fichier');
        return;
    }

    const fileInfo = metadata.files[fileID];
    if (!fileInfo) {
        console.error('Informations du fichier non trouvées');
        return;
    }

    // Construisez l'URL pour ouvrir le fichier en visualisation
    const viewerUrl = `/api/view-source?path=${encodeURIComponent(fileInfo.directory + '/' + fileInfo.source_file)}`;
    
    // Ouvrez le fichier dans un nouvel onglet ou une nouvelle fenêtre
    window.open(viewerUrl, '_blank');
}

async function showElementDetails(elementName) {
    const loadingSpinner = document.getElementById('loading-spinner');
    if (loadingSpinner) loadingSpinner.classList.remove('hidden');

    try {
        // Récupérer les détails de l'élément
        const element = await getElementDetails(elementName);
        console.log("Détails de l'élément reçus:", element);

        // Afficher les détails
        const detailsContainer = document.getElementById('element-details');
        if (detailsContainer) {
            detailsContainer.innerHTML = `
                <div class="element-info">
                    <h3>${escapeHtml(element.Name)}</h3>
                    <p><strong>Type:</strong> ${escapeHtml(element.Type || '')}</p>
                    <p><strong>Description:</strong> ${escapeHtml(element.Description || '')}</p>
                </div>
            `;
        }

        // Afficher les contextes
        displayElementContexts(element);
        
        // Gérer les relations
        const relations = await getElementRelations(elementName);
        if (relations && relations.length > 0) {
            createRelationsGraph(element, relations);
            displayRelationsList(relations);
        } else {
            document.getElementById('element-relations-graph').innerHTML = 
                '<div class="empty-state">Aucune relation à afficher</div>';
            document.getElementById('element-relations-list').innerHTML = 
                '<div class="empty-state">Aucune relation disponible</div>';
        }

    } catch (error) {
        console.error('Erreur lors de la récupération des détails:', error);
        showErrorMessage('Impossible de charger les détails de l\'élément.');
    } finally {
        if (loadingSpinner) loadingSpinner.classList.add('hidden');
    }
}

// Rendre les éléments marqués cliquables
function displayElementContexts(element) {
    const contextsContainer = document.getElementById('element-contexts');
    if (!contextsContainer) return;

    contextsContainer.innerHTML = '<h3>Contextes</h3>';

    if (!element.Contexts || element.Contexts.length === 0) {
        contextsContainer.innerHTML += '<div class="empty-state">Aucun contexte disponible</div>';
        return;
    }

    element.Contexts.forEach((ctx, index) => {
        const contextDiv = document.createElement('div');
        contextDiv.className = 'context';

        const highlightedBefore = highlightElement(ctx.before?.join(' ') || '', element.Name);
        const highlightedAfter = highlightElement(ctx.after?.join(' ') || '', element.Name);

        contextDiv.innerHTML = `
        <div class="context-content">
            <span class="before">${highlightElement(ctx.before?.join(' ') || '', element.Name)}</span>
            <mark class="element">${escapeHtml(ctx.element)}</mark>
            <span class="after">${highlightElement(ctx.after?.join(' ') || '', element.Name)}</span>
        </div>
        <div class="context-meta">
            <span class="file-id">FileID: ${ctx.file_id} - </span>
            <span class="file-position">Position: ${ctx.file_position}</span>
        </div>
    `;
        // Ajouter les gestionnaires d'événements pour les éléments marqués
        contextDiv.querySelectorAll('mark').forEach(mark => {
            mark.style.cursor = 'pointer';
            mark.addEventListener('click', () => {
                performSearch(mark.textContent);
            });
        });

        contextsContainer.appendChild(contextDiv);
    });
}

function displayRelationsList(relations) {
    const listContainer = document.getElementById('element-relations-list');
    if (!listContainer) return;

    listContainer.innerHTML = '';

    if (!relations || relations.length === 0) {
        listContainer.innerHTML = '<div class="empty-state">Aucune relation disponible</div>';
        return;
    }

    const ul = document.createElement('ul');
    ul.className = 'relations-list';

    relations.forEach(relation => {
        const li = document.createElement('li');
        li.className = 'relation-item';
        li.innerHTML = `
            <span class="relation-source">${escapeHtml(relation.Source)}</span>
            <span class="relation-type">${escapeHtml(relation.Type)}</span>
            <span class="relation-target">${escapeHtml(relation.Target)}</span>
        `;
        ul.appendChild(li);
    });

    listContainer.appendChild(ul);
}

// Fonction pour afficher le tooltip
function showMetadataTooltip(element) {
    const tooltip = document.getElementById('metadata-tooltip');
    if (!tooltip) {
        console.error("Tooltip element not found in the DOM");
        return;
    }
    if (!element) {
        console.error("Source element is undefined");
        return;
    }

    console.log("Element dataset:", element.dataset);

    const fileID = element.dataset.fileId;
    const sourceFile = element.dataset.sourceFile;
    const directory = element.dataset.directory;
    const fileDate = element.dataset.fileDate;
    const sha256Hash = element.dataset.sha256Hash;
    const ontologyFile = element.dataset.ontologyFile;
    const processingDate = element.dataset.processingDate;

    console.log("Showing tooltip for file:", fileID, sourceFile);

    const formattedDate = fileDate ? new Date(fileDate).toLocaleString() : 'Date non disponible';
    const formattedProcessingDate = processingDate ? new Date(processingDate).toLocaleString() : 'Date non disponible';

    tooltip.innerHTML = `
        <div class="tooltip-content">
            <h4>Informations du document</h4>
            <p><strong>Fichier ontologie :</strong> ${escapeHtml(ontologyFile || 'Non spécifié')}</p>
            <p><strong>Date de traitement :</strong> ${formattedProcessingDate}</p>
            <p><strong>File ID :</strong> ${escapeHtml(fileID || 'Non spécifié')}</p>
            <p><strong>Fichier source :</strong> ${escapeHtml(sourceFile || 'Non spécifié')}</p>
            <p><strong>Répertoire :</strong> ${escapeHtml(directory || 'Non spécifié')}</p>
            <p><strong>Date du fichier :</strong> ${formattedDate}</p>
            <p><strong>SHA256 :</strong> <span class="hash">${sha256Hash || 'Non disponible'}</span></p>
        </div>
    `;

    tooltip.style.display = 'block';
}

// Fonction pour mettre à jour la position du tooltip
function updateTooltipPosition(event) {
    const tooltip = document.getElementById('metadata-tooltip');
    if (!tooltip) return;

    const margin = 10;
    const tooltipRect = tooltip.getBoundingClientRect();
    const windowWidth = window.innerWidth;
    const windowHeight = window.innerHeight;

    let left = event.clientX + margin;
    let top = event.clientY + margin;

    if (left + tooltipRect.width > windowWidth) {
        left = windowWidth - tooltipRect.width - margin;
    }

    if (top + tooltipRect.height > windowHeight) {
        top = windowHeight - tooltipRect.height - margin;
    }

    tooltip.style.left = `${left}px`;
    tooltip.style.top = `${top}px`;
}

// Fonction pour cacher le tooltip
function hideMetadataTooltip() {
    const tooltip = document.getElementById('metadata-tooltip');
    if (tooltip) {
        tooltip.style.display = 'none';
    }
}

function highlightElement(text, elementName) {
    if (!text) return '';
    return text.replace(new RegExp(elementName, 'gi'), match => `<mark>${match}</mark>`);
}

// Export des fonctions nécessaires
export { showElementDetails };