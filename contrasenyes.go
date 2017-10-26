package main

/**
  Programa que fa un atac a les contrasenyes d'un fitxer
  shadow de Linux (Per ara només funciona en SHA-512)
*/

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/kless/osutil/user/crypt/sha512_crypt"
)

/**
  Emmagatzema les contrasenyes de cada un dels usuaris
*/
type usuari struct {
	nom  string
	hash string
}

/**
Obtenir una llista amb els usuaris amb contrasenya
d'element 'usuari' a partir d'un fitxer en format
shadow

@returns llista dels usuaris
*/
func obtenirElsusuariDeShadow(nomFitxer string) []usuari {
	var llista []usuari
	file, err := os.Open(nomFitxer)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var linia = scanner.Text()
		var nom = strings.Split(linia, ":")

		if strings.HasPrefix(nom[1], "$6$") && strings.Count(nom[1], "$") == 3 {
			llista = append(llista, usuari{nom[0], nom[1]})
		}

	}
	return llista
}

/**
Comprova si la contrasenya d'un usuari està entre les del
fitxer de contrasenyes que es proporciona en el paràmetre
diccionari

@returns a contrasenya hi haurà la contrasenya trobada o err en cas
d'error o que no s'hagi trobat
*/
func comprovaUsuari(user usuari, diccionari string) (contrasenya string, err error) {

	file, err := os.Open(diccionari)
	if err != nil {
		return "", err
	}

	defer file.Close()
	c := sha512_crypt.New()
	scanner := bufio.NewScanner(file)
	// compta := 0
	for scanner.Scan() {
		var paraula = scanner.Text()
		// fmt.Println("Provant l'usuari " + user.nom + ":" + paraula)
		salt := user.hash
		hashResultat, err := c.Generate([]byte(paraula), []byte(salt))
		if err != nil {
			return "", err
		}
		// Mirem si hem trobat l'error
		if hashResultat == user.hash {
			return user.nom + ":" + paraula, nil
		}
		// compta = (compta + 1)
		// if compta%4000 == 0 {
		// 	fmt.Println(compta, " ... provant "+paraula)
		// }
	}
	return "", errors.New(user.nom + ": contrasenya no trobada")
}

/**
Carrega el fitxer de contrasenyes i el diccionari a fer servir i
retorna les contrasenyes que ha trobat.
*/
func main() {

	// 1 - Processar la línia de comandes
	fitxerDiccionari := flag.String("diccionari", "file.txt", "diccionari de paraules")
	fitxerShadow := flag.String("shadow", "shadow", "Fitxer de contrasenyes")

	flag.Parse()

	// 2 - Carregar la llista d'usuaris
	var usuaris []usuari
	usuaris = obtenirElsusuariDeShadow(*fitxerShadow)

	// 3 - Processar els usuaris en PARAL·LEL
	messages, errc := make(chan string), make(chan error)
	fmt.Println("Els usuaris són", len(usuaris))
	fmt.Println("------------------------------------")

	start := time.Now()

	for _, user := range usuaris {
		go func(user usuari) {
			paraula, err := comprovaUsuari(user, *fitxerDiccionari)
			if err != nil {
				errc <- err
				return
			}
			messages <- paraula
		}(user)
	}
	// 4 - Imprimir els resultats PARAL·LEL
	for i := 0; i < len(usuaris); i++ {
		select {
		case res := <-messages:
			fmt.Println(res)
		case err := <-errc:
			fmt.Println(err)
		}
	}

	// 3b - Processar els usuaris en ORDRE
	// -----------------------------------------
	// for _, user := range usuaris {
	// 	contrasenya, err := comprovaUsuari(user, "file.txt")
	// 	if err == nil {
	// 		fmt.Println("-----------------------------------------------")
	// 		fmt.Println(user.nom + " té de contrasenya " + contrasenya)
	// 	} else {
	// 		fmt.Println("--------------------------------------------")
	// 		fmt.Println(user.nom + " : " + err.Error())
	// 	}
	// }

	// 5. Càlcula el temps que ha tardat
	elapsed := time.Since(start)
	log.Printf("Ha tardat %s", elapsed)
}
