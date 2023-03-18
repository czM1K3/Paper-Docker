package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/czM1K3/Paper-Docker/src/config"
	"github.com/czM1K3/Paper-Docker/src/download"
	"github.com/czM1K3/Paper-Docker/src/files"
	"github.com/czM1K3/Paper-Docker/src/user"
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

	cmd := exec.Command("runuser", "-u", "paper","--","bash", "-c"," cd " + config.DataPath+ " && echo 'eula=true' > eula.txt && java -Xmx"+memoryString+"M -Xms"+memoryString+"M -jar "+ config.PaperPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}
