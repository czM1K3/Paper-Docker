package main

import (
	"czm1k3/paper-launcher/config"
	"czm1k3/paper-launcher/download"
	"czm1k3/paper-launcher/files"
	"czm1k3/paper-launcher/user"
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

func main() {
	rawMemory := os.Getenv("MEMORY")
	memory, error := strconv.ParseInt(rawMemory, 10, 64)
	if error != nil {
		memory = 2048
	}
	memoryString := strconv.Itoa(int(memory))

	user.ValidateUser()

	files.MakeFolder(config.DataPath)

	download.DownloadPaper()

	files.Chmod()

	fmt.Println("Starting Paper")
	fmt.Println()

	cmd := exec.Command("runuser", "-u", "paper", "--", "cd", "/data", "&&", "java", "-Xmx"+memoryString+"M", "-Xms"+memoryString+"M", "-jar", config.PaperPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}
