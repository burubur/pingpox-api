package main

import (
	"context"
	"database/sql"
	"time"

	"github.com/burubur/pingpox-api/types"
	"github.com/google/uuid"
)

type Postgres struct {
	dbConn *sql.DB
}

func NewPostgres() Postgres {
	return Postgres{}
}

func (p Postgres) StoreBookingCreationData(ctx context.Context, bookingData types.Bookings) (uuid.UUID, error) {
	// construct sql insertion query
	// construct sql parameters/arguments value
	// sql exec context
	// handle error
	// handle happy path/flow
	query := `INSERT INTO bookings (id, coach_id, user_id, date_time, court_id, duration, last_status, created_at, updated_at) 
	values ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	// always override createdAt and updatedAt on postgres implementation layer
	timeStamp := time.Now()
	bookingData.CreatedAt = timeStamp
	bookingData.UpdatedAt = timeStamp

	args := []any{
		bookingData.ID,
		bookingData.CoachID,
		bookingData.UserID,
		bookingData.DateTime,
		bookingData.CourtID,
		bookingData.Duration,
		bookingData.LastStatus,
		bookingData.CreatedAt,
		bookingData.UpdatedAt,
	}

	_, err := p.dbConn.ExecContext(ctx, query, args)
	if err != nil {
		return uuid.Nil, err
	}

	return bookingData.ID, nil
}

func (p Postgres) FetchBookingData(context.Context, uuid.UUID) (booking types.Bookings, err error) {
	return
}

func (p Postgres) UpdateBookingStatus(context.Context, types.TypeStatus) (err error) {
	return
}
