package types

import (
	"time"

	"github.com/google/uuid"
)

type Bookings struct {
	ID           uuid.UUID // uuid
	CoachID      uuid.UUID // uuid
	UserID       uuid.UUID // uuid
	DateTime     time.Time
	CourtID      string // uuid
	LatestStatus string // created, paid, confirmed, started, completed
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type BookingHistories struct {
	ID        uuid.UUID
	BookingID uuid.UUID
	Timestamp time.Time
	Status    string
}
