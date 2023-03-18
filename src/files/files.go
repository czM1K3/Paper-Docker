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
	return err == nil
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

func Chmod() {
	cmd := exec.Command("chown", "-R", "paper:minecraft", config.DataPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}
