package storage

import (
	"archive/zip"
	"io/ioutil"
	"regexp"
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
		if f.Name == "word/document.xml" {
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

	// Eliminar tags XML simples
	re := regexp.MustCompile("<[^>]+>")
	text := re.ReplaceAllString(string(content), "")
	return text, nil
}
