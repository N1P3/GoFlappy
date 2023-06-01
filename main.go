package main

import (
	"GoFlappy/bird"
	"GoFlappy/pipe"
	rl "github.com/gen2brain/raylib-go/raylib"
	"math/rand"
	"strconv"
)

const (
	screenWidth  int32 = 1920
	screenHeight int32 = 1080
	maxPipes           = 2000
	pipesWidth         = 150
	pipesSpeedX        = 4
)

func main() {
	game := NewGame()
	rl.InitWindow(screenWidth, screenHeight, "GoFlappy")
	game.Load()
	rl.SetTargetFPS(60)

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
	prevPosY          float32
	FramesCounter     int
	TxSprites2        rl.Texture2D
	TxSprites3        rl.Texture2D
}

func NewGame() (g Game) {
	g.init()
	g.generatePipes()
	return
}

func (g *Game) init() {
	g.Floppy.PosX = 400
	g.Floppy.PosY = float32(screenHeight / 2)
	g.Dead = false
	g.Pause = false
	g.Score = 0
	g.Floppy.Radius = 70
	g.WindowShouldClose = false
	g.GameOver = false
	g.FrameRec = rl.NewRectangle(0, 0, float32(g.Floppy.Radius), float32(g.Floppy.Radius))
	g.FramesCounter = 0
	g.prevPosY = 0

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
		x += 600
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
	g.prevPosY = g.Floppy.PosY
	for i := 0; i < maxPipes; i++ {
		g.Pipes[i].X -= float32(pipesSpeedX)
	}
	if !g.Dead {
		if rl.IsKeyPressed(rl.KeySpace) == true {
			g.Floppy.PosY -= 22
		}
		if rl.IsKeyDown(rl.KeySpace) == true {
			g.Floppy.PosY -= 15
		}
		g.Floppy.PosY += 7.5
	}
	if g.Dead {
		g.Floppy.PosX -= pipesSpeedX
	}
	g.isAlive()
	if g.Pipes[g.Score*2].X == g.Floppy.PosX && !g.Dead {
		g.Score = g.Score + 1
	}
}

func (g *Game) Load() {
	g.TxSprites = rl.LoadTexture("images/sprite.png")
	g.TxSprites2 = rl.LoadTexture("images/sprite2.png")
	g.TxSprites3 = rl.LoadTexture("images/sprite3.png")
}
func (g *Game) Unload() {
	rl.UnloadTexture(g.TxSprites)
	rl.UnloadTexture(g.TxSprites2)
	rl.UnloadTexture(g.TxSprites3)
}

func (g *Game) draw() {
	rl.BeginDrawing()

	rl.ClearBackground(rl.SkyBlue)
	for i := 0; i < maxPipes; i++ {
		if g.Pipes[i].X > -100 && g.Pipes[i].X < 2000 {
			rl.DrawRectangle(int32(g.Pipes[i].X), int32(g.Pipes[i].Y), pipesWidth, int32(g.Pipes[i].Height), rl.Green)
		}
	}
	if g.Floppy.PosY > g.prevPosY && g.Floppy.PosX > -100 {
		rl.DrawTextureRec(g.TxSprites3, g.FrameRec, rl.Vector2{X: g.Floppy.PosX, Y: g.Floppy.PosY}, rl.RayWhite)

	} else {
		g.FramesCounter = 0
		rl.DrawTextureRec(g.TxSprites2, g.FrameRec, rl.Vector2{X: g.Floppy.PosX, Y: g.Floppy.PosY}, rl.RayWhite)
	}
	rl.DrawText("SCORE: "+strconv.Itoa(g.Score), 1600, 100, 60, rl.Black)

	if g.GameOver {
		if g.Score > g.HiScore {
			g.HiScore = g.Score
		}
		if !rl.IsKeyPressed(rl.KeySpace) {
			rl.DrawRectangle(410, 300, 1110, 480, rl.White)
			rl.DrawText("Your score was:"+strconv.Itoa(g.Score), 500, 400, 100, rl.Black)
			rl.DrawText("Press 'Spacebar' to restart", 450, 600, 70, rl.Black)
		}
		if rl.IsKeyPressed(rl.KeySpace) {
			g.init()
			g.GameOver = false
			g.Dead = false
			g.generatePipes()
		}
	}
	rl.EndDrawing()
}
