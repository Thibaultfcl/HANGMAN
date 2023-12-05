package hangmanclassic

import (
	"bufio"
	"fmt"
	"os"
)

func Asciiprinter(Word string) {
	for hauteur := 2; hauteur <= 10; hauteur++ {
		for _, letter := range Word {
			Showletter(letter, hauteur)
		}
		fmt.Println("")
	}
}


func Showletter(letter rune, line int) {
	var numero int
	file, err := os.Open("ascii.txt")
	if err != nil {
		fmt.Println("Erreur lors de l'ouverture du fichier :", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	lineCount := 1
	startLine := 0

	numero = int(letter) - 32

	startLine = (numero * 9) + line

	for scanner.Scan() {
		lineCount++
		if lineCount == startLine {
			fmt.Print(scanner.Text())
		}
	}

	if scanner.Err() != nil {
		fmt.Println("Erreur lors de la lecture du fichier :", scanner.Err())
	}
}
