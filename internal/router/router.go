package router

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"

	"github.com/kolllaka/_img_uploader/internal/service"
)

type router struct {
	service service.Service
}

var tmpl *template.Template

func init() {
	tmpl = template.Must(template.ParseFiles("templates/index.html"))
}

func New(service service.Service) *router {
	return &router{
		service: service,
	}
}

func (s *router) Init() *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("/", s.homeHandler)
	// api routs
	router.HandleFunc("/api/upload", s.uploadHandler)
	router.HandleFunc("/api/link", s.uploadLinkHandler)
	router.HandleFunc("/api/delete", s.deleteLinkHandler)

	return router
}

func (s *router) homeHandler(w http.ResponseWriter, r *http.Request) {
	images, err := s.service.GetAllImages()
	if err != nil {
		http.Error(w, "Unable to get images", http.StatusInternalServerError)

		return
	}

	jsonImages, err := json.Marshal(images)
	if err != nil {
		log.Println(err.Error())
	}

	tmpl.Execute(w, string(jsonImages))
}
func (s *router) uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)

		return
	}

	file, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Failed to read file", http.StatusBadRequest)

		return
	}
	defer file.Close()

	image, err := s.service.SaveImage(file)
	if err != nil {
		http.Error(w, "Failed to save image", http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(image)
}
func (s *router) uploadLinkHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)

		return
	}

	type Data struct {
		Link string `json:"link"`
	}

	data := Data{}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)

		return
	}
	defer r.Body.Close()

	json.Unmarshal(body, &data)

	req, err := http.NewRequest("GET", data.Link, nil)
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)

		return
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to fetch image", http.StatusInternalServerError)

		return
	}
	defer res.Body.Close()

	image, err := s.service.SaveImage(res.Body)
	if err != nil {
		http.Error(w, "Failed to save image", http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(image)
}
func (s *router) deleteLinkHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)

		return
	}

	params := r.URL.Query()
	id := params.Get("id")

	if err := s.service.DeleteImage(id); err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to delete link", http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
