package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"sync"
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

var (
	users       []User
	usersMutex  sync.RWMutex
	accountFile = "comptes.json"
)

func LoadUsers() error {
	usersMutex.Lock()
	defer usersMutex.Unlock()

	file, err := ioutil.ReadFile(accountFile)
	if err != nil {
		return err
	}

	return json.Unmarshal(file, &users)
}

func SaveUsers() error {
	usersMutex.Lock()
	defer usersMutex.Unlock()

	file, err := json.MarshalIndent(users, "", " ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(accountFile, file, 0644)
}

func AddFavoris(username, name, itemType string) error {
	usersMutex.Lock()
	defer usersMutex.Unlock()

	for i, user := range users {
		if user.Name == username {
			// Vérifier si le favori existe déjà
			for _, fav := range user.Favoris {
				if fav.Name == name && fav.Type == itemType {
					return fmt.Errorf("favori déjà existant")
				}
			}

			users[i].Favoris = append(users[i].Favoris, Favoris{
				Name: name,
				Type: itemType,
			})
			return SaveUsers()
		}
	}

	return fmt.Errorf("utilisateur non trouvé")
}

func RemoveFavoris(username, name, itemType string) error {
	usersMutex.Lock()
	defer usersMutex.Unlock()

	for i, user := range users {
		if user.Name == username {
			for j, fav := range user.Favoris {
				if fav.Name == name && fav.Type == itemType {
					users[i].Favoris = append(users[i].Favoris[:j], users[i].Favoris[j+1:]...)
					return SaveUsers()
				}
			}
			return fmt.Errorf("favori non trouvé")
		}
	}

	return fmt.Errorf("utilisateur non trouvé")
}

func GetUserFavoris(username string) ([]Favoris, error) {
	usersMutex.RLock()
	defer usersMutex.RUnlock()

	for _, user := range users {
		if user.Name == username {
			return user.Favoris, nil
		}
	}

	return nil, fmt.Errorf("utilisateur non trouvé")
}

func IsFavoris(username, name, itemType string) bool {
	favoris, err := GetUserFavoris(username)
	if err != nil {
		return false
	}

	for _, fav := range favoris {
		if fav.Name == name && fav.Type == itemType {
			return true
		}
	}

	return false
}

type Pagination struct {
	TotalItems   int
	ItemsPerPage int
	CurrentPage  int
	TotalPages   int
}

func NewPagination(totalItems, itemsPerPage, currentPage int) Pagination {
	if itemsPerPage <= 0 {
		itemsPerPage = 10 // Default
	}
	if currentPage <= 0 {
		currentPage = 1
	}

	totalPages := int(math.Ceil(float64(totalItems) / float64(itemsPerPage)))

	return Pagination{
		TotalItems:   totalItems,
		ItemsPerPage: itemsPerPage,
		CurrentPage:  currentPage,
		TotalPages:   totalPages,
	}
}

func (p *Pagination) GetOffset() int {
	return (p.CurrentPage - 1) * p.ItemsPerPage
}

func (p *Pagination) GetPaginatedItems(items interface{}) interface{} {
	// Cette méthode devrait être implémentée de manière générique selon le type d'items
	return items
}
