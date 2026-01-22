package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/gocolly/colly"
)

func analyzePageStructure(url string) (*PageStructure, error) {
	log.Println("🔍 Analisando estrutura da página...")
	
	c := colly.NewCollector(
		colly.AllowedDomains(extractDomain(url)),
		colly.MaxDepth(1),
	)
	
	structure := &PageStructure{
		Selectors: make(map[string]string),
	}
	
	// Mapas para contar frequência de classes/IDs
	containerCandidates := make(map[string]int)
	titleCandidates := make(map[string]int)
	imageCandidates := make(map[string]int)
	linkCandidates := make(map[string]int)
	
	// Coletar todas as divs e suas classes
	c.OnHTML("div, article, section, li", func(e *colly.HTMLElement) {
		class := e.Attr("class")
		id := e.Attr("id")
		
		// Verificar se parece ser um container de item
		if containsKeywords(class, []string{"album", "item", "product", "card", "box", "container", "grid"}) {
			containerCandidates[class]++
		}
		if containsKeywords(id, []string{"album", "item", "product", "card"}) {
			containerCandidates["#"+id]++
		}
		
		// Verificar elementos internos
		e.ForEach("*", func(_ int, child *colly.HTMLElement) {
			childClass := child.Attr("class")
			childTag := child.Name
			
			// Títulos
			if childTag == "h1" || childTag == "h2" || childTag == "h3" || 
			   childTag == "h4" || childTag == "h5" || childTag == "h6" ||
			   containsKeywords(childClass, []string{"title", "name", "heading", "caption"}) {
				titleCandidates[childClass]++
			}
			
			// Imagens
			if childTag == "img" || containsKeywords(childClass, []string{"image", "img", "photo", "thumbnail", "picture"}) {
				imageCandidates[childClass]++
			}
			
			// Links
			if childTag == "a" || containsKeywords(childClass, []string{"link", "url", "href", "button"}) {
				linkCandidates[childClass]++
			}
		})
	})
	
	// Analisar padrões de links para paginação
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		href := e.Attr("href")
		if strings.Contains(href, "page=") || strings.Contains(href, "p=") || 
		   regexp.MustCompile(`/page/\d+`).MatchString(href) {
			log.Printf("Possível padrão de paginação encontrado: %s", href)
		}
	})
	
	c.OnRequest(func(r *colly.Request) {
		log.Printf("Analisando: %s", r.URL)
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	})
	
	err := c.Visit(url)
	if err != nil {
		return nil, err
	}
	
	// Selecionar os candidatos mais prováveis
	structure.ContainerSelector = selectBestCandidate(containerCandidates, "div")
	structure.TitleSelector = selectBestCandidate(titleCandidates, "h2, h3, .title, .name")
	structure.ImageSelector = selectBestCandidate(imageCandidates, "img, .image, .thumbnail")
	structure.LinkSelector = selectBestCandidate(linkCandidates, "a, .link")
	
	// Se não encontrou específico, usar padrões comuns
	if structure.ContainerSelector == "" {
		structure.ContainerSelector = findContainerByPattern(url)
	}
	
	log.Printf("✅ Estrutura detectada:")
	log.Printf("   Container: %s", structure.ContainerSelector)
	log.Printf("   Título: %s", structure.TitleSelector)
	log.Printf("   Imagem: %s", structure.ImageSelector)
	log.Printf("   Link: %s", structure.LinkSelector)
	
	return structure, nil
}

func containsKeywords(text string, keywords []string) bool {
	text = strings.ToLower(text)
	for _, keyword := range keywords {
		if strings.Contains(text, keyword) {
			return true
		}
	}
	return false
}

func selectBestCandidate(candidates map[string]int, fallback string) string {
	var bestClass string
	maxCount := 0
	
	for class, count := range candidates {
		if count > maxCount && class != "" {
			maxCount = count
			bestClass = class
		}
	}
	
	if bestClass == "" {
		return fallback
	}
	
	// Formatar como seletor CSS
	if strings.HasPrefix(bestClass, "#") {
		return bestClass
	}
	return "." + strings.Split(bestClass, " ")[0] // Pega apenas a primeira classe
}

