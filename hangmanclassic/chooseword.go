package hangmanclassic

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

func ChooseWord() string {
	f, err := os.Open("wordlist.txt")

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	var words []string

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		word := strings.TrimSpace(scanner.Text())
		words = append(words, word)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	rand.Seed(time.Now().UnixNano())

	randomIndex := rand.Intn(len(words))

	choosenWord := removeAccents(words[randomIndex])

	return choosenWord
}

func removeAccents(input string) string {
	accentMap := map[rune]rune{
		'à': 'a', 'á': 'a', 'â': 'a', 'ã': 'a', 'ä': 'a', 'å': 'a',
		'ç': 'c',
		'è': 'e', 'é': 'e', 'ê': 'e', 'ë': 'e',
		'ì': 'i', 'í': 'i', 'î': 'i', 'ï': 'i',
		'ñ': 'n',
		'ò': 'o', 'ó': 'o', 'ô': 'o', 'õ': 'o', 'ö': 'o',
		'ù': 'u', 'ú': 'u', 'û': 'u', 'ü': 'u',
		'ý': 'y',
		'À': 'A', 'Á': 'A', 'Â': 'A', 'Ã': 'A', 'Ä': 'A', 'Å': 'A',
		'Ç': 'C',
		'È': 'E', 'É': 'E', 'Ê': 'E', 'Ë': 'E',
		'Ì': 'I', 'Í': 'I', 'Î': 'I', 'Ï': 'I',
		'Ñ': 'N',
		'Ò': 'O', 'Ó': 'O', 'Ô': 'O', 'Õ': 'O', 'Ö': 'O',
		'Ù': 'U', 'Ú': 'U', 'Û': 'U', 'Ü': 'U',
		'Ý': 'Y',
	}

	var result strings.Builder
	for _, r := range input {
		if replacement, ok := accentMap[r]; ok {
			result.WriteRune(replacement)
		} else {
			result.WriteRune(r)
		}
	}

	return result.String()
}

func ChooseWordAPI() (string, error) {
	var apiUrl string
	var lang string
	var choice string

	ClearTerminal()
	fmt.Scanln()
	fmt.Println("Veuillez choisir une langue :")
	fmt.Println("1. Français")
	fmt.Println("2. English")
	fmt.Scanln(&choice)

	if choice == "1" {
		lang = "FR"
	}
	if choice == "2" {
		lang = "EN"
	}

	if lang == "FR" {
		apiUrl = "https://trouve-mot.fr/api/random"
	} else {
		apiUrl = "https://random-word-api.herokuapp.com/word"
	}

	response, err := http.Get(apiUrl)
	if err != nil {
		return "", fmt.Errorf("erreur de requête à l'API: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("erreur de requête à l'API: %s", response.Status)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("erreur lors de la lecture de la réponse de l'API: %v", err)
	}

	if lang == "FR" {
		var data []struct {
			Name string `json:"name"`
		}

		err = json.Unmarshal(body, &data)
		if err != nil {
			return "", fmt.Errorf("erreur d'analyse JSON : %v", err)
		}

		if len(data) > 0 {
			return data[0].Name, nil
		}

		return "", fmt.Errorf("réponse de l'API ne contient pas de champ 'name'")
	}
	if lang == "EN" {
		var words []string

		err = json.Unmarshal(body, &words)
		if err != nil {
			return "", fmt.Errorf("erreur d'analyse JSON : %v", err)
		}

		if len(words) > 0 {
			return words[0], nil
		}

		return "", fmt.Errorf("réponse de l'API ne contient pas de mot")
	}

	return "", nil
}
