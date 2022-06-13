package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
)

const appPath string = "/app/"
const dataPath string = "/data/"
const defaultPath string = "/default/"
const paperPath string = appPath + "paper.jar"
const mcVersion string = "1.19"

func fetchVersions() int {
	response, err := http.Get("https://papermc.io/api/v2/projects/paper/versions/" + mcVersion)

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	var responseObject Versions
	json.Unmarshal(responseData, &responseObject)

	return responseObject.Builds[len(responseObject.Builds)-1]
}

type Versions struct {
	ProjectID   string `json:"project_id"`
	ProjectName string `json:"project_name"`
	Version     string `json:"version"`
	Builds      []int  `json:"builds"`
}

func fetchBuild(build int) (string, string) {
	url := "https://papermc.io/api/v2/projects/paper/versions/" + mcVersion + "/builds/" + strconv.Itoa(build)
	response, err := http.Get(url)

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	var responseObject Build
	json.Unmarshal(responseData, &responseObject)

	return url + "/downloads/" + responseObject.Downloads.Application.Name, responseObject.Downloads.Application.SHA256
}

type Build struct {
	ProjectID   string    `json:"project_id"`
	ProjectName string    `json:"project_name"`
	Version     string    `json:"version"`
	Build       int       `json:"build"`
	Time        string    `json:"time"`
	Channel     string    `json:"channel"`
	Promoted    string    `json:"ipromoted"`
	Changes     string    `json:"changes"`
	Downloads   Downloads `json:"downloads"`
}

type Downloads struct {
	Application    Download `json:"application"`
	MojangMappings Download `json:"mojang-mappings"`
}

type Download struct {
	Name   string `json:"name"`
	SHA256 string `json:"sha256"`
}

func DownloadFile(filepath string, url string) error {
	fmt.Println("Downloading Paper from: " + url)

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)

	fmt.Println("Paper downloaded")

	return err
}

func FileExists(name string) bool {
	_, err := os.Stat(name)
	if err == nil {
		return true
	}
	return false
}

func GetHash(path string) string {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

func MakeFolder(path string) {
	cmd := exec.Command("mkdir", "-p", path)
	cmd.Output()
}

func Link(source string, target string) {
	cmd := exec.Command("ln", "-s", source, target)
	cmd.Output()
}

func CopyFile(source string, target string) {
	cmd := exec.Command("cp", source, target)
	cmd.Output()
}

func LinkFile(appSource string, dataTarget string) {
	if !FileExists(dataPath + dataTarget) {
		CopyFile(defaultPath+appSource, dataPath+dataTarget)
		Link(dataPath+dataTarget, appPath+appSource)
	} else if !FileExists(appPath + appSource) {
		Link(dataPath+dataTarget, appPath+appSource)
	}
}

func LinkFolder(appSource string, dataTarget string) {
	if !FileExists(dataPath + dataTarget) {
		MakeFolder(dataPath + dataTarget)
		Link(dataPath+dataTarget, appPath+appSource)
	} else if !FileExists(appPath + appSource) {
		Link(dataPath+dataTarget, appPath+appSource)
	}
}

func main() {
	rawMemory := os.Getenv("MEMORY")
	memory, error := strconv.ParseInt(rawMemory, 10, 64)
	if error != nil {
		memory = 2048
	}
	memoryString := strconv.Itoa(int(memory))

	MakeFolder(appPath)
	MakeFolder(dataPath)
	latestBuild := fetchVersions()
	url, hash := fetchBuild(latestBuild)
	if FileExists(paperPath) {
		if GetHash(paperPath) == hash {
			fmt.Println("You are already running the latest version of Paper")
		} else {
			fmt.Println("Your version is outdated")
			os.Remove(paperPath)
			DownloadFile(paperPath, url)
		}
	} else {
		DownloadFile(paperPath, url)
	}

	MakeFolder(dataPath + "config/")
	MakeFolder(dataPath + "save/")

	LinkFile("banned-ips.json", "config/banned-ips.json")
	LinkFile("banned-players.json", "config/banned-players.json")
	LinkFile("bukkit.yml", "config/bukkit.yml")
	LinkFile("eula.txt", "config/eula.txt")
	LinkFile("ops.json", "config/ops.json")
	LinkFile("permissions.yml", "config/permissions.yml")
	LinkFile("server.properties", "config/server.properties")
	LinkFile("spigot.yml", "config/spigot.yml")
	LinkFile("whitelist.json", "config/whitelist.json")

	LinkFolder("world", "save/world")
	LinkFolder("world_nether", "save/world_nether")
	LinkFolder("world_the_end", "save/world_the_end")
	LinkFolder("plugins", "plugins")
	LinkFolder("logs", "logs")

	fmt.Println("Starting Paper")
	cmd := exec.Command("java", "-Xmx"+memoryString+"M", "-Xms"+memoryString+"M", "-jar", paperPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run()
}
