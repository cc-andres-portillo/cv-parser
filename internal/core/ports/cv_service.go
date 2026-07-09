package ports

import (
	"context"
	"github.com/cc-andres-portillo/cv-parser/internal/core/domain"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

type CVService struct {
	extractor ReaderExtractorPort
	repo      CVRepositoryPort
}

func NewCVService(extractor ReaderExtractorPort, repo CVRepositoryPort) CVServicePort {
	return &CVService{
		extractor: extractor,
		repo:      repo,
	}
}

func (s *CVService) ParseCV(ctx context.Context, file io.ReaderAt, size int64, extension string) (domain.FormularioCV, error) {
	var textoCV string
	var err error

	// Delegar la extracción al adaptador correspondiente según la extensión
	switch strings.ToLower(extension) {
	case ".pdf":
		textoCV, err = s.extractor.ExtractTextFromPDF(file, size)
	case ".docx":
		// Docx requiere archivo en disco temporal debido a limitaciones de su librería
		tempFile, tmpErr := os.CreateTemp("", "cv-core-*.docx")
		if tmpErr != nil {
			return domain.FormularioCV{}, tmpErr
		}
		defer os.Remove(tempFile.Name())
		defer tempFile.Close()
		if _, copyErr := io.Copy(tempFile, io.NewSectionReader(file, 0, size)); copyErr != nil {
			return domain.FormularioCV{}, copyErr
		}
		textoCV, err = s.extractor.ExtractTextFromDocx(tempFile.Name())
	default:
		return domain.FormularioCV{}, fmt.Errorf("extensión no soportada: %s", extension)
	}

	if err != nil {
		return domain.FormularioCV{}, err
	}

	resultado := s.ejecutarReglasDeNegocio(textoCV)

	// 🌟 CASO DE USO: Persistimos los datos estructurados en la BD usando el puerto
	if err := s.repo.Save(ctx, &resultado); err != nil {
		return domain.FormularioCV{}, fmt.Errorf("error al persistir en BD: %v", err)
	}

	return resultado, nil
}

func (s *CVService) ejecutarReglasDeNegocio(textoCV string) domain.FormularioCV {
	var resultado domain.FormularioCV

	lineas := strings.Split(textoCV, "\n")
	var lineasLimpias []string
	for _, l := range lineas {
		t := strings.TrimSpace(l)
		if t != "" {
			lineasLimpias = append(lineasLimpias, t)
		}
	}

	if len(lineasLimpias) > 0 {
		resultado.NombreCompleto = lineasLimpias[0]
	}

	reEmail := regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)
	resultado.Email = reEmail.FindString(textoCV)

	reTelefono := regexp.MustCompile(`\+?\d{1,4}?[\s-]?\(?\d{1,3}?\)?[\s-]?\d{3,4}[\s-]?\d{3,4}`)
	resultado.Telefono = reTelefono.FindString(textoCV)

	seccionActual := ""
	var expTemporal *domain.Experiencia

	esSeccionExperiencia := regexp.MustCompile(`(?i)(experiencia|historial laboral|work experience|trayectoria)`)
	esSeccionHabilidades := regexp.MustCompile(`(?i)(habilidades|skills|conocimientos|competencias)`)

	for _, linea := range lineasLimpias {
		if esSeccionExperiencia.MatchString(linea) {
			seccionActual = "experiencia"
			continue
		}
		if esSeccionHabilidades.MatchString(linea) {
			seccionActual = "habilidades"
			continue
		}

		switch seccionActual {
		case "habilidades":
			if strings.Contains(linea, ",") {
				elementos := strings.Split(linea, ",")
				for _, el := range elementos {
					resultado.Habilidades = append(resultado.Habilidades, strings.Trim(el, " •-*"))
				}
			} else {
				resultado.Habilidades = append(resultado.Habilidades, strings.Trim(linea, " •-*"))
			}
		case "experiencia":
			contieneFecha := regexp.MustCompile(`(?i)(\d{4}|enero|febrero|marzo|abril|mayo|junio|julio|agosto|septiembre|octubre|noviembre|diciembre|actualidad|presente)`)
			if contieneFecha.MatchString(linea) {
				if expTemporal != nil {
					expTemporal.Periodo = linea
					resultado.Experiencias = append(resultado.Experiencias, *expTemporal)
					expTemporal = nil
				}
			} else {
				if expTemporal == nil {
					expTemporal = &domain.Experiencia{}
					if strings.Contains(linea, " - ") {
						partes := strings.Split(linea, " - ")
						expTemporal.Puesto = partes[0]
						expTemporal.Empresa = partes[1]
					} else {
						expTemporal.Puesto = linea
						expTemporal.Empresa = "No identificada"
					}
				}
			}
		}
	}

	if expTemporal != nil {
		resultado.Experiencias = append(resultado.Experiencias, *expTemporal)
	}

	return resultado
}
