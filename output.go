package main

import (
	"encoding/csv"
	"encoding/json"
	"log"
	"os"
)

func saveToCSV(products []Product, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Header atualizado
	header := []string{"Título", "Preço", "Imagem", "Link", "Visualizações", "Usuário", "Tipo de Mídia"}
	if err := writer.Write(header); err != nil {
		return err
	}

	for _, product := range products {
		record := []string{
			product.Title,
			product.Price,
			product.Image,
			product.Link,
			product.Views,
			product.User,
			product.MediaType,
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	log.Printf("✅ %d produtos salvos em %s", len(products), filename)
	return nil
}

func saveToJSON(products []Product, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(products); err != nil {
		return err
	}

	log.Printf("✅ %d produtos salvos em %s", len(products), filename)
	return nil
}