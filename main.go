package main

import (
	"AutoClicker/windows"
	"github.com/moutend/go-hook/pkg/mouse"
	"github.com/moutend/go-hook/pkg/types"
	"math/rand"
	"os"
	"os/signal"
	"time"
)

var toggled bool
var pressed bool
var firstClick bool
var tickPercentage = make(map[uint8]uint8)
var timeUntilClickable time.Time
var firstClickDelay time.Time

func main() {
	go click()
	go randomizeChances()

	registerHooks()
}

func click() {
	for {
		time.Sleep(1 * 1000 * 1000)
		now := time.Now()
		if isMinecraftFocused() && toggled && pressed && now.After(timeUntilClickable) && now.After(firstClickDelay) {
			if !firstClick {
				mouseDown()
			}
			time.Sleep(time.Millisecond * 50)
			mouseUp()
			firstClick = false
			timeUntilClickable = time.Now().Add(time.Duration(int64(getDelay()) * 1000 * 1000))
		}
	}
}

// 1 tick = 28-55% 2 tick = 82-94% 3 tick = 97-99%
// 3 tick chance increased for more outliers, changed to 90, 99
func randomizeChances() {
	for {
		time.Sleep(time.Second * 10)
		tickPercentage[1] = uint8(rand.Intn(56-28) + 28)
		tickPercentage[2] = uint8(rand.Intn(95-82) + 82)
		tickPercentage[3] = uint8(rand.Intn(100-90) + 90)
	}
}

func getDelay() uint8 {
	if tryChance(tickPercentage[1]) {
		return 0
	} else if tryChance(tickPercentage[2]) {
		return 50
	} else if tryChance(tickPercentage[3]) {
		return 150
	} else {
		return uint8(rand.Intn(6) * 50)
	}
}

func mouseUp() {
	windows.PostMessage(windows.GetForegroundWindow(), 0x0202, 0, 0)
}

func mouseDown() {
	windows.PostMessage(windows.GetForegroundWindow(), 0x0201, 1, 0)
}

func registerHooks() {
	// Buffer size is depends on your need. The 100 is placeholder value.
	mouseChan := make(chan types.MouseEvent, 100)

	if err := mouse.Install(nil, mouseChan); err != nil {
		return
	}

	defer mouse.Uninstall()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	/*
	 * 513 = Left Down
	 * 514 = Left Up
	 * 519 = Scroll Wheel Down
	 */
	for {
		m := <-mouseChan
		if m.Message == 513 {
			pressed = true
			firstClick = true
			firstClickDelay = time.Now().Add(50 * 1000 * 1000)
		} else if m.Message == 514 {
			pressed = false
		} else if m.Message == 519 {
			toggled = !toggled
		}
	}
}

func tryChance(chance uint8) bool {
	return uint8(rand.Intn(101-0)+0) <= chance
}

// Skidded from Fyu
func isMinecraftFocused() bool {
	focusWindow := windows.GetForegroundWindow()
	windowName := windows.GetClassNameW(focusWindow)
	return windowName == "LWJGL" || windowName == "GLFW30"
}
