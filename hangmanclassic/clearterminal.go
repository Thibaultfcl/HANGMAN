package hangmanclassic

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

func ClearTerminal() {
	osName := runtime.GOOS

	var cmd *exec.Cmd

	switch osName {
	case "linux", "darwin":
		cmd = exec.Command("clear")
	case "windows":
		cmd = exec.Command("cmd", "/c", "cls")
	default:
		fmt.Println("Système d'exploitation non pris en charge")
		return
	}

	cmd.Stdout = os.Stdout
	cmd.Run()
}
