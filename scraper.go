package main

import (
	"fmt"
	"log"
	"regexp"  // ADICIONE ESTA LINHA
	"strings"
	"time"

	"github.com/gocolly/colly"
)

func scrapePage(url string, products chan<- Product, structure *PageStructure) (int, error) {
	c := colly.NewCollector(
		colly.AllowedDomains(extractDomain(url)),
		colly.Async(true),
	)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 2,
		Delay:       2 * time.Second,
	})

	c.SetRequestTimeout(30 * time.Second)

	productsCount := 0

	c.OnRequest(func(r *colly.Request) {
		log.Printf("📄 Visitando: %s", r.URL)
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Printf("❌ Erro em %s: %v", r.Request.URL, err)
	})

	// Usar os seletores detectados
	if structure.ContainerSelector != "" {
		c.OnHTML(structure.ContainerSelector, func(e *colly.HTMLElement) {
			product := extractProduct(e, structure)
			if product.Title != "" && product.Link != "" {
				products <- product
				productsCount++
				log.Printf("✅ Coletado: %s", product.Title)
			}
		})
	} else {
		// Fallback: tentar vários seletores comuns
		commonSelectors := []string{
			"div.album", "div.product", "div.item", "article", 
			"div.card", "li.product", "div.s-result-item",
		}
		
		for _, selector := range commonSelectors {
			c.OnHTML(selector, func(e *colly.HTMLElement) {
				product := extractProduct(e, structure)
				if product.Title != "" && product.Link != "" {
					products <- product
					productsCount++
					log.Printf("✅ Coletado: %s", product.Title)
				}
			})
		}
	}

	err := c.Visit(url)
	if err != nil {
		return 0, fmt.Errorf("erro visitando %s: %v", url, err)
	}

	c.Wait()
	return productsCount, nil
}

func extractProduct(e *colly.HTMLElement, structure *PageStructure) Product {
	product := Product{}
	
	// Extrair título
	if structure.TitleSelector != "" {
		product.Title = strings.TrimSpace(e.ChildText(structure.TitleSelector))
	}
	if product.Title == "" {
		// Tentar alternativas
		product.Title = strings.TrimSpace(e.ChildText("h1, h2, h3, .title, .name, [class*='title'], [class*='name']"))
		product.Title = e.ChildAttr("img", "alt")
	}
	
	// Extrair imagem
	if structure.ImageSelector != "" {
		product.Image = e.ChildAttr(structure.ImageSelector, "src")
		if product.Image == "" {
			product.Image = e.ChildAttr(structure.ImageSelector, "data-src")
		}
	}
	if product.Image == "" {
		product.Image = e.ChildAttr("img", "src")
		if product.Image == "" {
			product.Image = e.ChildAttr("img", "data-src")
		}
	}
	
	// Extrair link
	if structure.LinkSelector != "" {
		product.Link = e.ChildAttr(structure.LinkSelector, "href")
	}
	if product.Link == "" {
		product.Link = e.ChildAttr("a", "href")
	}
	
	// Extrair informações extras baseadas no domínio
	domain := extractDomain(e.Request.URL.String())
	
	switch {
	case strings.Contains(domain, "erome.com"):
		product.Views = e.ChildText("span.album-bottom-views")
		product.User = e.ChildText("span.album-user")
		product.MediaType = "mixed"
		if e.ChildText("span.album-images") != "" {
			product.MediaType = "images"
		} else if e.ChildText("span.album-videos") != "" {
			product.MediaType = "videos"
		}
		product.Price = "N/A"
		
	case strings.Contains(domain, "amazon") || 
	     strings.Contains(domain, "mercadolivre") || 
	     strings.Contains(domain, "aliexpress"):
		product.Price = e.ChildText(".price, [class*='price'], .prc, .value")
		product.Views = e.ChildText(".views, .sales, .rating")
		
	default:
		// Tentar encontrar preço genérico
		product.Price = findPrice(e.Text)
	}
	
	// Tornar URLs absolutas
	if product.Image != "" && !strings.HasPrefix(product.Image, "http") {
		product.Image = e.Request.AbsoluteURL(product.Image)
	}
	if product.Link != "" && !strings.HasPrefix(product.Link, "http") {
		product.Link = e.Request.AbsoluteURL(product.Link)
	}
	
	return product
}

func findPrice(text string) string {
	// Regex para encontrar preços (R$, $, €, etc.)
	re := regexp.MustCompile(`(R\$\s*[\d.,]+|\$\s*[\d.,]+|€\s*[\d.,]+|[\d.,]+\s*(reais|dollars|euros))`)
	if match := re.FindString(text); match != "" {
		return strings.TrimSpace(match)
	}
	return "N/A"
}

func buildPageURL(baseURL string, pattern string, page int) string {
	if pattern == "" {
		return baseURL
	}
	
	if strings.Contains(pattern, "?") && strings.Contains(baseURL, "?") {
		// Remover query string existente
		baseURL = strings.Split(baseURL, "?")[0]
	}
	
	return baseURL + fmt.Sprintf(pattern, page)
}