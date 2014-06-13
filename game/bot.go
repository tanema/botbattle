package game

import (
	"botbattle/conn"
	"encoding/json"
	"math/rand"
	"sort"
	"time"
)

const (
	MOVES_SPEED = 500
	SCAN_WAIT   = 500
	GUN_WAIT    = 1000
	CANNON_WAIT = 3000
  SHIELD_ON   = 3000
  SHIELD_POWER_UP = 5000
  GUN_DMG = 25
  CANNON_DMG = 50
  COLLISION_DMG = 5
)

type Bot struct {
	scene    *Scene       `json:"-"`
	client   *conn.Client `json:"-"`
	name     string       `json:"name"`
	rotation int          `json:"rotation"`
	x        int          `json:"x"`
	y        int          `json:"y"`
	health   int          `json:"health"`
  killcount   int       `json:"killcount"`
  ShieldOn    bool
  ShieldReady bool
}

func NewBot(scene *Scene, new_client *conn.Client, name string) *Bot {
	return &Bot{
		scene,
		new_client,
		name,
		[]int{0, 90, 180, 270}[rand.Intn(4)],
		rand.Intn(ARENA_WIDTH),
		rand.Intn(ARENA_HEIGHT),
		100,
    0,
    false,
    true,
	}
}

func (self *Bot) Status() *Status {
	return &Status{self.client.Id, self.x, self.y, self.rotation, self.name, self.health, self.killcount}
}

func (self *Bot) RotRight() *Status {
	if self.rotation == 270 {
		self.rotation = 0
	} else {
		self.rotation += 90
	}
	time.Sleep(MOVES_SPEED * time.Millisecond)
	return self.Status()
}

func (self *Bot) RotLeft() *Status {
	if self.rotation == 0 {
		self.rotation = 270
	} else {
		self.rotation -= 90
	}
	time.Sleep(MOVES_SPEED * time.Millisecond)
	return self.Status()
}

func (self *Bot) MoveForward() *Status {
	switch self.rotation {
	case 90:
		if self.y > 0 && self.At(self.x, self.y-1) == nil {
			self.y--
		} else if self.At(self.x, self.y-1) != nil {
      self.Hit(COLLISION_DMG)
    }
	case 270:
		if self.y < ARENA_HEIGHT && self.At(self.x, self.y+1) == nil {
			self.y++
		} else if self.At(self.x, self.y+1) != nil {
      self.Hit(COLLISION_DMG)
    }
	case 0:
		if self.x > 0 && self.At(self.x-1, self.y) == nil {
			self.x--
		} else if self.At(self.x-1, self.y) != nil {
      self.Hit(COLLISION_DMG)
    }
	case 180:
		if self.x < ARENA_WIDTH && self.At(self.x+1, self.y) == nil {
			self.x++
		} else if self.At(self.x+1, self.y) != nil {
      self.Hit(COLLISION_DMG)
    }
	}
	time.Sleep(MOVES_SPEED * time.Millisecond)
	return self.Status()
}

func (self *Bot) MoveBackward() *Status {
	switch self.rotation {
	case 90:
		if self.y < ARENA_HEIGHT && self.At(self.x, self.y+1) == nil {
			self.y++
		} else if self.At(self.x, self.y+1) != nil {
      self.Hit(COLLISION_DMG)
    }
	case 270:
		if self.y > 0 && self.At(self.x, self.y-1) == nil {
			self.y--
		}  else if self.At(self.x, self.y-1) != nil {
      self.Hit(COLLISION_DMG)
    }
	case 0:
		if self.x < ARENA_WIDTH && self.At(self.x+1, self.y) == nil {
			self.x++
		} else if self.At(self.x+1, self.y) != nil {
      self.Hit(COLLISION_DMG)
    }
	case 180:
		if self.x > 0  && self.At(self.x-1, self.y) == nil {
			self.x--
		} else if self.At(self.x-1, self.y) == nil {
      self.Hit(COLLISION_DMG)
    }
	}
	time.Sleep(MOVES_SPEED * time.Millisecond)
	return self.Status()
}

