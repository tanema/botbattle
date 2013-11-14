package scene

import (
	"github.com/vova616/GarageEngine/engine"
	"github.com/vova616/GarageEngine/engine/components"
	"github.com/vova616/chipmunk"
	"github.com/vova616/chipmunk/vect"
  "time"
  "math/rand"
  "net"
	"strconv"
  "fmt"
)

type MainScene struct {
	*engine.SceneData
	Layer1 *engine.GameObject
}

var (
	MainSceneGeneral  *MainScene
	atlas             *engine.ManagedAtlas
	backgroundTexture *engine.Texture
  wallTexture       *engine.Texture
	botTexture        *engine.Texture
	missleTexture     *engine.Texture
	scannerTexture    *engine.Texture
	ArialFont         *engine.Font
  Explosion_ID      engine.ID
	missle            *Missle
	scanner           *Scanner
	Explosion         *engine.GameObject
  Players           map[string]*BotController
)

const (
	MissleTag = "Missle"
  HP_A = 123
  HPGUI_A = 124
)

func LoadTextures() {
	ArialFont, _ = engine.NewFont("./data/arial.ttf", 24)
	ArialFont.Texture.SetReadOnly()

	backgroundTexture, _ = engine.LoadTexture("./data/background.png")
	botTexture, _ = engine.LoadTexture("./data/ship.png")
  wallTexture, _ = engine.LoadTexture("./data/wall.png")
  missleTexture, _ = engine.LoadTexture("./data/missile.png")
  scannerTexture, _ = engine.LoadTexture("./data/scanner.png")

	atlas = engine.NewManagedAtlas(2048, 1024)
	_, Explosion_ID = atlas.LoadGroupSheet("./data/Explosion.png", 128, 128, 6*8)
	atlas.LoadImageID("./data/HealthBar.png", HP_A)
	atlas.LoadImageID("./data/HealthBarGUI.png", HPGUI_A)
	atlas.BuildAtlas()
	atlas.BuildMipmaps()
	atlas.SetFiltering(engine.MipMapLinearNearest, engine.Nearest)
	atlas.Texture.SetReadOnly()
}

func random(min, max int) int {
  rand.Seed(time.Now().Unix())
  return rand.Intn(max - min) + min
}

func SpawnBot(name string, conn net.Conn) *BotController {
	newPlayerController, exists := Players[name]
	if exists {
    newPlayerController.OnDie(false)
	}
  newPlayer := engine.NewGameObject(name)
  newPlayer.AddComponent(engine.NewPhysics(false))
  newPlayer.AddComponent(engine.NewSprite(botTexture))
	newPlayer.AddComponent(NewDestoyable(float32(500)))
	newPlayer.Transform().SetWorldPositionf(float32(random(-500, 500)), float32(random(-300, 300)))
	newPlayer.Transform().SetScalef(50, 50)
	newPlayer.Transform().SetParent2(MainSceneGeneral.Layer1)
	newPlayer.Physics.Shape.SetElasticity(0)
	newPlayer.Physics.Shape.SetFriction(1)
	newPlayer.Physics.Body.SetMass(100)
	newPlayer.Physics.Body.IgnoreGravity = true
	newPlayer.Physics.Body.SetMoment(engine.Inf)

	Health := engine.NewGameObject("HP")
	Health.Transform().SetParent2(MainSceneGeneral.Camera.GameObject())

	HealthGUI := engine.NewGameObject("HPGUI")
	HealthGUI.AddComponent(engine.NewSprite2(atlas.Texture, engine.IndexUV(atlas, HPGUI_A)))
	HealthGUI.Transform().SetParent2(Health)
	HealthGUI.Transform().SetDepth(4)
	HealthGUI.Transform().SetPositionf(0, 0)
	HealthGUI.Transform().SetScalef(50, 50)

	HealthBar := engine.NewGameObject("HealthBar")
	HealthBar.Transform().SetParent2(Health)
	HealthBar.Transform().SetPositionf(-82, 0)
	HealthBar.Transform().SetScalef(100, 50)
	uvHP := engine.IndexUV(atlas, HP_A)

	HealthBarGUI := engine.NewGameObject("HealthBarGUI")
	HealthBarGUI.Transform().SetParent2(HealthBar)
	HealthBarGUI.AddComponent(engine.NewSprite2(atlas.Texture, uvHP))
	HealthBarGUI.Transform().SetScalef(0.52, 1)
	HealthBarGUI.Transform().SetDepth(3)
	HealthBarGUI.Transform().SetPositionf((uvHP.Ratio/2)*HealthBarGUI.Transform().Scale().X, 0)

  newPlayerController = newPlayer.AddComponent(NewBotController(name, conn, Health, HealthBar, missle, scanner)).(*BotController)
  Players[name] = newPlayerController

  return Players[name]
}

func ReorderHealthBars(){
  i := 0
  for k, v := range Players {
    i++
    if v.NameObject == nil {
      v.NameObject = engine.NewGameObject("Name")
      v.NameObject.AddComponent(components.NewUIText(ArialFont, k))
      v.NameObject.Transform().SetParent2(v.Health)
      v.NameObject.Transform().SetDepth(5)
      v.NameObject.Transform().SetPositionf(0, 1)
      v.NameObject.Transform().SetScalef(20, 20)
    }
	  v.Health.Transform().SetPositionf(-float32(engine.Width)/2+110, -float32(engine.Height)/2+(40*float32(i)))
  }
}

