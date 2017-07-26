package Global

import (
	"time"
)
//Configuracion
type Configuration struct {
	SensorAServerAddress string `json:"SensorServerAddress"` //Direccion tcp del servidor envia datos del A
	SensorBServerAddress string `json:"SensorServerAddress"` //Direccion tcp del servidor envia datos del B
	SensorCServerAddress string `json:"SensorServerAddress"` //Direccion tcp del servidor envia datos del C
	ServerEndpoint      string `json:"ServerEndpoint"`	//Endpoint addres where listen HTTP request
}

type ItemInfo struct {
	Value int	`json:"Value"`
	Date  time.Time  `json:"TimeReceived"`
}
type Store struct {
	Data map[string]*ItemInfo
}

type Global struct {
	Config Configuration
	Store Store
}

var Resources Global

//var DB  map[string]*ItemInfo

func init() {
	Resources.Store.Data = make(map[string]*ItemInfo)

	Resources.Store.Data["A"]= &ItemInfo{5,time.Now()}
	Resources.Store.Data["B"]= &ItemInfo{5,time.Now()}
	Resources.Store.Data["C"]= &ItemInfo{5,time.Now()}

	//DB= make(map[string]*ItemInfo)
	//
	//DB["A"]=&ItemInfo{Value:5,Date: time.Now()}
	//DB["B"]=&ItemInfo{Value:5,Date: time.Now()}
	//DB["C"]=&ItemInfo{Value:5,Date: time.Now()}
}
