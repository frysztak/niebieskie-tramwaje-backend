package News

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"net/http"
	"net/url"
	"strconv"
)

type Handler func(w http.ResponseWriter, r *http.Request)

func RecentNewsHandler(db *sqlx.DB) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		news := getNews(db, 1, 0)
		data := news[0]

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-cache")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(data)
	}
}

func NewsHandler(db *sqlx.DB) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		pageNumStr, err := url.QueryUnescape(vars["pageNum"])
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		page, err := strconv.ParseInt(pageNumStr, 10, 32)
		data := getNews(db, itemsPerPage, int(page))

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-cache")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(data)
	}
}
