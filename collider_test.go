package engine

import (
	"image/color"
	"log"
	"testing"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

func Test_isConvex(t *testing.T) {
	var convex bool
	// square, expect convex = true
	convex = isConvex([]*Vector2{{0.0, 0.0}, {1.0, 0.0}, {1.0, 1.0}, {0.0, 1.0}})
	if !convex {
		t.Errorf("isConvex (square): expected true, got false!")
	}

	// chevron, expect convex = false
	convex = isConvex([]*Vector2{{1.0, 0.0}, {2.0, 2.0}, {1.0, 1.0}, {0.0, 2.0}})
	if convex {
		t.Errorf("isConvex (chevron): expected false, got true!")
	}
}

var testCollider Collider
var testColliderImg *ebiten.Image

func collidersInteractiveUpdate(screen *ebiten.Image) error {
	for i := 0; i < len(testCollider.Vertices)-1; i++ {
		//log.Printf("(%.4f, %.4f); (%.4f, %.4f)", testCollider.Vertices[i].X, testCollider.Vertices[i].Y, testCollider.Vertices[i+1].X, testCollider.Vertices[i+1].Y)
		//ebitenutil.DrawLine(screen, testCollider.Vertices[i].X, testCollider.Vertices[i].Y, testCollider.Vertices[i+1].X, testCollider.Vertices[i+1].Y, color.Gray{90})
	}
	screen.DrawImage(testColliderImg, &ebiten.DrawImageOptions{})
	ebitenutil.DrawLine(screen, 1.0, 1.0, 20.0, 20.0, color.Gray{80})
	return nil
}

func Test_CollidersInteractive(t *testing.T) {
	diamCollider, err := NewCollider([]*Vector2{{10.0, 10.0}, {50.0, 50.0}})
	if err != nil {
		log.Fatal((err))
	}
	testCollider = *diamCollider
	testColliderImg, err = ebiten.NewImage(50, 50, ebiten.FilterDefault)
	if err != nil {
		log.Fatal((err))
	}

	if err := ebiten.Run(collidersInteractiveUpdate, 400, 240, 2, "Interactive Colliders Test"); err != nil {
		log.Fatal(err)
	}
}
