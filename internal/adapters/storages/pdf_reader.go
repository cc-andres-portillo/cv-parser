package storage

import (
	"io"
	"strings"

	"github.com/ledongthuc/pdf"
)

type DocumentExtractorAdapter struct{}

func NewDocumentExtractorAdapter() *DocumentExtractorAdapter {
	return &DocumentExtractorAdapter{}
}

func (a *DocumentExtractorAdapter) ExtractTextFromPDF(r io.ReaderAt, size int64) (string, error) {
	pdfReader, err := pdf.NewReader(r, size)
	if err != nil {
		return "", err
	}
	b, err := pdfReader.GetPlainText()
	if err != nil {
		return "", err
	}
	var buf strings.Builder
	_, err = io.Copy(&buf, b)
	return buf.String(), err
}
