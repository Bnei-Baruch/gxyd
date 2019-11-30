package main

const (
	EnvBackendUrl = "BACKEND_URL"
	EnvLogPath	  = "LOG_PATH"
	envListenPort = "LISTEN_PORT"
)

func main() {
	a := App{}
	a.Initialize()
	a.Run()
}
