package api

import (
	"encoding/json"
	"miniagi/entities"
	"miniagi/utils"
	"sync"
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

func (a *AGI) GetClosetDocument(text string) (entities.VectorDocument, float64, error) {
	vector, err := a.ttv.ConvertTextToVector(text, a.embeddingModel)
	if err != nil {
		return entities.VectorDocument{}, 0, err
	}

	file, err := utils.ReadFile(a.outputDocumentFilePath)
	if err != nil {
		return entities.VectorDocument{}, 0, err
	}
	var vectorDocuments []entities.VectorDocument
	err = json.Unmarshal(file, &vectorDocuments)
	if err != nil {
		return entities.VectorDocument{}, 0, err
	}

	var minDistance float64
	ClosetDocument := entities.VectorDocument{}

	for _, vectorDocument := range vectorDocuments {
		distanceResult, _ := a.vtmp.EuclideanDistance(vector[0], vectorDocument.DocumentData.Vector)
		if minDistance == 0 || distanceResult < minDistance {
			minDistance = distanceResult
			ClosetDocument = vectorDocument
		}
	}

	return ClosetDocument, minDistance, nil
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
		distanceResult, _ := a.vtmp.EuclideanDistance(vector[0], vectorDocument.DocumentData.Vector)
		distanceMap[distanceResult] = vectorDocument
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

		vectors, err := a.ttv.ConvertTextToVector(rawDoc.Description, a.embeddingModel)
		if err != nil {
			return err
		}

		vactorDocument := entities.VectorDocument{
			Command: command,
			DocumentData: entities.DocumentData{
				Description:       rawDoc.Description,
				Vector:            vectors[0],
				EmbedingModelName: a.embeddingModel,
			},
		}

		vectorDocuments = append(vectorDocuments, vactorDocument)
	}

	data, err := json.Marshal(vectorDocuments)
	if err != nil {
		return err
	}

	return utils.WriteToJson(data, a.outputDocumentFilePath)
}

func (a *AGI) GetClosetDocumentMultiCal(text string, numWorker int) (entities.VectorDocument, float64, error) {
	vector, err := a.ttv.ConvertTextToVector(text, a.embeddingModel)
	if err != nil {
		return entities.VectorDocument{}, 0, err
	}

	file, err := utils.ReadFile(a.outputDocumentFilePath)
	if err != nil {
		return entities.VectorDocument{}, 0, err
	}
	var vectorDocuments []entities.VectorDocument
	err = json.Unmarshal(file, &vectorDocuments)
	if err != nil {
		return entities.VectorDocument{}, 0, err
	}

	calChan := make(chan handleCal)

	go func() {
		for position, vectorDocument := range vectorDocuments {
			c := handleCal{
				position:    position,
				inputVector: vector[0],
				input2List:  vectorDocument.DocumentData.Vector,
			}
			calChan <- c
		}
		close(calChan)
	}()

	position, distance := a.MultiCalDistanceResult(vector[0], calChan, numWorker)

	return vectorDocuments[position], distance, nil
}

type result struct {
	position int
	distance float64
}

type handleCal struct {
	position    int
	inputVector []float64
	input2List  []float64
}

func (a *AGI) MultiCalDistanceResult(inputVector []float64, calChan chan handleCal, numWorker int) (int, float64) {
	respCh := make(chan result)

	var wg sync.WaitGroup

	for i := 1; i <= numWorker; i++ {
		wg.Add(1)
		go a.handleCalVector(i, calChan, respCh, &wg)
	}

	go func() {
		wg.Wait()
		close(respCh)
	}()

	minimumDistance := 0.0
	position := 0

	for result := range respCh {
		if minimumDistance == 0 || result.distance < minimumDistance {
			minimumDistance = result.distance
			position = result.position
		}
	}

	return position, minimumDistance
}

func (a *AGI) handleCalVector(worker int, calChan chan handleCal, respCh chan result, wg *sync.WaitGroup) {
	defer wg.Done()
	// fmt.Println("Start Worker : ", worker)

	for c := range calChan {
		distance, _ := a.vtmp.EuclideanDistance(c.inputVector, c.input2List)
		res := result{
			position: c.position,
			distance: distance,
		}
		// fmt.Println("Worker : ", worker, " Position : ", c.position, " Distance : ", distance)
		respCh <- res
	}

}
