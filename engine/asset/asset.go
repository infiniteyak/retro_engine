package asset

import (
    "image"
    "image/color"
    _ "image/png"
	"bytes"
    "log"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/goaseprite"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

type SpriteAsset struct {
    Path string
    Image *ebiten.Image
    ImageBytes []byte
    JsonBytes []byte
    File *goaseprite.File
}

type AudioAsset struct {
    DecodedAudio *wav.Stream
    RawAudio []byte
}

const (
    FontHeight = 8
    FontWidth = 8
)

var SpriteAssets = map[string]SpriteAsset {}

var AudioAssets = map[string]AudioAsset{}

var PolyImage *ebiten.Image

func LoadSpriteAsset(name string, json, png []byte) {
    img, _, err := image.Decode(bytes.NewReader(png)) 
    if err != nil {
        panic(err)
    }
    SpriteAssets[name] = SpriteAsset{
        JsonBytes: json, 
        ImageBytes: png,
        File: goaseprite.Read(json),
        Image: ebiten.NewImageFromImage(img),
    }
}

func InitPolyAssets() {
    PolyImage = ebiten.NewImage(3,3)
    PolyImage.Fill(color.White)
}

func LoadAudioAsset(name string, rawAudio []byte) {
    //var err error
    da, err := wav.DecodeWithoutResampling(bytes.NewReader(rawAudio))
	if err != nil {
		log.Fatal(err)
	}
    AudioAssets[name] = AudioAsset{
        DecodedAudio: da,
        RawAudio: rawAudio,
    }
}
