import * as d3 from 'https://cdn.jsdelivr.net/npm/d3@7/+esm';

export function createRelationsGraph(element, relations) {
    const graphContainer = document.getElementById('element-relations-graph');
    graphContainer.innerHTML = ''; // Clear previous content

    const width = graphContainer.clientWidth;
    const height = graphContainer.clientHeight;
    
    const svg = d3.select(graphContainer)
        .append("svg")
        .attr("width", width)
        .attr("height", height);

    // Créer un ensemble unique de noms de nœuds
    const nodeNames = new Set([element.Name, ...relations.flatMap(r => [r.Source, r.Target])]);

    // Créer les nœuds en tant qu'objets
    const nodes = Array.from(nodeNames).map(name => ({
        id: name,
        group: name === element.Name ? 1 : 2
    }));

    // Créer les liens
    const links = relations.map(r => ({
        source: r.Source,
        target: r.Target,
        type: r.Type
    }));

    const simulation = d3.forceSimulation(nodes)
        .force("link", d3.forceLink(links).id(d => d.id).distance(100))
        .force("charge", d3.forceManyBody().strength(-300))
        .force("center", d3.forceCenter(width / 2, height / 2));

    const link = svg.append("g")
        .selectAll("line")
        .data(links)
        .enter().append("line")
        .attr("stroke", "#999")
        .attr("stroke-width", 2);

    const node = svg.append("g")
        .selectAll("g")
        .data(nodes)
        .enter().append("g");

    node.append("circle")
        .attr("r", 20)
        .attr("fill", d => d.group === 1 ? "#ff9900" : "#66b3ff");

    node.append("text")
        .text(d => d.id)
        .attr("text-anchor", "middle")
        .attr("dy", ".35em")
        .attr("font-size", "12px");

    simulation.on("tick", () => {
        link
            .attr("x1", d => d.source.x)
            .attr("y1", d => d.source.y)
            .attr("x2", d => d.target.x)
            .attr("y2", d => d.target.y);

        node
            .attr("transform", d => `translate(${d.x},${d.y})`);
    });
}