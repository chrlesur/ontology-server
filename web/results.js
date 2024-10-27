// results.js
import { getElementDetails, getElementRelations } from './api.js';
import { createRelationsGraph } from './graph.js';
import { showErrorMessage } from './main.js';

const resultsList = document.getElementById('results-list');
const elementDetails = document.getElementById('element-details');
const elementContexts = document.getElementById('element-contexts');
const elementRelations = document.getElementById('element-relations');

export function displayResults(results) {
    const resultsList = document.getElementById('results-list');
    resultsList.innerHTML = '';
    
    results.forEach(result => {
        const resultItem = document.createElement('div');
        resultItem.className = 'result-item';
        resultItem.innerHTML = `
            <h3>${escapeHtml(result.ElementName)}</h3>
            <p>${escapeHtml(result.Description)}</p>
        `;
        resultItem.addEventListener('click', () => {
            // Supprimer la classe 'selected' de tous les éléments
            document.querySelectorAll('.result-item').forEach(item => item.classList.remove('selected'));
            // Ajouter la classe 'selected' à l'élément cliqué
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
        const element = await getElementDetails(elementName);
        console.log("Détails de l'élément reçus:", element);

        displayElementInfo(element);
        displayElementContexts(element);
        
        const relations = await getElementRelations(elementName);
        console.log("Relations de l'élément reçues:", relations);

        if (relations && relations.length > 0) {
            displayElementRelations(relations);
            createRelationsGraph(element, relations);
        } else {
            displayElementRelations([]);
            // Vous pouvez également choisir de cacher ou effacer le graphique ici
        }
    } catch (error) {
        console.error('Erreur lors de la récupération des détails ou des relations de l\'élément:', error);
        showErrorMessage('Impossible de charger les détails ou les relations de l\'élément.');
    } finally {
        if (loadingSpinner) loadingSpinner.classList.add('hidden');
    }
}

function displayElementRelations(relations) {
    const relationsList = document.getElementById('element-relations-list');
    if (!relationsList) {
        console.error("Element 'element-relations-list' not found");
        return;
    }

    relationsList.innerHTML = '';
    if (relations && relations.length > 0) {
        const ul = document.createElement('ul');
        relations.forEach(relation => {
            const li = document.createElement('li');
            li.textContent = `${relation.Source} ${relation.Type} ${relation.Target}`;
            ul.appendChild(li);
        });
        relationsList.appendChild(ul);
    } else {
        relationsList.innerHTML = '<p>Aucune relation disponible pour cet élément.</p>';
    }
}

function displayElementInfo(element) {
    const positions = element.Positions && Array.isArray(element.Positions) 
        ? element.Positions.join(', ') 
        : 'Non spécifié';

    elementDetails.innerHTML = `
        <h2>${escapeHtml(element.Name)}</h2>
        <p><strong>Type:</strong> ${escapeHtml(element.Type)}</p>
        <p><strong>Description:</strong> ${escapeHtml(element.Description)}</p>
        <p><strong>Positions:</strong> ${positions}</p>
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
                // Déclencher un événement personnalisé
                const event = new CustomEvent('performSearch', {
                    detail: { query: el.textContent }
                });
                document.dispatchEvent(event);
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

function displayRelationsList(relations) {
    const relationsList = document.getElementById('element-relations-list');
    relationsList.innerHTML = '';
    relations.forEach(relation => {
        const relationItem = document.createElement('div');
        relationItem.className = 'relation-item';
        relationItem.textContent = `${relation.Source} ${relation.Type} ${relation.Target}`;
        relationsList.appendChild(relationItem);
    });
}