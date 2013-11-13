package scene

import (
	"github.com/vova616/GarageEngine/engine"
	"math"
	"math/rand"
  "time"
)

type BotController struct {
	engine.BaseComponent
  Name                 string
	Missle               *Missle
	Health               *engine.GameObject
	HPBar                *engine.GameObject
	Destoyable           *Destoyable
	lastShoot            time.Time
}

func NewBotController(name string, health, healthbar *engine.GameObject, missle *Missle) *BotController {
  return &BotController{engine.NewComponent(), name, missle, health, healthbar, nil, time.Now()}
}
func (sp *BotController) Start() {
	sp.Destoyable = sp.GameObject().ComponentTypeOf(sp.Destoyable).(*Destoyable)
}

func (sp *BotController) OnHit(enemey *engine.GameObject, damager *DamageDealer) {
	if sp.HPBar != nil && sp.Destoyable != nil {
		hp := (float32(sp.Destoyable.HP) / float32(sp.Destoyable.FullHP)) * 100
		s := sp.HPBar.Transform().Scale()
		s.X = hp
		sp.HPBar.Transform().SetScale(s)
	}
}

func (sp *BotController) OnDie(byTimer bool) {
  if sp.GameObject() == nil {
    return
  }

	for i := 0; i < 20; i++ {
		n := Explosion.Clone()
		n.Transform().SetParent2(MainSceneGeneral.Layer1)
		n.Transform().SetWorldPosition(sp.Transform().WorldPosition())
		s := n.Transform().Scale()
		n.Transform().SetScale(s.Mul2(rand.Float32() * 8))
		n.AddComponent(engine.NewPhysics(false))

		n.Transform().SetRotationf(rand.Float32() * 360)
		rot := n.Transform().Direction()
		n.Physics.Body.SetVelocity(-rot.X*100, -rot.Y*100)

		n.Physics.Body.SetMass(1)
		n.Physics.Shape.Group = 1
		n.Physics.Shape.IsSensor = true
	}
  sp.Health.Destroy()
	sp.GameObject().Destroy()
  delete(Players, sp.Name)
}

func (sp *BotController) Shoot() {
	if sp.Missle != nil && time.Now().After(sp.lastShoot) {
		a := sp.Transform().Rotation()

    pos := engine.Vector{0, 37, 0}
    s := sp.Transform().DirectionTransform(engine.Vector{0,1,0})

    p := sp.Transform().WorldPosition()
    m := engine.Identity()
    m.Translate(pos.X, pos.Y, pos.Z)
    m.RotateZ(a.Z, -1)
    m.Translate(p.X, p.Y, p.Z)
    p = m.Translation()

    nfire := sp.Missle.GameObject().Clone()
    nfire.Transform().SetParent2(MainSceneGeneral.Layer1)
    nfire.Transform().SetWorldPosition(p)
    nfire.Physics.Body.IgnoreGravity = true
    nfire.Physics.Body.SetMass(0.1)
    nfire.Tag = MissleTag

    v := sp.GameObject().Physics.Body.Velocity()
    angle := float32(math.Atan2(float64(s.X), float64(s.Y))) * engine.DegreeConst

    nfire.Physics.Body.SetVelocity(float32(v.X), float32(v.Y))
    nfire.Physics.Body.AddForce(s.X*5000, s.Y*5000)

    nfire.Physics.Shape.Group = 1
    nfire.Physics.Body.SetMoment(engine.Inf)
    nfire.Transform().SetRotationf(180 - angle)

    sp.lastShoot = time.Now().Add(time.Millisecond * 200)
	}
}
