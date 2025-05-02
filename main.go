package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	// Solicitar URL do usuário
	baseURL := getURLFromUser()
	
	// Verificar se a URL é válida
	if !strings.HasPrefix(baseURL, "http") {
		fmt.Println("URL inválida. Deve começar com http:// ou https://")
		return
	}

	// Detectar ou solicitar o padrão de paginação
	pagePattern, err := detectPagePattern(baseURL)
	if err != nil {
		fmt.Println("Não foi possível detectar o padrão de paginação automaticamente.")
		pagePattern = getPagePatternFromUser()
	}

	config := ScrapingConfig{
		BaseURL:     baseURL,
		PagePattern: pagePattern,
		TotalPages:  getTotalPages(baseURL),
	}

	fmt.Printf("Iniciando scraping de %d páginas...\n", config.TotalPages)
	fmt.Println("Padrão de paginação detectado:", config.PagePattern)

	var allProducts []Product

	for i := 1; i <= config.TotalPages; i++ {
		pageURL := buildPageURL(config, i)
		fmt.Printf("\nScraping página %d/%d: %s\n", i, config.TotalPages, pageURL)
		
		products, err := scrapePage(pageURL)
		if err != nil {
			fmt.Printf("Erro na página %d: %v\n", i, err)
			continue
		}

		allProducts = append(allProducts, products...)
		showProgress(i, config.TotalPages, len(allProducts))
	}

	// Salvar resultados
	saveToCSV(allProducts, "produtos.csv")
	saveToJSON(allProducts, "produtos.json")

	fmt.Printf("\n\n✅ %d produtos salvos em produtos.csv e produtos.json\n", len(allProducts))
}

// Função para solicitar a URL
func getURLFromUser() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Digite a URL do site para scraping: ")
	url, _ := reader.ReadString('\n')
	return strings.TrimSpace(url)
}

// Função para mostrar progresso
func showProgress(currentPage, totalPages, totalProducts int) {
	fmt.Printf("\rPágina %d/%d | Produtos coletados: %d", currentPage, totalPages, totalProducts)
}
