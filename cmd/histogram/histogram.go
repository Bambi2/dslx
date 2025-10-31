package main

import (
	"dslx/internal/hogwarts"
	"dslx/internal/stats"
	"fmt"
	"image/color"
	"os"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
)

var houseColors = map[string]color.RGBA{
	"Gryffindor": {R: 116, G: 0, B: 1, A: 248},
	"Hufflepuff": {R: 97, G: 75, B: 58, A: 248},
	"Ravenclaw":  {R: 14, G: 26, B: 64, A: 248},
	"Slytherin":  {R: 26, G: 71, B: 42, A: 248},
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: histogram <csv_file_path>")
		os.Exit(1)
	}

	csvFilePath := os.Args[1]

	dataset, err := hogwarts.LoadDataset(csvFilePath, true)
	if err != nil {
		fmt.Println("Error loading dataset:", err)
		os.Exit(1)
	}

	rows := 4
	cols := 4
	plots := make([][]*plot.Plot, rows)
	for i := range rows {
		plots[i] = make([]*plot.Plot, cols)
	}

	for i, featureName := range dataset.FeatureNames {
		row := i / cols
		col := i % cols

		p := plot.New()
		p.Title.Text = featureName
		p.X.Label.Text = "Score"
		p.Y.Label.Text = "Frequency"

		for _, house := range dataset.Houses {
			values := stats.RemoveMissingValues(dataset.GetFeatureValuesByHouse(i, house))

			if len(values) > 0 {
				h, err := plotter.NewHist(plotter.Values(values), 20)
				if err != nil {
					fmt.Printf("Error creating histogram for %s - %s: %v\n", featureName, house, err)
					continue
				}

				h.FillColor = houseColors[house]
				h.LineStyle.Width = vg.Points(0.5)
				h.LineStyle.Color = color.RGBA{R: 0, G: 0, B: 0, A: 255}

				p.Add(h)
				p.Legend.Add(house, h)
			}
		}

		p.Legend.Top = true
		p.Legend.Left = false
		p.Legend.XOffs = -0.5 * vg.Centimeter

		plots[row][col] = p
	}

	img := vgimg.New(20*vg.Inch, 20*vg.Inch)
	dc := draw.New(img)

	t := draw.Tiles{
		Rows: rows,
		Cols: cols,
	}

	canvases := plot.Align(plots, t, dc)
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			if plots[i][j] != nil {
				plots[i][j].Draw(canvases[i][j])
			}
		}
	}

	w, err := os.Create("histograms.png")
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
