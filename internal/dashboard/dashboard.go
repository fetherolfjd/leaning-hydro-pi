package dashboard

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type DashboardDataItem struct {
	Temp            int
	SpecificGravity float32
}

type DashboardData struct {
	DataItems []DashboardDataItem
}

var dashboardTemplate *template.Template

func Configure(r *mux.Router) {
	log.Println("Configuring dashboard stuff")
	var err error

	dashboardTemplate, err = template.ParseFiles("web/templates/dashboard.html")
	// dashboardTemplate, err = template.New("DashboardTemplate").Parse(dbHtml)
	if err != nil {
		log.Fatalf("Failed to parse and create template 'DashboardTemplate': %v", err)
	}

	r.HandleFunc("/dashboard", dashboardHandler)
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Handle dashboard request")
	dbData := DashboardData{
		DataItems: []DashboardDataItem{
			DashboardDataItem{Temp: 69, SpecificGravity: 1.013},
			DashboardDataItem{Temp: 63, SpecificGravity: 1.010},
		},
	}
	err := dashboardTemplate.Execute(w, dbData)
	if err != nil {
		log.Fatalf("Error executing template: %v", err)
	}
}
