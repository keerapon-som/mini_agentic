package main

import (
	"fmt"
	"miniagi/api"
	"miniagi/ollama"
)

func TurnRawToVectorDocument(agi *api.AGI) {
	err := agi.TurnRawToVectorDocument()
	if err != nil {
		panic(err)
	}
}

func main() {
	ttv := api.NewTextToVector(
		ollama.NewOllama("http://localhost:11434"),
	)
	agi := api.NewAGI(ttv, api.NewVectorMapper(), "file/command.json", "file/output.json", api.MODEL_ALL_MINILM)

	// TurnRawToVectorDocument(agi)

	resp, err := agi.GetClosetDocument("Please remove some file")
	if err != nil {
		panic(err)
	}
	fmt.Println("Command is : ", resp.Command)
	fmt.Println("Document is : ", resp.DocumentData[0].Description)
}
