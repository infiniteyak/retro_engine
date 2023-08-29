package shape_courier_entity

import (
	"math/rand"
	"github.com/infiniteyak/retro_engine/engine/component"
	"github.com/infiniteyak/retro_engine/engine/entity"
	"github.com/infiniteyak/retro_engine/engine/layer"
	"github.com/infiniteyak/retro_engine/engine/utility"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
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
    view component.ViewData
    grid [][]GridSpace
    gridRows int
    gridColumns int
    StartR int
    StartC int
    SpawnR int
    SpawnC int
}

func (this *MazeData) GetNWCorner() utility.Point {
    return this.grid[0][0].Center
}

func (this *MazeData) GetNECorner() utility.Point {
    return this.grid[0][this.gridColumns-1].Center
}

func (this *MazeData) GetSECorner() utility.Point {
    return this.grid[this.gridRows-1][this.gridColumns-1].Center
}

func (this *MazeData) GetSWCorner() utility.Point {
    return this.grid[this.gridRows-1][0].Center
}

func (this *MazeData) GetStartPosition() (float64, float64) {
    return this.grid[this.StartR][this.StartC].Center.X, this.grid[this.StartR][this.StartC].Center.Y
}

func (this *MazeData) GetSpawnPosition() (float64, float64) {
    return this.grid[this.SpawnR][this.SpawnC].Center.X, this.grid[this.SpawnR][this.SpawnC].Center.Y
}

type Direction int
const (
    Undefined_direction Direction = iota
    South_direction
    East_direction
    North_direction
    West_direction
)

func distance(a, b utility.Point) float64 { //TODO move to utility
    return math.Sqrt(math.Pow(a.X - b.X, 2) + math.Pow(a.Y - b.Y, 2))
}

func (this *MazeData) GetRandomDirection(pos utility.Point, curDir Direction) Direction {
    r, c := this.FindCoordinates(pos)
    selected := Undefined_direction
    backwards := Undefined_direction
    switch curDir {
    case North_direction:
        backwards = South_direction
    case South_direction:
        backwards = North_direction
    case East_direction:
        backwards = West_direction
    case West_direction:
        backwards = East_direction
    }
    options := []Direction{}
    if !this.grid[r-1][c].Solid && curDir != South_direction {
        options = append(options, North_direction)
    }
    if !this.grid[r+1][c].Solid && curDir != North_direction {
        options = append(options, South_direction)
    }
    if !this.grid[r][c-1].Solid && curDir != East_direction {
        options = append(options, West_direction)
    }
    if !this.grid[r][c+1].Solid && curDir != West_direction {
        options = append(options, East_direction)
    }
    if len(options) > 0 {
        selected = options[rand.Intn(len(options))]
    }
    if selected == Undefined_direction {
        selected = backwards
    }
    return selected
}

// what direction is closest to target, avoiding turning around completely
func (this *MazeData) GetDirectionToTarget(pos utility.Point, target utility.Point, curDir Direction) Direction {
    r, c := this.FindCoordinates(pos)
    bestDst := -1.0
    selected := Undefined_direction
    backwards := Undefined_direction
    switch curDir {
    case North_direction:
        backwards = South_direction
    case South_direction:
        backwards = North_direction
    case East_direction:
        backwards = West_direction
    case West_direction:
        backwards = East_direction
    }
    if !this.grid[r-1][c].Solid && curDir != South_direction {
        newDst := distance(this.grid[r-1][c].Center, target)
        if bestDst < 0 || newDst < bestDst {
            selected = North_direction
            bestDst = newDst
        }
    }
    if !this.grid[r+1][c].Solid && curDir != North_direction {
        newDst := distance(this.grid[r+1][c].Center, target)
        if bestDst < 0 || newDst < bestDst {
            selected = South_direction
            bestDst = newDst
        }
    }
    if !this.grid[r][c-1].Solid && curDir != East_direction {
        newDst := distance(this.grid[r][c-1].Center, target)
        if bestDst < 0 || newDst < bestDst {
            selected = West_direction
            bestDst = newDst
        }
    }
    if !this.grid[r][c+1].Solid && curDir != West_direction {
        newDst := distance(this.grid[r][c+1].Center, target)
        if bestDst < 0 || newDst < bestDst {
            selected = East_direction
            bestDst = newDst
        }
    }
    //TODO what if selected dir is undefined?
    if selected == Undefined_direction {
        selected = backwards
    }
    return selected
}

