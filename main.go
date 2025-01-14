package main

import (
	"api/handlers"
	"log"
	"net/http"
)

func main() {

	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))

	http.HandleFunc("/topalbums", handlers.TopAlbumsHandler)
	http.HandleFunc("/toptracks", handlers.TopTracksHandler)
	http.HandleFunc("/header", handlers.HeaderHandler)
	http.HandleFunc("/", handlers.AccueilHandler)
	//http.HandleFunc("/topartists", topArtistsHandler)
	//http.HandleFunc("/recherche", Recherche )
	//http.HandleFunc("/resultat", Resultat )
	http.HandleFunc("/compte", handlers.CompteHandler)
	http.HandleFunc("/inscription", handlers.InscriptionHandler)
	http.HandleFunc("/connexion", handlers.ConnexionHandler)

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
