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
    const resultsList = document.getElementById('results-list');
    resultsList.innerHTML = '';

    if (!results || results.length === 0) {
        resultsList.innerHTML = '<div class="empty-state">Aucun résultat trouvé</div>';
        return;
    }
    
    results.forEach(result => {
        const resultItem = document.createElement('div');
        resultItem.className = 'result-item';
        
        resultItem.innerHTML = `
            <div class="result-item-content">
                <h3>${escapeHtml(result.ElementName)}</h3>
                <p>${escapeHtml(result.Description || '')}</p>
            </div>
            <div class="result-item-meta">
                <div class="file-name">${escapeHtml(result.sourceFile)}</div>
            </div>
        `;

        // Ajoutez un gestionnaire d'événements pour ouvrir le fichier
        const fileNameElement = resultItem.querySelector('.file-name');
        fileNameElement.style.cursor = 'pointer';
        fileNameElement.addEventListener('click', () => openSourceFile(result.sourceMetadata));
        
        // Gestion du tooltip
        resultItem.addEventListener('mouseenter', (e) => {
            if (result.Source) {
                showMetadataTooltip(e, result.Source);
            }
        });

        resultItem.addEventListener('mousemove', (e) => {
            updateTooltipPosition(e);
        });

        resultItem.addEventListener('mouseleave', () => {
            hideMetadataTooltip();
        });

        resultItem.addEventListener('click', () => {
            document.querySelectorAll('.result-item').forEach(item => 
                item.classList.remove('selected'));
            resultItem.classList.add('selected');
            showElementDetails(result.ElementName);
        });

        resultsList.appendChild(resultItem);
    });
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
                    ${element.Positions && element.Positions.length > 0 ? 
                        `<p><strong>Positions:</strong> ${element.Positions.join(', ')}</p>` : ''}
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
                <span class="before">${highlightedBefore}</span>
                <mark class="element">${escapeHtml(ctx.element)}</mark>
                <span class="after">${highlightedAfter}</span>
            </div>
            <div class="context-position">Position: ${ctx.position}</div>
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

function escapeHtml(unsafe) {
    if (!unsafe) return '';
    return unsafe
        .replace(/&/g, "&amp;")
        .replace(/</g, "&lt;")
        .replace(/>/g, "&gt;")
        .replace(/"/g, "&quot;")
        .replace(/'/g, "&#039;");
}

// Fonction pour afficher le tooltip
function showMetadataTooltip(event, metadata) {
    const tooltip = document.getElementById('metadata-tooltip');
    if (!tooltip) return;

    const formattedDate = new Date(metadata.file_date).toLocaleString();
    const formattedProcessingDate = new Date(metadata.processing_date).toLocaleString();

    tooltip.innerHTML = `
        <div class="tooltip-content">
            <h4>Informations du document</h4>
            <p><strong>Fichier source :</strong> ${escapeHtml(metadata.source_file)}</p>
            <p><strong>Répertoire :</strong> ${escapeHtml(metadata.directory)}</p>
            <p><strong>Date du fichier :</strong> ${formattedDate}</p>
            <p><strong>Date de traitement :</strong> ${formattedProcessingDate}</p>
            <p><strong>SHA256 :</strong> <span class="hash">${metadata.sha256_hash}</span></p>
            <p><strong>Fichier ontologie :</strong> ${escapeHtml(metadata.ontology_file)}</p>
            ${metadata.context_file ? 
                `<p><strong>Fichier contexte :</strong> ${escapeHtml(metadata.context_file)}</p>` : 
                ''}
        </div>
    `;

    tooltip.style.display = 'block';
    updateTooltipPosition(event);
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

    // Ajustement horizontal
    if (left + tooltipRect.width > windowWidth) {
        left = windowWidth - tooltipRect.width - margin;
    }

    // Ajustement vertical
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

function openSourceFile(metadata) {
    if (!metadata) {
        console.error('Métadonnées manquantes pour ouvrir le fichier');
        return;
    }

    // Construisez l'URL pour ouvrir le fichier en visualisation
    const viewerUrl = `/api/view-source?path=${encodeURIComponent(metadata.directory + '/' + metadata.source_file)}`;
    
    // Ouvrez le fichier dans un nouvel onglet ou une nouvelle fenêtre
    window.open(viewerUrl, '_blank');
}