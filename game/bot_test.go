package game

import (
	"testing"
	"encoding/json"
	"botbattle/conn"
)

func NewTestBot(name string, rot, x, y int) *Bot {
  scene := NewScene()
	bot := &Bot{scene, &conn.Client{Id: 1}, name, rot, x, y, 100, 0, false, true}
	scene.bots[bot.client.Id] = bot
  return bot
}

func TestRotateLeft(t *testing.T) {
	bot := NewTestBot("test", 0, 5, 5)
	if status := bot.RotLeft(); status.Rotation != 270 {
		t.Error("Rotation incorrect")
	}
	if status := bot.RotLeft(); status.Rotation != 180 {
		t.Error("Rotation incorrect")
	}
	if status := bot.RotLeft(); status.Rotation != 90 {
		t.Error("Rotation incorrect")
	}
	if status := bot.RotLeft(); status.Rotation != 0 {
		t.Error("Rotation incorrect")
	}
}

func TestRotateRight(t *testing.T) {
	bot := NewTestBot("test", 0, 5, 5)
	if status := bot.RotRight(); status.Rotation != 90 {
		t.Error("Rotation incorrect")
	}
	if status := bot.RotRight(); status.Rotation != 180 {
		t.Error("Rotation incorrect")
	}
	if status := bot.RotRight(); status.Rotation != 270 {
		t.Error("Rotation incorrect")
	}
	if status := bot.RotRight(); status.Rotation != 0 {
		t.Error("Rotation incorrect")
	}
}

func TestTopBoundary(t *testing.T) {
	bot := NewTestBot("test", 90, 5, 1)
	bot.MoveForward()
	status := bot.MoveForward()
	if status.X != 5 || status.Y != 0 {
		t.Error("TestTopBoundary Failed with ", status.X, status.Y)
	}
}

func TestBottomBoundary(t *testing.T) {
	bot := NewTestBot("test", 270, 5, 10)
	bot.MoveForward()
	status := bot.MoveForward()
	if status.X != 5 || status.Y != 11 {
		t.Error("TestBottomBoundary Failed with ", status.X, status.Y)
	}
}

func TestLeftBoundary(t *testing.T) {
	bot := NewTestBot("test", 0, 1, 5)
	bot.MoveForward()
	status := bot.MoveForward()
	if status.X != 0 || status.Y != 5 {
		t.Error("TestLeftBoundary Failed with ", status.X, status.Y)
	}
}


func TestRightBoundary(t *testing.T) {
	bot := NewTestBot("test", 180, 22, 5)
	bot.MoveForward()
	status := bot.MoveForward()
	if status.X != 23 || status.Y != 5 {
		t.Error("TestRightBoundary Failed with ", status.X, status.Y)
	}
}

func TestMoveForwardUp(t *testing.T) {
	bot := NewTestBot("test", 90, 5, 5)
	status := bot.MoveForward()
	if status.X != 5 || status.Y != 4 {
		t.Error("MoveForwardUp Failed with ", status.X, status.Y)
	}
}

func TestMoveForwardDown(t *testing.T) {
	bot := NewTestBot("test", 270, 5, 5)
	status := bot.MoveForward()
	if status.X != 5 || status.Y != 6 {
		t.Error("TestMoveForwardDown Failed with ", status.X, status.Y)
	}
}

func TestMoveForwardLeft(t *testing.T) {
	bot := NewTestBot("test", 0, 5, 5)
	status := bot.MoveForward()
	if status.X != 4 || status.Y != 5 {
		t.Error("TestMoveForwardLeft Failed with ", status.X, status.Y)
	}
}

func TestMoveForwardRigth(t *testing.T) {
	bot := NewTestBot("test", 180, 5, 5)
	status := bot.MoveForward()
	if status.X != 6 || status.Y != 5 {
		t.Error("TestMoveForwardRigth Failed with ", status.X, status.Y)
	}
}

func TestMoveBackwardUp(t *testing.T) {
	bot := NewTestBot("test", 90, 5, 5)
	status := bot.MoveBackward()
	if status.X != 5 || status.Y != 6 {
		t.Error("MoveBackwardUp Failed with ", status.X, status.Y)
	}
}

func TestMoveBackwardDown(t *testing.T) {
	bot := NewTestBot("test", 270, 5, 5)
	status := bot.MoveBackward()
	if status.X != 5 || status.Y != 4 {
		t.Error("TestMoveBackwardDown Failed with ", status.X, status.Y)
	}
}

func TestMoveBackwardLeft(t *testing.T) {
	bot := NewTestBot("test", 0, 5, 5)
	status := bot.MoveBackward()
	if status.X != 6 || status.Y != 5 {
		t.Error("TestMoveBackwardLeft Failed with ", status.X, status.Y)
	}
}

func TestMoveBackwardRigth(t *testing.T) {
	bot := NewTestBot("test", 180, 5, 5)
	status := bot.MoveBackward()
	if status.X != 4 || status.Y != 5 {
		t.Error("TestMoveBackwardRigth Failed with ", status.X, status.Y)
	}
}

//###########
//###########
//######4####
//######5####
//#####123###
//###########
func botSetup() (*Scene, map[int]*Bot) {
	scene := new(Scene)
	scene.bots = map[int]*Bot{
    1: &Bot{scene, &conn.Client{Id: 1}, "tester0", 180, 5, 5, 100, 0, false, true},
		2: &Bot{scene, &conn.Client{Id: 2}, "tester1", 90, 6, 5, 100, 0, false, true},
		3: &Bot{scene, &conn.Client{Id: 3}, "tester2", 0, 7, 5, 100, 0, false, true},
		4: &Bot{scene, &conn.Client{Id: 4}, "tester3", 270, 6, 3, 100, 0, false, true},
		5: &Bot{scene, &conn.Client{Id: 5}, "tester4", 0, 6, 4, 100, 0, false, true},
	}
	return scene, scene.bots
}

func TestLookingAt(t *testing.T) {
	_, bot := botSetup()
	if bots := bot[1].LookingAt(); len(bots) != 2 || bots[0] != bot[2] {
		t.Error("bot1 cant see bot 2")
		return
	}
	if bots := bot[2].LookingAt(); len(bots) != 2 || bots[0] != bot[5] {
		t.Error("bot2 cant see bot 5")
		return
	}
	if bots := bot[3].LookingAt(); len(bots) != 2 || bots[0] != bot[2] {
		t.Error("bot3 cant see bot 2")
		return
	}
	if bots := bot[4].LookingAt(); len(bots) != 2 || bots[0] != bot[5] {
		t.Error("bot4 cant see bot 5")
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
