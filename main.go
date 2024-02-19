package main

import (
	server "arifm_operations/server/megaServer"
	time_settings "arifm_operations/server/settings"
	agents "arifm_operations/server/agents"
	"net/http"
)

func main() {
	storage := server.Init()  // объявляю storage
	server.GetStorage(storage) // перекидываю storage в server
	main := http.HandlerFunc(server.MainPageHandler)  // накидываем обработчик главной страницы
	http.Handle("/", main)

	settings := time_settings.Init()  // инициализируем settings
	time_settings.SendSettings(*settings)  // отправляем settings в server/settings
	setPage := http.HandlerFunc(time_settings.SettingsPage)  // накидываем обработчик страницы настроек
	http.Handle("/settings", setPage)

	agents_page := http.HandlerFunc(agents.AgentsHandler)
	http.Handle("/workers", agents_page)
	http.ListenAndServe(":8080", nil)
}
