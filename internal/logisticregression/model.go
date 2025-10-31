package logisticregression

import (
	"dslx/internal/hogwarts"
	"encoding/json"
	"fmt"
	"math"
	"os"
)

type Model struct {
	LabelNames []string    `json:"label_names"`
	Weights    [][]float64 `json:"weights"`
	Means      []float64   `json:"means"`
	Stds       []float64   `json:"stds"`
}

func TrainNewModel(dataset *hogwarts.Dataset, alhpha float64, iteractions int) *Model {
	x := fillMissingValues(dataset.Features, dataset.Means)
	normalizeFeatures(x, dataset.Means, dataset.Stds)
	addBiasTerm(x)

	weights := make([][]float64, 0, len(dataset.Houses))
	labelNames := make([]string, 0, len(dataset.Houses))
	for _, house := range dataset.Houses {
		labelNames = append(labelNames, house)
		y := make([]float64, 0, len(dataset.Labels))
		for _, label := range dataset.Labels {
			if label == house {
				y = append(y, 1.0)
			} else {
				y = append(y, 0.0)
			}
		}
		weights = append(weights, gradientDescent(x, y, alhpha, iteractions))
	}

	return &Model{
		LabelNames: labelNames,
		Weights:    weights,
		Means:      dataset.Means,
		Stds:       dataset.Stds,
	}
}

func LoadModelFromFile(filePath string) (*Model, error) {
	modelsJSON, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var model *Model
	err = json.Unmarshal(modelsJSON, &model)
	if err != nil {
		return nil, err
	}

	return model, nil
}

func (m *Model) Predict(dataset *hogwarts.Dataset) []string {
	x := fillMissingValues(dataset.Features, m.Means)
	normalizeFeatures(x, m.Means, m.Stds)
	addBiasTerm(x)

	labelNames := make([]string, 0, len(x))
	for i := range x {
		maxPrediction := 0.0
		maxPredictionLabel := ""
		for j := range m.Weights {
			prediction := predict(x[i], m.Weights[j])
			if prediction > maxPrediction {
				maxPrediction = prediction
				maxPredictionLabel = m.LabelNames[j]
			}
		}
		labelNames = append(labelNames, maxPredictionLabel)
	}

	return labelNames
}

func gradientDescent(x [][]float64, y []float64, alhpha float64, iteractions int) []float64 {
	weights := make([]float64, len(x[0]))
	for iter := 0; iter < iteractions; iter++ {
		gradients := make([]float64, len(weights))
		for i := range len(y) {
			prediction := predict(x[i], weights)
			for j := range len(weights) {
				gradients[j] += (prediction - y[i]) * x[i][j]
			}
		}
		for j := range weights {
			gradients[j] /= float64(len(y))
			weights[j] -= alhpha * gradients[j]
		}

		if iter%100 == 0 {
			cost := computeCost(x, y, weights)
			fmt.Printf("  Iteration %d: Cost = %.6f\n", iter, cost)
		}
	}
	return weights
}

func computeCost(x [][]float64, y []float64, weights []float64) float64 {
	m := float64(len(y))
	cost := 0.0
	epsilon := 1e-15

	for i := 0; i < len(y); i++ {
		h := predict(x[i], weights)
		h = math.Max(epsilon, math.Min(1.0-epsilon, h))
		cost += -y[i]*math.Log(h) - (1.0-y[i])*math.Log(1.0-h)
	}

	return cost / m
}

func predict(x []float64, weights []float64) float64 {
	z := 0.0
	for i := range len(x) {
		z += x[i] * weights[i]
	}
	return sigmoid(z)
}

func sigmoid(z float64) float64 {
	return 1.0 / (1.0 + math.Exp(-z))
}

func fillMissingValues(rows [][]float64, means []float64) [][]float64 {
	x := make([][]float64, len(rows))
	for i := range rows {
		x[i] = make([]float64, len(rows[i]))
		for j := range rows[i] {
			if math.IsNaN(rows[i][j]) {
				x[i][j] = means[j]
			} else {
				x[i][j] = rows[i][j]
			}
		}
	}

	return x
}

func normalizeFeatures(x [][]float64, means []float64, stds []float64) {
	for i := range x {
		for j := range x[i] {
			if stds[j] < 1e-10 {
				x[i][j] = 0.0
			} else {
				x[i][j] = (x[i][j] - means[j]) / stds[j]
			}
		}
	}
}

func addBiasTerm(x [][]float64) {
	for i := range x {
		x[i] = append([]float64{1.0}, x[i]...)
	}
}
