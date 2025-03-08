package requests

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

var _httpClient = http.Client{
	Timeout: 5 * time.Second,
}

var ClientId string = "2df1c031f64f4ebcb65e3a7b604d8d58"
var ClientSecret string = "0dfac3f4cc324c36bc516d4d7f37db88"

var _token string = "Bearer BQDojtAoJKiy30WBtBG0jjZiLf2MmZE_StI1RoUn6fvxW_GWSW5iAoMSfnYlFnTpis9_w74xaHhx6tOwGEIbQI5wnOZbeXefwwFynmOM9D7lDvfAxA8"

type DataToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

type DataTracks struct {
	Tracks struct {
		Items []TrackItem `json:"items"`
	} `json:"tracks"`
}

type TrackItem struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	AlbumName   string `json:"album_name"`
	ExternalUrl string `json:"external_url"`
	Album       struct {
		Image []struct {
			Url string `json:"url"`
		} `json:"images"`
	} `json:"album"`
	ImageUrl    string
	ReleaseDate string `json:"release_date"`
	Artists     []struct {
		Name string `json:"name"`
	} `json:"artists"`
}
type SearchItem struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	AlbumName   string `json:"album_name"`
	ExternalUrl string `json:"external_url"`
	Album       struct {
		Image []struct {
			Url string `json:"url"`
		} `json:"images"`
	} `json:"album"`
	ImageUrl    string
	ReleaseDate string `json:"release_date"`
	Artists     []struct {
		Name string `json:"name"`
	} `json:"artists"`
}
type DataSearch struct {
	Albums  SearchItem `json:"albums"`
	Tracks  SearchItem `json:"tracks"`
	Artists SearchItem `json:"artists"`
}

