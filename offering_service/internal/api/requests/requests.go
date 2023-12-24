package requests

import "offering_service/internal/models"

type CreateOfferRequest struct {
	From      models.Position `json:"from"`
	To        models.Position `json:"to"`
	Client_id string          `json:"client_Id"`
}

type ParseOfferRequest struct {
	Offer_id string `json:"offer_Id"`
}
