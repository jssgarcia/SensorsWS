package main

import (
	_"path/filepath"
	_"github.com/tkanos/gonfig"
	"ty/csi/ws/SensorsWS/Global"
	"context"
	_"ty/csi/ws/SensorsWS/TcpClient"
	"github.com/kardianos/service"
	lg "ty/csi/ws/SensorsWS/lgg"
	"ty/csi/ws/SensorsWS/server"
	_ "github.com/Sirupsen/logrus"
	"log"
	"os"
	"ty/csi/ws/SensorsWS/TcpClient"
	"path/filepath"
	"github.com/tkanos/gonfig"
	"github.com/kardianos/osext"
	"strings"
	"path"
	"time"
)

//region Variables-Modulo
type ctxWrap struct {
	ctx context.Context
	cancel context.CancelFunc
}
var _ctx ctxWrap
var svclgg service.Logger

//endregion

//region Service-Functions
type program struct{}

func (p *program) Start(s service.Service) error {
	// Start should not block. Do the actual work async.
	go p.run()
	return nil
}
func (p *program) run() {
	initExecution()
}
func (p *program) Stop(s service.Service) error {
	dispose()
	return nil
}
//endregion


func main() {

	//lg.InitLogger(svclgg)

	initService()

	//lg.Lgdef.Info("MAIN: == Finish Main func ===")
}

/*Inicia la Ejecuion de todos los procesos necesarios.
Es llamado desde el metodo RUN */
func initExecution() {

	//Obtenemos la configuracion
	err := lg.InitLogger(svclgg)
	if err!=nil {
		panic(err)
	}
	//
	lg.Lgdef.Info("MAIN: == Init Main func ===\n\n")
	getConfig()

	defer func() {
		lg.Lgdef.Info("MAIN: DEFER called...")
		dispose()
	}()

	_ctx.ctx, _ctx.cancel = context.WithCancel(context.Background())

	//Iniciamos los clientes TPC para escuchar los sensores
	 go TcpClient.InitClient(_ctx.ctx, &TcpClient.ClientInfo{
		ServerName:    "A",
		ServerAddress: Global.Resources.Config.SensorAServerAddress, },
	)
	go TcpClient.InitClient(_ctx.ctx, &TcpClient.ClientInfo{
		ServerName:    "B",
		ServerAddress: Global.Resources.Config.SensorBServerAddress, },
	)
	go TcpClient.InitClient(_ctx.ctx, &TcpClient.ClientInfo{
		ServerName:    "C",
		ServerAddress: Global.Resources.Config.SensorCServerAddress, },
	)

	//Iniciamos el servidor HTTP server
	initHTTPServer()
}

func initHTTPServer() {

	 server.InitServer(
		server.HttpServerInfo{
			EndpointName:    "Default",
			EndpointAddress: Global.Resources.Config.ServerEndpoint,
		})
}


func initService(){

	//lg.Lgdef.Debugf("=== INIT SERVICE ===")

	svcConfig := &service.Config{
		Name:        "CSI-SensorsReaderAndServer",
		DisplayName: "CSI Sensors Reader and HTTP-Server",
		Description: "Permite leer los sensores a traves de TCP connection y ademas provee un Servidor HTTP Restful para servirlos.",
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}

	//Permite obtener el Logger con un Canal, para enviarlos
	chErrs := make(chan error, 5)
	if service.Interactive() {
		svclgg, err = s.Logger(chErrs)
	} else {
		svclgg, err = s.SystemLogger(chErrs)
	}
	if err != nil {
		svclgg.Error(err)
	}
	//Escribe los errores producidos en el Servido, desde el Canal
	go func() {
		for {
			err := <-chErrs
			if err != nil {
				svclgg.Error(err)
			}
		}
	}()

	if len(os.Args) > 1 {
		err = service.Control(s, os.Args[1])
		if err != nil {
			svclgg.Error(err)
		}
		return
	}

	err = s.Run()
	if err != nil {
		svclgg.Error(err)
	}

	//lg.Lgdef.Debugf("=== FINIST SERVICE ===")
}

func getConfig() {

	cnfg := Global.Configuration{}

	//Defaults
	cnfg.SensorAServerAddress ="127.0.0.1:3333"
	cnfg.SensorBServerAddress ="127.0.0.1:3334"
	cnfg.SensorCServerAddress ="127.0.0.1:3335"
	cnfg.ServerEndpoint = "127.0.0.1:9091"

	file := ""
	fname := "config.json"

	if service.Interactive() {
		file,_ = filepath.Abs("./" + fname)
	} else {
		pathexec,_ :=osext.ExecutableFolder();
		file = strings.Replace(path.Join(pathexec,fname),"/","\\",-1)
	}
	//Obtenemos la configuracion
	err := gonfig.GetConf(file,&cnfg)
	if err!=nil {
		panic("ERROR obtener configuracion: " + err.Error())
	}

	svclgg.Infof("Configuración obtenida \n: %s",cnfg)
	lg.Lgdef.Debugf("Configuración obtenida \n: %s",cnfg)

	Global.Resources.Config = cnfg

}

//region Aux Functions

func dispose(){
	lg.Lgdef.Info("MAIN: == Dispose start === ")

	if _ctx.cancel !=nil {
		_ctx.cancel()  //Provoca llamar a ctx.Done() channel
	}

	time.Sleep(2 * time.Second)

	lg.Lgdef.Info("MAIN: == Dispose end ==  ")


}

//endregion