package user

import (
	"os"
	"os/exec"
	"os/user"
	"strconv"
)

func ValidateUser() {
	userIdRaw := os.Getenv("USER_ID")
	userId, error := strconv.ParseInt(userIdRaw, 10, 32)
	if error != nil {
		userId = 1000
	}
	userIdString := strconv.Itoa(int(userId))

	groupIdRaw := os.Getenv("GROUP_ID")
	groupId, error := strconv.ParseInt(groupIdRaw, 10, 32)
	if error != nil {
		groupId = 1000
	}
	groupIdString := strconv.Itoa(int(groupId))

	_, groupError := user.LookupGroupId(groupIdString)
	if groupError != nil {
		cmd := exec.Command("groupadd", "-g", groupIdString, "minecraft")
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
	}

	_, userError := user.LookupId(userIdString)
	if userError != nil {
		cmd := exec.Command("useradd", "-mG", "minecraft", "-u", userIdString, "paper")
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
	}

}