func findContainerByPattern(url string) string {
	// Padrões conhecidos para sites comuns
	domain := extractDomain(url)
	
	switch {
	case strings.Contains(domain, "erome.com"):
		return "div.album"
	case strings.Contains(domain, "amazon") || strings.Contains(domain, "mercadolivre"):
		return "div.s-result-item, div.product"
	case strings.Contains(domain, "aliexpress"):
		return "div.product"
	case strings.Contains(domain, "magazineluiza"):
		return "div.product, li.product"
	default:
		return "div.product, article.item, li.product, div.card"
	}
}

func extractDomain(url string) string {
	url = strings.TrimPrefix(url, "http://")
	url = strings.TrimPrefix(url, "https://")
	url = strings.Split(url, "/")[0]
	return url
}

func detectPagePattern(url string) string {
	log.Println("🔍 Detectando padrão de paginação...")
	
	c := colly.NewCollector(
		colly.AllowedDomains(extractDomain(url)),
	)
	
	var patterns []string
	
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		href := e.Attr("href")
		
		// Verificar padrões comuns de paginação
		if strings.Contains(href, "page=") {
			patterns = append(patterns, "?page=%d")
		} else if strings.Contains(href, "p=") {
			patterns = append(patterns, "?p=%d")
		} else if match := regexp.MustCompile(`/page/(\d+)`).FindStringSubmatch(href); match != nil {
			patterns = append(patterns, "/page/%d")
		} else if match := regexp.MustCompile(`/p/(\d+)`).FindStringSubmatch(href); match != nil {
			patterns = append(patterns, "/p/%d")
		}
	})
	
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	})
	
	c.Visit(url)
	
	// Retornar o padrão mais comum
	if len(patterns) > 0 {
		// Contar frequência
		freq := make(map[string]int)
		for _, p := range patterns {
			freq[p]++
		}
		
		// Encontrar o mais frequente
		var bestPattern string
		maxFreq := 0
		for pattern, count := range freq {
			if count > maxFreq {
				maxFreq = count
				bestPattern = pattern
			}
		}
		
		log.Printf("✅ Padrão de paginação detectado: %s", bestPattern)
		return bestPattern
	}
	
	// Padrão padrão para erome.com
	if strings.Contains(url, "erome.com") {
		return "?page=%d"
	}
	
	log.Println("⚠️  Padrão de paginação não detectado, usando padrão genérico")
	return "?page=%d"
}

func getPagePatternFromUser(url string) string {
	// Auto-detecta, não precisa do usuário
	return detectPagePattern(url)
}

func getTotalPages(baseURL string, pattern string) int {
	log.Println("🔍 Detectando total de páginas...")
	
	if strings.Contains(baseURL, "erome.com") {
		return 50 // Valor padrão para erome.com
	}
	
	// Tentar detectar páginas visitando a primeira página
	c := colly.NewCollector(
		colly.AllowedDomains(extractDomain(baseURL)),
	)
	
	maxPage := 1
	
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		href := e.Attr("href")
		
		// Extrair número da página
		re := regexp.MustCompile(`[?&/](?:page|p)=(\d+)`)
		if matches := re.FindStringSubmatch(href); matches != nil {
			pageNum := matches[1]
			// Converter para int
			var num int
			fmt.Sscanf(pageNum, "%d", &num)
			if num > maxPage {
				maxPage = num
			}
		}
		
		// Outro padrão: /page/2
		re2 := regexp.MustCompile(`/page/(\d+)`)
		if matches := re2.FindStringSubmatch(href); matches != nil {
			pageNum := matches[1]
			var num int
			fmt.Sscanf(pageNum, "%d", &num)
			if num > maxPage {
				maxPage = num
			}
		}
	})
	
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	})
	
	c.Visit(baseURL)
	
	if maxPage > 1 {
		log.Printf("✅ Total de páginas detectado: %d", maxPage)
		return maxPage
	}
	
	log.Println("⚠️  Não foi possível detectar total de páginas, usando 10 como padrão")
	return 10
}