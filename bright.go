package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"time"
)

const (
	PATH       = "/sys/class/backlight/intel_backlight/brightness"
	MIN        = 0
	MAX        = 7500
	PCTMIN     = 2
	PCTMAX     = 100
	STEP   int = MAX / 100
	STEPS      = 20
)

var (
	LOGMIN = math.Log10(PCTMIN)
	LOGMAX = math.Log10(PCTMAX)
)

// GetBrightness returns the brightness in %.
func GetBrightness() (int, error) {
	f, err := os.ReadFile(PATH)
	if err != nil {
		return -1, fmt.Errorf("failed to read brightness: %v", err)
	}
	n, err := strconv.Atoi(string(f)[:len(f)-1])
	if err != nil {
		return -1, fmt.Errorf("didn't get a number: %v", err)
	}
	return n, nil
}

func GetBrightnessPct() (float64, error) {
	i, err := GetBrightness()
	if err != nil {
		return -1, err
	}
	return (float64(i) / float64(MAX)) * 100.0, nil
}

func SetBrightness(val int) error {
	if val > MAX {
		return fmt.Errorf("value greater than max brightness %d", MAX)
	}
	s := strconv.Itoa(val)
	return os.WriteFile(PATH, []byte(s), 0664)
}

func SetBrightnessPct(pct float64) error {
	return SetBrightness(int((pct / 100) * MAX))
}

func SetBrightnessSmooth(val int) error {
	current, err := GetBrightness()
	if err != nil {
		return err
	}
	if val > current {
		for ; current < val; current += STEP {
			SetBrightness(current)
			time.Sleep(10 * time.Millisecond)
		}
	} else {
		for ; current > val; current -= STEP {
			SetBrightness(current)
			time.Sleep(10 * time.Millisecond)
		}
	}
	return nil
}

func BrightnessPctToStep(pct float64) float64 {
	return math.Round(math.Log10(pct) / (LOGMAX - LOGMIN) * STEPS)
}

func StepToBrightnessPct(step float64) float64 {
	x := step / STEPS * (LOGMAX - LOGMIN)
	return math.Round(math.Max(math.Min(math.Pow(10, x), PCTMAX), PCTMIN))
}

// Usage prints the program usage, then exits with error code 1.
func Usage() {
	fmt.Printf("usage\n\t%s -U / -B / -D\n", os.Args[0])
	os.Exit(1)
}

func main() {
	if len(os.Args) < 2 || !(os.Args[1] == "-U" || os.Args[1] == "-D" || os.Args[1] == "-B") {
		Usage()
	}

	currentBrightness, err := GetBrightnessPct()
	if err != nil {
		fmt.Println(err)
	}

	if currentBrightness == 0 {
		SetBrightnessPct(PCTMIN + 1)
		currentBrightness, err = GetBrightnessPct()
		if err != nil {
			fmt.Println(err)
		}
	}

	currentStep := BrightnessPctToStep(currentBrightness)

	newStep := currentStep

	switch os.Args[1] {
	case "-U":
		newStep += 2
	case "-D":
		newStep -= 2
	case "-B":
		newStep += 5
	}

	newBrightness := StepToBrightnessPct(newStep)

	fmt.Printf("Current backlight: %.1f\nChanging to %.1f\n", currentBrightness, newBrightness)

	if newBrightness == 0 {
		os.Exit(1)
	}
	if os.Args[1] == "-U" && currentBrightness == 99.0 {
		newBrightness = 100.0
	}
	err = SetBrightnessSmooth(int((newBrightness / 100.0) * MAX))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
