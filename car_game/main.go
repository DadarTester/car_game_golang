package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	screenWidth  = 800
	screenHeight = 600
	carWidth     = 50
	carHeight    = 80
	roadWidth    = 400
	stripeSpeed  = 5.0
	enemySpeed   = 4.0
	minDistance  = 100.0
)

// Объявление глобальных переменных для спрайтов
var (
	playerCarSprite  *ebiten.Image
	blueEnemySprite  *ebiten.Image
	greenEnemySprite *ebiten.Image
	coinSprite       *ebiten.Image
	roadSprite       *ebiten.Image
	stripeSprite     *ebiten.Image
)

// Структура для позиционирования объектов
type Position struct {
	X, Y float64
}

// Структура для размеров объектов
type Size struct {
	Width, Height float64
}

// Структура для игровых объектов
type GameObject struct {
	Position
	Size
}

type Car struct {
	GameObject
	Speed float64
}

type Coin struct {
	GameObject
	Speed  float64
	Active bool
}

type EnemyCar struct {
	GameObject
	Speed  float64
	Sprite *ebiten.Image
}

type Game struct {
	playerCar      Car
	enemyCars      []EnemyCar
	coins          []Coin
	stripeY        float64
	gameOver       bool
	score          int
	spawnTimer     int
	coinsCollected int
	lastEnemyY     float64 // Для отслеживания позиции последней вражеской машины
}

// Создание спрайта машины с заданным цветом
func createCarSprite(mainColor, windowColor, wheelColor, lightColor color.RGBA) *ebiten.Image {
	img := ebiten.NewImage(carWidth, carHeight)
	img.Fill(mainColor)

	// Окна
	ebitenutil.DrawRect(img, 10, 10, 30, 20, windowColor)
	ebitenutil.DrawRect(img, 10, 40, 30, 20, windowColor)

	// Колеса
	ebitenutil.DrawRect(img, 5, 15, 5, 15, wheelColor)
	ebitenutil.DrawRect(img, 40, 15, 5, 15, wheelColor)
	ebitenutil.DrawRect(img, 5, 50, 5, 15, wheelColor)
	ebitenutil.DrawRect(img, 40, 50, 5, 15, wheelColor)

	// Фары
	ebitenutil.DrawRect(img, 5, 5, 5, 5, lightColor)
	ebitenutil.DrawRect(img, 40, 5, 5, 5, lightColor)

	return img
}

// Создание спрайта монетки
func createCoinSprite() *ebiten.Image {
	size := 30
	img := ebiten.NewImage(size, size)

	// Рисуем круглую монетку
	coinColor := color.RGBA{255, 215, 0, 255}
	darkCoinColor := color.RGBA{205, 173, 0, 255}
	center := float64(size / 2)

	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			dx := float64(x) - center
			dy := float64(y) - center
			distance := dx*dx + dy*dy

			if distance <= center*center {
				if distance <= (center/2)*(center/2) {
					img.Set(x, y, darkCoinColor)
				} else {
					img.Set(x, y, coinColor)
				}
			}
		}
	}

	return img
}

// Инициализация спрайтов
func initSprites() {
	// Спрайт игрока
	playerCarSprite = createCarSprite(
		color.RGBA{220, 20, 20, 255},
		color.RGBA{50, 100, 200, 255},
		color.RGBA{20, 20, 20, 255},
		color.RGBA{255, 255, 0, 255},
	)

	// Спрайты вражеских машин
	blueEnemySprite = createCarSprite(
		color.RGBA{20, 20, 220, 255},
		color.RGBA{100, 150, 255, 255},
		color.RGBA{20, 20, 20, 255},
		color.RGBA{255, 0, 0, 255},
	)

	greenEnemySprite = createCarSprite(
		color.RGBA{20, 220, 20, 255},
		color.RGBA{100, 255, 150, 255},
		color.RGBA{20, 20, 20, 255},
		color.RGBA{255, 0, 0, 255},
	)

	// Спрайт монетки
	coinSprite = createCoinSprite()

	// Спрайты дороги и разметки
	roadSprite = ebiten.NewImage(roadWidth, screenHeight)
	roadSprite.Fill(color.RGBA{50, 50, 50, 255})

	stripeSprite = ebiten.NewImage(10, 20)
	stripeSprite.Fill(color.RGBA{255, 255, 0, 255})
}

