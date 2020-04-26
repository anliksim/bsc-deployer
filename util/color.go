package util

import "fmt"

func SetDarkGray() {
	fmt.Printf("\033[1;30m")
}

func SetNoColor() {
	fmt.Printf("\033[0m")
}
