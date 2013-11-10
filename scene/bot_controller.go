package scene

import (
	"github.com/vova616/GarageEngine/engine"
)

type BotController struct {
	engine.BaseComponent `json:"-"`
	Missle               *Missle `json:"-"`
	HPBar                *engine.GameObject `json:"-"`
	Destoyable           *Destoyable        `json:"-"`
}

func NewBotController() *BotController {
  return &BotController{engine.NewComponent(), nil, nil, nil}
}

func (sp *BotController) OnComponentAdd() {
	sp.GameObject().AddComponent(engine.NewPhysicsCircle(false))
}

func (sp *BotController) Start() {
}

func (sp *BotController) OnHit(enemey *engine.GameObject, damager *DamageDealer) {
}

func (sp *BotController) OnDie(byTimer bool) {
	sp.GameObject().Destroy()
}

func (sp *BotController) Shoot() {
}

func (sp *BotController) Update() {
}
