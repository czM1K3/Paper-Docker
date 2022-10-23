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

	files.MakeFolder(config.AppPath)
	files.MakeFolder(config.DataPath)

	download.DownloadPaper()

	files.MakeFolder(config.DataPath + "config/")
	files.MakeFolder(config.DataPath + "save/")
	files.MakeFolder(config.DataPath + "bluemap/")

	files.LinkFile("banned-ips.json", "config/banned-ips.json")
	files.LinkFile("banned-players.json", "config/banned-players.json")
	files.LinkFile("bukkit.yml", "config/bukkit.yml")
	files.LinkFile("commands.yml", "config/commands.yml")
	files.LinkFile("eula.txt", "config/eula.txt")
	files.LinkFile("ops.json", "config/ops.json")
	files.LinkFile("permissions.yml", "config/permissions.yml")
	files.LinkFile("server.properties", "config/server.properties")
	files.LinkFile("spigot.yml", "config/spigot.yml")
	files.LinkFile("whitelist.json", "config/whitelist.json")
	files.LinkFile("server-icon.png", "config/server-icon.png")

	files.LinkFolder("world", "save/world")
	files.LinkFolder("world_nether", "save/world_nether")
	files.LinkFolder("world_the_end", "save/world_the_end")
	files.LinkFolder("plugins", "plugins")
	files.LinkFolder("bluemap", "bluemap")
	files.LinkFolder("logs", "logs")

	files.LinkFile("world/paper-world.yml", "save/world/paper-world.yml")
	files.LinkFile("world_nether/paper-world.yml", "save/world_nether/paper-world.yml")
	files.LinkFile("world_the_end/paper-world.yml", "save/world_the_end/paper-world.yml")

	files.Chmod()

	fmt.Println("Starting Paper")
	fmt.Println()

	cmd := exec.Command("runuser", "-u", "paper", "--", "java", "-Xmx"+memoryString+"M", "-Xms"+memoryString+"M", "-jar", config.PaperPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}
