package engine

import (
	"github.com/jwlarocque/engine/mech"
	"github.com/jwlarocque/engine/tiled"
)

type Level struct {
	terrainColliders []*mech.Collider
	tiled.Map
}

func NewLevelFromFile(filePath string) *Level {
	return &Level{}
}
