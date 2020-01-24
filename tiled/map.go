package tiled

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten"
)

// TODO: reduce use of log.Fatal lol

type mapXML struct {
	XMLName     xml.Name        `xml:"map"`
	MapTilesets []mapTilesetXML `xml:"tileset"`
	Layers      []mapLayerXML   `xml:"layer"`
	Width       string          `xml:"width,attr"`  // map width in tiles
	Height      string          `xml:"height,attr"` // map height in tiles
}

type mapTilesetXML struct {
	XMLName  xml.Name `xml:"tileset"`
	FilePath string   `xml:"source,attr"` // tileset path relative to .tmx file
}

type mapLayerXML struct {
	XMLName xml.Name `xml:"layer"`
	Data    string   `xml:"data"`
}

// NewTMXFromFile parses the given .tmx file into a mapXML
func newTMXFromFile(filePath string) mapXML {
	tmxFile, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}

	bytes, _ := ioutil.ReadAll(tmxFile)
	var mapRaw mapXML
	xml.Unmarshal(bytes, &mapRaw)

	defer tmxFile.Close()
	return mapRaw
}

// TODO: maybe Map can be just ebiten.Image
// (discard tileset etc. after running the constructor)
type Map struct {
	Image    *ebiten.Image
	Tileset  *Tileset
	tileData []uint32
	width    int // map width in tiles
	height   int // map height in tiles
}

func parseIntCSV(csv string) []uint32 {
	strs := strings.Split(csv, ",")
	ints := make([]uint32, len(strs))
	var fatInt uint64 // what
	var err error
	for i := range ints {
		fatInt, err = strconv.ParseUint(strs[i], 10, 32)
		if err != nil {
			log.Fatal(err)
		}
		ints[i] = uint32(fatInt)
	}
	return ints
}

// TODO: clean this up
func getTileImageAndOpts(newMap *Map, tileNum int) (*ebiten.Image, *ebiten.DrawImageOptions) {
	tileID := newMap.tileData[tileNum]
	opts := &ebiten.DrawImageOptions{}

	// look mom, I'm using bitwise operators
	localID := (tileID & 0x1FFFFFFF) - 1
	flipHoriz := (tileID & 0x80000000) > 0
	flipVert := (tileID & 0x40000000) > 0
	flipDiag := (tileID & 0x20000000) > 0
	(fmt.Sprintf("ID: %d, H: %t, V: %t, D: %t", localID, flipHoriz, flipVert, flipDiag)) // TODO: remove this
	img := newMap.Tileset.GetTileImage(int(localID))
	opts.GeoM.Translate(-float64(newMap.Tileset.tileWidth)/2, -float64(newMap.Tileset.tileHeight)/2)
	if flipHoriz {
		opts.GeoM.Scale(-1, 1)
	}
	if flipVert {
		opts.GeoM.Scale(1, -1)
	}
	if flipDiag {
		opts.GeoM.Rotate(0.5 * math.Pi) // TODO: use a more intelligent matrix transform instead of this naive rotation
	}
	opts.GeoM.Translate(float64((tileNum%newMap.width)*newMap.Tileset.tileWidth), float64((tileNum/newMap.width)*newMap.Tileset.tileHeight))
	opts.GeoM.Translate(float64(newMap.Tileset.tileWidth)/2, float64(newMap.Tileset.tileHeight)/2)
	return img, opts
}

func NewMapFromFile(filePath string) *Map {
	var err error
	newMap := Map{}
	tmx := newTMXFromFile(filePath)

	newMap.width, err = strconv.Atoi(tmx.Width)
	if err != nil {
		log.Fatal(err)
	}
	newMap.height, err = strconv.Atoi(tmx.Height)
	if err != nil {
		log.Fatal(err)
	}

	if len(tmx.MapTilesets) < 1 {
		log.Fatal(fmt.Sprintf("map at %s had no tilesets", filePath))
	}
	newMap.Tileset = NewTilesetFromFile(tmx.MapTilesets[0].FilePath)

	if len(tmx.Layers) < 1 {
		log.Fatal(fmt.Sprintf("map at %s had no layers (data)", filePath))
	}
	newMap.tileData = parseIntCSV(strings.Replace(tmx.Layers[0].Data, "\n", "", -1))

	newMap.Image, err = ebiten.NewImage(newMap.width*newMap.Tileset.tileWidth, newMap.height*newMap.Tileset.tileHeight, ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < len(newMap.tileData); i++ {
		newMap.Image.DrawImage(getTileImageAndOpts(&newMap, i))
	}

	return &newMap
}
