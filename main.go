package main

import (
	"IRC/ClientServer/Messages"
	"fmt"
)

func main() {

	fmt.Println("Hello world")

	fmt.Println(Messages.CreateErrorMessage(Messages.IRC_ERR_ILLEGALLENGTH))

	fmt.Println(Messages.MakeLabel("hello"))
}
