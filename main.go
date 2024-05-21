package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func main() {
	http.HandleFunc("/", home)
	http.HandleFunc("/export", export)

	fmt.Println("http://localhost:8090")

	if err := http.ListenAndServe(":8090", nil); err != nil {
		log.Fatal(err)
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 Page non trouvée", http.StatusNotFound)
		return
	}
	t := template.Must(template.ParseFiles("./html/ascii.html"))
	formatted := ""

	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "Erreur ParseForm() err=%v", err)
			return
		}

		text := r.FormValue("text") // extraction du texte dans le textarea
		for _, i := range text {
			if (i != 10 && i != 13 && i < 32) || i > 127 {
				formatted = "Les Caractères accentués ne sont pas permis"
				t.Execute(w, formatted)
				return
			}
		}

		format := r.FormValue("format") // extraction de la police sélectionnée
		file, err := ioutil.ReadFile(strings.TrimSpace(format) + ".txt")
		if err != nil {
			fmt.Println(err)
			return
		}

		fileStr := strings.Split(string(file), "\n")
		text = strings.ReplaceAll(string(text), "\r\n", "\\CR") // Remplacer les retours à la ligne(LF et CR) par un autre caractère
		Args := strings.Split(text, "\\CR")
		for _, val := range Args {
			for i := 0; i < 8; i++ {
				for j := 0; j < len(val); j++ {
					formatted += fileStr[int(val[j]-32)*9+1+i]
				}
				formatted += "\n"
			}
		}

	}

	t.Execute(w, formatted)
}

func export(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Disposition", "attachment; filename=ASCII-ART")
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", r.Header.Get("Content-Length"))
	http.ServeFile(w, r, "./output.txt")

}