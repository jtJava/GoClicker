package main

import (
	"github.com/gonutz/w32"
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
	go randomizeChances()
	go monitorToggle()
	go monitorLeftButton()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	click()
}

func click() {
	for {
		time.Sleep(1 * time.Millisecond)
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
	w32.PostMessage(w32.GetForegroundWindow(), 0x0202, 0, 0)
}

func mouseDown() {
	w32.PostMessage(w32.GetForegroundWindow(), 0x0201, 1, 0)
}

func tryChance(chance uint8) bool {
	return uint8(rand.Intn(101-0)+0) <= chance
}

func isMinecraftFocused() bool {
	focusWindow := w32.GetForegroundWindow()
	windowName, _ := w32.GetClassName(focusWindow)
	return windowName == "LWJGL" || windowName == "GLFW30"
}

func monitorLeftButton() {
	/*
	 * 0x01 = Left Mouse Button
	 * 0x04 = Middle Mouse Button
	 * 32768 = SHORT (Pressed)
	 * 0 = SHORT (Not Pressed)
	 */
	for {
		time.Sleep(1 * time.Millisecond)
		state := w32.GetAsyncKeyState(0x01) == 32768

		if state == pressed {
			continue
		}

		if state {
			pressed = true
			firstClick = true
			firstClickDelay = time.Now().Add(50 * time.Millisecond)
		} else {
			pressed = false
		}
	}
}

func monitorToggle() {
	/*
	 * 0x04 = Middle Mouse Button
	 * 32768 = SHORT (Pressed)
	 * 0 = SHORT (Not Pressed)
	 */
	for {
		skip := w32.GetAsyncKeyState(0x04)&0x1 == 0
		if !skip {
			toggled = !toggled
			println(toggled)
		}
	}
}
