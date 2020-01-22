package engine

import (
	"io/ioutil"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

type ImageProvider interface {
	GetImage() *ebiten.Image
}

// == Animation ================

type Animation struct {
	frameNum      int
	frameDuration time.Duration
	lastFrameTime time.Time
	Images        []*ebiten.Image
}

func (a *Animation) GetImage() *ebiten.Image {
	// TODO: global frame counter
	if time.Since(a.lastFrameTime) > a.frameDuration {
		a.lastFrameTime = time.Now()
		a.frameNum = (a.frameNum + 1) % len(a.Images)
	}
	return a.Images[a.frameNum]
}

// NewAnimationFromSheet returns an Animation given a folder of image files and a frameDuration
// WIP: function not complete, since no animation sheets in use
// ? needs additional args for position, size, ... ?
func NewAnimationFromSheet(path string, frameDuration time.Duration) *Animation {
	return &Animation{}
}

func NewAnimationFromFolder(path string, frameDuration time.Duration) *Animation {
	a := Animation{}
	a.lastFrameTime = time.Now()
	a.frameDuration = frameDuration

	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		img, _, err := ebitenutil.NewImageFromFile(path+"/"+f.Name(), ebiten.FilterDefault)
		if err != nil {
			log.Fatal(err)
		}
		a.Images = append(a.Images, img)
	}

	return &a
}

// == Animation ================

type Sprite struct {
	Mirrors int
	Image   *ebiten.Image
}

func (s *Sprite) GetImage() *ebiten.Image {
	return s.Image
}
