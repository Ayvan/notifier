package main

import (
	"fmt"
	"notifier/controllers"
	_ "notifier/routers"
	"os"
	"os/signal"
	"github.com/astaxie/beego"
)

func startService() {

	c := controllers.ServiceController{}
	// инициализируем сервис
	c.InitService()
	// запускаем сервис
	c.Run()

	// создаем канал, получающий событие ОС "завершение процесса" и подписываемся на событие
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, os.Kill)

	// останов сервера по Ctrl-C, если в этот канал пришел сигнал - требуется инициировать остановку сервиса
	<-sigChan
	c.Stop()
	fmt.Println("Сервер успешно остановлен.")
	os.Exit(0)
}

func main() {
	// запуск HTTP-сервера
	go beego.Run()
	// запуск сервиса оповещений
	startService()
}
