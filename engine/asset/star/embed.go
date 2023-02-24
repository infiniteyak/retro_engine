package asset_star

import (
    _ "embed"
)

var (
    //go:embed star.json 
    Json []byte
    //go:embed star.png
    Png []byte
)
