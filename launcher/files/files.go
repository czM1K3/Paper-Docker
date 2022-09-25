package files

import (
	"crypto/sha256"
	"czm1k3/paper-launcher/config"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
)

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
	if !FileExists(config.DataPath + dataTarget) {
		CopyFile(config.DefaultPath+appSource, config.DataPath+dataTarget)
		Link(config.DataPath+dataTarget, config.AppPath+appSource)
	} else if !FileExists(config.AppPath + appSource) {
		Link(config.DataPath+dataTarget, config.AppPath+appSource)
	}
}

func LinkFolder(appSource string, dataTarget string) {
	if !FileExists(config.DataPath + dataTarget) {
		MakeFolder(config.DataPath + dataTarget)
		Link(config.DataPath+dataTarget, config.AppPath+appSource)
	} else if !FileExists(config.AppPath + appSource) {
		Link(config.DataPath+dataTarget, config.AppPath+appSource)
	}
}

func Chmod() {
	cmd1 := exec.Command("chown", "-R", "paper:minecraft", config.DataPath)
	cmd1.Stdin = os.Stdin
	cmd1.Stdout = os.Stdout
	cmd1.Stderr = os.Stderr
	cmd1.Run()

	cmd2 := exec.Command("chown", "-R", "paper:minecraft", config.AppPath)
	cmd2.Stdin = os.Stdin
	cmd2.Stdout = os.Stdout
	cmd2.Stderr = os.Stderr
	cmd2.Run()
}
