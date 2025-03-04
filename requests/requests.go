package requests

import (
	"encoding/json"
	"fmt"
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

type DataAlbums struct {
	Albums struct {
		Items []AlbumItem `json:"items"`
	} `json:"albums"`
}

type AlbumItem struct {
	TotalTracks int    `json:"total_tracks"`
	ExternalUrl string `json:"external_url"`
	ImageUrl    string `json:"image_url"`
	Name        string `json:"name"`
	ReleaseDate string `json:"release_date"`
	Artists     []struct {
		Name string `json:"name"`
	} `json:"artists"`
}

type DataTracks struct {
	Tracks struct {
		Items []TrackItem `json:"items"`
	} `json:"tracks"`
}

type TrackItem struct {
	Name        string `json:"name"`
	AlbumName   string `json:"album_name"`
	ExternalUrl string `json:"external_url"`
	ImageUrl    string `json:"image_url"`
	ReleaseDate string `json:"release_date"`
	Artists     []struct {
		Name string `json:"name"`
	} `json:"artists"`
}
type SearchItem struct {
	Items []struct {
		Image []struct {
			Url string `json:"url"`
		} `json:"images"`
		Name string `json:"name"`
		Type string `json:"type"`
	} `json:"items"`
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

func RequestAlbums() (DataAlbums, int, error) {
	url := "https://api.spotify.com/v1/search?q=top+albums&type=album&market=FR&limit=50"
	req, reqErr := http.NewRequest(http.MethodGet, url, nil)
	if reqErr != nil {
		return DataAlbums{}, 500, fmt.Errorf("Erreur lors de l'initialisation de la réquête")
	}

	req.Header.Set("Authorization", _token)

	res, resErr := _httpClient.Do(req)
	if resErr != nil {
		return DataAlbums{}, 500, fmt.Errorf("Erreur lors de l'envois de la réquête")
	}

	defer res.Body.Close()

	if res.StatusCode == 401 {
		errToken := RequestToken()
		if errToken != nil {
			return DataAlbums{}, 500, fmt.Errorf("Erreur lors de la récupération du token")
		}
		return RequestAlbums()
	}

	if res.StatusCode != 200 {
		return DataAlbums{}, res.StatusCode, fmt.Errorf("RequestAlbums - Erreur dans la réponse de la requête code : %d", res.StatusCode)
	}

	var albums DataAlbums

	decodeErr := json.NewDecoder(res.Body).Decode(&albums)
	if decodeErr != nil {
		return DataAlbums{}, 500, fmt.Errorf("RequestAlbums - Erreur lors du décodage des données : %s", decodeErr.Error())
	}

	return albums, res.StatusCode, nil
}

func RequestTrack() (DataTracks, int, error) {
	url := "https://api.spotify.com/v1/search?q=top+tracks&type=track&market=FR&limit=50"
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
		return RequestTrack()
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
	url := fmt.Sprintf("https://api.spotify.com/v1/search?q=%s&type=%s&limit=50", Query, Type)
	url = strings.ReplaceAll(url, " ", "+")
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
	if res.StatusCode == 400 {
		return DataSearch{}, res.StatusCode, fmt.Errorf("RequestTrack - Erreur dans la réponse de la requête code : %d", res.StatusCode)

	}

	if res.StatusCode != 200 {
		return DataSearch{}, res.StatusCode, fmt.Errorf("RequestTrack code : %d", res.StatusCode)
	}

	var Datasearch DataSearch

	decodeErr := json.NewDecoder(res.Body).Decode(&Datasearch)
	if decodeErr != nil {
		return DataSearch{}, 500, fmt.Errorf("RequestTrack - Erreur lors du décodage des données : %s", decodeErr.Error())
	}

	return Datasearch, res.StatusCode, nil
}
