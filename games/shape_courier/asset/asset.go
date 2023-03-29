package shape_courier_assets

import (
    "github.com/infiniteyak/retro_engine/engine/asset"
    black_font "github.com/infiniteyak/retro_engine/engine/asset/black_font"
    white_font "github.com/infiniteyak/retro_engine/engine/asset/white_font"
    light_blue_font "github.com/infiniteyak/retro_engine/engine/asset/light_blue_font"
    red_font "github.com/infiniteyak/retro_engine/engine/asset/red_font"
    purple_font "github.com/infiniteyak/retro_engine/engine/asset/purple_font"
    white_select "github.com/infiniteyak/retro_engine/engine/asset/white_select"
    scifi_projectile "github.com/infiniteyak/retro_engine/engine/asset/scifi_projectile"
    generic_hit "github.com/infiniteyak/retro_engine/engine/asset/generic_hit"
    menu_noise "github.com/infiniteyak/retro_engine/engine/asset/menu_noise"
    player_ship_destroyed "github.com/infiniteyak/retro_engine/engine/asset/player_ship_destroyed"
    thruster "github.com/infiniteyak/retro_engine/engine/asset/thruster"
    wave "github.com/infiniteyak/retro_engine/engine/asset/wave"
    wall "github.com/infiniteyak/retro_engine/games/shape_courier/asset/wall"
    space_mandy "github.com/infiniteyak/retro_engine/games/shape_courier/asset/space_mandy"
    items "github.com/infiniteyak/retro_engine/games/shape_courier/asset/items"
    ghost "github.com/infiniteyak/retro_engine/games/shape_courier/asset/ghost"
)

func InitAssets() {
    InitSpriteAssets()
    InitAudioAssets()
    asset.InitPolyAssets()
}

func InitAudioAssets() {
    //Shared assets
    asset.InitAudioAssets()
    asset.LoadWavAudioAsset("SciFiProjectile", scifi_projectile.Wav)
    asset.LoadWavAudioAsset("GenericHit", generic_hit.Wav)
    asset.LoadWavAudioAsset("MenuNoise", menu_noise.Wav)
    asset.LoadWavAudioAsset("PlayerShipDestroyed", player_ship_destroyed.Wav)
    asset.LoadWavAudioAsset("Thruster", thruster.Wav)
    asset.LoadWavAudioAsset("Wave", wave.Wav)
}

func InitSpriteAssets() {
    //Shared assets
    asset.LoadSpriteAsset("TitleFont", black_font.Json, black_font.Png)
    asset.LoadSpriteAsset("WhiteFont", white_font.Json, white_font.Png)
    asset.LoadSpriteAsset("LightBlueFont", light_blue_font.Json, light_blue_font.Png)
    asset.LoadSpriteAsset("RedFont", red_font.Json, red_font.Png)
    asset.LoadSpriteAsset("PurpleFont", purple_font.Json, purple_font.Png)
    asset.LoadSpriteAsset("WhiteSelect", white_select.Json, white_select.Png)
    asset.LoadSpriteAsset("Wall", wall.Json, wall.Png)
    asset.LoadSpriteAsset("SpaceMandy", space_mandy.Json, space_mandy.Png)
    asset.LoadSpriteAsset("Items", items.Json, items.Png)
    asset.LoadSpriteAsset("Ghost", ghost.Json, ghost.Png)
}
