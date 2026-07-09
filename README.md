# cv-parser

Pequeño servicio en Go para extraer y parsear CVs (PDF / DOCX).

Requisitos
- Go 1.20+ (se usó 1.24 en go.mod)

Rápido inicio
1. Obtener dependencias:

```bash
go mod tidy
```

2. Compilar todo el módulo:

```bash
go build ./...
```

3. Ejecutar el servidor HTTP (desde la raíz del repo):

```bash
go run ./cmd/api
```

Variables de entorno (opcional)
- `PORT` — puerto en que el servidor escucha (por defecto `:8080`)
- `MONGO_URI` — URI de MongoDB. Si se establece junto a `MONGO_DB_NAME`, la aplicación usará MongoDB automáticamente.
- `MONGO_DB_NAME` — nombre de la base de datos en Mongo

En PowerShell:

```powershell
$env:PORT=":8080"
$env:MONGO_URI="mongodb://localhost:27017"
$env:MONGO_DB_NAME="cv_db"
go run ./cmd/api
```

Si no estableces `MONGO_URI` y `MONGO_DB_NAME`, el servicio funciona con `MockDatabaseAdapter` y no requiere MongoDB.

Ejemplo de petición `curl`

```bash
curl -X POST \
  -F "cv=@/ruta/a/tu/cv.docx" \
  http://localhost:8080/api/v1/parse-cv
```

Reemplaza `/ruta/a/tu/cv.docx` por la ruta de tu archivo `.docx` o `.pdf`.

Notas
- El módulo del proyecto está declarado como `github.com/cc-andres-portillo/cv-parser` en `go.mod`.
- El servicio extrae texto de `.pdf` y `.docx` y responde JSON.
- Para pruebas locales, sigue usando `go run ./cmd/api`.

Licencia
- Revisa y añade una licencia si es necesario.
