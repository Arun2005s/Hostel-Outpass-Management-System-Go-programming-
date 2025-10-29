package main
import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"
)

type Outpass struct {
	Name   string `json:"name"`
	RoomNo string `json:"roomNo"`
	Reason string `json:"reason"`
	Date   string `json:"date"`
}

type OutpassManager struct {
	Outpasses []Outpass
	FilePath  string
}

func (om *OutpassManager) Load() {
	file, err := os.ReadFile(om.FilePath)
	if err == nil {
		json.Unmarshal(file, &om.Outpasses)
	}
}

func (om *OutpassManager) Save() {
	data, err := json.MarshalIndent(om.Outpasses, "", "  ")
	if err != nil {
		fmt.Println("Failed to save:", err)
		return
	}
	err = os.WriteFile(om.FilePath, data, 0644)
	if err != nil {
		fmt.Println("Failed to write file:", err)
	}
}

func (om *OutpassManager) ApplyOutpass(name, roomNo, reason, date string) {
	om.Outpasses = append(om.Outpasses, Outpass{Name: name, RoomNo: roomNo, Reason: reason, Date: date})
	om.Save()
}

func (om *OutpassManager) SearchOutpass(name string) *Outpass {
	for _, o := range om.Outpasses {
		if strings.EqualFold(o.Name, name) {
			return &o
		}
	}
	return nil
}

func (om *OutpassManager) DeleteOutpass(name string) bool {
	for i, o := range om.Outpasses {
		if strings.EqualFold(o.Name, name) {
			om.Outpasses = append(om.Outpasses[:i], om.Outpasses[i+1:]...)
			om.Save()
			return true
		}
	}
	return false
}

var manager = OutpassManager{FilePath: "outpasses.json"}

func main() {
	manager.Load()

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", handleHome)
	http.HandleFunc("/apply", handleApply)
	http.HandleFunc("/view", handleView)
	http.HandleFunc("/search", handleSearch)
	http.HandleFunc("/delete", handleDelete)

	fmt.Println("Server starting on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, nil)
}

func handleApply(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		name := r.FormValue("name")
		roomNo := r.FormValue("roomNo")
		reason := r.FormValue("reason")
		date := r.FormValue("date")

		manager.ApplyOutpass(name, roomNo, reason, date)
		http.Redirect(w, r, "/view", http.StatusSeeOther)
		return
	}
	tmpl := template.Must(template.ParseFiles("templates/apply.html"))
	tmpl.Execute(w, nil)
}

func handleView(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/view.html"))
	tmpl.Execute(w, manager.Outpasses)
}

func handleSearch(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		name := r.FormValue("name")
		outpass := manager.SearchOutpass(name)
		tmpl := template.Must(template.ParseFiles("templates/search.html"))
		tmpl.Execute(w, outpass)
		return
	}
	tmpl := template.Must(template.ParseFiles("templates/search.html"))
	tmpl.Execute(w, nil)
}

func handleDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		name := r.FormValue("name")
		if manager.DeleteOutpass(name) {
			http.Redirect(w, r, "/view", http.StatusSeeOther)
			return
		}
	}
	tmpl := template.Must(template.ParseFiles("templates/delete.html"))
	tmpl.Execute(w, nil)
}

