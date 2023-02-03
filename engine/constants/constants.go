package constants


const (
    // For now I'm  going with resolutions in 16:9 ratio that yield roughly the
    // same area as original arcade resolution. (Tate means vertical mode.)

    // First era of arcade games, but at 16:9
    ArcadeAHeight = 180
    ArcadeAWidth = 320
    ArcadeATateHeight = 320
    ArcadeATateWidth = 180

    // First era of arcade games, but at 9:7
    ArcadeBHeight = 224
    ArcadeBWidth = 288
    ArcadeBTateHeight = 288
    ArcadeBTateWidth = 224

    // Engine defaults
    MaxTPS = 120 //TODO should this be 60?
)
