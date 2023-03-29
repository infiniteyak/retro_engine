package component

import (
	"github.com/hajimehoshi/ebiten/v2"
    "github.com/yohamta/donburi"
)

type InputTypeId int

const (
    Undefined_inputtypeid InputTypeId = iota
    Continuous_inputtypeid //triggers as often as it can
    Limited_inputtypeid //once per button press
    Hybrid_inputtypeid //once per button press OR at a certain frequency once held
)

type InputsData struct {
    KeyMap map[ActionId]ebiten.Key
    TypeMap map[ActionId]InputTypeId
    DelayMap map[ActionId]int
    FrequencyMap map[ActionId]int
}

func (this *InputsData) AddContinuousInput(actionId ActionId, key ebiten.Key) {
    this.KeyMap[actionId] = key
    this.TypeMap[actionId] = Continuous_inputtypeid
}

func (this *InputsData) AddLimitedInput(actionId ActionId, key ebiten.Key) {
    this.KeyMap[actionId] = key
    this.TypeMap[actionId] = Limited_inputtypeid
}

func (this *InputsData) AddHybridInput(actionId ActionId, key ebiten.Key, delay int, frequency int) {
    this.KeyMap[actionId] = key
    this.TypeMap[actionId] = Hybrid_inputtypeid
    this.DelayMap[actionId] = delay
    this.FrequencyMap[actionId] = frequency
}

func NewInput() InputsData {
    km := make(map[ActionId]ebiten.Key)
    tym := make(map[ActionId]InputTypeId)
    dm := make(map[ActionId]int)
    fm := make(map[ActionId]int)
    return InputsData {
        KeyMap: km,
        TypeMap: tym,
        DelayMap: dm,
        FrequencyMap: fm,
    }
}

var Inputs = donburi.NewComponentType[InputsData]()