//given a point, which grid space is that in
func (this *MazeData) FindCoordinates(pos utility.Point) (int, int) {
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

    return r, c
}

func (this *MazeData) findOpenSpaceDir(r, c int, 
                                    dir Direction, 
                                    breakPositive bool) Direction {
    distance := 1
    var a int
    var b int
    var bLimit int
    var forwardMod int
    var gridAccess func(int, int) GridSpace
    var posDir Direction
    var negDir Direction
    switch dir {
    case North_direction:
        bLimit = this.gridColumns
        a = r
        b = c
        forwardMod = -1
        gridAccess = func(a int, b int) GridSpace {
            return this.grid[a][b]
        }
        posDir = East_direction
        negDir = West_direction
    case South_direction:
        bLimit = this.gridColumns
        a = r
        b = c
        forwardMod = 1
        gridAccess = func(a int, b int) GridSpace {
            return this.grid[a][b]
        }
        posDir = East_direction
        negDir = West_direction
    case East_direction:
        bLimit = this.gridRows
        a = c
        b = r
        forwardMod = 1
        gridAccess = func(a int, b int) GridSpace {
            return this.grid[b][a]
        }
        posDir = South_direction
        negDir = North_direction
    case West_direction:
        bLimit = this.gridRows
        a = c
        b = r
        forwardMod = -1
        gridAccess = func(a int, b int) GridSpace {
            return this.grid[b][a]
        }
        posDir = South_direction
        negDir = North_direction
    }
    disqPos := false
    disqNeg := false
    for {
        if b+distance >= bLimit || gridAccess(a, b+distance).Solid { // would hit a wall
            disqPos = true
        }
        if b-distance < 0 || gridAccess(a, b-distance).Solid { // would hit a wall
            disqNeg = true
        }
        if disqPos && disqNeg { // just keep going in the same dir
            return Undefined_direction
        }

        foundPos := !disqPos && !gridAccess(a+forwardMod, b+distance).Solid //east
        foundNeg := !disqNeg && !gridAccess(a+forwardMod, b-distance).Solid // west

        if foundPos && foundNeg {
            if breakPositive {
                return posDir
            } else {
                return negDir
            }
        } else if foundPos {
            return posDir
        } else if foundNeg {
            return negDir
        }
        distance += 1
    }
    return Undefined_direction //TODO raise exception?
}

func (this *MazeData) PickDirection(pos utility.Point, 
                                    dirA, dirB Direction) Direction {
    r, c := this.FindCoordinates(pos)

    var aSolid bool

    switch dirA {
    case North_direction:
        aSolid = this.grid[r-1][c].Solid
    case South_direction:
        aSolid = this.grid[r+1][c].Solid
    case West_direction:
        aSolid = this.grid[r][c-1].Solid
    case East_direction:
        aSolid = this.grid[r][c+1].Solid
    }

    if aSolid {
        return dirA
    } else {
        return dirB
    }
}

