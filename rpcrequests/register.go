package rpcrequests

import (
	"fmt"
	"github.com/myOmikron/RustymonBackend/rpchandler"
	"github.com/myOmikron/echotools/color"
	"log"
	"net/rpc"
	"os"
)

func Register(sockPath, username, password, email, trainerName string) {
	if conn, err := rpc.DialHTTP("unix", sockPath); err != nil {
		log.Fatalln("dialing error: ", err)
	} else {
		req := rpchandler.RegisterRequest{
			Username:    username,
			Email:       email,
			Password:    password,
			TrainerName: trainerName,
		}
		var res rpchandler.RegisterResult

		if err := conn.Call("RPC.RegisterUser", req, &res); err != nil {
			color.Println(color.RED, "Error returned:")
			fmt.Println(err.Error())
			if res.ErrorMessage != "" {
				fmt.Println("Additional message:", res.ErrorMessage)
			}
			os.Exit(1)
		}

		color.Println(color.PURPLE, "User registered successfully!")
	}
}