func (s *MainScene) Load() {
	engine.SetTitle("Bot Battle!")
  Players = make(map[string]*BotController)
	LoadTextures()

	s.Camera = engine.NewCamera()
	cam := engine.NewGameObject("Camera")
	cam.AddComponent(s.Camera)
	cam.Transform().SetPosition(engine.NewVector3(0, 0, -50))
	cam.Transform().SetScalef(1, 1)

	background := engine.NewGameObject("Background")
	background.AddComponent(engine.NewSprite(backgroundTexture))
	background.Transform().SetScalef(800, 800)
	background.Transform().SetPositionf(0, 0)
	background.Transform().SetDepth(-1)

	hud := engine.NewGameObject("HUD")
	hud.Transform().SetParent2(cam)

	FPSDrawer := engine.NewGameObject("FPS")
	FPSDrawer.Transform().SetParent2(hud)
	txt := FPSDrawer.AddComponent(components.NewUIText(ArialFont, "")).(*components.UIText)
	fps := FPSDrawer.AddComponent(engine.NewFPS()).(*engine.FPS)
	fps.SetAction(func(fps float64) {
		txt.SetString("FPS: " + strconv.FormatFloat(fps, 'f', 2, 32))
	})
	txt.SetAlign(engine.AlignLeft)
	FPSDrawer.Transform().SetPositionf((float32(-engine.Width)/2)+20, (float32(engine.Height)/2)-20)
	FPSDrawer.Transform().SetScalef(20, 20)

	uvs, ind := engine.AnimatedGroupUVs(atlas, Explosion_ID)
	Explosion = engine.NewGameObject("Explosion")
	Explosion.AddComponent(engine.NewSprite3(atlas.Texture, uvs))
	Explosion.Sprite.BindAnimations(ind)
	Explosion.Sprite.AnimationSpeed = 25
	Explosion.Sprite.AnimationEndCallback = func(sprite *engine.Sprite) {
		sprite.GameObject().Destroy()
	}
	Explosion.Transform().SetScalef(30, 30)
	Explosion.Transform().SetDepth(1)

	missleGameObject := engine.NewGameObject("Missle")
	missleGameObject.AddComponent(engine.NewSprite(missleTexture))
	missleGameObject.AddComponent(engine.NewPhysics(false))
	missleGameObject.Transform().SetScalef(20, 20)
	missleGameObject.AddComponent(NewDamageDealer(250))
	missleGameObject.Physics.Shape.IsSensor = true
	missle = NewMissle(30000)
	missleGameObject.AddComponent(missle)
	missle.Explosion = Explosion
	ds := NewDestoyable(0)
	ds.SetDestroyTime(1)
	missleGameObject.AddComponent(ds)

	scannerGameObject := engine.NewGameObject("Scanner")
	scannerGameObject.AddComponent(engine.NewSprite(scannerTexture))
	scannerGameObject.AddComponent(engine.NewPhysics(false))
	scannerGameObject.Transform().SetScalef(20, 20)
	scannerGameObject.Physics.Shape.IsSensor = true
  scanner = NewScanner()
	scannerGameObject.AddComponent(scanner)
	ds = NewDestoyable(0)
	ds.SetDestroyTime(1)
	scannerGameObject.AddComponent(ds)

	//SPACCCEEEEE
	engine.Space.Gravity.Y = 0
	engine.Space.Iterations = 10

	wall := engine.NewGameObject("wall")
	wall.AddComponent(engine.NewSprite(wallTexture))
	wall.AddComponent(NewDestoyable(float32(engine.Inf)))
	wall.Transform().SetScalef(100, 100)
  wall.AddComponent(engine.NewPhysicsShape(true, chipmunk.NewBox(vect.Vect{0, 0}, 100, 100)))
	wall.Physics.Shape.SetElasticity(0)
	wall.Physics.Body.SetMass(10000000)
	wall.Physics.Body.SetMoment(engine.Inf)

  Wall := engine.NewGameObject("Wall")
	for i := -9; i < 9; i++ {
		c := wall.Clone()
		c.Transform().SetParent2(Wall)
	  c.Transform().SetPositionf(float32(i)*80, 400)
		c = wall.Clone()
		c.Transform().SetParent2(Wall)
	  c.Transform().SetPositionf(float32(i)*80, -400)
	}
	for i := -6; i < 6; i++ {
		c := wall.Clone()
		c.Transform().SetParent2(Wall)
	  c.Transform().SetPositionf(680, float32(i)*80)
		c = wall.Clone()
		c.Transform().SetParent2(Wall)
	  c.Transform().SetPositionf(-680, float32(i)*80)
	}

	s.Layer1 = engine.NewGameObject("Layer1")
	s.AddGameObject(cam)
	s.AddGameObject(s.Layer1)
	s.AddGameObject(Wall)
	s.AddGameObject(background)

	MainSceneGeneral = s
	fmt.Println("MainScene loaded")
}

func (m *MainScene) New() engine.Scene {
	gs := new(MainScene)
	gs.SceneData = engine.NewScene("GameScene")
	return gs
}
