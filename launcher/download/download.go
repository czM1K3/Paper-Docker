package download

import (
	"czm1k3/paper-launcher/config"
	"czm1k3/paper-launcher/files"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

func fetchVersions() int {
	response, err := http.Get("https://papermc.io/api/v2/projects/paper/versions/" + config.McVersion)

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
	url := "https://papermc.io/api/v2/projects/paper/versions/" + config.McVersion + "/builds/" + strconv.Itoa(build)
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

func downloadFile(filepath string, url string) error {
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

func DownloadPaper() {
	latestBuild := fetchVersions()
	url, hash := fetchBuild(latestBuild)
	if files.FileExists(config.PaperPath) {
		if files.GetHash(config.PaperPath) == hash {
			fmt.Println("You are already running the latest version of Paper")
		} else {
			fmt.Println("Your version is outdated")
			os.Remove(config.PaperPath)
			downloadFile(config.PaperPath, url)
		}
	} else {
		downloadFile(config.PaperPath, url)
	}
}
