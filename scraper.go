package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Estrutura para configuração do scraping
type ScrapingConfig struct {
	BaseURL     string
	PagePattern string // Padrão para as páginas (ex: "page-%d.html" ou "?page=%d")
	TotalPages  int
}

// Função para analisar a primeira página e detectar o padrão de paginação
func detectPagePattern(baseURL string) (string, error) {
	res, err := http.Get(baseURL)
	if err != nil {
		return "", fmt.Errorf("erro ao acessar a página: %v", err)
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return "", fmt.Errorf("erro ao processar HTML: %v", err)
	}

	// Procura por links de paginação
	var pagePattern string
	doc.Find("a[href*='page'], a[href*='p=']").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists {
			// Tenta identificar o padrão da URL
			if strings.Contains(href, "page-") {
				pagePattern = "page-%d.html"
			} else if strings.Contains(href, "?page=") {
				pagePattern = "?page=%d"
			} else if strings.Contains(href, "&p=") {
				pagePattern = "&p=%d"
			}
		}
	})

	if pagePattern == "" {
		return "", fmt.Errorf("não foi possível detectar o padrão de paginação automaticamente")
	}

	return pagePattern, nil
}

// Função para solicitar o padrão de paginação ao usuário
func getPagePatternFromUser() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Digite o padrão de paginação (use %d para o número da página, ex: 'page-%d.html' ou '?page=%d'): ")
	pattern, _ := reader.ReadString('\n')
	return strings.TrimSpace(pattern)
}

// Função para construir a URL da página
func buildPageURL(config ScrapingConfig, pageNum int) string {
	return fmt.Sprintf("%s%s", config.BaseURL, 
		fmt.Sprintf(config.PagePattern, pageNum))
}

func getTotalPages(baseURL string) int {
	res, err := http.Get(baseURL)
	if err != nil {
		log.Println("Erro ao acessar base:", err)
		return 1
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Println("Erro ao processar HTML:", err)
		return 1
	}

	lastPage := 1
	doc.Find("ul.pagination li a, .pagination a, .page-numbers a").Each(func(i int, s *goquery.Selection) {
		text := s.Text()
		if n, err := strconv.Atoi(text); err == nil && n > lastPage {
			lastPage = n
		}
	})

	return lastPage
}

func scrapePage(url string) ([]Product, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	var products []Product

	// Tenta diferentes seletores comuns para produtos
	selectors := []string{
		"li.product",
		".product",
		".item",
		".product-item",
	}

	for _, selector := range selectors {
		doc.Find(selector).Each(func(i int, s *goquery.Selection) {
			title := strings.TrimSpace(s.Find("h2, .title, .product-title").Text())
			price := strings.TrimSpace(s.Find(".price, .product-price").Text())
			image, _ := s.Find("img").Attr("src")
			link, _ := s.Find("a").Attr("href")

			if title != "" && price != "" {
				products = append(products, Product{
					Title:      title,
					Price:      price,
					ImageURL:   image,
					ProductURL: link,
				})
			}
		})

		if len(products) > 0 {
			break
		}
	}

	return products, nil
}
