package Global

//Configuracion
type Configuration struct {
	ServerAddress string `json:"ServerAddress"`
}

type gglobal struct {
	Config Configuration
}

var Resources gglobal
