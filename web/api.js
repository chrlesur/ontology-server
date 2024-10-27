// api.js
const API_BASE_URL = '/api';

// Charger la liste des ontologies
export async function loadOntologies() {
    const response = await fetch(`${API_BASE_URL}/ontologies`);
    if (!response.ok) throw new Error('Erreur lors du chargement des ontologies');
    return await response.json();
}

// Charger les types d'éléments (à implémenter côté serveur)
export async function loadElementTypes() {
    // TODO: Implémenter l'appel API réel
    return ['Concept', 'Relation', 'Instance'];
}

// Rechercher dans les ontologies
export async function searchOntologies(query, ontologyId, elementType) {
    let url = `${API_BASE_URL}/search?q=${encodeURIComponent(query)}`;
    if (ontologyId) url += `&ontology_id=${encodeURIComponent(ontologyId)}`;
    if (elementType) url += `&element_type=${encodeURIComponent(elementType)}`;

    const response = await fetch(url);
    if (!response.ok) throw new Error('Erreur lors de la recherche');
    const data = await response.json();
    
    return Array.isArray(data) ? data : [];
}

// api.js
export async function getElementDetails(elementName) {
    const url = `${API_BASE_URL}/elements/details/${encodeURIComponent(elementName)}`;
    console.log('Fetching element details from:', url);
    try {
        const response = await fetch(url);
        console.log('Response status:', response.status);
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        const data = await response.json();
        console.log('Received element details:', data);
        
        // Assurez-vous que toutes les propriétés attendues sont présentes
        return {
            Name: data.Name || '',
            Type: data.Type || '',
            Description: data.Description || '',
            Positions: data.Positions || [],
            Relations: data.Relations || [],
            Contexts: data.Contexts || []
        };
    } catch (error) {
        console.error('Erreur lors de la récupération des détails de l\'élément:', error);
        throw new Error('Erreur lors de la récupération des détails de l\'élément');
    }
}

// Charger une ontologie
export async function uploadOntology(formData) {
    const response = await fetch(`${API_BASE_URL}/ontologies/load`, {
        method: 'POST',
        body: formData
    });
    if (!response.ok) throw new Error('Erreur lors du chargement de l\'ontologie');
    return await response.json();
}

// api.js

export async function getElementRelations(elementName) {
    const url = `${API_BASE_URL}/elements/relations/${encodeURIComponent(elementName)}`;
    try {
        const response = await fetch(url);
        if (!response.ok) {
            if (response.status === 404) {
                return []; // Retourne un tableau vide si aucune relation n'est trouvée
            }
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        const data = await response.json();
        return data;
    } catch (error) {
        console.error('Erreur lors de la récupération des relations de l\'élément:', error);
        return []; // Retourne un tableau vide en cas d'erreur
    }
}