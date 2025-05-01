package main

import "fmt"

func main() {
	baseURL := "https://www.scrapingcourse.com/ecommerce/"
	totalPages := getTotalPages(baseURL)

	var allProducts []Product

	for i := 1; i <= totalPages; i++ {
		pageURL := fmt.Sprintf("https://www.scrapingcourse.com/ecommerce/page-%d.html", i)
		fmt.Println("Scraping:", pageURL)

		products, err := scrapePage(pageURL)
		if err != nil {
			fmt.Println("Erro:", err)
			continue
		}

		allProducts = append(allProducts, products...)
	}

	saveToCSV(allProducts, "produtos.csv")
	saveToJSON(allProducts, "produtos.json")

	fmt.Printf("\n✅ %d produtos salvos em produtos.csv e produtos.json\n", len(allProducts))
}
