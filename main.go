package main

import (
	"spotify/handlers"
	temp "spotify/templates"

	"log"
	"net/http"
)

func main() {

	temp.InitTemplates()

	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))

	http.HandleFunc("/topalbums", handlers.TopAlbumsHandler)
	http.HandleFunc("/toptracks", handlers.TopTracksHandler)
	http.HandleFunc("/header", handlers.HeaderHandler)

	//http.HandleFunc("/topartists", topArtistsHandler)
	http.HandleFunc("/recherche", handlers.Recherche)
	http.HandleFunc("/compte", handlers.CompteHandler)
	http.HandleFunc("/compte/deconnexion", handlers.DeconnexionHandler)
	http.HandleFunc("/inscription", handlers.InscriptionHandler)
	http.HandleFunc("/inscription/treatment", handlers.InscritiontreatmentHandler)

	http.HandleFunc("/connexion", handlers.ConnexionHandler)
	http.HandleFunc("/connexion/treatment", handlers.ConnexionTreatmentHandler)
	http.HandleFunc("/", handlers.AccueilHandler)

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
