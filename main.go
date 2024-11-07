package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"

	"golang.org/x/term"
)

type Snake struct {
	Position [][2]float32
	Velocity [2]float32
	Growing  bool
}

func (s *Snake) UpdatePosition(dt float32) {
	// new head
	H := [][2]float32{
		{
			s.Position[0][0] + dt*s.Velocity[0],
			s.Position[0][1] + dt*s.Velocity[1],
		},
	}

	// if the int of the current head and the int of the uodated arent the same then enqueue and push the new head on
	if int(s.Position[0][0]) != int(H[0][0]) || int(s.Position[0][1]) != int(H[0][1]) {
		// [ A, B , C]
		// remove the C segment
		if !s.Growing {
			s.Position = s.Position[:len(s.Position)-1]
		} else {
			s.Growing = false
		}

		// Add New Head to the top to get [H, A, B]
		s.Position = append(H, s.Position...)
	} else {
		s.Position[0] = H[0]
	}
	// else change th value of the head

	// Check if the snakes head ended up at Left Or Rigth wall
	switch int(s.Position[0][0]) {
	case 0:
		// Right Case
		s.Position[0][0] = 21
	case 22:
		// Left case
		s.Position[0][0] = 1
	}

	// Check if the snakes head ended up at Top or Bottom wall
	switch int(s.Position[0][1]) {
	case 0:
		// Bottom Case
		s.Position[0][1] = 10
	case 11:
		// Top Case
		s.Position[0][1] = 1
	}
}

// GAME STRUCT DATA
type GameWorld struct {
	InputChannel chan rune
	Snake        *Snake
	World        []byte
	Apple        [2]int32
	Score        int32
	GameOver     bool
}

func (g *GameWorld) Render() {
	// --------------------- DRAW --------------
	// Position x + y*WORLD_WIDTH
	ap := g.Apple[0] + g.Apple[1]*25

	// place apple on map
	g.World[ap] = byte('@')

	for _, segment := range g.Snake.Position {
		sp := int(segment[0]) + int(segment[1])*25
		g.World[sp] = byte('s')
	}

	// --------------------- PAINT --------------
	// place cursor in top left of terminal
	fmt.Print("\x1b[H")

	// Paint the world interpreting the bytes as a string
	fmt.Printf("%s", g.World)

	// render score
	fmt.Printf("\r\nScore: %v", g.Score)

	// --------------------- RESET --------------
	// redraw the empty char bytes to the map
	for _, segment := range g.Snake.Position {
		sp := int(segment[0]) + int(segment[1])*25
		g.World[sp] = byte(' ')
	}
}

func (g *GameWorld) QuitGame(s string) {
	fmt.Printf("\033[H\033[J\r\n%s\r\n", s)
	os.Exit(0)
}

func (g *GameWorld) CheckCollisions() {
	SNAKE_POS := g.Snake.Position[0]

	for _, seg := range g.Snake.Position[1:] {
		if int(SNAKE_POS[0]) == int(seg[0]) && int(SNAKE_POS[1]) == int(seg[1]) {
			g.QuitGame("Oh No You Died...Fool :)")
		}
	}

	// Check for contact on x and y coordinates
	APPLE_HIT_X := int(SNAKE_POS[0]) == int(g.Apple[0])
	APPLE_HIT_Y := int(SNAKE_POS[1]) == int(g.Apple[1])

	if APPLE_HIT_X && APPLE_HIT_Y {
		g.Score += 1

		// choose a rabdom spot for the apple
		a_x := int32(rand.Intn(20) + 1)
		a_y := int32(rand.Intn(8) + 1)
		// if the board spot is empty then place it
		// Do this until you find an empty spot
		for g.World[a_x+a_y*25] != byte(' ') {
			a_x = int32(rand.Intn(20) + 1)
			a_y = int32(rand.Intn(8) + 1)
		}
		g.Apple = [2]int32{a_x, a_y}

		g.Snake.Growing = true
	}
}

func (g *GameWorld) ProcessInput() {
	select {
	case val := <-g.InputChannel:
		switch val {
		case 'w':
			if g.Snake.Velocity != [2]float32{0, 4} {
				g.Snake.Velocity = [2]float32{0, -4} // Move Up
			}
		case 'a':
			if g.Snake.Velocity != [2]float32{4, 0} {
				g.Snake.Velocity = [2]float32{-4, 0} // Move Left
			}
		case 's':
			if g.Snake.Velocity != [2]float32{0, -4} {
				g.Snake.Velocity = [2]float32{0, 4} // Move Down
			}
		case 'd':
			if g.Snake.Velocity != [2]float32{-4, 0} {
				g.Snake.Velocity = [2]float32{4, 0} // Move Right
			}
		case 'q':
			g.QuitGame("Thanks For Playing!")
		}
	default:
		// No input continue without blocking
	}
}

func InputReader(ch chan<- rune) {
	reader := bufio.NewReader(os.Stdin)

	// wait for a character to be sent through stdin
	// the send it in the channell
	for {
		char, _, _ := reader.ReadRune()
		ch <- char
	}
}

func InitialiseGame() *GameWorld {
	world := []byte{}

	// Build initial world
	world = append(world, []byte("|=====================|\r\n")...)
	for range 10 {
		world = append(world, []byte("|                     |\r\n")...)
	}
	world = append(world, []byte("|=====================|\r\n")...)

	// Initailise the terminal
	fmt.Print("\033[H\033[J")
	fmt.Print("\x1b[?25l")

	return &GameWorld{
		GameOver: false,
		Score:    0,
		Apple:    [2]int32{2, 2},
		World:    world,
		Snake: &Snake{
			Position: [][2]float32{{5, 4}, {4, 4}, {3, 4}, {2, 4}, {1, 4}},
			Velocity: [2]float32{4, 0},
			Growing:  false,
		},
		InputChannel: make(chan rune),
	}
}

func (g *GameWorld) StartGame() {
	oldstate, _ := term.MakeRaw(int(os.Stdin.Fd()))
	defer term.Restore(int(os.Stderr.Fd()), oldstate)

	// Start Input Reader
	go InputReader(g.InputChannel)

	// Get Initail times
	t1 := time.Now()
	for {
		t2 := time.Now()                    // t2 is updated to the current time
		dt := float32(t2.Sub(t1).Seconds()) // get the delta t
		t1 = t2                             // t1 should now be set to t2

		g.ProcessInput()           // Process Input
		g.Snake.UpdatePosition(dt) // Update Snale Position
		g.CheckCollisions()        // Check if Apple is Eaten
		g.Render()                 // Render Map
	}
}

func main() {
	game := InitialiseGame()
	game.StartGame()
}
