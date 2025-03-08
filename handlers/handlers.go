package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"spotify/requests"
	"spotify/utils"
)

var (
	templates = template.Must(template.ParseGlob("./templates/*.html"))
	Query     string
	Type      string
)

func RandomHandler(w http.ResponseWriter, r *http.Request) {

	tracks, statusCode, err := requests.RequestRandom()
	if err != nil {
		http.Error(w, "Failed to fetch random", statusCode)
		return
	}
	for index, element := range tracks.Tracks.Items {
		tracks.Tracks.Items[index].ImageUrl = element.Album.Image[0].Url
	}

	templates.ExecuteTemplate(w, "random", tracks)
}
func DetailHandler(w http.ResponseWriter, r *http.Request) {
	ID := r.FormValue("id")

	track, statusCode, err := requests.RequestTrackByID(ID)
	if err != nil {
		fmt.Println("Failed to fetch track details:", err)
		http.Error(w, "Failed to fetch track details", statusCode)
		return
	}

	templates.ExecuteTemplate(w, "detail", track)
}
func RandomTreatmentHandler(w http.ResponseWriter, r *http.Request) {
	ID := r.FormValue("id")
	fmt.Println(ID)

	if !utils.IsAuthenticated() {
		http.Redirect(w, r, "/connexion", http.StatusSeeOther)
		return
	}

	if !utils.AddToFavorites(ID) {
		http.Error(w, "Failed to add to favorites", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/random", http.StatusSeeOther)
}

func RemoveFavoriteHandler(w http.ResponseWriter, r *http.Request) {
	if !utils.IsAuthenticated() {
		http.Redirect(w, r, "/connexion", http.StatusSeeOther)
		return
	}

	trackID := r.URL.Query().Get("id")
	if trackID == "" {
		http.Error(w, "Missing track ID", http.StatusBadRequest)
		return
	}

	err := utils.RemoveFavorite(trackID)
	if err != nil {
		http.Error(w, "Failed to remove favorite", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/random", http.StatusSeeOther)
}

func HeaderHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "header", nil)

}

func AccueilHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "accueil", nil)

}
func AboutHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "about", nil)

}

func CompteHandler(w http.ResponseWriter, r *http.Request) {
	if !utils.IsAuthenticated() {
		http.Redirect(w, r, "/connexion", http.StatusSeeOther)
		return
	}

	favorites, err := utils.GetFavorites()
	if err != nil {
		http.Error(w, "Failed to get favorites", http.StatusInternalServerError)
		return
	}

	var detailedFavorites []requests.Track
	for _, fav := range favorites {
		track, statusCode, err := requests.RequestTrackByID(fav.Name)
		if err != nil {
			fmt.Println("Failed to fetch track details:", err)
			http.Error(w, "Failed to fetch track details", statusCode)
			return
		}
		detailedFavorites = append(detailedFavorites, track)
	}

	data := struct {
		User      utils.User
		Favorites []requests.Track
	}{
		User:      utils.CurrentUser,
		Favorites: detailedFavorites,
	}
	for index, element := range data.Favorites {
		data.Favorites[index].ImageUrl = element.Album.Image[0].Url
	}
	templates.ExecuteTemplate(w, "compte", data)
}

func ConnexionHandler(w http.ResponseWriter, r *http.Request) {

	templates.ExecuteTemplate(w, "connexion", nil)

}

func InscriptionHandler(w http.ResponseWriter, r *http.Request) {
	value := r.FormValue("erreur")
	templates.ExecuteTemplate(w, "inscription", value)
}

func DeconnexionHandler(w http.ResponseWriter, r *http.Request) {
	utils.Deconnexion()
	http.Redirect(w, r, "/", http.StatusSeeOther)

}

type PageRecherche struct {
	SearchItem requests.SearchItem
	IsEmpty    bool
}

func Recherche(w http.ResponseWriter, r *http.Request) {
	Query = r.FormValue("query")
	Type = r.FormValue("type")
	var data PageRecherche = PageRecherche{IsEmpty: true}

	searchData, statusCode, err := requests.RequestRecherche(Query, Type)
	if err != nil && http.StatusBadRequest != statusCode {

		return
	}

	if statusCode != http.StatusBadRequest {
		data.IsEmpty = false
		switch Type {
		case "track":
			data.SearchItem = searchData.Tracks
		case "album":
			data.SearchItem = searchData.Albums
		case "artist":
			data.SearchItem = searchData.Artists
		}
	}
	templates.ExecuteTemplate(w, "recherche", data)

}

func ConnexionTreatmentHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	success, err := utils.Connexion(username, password)
	if err != nil {
		http.Redirect(w, r, "/connexion?erreur=Nom d'utilisateur ou mot de passe incorrect.", http.StatusSeeOther)
		return
	}

	if success {
		utils.CurrentUser = utils.User{Name: username}
		http.Redirect(w, r, "/compte", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/connexion", http.StatusSeeOther)
}

func InscritiontreatmentHandler(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("username") == "" || r.FormValue("password") == "" || r.FormValue("password2") == "" {
		http.Redirect(w, r, "/inscription?erreur=Valeurs manquantes, fais un effort le boss.", http.StatusSeeOther)
		return
	}
	if r.FormValue("password") != r.FormValue("password2") {
		http.Redirect(w, r, "/inscription?erreur=Les 2 mots de passe ne sont pas identiques t'es pas d'accord avec toi même là.", http.StatusSeeOther)
		return
	}
	_, userErr := utils.CreateUser(r.FormValue("username"), r.FormValue("password"))

	if userErr != nil {
		http.Redirect(w, r, "/inscription?erreur=Nom d'utilisateur déjà utilisé.", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/connexion", http.StatusSeeOther)
	}

}

func CompteRemoveFavoriteHandler(w http.ResponseWriter, r *http.Request) {
	if !utils.IsAuthenticated() {
		http.Redirect(w, r, "/connexion", http.StatusSeeOther)
		return
	}

	trackID := r.URL.Query().Get("id")
	if trackID == "" {
		http.Error(w, "Missing track ID", http.StatusBadRequest)
		return
	}

	err := utils.RemoveFavorite(trackID)
	if err != nil {
		http.Error(w, "Failed to remove favorite", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/compte", http.StatusSeeOther)
}
