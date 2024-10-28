import { performSearch } from './search.js'; // Assurez-vous que cette fonction existe et est exportée

export function createRelationsGraph(element, relations) {
    const graphContainer = document.getElementById('element-relations-graph');
    graphContainer.innerHTML = ''; // Clear previous content

    // Fonction pour échapper les chaînes contenant des espaces ou des caractères spéciaux
    const escapeString = (str) => {
        return str.replace(/[^a-zA-Z0-9]/g, '_');
    };

    // Créer un ensemble pour stocker les nœuds uniques
    const nodes = new Set([element.Name]);

    // Créer la définition du diagramme Mermaid
    let mermaidDef = 'graph TD\n';

    // Ajouter le nœud principal
    const escapedName = escapeString(element.Name);
    const escapedType = escapeString(element.Type);
    mermaidDef += `    ${escapedName}["${element.Name}<br/>${element.Type}"]\n`;

    // Ajouter les relations
    relations.forEach(relation => {
        const escapedSource = escapeString(relation.Source);
        const escapedTarget = escapeString(relation.Target);
        const escapedType = escapeString(relation.Type);
        
        nodes.add(relation.Source);
        nodes.add(relation.Target);
        mermaidDef += `    ${escapedSource}["${relation.Source}"] -- "${relation.Type}" --> ${escapedTarget}["${relation.Target}"]\n`;
    });

    // Ajouter les nœuds isolés
    nodes.forEach(node => {
        const escapedNode = escapeString(node);
        if (node !== element.Name && !relations.some(r => r.Source === node || r.Target === node)) {
            mermaidDef += `    ${escapedNode}["${node}"]\n`;
        }
    });

    // Créer un élément div pour le diagramme
    const mermaidDiv = document.createElement('div');
    mermaidDiv.className = 'mermaid';
    mermaidDiv.textContent = mermaidDef;
    graphContainer.appendChild(mermaidDiv);

    // Initialiser et rendre le diagramme Mermaid
    mermaid.initialize({ 
        startOnLoad: true, 
        theme: 'default',
        flowchart: {
            useMaxWidth: false,
            htmlLabels: true,
            curve: 'basis'
        }
    });

    mermaid.render('graphDiv', mermaidDef).then(result => {
        mermaidDiv.innerHTML = result.svg;

        // Appliquer panzoom au SVG
        const svg = mermaidDiv.querySelector('svg');
        panzoom(svg, {
            maxZoom: 5,
            minZoom: 0.5,
            initialZoom: 0.8
        });

        // Ajouter des événements de clic aux éléments du graphique
        svg.querySelectorAll('.nodeLabel, .edgeLabel').forEach(element => {
            element.style.cursor = 'pointer';
            element.addEventListener('click', (event) => {
                const text = event.target.textContent.split('\n')[0]; // Prendre seulement la première ligne
                performSearch(text);
            });
        });
    });
}