// graph.js
import * as d3 from 'https://cdn.jsdelivr.net/npm/d3@7/+esm';

export function createRelationsGraph(element) {
    const elementRelations = document.getElementById('element-relations');
    elementRelations.innerHTML = ''; // Clear previous content

    const width = elementRelations.clientWidth;
    const height = 400;
    
    const svg = d3.select(elementRelations)
        .append("svg")
        .attr("width", width)
        .attr("height", height);

    const simulation = d3.forceSimulation()
        .force("link", d3.forceLink().id(d => d.id).distance(100))
        .force("charge", d3.forceManyBody().strength(-500))
        .force("center", d3.forceCenter(width / 2, height / 2));

    const links = element.Relations.map(r => ({source: element.Name, target: r.Target, type: r.Type}));
    const nodes = [{id: element.Name, group: 1}, ...element.Relations.map(r => ({id: r.Target, group: 2}))];

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

    simulation
        .nodes(nodes)
        .on("tick", ticked);

    simulation.force("link")
        .links(links);

    function ticked() {
        link
            .attr("x1", d => d.source.x)
            .attr("y1", d => d.source.y)
            .attr("x2", d => d.target.x)
            .attr("y2", d => d.target.y);

        node
            .attr("transform", d => `translate(${d.x},${d.y})`);
    }
}