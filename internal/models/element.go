package models

type JSONContext struct {
	Position    int      `json:"position"`
	Before      []string `json:"before"`
	After       []string `json:"after"`
	Element     string   `json:"element"`
	Length      int      `json:"length"`
	StartOffset int      // Nouvelle: position de début de l'intervalle
	EndOffset   int      // Nouvelle: position de fin de l'intervalle
}
