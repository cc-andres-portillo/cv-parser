package storage

import (
	"context"
	"github.com/cc-andres-portillo/cv-parser/internal/core/domain"
	"log"
)

type MockDatabaseAdapter struct{}

func NewMockDatabaseAdapter() *MockDatabaseAdapter {
	return &MockDatabaseAdapter{}
}

// Cumple con el contrato de CVRepositoryPort
func (m *MockDatabaseAdapter) Save(ctx context.Context, cv *domain.FormularioCV) error {
	log.Printf("[MOCK DB] Guardando exitosamente el CV de: %s", cv.NombreCompleto)
	return nil
}
