# RecuperarContrasenyes
En la vostra empresa s’ha despatxat massivament a molts treballadors, entre ells l’administrador informàtic, i abans de fer-los marxar se’n van oblidar de demanar-los les contrasenyes. 

![El director té mal caràcter](https://raw.githubusercontent.com/utrescu/utrescu.github.io/master/images/kefe.png)

El problema és que ara no es poden recuperar els fitxers amb els que estaven treballant els treballadors abans que els despatxessin. I com que al director li fa vergonya trucar a la gent que va despatxar per demanar-los la contrasenya, es veu que van acabar malament, us ha demanat que li feu un programa que intenti descobrir les contrasenyes dels fitxers de contrasenyes dels despatxats.

El sistema on treballaven els treballadors despatxats és un Linux. En Linux les contrasenyes estan en el fitxer /etc/shadow i tenen un format estàndard:

![contrasenyes](https://raw.githubusercontent.com/utrescu/utrescu.github.io/master/images/shadow.png)

Tasca
------------------

Es tracta de fer un programa que vagi provant les contrasenyes que li proporcionarem contra els usuaris del sistema fins que els trobi o bé s'acabin les contrasenyes a provar.

El fitxer amb les contrasenyes a descobrir és aquest 

[Descarregar](https://drive.google.com/file/d/0BxakKCNfTojqbWplU1FfRldDVDA/view?usp=sharing "fitxer amb les contrasenyes")

> Compte que hi ha usuaris que no tenen contrasenya

El programa ha de permetre rebre un fitxer amb un diccionari de contrasenyes per provar: 

    bolet
    sabata
    aaaaaa
    sombrero
    patata
    barretina

Per exemple aquí hi ha exemples de diccionaris (on hi són les contrasenyes a trobar) 

[Descarregar](https://drive.google.com/file/d/0BxakKCNfTojqWkNJQ2luRldTM00/view?usp=sharing "diccionaris")

Programa
====================================
El programa s'executa des de línia de comandes. S'hi poden afegir paràmetres per indicar arxius que no són els que es fan servir per defecte:

~~~~bash
$ go build contrasenyes.go
[xavier@pilaf contrasenyes]$ ./contrasenyes -diccionari diccionaris/john.txt -shadow shadow
Els usuaris són 5
------------------------------------
manel:coffee
marcel:orange
pepet: contrasenya no trobada
manolo: contrasenya no trobada
frederic: contrasenya no trobada
2017/10/26 12:11:58 Ha tardat 16.20090154s
~~~~

