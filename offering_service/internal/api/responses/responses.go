package responses

import "offering_service/internal/models"

type CreateOfferResponse struct {
	Offer_id string `json:"offer_Id"`
}

type ParseOfferResponse struct {
	From      models.Position `json:"from"`
	To        models.Position `json:"to"`
	Client_id string          `json:"client_Id"`
	Price     models.Price    `json:"price"`
}
