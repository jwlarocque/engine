package tiled

import (
	"encoding/json"
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

// Tiled JSON Format: https://doc.mapeditor.org/en/stable/reference/json-map-format/
// Tiled TMX Format: https://doc.mapeditor.org/en/stable/reference/tmx-map-format/
//
// TODO: reduce use of log.Fatal

// Map represents the data about a level which can be found in a Tiled file
// TODO: maybe Map can be just ebiten.Image
// (discard tileset etc. after running the constructor)
type Map struct {
	Image    *ebiten.Image
	Tileset  *Tileset
	tileData []uint32
	width    int // map width in tiles
	height   int // map height in tiles
}

// TODO: clean this up
// helper used by both JSON and TMX Map constructors
func getTileImageAndOpts(newMap *Map, tileNum int) (*ebiten.Image, *ebiten.DrawImageOptions) {
	tileID := newMap.tileData[tileNum]
	opts := &ebiten.DrawImageOptions{}

	// bits 32, 31, and 30 store whether tiles are flipped
	localID := (tileID & 0x1FFFFFFF) - 1 // TODO: use firstID/firstgid instead of hardcoding 1
	flipHoriz := (tileID & 0x80000000) > 0
	flipVert := (tileID & 0x40000000) > 0
	flipDiag := (tileID & 0x20000000) > 0

	// apply tile flips/rotatoin
	opts.GeoM.Translate(-float64(newMap.Tileset.tileWidth)/2, -float64(newMap.Tileset.tileHeight)/2)
	if flipDiag {
		opts.GeoM.Rotate(0.5 * math.Pi)
	}
	if flipHoriz {
		opts.GeoM.Scale(-1, 1)
	}
	if flipVert {
		opts.GeoM.Scale(1, -1)
	}
	// translate to position relative to rest of map
	opts.GeoM.Translate(float64((tileNum%newMap.width)*newMap.Tileset.tileWidth), float64((tileNum/newMap.width)*newMap.Tileset.tileHeight))
	opts.GeoM.Translate(float64(newMap.Tileset.tileWidth)/2, float64(newMap.Tileset.tileHeight)/2)

	img := newMap.Tileset.GetTileImage(int(localID))
	return img, opts
}

// == JSON ========

type mapJSON struct {
	MapTilesets []mapTilesetJSON `json:"tilesets"`
	Layers      []mapLayerJSON   `json:"layers"`
	Width       int              `json:"width"`
	Height      int              `json:"height"`
}

type mapTilesetJSON struct {
	FilePath string `json:"source"`
	firstID  int    `json:"firstgid"`
}

type mapLayerJSON struct {
	Data []uint32 // TODO: dunno if unmarshaling straight to uint32 slice will work
}

// newJSONFromFile parses the given .json file into a mapJSON
func newMapJSONFromFile(filePath string) mapJSON {
	jsonFile, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}

	bytes, _ := ioutil.ReadAll(jsonFile)
	var mapRaw mapJSON
	json.Unmarshal(bytes, &mapRaw)

	defer jsonFile.Close()
	return mapRaw
}

// NewMapFromJSON returns a Map given a .json map file
func NewMapFromJSON(filePath string) *Map {
	var err error
	newMap := Map{}
	json := newMapJSONFromFile(filePath)

	newMap.width = json.Width
	newMap.height = json.Height

	if len(json.MapTilesets) < 1 {
		log.Fatal(fmt.Sprintf("map at %s had no tilesets", filePath))
	}
	newMap.Tileset = NewTilesetFromJSON(json.MapTilesets[0].FilePath)

	if len(json.Layers) < 1 {
		log.Fatal(fmt.Sprintf("map at %s had no layers (data)", filePath))
	}
	newMap.tileData = json.Layers[0].Data

	newMap.Image, err = ebiten.NewImage(newMap.width*newMap.Tileset.tileWidth, newMap.height*newMap.Tileset.tileHeight, ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < len(newMap.tileData); i++ {
		newMap.Image.DrawImage(getTileImageAndOpts(&newMap, i))
	}

	return &newMap
}

// == XML (TMX) ========

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
	firstID  string   `xml:"firstgid,attr"`
}

type mapLayerXML struct {
	XMLName xml.Name `xml:"layer"`
	Data    string   `xml:"data"`
}

// NewTMXFromFile parses the given .tmx file into a mapXML
func newMapTMXFromFile(filePath string) mapXML {
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

// parseIntCSV converts the data string in a TMX layer into []uint32
func parseIntCSV(csv string) []uint32 {
	strs := strings.Split(csv, ",")
	ints := make([]uint32, len(strs))
	var fatInt uint64 // ParseUint returns uint64, cast later
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

// NewMapFromTMX returns a Map given a .tmx map file
func NewMapFromTMX(filePath string) *Map {
	var err error
	newMap := Map{}
	tmx := newMapTMXFromFile(filePath)

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
	newMap.Tileset = NewTilesetFromTSX(tmx.MapTilesets[0].FilePath)

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
