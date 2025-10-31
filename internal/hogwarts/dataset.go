package hogwarts

import (
	"dslx/internal/stats"
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

const (
	featureStartIndex = 6
	featureEndIndex   = 19
	numberOfFeatures  = featureEndIndex - featureStartIndex
	houseIndex        = 1
)

type Dataset struct {
	Features     [][]float64
	Labels       []string
	Houses       []string
	FeatureNames []string
	Counts       []float64
	Means        []float64
	Stds         []float64
	Mins         []float64
	Maxs         []float64
	Q25s         []float64
	Q50s         []float64
	Q75s         []float64
}

func LoadDataset(filename string, skipEmptyHouses bool) (*Dataset, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	header := records[0]

	featureNames := header[featureStartIndex:featureEndIndex]

	features := make([][]float64, 0, len(records)-1)
	labels := make([]string, 0, len(records)-1)
	housesMap := make(map[string]struct{})

	for i := 1; i < len(records); i++ {
		row := records[i]
		house := row[houseIndex]

		if house == "" && skipEmptyHouses {
			continue
		}

		labels = append(labels, house)
		housesMap[house] = struct{}{}

		featureRow := make([]float64, numberOfFeatures)
		for j := range numberOfFeatures {
			featureStr := row[featureStartIndex+j]
			if featureStr == "" || featureStr == " " {
				featureRow[j] = math.NaN()
			} else {
				featureFloat, err := strconv.ParseFloat(strings.TrimSpace(featureStr), 64)
				if err != nil {
					featureRow[j] = math.NaN()
				} else {
					featureRow[j] = featureFloat
				}
			}
		}
		features = append(features, featureRow)
	}

	houses := make([]string, 0, len(housesMap))
	for house := range housesMap {
		houses = append(houses, house)
	}

	counts := make([]float64, numberOfFeatures)
	for j := range numberOfFeatures {
		count := 0.0
		for i := range features {
			if !math.IsNaN(features[i][j]) {
				count++
			}
		}
		counts[j] = count
	}

	means := make([]float64, numberOfFeatures)
	stds := make([]float64, numberOfFeatures)
	mins := make([]float64, numberOfFeatures)
	maxs := make([]float64, numberOfFeatures)
	q25s := make([]float64, numberOfFeatures)
	q50s := make([]float64, numberOfFeatures)
	q75s := make([]float64, numberOfFeatures)

	for i := range numberOfFeatures {
		means[i] = stats.Mean(getFeaureValues(i, features))
		stds[i] = stats.Std(getFeaureValues(i, features))
		if stds[i] < 1e-10 {
			stds[i] = 1.0
		}
		mins[i] = stats.Min(getFeaureValues(i, features))
		maxs[i] = stats.Max(getFeaureValues(i, features))
		q25s[i] = stats.Q25(getFeaureValues(i, features))
		q50s[i] = stats.Q50(getFeaureValues(i, features))
		q75s[i] = stats.Q75(getFeaureValues(i, features))
	}

	return &Dataset{
		Features:     features,
		Labels:       labels,
		Houses:       houses,
		FeatureNames: featureNames,
		Counts:       counts,
		Means:        means,
		Stds:         stds,
		Mins:         mins,
		Maxs:         maxs,
		Q25s:         q25s,
		Q50s:         q50s,
		Q75s:         q75s,
	}, nil
}

func (d *Dataset) String() string {
	const statColumnWidth = 15
	const minFeatureColumnWidth = 18
	const terminalWidth = 120

	featuresPerChunk := (terminalWidth - statColumnWidth) / minFeatureColumnWidth
	if featuresPerChunk < 1 {
		featuresPerChunk = 1
	}

	var result strings.Builder

	for chunkStart := 0; chunkStart < len(d.FeatureNames); chunkStart += featuresPerChunk {
		chunkEnd := chunkStart + featuresPerChunk
		if chunkEnd > len(d.FeatureNames) {
			chunkEnd = len(d.FeatureNames)
		}

		featureColumnWidth := minFeatureColumnWidth
		for i := chunkStart; i < chunkEnd; i++ {
			nameLen := len(d.FeatureNames[i])
			if nameLen > featureColumnWidth {
				featureColumnWidth = nameLen
			}
		}

		if chunkStart > 0 {
			result.WriteString("\n")
		}

		result.WriteString(fmt.Sprintf("%-*s", statColumnWidth, ""))
		for i := chunkStart; i < chunkEnd; i++ {
			result.WriteString(fmt.Sprintf("%-*s ", featureColumnWidth, d.FeatureNames[i]))
		}
		result.WriteString("\n")

		result.WriteString(fmt.Sprintf("%-*s", statColumnWidth, "Count"))
		for i := chunkStart; i < chunkEnd; i++ {
			result.WriteString(fmt.Sprintf("%-*.6f ", featureColumnWidth, d.Counts[i]))
		}
		result.WriteString("\n")

		result.WriteString(fmt.Sprintf("%-*s", statColumnWidth, "Mean"))
		for i := chunkStart; i < chunkEnd; i++ {
			result.WriteString(fmt.Sprintf("%-*.6f ", featureColumnWidth, d.Means[i]))
		}
		result.WriteString("\n")

		result.WriteString(fmt.Sprintf("%-*s", statColumnWidth, "Std"))
		for i := chunkStart; i < chunkEnd; i++ {
			result.WriteString(fmt.Sprintf("%-*.6f ", featureColumnWidth, d.Stds[i]))
		}
		result.WriteString("\n")

		result.WriteString(fmt.Sprintf("%-*s", statColumnWidth, "Min"))
		for i := chunkStart; i < chunkEnd; i++ {
			result.WriteString(fmt.Sprintf("%-*.6f ", featureColumnWidth, d.Mins[i]))
		}
		result.WriteString("\n")

		result.WriteString(fmt.Sprintf("%-*s", statColumnWidth, "25%"))
		for i := chunkStart; i < chunkEnd; i++ {
			result.WriteString(fmt.Sprintf("%-*.6f ", featureColumnWidth, d.Q25s[i]))
		}
		result.WriteString("\n")

		result.WriteString(fmt.Sprintf("%-*s", statColumnWidth, "50%"))
		for i := chunkStart; i < chunkEnd; i++ {
			result.WriteString(fmt.Sprintf("%-*.6f ", featureColumnWidth, d.Q50s[i]))
		}
		result.WriteString("\n")

		result.WriteString(fmt.Sprintf("%-*s", statColumnWidth, "75%"))
		for i := chunkStart; i < chunkEnd; i++ {
			result.WriteString(fmt.Sprintf("%-*.6f ", featureColumnWidth, d.Q75s[i]))
		}
		result.WriteString("\n")

		result.WriteString(fmt.Sprintf("%-*s", statColumnWidth, "Max"))
		for i := chunkStart; i < chunkEnd; i++ {
			result.WriteString(fmt.Sprintf("%-*.6f ", featureColumnWidth, d.Maxs[i]))
		}
		result.WriteString("\n")
	}

	return result.String()
}

func (d *Dataset) GetFeatureValuesByHouse(featureIndex int, house string) []float64 {
	values := make([]float64, 0)
	for i, label := range d.Labels {
		if label == house {
			values = append(values, d.Features[i][featureIndex])
		}
	}
	return values
}

func (d *Dataset) GetFeatureValues(featureIndex int) []float64 {
	return getFeaureValues(featureIndex, d.Features)
}

func getFeaureValues(featureIndex int, data [][]float64) []float64 {
	values := make([]float64, 0)
	for _, dataRow := range data {
		values = append(values, dataRow[featureIndex])
	}

	return values
}
