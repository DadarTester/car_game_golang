package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth  = 800
	screenHeight = 600
	carWidth     = 50
	carHeight    = 80
	roadWidth    = 400
	stripeSpeed  = 7.0 
)

// Загружаем спрайты
var (
	playerCarSprite  *ebiten.Image
	blueEnemySprite  *ebiten.Image
	greenEnemySprite *ebiten.Image
	coinSprite       *ebiten.Image
	roadSprite       *ebiten.Image
	stripeSprite     *ebiten.Image
)

type Car struct {
	X, Y    float64
	Speed   float64
	Width   float64
	Height  float64
}

type Coin struct {
	X, Y    float64
	Width   float64
	Height  float64
	Speed   float64
	Active  bool
}

type EnemyCar struct {
	X, Y    float64
	Width   float64
	Height  float64
	Speed   float64
	Color   string
}

type Game struct {
	playerCar     Car
	enemyCars     []EnemyCar
	coins         []Coin
	stripeY       float64 
	gameOver      bool
	score         int
	spawnTimer    int
	coinsCollected int
}

func createPlayerCarSprite() *ebiten.Image {
	img := ebiten.NewImage(carWidth, carHeight)
	
	img.Fill(color.RGBA{220, 20, 20, 255})
	
	windowColor := color.RGBA{50, 100, 200, 255}
	ebitenutil.DrawRect(img, 10, 10, 30, 20, windowColor)
	ebitenutil.DrawRect(img, 10, 40, 30, 20, windowColor)
	
	wheelColor := color.RGBA{20, 20, 20, 255}
	ebitenutil.DrawRect(img, 5, 15, 5, 15, wheelColor)
	ebitenutil.DrawRect(img, 40, 15, 5, 15, wheelColor)
	ebitenutil.DrawRect(img, 5, 50, 5, 15, wheelColor)
	ebitenutil.DrawRect(img, 40, 50, 5, 15, wheelColor)
	
	lightColor := color.RGBA{255, 255, 0, 255}
	ebitenutil.DrawRect(img, 5, 5, 5, 5, lightColor)
	ebitenutil.DrawRect(img, 40, 5, 5, 5, lightColor)
	
	return img
}

func createBlueEnemySprite() *ebiten.Image {
	img := ebiten.NewImage(carWidth, carHeight)
	
	img.Fill(color.RGBA{20, 20, 220, 255})
	
	windowColor := color.RGBA{100, 150, 255, 255}
	ebitenutil.DrawRect(img, 10, 10, 30, 20, windowColor)
	ebitenutil.DrawRect(img, 10, 40, 30, 20, windowColor)
	
	wheelColor := color.RGBA{20, 20, 20, 255}
	ebitenutil.DrawRect(img, 5, 15, 5, 15, wheelColor)
	ebitenutil.DrawRect(img, 40, 15, 5, 15, wheelColor)
	ebitenutil.DrawRect(img, 5, 50, 5, 15, wheelColor)
	ebitenutil.DrawRect(img, 40, 50, 5, 15, wheelColor)
	
	lightColor := color.RGBA{255, 0, 0, 255}
	ebitenutil.DrawRect(img, 5, 70, 5, 5, lightColor)
	ebitenutil.DrawRect(img, 40, 70, 5, 5, lightColor)
	
	return img
}

func createGreenEnemySprite() *ebiten.Image {
	img := ebiten.NewImage(carWidth, carHeight)
	
	// Корпус машинки (зеленый)
	img.Fill(color.RGBA{20, 220, 20, 255})
	
	windowColor := color.RGBA{100, 255, 150, 255}
	ebitenutil.DrawRect(img, 10, 10, 30, 20, windowColor)
	ebitenutil.DrawRect(img, 10, 40, 30, 20, windowColor)
	
	wheelColor := color.RGBA{20, 20, 20, 255}
	ebitenutil.DrawRect(img, 5, 15, 5, 15, wheelColor)
	ebitenutil.DrawRect(img, 40, 15, 5, 15, wheelColor)
	ebitenutil.DrawRect(img, 5, 50, 5, 15, wheelColor)
	ebitenutil.DrawRect(img, 40, 50, 5, 15, wheelColor)

	lightColor := color.RGBA{255, 0, 0, 255}
	ebitenutil.DrawRect(img, 5, 70, 5, 5, lightColor)
	ebitenutil.DrawRect(img, 40, 70, 5, 5, lightColor)
	
	return img
}

func createCoinSprite() *ebiten.Image {
	size := 30
	img := ebiten.NewImage(size, size)
	coinColor := color.RGBA{255, 215, 0, 255}
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			dx := float64(x - size/2)
			dy := float64(y - size/2)
			distance := dx*dx + dy*dy
			if distance <= float64(size/2*size/2) {
				img.Set(x, y, coinColor)
			}
		}
	}
	
	darkCoinColor := color.RGBA{205, 173, 0, 255}
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			dx := float64(x - size/2)
			dy := float64(y - size/2)
			distance := dx*dx + dy*dy
			
			if distance <= float64(size/4*size/4) {
				img.Set(x, y, darkCoinColor)
			}
		}
	}
	
	return img
}

func init() {
	playerCarSprite = createPlayerCarSprite()
	blueEnemySprite = createBlueEnemySprite()
	greenEnemySprite = createGreenEnemySprite()
	coinSprite = createCoinSprite()
	roadSprite = ebiten.NewImage(roadWidth, screenHeight)
	roadSprite.Fill(color.RGBA{50, 50, 50, 255})
	
	stripeSprite = ebiten.NewImage(10, 20)
	stripeSprite.Fill(color.RGBA{255, 255, 0, 255})
}

