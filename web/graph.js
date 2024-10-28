import { performSearch } from './search.js';

export function createRelationsGraph(element, relations) {
    const graphContainer = document.getElementById('element-relations-graph');
    graphContainer.innerHTML = '';

    const escapeId = (str) => {
        return str.replace(/[^\w\s-]/g, '_').replace(/\s+/g, '_');
    };

    let mermaidDef = 'graph LR\n';

    // Créer un Map pour stocker les nœuds uniques
    const nodes = new Map();

    // Ajouter le nœud principal
    const mainNodeId = `node_${escapeId(element.Name)}`;
    nodes.set(mainNodeId, element.Name);
    mermaidDef += `    ${mainNodeId}["${element.Name}<br><i>${element.Type}</i>"]\n`;

    // Ajouter les relations
    relations.forEach(relation => {
        const sourceId = `node_${escapeId(relation.Source)}`;
        const targetId = `node_${escapeId(relation.Target)}`;

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
        const instance = panzoom(svg, {
            maxZoom: 5,
            minZoom: 0.1,
            initialZoom: 1
        });

        const bbox = svg.getBBox();
        const viewportWidth = graphContainer.clientWidth;
        const viewportHeight = graphContainer.clientHeight;
        instance.moveTo(viewportWidth/2 - bbox.width/2, viewportHeight/2 - bbox.height/2);

        svg.querySelectorAll('.node').forEach(node => {
            node.style.cursor = 'pointer';
            node.addEventListener('click', (event) => {
                const nodeId = node.id;
                const originalName = nodes.get(nodeId);
                if (originalName) {
                    performSearch(originalName);
                }
            });
        });
    });
}