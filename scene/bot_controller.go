package scene

import (
	"github.com/vova616/GarageEngine/engine"
	"github.com/vova616/GarageEngine/engine/input"
	"math"
	"math/rand"
  "time"
)

type BotController struct {
	engine.BaseComponent                    `json:"-"`
	Missle               *Missle            `json:"-"`
	HPBar                *engine.GameObject `json:"-"`
	Destoyable           *Destoyable        `json:"-"`
	lastShoot            time.Time          `json:"-"`
}

func NewBotController() *BotController {
  return &BotController{engine.NewComponent(), nil, nil, nil, time.Now()}
}

func (sp *BotController) OnComponentAdd() {
	sp.GameObject().AddComponent(engine.NewPhysicsCircle(false))
}

func (sp *BotController) Start() {
	ph := sp.GameObject().Physics
	ph.Body.SetMass(50)
	ph.Shape.Group = 1
	sp.Destoyable = sp.GameObject().ComponentTypeOf(sp.Destoyable).(*Destoyable)
	sp.OnHit(nil, nil)
}

func (sp *BotController) OnHit(enemey *engine.GameObject, damager *DamageDealer) {
  println("on hit")
	if sp.HPBar != nil && sp.Destoyable != nil {
		hp := (float32(sp.Destoyable.HP) / float32(sp.Destoyable.FullHP)) * 100
		s := sp.HPBar.Transform().Scale()
		s.X = hp
		sp.HPBar.Transform().SetScale(s)
	}
}

func (sp *BotController) OnDie(byTimer bool) {
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
	sp.GameObject().Destroy()
}

func (sp *BotController) Shoot() {
	if sp.Missle != nil {
		a := sp.Transform().Rotation()

    pos := engine.Vector{0, 0, 0}
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
    nfire.Physics.Body.AddForce(s.X*3000, s.Y*3000)

    nfire.Physics.Shape.Group = 1
    nfire.Physics.Body.SetMoment(engine.Inf)
    nfire.Transform().SetRotationf(180 - angle)
	}
}

func (sp *BotController) Update() {
  Speed := float32(500000)
	rotSpeed := float32(250)

	delta := float32(engine.DeltaTime())
	r2 := sp.Transform().DirectionTransform(engine.Up)
	ph := sp.GameObject().Physics
	rx, ry := r2.X*delta, r2.Y*delta

	if input.KeyDown('W') {
		ph.Body.AddForce(Speed*rx, Speed*ry)
	}
	if input.KeyDown('S') {
		ph.Body.AddForce(-Speed*rx, -Speed*ry)
	}
	r := sp.Transform().Rotation()
  if input.KeyDown('D') {
    sp.Transform().SetRotationf(r.Z - rotSpeed*delta)
  }
	if input.KeyDown('A') {
		sp.Transform().SetRotationf(r.Z + rotSpeed*delta)
  }
  if input.KeyDown('G') {
    sp.OnDie(false)
  }

	if input.MouseDown(input.MouseLeft) {
		if time.Now().After(sp.lastShoot) {
			sp.Shoot()
			sp.lastShoot = time.Now().Add(time.Millisecond * 200)
		}
	}
}