func (this *MazeData) ResolveMove(pos utility.Point, 
                                  dir Direction, 
                                  speed float64) (utility.Point, Direction) {
    r, c := this.FindCoordinates(pos)
    
    //First check if we are moving into a wall, and need to turn
    turning := false

    moveDir := dir
    var forwardSpaceIsSolid bool
    switch dir {
    case North_direction:
        forwardSpaceIsSolid = this.grid[r-1][c].Solid
    case South_direction:
        forwardSpaceIsSolid = this.grid[r+1][c].Solid
    case West_direction:
        forwardSpaceIsSolid = this.grid[r][c-1].Solid
    case East_direction:
        forwardSpaceIsSolid = this.grid[r][c+1].Solid
    }

    var aligned bool
    var breakPos bool
    var leaving bool
    if dir == North_direction || dir == South_direction {
        aligned = pos.X == this.grid[r][c].Center.X

        breakPos = pos.X > this.grid[r][c].Center.X
        if dir == North_direction {
            leaving = pos.Y <= this.grid[r][c].Center.Y
        } else {
            leaving = pos.Y >= this.grid[r][c].Center.Y
        }
        if breakPos {
            moveDir = West_direction
        } else {
            moveDir = East_direction
        }
    } else if dir == West_direction || dir == East_direction {
        aligned = pos.Y == this.grid[r][c].Center.Y

        breakPos = pos.Y > this.grid[r][c].Center.Y
        //leaving = (dir == West_direction) == breakPos
        if dir == West_direction {
            leaving = pos.X <= this.grid[r][c].Center.X
        } else {
            leaving = pos.X >= this.grid[r][c].Center.X
        }
        if breakPos {
            moveDir = North_direction
        } else {
            moveDir = South_direction
        }
    }

    //TODO I think there's a bug here where if we are moving fast enough
    // you can overshoot. I think this would happen if the move speed is
    // greater than (or equal to?) half the height of a space.
    // That would be so fast the game would be unplayable though (at least
    // at the current resolution). So I'm going to ignore it for now.
    var stop bool
    if !forwardSpaceIsSolid && !aligned { //move towards center of current space
        turning = true //meaning we need to stop at the center point
        dir = moveDir // TODO do I really need movedir?
    } else if forwardSpaceIsSolid { // find direction to nearest open space
        openDir := this.findOpenSpaceDir(r, c, dir, breakPos)
        if openDir != Undefined_direction && leaving {
            dir = openDir
        } else { //there's no where to go so stop at center point
            stop = true
        }
    }
    if leaving && stop { // we can't go farther in this direction, so stop
        if dir == North_direction || dir == South_direction {
            return utility.Point{X:pos.X, Y:this.grid[r][c].Center.Y}, dir
        } else {
            return utility.Point{X:this.grid[r][c].Center.X, Y:pos.Y}, dir
        }
    }
    if (!leaving && stop) || turning { //move up to the center but no farther
        switch dir {
        case North_direction:
            if pos.Y - speed < this.grid[r][c].Center.Y { //moving past the center
                return utility.Point{X:pos.X, Y:this.grid[r][c].Center.Y}, dir
            } else { // didn't reach the center
                return utility.Point{X:pos.X, Y:pos.Y-speed}, dir
            }
        case South_direction:
            if pos.Y + speed > this.grid[r][c].Center.Y { //moving past the center
                return utility.Point{X:pos.X, Y:this.grid[r][c].Center.Y}, dir
            } else { // didn't reach the center
                return utility.Point{X:pos.X, Y:pos.Y+speed}, dir
            }
        case West_direction:
            if pos.X - speed < this.grid[r][c].Center.X { //moving past the center
                return utility.Point{X:this.grid[r][c].Center.X, Y:pos.Y}, dir
            } else { // didn't reach the center
                return utility.Point{X:pos.X-speed, Y:pos.Y}, dir
            }
        case East_direction:
            if pos.X + speed > this.grid[r][c].Center.X { //moving past the center
                return utility.Point{X:this.grid[r][c].Center.X, Y:pos.Y}, dir
            } else { // didn't reach the center
                return utility.Point{X:pos.X+speed, Y:pos.Y}, dir
            }
        }
    }

    switch dir {
    case North_direction:
        return utility.Point{X:pos.X, Y:pos.Y-speed}, dir
    case South_direction:
        return utility.Point{X:pos.X, Y:pos.Y+speed}, dir
    case West_direction:
        return utility.Point{X:pos.X-speed, Y:pos.Y}, dir
    case East_direction:
        return utility.Point{X:pos.X+speed, Y:pos.Y}, dir
    }
    
    return utility.Point{X:pos.X, Y:pos.Y}, dir // TODO exception?
}

