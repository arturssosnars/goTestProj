package DataModules

// Go struktiem var norādīt vēlamo JSON lauka atslēgu ar anotāciju, piemēram, `json:"message"`
// Labāk būtu atgriezt json atslēgas ar camelCase.
type Error struct {
	Message string
}