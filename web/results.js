// results.js
import { getElementDetails } from './api.js';
import { createRelationsGraph } from './graph.js';
import { showErrorMessage } from './main.js';

const resultsList = document.getElementById('results-list');
const elementDetails = document.getElementById('element-details');
const elementContexts = document.getElementById('element-contexts');
const elementRelations = document.getElementById('element-relations');

export function displayResults(results) {
    resultsList.innerHTML = '';
    results.forEach(result => {
        const resultItem = document.createElement('div');
        resultItem.className = 'result-item';
        resultItem.innerHTML = `
            <h3>${escapeHtml(result.ElementName)}</h3>
            <p>${escapeHtml(result.Description)}</p>
        `;
        resultItem.addEventListener('click', () => showElementDetails(result.ElementName));
        resultsList.appendChild(resultItem);
    });
}

async function showElementDetails(elementName) {
    const loadingSpinner = document.getElementById('loading-spinner');
    loadingSpinner.classList.remove('hidden');

    try {
        const element = await getElementDetails(elementName);
        console.log("Détails de l'élément reçus:", element);

        displayElementInfo(element);
        displayElementContexts(element);
        
        if (element.Relations && element.Relations.length > 0) {
            createRelationsGraph(element);
        } else {
            elementRelations.innerHTML = '<p>Aucune relation disponible pour cet élément.</p>';
        }
    } catch (error) {
        console.error('Erreur lors de la récupération des détails de l\'élément:', error);
        showErrorMessage('Impossible de charger les détails de l\'élément.');
    } finally {
        loadingSpinner.classList.add('hidden');
    }
}

function displayElementInfo(element) {
    elementDetails.innerHTML = `
        <h2>${escapeHtml(element.Name)}</h2>
        <p><strong>Type:</strong> ${escapeHtml(element.Type)}</p>
        <p><strong>Description:</strong> ${escapeHtml(element.Description)}</p>
        <p><strong>Positions:</strong> ${element.Positions.join(', ')}</p>
    `;
}

function displayElementContexts(element) {
    elementContexts.innerHTML = '<h3>Contextes</h3>';
    if (element.Contexts && element.Contexts.length > 0) {
        element.Contexts.forEach((ctx, index) => {
            const highlightedBefore = highlightElement(ctx.before.join(' '), element.Name);
            const highlightedAfter = highlightElement(ctx.after.join(' '), element.Name);
            
            elementContexts.innerHTML += `
                <div class="context">
                    <h4>Contexte ${index + 1}</h4>
                    <p><strong>Avant:</strong> ${highlightedBefore}</p>
                    <p><strong>Élément du contexte:</strong> <mark>${escapeHtml(ctx.element)}</mark></p>
                    <p><strong>Après:</strong> ${highlightedAfter}</p>
                    <p><strong>Position:</strong> ${ctx.position}</p>
                </div>
            `;
        });

        // Ajouter des écouteurs d'événements pour les éléments surlignés
        const highlightedElements = elementContexts.querySelectorAll('mark');
        highlightedElements.forEach(el => {
            el.style.cursor = 'pointer';
            el.addEventListener('click', () => {
                const searchInput = document.getElementById('search-input');
                searchInput.value = el.textContent;
                handleSearch(); // Assurez-vous que cette fonction est accessible globalement
            });
        });
    } else {
        elementContexts.innerHTML += '<p>Aucun contexte disponible</p>';
    }
}

function highlightElement(text, elementName) {
    const regex = new RegExp(elementName, 'gi');
    return text.replace(regex, match => `<mark>${match}</mark>`);
}

function escapeHtml(unsafe) {
    return unsafe
         .replace(/&/g, "&amp;")
         .replace(/</g, "&lt;")
         .replace(/>/g, "&gt;")
         .replace(/"/g, "&quot;")
         .replace(/'/g, "&#039;");
}