func AddMaze( ecs *ecs.ECS,
              x, y float64,
              view *utility.View) *MazeData {
    this := &MazeData{}
    this.ecs = ecs
    
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
        {1, 1, 1, 1, 1, 1, 0, 1, 1, 0, 0, 0, 0, 7, 0, 0, 0, 0, 0, 1, 1, 0, 1, 1, 1, 1, 1, 1,},
        {1, 1, 1, 1, 1, 1, 0, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 0, 1, 1, 1, 1, 1, 1,},
        {1, 1, 1, 1, 1, 1, 0, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 0, 1, 1, 1, 1, 1, 1,},
        {1, 9, 5, 6, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 6, 3, 9, 1,},
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

    var teleportAOffsetX float64
    var teleportAOffsetY float64
    var teleportBOffsetX float64
    var teleportBOffsetY float64

    for r := 0; r < this.gridRows; r++ {
        var curOffsetX float64 = initialMazeOffsetX
        this.grid[r] = make([]GridSpace, this.gridColumns)
        for c := 0; c < this.gridColumns; c++ {
            if pattern[r][c] == 1 {
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
                    AddDot(ecs, curOffsetX, curOffsetY, view)
                }
                if pattern[r][c] == 4 { //TODO fix
                    AddPower(ecs, curOffsetX, curOffsetY, view)
                }
                if pattern[r][c] == 3 { // teleporter A
                    teleportAOffsetX = curOffsetX
                    teleportAOffsetY = curOffsetY
                }
                if pattern[r][c] == 5 { // teleporter B
                    teleportBOffsetX = curOffsetX
                    teleportBOffsetY = curOffsetY
                }
                if pattern[r][c] == 6 { // Allow tp
                    AddActionTrigger(
                        ecs, 
                        curOffsetX, 
                        curOffsetY, 
                        component.ReadyTeleport_actionid, 
                        view,
                    )
                }
                if pattern[r][c] == 7 { // ghost
                    this.SpawnR = r
                    this.SpawnC = c
                }
            } 
            curOffsetX += sw
        }
        curOffsetY += sh
    }
    AddTeleporter(
        ecs, 
        teleportAOffsetX, 
        teleportAOffsetY, 
        teleportBOffsetX, 
        teleportBOffsetY, 
        4.0,
        0,
        view,
    )
    AddTeleporter(
        ecs, 
        teleportBOffsetX, 
        teleportBOffsetY, 
        teleportAOffsetX, 
        teleportAOffsetY, 
        -4.0,
        0,
        view,
    )

    return this
}

