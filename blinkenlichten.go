package main

import (
	"bufio"
	"io"
	"os"
	"strconv"
	
	"github.com/codegangsta/cli"
	"github.com/tarm/goserial"
	"fmt"
)

var (
	port io.ReadWriteCloser
	responses *bufio.Scanner
)

func main() {
	app := cli.NewApp()
	app.Name = "blinkenlichten"
	app.Usage = "blink lights"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name: "device, d",
			Usage: "device file to use for serial communications",
			EnvVar: "ARDUINO_PORT",
		},
	}
	app.Commands = []cli.Command{
		{
			Name: "count-leds",
			ShortName: "c",
			Usage: "prints the number of LEDs available",
			Action: countLEDs,
		},
		{
			Name: "on",
			ShortName: "+",
			Usage: "turn on all the lights",
			Action: handleOn,
		},
		{
			Name: "off",
			ShortName: "-",
			Usage: "turn off all the lights",
			Action: handleOff,
		},
	}

	app.Before = prepareConnection

	app.Run(os.Args)
}

// Wrapper for printing out the LED count as a CLI command.
func countLEDs(c *cli.Context) {
	println(GetLEDCount())
}
func GetLEDCount() uint16 {
	port.Write([]byte("c\n"))
	response, _ := getLine()
	num, _ := strconv.ParseUint(response, 10, 16)
	return uint16(num)
}

func handleOn(c *cli.Context) {
	SetAllLEDColors(255, 255, 255)
}
func handleOff(c *cli.Context) {
	SetAllLEDColors(0, 0, 0)
}

func SetAllLEDColors(r,g,b uint8) {
	nleds := GetLEDCount()
	var led uint16
	for led = 0; led < nleds; led++ {
		SetLEDColor(led, r, g, b)
	}
}
func SetLEDColor(n uint16, r,g,b uint8) {
	fmt.Fprintf(port, "s %d %d %d %d\n", n, r, g, b)
}

func getLine() (string, error) {
	if responses.Scan() {
		return responses.Text(), nil
	}
	return "", responses.Err()
}

func prepareConnection(c *cli.Context) error {
	device := c.String("device")
	if device == "" {
		device = "/dev/ttyACM0" // device Arduino shows up as on my Raspberry Pi. #laziness
	}
	baud := c.Int("baud")
	if baud == 0 {
		baud = 9600
	}

	var err error
	port, err = connect(device, baud)
	if err != nil {
		return err
	}
	
	responses = bufio.NewScanner(port)
	getLine() // handle reading in the initial "ready" (FIXME when that doesn't erroneously happen anymore)
	return err
}

func connect(device string, baud int) (io.ReadWriteCloser, error) {
	config := &serial.Config{Name: device, Baud: baud}
	return serial.OpenPort(config)
}
