package handlers

import (
	"bookings/internals/config"
	"bookings/internals/driver"
	"bookings/internals/forms"
	"bookings/internals/helpers"
	"bookings/internals/models"
	"bookings/internals/render"
	"bookings/internals/repository"
	"bookings/internals/repository/dbrepo"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

var Repo *Repository

type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

type jsonResponse struct {
	OK        bool   `json:"ok"`
	Message   string `json:"message"`
	RoomID    string `json:"room_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgressRepo(db.SQL, a),
	}
}

func NewHandler(r *Repository) {
	Repo = r
}

func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)
	render.Template(w, r, "home.page.tmpl", &models.TemplateData{})
}

func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello User!!!"
	remoteIp := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIp
	render.Template(w, r, "about.page.tmpl", &models.TemplateData{})
}

func (m *Repository) Luxury(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "luxury.page.tmpl", &models.TemplateData{})
}

func (m *Repository) General(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "general.page.tmpl", &models.TemplateData{})
}

func (m *Repository) Reservations(w http.ResponseWriter, r *http.Request) {
	res, ok := m.App.Session.Get(r.Context(), "reservations").(models.Reservation)
	if !ok {
		helpers.ServerError(w, errors.New("Cannot get reservations details in Get Reservations"))
		return
	}
	sd := res.StartDate.Format("2006-01-02")
	ed := res.EndDate.Format("2006-01-02")

	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	data := make(map[string]interface{})
	data["reservations"] = res
	render.Template(w, r, "reservations.page.tmpl", &models.TemplateData{
		Form:      forms.New(nil),
		DataMap:   data,
		StringMap: stringMap,
	})
}

// Reservation submit, form validation and redirection
func (m *Repository) PostReservations(w http.ResponseWriter, r *http.Request) {
	res, ok := m.App.Session.Get(r.Context(), "reservations").(models.Reservation)
	if !ok {
		helpers.ServerError(w, errors.New("Cannot get reservations details in Post Reservations"))
		return
	}

	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	res.FirstName = r.Form.Get("FirstName")
	res.LastName = r.Form.Get("LastName")
	res.Email = r.Form.Get("Email")
	res.Phone = r.Form.Get("Phone")

	form := forms.New(r.PostForm)
	form.Required("FirstName", "LastName", "Email", "Phone")
	form.Minlength("FirstName", 3, r)
	form.Minlength("LastName", 3, r)
	form.IsNumeric("Phone", r)
	form.IsEmail("Email", r)

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservations"] = res
		render.Template(w, r, "reservations.page.tmpl", &models.TemplateData{
			Form:    form,
			DataMap: data,
		})
		return
	}

	reservation_id, err := m.DB.InsertReservation(res)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	room_restriction := models.RoomRestriction{
		StartDate:     res.StartDate,
		EndDate:       res.EndDate,
		RoomID:        res.RoomID,
		ReservationID: reservation_id,
		RestrictionID: 1,
	}

	err = m.DB.InsertRoomRestriction(room_restriction)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	m.App.Session.Put(r.Context(), "reservations", res)
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

// Reservation Summary shows the summary of successful reservation
func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservations, ok := m.App.Session.Get(r.Context(), "reservations").(models.Reservation)
	if !ok {
		log.Printf("Error on getting values from session")
		m.App.Session.Put(r.Context(), "Error", "Cannot get session values for resrevations")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}
	m.App.Session.Remove(r.Context(), "reservations")
	data := make(map[string]interface{})
	data["reservations"] = reservations

	sd := reservations.StartDate.Format("2006-01-02")
	ed := reservations.EndDate.Format("2006-01-02")

	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	render.Template(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		DataMap:   data,
		StringMap: stringMap,
	})
}

// Search Get method
func (m *Repository) Search(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "search.page.tmpl", &models.TemplateData{})
}

// Post Search through JSON
func (m *Repository) PostSearchJSON(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get("end")
	layout := "2006-01-02"

	startDate, err := time.Parse(layout, start)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	endDate, err := time.Parse(layout, end)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	roomid, err := strconv.Atoi(r.Form.Get("room_id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	available, err := m.DB.SearchAvailabilityByDatesByRoomId(startDate, endDate, roomid)
	resp := jsonResponse{
		OK:        available,
		Message:   "",
		RoomID:    strconv.Itoa(roomid),
		StartDate: start,
		EndDate:   end,
	}

	out, err := json.MarshalIndent(resp, "", "     ")
	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

// Post Search method logic, redirection
func (m *Repository) PostSearch(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get("end")
	var rooms []models.Room
	layout := "2006-01-02"

	startDate, err := time.Parse(layout, start)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	endDate, err := time.Parse(layout, end)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	rooms, err = m.DB.SearchAvailabilityForAllRooms(startDate, endDate)
	if err != nil {
		fmt.Println("Search Failed!!!")
		helpers.ServerError(w, err)
		return
	}
	if len(rooms) == 0 {
		fmt.Println("No Rooms block")
		m.App.Session.Put(r.Context(), "error", "No Room Available")
		http.Redirect(w, r, "/search", http.StatusSeeOther)
		return
	}
	reservations := models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
	}

	m.App.Session.Put(r.Context(), "reservations", reservations)
	m.App.Session.Put(r.Context(), "rooms", rooms)
	http.Redirect(w, r, "/search-summary", http.StatusSeeOther)
}

// Get Search summary shows the searched results
func (m *Repository) SearchSummary(w http.ResponseWriter, r *http.Request) {
	rooms, ok := m.App.Session.Get(r.Context(), "rooms").([]models.Room)
	if !ok {
		log.Printf("Error on getting values from rooms session")
		m.App.Session.Put(r.Context(), "Error", "Cannot get session values for rooms")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}
	m.App.Session.Remove(r.Context(), "rooms")
	data := make(map[string]interface{})
	data["rooms"] = rooms
	render.Template(w, r, "search-summary.page.tmpl", &models.TemplateData{
		DataMap: data,
	})
}

// Get Search summary shows the searched results
func (m *Repository) ChooseRoomById(w http.ResponseWriter, r *http.Request) {
	roomID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	res, ok := m.App.Session.Get(r.Context(), "reservations").(models.Reservation)
	if !ok {
		helpers.ServerError(w, err)
		return
	}
	room, err := m.DB.GetRoomDetailsByID(roomID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	res.RoomID = roomID
	res.Room = room
	m.App.Session.Put(r.Context(), "reservations", res)
	http.Redirect(w, r, "/reservations", http.StatusSeeOther)
}

func (m *Repository) BookRoom(w http.ResponseWriter, r *http.Request) {
	ID, err := strconv.Atoi(r.URL.Query().Get("id"))
	sd := r.URL.Query().Get("s")
	ed := r.URL.Query().Get("e")
	log.Println(ID, sd, ed)
	layout := "2006-01-02"

	startDate, err := time.Parse(layout, sd)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	endDate, err := time.Parse(layout, ed)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	res := models.Reservation{
		RoomID:    ID,
		StartDate: startDate,
		EndDate:   endDate,
	}

	m.App.Session.Put(r.Context(), "reservations", res)
	http.Redirect(w, r, "/reservations", http.StatusSeeOther)

}
