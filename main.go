package main

import (
	"log"
	"net/http"
	"os"

	"gopkg.in/boj/redistore.v1"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/yosssi/ace"
	"github.com/yosssi/ace-proxy"
)

var p = proxy.New(&ace.Options{BaseDir: "views"})

var store *redistore.RediStore

var router *mux.Router

const sessionName = "535510N"

var (
	HttpPort  string
	HttpsPort string
	CertFile  string
	KeyFile   string
)

func main() {
	var err error
	store, err = redistore.NewRediStore(10, "tcp", "127.0.0.1:6379", "", []byte("sessions"))
	if err != nil {
		log.Fatal(err)
	}
	defer store.Close()

	if HttpPort == "" {
		HttpPort = "8080"
	}

	router = mux.NewRouter()
	router.NotFoundHandler = http.HandlerFunc(NotFoundHandler)
	router.HandleFunc("/", HomeHandler).Methods("GET")
	router.HandleFunc("/homework", HomeworkHandler).Methods("GET")
	router.HandleFunc("/homework/aidan", HomeworkAidanHandler).Methods("GET")
	router.HandleFunc("/homework/assignments", HomeworkAssignmentsHandler).Methods("GET")
	router.HandleFunc("/homework/classes", HomeworkPUTClassesHandler).Methods("PUT")
	router.HandleFunc("/homework/classes", HomeworkGETClassesHandler).Methods("GET")
	router.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./static/"))))
	router.PathPrefix("/.well-known/").Handler(http.StripPrefix("/.well-known/", http.FileServer(http.Dir("./.well-known/"))))

	allHandler := handlers.CompressHandler(handlers.LoggingHandler(os.Stdout, router))

	if CertFile != "" && KeyFile != "" && HttpsPort != "" {
		go http.ListenAndServe(":"+HttpPort, allHandler)

		log.Fatal(http.ListenAndServeTLS(":"+HttpsPort, CertFile, KeyFile, allHandler))
	} else {
		log.Fatal(http.ListenAndServe(":"+HttpPort, allHandler))
	}

}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	ErrorHandler(w, r, "Not Found: The page you requested could not be found.", 404)
}

func ErrorHandler(w http.ResponseWriter, r *http.Request, errStr string, errNum int) {
	tpl, err := p.Load("base", "error", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Status":  errNum,
		"Message": errStr,
	}

	if err := tpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	tpl, err := p.Load("base", "index", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := tpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
