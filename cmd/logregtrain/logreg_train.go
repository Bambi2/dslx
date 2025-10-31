package main

import (
	"dslx/internal/hogwarts"
	"dslx/internal/logisticregression"
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: logreg_train <csv_file_path>")
		os.Exit(1)
	}

	csvFilePath := os.Args[1]

	dataset, err := hogwarts.LoadDatasetWithFeatures(csvFilePath, true, []string{
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

	model := logisticregression.TrainNewModel(dataset, 0.01, 1000)

	// Saving models to a file
	modelsJSON, err := json.Marshal(model)
	if err != nil {
		fmt.Println("Error marshalling model:", err)
		os.Exit(1)
	}

	err = os.WriteFile("model.json", modelsJSON, 0644)
	if err != nil {
		fmt.Println("Error saving model:", err)
		os.Exit(1)
	}
}
