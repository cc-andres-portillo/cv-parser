package main

import (
	"archive/zip"
	"fmt"
	"github.com/cc-andres-portillo/cv-parser/internal/adapters/storages"
)

func main() {
	a := storage.NewDocumentExtractorAdapter()

	// Abrir zip y listar archivos
	zr, err := zip.OpenReader("sample_good.zip")
	if err != nil {
		fmt.Println("zip open err:", err)
		return
	}
	defer zr.Close()
	fmt.Println("Files in zip:")
	for _, f := range zr.File {
		fmt.Println(" -", f.Name)
	}

	text, err := a.ExtractTextFromDocx("sample_good.zip")
	fmt.Println("\nExtract err:", err)
	fmt.Printf("Extracted text: '%s'\n", text)
}
