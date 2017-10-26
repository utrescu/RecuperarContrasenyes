package main

/**
  Programa que fa un atac a les contrasenyes d'un fitxer
  shadow de Linux (Per ara només funciona en SHA-512)
*/

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	u "github.com/utrescu/gontrasenya"
)

/**
Carrega el fitxer de contrasenyes i el diccionari a fer servir i
retorna les contrasenyes que ha trobat.
*/
func main() {

	// 1 - Processar la línia de comandes
	fitxerDiccionari := flag.String("diccionari", "file.txt", "diccionari de paraules")
	fitxerShadow := flag.String("shadow", "shadow", "Fitxer de contrasenyes")

	flag.Parse()

	if _, err := os.Stat(*fitxerDiccionari); os.IsNotExist(err) {
		fmt.Println("El fitxer", *fitxerDiccionari, " no existeix")
		return
	}
	if _, err := os.Stat(*fitxerShadow); os.IsNotExist(err) {
		fmt.Println("El fitxer", *fitxerShadow, " no existeix")
		return
	}

	// 2 - Carregar la llista d'usuaris
	var usuaris []u.Usuari
	usuaris = u.ObtenirElsusuariDeShadow(*fitxerShadow)

	// 3 - Processar els usuaris en PARAL·LEL
	messages, errc := make(chan string), make(chan error)
	fmt.Println("Els usuaris són", len(usuaris))
	fmt.Println("------------------------------------")

	start := time.Now()

	for _, user := range usuaris {
		go func(user u.Usuari) {
			paraula, err := u.ComprovaUsuari(user, *fitxerDiccionari)
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
