package main

import (
	server "arifm_operations/server/megaServer"
	time_settings "arifm_operations/server/settings"
	"net/http"
)

func main() {
	storage := server.Init()  // объявляю storage
	server.GetStorage(*storage) // перекидываю storage в server
	main := http.HandlerFunc(server.GettingResponse)
	http.Handle("/", main)

	settings := time_settings.Init()
	time_settings.SendSettings(*settings)
	setPage := http.HandlerFunc(time_settings.SettingsPage)
	http.Handle("/settings", setPage)
	http.ListenAndServe(":8080", nil)
}
