package asset

import (
	"image"
)

const (
    FontHeight = 8
    FontWidth = 8
)

//ABCDEFGHIJKLM
//NOPQRSTUVWXYZ
//0123456789_.^
var FontMasks = map[string]image.Rectangle{
    "A": image.Rect(0,0,8,8),
    "B": image.Rect(8,0,16,8),
    "C": image.Rect(16,0,24,8),
    "D": image.Rect(24,0,32,8),
    "E": image.Rect(32,0,40,8),
    "F": image.Rect(40,0,48,8),
    "G": image.Rect(48,0,56,8),
    "H": image.Rect(56,0,64,8),
    "I": image.Rect(64,0,72,8),
    "J": image.Rect(72,0,80,8),
    "K": image.Rect(80,0,88,8),
    "L": image.Rect(88,0,96,8),
    "M": image.Rect(96,0,104,8),
    "N": image.Rect(0,8,8,16),
    "O": image.Rect(8,8,16,16),
    "P": image.Rect(16,8,24,16),
    "Q": image.Rect(24,8,32,16),
    "R": image.Rect(32,8,40,16),
    "S": image.Rect(40,8,48,16),
    "T": image.Rect(48,8,56,16),
    "U": image.Rect(56,8,64,16),
    "V": image.Rect(64,8,72,16),
    "W": image.Rect(72,8,80,16),
    "X": image.Rect(80,8,88,16),
    "Y": image.Rect(88,8,96,16),
    "Z": image.Rect(96,8,104,16),
    "0": image.Rect(0,16,8,24),
    "1": image.Rect(8,16,16,24),
    "2": image.Rect(16,16,24,24),
    "3": image.Rect(24,16,32,24),
    "4": image.Rect(32,16,40,24),
    "5": image.Rect(40,16,48,24),
    "6": image.Rect(48,16,56,24),
    "7": image.Rect(56,16,64,24),
    "8": image.Rect(64,16,72,24),
    "9": image.Rect(72,16,80,24),
    "_": image.Rect(80,16,88,24),
    ".": image.Rect(88,16,96,24),
    "^": image.Rect(96,16,104,24),
    "a": image.Rect(0,24,8,32),
    "b": image.Rect(8,24,16,32),
    "c": image.Rect(16,24,24,32),
    "d": image.Rect(24,24,32,32),
    "e": image.Rect(32,24,40,32),
    "f": image.Rect(40,24,48,32),
    "g": image.Rect(48,24,56,32),
    "h": image.Rect(56,24,64,32),
    "i": image.Rect(64,24,72,32),
    "j": image.Rect(72,24,80,32),
    "k": image.Rect(80,24,88,32),
    "l": image.Rect(88,24,96,32),
    "m": image.Rect(96,24,104,32),
    "n": image.Rect(0,32,8,40),
    "o": image.Rect(8,32,16,40),
    "p": image.Rect(16,32,24,40),
    "q": image.Rect(24,32,32,40),
    "r": image.Rect(32,32,40,40),
    "s": image.Rect(40,32,48,40),
    "t": image.Rect(48,32,56,40),
    "u": image.Rect(56,32,64,40),
    "v": image.Rect(64,32,72,40),
    "w": image.Rect(72,32,80,40),
    "x": image.Rect(80,32,88,40),
    "y": image.Rect(88,32,96,40),
    "z": image.Rect(96,32,104,40),
    "?": image.Rect(0,40,8,48),
    ",": image.Rect(8,40,16,48),
    "'": image.Rect(16,40,24,48),
}
