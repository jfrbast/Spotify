package _func

import (
	"fmt"
	"os"
)

func EcrireCompte(username string, password string) {

	file, err := os.OpenFile("compte.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("Identifiant:%sMdp:%s\n", username, password))
	if err != nil {
		fmt.Println("Error writing to file:", err)
	}
}
