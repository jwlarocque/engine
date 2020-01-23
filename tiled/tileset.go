package tiled

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

// == TSX file parsing ================

type tilesetXML struct {
	XMLName    xml.Name   `xml:"tileset"`
	Images     []imageXML `xml:"image"`
	TileWidth  string     `xml:"tilewidth,attr"`
	TileHeight string     `xml:"tileheight,attr"`
	NumTiles   string     `xml:"tilecount,attr"`
	NumCols    string     `xml:"columns,attr"`
}

type imageXML struct {
	XMLName  xml.Name `xml:"image"`
	FilePath string   `xml:"source,attr"`
}

// NewTSXFromFile parses the given .tsx file into a tilsetXML
func newTSXFromFile(filePath string) tilesetXML {
	tsxFile, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}

	bytes, _ := ioutil.ReadAll(tsxFile)

	var tileset tilesetXML
	xml.Unmarshal(bytes, &tileset)

	defer tsxFile.Close()

	return tileset
}

// ----------------

type Tileset struct {
	TilesImage *ebiten.Image
	tileWidth  int
	tileHeight int
	numTiles   int
	numCols    int
}

// NewTilesetFromFile creates a tileset from a .tsx file
func NewTilesetFromFile(filePath string) *Tileset {
	tileset := Tileset{}

	tsx := newTSXFromFile(filePath)

	if len(tsx.Images) < 1 {
		log.Fatal(fmt.Sprintf("TSX XML at %s had no image", filePath))
	}

	// TODO: reduce repeated code
	var err error
	tileset.TilesImage, _, err = ebitenutil.NewImageFromFile(tsx.Images[0].FilePath, ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}

	tileset.tileWidth, err = strconv.Atoi(tsx.TileWidth)
	if err != nil {
		log.Fatal(err)
	}
	tileset.tileHeight, err = strconv.Atoi(tsx.TileHeight)
	if err != nil {
		log.Fatal(err)
	}
	tileset.numTiles, err = strconv.Atoi(tsx.NumTiles)
	if err != nil {
		log.Fatal(err)
	}
	tileset.numCols, err = strconv.Atoi(tsx.NumCols)
	if err != nil {
		log.Fatal(err)
	}

	return &tileset
}
