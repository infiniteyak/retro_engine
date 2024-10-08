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
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/audio"
    "io"
)

type SpriteAsset struct {
    Path string
    Image *ebiten.Image
    ImageBytes []byte
    JsonBytes []byte
    File *goaseprite.File
}

type audioStream interface {
    io.ReadSeeker
    Length() int64
}

const (
    Mp3_audiotype int = iota
    Wav_audiotype

)
type AudioAsset struct {
    RawAudio []byte
    AudioType int
}

var SpriteAssets = map[string]SpriteAsset {}

var AudioAssets = map[string]AudioAsset{}

var PolyImage *ebiten.Image

var AudioContext *audio.Context
var MusicPlayer *audio.Player
var CurrentMusic string

const (
    DefaultMusicVolume = 0.5
    DefaultSFXVolume = 0.5
)
var MusicVolume float64
var SFXVolume float64

// TODO move to utilities?
func isColorMatch(a color.Color, b color.Color) bool {
    aR, aG, aB, aA := a.RGBA()
    bR, bG, bB, bA := b.RGBA()
    if aR == bR &&
       aG == bG &&
       aB == bB &&
       aA == bA {
        return true
    }
    return false
}

//For some reason can't be run immediately?
func SwapColor(name string, oldColor color.Color, newColor color.Color) {
    maxX := SpriteAssets[name].Image.Bounds().Max.X
    maxY := SpriteAssets[name].Image.Bounds().Max.Y

    for x := SpriteAssets[name].Image.Bounds().Min.X; x < maxX; x++ {
        for y := SpriteAssets[name].Image.Bounds().Min.Y; y < maxY; y++ {
            if isColorMatch(SpriteAssets[name].Image.At(x , y), oldColor) {
                SpriteAssets[name].Image.Set(x , y, newColor)
            }
        }
    }
}

func DuplicateSpriteAsset(oldName string, newName string) {
    SpriteAssets[newName] = SpriteAsset{
        JsonBytes: SpriteAssets[oldName].JsonBytes, 
        ImageBytes: SpriteAssets[oldName].ImageBytes,
        File: SpriteAssets[oldName].File,
        Image: ebiten.NewImageFromImage(SpriteAssets[oldName].Image),
    }
}

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

func InitAudioAssets() {
    MusicVolume = DefaultMusicVolume
    SFXVolume = DefaultSFXVolume
    AudioContext = audio.NewContext(48000)
}

func LoadWavAudioAsset(name string, rawAudio []byte) {
    AudioAssets[name] = AudioAsset{
        RawAudio: rawAudio,
        AudioType: Wav_audiotype,
    }
}

func LoadMp3AudioAsset(name string, rawAudio []byte) {
    AudioAssets[name] = AudioAsset{
        RawAudio: rawAudio,
        AudioType: Mp3_audiotype,
    }
}

func PlaySound(name string) {
    var err error
    var s audioStream
    switch AudioAssets[name].AudioType {
    case Mp3_audiotype:
        s, err = mp3.DecodeWithoutResampling(bytes.NewReader(AudioAssets[name].RawAudio))
    case Wav_audiotype:
        s, err = wav.DecodeWithoutResampling(bytes.NewReader(AudioAssets[name].RawAudio))
    }
    if err != nil {
        log.Fatal(err)
    }
    player, err := AudioContext.NewPlayer(s)
    if err != nil {
        log.Fatal(err)
    }
    player.SetVolume(SFXVolume)
    player.Rewind()
    player.Play()
}

func PlayMusic(name string) {
    if MusicPlayer != nil && MusicPlayer.IsPlaying() && name == CurrentMusic {
        return
    }

    println("playing music")

    var err error
    var s audioStream
    switch AudioAssets[name].AudioType {
    case Mp3_audiotype:
        s, err = mp3.DecodeWithoutResampling(bytes.NewReader(AudioAssets[name].RawAudio))
    case Wav_audiotype:
        s, err = wav.DecodeWithoutResampling(bytes.NewReader(AudioAssets[name].RawAudio))
    }
    if err != nil {
        log.Fatal(err)
    }
    infLoop := audio.NewInfiniteLoop(s, s.Length())
    MusicPlayer, err = AudioContext.NewPlayer(infLoop)
    if err != nil {
        log.Fatal(err)
    }
    MusicPlayer.SetVolume(MusicVolume)
    MusicPlayer.Rewind()
    MusicPlayer.Play()
    CurrentMusic = name
}

func StopMusic() {
    if MusicPlayer == nil {
        return
    }
    println("stopping music")
    MusicPlayer.Close()
    CurrentMusic = ""
}

func SetMusicVolume(v float64) {
    if MusicPlayer != nil && 0 <= v && 1 >= v {
        MusicVolume = v
        MusicPlayer.SetVolume(v)
    }
}

func SetSFXVolume(v float64) {
    if 0 <= v && 1 >= v {
        SFXVolume = v
    }
}

func GetMusicVolume() float64 {
    return MusicVolume
}

func GetSFXVolume() float64 {
    return SFXVolume
}
