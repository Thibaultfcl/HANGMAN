package hangmanclassic

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"time"
)

type SaveData struct {
	WordToGuess    string
	GuessedWord    []string
	AttemptsLeft   int
	GuessedLetters []string
}

func PlayGame() {
	Menu()

	var input string
	fmt.Scan(&input)

	if input == "1" {
		var wordToGuess string
		var guessedWord []string
		var attemptsLeft int
		var guessedLetters []string

		startWithFlag := false

		if len(os.Args) > 1 && os.Args[1] == "--startWith" {
			startWithFlag = true
			saveData, err := ioutil.ReadFile("save.txt")
			if err == nil {
				var savedGame SaveData
				if err := json.Unmarshal(saveData, &savedGame); err == nil {
					wordToGuess = savedGame.WordToGuess
					guessedWord = savedGame.GuessedWord
					attemptsLeft = savedGame.AttemptsLeft
					guessedLetters = savedGame.GuessedLetters
				}
			}
		}

		if !startWithFlag {
			wordToGuess, _ = ChooseWordAPI()
			wordToGuess = removeAccents(wordToGuess)
			guessedWord = make([]string, len(wordToGuess))
			for i := range guessedWord {
				guessedWord[i] = "_"
			}
			attemptsLeft = 10
			guessedLetters = make([]string, 0)

			NbRevealedLetters := len(wordToGuess)/2 - 1
			if NbRevealedLetters < 1 {
				NbRevealedLetters = 1
			}
			indexLetters := make([]int, 0)
			rand.Seed(time.Now().UnixNano())
			for len(indexLetters) < NbRevealedLetters {
				randomIndex := rand.Intn(len(wordToGuess))
				if !contains(indexLetters, randomIndex) {
					indexLetters = append(indexLetters, randomIndex)
				}
			}
			for _, i := range indexLetters {
				guessedWord[i] = string(wordToGuess[i])
			}
		}

		for attemptsLeft > 0 {
			ClearTerminal()
			fmt.Print("Mot à deviner: ")
			Asciiprinter(strings.Join(guessedWord, " "))
			fmt.Print("\n")
			fmt.Printf("Tentatives restantes: %d\n", attemptsLeft)
			fmt.Printf("Lettres déjà essayées: %s\n", strings.Join(guessedLetters, ", "))
			fmt.Print("\n")
			fmt.Println("Entrez STOP pour sauvgarder et quitter.")
			NthHangman(attemptsLeft)

			var guess string
			fmt.Scanln()
			fmt.Print("Devinez une lettre : ")
			fmt.Scanln(&guess)

			if guess == "STOP" {
				saveData := SaveData{
					WordToGuess:    wordToGuess,
					GuessedWord:    guessedWord,
					AttemptsLeft:   attemptsLeft,
					GuessedLetters: guessedLetters,
				}
				jsonData, _ := json.Marshal(saveData)
				ioutil.WriteFile("save.txt", jsonData, 0644)
				fmt.Println("Partie sauvegardée. Au revoir!")
				fmt.Println("Pour reprendre votre partie tapez : go run . --startWith save.txt")
				return
			}

			containsOnlyLetters, _ := regexp.MatchString("^[a-zA-Z]*$", guess)
			if !containsOnlyLetters {
				fmt.Println("Veuillez entrer seulement des lettres")
				fmt.Print("\n")
				fmt.Println("Appuiez sur une touche pour continuer.")
				GetUserInput()
				continue
			}

			guess = strings.ToLower(guess)

			if len(guess) == 1 {
				if containsStr(guessedLetters, guess) {
					fmt.Println("Vous avez déjà deviné cette lettre.")
					fmt.Print("\n")
					fmt.Println("Appuiez sur une touche pour continuer.")
					GetUserInput()
					continue
				}
				guessedLetters = append(guessedLetters, guess)

				if strings.Contains(wordToGuess, guess) {
					fmt.Println("Bien joué !")
					for i, letter := range wordToGuess {
						if string(letter) == guess {
							guessedWord[i] = guess
							if strings.Join(guessedWord, "") == wordToGuess {
								fmt.Println("Félicitations, vous avez gagné! Le mot était:", wordToGuess)
								return
							}
						}
					}
				} else {
					fmt.Println("Lettre incorrecte...")
					attemptsLeft--
				}
			} else if guess == wordToGuess {
				fmt.Println("Félicitations, vous avez gagné! Le mot était:", wordToGuess)
				return
			} else {
				if containsStr(guessedLetters, guess) {
					fmt.Println("Vous avez déjà essayé ce mot.")
				} else {
					guessedLetters = append(guessedLetters, guess)
					fmt.Println("Mot incorrect...")
					attemptsLeft -= 2
				}
			}

			fmt.Println("Appuiez sur une touche pour continuer.")
			GetUserInput()
		}
		if strings.Join(guessedWord, "") != wordToGuess {
			ClearTerminal()
			fmt.Println("Dommage, vous avez épuisé toutes vos tentatives. Le mot était:", wordToGuess)
			file, err := os.Open("hangman.txt")
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			defer file.Close()
			PrintLine(file, 72, 7)
		}
	}

	if input == "2" {
		return
	} else {
		fmt.Println("Veuillez choisir un des numéros proposés")
	}
}
