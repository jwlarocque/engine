package engine

import (
	"github.com/hajimehoshi/ebiten"
)

/* TODO:
 *     remappable inputs
 *     multiple input buttons/keys for same
 *     gamepad input
 */
func GetMovementVector() Vector2 {
	v := Vector2{}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		v.Y -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		v.Y += 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		v.X -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		v.X += 1
	}

	return v.Normalize()
}
