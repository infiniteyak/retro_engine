package shape_courier_entity

import (
	//gMath "math"
	//"math/rand"
	//"strconv"

	"github.com/infiniteyak/retro_engine/engine/component"
	"github.com/infiniteyak/retro_engine/engine/entity"
	//"github.com/infiniteyak/retro_engine/engine/event"
	"github.com/infiniteyak/retro_engine/engine/layer"
	"github.com/infiniteyak/retro_engine/engine/utility"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	//"github.com/yohamta/donburi/features/math"
    "math"
)

const (
    wallSpriteHeight = 8.0
    wallSpriteWidth = 8.0
    initialMazeOffsetX = 3.0
    initialMazeOffsetY = 15.0
)

type WallType int
const (
    Undefined_walltype WallType = iota
    Solid_walltype
    Empty_walltype
)

type GridSpace struct {
    Center utility.Point
    Solid bool
}

type MazeData struct {
    ecs *ecs.ECS
    entry *donburi.Entry
    entity *donburi.Entity

    //position component.PositionData
    view component.ViewData

    grid [][]GridSpace
    gridRows int
    gridColumns int

    StartR int
    StartC int
}

func (this *MazeData) GetStartPosition() (float64, float64) {
    return this.grid[this.StartR][this.StartC].Center.X, this.grid[this.StartR][this.StartC].Center.Y
}

type Direction int
const (
    North_direction Direction = iota
    East_direction
    South_direction
    West_direction
)

//given a point, which grid space is that in
func (this *MazeData) findCoordinates(pos utility.Point) (int, int) {
    /*
    r_old := int((pos.Y - initialMazeOffset)/wallSpriteHeight)
    c_old := int((pos.X - initialMazeOffset)/wallSpriteWidth)
    */

    //println("c map: ", pos.X, pos.Y, r, c)
    var r int
    var c int
    for r = 0; r < this.gridRows; r++ {
        if math.Abs(this.grid[r][0].Center.Y - pos.Y) <= wallSpriteHeight/2.0 {
            break
        }
    }
    for c = 0; c < this.gridColumns; c++ {
        if math.Abs(this.grid[r][c].Center.X - pos.X) <= wallSpriteWidth/2.0 {
            break
        }
    }
    /*
    if r_old != r || c_old != c {
        println("Bad comp: ",r_old, c_old, r, c)
    }
    */

    return r, c
}

