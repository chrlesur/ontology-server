import { performSearch } from './search.js';

export function createRelationsGraph(element, relations) {
    console.log("Création du graphique pour:", element.Name);
    const graphContainer = document.getElementById('element-relations-graph');
    graphContainer.innerHTML = '';

    const escapeId = (str) => {
        return str.replace(/[^\w\s-]/g, '_').replace(/\s+/g, '_');
    };

    let mermaidDef = 'graph LR\n';

    // Créer un Map pour stocker les nœuds uniques
    const nodes = new Map();

    // Ajouter le nœud principal
    const mainNodeId = escapeId(element.Name);
    nodes.set(mainNodeId, element.Name);
    mermaidDef += `    ${mainNodeId}["${element.Name}<br><i>${element.Type}</i>"]\n`;

    // Ajouter les relations
    relations.forEach(relation => {
        const sourceId = escapeId(relation.Source);
        const targetId = escapeId(relation.Target);

        // Ajouter les nœuds s'ils n'existent pas déjà
        if (!nodes.has(sourceId)) {
            nodes.set(sourceId, relation.Source);
            mermaidDef += `    ${sourceId}["${relation.Source}"]\n`;
        }
        if (!nodes.has(targetId)) {
            nodes.set(targetId, relation.Target);
            mermaidDef += `    ${targetId}["${relation.Target}"]\n`;
        }

        // Ajouter la relation
        mermaidDef += `    ${sourceId} -->|"${relation.Type}"| ${targetId}\n`;
    });

    const mermaidDiv = document.createElement('div');
    mermaidDiv.className = 'mermaid';
    mermaidDiv.textContent = mermaidDef;
    graphContainer.appendChild(mermaidDiv);

    mermaid.initialize({ 
        startOnLoad: true, 
        theme: 'default',
        flowchart: {
            useMaxWidth: true,
            htmlLabels: true,
            curve: 'basis'
        },
        securityLevel: 'loose'
    });

    mermaid.render('graphDiv', mermaidDef).then(result => {
        mermaidDiv.innerHTML = result.svg;

        const svg = mermaidDiv.querySelector('svg');
        
        // Réappliquer panzoom
        const instance = panzoom(svg, {
            maxZoom: 5,
            minZoom: 0.1,
            initialZoom: 1
        });

        // Centrer le graphique
        const bbox = svg.getBBox();
        const viewportWidth = graphContainer.clientWidth;
        const viewportHeight = graphContainer.clientHeight;
        instance.moveTo(viewportWidth/2 - bbox.width/2, viewportHeight/2 - bbox.height/2);

        // Ajouter des gestionnaires d'événements de clic aux nœuds
        svg.querySelectorAll('.node').forEach(node => {
            node.style.cursor = 'pointer';
            node.addEventListener('click', function(event) {
                event.preventDefault();
                event.stopPropagation();
                const nodeText = this.querySelector('tspan')?.textContent || this.textContent;
                console.log("Nœud cliqué, texte:", nodeText);
                if (nodeText) {
                    // Extraire uniquement le nom de l'élément (première partie avant un espace ou un caractère spécial)
                    const elementName = nodeText.split(/\s/)[0].trim();
                    console.log("Recherche pour l'élément:", elementName);
                    performSearch(elementName);
                }
            });
        });

        console.log("Graphique créé avec", nodes.size, "nœuds");
    }).catch(error => {
        console.error("Erreur lors de la création du graphique:", error);
    });
}