package main

import (
	"encoding/csv"
	"encoding/json"
	"log"
	"os"
)

func saveToCSV(products []Product, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal("Erro ao criar arquivo CSV:", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"Título", "Preço", "Imagem", "Link"})

	for _, p := range products {
		writer.Write([]string{p.Title, p.Price, p.ImageURL, p.ProductURL})
	}
}

func saveToJSON(products []Product, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal("Erro ao criar arquivo JSON:", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(products); err != nil {
		log.Fatal("Erro ao codificar JSON:", err)
	}
}
