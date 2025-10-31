package main

import (
	"dslx/internal/hogwarts"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: describe <csv_file_path>")
		os.Exit(1)
	}

	csvFilePath := os.Args[1]

	dataset, err := hogwarts.LoadDataset(csvFilePath, true)
	if err != nil {
		fmt.Println("Error loading dataset:", err)
		os.Exit(1)
	}

	fmt.Println(dataset)
}
