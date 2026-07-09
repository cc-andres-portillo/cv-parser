package storage

import (
	"archive/zip"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
)

// ExtractTextFromDocx realiza una extracción básica de texto desde el archivo .docx
// Abre el .docx como zip, lee word/document.xml y elimina etiquetas XML.
func (a *DocumentExtractorAdapter) ExtractTextFromDocx(filePath string) (string, error) {
	r, err := zip.OpenReader(filePath)
	if err != nil {
		return "", err
	}
	defer r.Close()

	var content []byte
	for _, f := range r.File {
		if filepath.ToSlash(f.Name) == "word/document.xml" {
			rc, err := f.Open()
			if err != nil {
				return "", err
			}
			data, err := ioutil.ReadAll(rc)
			rc.Close()
			if err != nil {
				return "", err
			}
			content = data
			break
		}
	}

	if content == nil {
		return "", nil
	}

	// Reemplazar cierres de párrafo por nueva línea y eliminar tags XML simples
	contentStr := string(content)
	contentStr = strings.ReplaceAll(contentStr, "</w:p>", "\n")
	contentStr = strings.ReplaceAll(contentStr, "</p>", "\n")
	re := regexp.MustCompile("<[^>]+>")
	text := strings.TrimSpace(re.ReplaceAllString(contentStr, ""))
	return text, nil
}
