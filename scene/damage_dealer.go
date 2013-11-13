package scene

import (
	"github.com/vova616/GarageEngine/engine"
)

type DamageDealer struct {
	engine.BaseComponent
	Damage float32
}

func NewDamageDealer(dmg float32) *DamageDealer {
	return &DamageDealer{BaseComponent: engine.NewComponent(), Damage: dmg}
}