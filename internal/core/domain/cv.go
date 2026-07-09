package domain

type Experiencia struct {
	Empresa string `json:"empresa"`
	Puesto  string `json:"puesto"`
	Periodo string `json:"periodo"`
}

type FormularioCV struct {
	NombreCompleto string        `json:"nombre_completo"`
	Email          string        `json:"email"`
	Telefono       string        `json:"telefono"`
	Habilidades    []string      `json:"habilidades"`
	Experiencias   []Experiencia `json:"experiencias"`
}
