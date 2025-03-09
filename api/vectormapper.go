package api

import (
	"fmt"
	"math"
)

type VectorMapper struct {
}

func NewVectorMapper() *VectorMapper {
	return &VectorMapper{}
}

// EuclideanDistance calculates the Euclidean distance between two vectors
func (v *VectorMapper) EuclideanDistance(vec1, vec2 []float64) (float64, error) {
	if len(vec1) != len(vec2) {
		return 0, fmt.Errorf("vectors must be of the same length")
	}

	var sum float64
	for i := range vec1 {
		diff := vec1[i] - vec2[i]
		sum += diff * diff
	}

	return math.Sqrt(sum), nil
}

func (v *VectorMapper) GetMostSimilarVector(inputVector []float64, vectors [][]float64) ([]float64, error) {
	if len(vectors) == 0 {
		return nil, fmt.Errorf("no vectors provided")
	}

	var minDistance float64
	var mostSimilarVector []float64
	for _, vector := range vectors {
		distance, err := v.EuclideanDistance(inputVector, vector)
		if err != nil {
			return nil, err
		}

		if minDistance == 0 || distance < minDistance {
			minDistance = distance
			mostSimilarVector = vector
		}
	}

	return mostSimilarVector, nil
}
