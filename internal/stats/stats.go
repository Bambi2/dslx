package stats

import (
	"math"
	"sort"
)

func Sum(values []float64) float64 {
	sum := 0.0
	for _, value := range values {
		if math.IsNaN(value) {
			continue
		}

		sum += value
	}
	return sum
}

func Mean(values []float64) float64 {
	count := 0
	for _, value := range values {
		if math.IsNaN(value) {
			continue
		}

		count++
	}

	return Sum(values) / float64(count)
}

func Std(values []float64) float64 {
	mean := Mean(values)
	sum := 0.0
	count := 0
	for _, value := range values {
		if math.IsNaN(value) {
			continue
		}

		sum += math.Pow(value-mean, 2)
		count++
	}
	if count == 0 {
		return 0.0
	}
	return math.Sqrt(sum / float64(count))
}

func Min(values []float64) float64 {
	min := values[0]
	for _, value := range values {
		if math.IsNaN(value) {
			continue
		}

		if value < min {
			min = value
		}
	}
	return min
}

func Max(values []float64) float64 {
	max := values[0]
	for _, value := range values {
		if math.IsNaN(value) {
			continue
		}

		if value > max {
			max = value
		}
	}
	return max
}

func Q25(values []float64) float64 {
	return Percentile(values, 0.25)
}

func Q50(values []float64) float64 {
	return Percentile(values, 0.5)
}

func Q75(values []float64) float64 {
	return Percentile(values, 0.75)
}

func Percentile(values []float64, p float64) float64 {
	values = RemoveMissingValues(values)

	if p <= 0.0 {
		return values[0]
	}
	if p >= 1.0 {
		return values[len(values)-1]
	}

	if len(values) == 0 {
		return 0.0
	}

	if len(values) == 1 {
		return values[0]
	}

	sort.Float64s(values)

	exactIndex := p * float64(len(values)+1)
	lowerIndex := int(exactIndex)
	upperIndex := lowerIndex + 1
	if upperIndex-1 >= len(values) {
		return values[len(values)-1]
	}
	if lowerIndex <= 0 {
		return values[0]
	}
	if float64(lowerIndex) == exactIndex {
		return values[lowerIndex-1]
	}

	weight := exactIndex - float64(lowerIndex)
	return values[lowerIndex-1] + (values[upperIndex-1]-values[lowerIndex-1])*weight
}

func FillMissingValuesWithMean(values []float64) []float64 {
	mean := Mean(values)
	for i := range values {
		if math.IsNaN(values[i]) {
			values[i] = mean
		}
	}
	return values
}

func RemoveMissingValues(values []float64) []float64 {
	newValues := make([]float64, 0, len(values))
	for _, value := range values {
		if !math.IsNaN(value) {
			newValues = append(newValues, value)
		}
	}
	return newValues
}

func CalculateCorrelation(xValues []float64, yValues []float64) float64 {
	filteredXValues := make([]float64, 0, len(xValues))
	filteredYValues := make([]float64, 0, len(yValues))
	for i := range xValues {
		if math.IsNaN(xValues[i]) || math.IsNaN(yValues[i]) {
			continue
		}

		filteredXValues = append(filteredXValues, xValues[i])
		filteredYValues = append(filteredYValues, yValues[i])
	}

	xMean := Mean(filteredXValues)
	yMean := Mean(filteredYValues)

	covariance := 0.0
	for i := range filteredXValues {
		covariance += (filteredXValues[i] - xMean) * (filteredYValues[i] - yMean)
	}
	covariance = covariance / float64(len(filteredXValues))

	return covariance / (Std(filteredXValues) * Std(filteredYValues))
}
