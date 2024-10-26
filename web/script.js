// Constantes pour l'API
const API_BASE_URL = '/api';
const SEARCH_ENDPOINT = `${API_BASE_URL}/search`;
const ONTOLOGIES_ENDPOINT = `${API_BASE_URL}/ontologies`;
const ELEMENT_DETAILS_ENDPOINT = `${API_BASE_URL}/elements`;
const UPLOAD_ENDPOINT = `${API_BASE_URL}/ontologies/load`;

// Éléments du DOM
const uploadForm = document.getElementById('upload-form');
const searchForm = document.getElementById('search-form');
const searchInput = document.getElementById('search-input');
const ontologySelect = document.getElementById('ontology-select');
const elementTypeSelect = document.getElementById('element-type-select');
const resultsList = document.getElementById('results-list');
const pagination = document.getElementById('pagination');
const elementDetails = document.getElementById('element-details');
const elementContext = document.getElementById('element-context');
const elementRelations = document.getElementById('element-relations');

// Variables globales
let currentPage = 1;
const resultsPerPage = 10;

// Initialisation
document.addEventListener('DOMContentLoaded', init);

async function init() {
    try {
        await loadOntologies();
        await loadElementTypes();
        uploadForm.addEventListener('submit', handleUpload);
        searchForm.addEventListener('submit', handleSearch);
        searchInput.addEventListener('input', debounce(handleSearch, 300));
    } catch (error) {
        console.error('Erreur lors de l\'initialisation:', error);
        showErrorMessage('Une erreur est survenue lors de l\'initialisation de l\'application.');
    }
}

// Chargement des ontologies
async function loadOntologies() {
    try {
        const response = await fetch(ONTOLOGIES_ENDPOINT);
        if (!response.ok) throw new Error('Erreur lors du chargement des ontologies');
        const ontologies = await response.json();
        ontologySelect.innerHTML = '<option value="">Toutes les ontologies</option>';
        ontologies.forEach(ontology => {
            const option = document.createElement('option');
            option.value = ontology.id;
            option.textContent = ontology.name;
            ontologySelect.appendChild(option);
        });
    } catch (error) {
        console.error('Erreur lors du chargement des ontologies:', error);
        showErrorMessage('Impossible de charger la liste des ontologies.');
    }
}

// Chargement des types d'éléments (à implémenter côté serveur si nécessaire)
async function loadElementTypes() {
    // Exemple statique, à remplacer par un appel API si implémenté côté serveur
    const types = ['Concept', 'Relation', 'Instance'];
    types.forEach(type => {
        const option = document.createElement('option');
        option.value = type;
        option.textContent = type;
        elementTypeSelect.appendChild(option);
    });
}

// Gestion du chargement de fichiers
async function handleUpload(event) {
    event.preventDefault();
    const formData = new FormData(uploadForm);
    try {
        const response = await fetch(UPLOAD_ENDPOINT, {
            method: 'POST',
            body: formData
        });
        if (!response.ok) throw new Error('Erreur lors du chargement des fichiers');
        const result = await response.json();
        showMessage(result.message);
        await loadOntologies(); // Recharger la liste des ontologies
    } catch (error) {
        console.error('Erreur lors du chargement des fichiers:', error);
        showErrorMessage('Impossible de charger les fichiers.');
    }
}

// Gestion de la recherche
async function handleSearch(event) {
    if (event) event.preventDefault();
    currentPage = 1;
    await performSearch();
}

async function performSearch() {
    const query = searchInput.value.trim();
    const ontologyId = ontologySelect.value;
    const elementType = elementTypeSelect.value;

    if (query === "") {
        displayResults([]);
        return;
    }

    try {
        const url = new URL(SEARCH_ENDPOINT, window.location.origin);
        url.searchParams.append('q', query);
        if (ontologyId) url.searchParams.append('ontology_id', ontologyId);
        if (elementType) url.searchParams.append('element_type', elementType);
        url.searchParams.append('page', currentPage);
        url.searchParams.append('per_page', resultsPerPage);

        const response = await fetch(url);
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        const data = await response.json();
        console.log("Données reçues de l'API:", data); // Ajout de log

        if (!Array.isArray(data)) {
            console.error('La réponse de l\'API n\'est pas un tableau:', data);
            throw new Error('Format de réponse inattendu');
        }

        displayResults(data);
    } catch (error) {
        console.error('Erreur lors de la recherche:', error);
        showErrorMessage('Une erreur est survenue lors de la recherche.');
        displayResults([]); // Afficher une liste vide en cas d'erreur
    }
}

// Affichage des résultats
function displayResults(results) {
    console.log("Affichage des résultats:", results); // Ajout de log
    resultsList.innerHTML = '';
    if (!Array.isArray(results) || results.length === 0) {
        resultsList.innerHTML = '<p>Aucun résultat trouvé.</p>';
        return;
    }

    results.forEach(result => {
        const resultItem = document.createElement('div');
        resultItem.className = 'result-item';
        resultItem.innerHTML = `
            <h3>${escapeHtml(result.ElementName)}</h3>
            <p>Type: ${escapeHtml(result.ElementType)}</p>
            <p>Ontologie: ${escapeHtml(result.OntologyID)}</p>
        `;
        resultItem.addEventListener('mouseover', () => showElementDetails(result));
        resultsList.appendChild(resultItem);
    });
}



