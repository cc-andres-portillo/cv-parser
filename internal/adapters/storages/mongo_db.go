package storage

import (
	"context"
	"github.com/cc-andres-portillo/cv-parser/internal/core/domain"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDatabaseAdapter struct {
	client     *mongo.Client
	dbName     string
	collection string
}

// ConectarAMongo inicializa el cliente con la URI del .env
func ConectarAMongo(ctx context.Context, uri, dbName string) (*MongoDatabaseAdapter, error) {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("error al conectar a MongoDB: %w", err)
	}

	// Verificar la conexión (Ping)
	ctxPing, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	if err := client.Ping(ctxPing, nil); err != nil {
		return nil, fmt.Errorf("ping fallido a MongoDB: %w", err)
	}

	log.Println("[📌 INFRAESTRUCTURA] Conexión exitosa a MongoDB")

	return &MongoDatabaseAdapter{
		client:     client,
		dbName:     dbName,
		collection: "cv_registrados", // Nombre de la colección en la base de datos
	}, nil
}

// Save cumple estrictamente con el contrato de CVRepositoryPort
func (m *MongoDatabaseAdapter) Save(ctx context.Context, cv *domain.FormularioCV) error {
	coll := m.client.Database(m.dbName).Collection(m.collection)
	
	ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := coll.InsertOne(ctxTimeout, cv)
	if err != nil {
		return fmt.Errorf("error al insertar documento en Mongo: %w", err)
	}

	log.Printf("[MongoDB] CV de %s guardado correctamente", cv.NombreCompleto)
	return nil
}
