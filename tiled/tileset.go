package tiled

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"image"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

type Tileset struct {
	tilesImage *ebiten.Image
	tileWidth  int
	tileHeight int
	numTiles   int
	numCols    int
}

// == JSON ========

type tilesetJSON struct {
	Name       string
	Image      string
	TileHeight int `json:"tileheight"`
	TileWidth  int `json:"tilewidth"`
	NumTiles   int `json:"tilecount"`
	NumCols    int `json:"columns"`
}

func newTilesetJSONFromFile(filePath string) tilesetJSON {
	jsonFile, err := os.Open(filePath)
	if err != nil {
		log.Println("err loading json tileset file")
		log.Fatal(err)
	}
	bytes, _ := ioutil.ReadAll(jsonFile)

	var tileset tilesetJSON
	json.Unmarshal(bytes, &tileset)

	defer jsonFile.Close()
	return tileset
}

func NewTilesetFromJSON(filePath string) *Tileset {
	tileset := Tileset{}

	json := newTilesetJSONFromFile(filePath)

	var err error
	tileset.tilesImage, _, err = ebitenutil.NewImageFromFile(json.Image, ebiten.FilterDefault)
	if err != nil {
		log.Println(filePath)
		log.Println(json.Name)
		log.Println(json.Image)
		log.Println("err loading tileset into ebiten image")
		log.Fatal(err)
	}

	// ¯\_(ツ)_/¯
	tileset.tileHeight = json.TileHeight
	tileset.tileWidth = json.TileWidth
	tileset.numTiles = json.NumTiles
	tileset.numCols = json.NumCols

	return &tileset
}

// == TSX ========

type tilesetXML struct {
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

// NewTilesetFromTSX creates a tileset from a .tsx file
func NewTilesetFromTSX(filePath string) *Tileset {
	tileset := Tileset{}

	tsx := newTSXFromFile(filePath)

	if len(tsx.Images) < 1 {
		log.Fatal(fmt.Sprintf("TSX XML at %s had no image", filePath))
	}

	// TODO: reduce repeated code
	var err error
	// TODO: construct file paths relative to .tsx file instead of this file
	tileset.tilesImage, _, err = ebitenutil.NewImageFromFile(tsx.Images[0].FilePath, ebiten.FilterDefault)
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

// GetTileImage takes a tile ID and returns the corresponding ebiten.Image
// from its tileset
// !!! NOTE: the global tile ID in a .tmx file and local ID used by the
//           tileset are generally not the same.  It is up to you to do
//           the conversion. !!!
// TODO: consider returning render opts? (would probably require global ID)
func (ts Tileset) GetTileImage(localTileID int) *ebiten.Image {
	subX := (localTileID % ts.numCols) * ts.tileWidth
	subY := (localTileID / ts.numCols) * ts.tileHeight
	return ts.tilesImage.SubImage(image.Rect(subX, subY, subX+ts.tileWidth, subY+ts.tileWidth)).(*ebiten.Image)
}
