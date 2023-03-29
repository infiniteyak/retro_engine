package shape_courier_asset_items

import (
    _ "embed"
)

// Sprites
var (
    //go:embed items.json 
    Json []byte
    //go:embed items.png
    Png []byte
)
