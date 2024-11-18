// api.js
const API_BASE_URL = '/api';

// Charger la liste des ontologies
export async function loadOntologies() {
    try {
        const response = await fetch(`${API_BASE_URL}/ontologies`);
        if (!response.ok) {
            throw new Error('Error loading ontologies');
        }
        const ontologies = await response.json();
        console.log("Raw ontologies data:", ontologies);
        
        // Récupérer les métadonnées pour chaque ontologie
        const ontologiesWithMetadata = await Promise.all(ontologies.map(async (ontology) => {
            try {
                const metadataResponse = await fetch(`${API_BASE_URL}/ontologies/${ontology.id}/metadata`);
                if (metadataResponse.ok) {
                    const metadata = await metadataResponse.json();
                    return { ...ontology, Source: metadata };
                }
            } catch (error) {
                console.error(`Error fetching metadata for ontology ${ontology.id}:`, error);
            }
            return ontology;
        }));
        
        console.log("Ontologies with metadata:", ontologiesWithMetadata);
        return ontologiesWithMetadata;
    } catch (error) {
        console.error('Error loading ontologies:', error);
        throw error;
    }
}

// Charger les types d'éléments (à implémenter côté serveur)
export async function loadElementTypes() {
    // TODO: Implémenter l'appel API réel
    return ['Concept', 'Relation', 'Instance'];
}

// Rechercher dans les ontologies
export async function searchOntologies(query, fileId, elementType) {
    if (!query) {
        throw new Error('Un terme de recherche est requis');
    }

    let url = `${API_BASE_URL}/search?q=${encodeURIComponent(query)}`;
    if (fileId) url += `&file_id=${encodeURIComponent(fileId)}`;
    if (elementType) url += `&element_type=${encodeURIComponent(elementType)}`;

    console.log("Search URL:", url);

    const response = await fetch(url);
    if (!response.ok) {
        const errorText = await response.text();
        console.error("Search API error:", errorText);
        throw new Error('Erreur lors de la recherche: ' + errorText);
    }
    
    const data = await response.json();
    console.log('Résultats bruts de la recherche:', data);

    return Array.isArray(data) ? data.map(item => ({
        ...item,
        sourceFile: item.Source?.source_file || 'Unknown',
        sourceMetadata: item.Source || null
    })) : [];
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

// Charger une ontologie avec ses métadonnées
export async function uploadOntology(formData) {
    try {
        const response = await fetch(`${API_BASE_URL}/ontologies/load`, {
            method: 'POST',
            body: formData
        });

        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.error || 'Failed to upload ontology');
        }

        return await response.json();
    } catch (error) {
        console.error('Error uploading ontology:', error);
        throw error;
    }
}

export async function getElementRelations(elementName) {
    try {
        // Encoder proprement le nom de l'élément pour l'URL
        const encodedName = encodeURIComponent(elementName);
        const url = `${API_BASE_URL}/elements/relations/${encodedName}`;
        
        const response = await fetch(url);
        if (!response.ok && response.status !== 404) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        
        // En cas de 404 ou si la réponse est ok, parser le JSON
        const data = await response.json();
        
        // S'assurer de toujours retourner un tableau
        return Array.isArray(data) ? data : [];
        
    } catch (error) {
        console.error('Erreur lors de la récupération des relations:', error);
        return []; // Retourner un tableau vide en cas d'erreur
    }
}

// Récupérer les métadonnées d'une ontologie
export async function getOntologyMetadata(ontologyId) {
    try {
        const response = await fetch(`${API_BASE_URL}/ontologies/${ontologyId}/metadata`);
        if (!response.ok) {
            throw new Error('Failed to fetch ontology metadata');
        }
        return await response.json();
    } catch (error) {
        console.error('Error fetching ontology metadata:', error);
        throw error;
    }
}