package app

import (
	"errors"
	"github.com/myOmikron/RustymonBackend/rpchandler"
	"github.com/myOmikron/echotools/logging"
	"gorm.io/gorm"
	"net"
	"net/http"
	"net/rpc"
)

func InitializeRPC(cliSock *net.Listener, socketPath string, db *gorm.DB, isReloading bool) {
	log := logging.GetLogger("rpc")

	if !isReloading {
		if err := rpc.Register(&rpchandler.RPC{DB: db}); err != nil {
			log.Error(err.Error())
		}
		rpc.HandleHTTP()
	}

	var err error
	if *cliSock, err = net.Listen("unix", socketPath); err != nil {
		log.Error(err.Error())
	} else {
		if err := http.Serve(*cliSock, nil); !errors.Is(err, net.ErrClosed) {
			log.Error(err.Error())
		}
	}

}
