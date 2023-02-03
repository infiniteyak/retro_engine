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

const (
    FontHeight = 8
    FontWidth = 8
)

var SpriteAssets = map[string]SpriteAsset {
    "WhiteFont":      SpriteAsset{JsonBytes: WhiteFont_json, ImageBytes: WhiteFont_png},
    "BlackFont":      SpriteAsset{JsonBytes: BlackFont_json, ImageBytes: BlackFont_png},
    "LightBlueFont":      SpriteAsset{JsonBytes: LightBlueFont_json, ImageBytes: LightBlueFont_png},
    "RedFont":      SpriteAsset{JsonBytes: RedFont_json, ImageBytes: RedFont_png},
    "PurpleFont":      SpriteAsset{JsonBytes: PurpleFont_json, ImageBytes: PurpleFont_png},
    "TitleFont":      SpriteAsset{JsonBytes: TitleFont_json, ImageBytes: TitleFont_png},
    "SmallAsteroid":  SpriteAsset{JsonBytes: SmallAsteroid_json, ImageBytes: SmallAsteroid_png},
    "MediumAsteroid": SpriteAsset{JsonBytes: MediumAsteroid_json, ImageBytes: MediumAsteroid_png},
    "LargeAsteroid":  SpriteAsset{JsonBytes: LargeAsteroid_json, ImageBytes: LargeAsteroid_png},
    "SimpleBullet":   SpriteAsset{JsonBytes: SimpleBullet_json, ImageBytes: SimpleBullet_png},
    "SmallExplosion": SpriteAsset{JsonBytes: SmallExplosion_json, ImageBytes: SmallExplosion_png},
    "Star": SpriteAsset{JsonBytes: Star_json, ImageBytes: Star_png},
    "PlayerShip": SpriteAsset{JsonBytes: PlayerShip_json, ImageBytes: PlayerShip_png},
    "Boomerang": SpriteAsset{JsonBytes: Boomerang_json, ImageBytes: Boomerang_png},
    "Alien1": SpriteAsset{JsonBytes: Alien1_json, ImageBytes: Alien1_png},
}

var PolyImage *ebiten.Image

func InitSpriteAssets() {
    PolyImage = ebiten.NewImage(3,3)
    PolyImage.Fill(color.White)

    for i, s := range SpriteAssets {
        s.File = goaseprite.Read(s.JsonBytes) //file keeps track of frame stuff so don't want global

        img, _, err := image.Decode(bytes.NewReader(s.ImageBytes)) 
        if err != nil {
            panic(err)
        }
        s.Image = ebiten.NewImageFromImage(img)
        SpriteAssets[i] = s
    }
}

var (
    FireD *wav.Stream
    HitD *wav.Stream
    MenuD *wav.Stream
    DestroyedD *wav.Stream
    ThrusterD *wav.Stream
    WaveD *wav.Stream
)

func InitAudioAssets() {
    var err error
    FireD, err = wav.DecodeWithoutResampling(bytes.NewReader(Fire_wav))
	if err != nil {
		log.Fatal(err)
	}

    HitD, err = wav.DecodeWithoutResampling(bytes.NewReader(Hit_wav))
	if err != nil {
		log.Fatal(err)
	}

    MenuD, err = wav.DecodeWithoutResampling(bytes.NewReader(Menu_wav))
	if err != nil {
		log.Fatal(err)
	}

    DestroyedD, err = wav.DecodeWithoutResampling(bytes.NewReader(Destroyed_wav))
	if err != nil {
		log.Fatal(err)
	}

    ThrusterD, err = wav.DecodeWithoutResampling(bytes.NewReader(Thruster_wav))
	if err != nil {
		log.Fatal(err)
	}

    WaveD, err = wav.DecodeWithoutResampling(bytes.NewReader(Wave_wav))
	if err != nil {
		log.Fatal(err)
	}
}
