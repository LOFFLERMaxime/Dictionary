package main

import (
	"encoding/json"
	"estiam/dictionary"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

const dictionaryFilePath = "data/dictionary.json"

var logger *log.Logger

const validToken = "1234"

func main() {
	dict := dictionary.New(dictionaryFilePath)
	router := mux.NewRouter()

	// Ajout des middlewares
	router.Use(loggingMiddleware)
	router.Use(authenticationMiddleware)

	// Routes pour l'API REST
	router.HandleFunc("/ajouter", AjouterMot(dict)).Methods("POST")
	router.HandleFunc("/mots", GetMots(dict)).Methods("GET")
	router.HandleFunc("/definition/{mot}", GetDefinition(dict)).Methods("GET")
	router.HandleFunc("/remove/{mot}", RemoveMot(dict)).Methods("DELETE")

	// Démarrez le serveur
	fmt.Println("Serveur sur le port 8080...")

	err := http.ListenAndServe(":8080", handlers.LoggingHandler(os.Stdout, router))
	if err != nil {
		log.Printf("Erreur lors du démarrage du serveur : %v\n", err)
	}
}

func loadDataFromJSON(d *dictionary.Dictionary) {
	data, err := ioutil.ReadFile(dictionaryFilePath)
	if err != nil {
		fmt.Printf("Erreur lors de la lecture du fichier JSON : %v\n", err)
		return
	}

	err = json.Unmarshal(data, d.Entries())
	if err != nil {
		fmt.Printf("Erreur lors de la désérialisation des données JSON : %v\n", err)
		return
	}

	fmt.Println("Données chargées à partir du fichier JSON.")
}

// ROUTES PACKAGE
func AjouterMot(d *dictionary.Dictionary) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Sample implementation: Add a new entry to the dictionary
		word := r.FormValue("word")
		definition := r.FormValue("definition")

		d.Add(word, definition)

		// Log the action
		log.Printf("Ajout du mot '%s' avec comme définition '%s' au dictionnaire\n", word, definition)
		w.Write([]byte("Mot ajouté avec succès"))
	}
}

func GetDefinition(d *dictionary.Dictionary) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Sample implementation: Get the definition of a word
		word := mux.Vars(r)["mot"]

		entry, err := d.Get(word)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		// Log the action
		log.Printf("Affiche la définition du mot '%s': %s\n", word, entry)
		w.Write([]byte(entry.String()))
	}
}

func RemoveMot(d *dictionary.Dictionary) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Sample implementation: Remove a word from the dictionary
		word := mux.Vars(r)["mot"]

		d.Remove(word)

		// Log the action
		log.Printf("Supprime le mot '%s' du dictionnaire\n", word)
		w.Write([]byte("Mot supprimé avec succès"))
	}
}

func GetMots(d *dictionary.Dictionary) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Sample implementation: Get the list of words in the dictionary
		words, _ := d.List()

		// Log the action
		log.Println("Affiche les mots du dictionnaire")

		// Convert the list of words to a JSON response
		response, _ := json.Marshal(words)
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}
}

// /
// /
// / LOG HANDLER
func init() {
	// Open the log file for writing
	file, err := os.OpenFile("api.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Erreur lors de l'ouverture du fichier de log: ", err)
	}

	// Set the logger to write to the file
	logger = log.New(file, "API: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Requête reçue - Méthode: %s, Chemin: %s\n", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

// /
// /
// /Auth Package
func authenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token != validToken {
			http.Error(w, "Accès non autorisé", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
