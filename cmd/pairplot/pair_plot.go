package main

import (
	"dslx/internal/hogwarts"
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

var houseColors = map[string]color.Color{
	"Gryffindor": color.RGBA{R: 116, G: 0, B: 1, A: 248},
	"Hufflepuff": color.RGBA{R: 97, G: 75, B: 58, A: 248},
	"Ravenclaw":  color.RGBA{R: 14, G: 26, B: 64, A: 248},
	"Slytherin":  color.RGBA{R: 26, G: 71, B: 42, A: 248},
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: pairplot <csv_file_path>")
		os.Exit(1)
	}

	csvFilePath := os.Args[1]

	dataset, err := hogwarts.LoadDataset(csvFilePath, true)
	if err != nil {
		fmt.Println("Error loading dataset:", err)
		os.Exit(1)
	}

	numFeatures := len(dataset.FeatureNames)

	plotSize := 2.5 * vg.Inch
	imgWidth := vg.Length(numFeatures) * plotSize
	imgHeight := vg.Length(numFeatures) * plotSize

	img := vgimg.New(imgWidth, imgHeight)
	dc := draw.New(img)

	plots := make([][]*plot.Plot, numFeatures)
	for row := range numFeatures {
		plots[row] = make([]*plot.Plot, numFeatures)
		for col := range numFeatures {
			p := plot.New()

			if row == numFeatures-1 {
				p.X.Label.Text = dataset.FeatureNames[col]
			}
			if col == 0 {
				p.Y.Label.Text = dataset.FeatureNames[row]
			}

			if row == col {
				p.Title.Text = dataset.FeatureNames[row]

				for _, house := range dataset.Houses {
					values := dataset.GetFeatureValuesByHouse(col, house)

					filteredValues := make(plotter.Values, 0, len(values))
					for _, v := range values {
						if !math.IsNaN(v) {
							filteredValues = append(filteredValues, v)
						}
					}

					if len(filteredValues) > 0 {
						hist, err := plotter.NewHist(filteredValues, 20)
						if err != nil {
							fmt.Println("Error creating histogram:", err)
							os.Exit(1)
						}
						hist.FillColor = houseColors[house]
						hist.LineStyle.Width = vg.Points(0.5)
						hist.LineStyle.Color = color.RGBA{R: 0, G: 0, B: 0, A: 255}
						p.Add(hist)

						if row == 0 {
							p.Legend.Add(house, hist)
						}
					}
				}

				if row == 0 {
					p.Legend.Top = true
					p.Legend.Left = false
				}
			} else {
				for _, house := range dataset.Houses {
					xValues := dataset.GetFeatureValuesByHouse(col, house)
					yValues := dataset.GetFeatureValuesByHouse(row, house)

					xys := make(plotter.XYs, 0, len(xValues))
					for i := range xValues {
						if !math.IsNaN(xValues[i]) && !math.IsNaN(yValues[i]) {
							xys = append(xys, plotter.XY{X: xValues[i], Y: yValues[i]})
						}
					}

					if len(xys) > 0 {
						scatter, err := plotter.NewScatter(xys)
						if err != nil {
							fmt.Println("Error creating scatter:", err)
							os.Exit(1)
						}
						scatter.GlyphStyle.Color = houseColors[house]
						scatter.GlyphStyle.Radius = vg.Points(1.5)
						scatter.GlyphStyle.Shape = draw.CircleGlyph{}
						p.Add(scatter)
					}
				}
			}

			plots[row][col] = p
		}
	}

	t := draw.Tiles{
		Rows: numFeatures,
		Cols: numFeatures,
	}

	canvases := plot.Align(plots, t, dc)

	for row := range numFeatures {
		for col := range numFeatures {
			plots[row][col].Draw(canvases[row][col])
		}
	}

	w, err := os.Create("scatter_plot_matrix.png")
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

	fmt.Println("Scatter plot matrix saved to scatter_plot_matrix.png")
}
