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

type User struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Favoris  []Fav  `json:"favoris"`
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

	var newUser User = User{Name: name, Password: HashPassword(password), Favoris: []Fav{}}

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
