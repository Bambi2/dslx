package main

import (
	"dslx/internal/hogwarts"
	"dslx/internal/stats"
	"fmt"
	"image/color"
	"math"
	"os"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: scatterplot <csv_file_path>")
		os.Exit(1)
	}

	csvFilePath := os.Args[1]

	dataset, err := hogwarts.LoadDataset(csvFilePath, true)
	if err != nil {
		fmt.Println("Error loading dataset:", err)
		os.Exit(1)
	}

	maxAbsoluteCorrelationFirstFeatureIndex := 0
	maxAbsoluteCorrelationSecondFeatureIndex := 0
	maxAbsoluteCorrelation := 0.0
	minAbsoluteCorrelation := 1.0
	minAbsoluteCorrelationFirstFeatureIndex := 0
	minAbsoluteCorrelationSecondFeatureIndex := 0
	for featureIndex := range dataset.FeatureNames {
		for otherFeatureIndex := featureIndex + 1; otherFeatureIndex < len(dataset.FeatureNames); otherFeatureIndex++ {
			correlation := stats.CalculateCorrelation(dataset.GetFeatureValues(featureIndex), dataset.GetFeatureValues(otherFeatureIndex))
			absoluteCorrelation := math.Abs(correlation)
			if absoluteCorrelation > maxAbsoluteCorrelation {
				maxAbsoluteCorrelation = absoluteCorrelation
				maxAbsoluteCorrelationFirstFeatureIndex = featureIndex
				maxAbsoluteCorrelationSecondFeatureIndex = otherFeatureIndex
			}
			if absoluteCorrelation < minAbsoluteCorrelation {
				minAbsoluteCorrelation = absoluteCorrelation
				minAbsoluteCorrelationFirstFeatureIndex = featureIndex
				minAbsoluteCorrelationSecondFeatureIndex = otherFeatureIndex
			}
		}
	}

	maxFirstValues := dataset.GetFeatureValues(maxAbsoluteCorrelationFirstFeatureIndex)
	maxSecondValues := dataset.GetFeatureValues(maxAbsoluteCorrelationSecondFeatureIndex)
	maxXYs := make(plotter.XYs, 0, len(maxFirstValues))
	for i := range maxFirstValues {
		if math.IsNaN(maxFirstValues[i]) || math.IsNaN(maxSecondValues[i]) {
			continue
		}

		maxXYs = append(maxXYs, plotter.XY{
			X: maxFirstValues[i],
			Y: maxSecondValues[i],
		})
	}
	maxCorrelationScatter, err := plotter.NewScatter(maxXYs)
	if err != nil {
		fmt.Println("Error creating max correlation scatter:", err)
		os.Exit(1)
	}
	maxCorrelationScatter.GlyphStyle.Color = color.RGBA{R: 255, B: 128, A: 255}
	highestPlot := plot.New()
	highestPlot.Title.Text = "Highest correlation"
	highestPlot.X.Label.Text = dataset.FeatureNames[maxAbsoluteCorrelationFirstFeatureIndex]
	highestPlot.Y.Label.Text = dataset.FeatureNames[maxAbsoluteCorrelationSecondFeatureIndex]
	highestPlot.Add(plotter.NewGrid())
	highestPlot.Add(maxCorrelationScatter)

	minFirstValues := dataset.GetFeatureValues(minAbsoluteCorrelationFirstFeatureIndex)
	minSecondValues := dataset.GetFeatureValues(minAbsoluteCorrelationSecondFeatureIndex)
	minXYs := make(plotter.XYs, 0, len(minFirstValues))
	for i := range minFirstValues {
		if math.IsNaN(minFirstValues[i]) || math.IsNaN(minSecondValues[i]) {
			continue
		}

		minXYs = append(minXYs, plotter.XY{
			X: minFirstValues[i],
			Y: minSecondValues[i],
		})
	}
	minCorrelationScatter, err := plotter.NewScatter(minXYs)
	if err != nil {
		fmt.Println("Error creating min correlation scatter:", err)
		os.Exit(1)
	}
	minCorrelationScatter.GlyphStyle.Color = color.RGBA{R: 128, B: 255, A: 255}
	lowestPlot := plot.New()
	lowestPlot.Title.Text = "Lowest correlation"
	lowestPlot.X.Label.Text = dataset.FeatureNames[minAbsoluteCorrelationFirstFeatureIndex]
	lowestPlot.Y.Label.Text = dataset.FeatureNames[minAbsoluteCorrelationSecondFeatureIndex]
	lowestPlot.Add(plotter.NewGrid())
	lowestPlot.Add(minCorrelationScatter)

	img := vgimg.New(4*vg.Inch, 10*vg.Inch)
	dc := draw.New(img)

	t := draw.Tiles{
		Rows: 2,
		Cols: 1,
	}

	canvases := plot.Align([][]*plot.Plot{
		{highestPlot},
		{lowestPlot},
	}, t, dc)

	highestPlot.Draw(canvases[0][0])
	lowestPlot.Draw(canvases[1][0])

	w, err := os.Create("scatterplot.png")
	if err != nil {
		fmt.Println("Error creating output file:", err)
		os.Exit(1)
	}
	defer w.Close()

	png := vgimg.PngCanvas{Canvas: img}
	if _, err := png.WriteTo(w); err != nil {
		fmt.Println("Error writing PNG:", err)
		os.Exit(1)
	}
}
