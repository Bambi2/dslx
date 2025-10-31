package main

import (
	"dslx/internal/hogwarts"
	"dslx/internal/logisticregression"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: logreg_predict <csv_file_path> <models_file_path>")
		os.Exit(1)
	}

	csvFilePath := os.Args[1]
	modelsFilePath := os.Args[2]

	dataset, err := hogwarts.LoadDatasetWithFeatures(csvFilePath, false, []string{
		"Astronomy",
		"Herbology",
		"Defense Against the Dark Arts",
		"Divination",
		"Muggle Studies",
		"Ancient Runes",
		"Charms",
	})
	if err != nil {
		fmt.Println("Error loading dataset:", err)
		os.Exit(1)
	}

	model, err := logisticregression.LoadModelFromFile(modelsFilePath)
	if err != nil {
		fmt.Println("Error loading model:", err)
		os.Exit(1)
	}

	predictions := model.Predict(dataset)

	outputFile, err := os.Create("houses.csv")
	if err != nil {
		fmt.Println("Error creating output file:", err)
		os.Exit(1)
	}
	defer outputFile.Close()

	fmt.Fprintln(outputFile, "Index,Hogwarts House")
	for i, prediction := range predictions {
		fmt.Fprintf(outputFile, "%d,%s\n", i, prediction)
	}
}
