package engine

import (
	"github.com/jwlarocque/engine/mechanism"
	"github.com/jwlarocque/engine/tiled"
)

type Level struct {
	terrainColliders []*mechanism.Collider
	tiled.Map
}

func NewLevelFromFile(filePath string) *Level {
	return &Level{}
}
