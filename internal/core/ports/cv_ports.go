package ports

import (
	"context"
	"github.com/cc-andres-portillo/cv-parser/internal/core/domain"
	"io"
)

// ReaderExtractorPort es un puerto de salida (Driven Adapter)
// Define cómo cualquier adaptador externo debe extraer texto de un archivo
type ReaderExtractorPort interface {
	ExtractTextFromPDF(r io.ReaderAt, size int64) (string, error)
	ExtractTextFromDocx(filePath string) (string, error)
}

// Define el contrato para persistir datos
type CVRepositoryPort interface {
	Save(ctx context.Context, cv *domain.FormularioCV) error
}

// CVServicePort es un puerto de entrada (Driver Adapter)
// Define los casos de uso disponibles para los clientes (HTTP, CLI, etc.)
type CVServicePort interface {
	ParseCV(ctx context.Context, file io.ReaderAt, size int64, extension string) (domain.FormularioCV, error)
}
