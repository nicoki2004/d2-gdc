package cache

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// Manager abstrae la gestión de caché
type Manager interface {
	Load(filename string, target interface{}) error
	Save(filename string, data interface{}) error
}

// FileCache implementa Manager usando el sistema de archivos
type FileCache struct{}

// NewFileCache crea un nuevo gestor de caché basado en archivos
func NewFileCache() *FileCache {
	return &FileCache{}
}

// Load intenta leer datos del caché desde un archivo
func (fc *FileCache) Load(filename string, target interface{}) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("cache file not found: %w", err)
	}

	if err := json.Unmarshal(data, target); err != nil {
		return fmt.Errorf("error unmarshaling cache: %w", err)
	}

	return nil
}

// Save guarda datos en el caché
func (fc *FileCache) Save(filename string, data interface{}) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Printf("Error serializando caché: %v", err)
		return err
	}

	err = os.WriteFile(filename, jsonData, 0o644)
	if err != nil {
		log.Printf("Error escribiendo archivo de caché: %v", err)
		return err
	}

	return nil
}
