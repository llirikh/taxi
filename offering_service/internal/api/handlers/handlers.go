package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"io"
	"net/http"
	"offering_service/internal/api/requests"
	"offering_service/internal/api/responses"
	"offering_service/internal/models"
	"offering_service/internal/service"
)

type OfferingHandler struct {
	Service *service.Offer_service
	Server  *http.Server
}

func NewHandler() *OfferingHandler {
	offerService := service.NewService()

	handler := OfferingHandler{Service: offerService}

	router := chi.NewRouter()
	router.Post("/offers", handler.CreateOffer)
	router.Get("/offers/{offerID}", handler.ParseOffer)

	handler.Server = &http.Server{
		Addr:    offerService.Config.Port,
		Handler: router,
	}

	return &handler
}

func (h *OfferingHandler) CreateOffer(w http.ResponseWriter, r *http.Request) {
	requestBodyJson, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	defer r.Body.Close()

	var request requests.CreateOfferRequest
	err = json.Unmarshal(requestBodyJson, &request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	offer := models.Offer{From: request.From, To: request.To, Client_id: request.Client_id}
	offer.Price = *h.Service.CountPrice(offer.From, offer.To)

	jwtOffer, err := h.Service.OfferToJwt(&offer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	response := responses.CreateOfferResponse{Offer_id: jwtOffer}
	resp, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *OfferingHandler) ParseOffer(w http.ResponseWriter, r *http.Request) {
	offerID := chi.URLParam(r, "offerID")

	parsedOffer, err := h.Service.JwtToOffer(offerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := responses.ParseOfferResponse{From: parsedOffer.From, To: parsedOffer.To, Price: parsedOffer.Price, Client_id: parsedOffer.Client_id}
	resp, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
