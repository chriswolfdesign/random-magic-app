package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	"github.com/rs/cors"
)

type Dragon struct {
	Name      string `json:"name"`
	Color     string `json:"color"`
	Manacost  string `json:"manaCost"`
	Power     int    `json:"power"`
	Toughness int    `json:"toughness"`
	Ability   string `json:"ability"`
}

func main() {
	PSQL_PASSWORD := os.Getenv("PSQL_PASSWORD")
	PSQL_USER := os.Getenv("PSQL_USER")
	PSQL_PORT := os.Getenv("PSQL_PORT")
	PSQL_DB := os.Getenv("PSQL_DB")
	PSQL_TABLE := os.Getenv("PSQL_TABLE")
	PSQL_HOST := os.Getenv("PSQL_HOST")

	db, err := createDB(PSQL_HOST, PSQL_PORT, PSQL_USER, PSQL_PASSWORD, PSQL_DB)
	defer db.Close()

	if err != nil {
		log.Println("ERROR ACCESSING DATABASE:", err)
		os.Exit(1)
	}

	log.Println("Contacted PSQL successfully")

	fmt.Println("Listening on port 8000")
	fmt.Println("Use Ctrl+C to quit")

	mux := http.NewServeMux()

	mux.HandleFunc("/greet", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to our Magic Database")
	})

	mux.HandleFunc("/dragons", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET")
		w.Header().Set("Allow-Access-Control-Headers", "text/plain; application/json")
		dragons, err := getDragons(db, PSQL_TABLE)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		json.NewEncoder(w).Encode(dragons)
	})

	mux.HandleFunc("/dragon", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET")
		w.Header().Set("Allow-Access-Control-Headers", "text/plain; application/json")
		color := r.URL.Query().Get("color")

		dragons, err := getDragon(db, PSQL_TABLE, color)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(dragons)

	})

	//handler := cors.Default().Handler(mux)
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://foo.com", "http://foo.com:8080"},
		AllowCredentials: true,
		// Enable Debugging for testing, consider disabling in production
		Debug: true,
	})

	handler := c.Handler(mux)

	err = http.ListenAndServe(":8000", handler)
	if err != nil {
		log.Fatal("ERROR STARTING SERVER:", err)
		os.Exit(1)
	}

}

func createDB(host string, port string, user string, password string, dbName string) (*sql.DB, error) {
	psqlconn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbName)
	return sql.Open("postgres", psqlconn)
}

func getDragons(db *sql.DB, table string) ([]Dragon, error) {
	rows, err := db.Query(fmt.Sprintf("select * from %s", table))
	if err != nil {
		return nil, err
	}

	dragons := []Dragon{}

	for rows.Next() {
		var dragon Dragon

		err := rows.Scan(&dragon.Name, &dragon.Color, &dragon.Manacost, &dragon.Power, &dragon.Toughness, &dragon.Ability)
		if err != nil {
			return nil, err
		}

		dragons = append(dragons, dragon)
	}

	return dragons, nil
}

func getDragon(db *sql.DB, table string, color string) ([]Dragon, error) {
	rows, err := db.Query(fmt.Sprintf("select * from %s where color='%s'", table, color))
	if err != nil {
		return nil, err
	}

	dragons := []Dragon{}

	for rows.Next() {
		var dragon Dragon

		err := rows.Scan(&dragon.Name, &dragon.Color, &dragon.Manacost, &dragon.Power, &dragon.Toughness, &dragon.Ability)
		if err != nil {
			return nil, err
		}

		dragons = append(dragons, dragon)
	}

	return dragons, nil
}