// Создание вражеской машины с указанным спрайтом
func createEnemyCar(sprite *ebiten.Image, x, y, speed float64) EnemyCar {
	return EnemyCar{
		GameObject: GameObject{
			Position: Position{X: x, Y: y},
			Size:     Size{Width: carWidth, Height: carHeight},
		},
		Speed:  speed,
		Sprite: sprite,
	}
}

// Проверка столкновений между двумя объектами
func checkCollision(a, b GameObject) bool {
	return a.X < b.X+b.Width &&
		a.X+a.Width > b.X &&
		a.Y < b.Y+b.Height &&
		a.Y+a.Height > b.Y
}

// Проверка состояния игры
func (g *Game) checkGameState() {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		os.Exit(0)
	}

	// Перезапуск игры только при gameOver
	if g.gameOver && inpututil.IsKeyJustPressed(ebiten.KeyR) {
		g.resetGame()
	}
}

// Сброс игры
func (g *Game) resetGame() {
	g.playerCar = Car{
		GameObject: GameObject{
			Position: Position{
				X: float64(screenWidth-carWidth) / 2,
				Y: float64(screenHeight - carHeight - 20),
			},
			Size: Size{Width: carWidth, Height: carHeight},
		},
		Speed: 0,
	}
	g.enemyCars = nil
	g.coins = nil
	g.stripeY = 0
	g.gameOver = false
	g.score = 0
	g.spawnTimer = 0
	g.coinsCollected = 0
	g.lastEnemyY = -carHeight // Сбрасываем позицию последней машины врага
}

// Обновление позиции игрока
func (g *Game) updatePlayer() {
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

	// Ограничение скорости
	if g.playerCar.Speed > 10 {
		g.playerCar.Speed = 10
	}
	if g.playerCar.Speed < 0 {
		g.playerCar.Speed = 0
	}

	// Ограничение движения по дороге
	roadLeft := float64(screenWidth-roadWidth) / 2
	roadRight := roadLeft + float64(roadWidth)
	if g.playerCar.X < roadLeft {
		g.playerCar.X = roadLeft
	}
	if g.playerCar.X > roadRight-float64(carWidth) {
		g.playerCar.X = roadRight - float64(carWidth)
	}
}

// Создание нового врага или монетки
func (g *Game) spawnObject() {
	roadLeft := float64(screenWidth-roadWidth) / 2
	// Проверяем, достаточно ли расстояние от последней вражеской машины
	if g.lastEnemyY > -float64(carHeight) && g.lastEnemyY < minDistance {
		return
	}

	if rand.Intn(100) < 30 {
		// Создаем монетку
		coinSize := 30.0
		g.coins = append(g.coins, Coin{
			GameObject: GameObject{
				Position: Position{
					X: roadLeft + rand.Float64()*(float64(roadWidth)-coinSize),
					Y: -coinSize,
				},
				Size: Size{Width: coinSize, Height: coinSize},
			},
			Speed:  g.playerCar.Speed,
			Active: true,
		})
	} else {
		// Создаем вражескую машину
		var enemySprite *ebiten.Image
		if rand.Intn(100) < 50 {
			enemySprite = blueEnemySprite
		} else {
			enemySprite = greenEnemySprite
		}

		enemy := createEnemyCar(
			enemySprite,
			roadLeft+rand.Float64()*(float64(roadWidth)-float64(carWidth)),
			-float64(carHeight),
			g.playerCar.Speed,
		)

		g.enemyCars = append(g.enemyCars, enemy)
		g.lastEnemyY = enemy.Y // Запоминаем позицию новой машины
	}
}

