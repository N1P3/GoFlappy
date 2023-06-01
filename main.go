package main

import (
	"GoFlappy/bird"
	"GoFlappy/pipe"
	rl "github.com/gen2brain/raylib-go/raylib"
	"math/rand"
)

const (
	screenWidth  int32 = 1920
	screenHeight int32 = 1080
	maxPipes           = 2000
	pipesWidth         = 70
	pipesSpeedX        = 8
)

func main() {
	game := NewGame()
	rl.InitWindow(screenWidth, screenHeight, "GoFlappy")
	game.Load()
	rl.SetTargetFPS(30)

	for !rl.WindowShouldClose() {
		game.Update()
		game.draw()
	}
	game.Unload()
	rl.CloseWindow()
}

type Game struct {
	GameOver          bool
	Dead              bool
	Pause             bool
	Score             int
	HiScore           int
	WindowShouldClose bool
	Floppy            bird.Bird
	Pipes             [maxPipes]pipe.Pipe
	TxSprites         rl.Texture2D
	FrameRec          rl.Rectangle
}

func NewGame() (g Game) {
	g.init()
	g.generatePipes()
	return
}

func (g *Game) init() {
	g.Floppy.PosX = float32(screenWidth / 5)
	g.Floppy.PosY = float32(screenHeight / 2)
	g.Dead = false
	g.Pause = false
	g.Score = 0
	g.Floppy.Radius = 50
	g.WindowShouldClose = false
	g.GameOver = false
	g.FrameRec = rl.NewRectangle(0, 0, float32(g.Floppy.Radius), float32(g.Floppy.Radius))

}

func (g *Game) generatePipes() {
	var x float32
	x = 1600
	for i := 0; i < maxPipes; i += 2 {
		g.Pipes[i].X = x
		g.Pipes[i].Width = pipesWidth
		g.Pipes[i+1].X = x
		g.Pipes[i].Y = 0
		g.Pipes[i].Height = float32(rand.Intn(800))
		g.Pipes[i+1].Y = g.Pipes[i].Height + 300
		g.Pipes[i+1].Height = float32(screenHeight - int32(g.Pipes[i+1].Y))
		x += 500
		g.Pipes[i].Width = pipesWidth
		g.Pipes[i+1].Width = pipesWidth
	}
}

func (g *Game) isAlive() {
	for i := 0; i < maxPipes; i++ {
		if rl.CheckCollisionRecs(rl.NewRectangle(g.Floppy.PosX, g.Floppy.PosY, float32(g.Floppy.Radius), float32(g.Floppy.Radius)), g.Pipes[i].Rectangle) {
			g.Dead = true
			g.GameOver = true
		}
	}
	if g.Floppy.PosY <= 0 || g.Floppy.PosY >= float32(screenHeight) {
		g.Dead = true
		g.GameOver = true
	}
}

func (g *Game) Update() {
	for i := 0; i < maxPipes; i++ {
		g.Pipes[i].X -= float32(pipesSpeedX)
	}
	if !g.Dead {
		if rl.IsKeyPressed(rl.KeySpace) == true {
			g.Floppy.PosY -= 45
		}
		if rl.IsKeyDown(rl.KeySpace) == true {
			g.Floppy.PosY -= 30
		}
		g.Floppy.PosY += 15
	}
	if g.Dead {
		g.Floppy.PosX -= pipesSpeedX
	}
	g.isAlive()
}

func (g *Game) Load() {
	g.TxSprites = rl.LoadTexture("images/sprite.png")
}
func (g *Game) Unload() {
	rl.UnloadTexture(g.TxSprites)
}

func (g *Game) draw() {
	rl.BeginDrawing()

	rl.ClearBackground(rl.SkyBlue)

	for i := 0; i < maxPipes; i++ {
		if g.Pipes[i].X > -100 && g.Pipes[i].X < 2000 {
			rl.DrawRectangle(int32(g.Pipes[i].X), int32(g.Pipes[i].Y), pipesWidth, int32(g.Pipes[i].Height), rl.Green)
		}
	}
	if g.Floppy.PosX > -100 {
		rl.DrawTextureRec(g.TxSprites, g.FrameRec, rl.Vector2{X: g.Floppy.PosX, Y: g.Floppy.PosY}, rl.RayWhite)
	}
	if g.GameOver {
		rl.DrawText("loser", 300, 300, 100, rl.White)
	}
	rl.EndDrawing()
}
