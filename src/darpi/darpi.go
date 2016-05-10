package main

import (
	"fmt"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/gpio"
	"github.com/hybridgroup/gobot/platforms/raspi"
	"time"
)

const (
	E_DELAY = 500000 // 0.0005 seconds
	E_PULSE = 500000

	LCD_RS = "26" // gpio 7
	LCD_E  = "24" // gpio 8
	LCD_D4 = "22" // gpio 25
	LCD_D5 = "18" // gpio 24
	LCD_D6 = "16" // gpio 23
	LCD_D7 = "12" // gpio 18
	LED_ON = "10" // gpio 15  backlight

	LCD_WIDTH = 20 // Maximum characters per line
	LCD_CHR   = 1
	LCD_CMD   = 0

	LCD_LINE_1 = 0x80 // LCD RAM address for the 1st line
	LCD_LINE_2 = 0xC0 // LCD RAM address for the 2nd line
	LCD_LINE_3 = 0x94 // LCD RAM address for the 3rd line
	LCD_LINE_4 = 0xD4 // LCD RAM address for the 4th line
)

type Pins struct {
	E         *gpio.DirectPinDriver
	RS        *gpio.DirectPinDriver
	D4        *gpio.DirectPinDriver
	D5        *gpio.DirectPinDriver
	D6        *gpio.DirectPinDriver
	D7        *gpio.DirectPinDriver
	BackLight *gpio.DirectPinDriver
}

func main() {
	fmt.Printf("Starting darpi!\n")

	gbot := gobot.NewGobot()

	r := raspi.NewRaspiAdaptor("raspi")

	pins := &Pins{}

	pins.E = gpio.NewDirectPinDriver(r, "e", LCD_E)
	pins.RS = gpio.NewDirectPinDriver(r, "rs", LCD_RS)
	pins.D4 = gpio.NewDirectPinDriver(r, "d4", LCD_D4)
	pins.D5 = gpio.NewDirectPinDriver(r, "d5", LCD_D5)
	pins.D6 = gpio.NewDirectPinDriver(r, "d6", LCD_D6)
	pins.D7 = gpio.NewDirectPinDriver(r, "d7", LCD_D7)
	pins.BackLight = gpio.NewDirectPinDriver(r, "backlight", LED_ON)

	work := func() {
		fmt.Printf("Doing work\n")

		r.Connect()

		lcdInit(pins)

		pins.BackLight.DigitalWrite(0)
		time.Sleep(1000 * 1000 * 1000)
		pins.BackLight.DigitalWrite(1)
		time.Sleep(1000 * 1000 * 1000)

		lcdString(pins, "Hello GoLangPhilly!", LCD_LINE_1, 1)
		lcdString(pins, "Hello DramaFever!", LCD_LINE_2, 1)
		lcdString(pins, "Are you having fun!", LCD_LINE_3, 1)
		lcdString(pins, "I hope so!", LCD_LINE_4, 1)

		fmt.Println("Exiting work")
	}

	robot := gobot.NewRobot("mypi",
		[]gobot.Connection{r},
		[]gobot.Device{pins.E, pins.RS, pins.D4, pins.D5, pins.D6, pins.D7, pins.BackLight},
		work,
	)

	gbot.AddRobot(robot)

	gbot.Start()

}

// lcdInit initializes the display
func lcdInit(pins *Pins) {
	fmt.Println("lcdByte 0x33")
	lcdByte(pins, 0x33, 0) // 110011 Initialise
	fmt.Println("lcdByte 0x32")
	lcdByte(pins, 0x32, LCD_CMD) // 110010 Initialise
	fmt.Println("lcdByte 0x06")
	lcdByte(pins, 0x06, LCD_CMD) // 000110 Cursor move direction
	fmt.Println("lcdByte 0x0C")
	lcdByte(pins, 0x0C, LCD_CMD) // 001100 Display On,Cursor Off, Blink Off
	fmt.Println("lcdByte 0x28")
	lcdByte(pins, 0x28, LCD_CMD) // 101000 Data length, number of lines, font size
	fmt.Println("lcdByte 0x01")
	lcdByte(pins, 0x01, LCD_CMD) // 000001 Clear display
	time.Sleep(E_DELAY)
}

// flip the Enable pin
func lcdToggleEnable(pins *Pins) {
	time.Sleep(E_DELAY)
	pins.E.DigitalWrite(1)
	time.Sleep(E_PULSE)
	pins.E.DigitalWrite(0)
	time.Sleep(E_DELAY)
}

// lcdByte sends byte to data pins.  bits = data, mode=true for character, false for command
func lcdByte(pins *Pins, bits int, mode byte) {

	pins.RS.DigitalWrite(mode)

	// high bits

	pins.D4.DigitalWrite(0)
	pins.D5.DigitalWrite(0)
	pins.D6.DigitalWrite(0)
	pins.D7.DigitalWrite(0)
	if bits&0x10 == 0x10 {
		pins.D4.DigitalWrite(1)
	}
	if bits&0x20 == 0x20 {
		pins.D5.DigitalWrite(1)
	}
	if bits&0x40 == 0x40 {
		pins.D6.DigitalWrite(1)
	}
	if bits&0x80 == 0x80 {
		pins.D7.DigitalWrite(1)
	}

	// toggle enable
	lcdToggleEnable(pins)

	// low bits

	pins.D4.DigitalWrite(0)
	pins.D5.DigitalWrite(0)
	pins.D6.DigitalWrite(0)
	pins.D7.DigitalWrite(0)
	if bits&0x01 == 0x01 {
		pins.D4.DigitalWrite(1)
	}
	if bits&0x02 == 0x02 {
		pins.D5.DigitalWrite(1)
	}
	if bits&0x04 == 0x04 {
		pins.D6.DigitalWrite(1)
	}
	if bits&0x08 == 0x08 {
		pins.D7.DigitalWrite(1)
	}

	lcdToggleEnable(pins)
}

func lcdString(pins *Pins, message string, line int, style int) {
	lcdByte(pins, line, LCD_CMD)

	// pad the message with spaces to the right
	message = fmt.Sprintf("%-20s", message)

	for i := 0; i < LCD_WIDTH; i++ {
		lcdByte(pins, int(message[i]), LCD_CHR)
	}

}