func (self *Bot) Hit(dmg int) bool {
  if self.ShieldOn {
    self.health -= dmg/2
  } else {
    self.health -= dmg
  }
	if self.scene.serv != nil { //check this for testing
		self.scene.serv.Broadcast("bot hit", self.client.Id, dmg)
	}
	if self.health <= 0 {
		self.scene.KillBot(self)
    return true
	}
  return false
}

func (self *Bot) FireGun() bool {
	targets := self.LookingAt()
	if len(targets) > 0 {
		if targets[0].Hit(GUN_DMG) {
      self.killcount++
		  self.scene.serv.Broadcast("bot notch", self.client.Id, self.killcount)
    }
	}
	time.Sleep(GUN_WAIT * time.Millisecond)
	return len(targets) > 0
}

func (self *Bot) FireCannon() bool {
	targets := self.LookingAt()
	if len(targets) > 0 {
		if targets[0].Hit(CANNON_DMG) {
      self.killcount++
		  self.scene.serv.Broadcast("bot notch", self.client.Id, self.killcount)
    }
	}
	time.Sleep(CANNON_WAIT * time.Millisecond)
	return len(targets) > 0
}

func (self *Bot) Scan() []string {
	var result []string
	bots := self.LookingAt()
	for _, bot := range bots {
		json_status, _ := json.Marshal(bot.Status())
		result = append(result, string(json_status))
	}
	time.Sleep(SCAN_WAIT * time.Millisecond)
	return result
}

func (self *Bot) Shield() bool {
  if(self.ShieldReady){
    self.ShieldOn = true;
    self.ShieldReady = false;
    self.scene.serv.Broadcast("shield", self.client.Id, true)
    go func(){
      time.Sleep(SHIELD_ON * time.Millisecond)
      self.ShieldOn = false;
      self.scene.serv.Broadcast("shield", self.client.Id, false)
      time.Sleep(SHIELD_POWER_UP * time.Millisecond)
      self.ShieldReady = true;
    }()
    return true
  } else {
    return false
  }
}

func (self *Bot) At(x, y int) *Bot {
	for _, bot := range self.scene.bots {
    if bot.x == x && bot.y == y {
      return bot
    }
	}
  return nil
}

func (self *Bot) LookingAt() []*Bot {
	result := []*Bot{}
	for _, bot := range self.scene.bots {
		switch self.rotation {
		case 90:
			if self.y > bot.y && self.x == bot.x {
				result = append(result, bot)
			}
		case 270:
			if self.y < bot.y && self.x == bot.x {
				result = append(result, bot)
			}
		case 0:
			if self.y == bot.y && self.x > bot.x {
				result = append(result, bot)
			}
		case 180:
			if self.y == bot.y && self.x < bot.x {
				result = append(result, bot)
			}
		}
	}

	var sorter BotsBy
	switch self.rotation {
	case 90:
		sorter = func(b1, b2 *Bot) bool {
			return b1.y > b2.y
		}
	case 270:
		sorter = func(b1, b2 *Bot) bool {
			return b1.y < b2.y
		}
	case 0:
		sorter = func(b1, b2 *Bot) bool {
			return b1.x > b2.x
		}
	case 180:
		sorter = func(b1, b2 *Bot) bool {
			return b1.x < b2.x
		}
	}
	BotsBy(sorter).Sort(result)
	return result
}

type BotsBy func(b1, b2 *Bot) bool

func (by BotsBy) Sort(bots []*Bot) {
	ps := &botSorter{
		bots: bots,
		by:   by, // The Sort method's receiver is the function (closure) that defines the sort order.
	}
	sort.Sort(ps)
}

type botSorter struct {
	bots []*Bot
	by   func(b1, b2 *Bot) bool // Closure used in the Less method.
}

func (s *botSorter) Len() int           { return len(s.bots) }
func (s *botSorter) Swap(i, j int)      { s.bots[i], s.bots[j] = s.bots[j], s.bots[i] }
func (s *botSorter) Less(i, j int) bool { return s.by(s.bots[i], s.bots[j]) }
