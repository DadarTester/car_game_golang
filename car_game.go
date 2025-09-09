package main
import (
    "fmt"
    "image/color"
    "log"
    "github.com/hajimehoshi/ebiten/v2"
    "github.com/hajimehoshi/ebiten/v2/ebitenutil"
    "github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
    screenWidth = 800
    screenHeight = 600
    carWidth = 50
    carHeight = 80
    roadWidth = 400
)

type Car struct {
    X, Y   float64
    Speed  float64
    Width  float64
    Height float64
}

type Game struct {
    car      Car
    roadX    float64
    gameOver bool
    score    int
}

func (g *Game) Update() error {
    if g.gameOver {
        return nil
    }

    if ebiten.IsKeyPressed(ebiten.KeyLeft) {
        g.car.X -= 5
    }
    if ebiten.IsKeyPressed(ebiten.KeyRight) {
        g.car.X += 5
    }
    if ebiten.IsKeyPressed(ebiten.KeyUp) {
        g.car.Speed += 0.1
    }
    if ebiten.IsKeyPressed(ebiten.KeyDown) {
        g.car.Speed -= 0.2
    }

    if g.car.Speed > 10 {
        g.car.Speed = 10
    }
    if g.car.Speed < 0 {
        g.car.Speed = 0
    }

    g.roadX -= g.car.Speed
    if g.roadX < roadWidth {
        g.roadX += roadWidth
    }

    if g.car.X < (screenWidth-roadWidth)/2 {
        g.car.X = (screenWidth - roadWidth) / 2
    }
    if g.car.X > (screenWidth+roadWidth) / 2 - carWidth {
        g.car.X = (screenWidth+roadWidth) / 2 - carWidth
    }

    g.score += int(g.car.Speed)

    return nil
}

func(g *Game) Draw(screen *ebiten.Image) {
    screen.Fill(color.RGBA{100, 100, 100, 255})             // Серый

    roadRect := ebiten.NewImage(roadWidth, screenHeight)
    roadRect.Fill(color.RGBA{50, 50, 50, 255})              // Темно Серый
    op := &ebiten.DrawImageOptions{}
    op.GeoM.Translate((screenWidth-roadWidth)/2, 0)
    screen.DrawImage(roadRect, op)
    stripeWidth := 10.0
    for y := 0; y < screenHeight; y += 40 {                //Вот тут хуйня какая-то, отображает часть строки комментарием
        ebitenutil.DrawRect(screen, 
            screenWidth/2 - stripeWidth/2, 
            float64(y) + g.roadX, 
            stripeWidth, 20, color.RGBA{255, 255, 0, 255})
    }

    ebitenutil.DrawRect(screen, g.car.X, g.car.Y, carWidth, carHeight, color.RGBA{255, 0, 0, 255})

    ebitenutil.DebugPrint(screen, fmt.Sprintf("Speed: %.1f\nScore: %d", g.car.Speed, g.score))

    if g.gameOver {
        ebitenutil.DebugPrintAt(screen, "GAME OVER", screenWidth/2-40, screenHeight/2)
        ebitenutil.DebugPrintAt(screen, "Press R to restart", screenWidth/2-60, screenHeight/2+20)
    }
}

func(g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
    return screenWidth, screenHeight
}

type Obstacle struct {
    X, Y   float64
    Width  float64
    Height float64
    Speed  float64
}

type Game struct {
    car        Car
    obstacles  []Obstacle
    roadX      float64
    gameOver   bool
    score      int
    obstacleTimer int
}

func (g *Game) Update() error {

    g.obstacleTimer++
    if g.obstacleTimer > 60 { 
        width := 40.0
        g.obstacles = append(g.obstacles, Obstacle{
            X:      (screenWidth-roadWidth)/2 + rand.Float64()*(roadWidth-width),
            Y:      -50,
            Width:  width,
            Height: 60,
            Speed:  g.car.Speed,
        })
        g.obstacleTimer = 0
    }

    for i, _ := range g.obstacles {
        g.obstacles[i].Y += g.obstacles[i].Speed
        
        if g.car.X < g.obstacles[i].X+g.obstacles[i].Width && g.car.X+carWidth > g.obstacles[i].X && .car.Y < g.obstacles[i].Y+g.obstacles[i].Height &&
            g.car.Y+carHeight > g.obstacles[i].Y {
            g.gameOver = true
        }
        
        if g.obstacles[i].Y > screenHeight {
            g.obstacles = append(g.obstacles[:i], g.obstacles[i+1:]...)
            break
        }
    }

    return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
    
    for _, obs := range g.obstacles {
        ebitenutil.DrawRect(screen, obs.X, obs.Y, obs.Width, obs.Height, color.RGBA{0, 0, 255, 255})
    }
    
}

func main() {
    game := &Game {
        car: Car {
            X: (screenWidth - carWidth) / 2,
            Y: screenHeight - carHeight - 20,
            Width: carWidth,
            Height: carHeight,
        }
    }

    ebiten.SetWindowSize(screenWidth, screenHeight)
    ebiten.SetWindowTitle("Car Game")

    if err := ebiten.RunGame(game); err != nil {                      //Такая же проблема с отображением комментария
        log.Fatal(err)
    }
}

//Вот тут должен быть запуск, но вместо запуска...ни-че-го