func AddDecorativeMaze( ecs *ecs.ECS,
              x, y float64,
              view *utility.View) *MazeData {
    this := &MazeData{}
    this.ecs = ecs
    
    sw := wallSpriteWidth
    sh := wallSpriteHeight
    this.gridColumns = 28
    this.gridRows = 33

    this.grid = make([][]GridSpace, this.gridRows)

    pattern := [][]int {
        {0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,},
        {0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,},
        {0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0,},
        {0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0,},
        {0, 0, 0, 0, 0, 0, 1, 1, 9, 9, 9, 9, 4, 1, 1, 4, 9, 9, 9, 9, 1, 1, 0, 0, 0, 0, 0, 0,},
        {0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 9, 1, 1, 9, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0,},
        {0, 0, 0, 1, 1, 0, 1, 1, 1, 1, 1, 1, 9, 1, 1, 9, 1, 1, 1, 1, 1, 1, 0, 1, 1, 0, 0, 0,},
        {0, 0, 0, 1, 1, 9, 9, 9, 9, 9, 1, 1, 9, 1, 1, 9, 1, 1, 9, 9, 9, 9, 9, 1, 1, 0, 0, 0,},
        {0, 0, 0, 1, 1, 1, 1, 1, 1, 9, 1, 1, 9, 1, 1, 9, 1, 1, 9, 1, 1, 1, 1, 1, 1, 0, 0, 0,},
        {0, 0, 0, 1, 1, 1, 1, 1, 1, 9, 1, 1, 9, 1, 1, 9, 1, 1, 9, 1, 1, 1, 1, 1, 1, 0, 0, 0,},
        {0, 0, 0, 9, 9, 9, 9, 9, 9, 4, 1, 1, 9, 1, 1, 9, 1, 1, 4, 9, 9, 9, 9, 9, 9, 0, 0, 0,},
        {0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 9, 1, 1, 9, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0,},
        {0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 9, 1, 1, 9, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0,},
        {0, 9, 1, 1, 4, 9, 9, 9, 9, 9, 9, 9, 9, 1, 1, 9, 9, 9, 9, 9, 9, 9, 9, 4, 1, 1, 9, 0,},
        {1, 9, 1, 1, 9, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 9, 1, 1, 9, 1,},
        {1, 9, 1, 1, 9, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 9, 1, 1, 9, 1,},
        {1, 4, 1, 1, 9, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 9, 1, 1, 4, 1,},
        {1, 9, 1, 1, 9, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 9, 1, 1, 9, 1,},
        {1, 9, 1, 1, 9, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 9, 1, 1, 9, 1,},
        {0, 9, 1, 1, 4, 9, 9, 9, 9, 9, 9, 9, 9, 1, 1, 9, 9, 9, 9, 9, 9, 9, 9, 4, 1, 1, 9, 0,},
        {0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 9, 1, 1, 9, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0,},
        {0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 9, 1, 1, 9, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0,},
        {0, 0, 0, 9, 9, 9, 9, 9, 9, 4, 1, 1, 9, 1, 1, 9, 1, 1, 4, 9, 9, 9, 9, 9, 9, 0, 0, 0,},
        {0, 0, 0, 1, 1, 1, 1, 1, 1, 9, 1, 1, 9, 1, 1, 9, 1, 1, 9, 1, 1, 1, 1, 1, 1, 0, 0, 0,},
        {0, 0, 0, 1, 1, 1, 1, 1, 1, 9, 1, 1, 9, 1, 1, 9, 1, 1, 9, 1, 1, 1, 1, 1, 1, 0, 0, 0,},
        {0, 0, 0, 1, 1, 9, 9, 9, 9, 9, 1, 1, 9, 1, 1, 9, 1, 1, 9, 9, 9, 9, 9, 1, 1, 0, 0, 0,},
        {0, 0, 0, 1, 1, 0, 1, 1, 1, 1, 1, 1, 9, 1, 1, 9, 1, 1, 1, 1, 1, 1, 0, 1, 1, 0, 0, 0,},
        {0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 9, 1, 1, 9, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0,},
        {0, 0, 0, 0, 0, 0, 1, 1, 9, 9, 9, 9, 4, 1, 1, 4, 9, 9, 9, 9, 1, 1, 0, 0, 0, 0, 0, 0,},
        {0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0,},
        {0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0,},
        {0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,},
        {0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,},
    }

    var curOffsetY float64 = initialMazeOffsetY
    for r := 0; r < this.gridRows; r++ {
        var curOffsetX float64 = initialMazeOffsetX
        this.grid[r] = make([]GridSpace, this.gridColumns)
        for c := 0; c < this.gridColumns; c++ {
            if pattern[r][c] == 1 {
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
                if pattern[r][c] == 9 { //TODO fix
                    AddDot(ecs, curOffsetX, curOffsetY, view)
                }
                if pattern[r][c] == 4 { //TODO fix
                    AddPower(ecs, curOffsetX, curOffsetY, view)
                }
            } 
            curOffsetX += sw
        }
        curOffsetY += sh
    }

    return this
}
