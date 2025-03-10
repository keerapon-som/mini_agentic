package entities

type VectorDocument struct {
	Command      string       `json:"command"`
	DocumentData DocumentData `json:"document_data"`
}

type RawDocument struct {
	Command     string `json:"command"`
	Description string `json:"description"`
}

type DocumentData struct {
	Description       string    `json:"description"`
	Vector            []float64 `json:"vector_lists"`
	EmbedingModelName string    `json:"embeding_model_name"`
}
