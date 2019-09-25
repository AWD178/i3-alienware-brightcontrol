package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

func ExecuteXrands(cmdArgs []string) (string, error) {
	path, err := exec.LookPath("xrandr")
	if err != nil {
		return "", err
	}
	cmd := exec.Command(path, cmdArgs...)
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	err = cmd.Run()
	if err != nil {
		return "", err
	}
	return outb.String(), nil
}

func GetCurrentBrightness() float64 {
	outputData, err := ExecuteXrands([]string{"--verbose"})
	if err != nil {
		log.Fatal(`error execute xrandr: %s`, err)
		return 0
	}

	for _, l := range strings.Split(outputData, "\n") {
		if strings.Index(strings.ToLower(l), "brightness") >= 0 {
			re := regexp.MustCompile(`(\d.\d)`)
			rResults := re.FindAllString(l, -1)
			if len(rResults) > 0 {
				d, e := strconv.ParseFloat(rResults[0], 64)
				if e != nil {
					log.Fatal(`error convert brightness to number: %s`, e)
				}
				return d
			}
		}
	}

	return 0
}

//xrandr --output eDP1 --brightness 1
func main() {
	scriptArgs := os.Args

	if len(scriptArgs) == 3 {
		currentBright := GetCurrentBrightness()
		percentArgs, err := strconv.Atoi(scriptArgs[2])
		if err != nil {
			log.Fatal("brightness must be a number")
		} else {
			switch scriptArgs[1] {
			case "inc":
				currentBright = currentBright + (float64(percentArgs) / 100)
			case "dec":
				currentBright = currentBright - (float64(percentArgs) / 100)
			}

			if currentBright > 1 {
				currentBright = 1
			}

			if currentBright < 0 {
				currentBright = 0
			}

			xrandrParams := []string{"--output", "eDP1", "--brightness", fmt.Sprintf("%f", currentBright)}
			_, err = ExecuteXrands(xrandrParams)
			if err != nil {
				log.Fatal(`error execute xrandr: %s`, err)
			}
		}

	}

	if len(scriptArgs) == 2 {
		percentArgs, err := strconv.Atoi(scriptArgs[1])
		if err != nil {
			log.Fatal("brightness must be a number")
		} else {
			brightness := float64(percentArgs) / 100

			if brightness > 1 {
				brightness = 1
			}

			if brightness < 0 {
				brightness = 0
			}

			xrandrParams := []string{"--output", "eDP1", "--brightness", fmt.Sprintf("%f", brightness)}
			_, err = ExecuteXrands(xrandrParams)
			if err != nil {
				log.Fatal(`error execute xrandr: %s`, err)
			}
		}
	}
}
