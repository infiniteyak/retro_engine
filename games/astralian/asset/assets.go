package astralian_assets

import (
    "github.com/infiniteyak/retro_engine/engine/asset"
    white_font "github.com/infiniteyak/retro_engine/engine/asset/white_font"
    light_blue_font "github.com/infiniteyak/retro_engine/engine/asset/light_blue_font"
    red_font "github.com/infiniteyak/retro_engine/engine/asset/red_font"
    purple_font "github.com/infiniteyak/retro_engine/engine/asset/purple_font"
    star "github.com/infiniteyak/retro_engine/engine/asset/star"
    scifi_projectile "github.com/infiniteyak/retro_engine/engine/asset/scifi_projectile"
    generic_hit "github.com/infiniteyak/retro_engine/engine/asset/generic_hit"
    menu_noise "github.com/infiniteyak/retro_engine/engine/asset/menu_noise"
    player_ship_destroyed "github.com/infiniteyak/retro_engine/engine/asset/player_ship_destroyed"
    thruster "github.com/infiniteyak/retro_engine/engine/asset/thruster"
    wave "github.com/infiniteyak/retro_engine/engine/asset/wave"
    boomerang "github.com/infiniteyak/retro_engine/games/astralian/asset/boomerang"
    laser "github.com/infiniteyak/retro_engine/games/astralian/asset/laser"
    ship "github.com/infiniteyak/retro_engine/games/astralian/asset/ship"
    alien_a "github.com/infiniteyak/retro_engine/games/astralian/asset/alien_a"
    alien_b "github.com/infiniteyak/retro_engine/games/astralian/asset/alien_b"
    alien_c "github.com/infiniteyak/retro_engine/games/astralian/asset/alien_c"
    alien_d "github.com/infiniteyak/retro_engine/games/astralian/asset/alien_d"
    alien_bullet "github.com/infiniteyak/retro_engine/games/astralian/asset/alien_bullet"
    title_font "github.com/infiniteyak/retro_engine/games/astralian/asset/title_font"
    bg_music "github.com/infiniteyak/retro_engine/games/astralian/asset/bg_music"
)

func InitAssets() {
    InitSpriteAssets()
    InitAudioAssets()
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

    asset.LoadMp3AudioAsset("Music", bg_music.Mp3)
}

func InitSpriteAssets() {
    //Game specific assets
    asset.LoadSpriteAsset("Boomerang", boomerang.Json, boomerang.Png)
    asset.LoadSpriteAsset("Laser", laser.Json, laser.Png)
    asset.LoadSpriteAsset("AstralianShip", ship.Json, ship.Png)
    asset.LoadSpriteAsset("AlienA", alien_a.Json, alien_a.Png)
    asset.LoadSpriteAsset("AlienB", alien_b.Json, alien_b.Png)
    asset.LoadSpriteAsset("AlienC", alien_c.Json, alien_c.Png)
    asset.LoadSpriteAsset("AlienD", alien_d.Json, alien_d.Png)
    asset.LoadSpriteAsset("AlienBullet", alien_bullet.Json, alien_bullet.Png)
    asset.LoadSpriteAsset("TitleFont", title_font.Json, title_font.Png)

    //Shared assets
    asset.LoadSpriteAsset("WhiteFont", white_font.Json, white_font.Png)
    asset.LoadSpriteAsset("LightBlueFont", light_blue_font.Json, light_blue_font.Png)
    asset.LoadSpriteAsset("RedFont", red_font.Json, red_font.Png)
    asset.LoadSpriteAsset("PurpleFont", purple_font.Json, purple_font.Png)
    asset.LoadSpriteAsset("Star", star.Json, star.Png)
}
