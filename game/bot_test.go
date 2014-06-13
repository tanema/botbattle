package game

import (
	"testing"
	"encoding/json"
)

var scene = NewScene()

func TestRotateLeft(t *testing.T) {
	bot := &Bot{scene, nil, "tester", 0, 5, 5, 100, false, true}
	if rot := bot.RotLeft(); rot != 270 {
		t.Error("Rotation incorrect")
	}
	if rot := bot.RotLeft(); rot != 180 {
		t.Error("Rotation incorrect")
	}
	if rot := bot.RotLeft(); rot != 90 {
		t.Error("Rotation incorrect")
	}
	if rot := bot.RotLeft(); rot != 0 {
		t.Error("Rotation incorrect")
	}
}

func TestRotateRight(t *testing.T) {
	bot := &Bot{scene, nil, "tester", 0, 5, 5, 100, false, true}
	if rot := bot.RotRight(); rot != 90 {
		t.Error("Rotation incorrect")
	}
	if rot := bot.RotRight(); rot != 180 {
		t.Error("Rotation incorrect")
	}
	if rot := bot.RotRight(); rot != 270 {
		t.Error("Rotation incorrect")
	}
	if rot := bot.RotRight(); rot != 0 {
		t.Error("Rotation incorrect")
	}
}

func TestMoveForwardUp(t *testing.T) {
	bot := &Bot{scene, nil, "tester", 90, 5, 5, 100, false, true}
	x, y := bot.MoveForward()
	if x != 5 || y != 4 {
		t.Error("MoveForwardUp Failed with ", x, y)
	}
}

func TestMoveForwardDown(t *testing.T) {
	bot := &Bot{scene, nil, "tester", 270, 5, 5, 100, false, true}
	x, y := bot.MoveForward()
	if x != 5 || y != 6 {
		t.Error("TestMoveForwardDown Failed with ", x, y)
	}
}

func TestMoveForwardLeft(t *testing.T) {
	bot := &Bot{scene, nil, "tester", 0, 5, 5, 100, false, true}
	x, y := bot.MoveForward()
	if x != 4 || y != 5 {
		t.Error("TestMoveForwardLeft Failed with ", x, y)
	}
}

func TestMoveForwardRigth(t *testing.T) {
	bot := &Bot{scene, nil, "tester", 180, 5, 5, 100, false, true}
	x, y := bot.MoveForward()
	if x != 6 || y != 5 {
		t.Error("TestMoveForwardRigth Failed with ", x, y)
	}
}

func TestMoveBackwardUp(t *testing.T) {
	bot := &Bot{scene, nil, "tester", 90, 5, 5, 100, false, true}
	x, y := bot.MoveBackward()
	if x != 5 || y != 6 {
		t.Error("MoveBackwardUp Failed with ", x, y)
	}
}

func TestMoveBackwardDown(t *testing.T) {
	bot := &Bot{scene, nil, "tester", 270, 5, 5, 100, false, true}
	x, y := bot.MoveBackward()
	if x != 5 || y != 4 {
		t.Error("TestMoveBackwardDown Failed with ", x, y)
	}
}

func TestMoveBackwardLeft(t *testing.T) {
	bot := &Bot{scene, nil, "tester", 0, 5, 5, 100, false, true}
	x, y := bot.MoveBackward()
	if x != 6 || y != 5 {
		t.Error("TestMoveBackwardLeft Failed with ", x, y)
	}
}

func TestMoveBackwardRigth(t *testing.T) {
	bot := &Bot{scene, nil, "tester", 180, 5, 5, 100, false, true}
	x, y := bot.MoveBackward()
	if x != 4 || y != 5 {
		t.Error("TestMoveBackwardRigth Failed with ", x, y)
	}
}

func botSetup() (*Scene, map[int]*Bot) {
	scene := new(Scene)
	scene.bots = map[int]*Bot{
		1: &Bot{scene, nil, "tester0", 180, 5, 5, 100, false, true},
		2: &Bot{scene, nil, "tester1", 90, 6, 5, 100, false, true},
		3: &Bot{scene, nil, "tester2", 0, 7, 5, 100, false, true},
		4: &Bot{scene, nil, "tester3", 270, 6, 3, 100, false, true},
		5: &Bot{scene, nil, "tester4", 0, 6, 4, 100, false, true},
	}
	return scene, scene.bots
}

func TestLookingAt(t *testing.T) {
	_, bot := botSetup()
	if bots := bot[1].LookingAt(); len(bots) != 2 || bots[0] != bot[2] {
		t.Error("bot1 cants see bot 2")
		return
	}
	if bots := bot[2].LookingAt(); len(bots) != 2 || bots[0] != bot[5] {
		t.Error("bot2 cants see bot 5")
		return
	}
	if bots := bot[3].LookingAt(); len(bots) != 2 || bots[0] != bot[2] {
		t.Error("bot3 cants see bot 2")
		return
	}
	if bots := bot[4].LookingAt(); len(bots) != 2 || bots[0] != bot[5] {
		t.Error("bot4 cants see bot 5")
		return
	}
}

func TestScan(t *testing.T) {
	_, bot := botSetup()
	result := bot[1].Scan()
	if len(result) != 2 && len(result[0]) != 4 {
		t.Error("not enough results")
		return
	}

  my_status := &Status{}
	json.Unmarshal([]byte(result[0]), my_status)

	if my_status.X != bot[2].x || my_status.Y != bot[2].y {
		t.Error("got incorrect ordering")
		return
	}
}

func TestFireGun(t *testing.T) {
	_, bot := botSetup()
	result := bot[1].FireGun()
	if result != true {
		t.Error("no hit")
		return
	}
	if bot[2].health != 75 {
		t.Error("Did not damage bot2")
		return
	}

	result = bot[5].FireGun()
	if result != false {
		t.Error("false hit")
		return
	}
}

func TestFireCannon(t *testing.T) {
	_, bot := botSetup()
	result := bot[1].FireCannon()
	if result != true {
		t.Error("no hit")
		return
	}
	if bot[2].health != 50 {
		t.Error("Did not damage bot2")
		return
	}

	result = bot[5].FireCannon()
	if result != false {
		t.Error("false hit")
		return
	}
}
