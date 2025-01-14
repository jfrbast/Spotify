package handlers

import (
	_func "api/func"
	"api/requests"
	"fmt"
	"html/template"
	"net/http"
)

var (
	templates = template.Must(template.ParseGlob("./templates/*.html"))
)

func TopAlbumsHandler(w http.ResponseWriter, r *http.Request) {

	albums, statusCode, err := requests.RequestAlbums()
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to fetch top albums", statusCode)
		return
	}

	templates.ExecuteTemplate(w, "topalbums", albums)
}

func TopTracksHandler(w http.ResponseWriter, r *http.Request) {

	tracks, statusCode, err := requests.RequestTrack()
	if err != nil {
		http.Error(w, "Failed to fetch top tracks", statusCode)
		return
	}
	templates.ExecuteTemplate(w, "toptracks", tracks)
}

/*
func topArtistsHandler(w http.ResponseWriter, r *http.Request) {


	tracks, statusCode, err := requests.Request()
	if err != nil {
		http.Error(w, "Failed to fetch top tracks", statusCode)
		return
	}

	templates.ExecuteTemplate(w, "topartists", tracks)
}
*/

func HeaderHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "header", nil)

}

func AccueilHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "accueil", nil)

}

func CompteHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "compte", nil)

}

func ConnexionHandler(w http.ResponseWriter, r *http.Request) {

	templates.ExecuteTemplate(w, "connexion", nil)

}

func InscriptionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		r.ParseForm()
		username := r.FormValue("username")
		password := r.FormValue("password")
		password2 := r.FormValue("password2")
		if password == password2 {
			_func.EcrireCompte(username, password)
			http.Redirect(w, r, "/compte", http.StatusSeeOther)
			return
		}

	}
	templates.ExecuteTemplate(w, "inscription", nil)
}
