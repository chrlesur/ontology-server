package models

// JSONContext représente le contexte d'un élément
type JSONContext struct {
	Position     int      `json:"position"`
	FileID       string   `json:"file_id"`
	FilePosition int      `json:"file_position"`
	Before       []string `json:"before"`
	After        []string `json:"after"`
	Element      string   `json:"element"`
	Length       int      `json:"length"`
	StartOffset  int
	EndOffset    int
}
