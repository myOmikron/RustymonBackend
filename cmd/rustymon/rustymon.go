package main

import (
	"fmt"
	"github.com/hellflame/argparse"
	"github.com/myOmikron/RustymonBackend/rpcrequests"
)

func main() {
	parser := argparse.NewParser("rustymon", "CLI tool for rustymon-server", &argparse.ParserConfig{})

	registerParser := parser.AddCommand("register", "Register a user", &argparse.ParserConfig{})
	usernameRegister := registerParser.String("u", "username", &argparse.Option{
		Required: true,
		Help:     "Username used for logging in the user",
	})
	passwordRegister := registerParser.String("p", "password", &argparse.Option{
		Required: true,
		Help:     "Initial password of the user",
	})
	emailRegister := registerParser.String("m", "mail", &argparse.Option{
		Required: true,
		Help:     "Mail address of the user",
	})
	trainerNameRegister := registerParser.String("t", "trainer-name", &argparse.Option{
		Required: true,
		Help:     "Display name of the user",
	})

	if err := parser.Parse(nil); err != nil {
		fmt.Println(err.Error())
		return
	}

	switch {
	case registerParser.Invoked:
		rpcrequests.Register(*usernameRegister, *passwordRegister, *emailRegister, *trainerNameRegister)
	}
}
