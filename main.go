package main

import (
	"html/template"
	"log"
	"net/http"
)


var tpl *template.Template

func init() {

	tpl = template.Must(template.ParseGlob("./website/*.html"))
}

func main() {

	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/logout", handleLogout)
	http.HandleFunc("/chat/", handleChat)
	http.HandleFunc("/chat/ws", handleConnections)
	http.HandleFunc("/news/", handleNews)
	http.HandleFunc("/files/", handleFiles)
	http.HandleFunc("/options/", handleOptions)
	http.HandleFunc("/options/change_password", handleChangePassword)

	http.Handle("/pages/", http.StripPrefix("/pages", http.FileServer(http.Dir("./website"))))
	http.Handle("/styles/", http.StripPrefix("/styles", http.FileServer(http.Dir("./website/css"))))
	http.Handle("/scripts/", http.StripPrefix("/scripts", http.FileServer(http.Dir("./website/javascript"))))
	http.Handle("/images/", http.StripPrefix("/images", http.FileServer(http.Dir("./website/pictures"))))
	http.Handle("/sounds/", http.StripPrefix("/sounds", http.FileServer(http.Dir("./website/audio"))))
	http.Handle("/file/", http.StripPrefix("/file", http.FileServer(http.Dir("./website/storage"))))

	go handleMessages()

	log.Fatal(http.ListenAndServe(":8080", nil))

}