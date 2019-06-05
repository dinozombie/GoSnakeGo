package main

import (
	"github.com/rthornton128/goncurses"
	"log"
	"math/rand"
	"os"
	"time"
)

const speed = time.Second / 16

//Snake struct
type Snake struct {
	nodes []*SnakeNode
}

//SnakeNode struct
type SnakeNode struct {
	Direction byte //0 - up, 1 - right, 2 - down, 3 - left
	X         int
	Y         int
}

//Food struct
type Food struct {
	X int
	Y int
}

func newSnake(x, y int) *Snake {
	head := newSnakeNode(1, x, y)
	head2 := newSnakeNode(1, x-1, y)
	head3 := newSnakeNode(1, x-2, y)

	nodes := []*SnakeNode{head, head2, head3}
	return &Snake{nodes}
}

func newSnakeNode(direction byte, x, y int) *SnakeNode {
	return &SnakeNode{direction, x, y}
}

func newFood(s *Snake, maxX, maxY int) *Food {
	rand.Seed(time.Now().UnixNano())
	x, y := rand.Intn(maxX), rand.Intn(maxY)
	for _, n := range s.nodes {
		if n.X == x && n.Y == y {
			return newFood(s, maxX, maxY)
		}
	}
	return &Food{x, y}
}

func (f *Food) update(s *Snake, maxX, maxY int) {
	rand.Seed(time.Now().UnixNano())
	x, y := rand.Intn(maxX), rand.Intn(maxY)
	for _, n := range s.nodes {
		if n.X == x && n.Y == y {
			f.update(s, maxX, maxY)
			return
		}
	}
	f.X, f.Y = x, y
}

func (s *Snake) update(newDirection int) {
	for i := len(s.nodes) - 1; i >= 0; i-- {
		switch s.nodes[i].Direction {
		case 0:
			s.nodes[i].Y--
		case 1:
			s.nodes[i].X++
		case 2:
			s.nodes[i].Y++
		case 3:
			s.nodes[i].X--
		}

		if i > 0 {
			s.nodes[i].Direction = s.nodes[i-1].Direction
		}
	}
}

func (s *Snake) grow() {
	var x, y int

	dir := s.nodes[len(s.nodes)-1].Direction

	switch dir {
	case 0:
		x = s.nodes[len(s.nodes)-1].X
		y = s.nodes[len(s.nodes)-1].Y + 1
	case 1:
		x = s.nodes[len(s.nodes)-1].X - 1
		y = s.nodes[len(s.nodes)-1].Y
	case 2:
		x = s.nodes[len(s.nodes)-1].X
		y = s.nodes[len(s.nodes)-1].Y - 1
	case 3:
		x = s.nodes[len(s.nodes)-1].X + 1
		y = s.nodes[len(s.nodes)-1].Y
	}
	s.nodes = append(s.nodes, newSnakeNode(dir, x, y))
}

func checkCollisions(s *Snake, f *Food, maxX, maxY int) bool {
	if (s.nodes[0].Direction == 0 && s.nodes[0].Y == 0) ||
		(s.nodes[0].Direction == 1 && s.nodes[0].X == maxX-1) ||
		(s.nodes[0].Direction == 2 && s.nodes[0].Y == maxY-1) ||
		(s.nodes[0].Direction == 3 && s.nodes[0].X == 0) {
		return false
	}

	for i := 3; i < len(s.nodes); i++ {
		switch s.nodes[0].Direction {
		case 0:
			if s.nodes[0].X == s.nodes[i].X && s.nodes[0].Y == s.nodes[i].Y+1 {
				return false
			}
		case 1:
			if s.nodes[0].Y == s.nodes[i].Y && s.nodes[0].X == s.nodes[i].X-1 {
				return false
			}
		case 2:
			if s.nodes[0].X == s.nodes[i].X && s.nodes[0].Y == s.nodes[i].Y-1 {
				return false
			}
		case 3:
			if s.nodes[0].Y == s.nodes[i].Y && s.nodes[0].X == s.nodes[i].X+1 {
				return false
			}
		}
	}

	switch s.nodes[0].Direction {
	case 0:
		if s.nodes[0].X == f.X && s.nodes[0].Y == f.Y+1 {
			s.grow()
			f.update(s, maxX, maxY)
		}
	case 1:
		if s.nodes[0].X == f.X-1 && s.nodes[0].Y == f.Y {
			s.grow()
			f.update(s, maxX, maxY)
		}

	case 2:
		if s.nodes[0].X == f.X && s.nodes[0].Y == f.Y-1 {
			s.grow()
			f.update(s, maxX, maxY)
		}

	case 3:
		if s.nodes[0].X == f.X+1 && s.nodes[0].Y == f.Y {
			s.grow()
			f.update(s, maxX, maxY)
		}
	}
	return true
}

func handleInput(stdscr *goncurses.Window, snake *Snake) bool {
	k := stdscr.GetChar()

	switch goncurses.KeyString(k) {
	case "q":
		return false
	case "left":
		if snake.nodes[1].Direction != 1 {
			snake.nodes[0].Direction = 3
		}
	case "right":
		if snake.nodes[1].Direction != 3 {
			snake.nodes[0].Direction = 1
		}
	case "down":
		if snake.nodes[1].Direction != 0 {
			snake.nodes[0].Direction = 2
		}
	case "up":
		if snake.nodes[1].Direction != 2 {
			snake.nodes[0].Direction = 0
		}
	case "x":
		snake.grow()
	}
	return true
}

func updateScreen(screen *goncurses.Window, snake *Snake, f *Food) {
	for _, n := range snake.nodes {
		screen.MovePrint(n.Y, n.X, `o`)
	}
	screen.MovePrint(f.Y, f.X, `x`)
}

func main() {
	f, err := os.Create("err.log")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	log.SetOutput(f)

	stdscr, err := goncurses.Init()
	if err != nil {
		log.Fatal("init", err)
	}
	defer goncurses.End()

	goncurses.Cursor(0)
	goncurses.Echo(false)
	goncurses.HalfDelay(1)

	stdscr.Keypad(true)

	if err := goncurses.StartColor(); err != nil {
		log.Fatal(err)
	}

	goncurses.InitPair(1, goncurses.C_RED, goncurses.C_BLACK)
	stdscr.ColorOn(1)

	maxY, maxX := stdscr.MaxYX()
	snake := newSnake(maxX/2, maxY/2)
	food := newFood(snake, maxX, maxY)

	ticker := time.NewTicker(speed)
	updateScreen(stdscr, snake, food)

mainloop:
	for {
		stdscr.Refresh()
		select {
		case <-ticker.C:
			if !checkCollisions(snake, food, maxX, maxY) {
				break mainloop
			}
			snake.update(-1)
			stdscr.Erase()
			updateScreen(stdscr, snake, food)

		default:
			if !handleInput(stdscr, snake) {
				break mainloop
			}
		}
	}
}
