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

// Verifica se é HTML e extrai conteúdo
func checkHTML(file fs.DirEntry, path string) {
    if strings.HasSuffix(strings.ToLower(file.Name()), ".html") {
        fmt.Println("✅ Arquivo HTML:", path)
        result, err := extractHTMLContent(path)
        if err != nil {
            fmt.Println("Erro ao extrair conteúdo:", err)
            return
        }
        fmt.Println("📝 Conteúdo extraído:", result)
    } else {
        fmt.Println("📄 Outro arquivo:", path)
    }
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

    // Percorre os arquivos e subpastas
    err = filepath.WalkDir(pasta, func(path string, d fs.DirEntry, err error) error {
        if err != nil {
            fmt.Println("Erro ao acessar:", path, err)
            return err
        }
        if d.IsDir() {
            fmt.Println("📁 Pasta:", path)
        } else {
			checkHTML(d, path)
        }
        return nil
    })

    if err != nil {
        fmt.Println("Erro ao percorrer a pasta:", err)
    }
}
