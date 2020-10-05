package main

import (
	"fmt"
	"regexp"
)

func keepAlive(msg string) {
	fmt.Printf("%v\n", msg)
	fmt.Scanln()
}

func separateIntoCommand(command string) (cmd string, user string, params string) {
	commandRegex := regexp.MustCompile(".*")
	soleCommandRegex := regexp.MustCompile(" _([a-z]|[0-9])* ?")
	clearupRegex := regexp.MustCompile("_[a-z]*")
	userRegex := regexp.MustCompile("<.*>")
	userIDRegex := regexp.MustCompile("[0-9]+")
	paramRegex := regexp.MustCompile("([a-z]|[0-9])*$")

	res := commandRegex.Find([]byte(command))
	cmd = string(clearupRegex.Find(soleCommandRegex.Find(res)))
	user = string(userIDRegex.Find(userRegex.Find(res)))
	params = string(paramRegex.Find(res))
	return
}
