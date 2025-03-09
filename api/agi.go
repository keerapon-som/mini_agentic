package api

import (
	"encoding/json"
	"miniagi/entities"
	"miniagi/utils"
)

type AGI struct {
	ttv                    *TextToVector
	vtmp                   *VectorMapper
	rawDocumentFilePath    string
	outputDocumentFilePath string
	embeddingModel         string
}

func NewAGI(ttv *TextToVector, vtmp *VectorMapper, rawDocumentFilePath string, outputDocumentFilePath string, embeddingModel string) *AGI {
	return &AGI{
		ttv:                    ttv,
		vtmp:                   vtmp,
		rawDocumentFilePath:    rawDocumentFilePath,
		outputDocumentFilePath: outputDocumentFilePath,
		embeddingModel:         embeddingModel,
	}
}

func (a *AGI) GetClosetDocument(text string) (entities.VectorDocument, error) {
	vector, err := a.ttv.ConvertTextToVector(text, a.embeddingModel)
	if err != nil {
		return entities.VectorDocument{}, err
	}

	file, err := utils.ReadFile(a.outputDocumentFilePath)
	if err != nil {
		return entities.VectorDocument{}, err
	}
	var vectorDocuments []entities.VectorDocument
	err = json.Unmarshal(file, &vectorDocuments)
	if err != nil {
		return entities.VectorDocument{}, err
	}

	var minDistance float64
	ClosetDocument := entities.VectorDocument{}

	for _, vectorDocument := range vectorDocuments {
		for _, docData := range vectorDocument.DocumentData {
			distanceResult, _ := a.vtmp.EuclideanDistance(vector[0], docData.Vector)
			if minDistance == 0 || distanceResult < minDistance {
				minDistance = distanceResult
				ClosetDocument = vectorDocument
			}
		}
	}

	return ClosetDocument, nil
}

func (a *AGI) GetNearDocuments(text string, limit int) ([]entities.VectorDocument, error) {
	vector, err := a.ttv.ConvertTextToVector(text, a.embeddingModel)
	if err != nil {
		return []entities.VectorDocument{}, err
	}

	file, err := utils.ReadFile(a.outputDocumentFilePath)
	if err != nil {
		return []entities.VectorDocument{}, err
	}
	var vectorDocuments []entities.VectorDocument
	err = json.Unmarshal(file, &vectorDocuments)
	if err != nil {
		return []entities.VectorDocument{}, err
	}

	var distanceMap = map[float64]entities.VectorDocument{}
	for _, vectorDocument := range vectorDocuments {
		for _, docData := range vectorDocument.DocumentData {
			distanceResult, _ := a.vtmp.EuclideanDistance(vector[0], docData.Vector)
			distanceMap[distanceResult] = vectorDocument
		}
	}

	var distances []float64
	for distance := range distanceMap {
		distances = append(distances, distance)
	}

	var sortedDistances []float64
	for i := 0; i < limit; i++ {
		var minDistance float64
		for _, distance := range distances {
			if minDistance == 0 || distance < minDistance {
				minDistance = distance
			}
		}
		sortedDistances = append(sortedDistances, minDistance)
		delete(distanceMap, minDistance)
	}

	var nearDocuments []entities.VectorDocument
	for _, distance := range sortedDistances {
		nearDocuments = append(nearDocuments, distanceMap[distance])
	}

	return nearDocuments, nil
}

func (a *AGI) TurnRawToVectorDocument() error {
	resp, err := utils.ReadFile(a.rawDocumentFilePath)
	if err != nil {
		return err
	}

	var rawDocs []entities.RawDocument
	err = json.Unmarshal(resp, &rawDocs)
	if err != nil {
		return err
	}

	vectorDocuments := []entities.VectorDocument{}

	for _, rawDoc := range rawDocs {
		command := rawDoc.Command
		description := rawDoc.Description
		// listVectors := [][]float64{}
		DocumentData := []entities.DocumentData{}
		for _, desc := range description {
			vectors, err := a.ttv.ConvertTextToVector(desc, a.embeddingModel)
			if err != nil {
				return err
			}
			DocumentData = append(DocumentData, entities.DocumentData{Description: desc, Vector: vectors[0], EmbedingModelName: a.embeddingModel})
		}
		vactorDocument := entities.VectorDocument{
			Command:      command,
			DocumentData: DocumentData,
		}
		vectorDocuments = append(vectorDocuments, vactorDocument)
	}

	data, err := json.Marshal(vectorDocuments)
	if err != nil {
		return err
	}

	return utils.WriteToJson(data, a.outputDocumentFilePath)
}
