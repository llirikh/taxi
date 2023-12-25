package service

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"go.uber.org/zap"
	"log"
	"math"
	"offering_service/internal/config"
	"offering_service/internal/models"
	"time"
)

const (
	EarthRad   = 6371
	MoneyPerKm = 15
	MinPrice   = 200
	ExpireTime = 12 * time.Hour
)

type Offer_service struct {
	Config *models.Config
	Logger *zap.Logger
}

func NewService() *Offer_service {

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Logger init error. %v", err)
		return nil
	}

	cfg, err := config.InitConfig()
	if err != nil {
		logger.Warn("Error initialization configuration")
	}

	return &Offer_service{Config: cfg, Logger: logger}
}

func degreeToRadians(degree float64) float64 {
	return degree * (math.Pi / 180)
}

func getDistanceFromLatLon(from models.Position, to models.Position) float64 {
	dLat := degreeToRadians(from.Lat - to.Lat)
	dLng := degreeToRadians(from.Lng - to.Lng)
	alpha := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(degreeToRadians(from.Lat))*math.Cos(degreeToRadians(to.Lat))*
			math.Sin(dLng/2)*math.Sin(dLng/2)
	distance := EarthRad * 2 * math.Atan2(math.Sqrt(alpha), math.Sqrt(1-alpha))
	return math.Round(math.Abs(distance))
}

func (o *Offer_service) CountPrice(from models.Position, to models.Position) *models.Price {
	price := &models.Price{Amount: math.Max(getDistanceFromLatLon(from, to)*MoneyPerKm, MinPrice), Currency: "RUB"}
	return price
}

func (o *Offer_service) OfferToJwt(offer *models.Offer) (string, error) {
	bytes, err := json.Marshal(offer)
	if err != nil {
		return "", err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"offer": string(bytes),
		"exp":   time.Now().Add(ExpireTime).Unix(),
	})

	return jwtToken.SignedString([]byte(o.Config.PrivateKey))
}

func (o *Offer_service) JwtToOffer(jwtToken string) (*models.Offer, error) {
	parsedToken, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(o.Config.PrivateKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !parsedToken.Valid {
		return nil, fmt.Errorf("Invalid jwt-token")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("Broken jwt-token")
	}

	offer, ok := claims["offer"].(string)
	if !ok {
		return nil, fmt.Errorf("No offer in jwt-token")
	}

	var jsonOffer models.Offer
	err = json.Unmarshal([]byte(offer), &jsonOffer)
	if err != nil {
		return nil, err
	}

	return &jsonOffer, nil
}
