package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/net/html"
)

// Tags HTML que serão extraídas
var tagsExtracao = map[string]bool{
	"title": true,
	"h1":    true,
	"h2":    true,
	"p":     true,
}

func main() {
	pasta := "./aulas"
	
	// Verifica se a pasta existe
	if info, err := os.Stat(pasta); err != nil || !info.IsDir() {
		fmt.Println("Erro: Pasta não encontrada ou inválida")
		return
	}

	// Processa a pasta e mostra resultados
	conteudoPorPasta := make(map[string]string)
	pastaAtual := ""

	filepath.WalkDir(pasta, func(caminho string, entrada fs.DirEntry, err error) error {
		if err != nil {
			return nil // Ignora erros e continua
		}

		if entrada.IsDir() {
			pastaAtual = caminho
			//fmt.Println("📁 Pasta:", caminho)
		} else if strings.HasSuffix(entrada.Name(), ".html") {
			processarArquivoHTML(caminho, pastaAtual, conteudoPorPasta)
		} else {
			fmt.Println("📄 Outro arquivo:", caminho)
		}
		return nil
	})

	mostrarResultados(conteudoPorPasta)
}

// Processa um arquivo HTML e adiciona ao conteúdo da pasta
func processarArquivoHTML(caminho, pasta string, conteudoPorPasta map[string]string) {
	//fmt.Println("✅ Arquivo HTML:", caminho)
	
	conteudo, err := extrairConteudoHTML(caminho)
	if err != nil {
		fmt.Println("   Erro:", err)
		return
	}

	//fmt.Println("📝 Conteúdo:", conteudo)
	
	// Adiciona ao conteúdo da pasta
	if conteudoPorPasta[pasta] != "" {
		conteudoPorPasta[pasta] += " | " + conteudo
	} else {
		conteudoPorPasta[pasta] = conteudo
	}
}

// Extrai conteúdo das tags HTML do arquivo
func extrairConteudoHTML(caminho string) (string, error) {
	file, err := os.Open(caminho)
	if err != nil {
		return "", err
	}
	defer file.Close()

	doc, err := html.Parse(file)
	if err != nil {
		return "", err
	}

	var textos []string
	var extrair func(*html.Node)
	
	extrair = func(n *html.Node) {
		if n.Type == html.ElementNode && tagsExtracao[n.Data] && n.FirstChild != nil {
			if texto := strings.TrimSpace(n.FirstChild.Data); texto != "" {
				textos = append(textos, texto)
			}
		}
		
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extrair(c)
		}
	}
	
	extrair(doc)
	return strings.Join(textos, ", "), nil
}

// Mostra os resultados organizados por pasta
func mostrarResultados(conteudoPorPasta map[string]string) {
	/*fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("📚 CONTEÚDO POR PASTA")
	fmt.Println(strings.Repeat("=", 50))*/

	for pasta, conteudo := range conteudoPorPasta {
		fmt.Println(/*"\n📁 %s\n",*/ pasta)
		fmt.Println(/*"   %s\n",*/ conteudo)
		//fmt.Println(strings.Repeat("-", 50))
	}
}