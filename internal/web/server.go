package web

import (
	"fmt"
	"log"
	"net/http"

	"github.com/troptropcontent/simon_bsc_checker/internal/db"
	"gorm.io/gorm"
)

type Server struct {
	db   *gorm.DB
	port int
}

func NewServer(db *gorm.DB, port int) *Server {
	return &Server{db: db, port: port}
}

func HandleRequest(database *gorm.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var logs []db.Log
		database.Order("created_at DESC").Find(&logs)
		for _, l := range logs {
			fmt.Fprintf(w, "[%s] %s - %s\n", l.CreatedAt.Format("2006-01-02 15:04:05"), l.Level, l.Message)
		}
	}
}

func (s *Server) Initialize() {
	http.HandleFunc("/", HandleRequest(s.db))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", s.port), nil))
}
