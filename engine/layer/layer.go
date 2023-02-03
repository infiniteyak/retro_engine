package layer

import "github.com/yohamta/donburi/ecs"

const (
	Background ecs.LayerID = iota
	Foreground
	HudBackground
	HudForeground
)
