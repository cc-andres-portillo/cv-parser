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

	puertoServidor := os.Getenv("PORT")
	if puertoServidor == "" {
		puertoServidor = ":8080"
	} else if puertoServidor[0] != ':' {
		puertoServidor = ":" + puertoServidor
	}
	mongoURI := os.Getenv("MONGO_URI")
	mongoDBName := os.Getenv("MONGO_DB_NAME")

	extractorAdapter := storage.NewDocumentExtractorAdapter()

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

	cvService := ports.NewCVService(extractorAdapter, dbAdapter)
	httpHandler := handlers.NewHTTPCVHandler(cvService)

	http.HandleFunc("/api/v1/parse-cv", httpHandler.ParseCVHandler)

	log.Printf("[Hexagonal MongoDB-Ready] Servidor escuchando en %s", puertoServidor)
	log.Fatal(http.ListenAndServe(puertoServidor, nil))
}
