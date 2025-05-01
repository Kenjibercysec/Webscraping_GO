package main

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

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
	doc.Find("ul.pagination li a").Each(func(i int, s *goquery.Selection) {
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

	doc.Find("li.product").Each(func(i int, s *goquery.Selection) {
		title := strings.TrimSpace(s.Find("h2.woocommerce-loop-product__title").Text())
		price := strings.TrimSpace(s.Find("span.price").Text())
		image, _ := s.Find("img").Attr("src")
		link, _ := s.Find("a").Attr("href")

		products = append(products, Product{
			Title:      title,
			Price:      price,
			ImageURL:   image,
			ProductURL: link,
		})
	})

	return products, nil
}
