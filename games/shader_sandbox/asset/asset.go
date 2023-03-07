package shader_sb_assets

import (
    "github.com/infiniteyak/retro_engine/engine/asset"
    white_font "github.com/infiniteyak/retro_engine/engine/asset/white_font"
)

func InitAssets() {
    InitSpriteAssets()
}

func InitSpriteAssets() {
    //Shared assets
    asset.LoadSpriteAsset("WhiteFont", white_font.Json, white_font.Png)
}
