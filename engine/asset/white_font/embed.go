package asset_white_font

import (
    _ "embed"
)

// Sprites
var (
    //go:embed whitefont.json
    Json []byte
    //go:embed whitefont.png
    Png []byte
)
