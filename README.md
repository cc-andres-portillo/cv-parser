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
- `MONGO_URI` — URI de MongoDB (si quieres usar Mongo real)
- `MONGO_DB_NAME` — nombre de la base de datos en Mongo

En PowerShell:

```powershell
$env:PORT=":8080"
$env:MONGO_URI="mongodb://localhost:27017"
$env:MONGO_DB_NAME="cv_db"
go run ./cmd/api
```

Notas
- El módulo del proyecto está declarado como `github.com/cc-andres-portillo/cv-parser` en `go.mod`.
- El proyecto ya incluye un adaptador `MockDatabaseAdapter` para desarrollo sin MongoDB.

Licencia
- Revisa y añade una licencia si es necesario.
