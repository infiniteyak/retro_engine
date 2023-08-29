package asset_wave

import (
    _ "embed"
)

var (
    //go:embed wave.wav
    Wav []byte
)
