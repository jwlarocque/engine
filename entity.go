// struct for entities - player, enemies, projectiles, etc.

package engine

import (
	"github.com/hajimehoshi/ebiten"
)

// Entity struct composes things which move: player, enemies, projectiles, etc.
type Entity struct {
	Position       Vector2
	Velocity       Vector2
	CurrentState   int
	ImageProviders map[int]ImageProvider
}

// GetImage grabs an animation frame or static image for the entity's current state
func (e Entity) GetImage() *ebiten.Image {
	return e.ImageProviders[e.CurrentState].GetImage()
}