// Обновление вражеских машин
func (g *Game) updateEnemies() {
	for i := 0; i < len(g.enemyCars); i++ {
		g.enemyCars[i].Y += g.enemyCars[i].Speed

		// Проверка столкновения с игроком
		if checkCollision(g.enemyCars[i].GameObject, g.playerCar.GameObject) {
			g.gameOver = true
		}

		// Удаление машин, уехавших за экран
		if g.enemyCars[i].Y > float64(screenHeight) {
			g.enemyCars = append(g.enemyCars[:i], g.enemyCars[i+1:]...)
			i--
		}
	}
}

// Обновление монеток
func (g *Game) updateCoins() {
	for i := 0; i < len(g.coins); i++ {
		if !g.coins[i].Active {
			continue
		}

		g.coins[i].Y += g.coins[i].Speed

		// Проверка сбора монетки
		if checkCollision(g.coins[i].GameObject, g.playerCar.GameObject) {
			g.score += 100
			g.coinsCollected++
			g.coins[i].Active = false
			g.coins = append(g.coins[:i], g.coins[i+1:]...)
			i--
			continue
		}

		// Удаление монеток, уехавших за экран
		if g.coins[i].Y > float64(screenHeight) {
			g.coins = append(g.coins[:i], g.coins[i+1:]...)
			i--
		}
	}
}

func (g *Game) Update() error {
	g.checkGameState()

	if g.gameOver {
		return nil
	}

	g.updatePlayer()

	// Обновление разметки
	g.stripeY += stripeSpeed
	if g.stripeY > 40 {
		g.stripeY = 0
	}

	// Создание новых объектов
	g.spawnTimer++
	if g.spawnTimer > 60 && g.playerCar.Speed > 0 {
		g.spawnObject()
		g.spawnTimer = 0
	}

	g.updateEnemies()
	g.updateCoins()

	// Увеличение счета в зависимости от скорости
	g.score += int(g.playerCar.Speed)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Фон
	screen.Fill(color.RGBA{100, 100, 100, 255})

	// Дорога
	roadOp := &ebiten.DrawImageOptions{}
	roadOp.GeoM.Translate(float64(screenWidth-roadWidth)/2, 0)
	screen.DrawImage(roadSprite, roadOp)

	// Разметка
	for y := -40; y < screenHeight+40; y += 40 {
		posY := float64(y) + g.stripeY
		if posY < -20 || posY > float64(screenHeight) {
			continue
		}
		stripeOp := &ebiten.DrawImageOptions{}
		stripeOp.GeoM.Translate(float64(screenWidth)/2-5, posY)
		screen.DrawImage(stripeSprite, stripeOp)
	}

	// Монетки
	for _, coin := range g.coins {
		if coin.Active {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(coin.X, coin.Y)
			screen.DrawImage(coinSprite, op)
		}
	}

	// Вражеские машины
	for _, enemy := range g.enemyCars {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(enemy.X, enemy.Y)
		screen.DrawImage(enemy.Sprite, op)
	}

	// Машина игрока
	carOp := &ebiten.DrawImageOptions{}
	carOp.GeoM.Translate(g.playerCar.X, g.playerCar.Y)
	screen.DrawImage(playerCarSprite, carOp)

	// Интерфейс
	ebitenutil.DebugPrint(screen,
		fmt.Sprintf("Speed: %.1f\nScore: %d\nCoins: %d", g.playerCar.Speed, g.score, g.coinsCollected))

	// Экран Game Over
	if g.gameOver {
		ebitenutil.DebugPrintAt(screen, "GAME OVER", screenWidth/2-40, screenHeight/2)
		ebitenutil.DebugPrintAt(screen, "Press R to restart", screenWidth/2-60, screenHeight/2+20)
		ebitenutil.DebugPrintAt(screen, "Press ESC to exit", screenWidth/2-60, screenHeight/2+40)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	rand.Seed(time.Now().UnixNano())
	initSprites()

	game := &Game{}
	game.resetGame()

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Car Game - Collect Coins, Avoid Enemies")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
