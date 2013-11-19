package scene

import (
	"github.com/vova616/GarageEngine/engine"
	"time"
)

type Destoyable struct {
	engine.BaseComponent
	Alive           bool
	HP              float32
	FullHP          float32
	destoyableFuncs DestoyableFuncs

	createTime    time.Time
	aliveDuration time.Duration
	autoDestory   bool
}

func NewDestoyable(hp float32) *Destoyable {
	return &Destoyable{BaseComponent: engine.NewComponent(), FullHP: hp, Alive: true, HP: hp}
}

type DestoyableFuncs interface {
	OnDie(byTimer bool)
	OnHit(*engine.GameObject, *DamageDealer)
}

func (ds *Destoyable) Start() {
	ds.createTime = time.Now()
	ds.destoyableFuncs, _ = ds.GameObject().ComponentImplements(&ds.destoyableFuncs).(DestoyableFuncs)
}

func (ds *Destoyable) SetDestroyTime(sec float32) {
	ds.autoDestory = true
	ds.aliveDuration = time.Millisecond * time.Duration(1000*sec)
}

func (ds *Destoyable) Update() {
	if ds.autoDestory && ds.GameObject() != nil {
		if time.Now().After(ds.createTime.Add(ds.aliveDuration)) {
			if ds.destoyableFuncs != nil {
				ds.destoyableFuncs.OnDie(true)
			} else {
				ds.GameObject().Destroy()
			}
		}
	}
}

func (ds *Destoyable) OnCollisionEnter(arbiter engine.Arbiter) bool {
	if !ds.Alive {
		return true
	}
	var dmg *DamageDealer = nil
	var enemyDestoyable *Destoyable

	enemy := arbiter.GameObjectB()

	if enemy == nil {
		return true
	}

	dmg, _ = enemy.ComponentTypeOf(dmg).(*DamageDealer)
	enemyDestoyable, _ = enemy.ComponentTypeOf(enemyDestoyable).(*Destoyable)

	if enemyDestoyable == nil {
		return true
	}

	if dmg != nil {
		ds.HP -= dmg.Damage
	}
	if ds.destoyableFuncs != nil {
		ds.destoyableFuncs.OnHit(enemy, dmg)
	}

	if ds.HP <= 0 {
		ds.Alive = false
		if ds.destoyableFuncs != nil {
			ds.destoyableFuncs.OnDie(false)
		} else {
			ds.GameObject().Destroy()
		}
	}

	return true
}
