package TcpClient

import (
	"net"
	"bufio"
	"fmt"
	//"os"
	"log"
	"ty/csi/ws/SensorsWS/Global"
	"context"
	"time"
	_ "errors"
	"errors"
	lg "ty/csi/ws/SensorsWS/lgg"
)

type wsErr struct {
	code int
	err  string
}

func (e *wsErr) Error() string {
	return fmt.Sprintf("%d: %s",e.code,e.err)
}

type vars struct {
	conn net.Conn
}

var _vars vars

func InitClient2(ctx context.Context) {

	for {
		select {
		case <-ctx.Done():
			lg.Lgdef.Warn("Cancelación recibida: Salimos")
		dispose()
		}
	}

}

func InitClient(ctx context.Context) error {

	//var cancelation struct{}
	//var bContinue bool = true
	ticker := time.NewTicker(5 * time.Second)

	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			lg.Lgdef.Warn("Cancelación recibida: Salimos")
			dispose()
			return errors.New("Cancelacion")

		case <- ticker.C:
				err := starClient(ctx)
				if err!=nil {
					lg.Lgdef.Error(err)
				}
				//
				//switch errws := err.(*wsErr); errws.code {
				//case 999:
				//	//Error de conexion
				//	time.Sleep(5 * time.Second)
				//	break
				//}
			break
		}
	}

	return nil
}

func starClient(ctx context.Context) error {

	if Global.Resources.Config.ServerAddress=="" {
		panic("ERROR Initiation: Server adddres is not defined")
	}

	conn, err := net.Dial("tcp", Global.Resources.Config.ServerAddress)
	_vars.conn = conn

	if (err!=nil) {
		return &wsErr{100, "TCP:Server '" + Global.Resources.Config.ServerAddress + "' is not available"}
	}

	log.Println(">>> TCP CLIENT connected to ... " + Global.Resources.Config.ServerAddress )

	defer func(){
		dispose()
	}()

	for {
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

func dispose() {
	lg.Lgdef.Debug(">> Dispose INIT >>> ")

	if _vars.conn!=nil {
		_vars.conn.Close()
	}

	lg.Lgdef.Debug("<<< Dispose FINISH <<< ")

}
