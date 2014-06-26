package main

import (
	"fmt"
	"iforgetgo/controllers"
	_ "iforgetgo/routers"
	"os"
	"os/signal"
)

func startService() {

	c := controllers.ServiceController{}
	c.InitService()
	// for i:=0;i<N;i++ { запуск нескольких горутин воркеров

	c.Run()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, os.Kill)

	// останов сервера по Ctrl-C
	<-sigChan
	c.Stop()
	fmt.Println("Сервер успешно остановлен.")
	os.Exit(0)
}

func main() {
	startService() //
	//beego.Run()
}
