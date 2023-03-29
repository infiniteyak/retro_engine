package asset_black_font

import (
    _ "embed"
)

// Sprites
var (
    //go:embed blackfont.json 
    Json []byte
    //go:embed blackfont.png
    Png []byte
)
