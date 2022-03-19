package main

import (
	"github.com/hellflame/argparse"
	"github.com/myOmikron/RustymonBackend/server"
	"os"
)

func main() {
	parser := argparse.NewParser("", "", nil)

	startParser := parser.AddCommand("start", "Start the server", &argparse.ParserConfig{
		DisableDefaultShowHelp: true,
	})

	if err := parser.Parse(nil); err != nil {
		os.Exit(0)
	}

	switch {
	case startParser.Invoked:
		server.StartServer()
	}
}
