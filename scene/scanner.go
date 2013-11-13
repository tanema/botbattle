package scene

import (
	"github.com/vova616/GarageEngine/engine"
)

type Scanner struct {
	engine.BaseComponent
  Player *BotController
}

func NewScanner() *Scanner {
	return &Scanner{BaseComponent: engine.NewComponent()}
}

func (ms *Scanner) OnComponentAdd() {
}

func (ms *Scanner) OnHit(enemy *engine.GameObject, damager *DamageDealer) {
  player, exists := Players[ms.GameObject().Tag]
  if exists {
    player.OnScan(enemy.Name(), enemy.Transform().WorldPosition())
  }
}

func (ms *Scanner) OnDie(byTimer bool) {
	ms.GameObject().Destroy()
}
