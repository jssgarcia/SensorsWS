package TcpClient

import (
	"net"
	"bufio"
	"fmt"
	//"os"
	"log"
	"context"
	"time"
	_ "errors"
	lg "ty/csi/ws/SensorsWS/lgg"
	g "ty/csi/ws/SensorsWS/Global"
	"math/rand"
)

type ClientInfo struct {
	ServerName string
	ServerAddress string
	servertx string
	conn net.Conn
}

type wsErr struct {
	code int
	err  string
}

func (e *wsErr) Error() string {
	return fmt.Sprintf("%d: %s",e.code,e.err)
}

type vars struct {
	conn net.Conn
	servertx string
}

var _vars vars

//func InitClient2(ctx context.Context) {
//
//	for {
//		select {
//		case <-ctx.Done():
//			lg.Lgdef.Warn("Cancelación recibida: Salimos")
//			dispose(info)
//			return
//		}
//	}
//
//}

func InitClient(ctx context.Context,info *ClientInfo) {

	_vars.servertx = info.ServerName + " [" + info.ServerAddress + "]"
	info.servertx = info.ServerName + " [" + info.ServerAddress + "]";

	lg.Lgdef.Info("TCPClient: INIT " + info.servertx)

	//var cancelation struct{}
	//var bContinue bool = true
	ticker := time.NewTicker(5 * time.Second)

	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			lg.Lgdef.Warnf("TCPClient(%s): Cancelación recibida: Salimos",info.servertx)
			dispose(info)
			return

		case <- ticker.C:
				err := starClient(ctx,info)
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
}

func starClient(ctx context.Context,info *ClientInfo) error {

	if info.ServerAddress =="" {
		panic("ERROR Initiation: Server adddres is not defined. " + info.servertx)
	}

	conn, err := net.Dial("tcp", info.ServerAddress)
	_vars.conn = conn

	if (err!=nil) {
		return &wsErr{100, "TCPClient: Server '" + info.servertx + "' is not available"}
	}

	log.Println("=== TCPCLIENT connected to ... " + info.servertx)

	defer func(){
		dispose(info)
	}()

	for {
		select {
		case <-ctx.Done():
			return &wsErr{999, "TCPClient. Connection is closed. " + _vars.servertx}
			break

		default:
			message, err := bufio.NewReader(conn).ReadString('\n')
			if err != nil {
				log.Printf("ERROR %v", err)
				conn.Close()

				return &wsErr{100, "TCPClient. Connection is closed. " + _vars.servertx}

			} else {
				processInputData(info,message)
			}
			break
		}
	}
}

func processInputData(info *ClientInfo,msg string) {
	rand.Seed(time.Now().Unix())

	item :=&g.ItemInfo{Value:rand.Intn(1000-100)+100,Date: time.Now()}
	g.Resources.Store.Data[info.ServerName]=item

	fmt.Printf("TCPClient: Message from Server(%s): %s\n",info.ServerName, item)
}

func closeConn(info *ClientInfo) {
	if info.conn!=nil {
		info.conn.Close()
	}
}

func dispose(info *ClientInfo) {
	lg.Lgdef.Debugf("TCPLient(%s): == Dispose INIT === ",_vars.servertx)
	closeConn(info)
	lg.Lgdef.Debugf("TCPLient(%s): === Dispose FINISH === ",_vars.servertx)
}
