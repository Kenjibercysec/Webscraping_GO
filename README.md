# Scraper Go

Este é um projeto Go simples para extrair informações de produtos de um site de comércio eletrônico e salvar os dados em arquivos CSV e JSON.

## Visão geral

O projeto consiste em vários arquivos Go:

- `main.go`: O ponto de entrada do aplicativo. Ele define a URL base, itera pelas páginas, extrai os produtos e salva os dados.
- `scraper.go`: Contém funções para obter o número total de páginas e extrair produtos de uma única página.
- `output.go`: Contém funções para salvar os dados extraídos em arquivos CSV e JSON.
- `models.go`: Define a estrutura de dados `Product`.
- `go.mod`: Define as dependências do módulo Go.

## Como usar

1.  **Pré-requisitos:**
    *   Go instalado (versão 1.24 ou superior)

2.  **Instalação:**

    ```bash
    git clone <repository_url>
    cd scrapergo
    go mod download
    ```

3.  **Execução:**

    ```bash
    go run main.go
    ```

    Isso irá extrair os dados do site configurado e salvar os resultados em `produtos.csv` e `produtos.json`.

## Adaptando para outros sites

Para usar este scraper com outros sites, você precisará ajustar os seletores CSS e a lógica de extração para corresponder à estrutura HTML do site de destino. Aqui estão os passos:

1.  **Inspecione o HTML do site de destino:**

    *   Use as ferramentas de desenvolvedor do seu navegador para examinar o HTML do site que você deseja extrair.
    *   Identifique os seletores CSS para os elementos que contêm o título do produto, preço, URL da imagem e URL do produto.

2.  **Atualize o arquivo `scraper.go`:**

    *   Modifique a função [`scrapePage`](scraper.go) para usar os seletores CSS corretos para o site de destino.
    *   Altere a lógica de extração conforme necessário para lidar com quaisquer formatos de dados específicos do site de destino.

    Exemplo de modificação da função [`scrapePage`](scraper.go):

    ```go
    // filepath: scraper.go
    // ...existing code...
    doc.Find("li.product").Each(func(i int, s *goquery.Selection) {
        title := strings.TrimSpace(s.Find("h2.product-title").Text()) // Alterado o seletor
        price := strings.TrimSpace(s.Find("span.item-price").Text())   // Alterado o seletor
        image, _ := s.Find("img").Attr("src")
        link, _ := s.Find("a").Attr("href")

        products = append(products, Product{
            Title:      title,
            Price:      price,
            ImageURL:   image,
            ProductURL: link,
        })
    })
    // ...existing code...
    ```

3.  **Atualize a função `getTotalPages` (se necessário):**

    *   Se o site de destino tiver uma estrutura de paginação diferente, você precisará atualizar a função [`getTotalPages`](scraper.go) para extrair o número total de páginas corretamente.

    Exemplo de modificação da função [`getTotalPages`](scraper.go):

    ```go
    // filepath: scraper.go
    // ...existing code...
    doc.Find("div.pagination a").Each(func(i int, s *goquery.Selection) { // Alterado o seletor
        text := s.Text()
        if n, err := strconv.Atoi(text); err == nil && n > lastPage {
            lastPage = n
        }
    })
    // ...existing code...
    ```

4.  **Atualize a URL base no arquivo `main.go`:**

    *   Altere a variável `baseURL` na função [`main`](main.go) para a URL base do site de destino.

    ```go
    // filepath: main.go
    // ...existing code...
    baseURL := "https://www.example.com/ecommerce/" // Alterado a URL base
    // ...existing code...
    ```

5.  **Execute o scraper:**

    ```bash
    go run main.go
    ```

## Considerações

*   **Respeite o arquivo `robots.txt`:** Verifique o arquivo `robots.txt` do site de destino para garantir que você tenha permissão para extrair os dados.
*   **Defina um tempo limite de requisição:** Para evitar sobrecarregar o servidor, defina um tempo limite de requisição entre as solicitações.
*   **Lidar com diferentes formatos de dados:** Esteja preparado para lidar com diferentes formatos de dados, como diferentes formatos de moeda ou diferentes estruturas HTML.
*   **Tratamento de erros:** Implemente um tratamento de erros robusto para lidar com quaisquer erros que possam ocorrer durante o processo de raspagem.

## Disclaimer

Raspagem de dados pode ser contra os termos de serviço de alguns sites. Certifique-se de ter permissão para extrair dados antes de executar este scraper. O autor deste projeto não é responsável por qualquer uso indevido deste scraper.