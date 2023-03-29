package asset_white_select

import (
    _ "embed"
)

// Sprites
var (
    //go:embed whiteselect.json
    Json []byte
    //go:embed whiteselect.png
    Png []byte
)