//given a position and move speed, return a new position
func (this *MazeData) ResolveMove(pos utility.Point, dir Direction, speed float64) utility.Point {
    r, c := this.findCoordinates(pos)

    //correct direction if needed
    turning := false
    switch dir {
    case East_direction:
        if !this.grid[r][c+1].Solid && pos.Y != this.grid[r][c].Center.Y {
            turning = true
            if pos.Y > this.grid[r][c].Center.Y {
                dir = North_direction
            } else {
                dir = South_direction
            }
        } else if this.grid[r][c+1].Solid && 
                  !this.grid[r+1][c+1].Solid &&
                  pos.Y >= this.grid[r][c].Center.Y {
            dir = South_direction
        } else if this.grid[r][c+1].Solid && 
                  !this.grid[r-1][c+1].Solid &&
                  pos.Y <= this.grid[r][c].Center.Y {
            dir = North_direction
        } 
    case West_direction:
        if !this.grid[r][c-1].Solid && pos.Y != this.grid[r][c].Center.Y {
            turning = true
            if pos.Y > this.grid[r][c].Center.Y {
                dir = North_direction
            } else {
                dir = South_direction
            }
        } else if this.grid[r][c-1].Solid && 
                  !this.grid[r+1][c-1].Solid &&
                  pos.Y >= this.grid[r][c].Center.Y {
            dir = South_direction
        } else if this.grid[r][c-1].Solid && 
                  !this.grid[r-1][c-1].Solid &&
                  pos.Y <= this.grid[r][c].Center.Y {
            dir = North_direction
        } 
    case South_direction:
        if !this.grid[r+1][c].Solid && pos.X != this.grid[r][c].Center.X {
            turning = true
            if pos.X > this.grid[r][c].Center.X {
                dir = West_direction
            } else {
                dir = East_direction
            }
        } else if this.grid[r+1][c].Solid && 
                  !this.grid[r+1][c-1].Solid &&
                  pos.X <= this.grid[r][c].Center.X {
            dir = West_direction
        } else if this.grid[r+1][c].Solid && 
                  !this.grid[r+1][c+1].Solid &&
                  pos.X >= this.grid[r][c].Center.X {
            dir = East_direction
        } 
    case North_direction:
        if !this.grid[r-1][c].Solid && pos.X != this.grid[r][c].Center.X {
            turning = true
            if pos.X > this.grid[r][c].Center.X {
                dir = West_direction
            } else {
                dir = East_direction
            }
        } else if this.grid[r-1][c].Solid && 
                  !this.grid[r-1][c-1].Solid &&
                  pos.X <= this.grid[r][c].Center.X {
            dir = West_direction
        } else if this.grid[r-1][c].Solid && 
                  !this.grid[r-1][c+1].Solid &&
                  pos.X >= this.grid[r][c].Center.X {
            dir = East_direction
        } 
    }
    //TODO use switches?
    if dir == East_direction {
        if this.grid[r][c].Center.X <= pos.X { //moving out of space
            if this.grid[r][c+1].Solid { //TODO add some magic to make it easier to turn
                return utility.Point{X:this.grid[r][c].Center.X, Y:pos.Y}
            } else {
                return utility.Point{X:pos.X+speed, Y:pos.Y}
            }
        } else if this.grid[r][c].Center.X > pos.X { //moving into space
            if pos.X + speed > this.grid[r][c].Center.X { //moving past center
                if this.grid[r][c+1].Solid || turning { //next space is solid so stop
                    return utility.Point{X:this.grid[r][c].Center.X, Y:pos.Y}
                } else {
                    return utility.Point{X:pos.X+speed, Y:pos.Y}
                }
            } else {
                return utility.Point{X:pos.X+speed, Y:pos.Y}
            }
        }
    } else if dir == West_direction {
        if this.grid[r][c].Center.X >= pos.X { //moving out of space
            if this.grid[r][c-1].Solid { //TODO add some magic to make it easier to turn
                return utility.Point{X:this.grid[r][c].Center.X, Y:pos.Y}
            } else {
                return utility.Point{X:pos.X-speed, Y:pos.Y}
            }
        } else if this.grid[r][c].Center.X < pos.X { //moving into space
            if pos.X - speed < this.grid[r][c].Center.X { //moving past center
                if this.grid[r][c-1].Solid || turning { //next space is solid so stop
                    return utility.Point{X:this.grid[r][c].Center.X, Y:pos.Y}
                } else {
                    return utility.Point{X:pos.X-speed, Y:pos.Y}
                }
            } else {
                return utility.Point{X:pos.X-speed, Y:pos.Y}
            }
        }
    } else if dir == North_direction {
        if this.grid[r][c].Center.Y >= pos.Y { //moving out of space
            if this.grid[r-1][c].Solid { //TODO add some magic to make it easier to turn
                return utility.Point{X:pos.X, Y:this.grid[r][c].Center.Y}
            } else {
                return utility.Point{X:pos.X, Y:pos.Y-speed}
            }
        } else if this.grid[r][c].Center.Y < pos.Y { //moving into space
            if pos.Y - speed < this.grid[r][c].Center.Y { //moving past center
                if this.grid[r-1][c].Solid || turning { //next space is solid so stop
                    return utility.Point{X:pos.X, Y:this.grid[r][c].Center.Y}
                } else {
                    return utility.Point{X:pos.X, Y:pos.Y-speed}
                }
            } else {
                return utility.Point{X:pos.X, Y:pos.Y-speed}
            }
        }
    } else if dir == South_direction {
        if this.grid[r][c].Center.Y <= pos.Y { //moving out of space
            if this.grid[r+1][c].Solid { //TODO add some magic to make it easier to turn
                return utility.Point{X:pos.X, Y:this.grid[r][c].Center.Y}
            } else {
                return utility.Point{X:pos.X, Y:pos.Y+speed}
            }
        } else if this.grid[r][c].Center.Y > pos.Y { //moving into space
            if pos.Y + speed < this.grid[r][c].Center.Y { //moving past center
                if this.grid[r+1][c].Solid || turning { //next space is solid so stop
                    return utility.Point{X:pos.X, Y:this.grid[r][c].Center.Y}
                } else {
                    return utility.Point{X:pos.X, Y:pos.Y+speed}
                }
            } else {
                return utility.Point{X:pos.X, Y:pos.Y+speed}
            }
        }
    }
    return pos
}