func RequestToken() error {
	body := strings.NewReader(fmt.Sprintf("grant_type=client_credentials&client_id=%s&client_secret=%s", ClientId, ClientSecret))

	req, reqErr := http.NewRequest(http.MethodPost, "https://accounts.spotify.com/api/token", body)
	if reqErr != nil {
		return fmt.Errorf("RequestToken - Erreur lors de l'initialisation de la réquête")
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, resErr := _httpClient.Do(req)
	if resErr != nil {
		return fmt.Errorf("RequestToken - Erreur lors de l'envois de la réquête")
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		return fmt.Errorf("RequestToken - Erreur dans la réponse de la requête code : %d", res.StatusCode)
	}

	var data DataToken

	decodeErr := json.NewDecoder(res.Body).Decode(&data)
	if decodeErr != nil {
		return fmt.Errorf("RequestToken - Erreur lors du décodage des données")
	}

	_token = fmt.Sprintf("%s %s", data.TokenType, data.AccessToken)
	return nil
}

func RequestRandom() (DataTracks, int, error) {

	musique := []string{
		"Michael+Jackson", "Beyoncé", "Rihanna", "Eminem", "Tupac+Shakur", "Drake",
		"The+Beatles", "Queen", "Nirvana", "The+Rolling+Stones", "Metallica",
		"Daft+Punk", "David+Guetta", "Skrillex", "The+Weeknd", "Calvin+Harris",
		"Whitney+Houston", "Alicia+Keys", "Usher", "Stromae", "Mylène+Farmer",
		"Angèle", "Edith+Piaf", "Jacques+Brel", "Johnny+Hallyday", "Madonna",
		"Bob+Marley", "Elvis+Presley", "Travis+Scott", "BTS", "Kanye+West",
		"AC/DC", "The+Ramones", "Marvin+Gaye", "Johnny+Cash", "Carl+Cox",
		"Pink+Floyd", "Lana+Del+Rey", "Frank+Ocean", "Céline+Dion", "Bruno+Mars",
		"Coldplay", "Imagine+Dragons", "Post+Malone", "Doja+Cat",
		"Pop", "Rock", "Rap/Hip-Hop", "R&B", "Électro/House", "Jazz", "Blues",
		"Reggae", "Funk", "Soul", "Métal", "Punk", "Chanson+française", "K-pop",
		"Musique+classique", "Country", "Trap", "Techno", "Hard+Rock", "Disco",
		"Grunge", "Dubstep", "Dancehall", "Afrobeat", "Indie+Rock",
	}

	rand.Shuffle(len(musique), func(i, j int) { musique[i], musique[j] = musique[j], musique[i] })

	randomIndex := rand.Intn(len(musique))
	Query := musique[randomIndex]
	url := fmt.Sprintf("https://api.spotify.com/v1/search?q=%s&type=track&market=FR&limit=50", Query)
	req, reqErr := http.NewRequest(http.MethodGet, url, nil)
	if reqErr != nil {
		return DataTracks{}, 500, fmt.Errorf("Erreur lors de l'initialisation de la réquête")
	}

	req.Header.Set("Authorization", _token)

	res, resErr := _httpClient.Do(req)
	if resErr != nil {
		return DataTracks{}, 500, fmt.Errorf("Erreur lors de l'envois de la réquête")
	}

	defer res.Body.Close()

	if res.StatusCode == 401 {
		errToken := RequestToken()
		if errToken != nil {
			return DataTracks{}, 500, fmt.Errorf("Erreur lors de la récupération du token")
		}
		return RequestRandom()
	}

	if res.StatusCode != 200 {
		return DataTracks{}, res.StatusCode, fmt.Errorf("RequestTrack - Erreur dans la réponse de la requête code : %d", res.StatusCode)
	}

	var data DataTracks

	decodeErr := json.NewDecoder(res.Body).Decode(&data)
	if decodeErr != nil {
		return DataTracks{}, 500, fmt.Errorf("RequestTrack - Erreur lors du décodage des données : %s", decodeErr.Error())
	}
	return data, res.StatusCode, nil
}

func RequestRecherche(Query string, Type string) (DataSearch, int, error) {
	url := fmt.Sprintf("https://api.spotify.com/v1/search?q=%s&type=%s", Query, Type)

	req, reqErr := http.NewRequest(http.MethodGet, url, nil)
	if reqErr != nil {
		return DataSearch{}, 500, fmt.Errorf("Erreur lors de l'initialisation de la réquête")
	}

	req.Header.Set("Authorization", _token)

	res, resErr := _httpClient.Do(req)
	if resErr != nil {
		return DataSearch{}, 500, fmt.Errorf("Erreur lors de l'envois de la réquête")
	}
	defer res.Body.Close()

	if res.StatusCode == 401 {
		errToken := RequestToken()
		if errToken != nil {
			return DataSearch{}, 500, fmt.Errorf("Erreur lors de la récupération du token")
		}
		return RequestRecherche(Query, Type)
	}

	if res.StatusCode != 200 {
		return DataSearch{}, res.StatusCode, fmt.Errorf("RequestTrack - Erreur dans la réponse de la requête code : %d", res.StatusCode)
	}

	var Datasearch DataSearch

	decodeErr := json.NewDecoder(res.Body).Decode(&Datasearch)
	if decodeErr != nil {
		return DataSearch{}, 500, fmt.Errorf("RequestTrack - Erreur lors du décodage des données : %s", decodeErr.Error())
	}

	return Datasearch, res.StatusCode, nil
}

type Track struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Album struct {
		Image []struct {
			Url string `json:"url"`
		} `json:"images"`
	} `json:"album"`
	ImageUrl    string
	ReleaseDate string `json:"release_date"`
	Artists     []struct {
		Name string `json:"name"`
	} `json:"artists"`
}

func RequestTrackByID(id string) (Track, int, error) {
	url := fmt.Sprintf("https://api.spotify.com/v1/tracks/%s", id)
	req, reqErr := http.NewRequest(http.MethodGet, url, nil)
	if reqErr != nil {
		return Track{}, 500, fmt.Errorf("Erreur lors de l'initialisation de la réquête")
	}

	req.Header.Set("Authorization", _token)

	res, resErr := _httpClient.Do(req)
	if resErr != nil {
		return Track{}, 500, fmt.Errorf("Erreur lors de l'envois de la réquête")
	}
	defer res.Body.Close()

	if res.StatusCode == 401 {
		errToken := RequestToken()
		if errToken != nil {
			return Track{}, 500, fmt.Errorf("Erreur lors de la récupération du token")
		}
		return RequestTrackByID(id)
	}

	if res.StatusCode != 200 {
		return Track{}, res.StatusCode, fmt.Errorf("RequestTrackByID - Erreur dans la réponse de la requête code : %d", res.StatusCode)
	}

	var track Track
	decodeErr := json.NewDecoder(res.Body).Decode(&track)
	if decodeErr != nil {
		return Track{}, 500, fmt.Errorf("RequestTrackByID - Erreur lors du décodage des données : %s", decodeErr.Error())
	}

	return track, res.StatusCode, nil
}
