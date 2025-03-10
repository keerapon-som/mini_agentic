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

	// resp, avgDistance, err := agi.GetClosetDocument("I would like to go to other directory")
	// if err != nil {
	// 	panic(err)
	// }
	resp, avgDistance, err := agi.GetClosetDocumentMultiCal("I would like to go to other directory", 10)
	if err != nil {
		panic(err)
	}

	fmt.Println("Command is : ", resp.Command)
	fmt.Println("Document is : ", resp.DocumentData.Description)
	fmt.Println("Distance is : ", avgDistance)
}

// 0.8397428422016603
