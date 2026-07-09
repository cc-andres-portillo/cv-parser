package main

import (
	"context"
	"github.com/cc-andres-portillo/cv-parser/internal/adapters/handlers"
	"github.com/cc-andres-portillo/cv-parser/internal/adapters/storages"
	"github.com/cc-andres-portillo/cv-parser/internal/core/ports"
	"log"
	"net/http"
	"os"
)

func main() {
	ctx := context.Background()

	// 1. Leer variables de entorno (puedes usar librerías como 'godotenv' para cargarlas)
	puertoServidor := os.Getenv("PORT")
	if puertoServidor == "" {
		puertoServidor = ":8080" // Valor por defecto
	}
	mongoURI := os.Getenv("MONGO_URI")
	mongoDBName := os.Getenv("MONGO_DB_NAME")

	// 2. Adaptador de archivos
	extractorAdapter := storage.NewDocumentExtractorAdapter()

	// 3. Adaptador de base de datos
	var dbAdapter ports.CVRepositoryPort
	if mongoURI != "" && mongoDBName != "" {
		mongoAdapter, err := storage.ConectarAMongo(ctx, mongoURI, mongoDBName)
		if err != nil {
			log.Fatalf("No se pudo iniciar el almacenamiento en MongoDB: %v", err)
		}
		dbAdapter = mongoAdapter
		log.Println("Usando MongoDB como almacenamiento")
	} else {
		dbAdapter = storage.NewMockDatabaseAdapter()
		log.Println("Usando MockDatabaseAdapter (sin MongoDB)")
	}

	// 4. Inyección de dependencias al Núcleo de Negocio
	cvService := ports.NewCVService(extractorAdapter, dbAdapter)

	// 4. Adaptador de Entrada HTTP
	httpHandler := handlers.NewHTTPCVHandler(cvService)

	http.HandleFunc("/api/v1/parse-cv", httpHandler.ParseCVHandler)

	log.Printf("[Hexagonal MongoDB-Ready] Servidor escuchando en %s", puertoServidor)
	log.Fatal(http.ListenAndServe(puertoServidor, nil))
}
