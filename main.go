package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Weather struct {
	CurrentCondition []struct {
		TempC string `json:"temp_C"`
	} `json:"current_condition"`
}

func main() {
	fmt.Printf("Data uruchomienia: %s\n", time.Now().Format(time.RFC1123))
	fmt.Println("Autor programu: Wojciech Makówka")
	fmt.Println("Aplikacja nasłuchuje na porcie: 8080")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		city := r.URL.Query().Get("city")
		
		if city == "" {
			html := `<html><body>
				<h2>Sprawdz pogode</h2>
				<form action="/" method="GET">
					Miasto: <input type="text" name="city" placeholder="np. Warsaw">
					<input type="submit" value="Sprawdz">
				</form>
			</body></html>`
			fmt.Fprint(w, html)
			return
		}

		resp, err := http.Get(fmt.Sprintf("https://wttr.in/%s?format=j1", city))
		if err != nil || resp.StatusCode != 200 {
			fmt.Fprintf(w, "Błąd pobierania pogody dla: %s", city)
			return
		}
		defer resp.Body.Close()

		var weather Weather
		json.NewDecoder(resp.Body).Decode(&weather)

		if len(weather.CurrentCondition) > 0 {
			temp := weather.CurrentCondition[0].TempC
			fmt.Fprintf(w, "<html><body><h2>Aktualna pogoda w %s</h2><h3>Temperatura: %s °C</h3><a href='/'>Powrót</a></body></html>", city, temp)
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}