// Affichage des détails d'un élément
async function showElementDetails(element) {
    console.log("Détails de l'élément:", element);

    elementDetails.innerHTML = `
        <h3>${escapeHtml(element.ElementName)}</h3>
        <p><strong>Type:</strong> ${escapeHtml(element.ElementType)}</p>
        <p><strong>Description:</strong> ${escapeHtml(element.Description)}</p>
        <p><strong>Ontologie:</strong> ${escapeHtml(element.OntologyID)}</p>
    `;

    // Affichage du contexte
    if (element.Context) {
        const contextParts = element.Context.split(`[${element.ElementName}]`);
        let beforeContext = '';
        let afterContext = '';
        
        if (contextParts.length > 1) {
            beforeContext = contextParts[0].trim();
            afterContext = contextParts[1].trim();
        } else {
            // Si le split n'a pas fonctionné, on affiche tout le contexte
            afterContext = element.Context.trim();
        }

        elementContext.innerHTML = `
            <h4>Contexte</h4>
            <p><strong>Avant:</strong> ${escapeHtml(beforeContext)}</p>
            <p><strong>Élément:</strong> <mark>${escapeHtml(element.ElementName)}</mark></p>
            <p><strong>Après:</strong> ${escapeHtml(afterContext)}</p>
        `;
    } else {
        elementContext.innerHTML = '<p>Aucun contexte disponible</p>';
    }

    // Récupérer les relations
    try {
        const response = await fetch(`/api/elements/relations/${encodeURIComponent(element.ElementName)}`);
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        const relations = await response.json();
        element.Relations = relations;
        console.log("Relations récupérées:", relations);
    } catch (error) {
        console.error('Erreur lors de la récupération des relations:', error);
        element.Relations = [];
    }

    // Création du graphique SVG pour les relations
    createRelationsGraph(element);
}