func AddMaze( ecs *ecs.ECS,
              x, y float64,
              view *utility.View) *MazeData {
    this := &MazeData{}
    this.ecs = ecs
    
    //start_offset := initialMazeOffset
    sw := wallSpriteWidth
    sh := wallSpriteHeight
    this.gridColumns = 28
    this.gridRows = 31

    this.grid = make([][]GridSpace, this.gridRows)

    pattern := [][]int { //This is identical to the pac man maze
        {1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,},
        {1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,},
        {1, 0, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 0, 1, 1, 0, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 0, 1,},
        {1, 4, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 0, 1, 1, 0, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 4, 1,},
        {1, 0, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 0, 1, 1, 0, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 0, 1,},
        {1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,},
        {1, 0, 1, 1, 1, 1, 0, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 0, 1, 1, 1, 1, 0, 1,},
        {1, 0, 1, 1, 1, 1, 0, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 0, 1, 1, 1, 1, 0, 1,},
        {1, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 1,},
        {1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 0, 1, 1, 0, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1,},
        {1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 0, 1, 1, 0, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1,},
        {1, 1, 1, 1, 1, 1, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 1, 1, 1, 1, 1, 1,},
        {1, 1, 1, 1, 1, 1, 0, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 0, 1, 1, 1, 1, 1, 1,},
        {1, 1, 1, 1, 1, 1, 0, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 0, 1, 1, 1, 1, 1, 1,},
        {1, 9, 3, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 3, 9, 1,},
        {1, 1, 1, 1, 1, 1, 0, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 0, 1, 1, 1, 1, 1, 1,},
        {1, 1, 1, 1, 1, 1, 0, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 0, 1, 1, 1, 1, 1, 1,},
        {1, 1, 1, 1, 1, 1, 0, 1, 1, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 1, 1, 0, 1, 1, 1, 1, 1, 1,},
        {1, 1, 1, 1, 1, 1, 0, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 0, 1, 1, 1, 1, 1, 1,},
        {1, 1, 1, 1, 1, 1, 0, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 0, 1, 1, 1, 1, 1, 1,},
        {1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,},
        {1, 0, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 0, 1, 1, 0, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 0, 1,},
        {1, 0, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 0, 1, 1, 0, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 0, 1,},
        {1, 4, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 4, 1,},
        {1, 1, 1, 0, 1, 1, 0, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 0, 1, 1, 0, 1, 1, 1,},
        {1, 1, 1, 0, 1, 1, 0, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 0, 1, 1, 0, 1, 1, 1,},
        {1, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 1,},
        {1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1,},
        {1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1,},
        {1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,},
        {1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,},
    }

    var curOffsetY float64 = initialMazeOffsetY
    for r := 0; r < this.gridRows; r++ {
        var curOffsetX float64 = initialMazeOffsetX
        this.grid[r] = make([]GridSpace, this.gridColumns)
        for c := 0; c < this.gridColumns; c++ {
            //wType := Undefined_walltype
            if pattern[r][c] == 1 {
                //wType = Solid_walltype
                isSolid := func(y int, x int) bool {
                    if x < 0 || y < 0 || x >= this.gridColumns || y >= this.gridRows {
                        return true
                    }
                    if pattern[y][x] == 1 {
                        return true
                    }
                    return false
                }
                
                north := isSolid(r-1, c)
                east := isSolid(r, c+1)
                south := isSolid(r+1, c)
                west := isSolid(r, c-1)

                tag := ""
                if north {
                    tag = tag + "N"
                }
                if east {
                    tag = tag + "E"
                }
                if south {
                    tag = tag + "S"
                }
                if west {
                    tag = tag + "W"
                }
                if tag == "NESW" {
                    if !isSolid(r-1, c-1) {
                        tag = "NW"
                    } else if !isSolid(r-1, c+1) {
                        tag = "NE"
                    } else if !isSolid(r+1, c+1) {
                        tag = "ES"
                    } else if !isSolid(r+1, c-1) {
                        tag = "SW"
                    } else {
                        tag = ""
                    }
                }
                if tag == "" {
                    tag = "_"
                }

                this.grid[r][c] = GridSpace{
                    Center: utility.Point{
                        X: curOffsetX,
                        Y: curOffsetY,
                    },
                    Solid: true,
                }

                entity.AddSpriteObject(
                    ecs, 
                    layer.Background, 
                    curOffsetX, 
                    curOffsetY, 
                    "Wall", 
                    tag, 
                    view,
                )
            } else {
                this.grid[r][c] = GridSpace{
                    Center: utility.Point{
                        X: curOffsetX,
                        Y: curOffsetY,
                    },
                    Solid: false,
                }
                if pattern[r][c] == 2 {
                    this.StartR = r
                    this.StartC = c
                }
                if pattern[r][c] == 0 { //TODO fix
                    /*
                    entity.AddSpriteObject(
                        ecs, 
                        layer.Background, 
                        curOffsetX, 
                        curOffsetY, 
                        "Items", 
                        "basic_dot", 
                        view,
                    )
                    */
                    AddDot(ecs, curOffsetX, curOffsetY, view)
                }
                if pattern[r][c] == 4 { //TODO fix
                    entity.AddSpriteObject(
                        ecs, 
                        layer.Background, 
                        curOffsetX, 
                        curOffsetY, 
                        "Items", 
                        "dot", 
                        view,
                    )
                }
                if pattern[r][c] == 3 {
                    entity.AddSpriteObject(
                        ecs, 
                        layer.Background, 
                        curOffsetX, 
                        curOffsetY, 
                        "Items", 
                        "teleporter", 
                        view,
                    )
                }
            } 
            curOffsetX += sw
        }
        curOffsetY += sh
    }

    //entity.AddSpriteObject(ecs, layer.Background, x, y, "Wall", "NES", view)

    return this
}
