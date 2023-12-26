package triprepo

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"trip/internal/models"

	sq "github.com/Masterminds/squirrel"
)

type TripSearchCriteria struct {
	ID      string
	UserID  string
	OfferID string
	Status  string
}

type TripRepository struct {
	db *sqlx.DB
}

func NewTripRepository(db *sqlx.DB) *TripRepository {
	return &TripRepository{db: db}
}

func (r *TripRepository) BeginTransaction(ctx context.Context, opts *sql.TxOptions, f func(ctx context.Context) error) error {
	tx, err := r.db.BeginTx(ctx, opts)
	if err != nil {
		return err
	}

	defer func() {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil && !errors.Is(rollbackErr, sql.ErrTxDone) {
			err = rollbackErr
		}
	}()

	fErr := f(ctx)
	if fErr != nil {
		_ = tx.Rollback()
		return fErr
	}

	return tx.Commit()
}

func (r *TripRepository) SearchTrips(ctx context.Context, criteria *TripSearchCriteria) ([]models.Trip, error) {
	query := sq.Select("*").
		From("trips").
		PlaceholderFormat(sq.Dollar)

	if criteria.ID != "" {
		query = query.Where(sq.Eq{"author": criteria.ID})
	}

	if criteria.UserID != "" {
		query = query.Where(sq.Eq{"user_id": criteria.UserID})
	}
	if criteria.OfferID != "" {
		query = query.Where(sq.Eq{"offer_id": criteria.OfferID})
	}
	if criteria.Status != "" {
		query = query.Where(sq.Eq{"status": criteria.Status})
	}

	sql, args, err := query.ToSql()
	rows, err := r.db.QueryxContext(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	trips := make([]models.Trip, 0)

	for rows.Next() {
		trip := models.Trip{}

		if err = rows.StructScan(&trip); err != nil {
			return nil, err
		}

		trips = append(trips, trip)
	}

	return trips, nil
}

func (r *TripRepository) AddTrip(ctx context.Context, trip *models.Trip) error {
	sql, args, err := sq.
		Insert("trips").Columns("id", "user_id", "offer_id", "status").
		Values(trip.Id, trip.UserId, trip.OfferId, trip.Status).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	var id int
	row := r.db.QueryRowContext(ctx, sql, args...)
	if err = row.Scan(&id); err != nil {
		return err
	}

	return nil
}
