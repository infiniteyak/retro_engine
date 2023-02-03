package component

import (
    "github.com/yohamta/donburi"
)

type TextInputData struct {
    String *string
    Length int
    Function func()
}

var TextInput = donburi.NewComponentType[TextInputData]()
