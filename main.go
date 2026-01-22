package main

import (
	"fmt"
	"log"
	"sync"
)

func main() {
	fmt.Print("Digite a URL do site para scraping: ")
	var url string
	fmt.Scanln(&url)
	
	if url == "" {
		url = "https://www.erome.com/explore"
		fmt.Printf("Usando URL padrão: %s\n", url)
	}
	
	log.Println("🚀 Iniciando análise do site...")
	
	// 1. Detectar estrutura da página
	structure, err := analyzePageStructure(url)
	if err != nil {
		log.Fatalf("❌ Erro ao analisar estrutura: %v", err)
	}
	
	// 2. Detectar padrão de paginação
	pattern := detectPagePattern(url)
	
	// 3. Detectar total de páginas
	totalPages := getTotalPages(url, pattern)
	
	// 4. Configurar scraping
	config := &ScrapingConfig{
		BaseURL:       url,
		PagePattern:   pattern,
		TotalPages:    totalPages,
		ContainerClass: structure.ContainerSelector,
		TitleClass:    structure.TitleSelector,
		ImageClass:    structure.ImageSelector,
		LinkClass:     structure.LinkSelector,
	}
	
	fmt.Printf("\n🔧 Configuração detectada:\n")
	fmt.Printf("   Padrão de paginação: %s\n", config.PagePattern)
	fmt.Printf("   Total de páginas: %d\n", config.TotalPages)
	fmt.Printf("   Container: %s\n", config.ContainerClass)
	fmt.Printf("   Título: %s\n", config.TitleClass)
	fmt.Printf("   Imagem: %s\n", config.ImageClass)
	fmt.Printf("   Link: %s\n", config.LinkClass)
	
	fmt.Printf("\n🎯 Iniciando scraping de %d páginas...\n\n", config.TotalPages)
	
	// Canal para produtos
	productsChan := make(chan Product, 100)
	var allProducts []Product
	
	// WaitGroup para sincronizar goroutines
	var wg sync.WaitGroup
	
	// Goroutine para coletar produtos
	go func() {
		for product := range productsChan {
			allProducts = append(allProducts, product)
		}
	}()
	
	// Scraping de múltiplas páginas
	for page := 1; page <= config.TotalPages; page++ {
		wg.Add(1)
		go func(pageNum int) {
			defer wg.Done()
			
			pageURL := buildPageURL(config.BaseURL, config.PagePattern, pageNum)
			log.Printf("📄 Página %d/%d: %s", pageNum, config.TotalPages, pageURL)
			
			count, err := scrapePage(pageURL, productsChan, structure)
			if err != nil {
				log.Printf("⚠️  Erro na página %d: %v", pageNum, err)
				return
			}
			
			log.Printf("✅ Página %d/%d | %d produtos coletados", pageNum, config.TotalPages, count)
		}(page)
	}
	
	// Aguardar todas as páginas
	wg.Wait()
	close(productsChan)
	
	log.Printf("\n🎉 Scraping concluído! Total de %d produtos coletados.\n", len(allProducts))
	
	// Salvar resultados
	if len(allProducts) > 0 {
		if err := saveToCSV(allProducts, "produtos.csv"); err != nil {
			log.Printf("❌ Erro ao salvar CSV: %v", err)
		}
		
		if err := saveToJSON(allProducts, "produtos.json"); err != nil {
			log.Printf("❌ Erro ao salvar JSON: %v", err)
		}
		
		// Mostrar amostra
		fmt.Println("\n📋 Amostra dos primeiros 5 produtos:")
		for i := 0; i < len(allProducts) && i < 5; i++ {
			fmt.Printf("%d. %s\n", i+1, allProducts[i].Title)
		}
	} else {
		fmt.Println("⚠️  Nenhum produto encontrado. Verifique:")
		fmt.Println("   - A URL está correta?")
		fmt.Println("   - O site requer JavaScript?")
		fmt.Println("   - Há proteção contra scraping?")
	}
}