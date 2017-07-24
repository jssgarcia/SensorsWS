package main

import (
	"path/filepath"
	"github.com/tkanos/gonfig"
	"ty/csi/ws/SensorsWS/Global"
	"context"
	"time"
	"ty/csi/ws/SensorsWS/TcpClient"
	"github.com/kardianos/service"
	"log"
	lg "ty/csi/ws/SensorsWS/lgg"
)

//region Variables-Modulo
type ctxWrap struct {
	ctx context.Context
	cancel context.CancelFunc
}
var _ctx ctxWrap
//endregion

//region Service-Functions
type program struct{}

func (p *program) Start(s service.Service) error {
	// Start should not block. Do the actual work async.
	go p.run()
	return nil
}
func (p *program) run() {
	// Do work here
}
func (p *program) Stop(s service.Service) error {
	// Stop should not block. Return with a few seconds.
	dispose()
	return nil
}
//endregion


func main() {

	//Obtenemos la configuracion
	lg.InitLogger()

	lg.Lgdef.Info(">> Init Main func >>>")
	getConfig()

	//TODO: initService (Uncomment)
    ctx,cancel := context.WithCancel(context.Background())
	_ctx.ctx=ctx
	_ctx.cancel = cancel

	defer dispose()

	go TcpClient.InitClient2(_ctx.ctx)

	for n:=0;n<5;n++ {
		time.Sleep(time.Second)
	}

	lg.Lgdef.Info("<< Finish Main func <<<")
}

func initService(){

	svcConfig := &service.Config{
		Name:        "GoServiceExampleSimple",
		DisplayName: "Go Service Example",
		Description: "This is an example Go service.",
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}

	err = s.Run()
	if err != nil {
		lg.Lgdef.Error(err)
	}
}

func getConfig() {

	cnfg := Global.Configuration{}


	fconfig := "config.json"

	//Defaults
	cnfg.ServerAddress="127.0.0.1:8888"

	f,_ := filepath.Abs("./" + fconfig)

	//Obtenemos la configuracion
	err := gonfig.GetConf(f,&cnfg)
	if err!=nil {
		panic("ERROR obtener configuracion: " + err.Error())
	}

	Global.Resources.Config = cnfg

}

//region Aux Functions

func dispose(){
	lg.Lgdef.Info(">> Dispose start >>> ")

	_ctx.cancel()  //Provoca llamar a ctx.Done() channel

	lg.Lgdef.Info("<< Dispose end <<  ")


}

//endregion