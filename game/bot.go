package game

import (
	"botbattle/conn"
	"math/rand"
	"sort"
	"time"
)

const (
	MOVES_SPEED = 500
	SCAN_WAIT   = 500
	GUN_WAIT    = 1000
	CANNON_WAIT = 3000
)

type Bot struct {
	scene    *Scene       `json:"-"`
	client   *conn.Client `json:"-"`
	name     string       `json:"name"`
	rotation int          `json:"rotation"`
	x        int          `json:"x"`
	y        int          `json:"y"`
	health   int          `json:"health"`
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
	}
}

func (self *Bot) Status() (int, int, int, int) {
	return self.x, self.y, self.rotation, self.health
}

func (self *Bot) RotRight() int {
	if self.rotation == 270 {
		self.rotation = 0
	} else {
		self.rotation += 90
	}
	time.Sleep(MOVES_SPEED * time.Millisecond)
	return self.rotation
}

func (self *Bot) RotLeft() int {
	if self.rotation == 0 {
		self.rotation = 270
	} else {
		self.rotation -= 90
	}
	time.Sleep(MOVES_SPEED * time.Millisecond)
	return self.rotation
}

func (self *Bot) MoveForward() (int, int) {
	switch self.rotation {
	case 90:
		if self.y > 0 {
			self.y--
		}
	case 270:
		if self.y < ARENA_HEIGHT-1 {
			self.y++
		}
	case 0:
		if self.x > 0 {
			self.x--
		}
	case 180:
		if self.x < ARENA_WIDTH-1 {
			self.x++
		}
	}
	time.Sleep(MOVES_SPEED * time.Millisecond)
	return self.x, self.y
}

func (self *Bot) MoveBackward() (int, int) {
	switch self.rotation {
	case 90:
		if self.y < ARENA_HEIGHT-1 {
			self.y++
		}
	case 270:
		if self.y > 0 {
			self.y--
		}
	case 0:
		if self.x < ARENA_WIDTH-1 {
			self.x++
		}
	case 180:
		if self.x > 0 {
			self.x--
		}
	}
	time.Sleep(MOVES_SPEED * time.Millisecond)
	return self.x, self.y
}

func (self *Bot) Hit(dmg int) {
	self.health = self.health - dmg
	if self.scene.serv != nil { //check this for testing
		self.scene.serv.Broadcast("bot hit", self.client.Id)
	}
	if self.health <= 0 {
		self.scene.KillBot(self)
	}
}

func (self *Bot) FireGun() bool {
	targets := self.LookingAt()
	if len(targets) > 0 {
		targets[0].Hit(25)
	}
	time.Sleep(GUN_WAIT * time.Millisecond)
	return len(targets) > 0
}

func (self *Bot) FireCannon() bool {
	targets := self.LookingAt()
	if len(targets) > 0 {
		targets[0].Hit(50)
	}
	time.Sleep(CANNON_WAIT * time.Millisecond)
	return len(targets) > 0
}

func (self *Bot) Scan() [][]int {
	var result [][]int
	bots := self.LookingAt()
	for _, bot := range bots {
		x, y, r, h := bot.Status()
		result = append(result, []int{x, y, r, h})
	}
	time.Sleep(SCAN_WAIT * time.Millisecond)
	return result
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
