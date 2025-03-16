package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"spotify/requests"
	"spotify/utils"
	"strconv"
)

var (
	templates = template.Must(template.ParseGlob("./templates/*.html"))
	Query     string
	Type      string
)

type RandomPageData struct {
	Tracks     []requests.TrackItem
	PrevOffset int
	NextOffset int
	Offset     int
}

func RandomHandler(w http.ResponseWriter, r *http.Request) {
	offsetStr := r.URL.Query().Get("offset")
	offset := 0

	if offsetStr != "" {
		parsedOffset, err := strconv.Atoi(offsetStr)
		if err == nil {
			offset = parsedOffset
		}
	}
	data, statusCode, err := requests.RequestRandom(offset)
	if err != nil {
		http.Error(w, "Failed to fetch random", statusCode)
		return
	}
	if offset < 0 {
		offset = 10
	}
	if offset > 300 {
		offset = 200
	}
	prevOffset := offset - 10
	if prevOffset < 0 {
		prevOffset = 0
	}
	nextOffset := offset + 10
	for index, element := range data.Tracks.Items {
		fmt.Println(element)
		data.Tracks.Items[index].ImageUrl = element.Album.Image[0].Url
	}
	pageData := RandomPageData{
		Tracks:     data.Tracks.Items,
		PrevOffset: prevOffset,
		NextOffset: nextOffset,
		Offset:     offset,
	}
	fmt.Println(pageData.Tracks[0])

	templates.ExecuteTemplate(w, "random", pageData)

}
func DetailHandler(w http.ResponseWriter, r *http.Request) {
	ID := r.FormValue("id")
	Type := r.FormValue("type")

	track, statusCode, err := requests.RequestTrackByID(ID, Type)
	if err != nil {
		fmt.Println("Failed to fetch track details:", err)
		http.Error(w, "Failed to fetch track details", statusCode)
		return

	}

	track.ImageUrl = track.Album.Image[0].Url

	templates.ExecuteTemplate(w, "detail", track)
}
func RandomTreatmentHandler(w http.ResponseWriter, r *http.Request) {
	ID := r.FormValue("id")
	Type := r.FormValue("type")
	fmt.Println(ID)

	if !utils.IsAuthenticated() {
		http.Redirect(w, r, "/connexion", http.StatusSeeOther)
		return
	}

	if !utils.AddToFavorites(ID, Type) {
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
		if fav.Type == "unknown" {
			fav.Type = "tracks"
		}
		track, statusCode, err := requests.RequestTrackByID(fav.Name, fav.Type)
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
	SearchItemT  requests.SearchItemT
	SearchItemAL requests.SearchItemAL
	SearchItemA  requests.SearchItemA
	IsEmpty      bool
}

func Recherche(w http.ResponseWriter, r *http.Request) {
	Query = r.FormValue("query")
	Type = r.FormValue("type")
	var data PageRecherche = PageRecherche{IsEmpty: true}

	searchData, statusCode, err := requests.RequestRecherche(Query, Type)
	if err != nil && http.StatusBadRequest != statusCode {
		data.IsEmpty = true

		return
	}

	if statusCode != http.StatusBadRequest {
		data.IsEmpty = false
		switch Type {
		case "track":
			data.SearchItemT = searchData.Tracks
			templates.ExecuteTemplate(w, "rechercheT", data)
		case "album":
			data.SearchItemAL = searchData.Albums
			templates.ExecuteTemplate(w, "rechercheAL", data)
		case "artist":
			data.SearchItemA = searchData.Artists
			templates.ExecuteTemplate(w, "rechercheA", data)
		default:
			templates.ExecuteTemplate(w, "rechercheT", data)
		}

	}

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