func (g *Game) Update() error {
	if g.gameOver {
		return nil
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.playerCar.X -= 5
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.playerCar.X += 5
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.playerCar.Speed += 0.1
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		g.playerCar.Speed -= 0.2
	}
	if g.playerCar.Speed > 10 {
		g.playerCar.Speed = 10
	}
	if g.playerCar.Speed < 0 {
		g.playerCar.Speed = 0
	}

	g.stripeY += stripeSpeed
	if g.stripeY > 40 { 
		g.stripeY = 0
	}
	var roadLeft float64
	roadLeft = (screenWidth - roadWidth) / 2
	roadRight := roadLeft + roadWidth
	if g.playerCar.X < roadLeft {
		g.playerCar.X = roadLeft
	}
	if g.playerCar.X > roadRight-carWidth {
		g.playerCar.X = roadRight - carWidth
	}

	g.spawnTimer++
	if g.spawnTimer > 60 && g.playerCar.Speed > 0 {
		if rand.Intn(100) < 40 { 
			coinWidth := 30.0
			coinHeight := 30.0
			
			g.coins = append(g.coins, Coin{
				X:      roadLeft + rand.Float64()*(roadWidth-coinWidth),
				Y:      -coinHeight,
				Width:  coinWidth,
				Height: coinHeight,
				Speed:  g.playerCar.Speed,
				Active: true,
			})
		} else {
			var enemyHeight, enemyWidth float64
			enemyWidth = carWidth
			enemyHeight = carHeight
			
			colorChoice := "blue"
			if rand.Intn(100) < 50 {
				colorChoice = "green"
			}
			
			g.enemyCars = append(g.enemyCars, EnemyCar{
				X:      roadLeft + rand.Float64()*(roadWidth-enemyWidth),
				Y:      -enemyHeight,
				Width:  enemyWidth,
				Height: enemyHeight,
				Speed:  g.playerCar.Speed,
				Color:  colorChoice,
			})
		}
		g.spawnTimer = 0
	}

	for i := 0; i < len(g.enemyCars); i++ {
		g.enemyCars[i].Y += g.enemyCars[i].Speed

		if g.playerCar.X < g.enemyCars[i].X+g.enemyCars[i].Width &&
			g.playerCar.X+carWidth > g.enemyCars[i].X &&
			g.playerCar.Y < g.enemyCars[i].Y+g.enemyCars[i].Height &&
			g.playerCar.Y+carHeight > g.enemyCars[i].Y {
			g.gameOver = true
		}
		
		if g.enemyCars[i].Y > screenHeight {
			g.enemyCars = append(g.enemyCars[:i], g.enemyCars[i+1:]...)
			i--
		}
	}

	for i := 0; i < len(g.coins); i++ {
		if !g.coins[i].Active {
			continue
		}
		
		g.coins[i].Y += g.coins[i].Speed
	
		if g.playerCar.X < g.coins[i].X+g.coins[i].Width &&
			g.playerCar.X+carWidth > g.coins[i].X &&
			g.playerCar.Y < g.coins[i].Y+g.coins[i].Height &&
			g.playerCar.Y+carHeight > g.coins[i].Y {
			g.score += 100
			g.coinsCollected++
			g.coins[i].Active = false
			g.coins = append(g.coins[:i], g.coins[i+1:]...)
			i--
			continue
		}
		if g.coins[i].Y > screenHeight {
			g.coins = append(g.coins[:i], g.coins[i+1:]...)
			i--
		}
	}
	g.score += int(g.playerCar.Speed)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{100, 100, 100, 255})

	roadOp := &ebiten.DrawImageOptions{}
	roadOp.GeoM.Translate((screenWidth-roadWidth)/2, 0)
	screen.DrawImage(roadSprite, roadOp)

	for y := -40; y < screenHeight+40; y += 40 {
		posY := float64(y) + g.stripeY
		if posY < -20 || posY > screenHeight {
			continue
		}
		stripeOp := &ebiten.DrawImageOptions{}
		stripeOp.GeoM.Translate(screenWidth/2-5, posY)
		screen.DrawImage(stripeSprite, stripeOp)
	}

	for _, coin := range g.coins {
		if coin.Active {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(coin.X, coin.Y)
	
			scaleX := coin.Width / float64(coinSprite.Bounds().Dx())
			scaleY := coin.Height / float64(coinSprite.Bounds().Dy())
			op.GeoM.Scale(scaleX, scaleY)
			
			screen.DrawImage(coinSprite, op)
		}
	}

	for _, enemy := range g.enemyCars {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(enemy.X, enemy.Y)
		
		if enemy.Color == "blue" {
			screen.DrawImage(blueEnemySprite, op)
		} else {
			screen.DrawImage(greenEnemySprite, op)
		}
	}

	carOp := &ebiten.DrawImageOptions{}
	carOp.GeoM.Translate(g.playerCar.X, g.playerCar.Y)
	screen.DrawImage(playerCarSprite, carOp)

	ebitenutil.DebugPrint(screen, 
		fmt.Sprintf("Speed: %.1f\nScore: %d\nCoins: %d", g.playerCar.Speed, g.score, g.coinsCollected))

	if g.gameOver {
		ebitenutil.DebugPrintAt(screen, "GAME OVER", screenWidth/2-40, screenHeight/2)
		ebitenutil.DebugPrintAt(screen, "Press R to restart", screenWidth/2-60, screenHeight/2+20)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	rand.Seed(42)

	game := &Game{
		playerCar: Car{
			X:      (screenWidth - carWidth) / 2,
			Y:      screenHeight - carHeight - 20,
			Width:  carWidth,
			Height: carHeight,
		},
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Car Game - Collect Coins, Avoid Enemies")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}