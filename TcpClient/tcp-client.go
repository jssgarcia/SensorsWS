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
	"bytes"
	"strings"
	"strconv"
	"ty/csi/ws/SensorsWS/Utils"
)

type ClientInfo struct {
	ServerName string
	ServerAddress string
	servertx string
	conn net.Conn
 	lastValue float64  //Ultimo valor recibido
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

	defer func() {
		if r := recover(); r != nil {
			lg.Lgdef.Errorf("UNHANDLER ERROR [TCPClient %s] %s",info.servertx,r)
		}
	}()

	log.Println("=== TCPCLIENT connected to ... " + info.servertx)

	defer func(){
		dispose(info)
	}()

	//Last Value
	info.lastValue = -1

	for {
		select {
		case <-ctx.Done():
			return &wsErr{999, "TCPClient. Connection is closed. " + _vars.servertx}
			break

		default:
			//original
			msg, err := bufio.NewReader(conn).ReadString('\r')

			if err != nil {
				log.Printf("ERROR %v", err)
				conn.Close()

				return &wsErr{100, "TCPClient. Connection is closed. " + _vars.servertx}

			} else {
				processInputData(info,msg)
				//fmt.Printf( msg)
			}
			break
		}
	}
}

func dropCR(data []byte) []byte {
	if len(data) > 0 && data[len(data)-1] == '\r' {
		return data[0 : len(data)-1]
	}
	return data
}


func ScanCRLF(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.Index(data, []byte{'\r'}); i >= 0 {
		// We have a full newline-terminated line.
		return i + 2, dropCR(data[0:i]), nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), dropCR(data), nil
	}
	// Request more data.
	return 0, nil, nil
}

func removeNoValidCharacters(msg string) string {
	b := make([]byte, len(msg))
	var bl int
	for i := 0; i < len(msg); i++ {
		c := msg[i]
		if c >= 32 && c < 127 {
			b[bl] = c
			bl++
		}
	}
	return string(b[:bl])
}

func processInputData(info *ClientInfo,msg string) {

	msg = strings.TrimSpace(removeNoValidCharacters(msg))

	switch info.ServerName {
	case "A","B":

		msg = strings.Replace(msg,".","",-1)  //Eliminamos el .

		num,err := strconv.ParseFloat(msg,16)
		if err!=nil {
			lg.Lgdef.Errorf("[TCPCliente: %s] ERR to convert to number. DataLast received '%s'. ERR: %s\n",info.servertx,msg,err)
		}

		//COMPROBAMOS SI TODOS LOS ELEMENTOS DE LA COLA SON DEL MISMO TIPO.
		// Si es asi, no encolamos mas, hasta que alguno nueva recepcion llegue distinta a todas las posiciones de la Cola
		if !g.Resources.Store.DataQueue[info.ServerName].AllEqual(num) {
			lg.Lgdef.Infof("[TCPCliente: %s] New value data received '%d'\n",info.servertx,num)
			info.lastValue = num

			item :=&g.ItemInfo{Value:info.lastValue,Date: time.Now().Local().Format("02-01-2006 15:04:05")}
			g.Resources.Store.DataLast[info.ServerName]=item
			g.Resources.Store.DataQueue[info.ServerName].Push(item)

			lg.Lgdef.Printf("Queue %s\n",Utils.PrettyPrint(g.Resources.Store.DataQueue[info.ServerName]))
		}

		break

	case "C":
		//Los grados no necesitan que se elimine el punto, es importante. Reemplazamos el
		//msg = strings.Replace(msg,".",",",-1)  //Eliminamos el .

		num,err := strconv.ParseFloat(msg,16)
		if err!=nil {
			lg.Lgdef.Errorf("[TCPCliente: %s] ERR to convert to number. DataLast received '%s'. ERR: %s\n",info.servertx,msg,err)
		}

		//COMPROBAMOS SI TODOS LOS ELEMENTOS DE LA COLA SON DEL MISMO TIPO.
		// Si es asi, no encolamos mas, hasta que alguno nueva recepcion llegue distinta a todas las posiciones de la Cola
		lg.Lgdef.Infof("[TCPCliente: %s] New value data received '%d'\n",info.servertx,num)
		info.lastValue = num

		item :=&g.ItemInfo{Value:info.lastValue,Date: time.Now().Local().Format("02-01-2006 15:04:05")}
		g.Resources.Store.DataLast[info.ServerName]=item
		g.Resources.Store.DataQueue[info.ServerName].Push(item)

		lg.Lgdef.Printf("Queue %s\n",Utils.PrettyPrint(g.Resources.Store.DataQueue[info.ServerName]))

		break
	}




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
