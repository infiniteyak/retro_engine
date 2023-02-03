package utility

type Point struct {
    X float64 
    Y float64 
} 

type Rectangle struct {
    Min Point
    Max Point
} 

type View struct {
    Area Rectangle
    Offset Point
} 

func NewView(offsetX, offsetY, maxX, maxY float64) *View {
    return &View{
        Area: Rectangle {
            Min: Point {
                X: 0.0,
                Y: 0.0,
            },
            Max: Point {
                X: maxX,
                Y: maxY,
            },
        },
        Offset: Point{
            X: offsetX, 
            Y: offsetY,
        },
    }
}
