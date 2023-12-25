package handlers

import (
	"client_service/internal/api/requests"
	"client_service/internal/models"
	"client_service/internal/mongodb"
	"encoding/json"
	"github.com/go-chi/chi"
	"io"
	"net/http"
)

type ClientHandler struct {
	Server   *http.Server
	Config   *models.Config
	Database *mongodb.Database
}

func NewHandler(db *mongodb.Database, cfg *models.Config) *ClientHandler {
	handler := ClientHandler{Config: cfg, Database: db}

	router := chi.NewRouter()
	router.Post("/trips", handler.createTrip)
	router.Get("/trips", handler.listTrips)
	router.Get("/trips/{trip_id}", handler.getTrip)
	router.Post("/trip/{trip_id}/cancel", handler.cancelTrip)

	handler.Server = &http.Server{
		Addr:    handler.Config.Port,
		Handler: router,
	}

	return &handler
}

func (h *ClientHandler) getTrip(w http.ResponseWriter, r *http.Request) {
	tripID := chi.URLParam(r, "trip_id")
	trip, err := h.Database.GetTripByID(tripID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	resp, err := json.Marshal(trip)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = w.Write(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *ClientHandler) listTrips(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("user_id")
	trips, err := h.Database.GetTripsByUserID(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	resp, err := json.Marshal(trips)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = w.Write(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *ClientHandler) cancelTrip(w http.ResponseWriter, r *http.Request) {
	tripID := chi.URLParam(r, "trip_id")
	if err := h.Database.CancelTripByID(tripID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *ClientHandler) createTrip(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("user_id")
	requestBodyJson, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	defer r.Body.Close()

	var request requests.RequestCreateTrip
	err = json.Unmarshal(requestBodyJson, &request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := http.Get("http://127.0.0.1:8080/offers/" + request.Offer_id)
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	var offer models.Offer
	err = json.Unmarshal(bytes, &offer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	defer resp.Body.Close()

	trip := mongodb.Trip{Offer_id: request.Offer_id, Client_id: userID,
		From:   mongodb.Position{Lat: offer.From.Lat, Lng: offer.From.Lng},
		To:     mongodb.Position{Lat: offer.To.Lat, Lng: offer.To.Lng},
		Price:  mongodb.Price{Amount: offer.Price.Amount, Currency: offer.Price.Currency},
		Status: "DRIVER_SEARCH"}

	err = h.Database.CreateTrip(&trip)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusOK)
}
