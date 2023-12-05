package hangmanclassic

import (
	"fmt"
)

func logo() {
	Asciiprinter("HANGMAN")
}

func Menu() {
	ClearTerminal()
	logo()
	fmt.Print("\n")
	fmt.Println("1. Nouvelle partie")
	fmt.Println("2. Quitter")
	fmt.Print("\n")
}
