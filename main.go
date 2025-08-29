package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/net/html"
)

// Extrai e concatena conteúdo das tags <title>, <h1>, <h2>, <p>
func extractHTMLContent(path string) (string, error) {
	tags := map[string]bool{
		"title": true,
		"h1":    true,
		"h2":    true,
		"p":     true,
	}

	var contents []string

	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	doc, err := html.Parse(file)
	if err != nil {
		return "", err
	}

	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode && tags[n.Data] {
			if n.FirstChild != nil {
				text := strings.TrimSpace(n.FirstChild.Data)
				if text != "" {
					contents = append(contents, text)
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}

	traverse(doc)
	return strings.Join(contents, ", "), nil
}

func main() {
	// Caminho da pasta que você quer ler
	pasta := "./aulas"

	// Verifica se a pasta existe
	info, err := os.Stat(pasta)
	if err != nil {
		fmt.Println("Erro ao acessar a pasta:", err)
		return
	}
	if !info.IsDir() {
		fmt.Println("O caminho especificado não é uma pasta.")
		return
	}

	// Mapa para armazenar conteúdo por pasta (usando ponteiros para Builder)
	conteudoPorPasta := make(map[string]*strings.Builder)

	// Variável para controlar a pasta atual
	pastaAtual := ""

	// Percorre os arquivos e subpastas
	err = filepath.WalkDir(pasta, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Println("Erro ao acessar:", path, err)
			return err
		}
		
		if d.IsDir() {
			fmt.Println("📁 Pasta:", path)
			pastaAtual = path
		} else {
			if strings.HasSuffix(strings.ToLower(d.Name()), ".html") {
				fmt.Println("✅ Arquivo HTML:", path)
				result, err := extractHTMLContent(path)
				if err != nil {
					fmt.Println("Erro ao extrair conteúdo:", err)
					return nil // Continua processando outros arquivos
				}
				fmt.Println("📝 Conteúdo extraído:", result)
				
				// Adiciona o conteúdo à pasta correspondente
				if builder, exists := conteudoPorPasta[pastaAtual]; exists {
					if builder.Len() > 0 {
						builder.WriteString(" | ")
					}
					builder.WriteString(result)
				} else {
					builder := &strings.Builder{}
					builder.WriteString(result)
					conteudoPorPasta[pastaAtual] = builder
				}
			} else {
				fmt.Println("📄 Outro arquivo:", path)
			}
		}
		return nil
	})

	if err != nil {
		fmt.Println("Erro ao percorrer a pasta:", err)
	}

	// Apresenta todo o conteúdo agrupado por pasta
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📚 CONTEÚDO EXTRAÍDO POR PASTA:")
	fmt.Println(strings.Repeat("=", 80))
	
	for pasta, conteudo := range conteudoPorPasta {
		fmt.Printf("\n📁 PASTA: %s\n", pasta)
		fmt.Printf("📝 CONTEÚDO: %s\n", conteudo.String())
		fmt.Println(strings.Repeat("-", 80))
	}
}