// Création du graphique SVG pour les relations
function createRelationsGraph(element) {
    const width = 900;
    const height = 700;
    const svg = d3.select(elementRelations)
        .html("")
        .append("svg")
        .attr("width", "100%")
        .attr("height", height)
        .attr("viewBox", [-width / 2, -height / 2, width, height]);

    const nodes = [
        { id: element.ElementName, group: 1, main: true },
        ...element.Relations.map(r => ({ id: r.Source === element.ElementName ? r.Target : r.Source, group: 2 }))
    ];

    const links = element.Relations.map(r => ({
        source: element.ElementName,
        target: r.Source === element.ElementName ? r.Target : r.Source,
        type: r.Type
    }));

    const colorScale = d3.scaleOrdinal(d3.schemePastel1);

    const simulation = d3.forceSimulation(nodes)
        .force("link", d3.forceLink(links).id(d => d.id).distance(300))
        .force("charge", d3.forceManyBody().strength(-4000))
        .force("center", d3.forceCenter())
        .force("collision", d3.forceCollide().radius(d => Math.max(d.width, d.height) / 2 + 20));

    const link = svg.append("g")
        .selectAll("line")
        .data(links)
        .join("line")
        .attr("stroke", "#999")
        .attr("stroke-opacity", 0.6)
        .attr("stroke-width", 2);

    const node = svg.append("g")
        .selectAll("g")
        .data(nodes)
        .join("g")
        .call(drag(simulation));

    node.each(function(d) {
        const el = d3.select(this);
        const text = el.append("text")
            .text(d.id)
            .attr("text-anchor", "middle")
            .attr("dy", ".35em")
            .attr("fill", "black")
            .style("font-size", "16px")
            .style("font-weight", "bold");

        const bbox = text.node().getBBox();
        d.width = bbox.width + 40;
        d.height = bbox.height + 30;

        el.insert("rect", "text")
            .attr("x", -d.width / 2)
            .attr("y", -d.height / 2)
            .attr("width", d.width)
            .attr("height", d.height)
            .attr("rx", 10)
            .attr("ry", 10)
            .attr("fill", d => d.main ? colorScale(0) : colorScale(1))
            .attr("stroke", "#fff")
            .attr("stroke-width", 2);

        text.call(wrap, d.width - 20);
    });

    const linkText = svg.append("g")
        .selectAll("g")
        .data(links)
        .join("g");

    linkText.append("rect")
        .attr("fill", "white")
        .attr("opacity", 0.8)
        .attr("rx", 5)
        .attr("ry", 5);

    linkText.append("text")
        .attr("text-anchor", "middle")
        .attr("dy", "0.35em")
        .text(d => d.type)
        .style("font-size", "14px")
        .style("font-weight", "bold")
        .each(function(d) {
            const bbox = this.getBBox();
            d3.select(this.parentNode).select("rect")
                .attr("x", -bbox.width / 2 - 4)
                .attr("y", -bbox.height / 2 - 2)
                .attr("width", bbox.width + 8)
                .attr("height", bbox.height + 4);
        });

    simulation.on("tick", () => {
        link
            .attr("x1", d => d.source.x)
            .attr("y1", d => d.source.y)
            .attr("x2", d => d.target.x)
            .attr("y2", d => d.target.y);

        node
            .attr("transform", d => `translate(${d.x},${d.y})`);

        linkText
            .attr("transform", d => {
                const dx = d.target.x - d.source.x;
                const dy = d.target.y - d.source.y;
                const angle = Math.atan2(dy, dx) * 180 / Math.PI;
                const midX = (d.source.x + d.target.x) / 2;
                const midY = (d.source.y + d.target.y) / 2;
                return `translate(${midX},${midY}) rotate(${angle})`;
            });
    });

    function drag(simulation) {
        function dragstarted(event) {
            if (!event.active) simulation.alphaTarget(0.3).restart();
            event.subject.fx = event.subject.x;
            event.subject.fy = event.subject.y;
        }
        
        function dragged(event) {
            event.subject.fx = event.x;
            event.subject.fy = event.y;
        }
        
        function dragended(event) {
            if (!event.active) simulation.alphaTarget(0);
            event.subject.fx = null;
            event.subject.fy = null;
        }
        
        return d3.drag()
            .on("start", dragstarted)
            .on("drag", dragged)
            .on("end", dragended);
    }

    function wrap(text, width) {
        text.each(function() {
            var text = d3.select(this),
                words = text.text().split(/\s+/).reverse(),
                word,
                line = [],
                lineNumber = 0,
                lineHeight = 1.1,
                y = text.attr("y"),
                dy = parseFloat(text.attr("dy")),
                tspan = text.text(null).append("tspan").attr("x", 0).attr("y", y).attr("dy", dy + "em");
            while (word = words.pop()) {
                line.push(word);
                tspan.text(line.join(" "));
                if (tspan.node().getComputedTextLength() > width) {
                    line.pop();
                    tspan.text(line.join(" "));
                    line = [word];
                    tspan = text.append("tspan").attr("x", 0).attr("y", y).attr("dy", ++lineNumber * lineHeight + dy + "em").text(word);
                }
            }
        });
    }
}
function createCircle(cx, cy, r, fill, text) {
    const group = document.createElementNS("http://www.w3.org/2000/svg", "g");
    
    const circle = document.createElementNS("http://www.w3.org/2000/svg", "circle");
    circle.setAttribute("cx", cx);
    circle.setAttribute("cy", cy);
    circle.setAttribute("r", r);
    circle.setAttribute("fill", fill);
    
    const textElement = document.createElementNS("http://www.w3.org/2000/svg", "text");
    textElement.setAttribute("x", cx);
    textElement.setAttribute("y", cy);
    textElement.setAttribute("text-anchor", "middle");
    textElement.setAttribute("dominant-baseline", "middle");
    textElement.setAttribute("fill", "white");
    textElement.textContent = text;
    
    group.appendChild(circle);
    group.appendChild(textElement);
    return group;
}

function createLine(x1, y1, x2, y2) {
    const line = document.createElementNS("http://www.w3.org/2000/svg", "line");
    line.setAttribute("x1", x1);
    line.setAttribute("y1", y1);
    line.setAttribute("x2", x2);
    line.setAttribute("y2", y2);
    line.setAttribute("stroke", "#999");
    line.setAttribute("stroke-width", "2");
    return line;
}

function createText(x, y, text) {
    const textElement = document.createElementNS("http://www.w3.org/2000/svg", "text");
    textElement.setAttribute("x", x);
    textElement.setAttribute("y", y);
    textElement.setAttribute("text-anchor", "middle");
    textElement.setAttribute("dominant-baseline", "middle");
    textElement.setAttribute("fill", "#333");
    textElement.textContent = text;
    return textElement;
}

// Utilitaires
function showErrorMessage(message) {
    console.error(message);
    const errorDiv = document.createElement('div');
    errorDiv.className = 'error-message';
    errorDiv.textContent = message;
    resultsList.innerHTML = '';
    resultsList.appendChild(errorDiv);
}

function showMessage(message) {
    // Implémentez cette fonction pour afficher les messages à l'utilisateur
    alert(message);
}

function escapeHtml(unsafe) {
    return unsafe
         .replace(/&/g, "&amp;")
         .replace(/</g, "&lt;")
         .replace(/>/g, "&gt;")
         .replace(/"/g, "&quot;")
         .replace(/'/g, "&#039;");
}

function debounce(func, wait) {
    let timeout;
    return function executedFunction(...args) {
        const later = () => {
            clearTimeout(timeout);
            func(...args);
        };
        clearTimeout(timeout);
        timeout = setTimeout(later, wait);
    };
}