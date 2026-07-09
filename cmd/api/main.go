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

	// Evitar errores de "declared and not used" mientras está en modo mock
	_ = ctx
	_ = mongoURI
	_ = mongoDBName

	// 2. Adaptador de archivos
	extractorAdapter := storage.NewDocumentExtractorAdapter()

	// =========================================================================
	// 🌟 INTERRUPTOR DE BASE DE DATOS (MOCK vs MONGO)
	// =========================================================================
	
	// Opción A: Modo simulado (Para pruebas sin base de datos)
	dbAdapter := storage.NewMockDatabaseAdapter()

	// Opción B: Modo MongoDB Real (Para producción)
	// Para activar MongoDB, se debe COMENTAR la Opción A anterior
	// y DESCOMENTAR las siguientes líneas:
	/*
	if mongoURI == "" || mongoDBName == "" {
		log.Fatal("MONGO_URI y MONGO_DB_NAME son requeridos en el .env")
	}
	mongoAdapter, err := storage.ConectarAMongo(ctx, mongoURI, mongoDBName)
	if err != nil {
		log.Fatalf("No se pudo iniciar el almacenamiento: %v", err)
	}
	dbAdapter := mongoAdapter
	*/
	// =========================================================================

	// 3. Inyección de dependencias al Núcleo de Negocio
	cvService := ports.NewCVService(extractorAdapter, dbAdapter)

	// 4. Adaptador de Entrada HTTP
	httpHandler := handlers.NewHTTPCVHandler(cvService)

	http.HandleFunc("/api/v1/parse-cv", httpHandler.ParseCVHandler)

	log.Printf("[Hexagonal MongoDB-Ready] Servidor escuchando en %s", puertoServidor)
	log.Fatal(http.ListenAndServe(puertoServidor, nil))
}
