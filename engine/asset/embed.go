package asset

import (
    _ "embed"
)

// Sounds
var (
    //go:embed fire.wav
    Fire_wav []byte

    //go:embed jab.wav
    Hit_wav []byte

    //go:embed menu.wav
    Menu_wav []byte

    //go:embed ship_destroy.wav
    Destroyed_wav []byte

    //go:embed thruster.wav
    Thruster_wav []byte

    //go:embed wave.wav
    Wave_wav []byte
)

// Sprites
var (
    //go:embed whitefont.json
    WhiteFont_json []byte
    //go:embed whitefont.png
    WhiteFont_png []byte

    //go:embed lightbluefont.json
    LightBlueFont_json []byte
    //go:embed lightbluefont.png
    LightBlueFont_png []byte

    //go:embed redfont.json
    RedFont_json []byte
    //go:embed redfont.png
    RedFont_png []byte

    //go:embed purplefont.json
    PurpleFont_json []byte
    //go:embed purplefont.png
    PurpleFont_png []byte

    //go:embed blackfont.json
    BlackFont_json []byte
    //go:embed blackfont.png
    BlackFont_png []byte

    //go:embed titlefont.json
    TitleFont_json []byte
    //go:embed titlefont.png
    TitleFont_png []byte

    //go:embed small_asteroid.json 
    SmallAsteroid_json []byte
    //go:embed small_asteroid.png 
    SmallAsteroid_png []byte

    //go:embed medium_asteroid.json 
    MediumAsteroid_json []byte
    //go:embed medium_asteroid.png 
    MediumAsteroid_png []byte

    //go:embed large_asteroid.json 
    LargeAsteroid_json []byte
    //go:embed large_asteroid.png 
    LargeAsteroid_png []byte

    //go:embed small_explosion.json 
    SmallExplosion_json []byte
    //go:embed small_explosion.png 
    SmallExplosion_png []byte

    //go:embed simple_bullet.json 
    SimpleBullet_json []byte
    //go:embed simple_bullet.png
    SimpleBullet_png []byte

    //go:embed star.json 
    Star_json []byte
    //go:embed star.png
    Star_png []byte

    //go:embed player_ship1.json 
    PlayerShip_json []byte
    //go:embed player_ship1.png
    PlayerShip_png []byte

    //go:embed boomerang.json 
    Boomerang_json []byte
    //go:embed boomerang.png
    Boomerang_png []byte

    //go:embed alien_ship1.json 
    Alien1_json []byte
    //go:embed alien_ship1.png
    Alien1_png []byte
)
