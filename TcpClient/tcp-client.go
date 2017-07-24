package TcpClient

import (
	"net"
	"bufio"
	"fmt"
	//"os"
	"log"
	"Learn/SimpleTcpServer/Global"
	"context"
	"time"
	_ "errors"
)

type wsErr struct {
	code int
	err  string
}

func (e *wsErr) Error() string {
	return fmt.Sprintf("%d: %s",e.code,e.err)
}


func InitClient(ctx context.Context) error {

	//var cancelation struct{}
	//var bContinue bool = true

	//for {
	//	select {
	//	case <-ctx.Done():
	//		log.Println("Cancelacion")
	//		return errors.New("Cancelacion")
	//	default:
	//		{
				err := starClient(ctx)

				switch errws := err.(*wsErr); errws.code {
				case 100:
					//Error de conexion
					time.Sleep(5 * time.Second)
					break
		//		}
		//	}
		//}
	}

	return nil
}

	//for bContinue {
	//
	//	//cancelation = <-ctx.Done()
	//	//if cancelation
	//
	//	err := starClient(ctx)
	//
	//	switch errws := err.(*wsErr); errws.code {
	//	case 100:
	//		//Error de conexion
	//		time.Sleep(5 * time.Second)
	//		break
	//
	//
	//	case 999:
	//		//Cancelacion
	//		bContinue=false
	//		break
	//	}
	//}



func starClient(ctx context.Context) error {

	if Global.Resources.Config.ServerAddress=="" {
		panic("ERROR Initiation: Server adddres is not defined")
	}

	conn, err := net.Dial("tcp", Global.Resources.Config.ServerAddress)
	if (err!=nil) {
		return &wsErr{100, "TCP:Server '" + Global.Resources.Config.ServerAddress + "' is not available"}
	}

	defer func(){
		if conn!=nil {
			conn.Close()
		}
	}()

	for {

		log.Println(">>> TCP CLIENT connected ... ")

		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			log.Printf("ERROR %v", err)
			conn.Close()

			return &wsErr{100,"TCP:Server. Connection is closed"}

		} else {
			fmt.Println("TCP CLIENT: Message from Server: " + message)
		}

		//if <-ctx.Done() {
		//	return &wsErr{999,"TCP:Server. Request cancelation. Exit"}
		//}
	}
}
