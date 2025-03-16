package utils

import (
	"encoding/json"
	"fmt"
	"os"
)

type Fav struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

var CurrentUser User

func CreateUser(name string, password string) (User, error) {
	res, err := os.ReadFile("comptes.json")
	if err != nil {
		fmt.Println(err)
		return User{}, fmt.Errorf("CreateUser - erreur de lecture du fichier : %s", err.Error())
	}

	var liste []User
	err = json.Unmarshal(res, &liste)
	if err != nil {
		fmt.Println(err)
		return User{}, fmt.Errorf("CreateUser - erreur de décodage du fichier : %s", err.Error())
	}

	var newUser User = User{Name: name, Password: HashPassword(password), Favoris: []Favoris{}}

	if IsFree := Isusable(newUser.Name, liste); !IsFree {
		return User{}, fmt.Errorf("CreateUser - nom d'utilisateur déjà utilisé")
	}

	liste = append(liste, newUser)

	dataEncode, err := json.Marshal(liste)
	if err != nil {
		fmt.Println(err)
	}

	err = os.WriteFile("comptes.json", dataEncode, 0644)
	if err != nil {
		fmt.Println(err)
	}

	return newUser, nil
}

func Isusable(username string, listeUsers []User) bool {
	for _, User := range listeUsers {
		if User.Name == username {
			return false
		}
	}

	return true
}

func Connexion(username string, password string) (bool, error) {
	res, err := os.ReadFile("comptes.json")
	if err != nil {
		return false, fmt.Errorf("Connexion - erreur de lecture du fichier : %s", err.Error())
	}

	var listUsers []User
	err = json.Unmarshal(res, &listUsers)
	if err != nil {
		return false, fmt.Errorf("Connexion - erreur de décodage du fichier : %s", err.Error())
	}

	for _, user := range listUsers {
		if user.Name == username {
			if user.Password == HashPassword(password) {
				CurrentUser = user
				return true, nil
			}
			return false, fmt.Errorf("Connexion - mot de passe incorrect")
		}
	}
	return false, fmt.Errorf("Connexion - nom d'utilisateur non trouvé")
}

func Deconnexion() {
	CurrentUser = User{}
}

func IsAuthenticated() bool {
	if CurrentUser.Name == "" {
		return false
	}
	return true
}

type Favoris struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type User struct {
	Name     string    `json:"name"`
	Password string    `json:"password"`
	Favoris  []Favoris `json:"favoris"`
}

func AddToFavorites(itemName string, Type string) bool {
	res, err := os.ReadFile("comptes.json")
	if err != nil {
		fmt.Println("Erreur de lecture du fichier:", err)
		return false
	}

	var liste []User
	err = json.Unmarshal(res, &liste)
	if err != nil {
		fmt.Println("Erreur de décodage du fichier:", err)
		return false
	}

	for i, user := range liste {
		if user.Name == CurrentUser.Name {
			liste[i].Favoris = append(liste[i].Favoris, Favoris{Name: itemName, Type: Type})
			CurrentUser.Favoris = liste[i].Favoris
			break
		}
	}

	dataEncode, err := json.Marshal(liste)
	if err != nil {
		fmt.Println("Erreur d'encodage des données:", err)
		return false
	}

	err = os.WriteFile("comptes.json", dataEncode, 0644)
	if err != nil {
		fmt.Println("Erreur d'écriture du fichier:", err)
		return false
	}

	return true
}

func GetFavorites() ([]Favoris, error) {
	res, err := os.ReadFile("comptes.json")
	if err != nil {
		return nil, fmt.Errorf("Erreur de lecture du fichier: %s", err.Error())
	}

	var liste []User
	err = json.Unmarshal(res, &liste)
	if err != nil {
		return nil, fmt.Errorf("Erreur de décodage du fichier: %s", err.Error())
	}

	for _, user := range liste {
		if user.Name == CurrentUser.Name {
			return user.Favoris, nil
		}
	}

	return nil, fmt.Errorf("Utilisateur non trouvé")
}

func SearchFavorite(itemName string) (Favoris, error) {
	res, err := os.ReadFile("comptes.json")
	if err != nil {
		return Favoris{}, fmt.Errorf("Erreur de lecture du fichier: %s", err.Error())
	}

	var liste []User
	err = json.Unmarshal(res, &liste)
	if err != nil {
		return Favoris{}, fmt.Errorf("Erreur de décodage du fichier: %s", err.Error())
	}

	for _, user := range liste {
		if user.Name == CurrentUser.Name {
			for _, favori := range user.Favoris {
				if favori.Name == itemName {
					return favori, nil
				}
			}
			break
		}
	}

	return Favoris{}, fmt.Errorf("Favori non trouvé")
}

func RemoveFavorite(itemName string) error {
	res, err := os.ReadFile("comptes.json")
	if err != nil {
		return fmt.Errorf("Erreur de lecture du fichier: %s", err.Error())
	}

	var liste []User
	err = json.Unmarshal(res, &liste)
	if err != nil {
		return fmt.Errorf("Erreur de décodage du fichier: %s", err.Error())
	}

	for i, user := range liste {
		if user.Name == CurrentUser.Name {
			for j, favori := range user.Favoris {
				if favori.Name == itemName {
					liste[i].Favoris = append(user.Favoris[:j], user.Favoris[j+1:]...)
					CurrentUser.Favoris = liste[i].Favoris
					break
				}
			}
			break
		}
	}

	dataEncode, err := json.Marshal(liste)
	if err != nil {
		return fmt.Errorf("Erreur d'encodage des données: %s", err.Error())
	}

	err = os.WriteFile("comptes.json", dataEncode, 0644)
	if err != nil {
		return fmt.Errorf("Erreur d'écriture du fichier: %s", err.Error())
	}

	return nil
}
