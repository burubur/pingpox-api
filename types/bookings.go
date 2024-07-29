package types

import (
	"time"

	"github.com/google/uuid"
)

type Bookings struct {
	ID                uuid.UUID // uuid
	UserID            uuid.UUID // uuid
	CoachID           uuid.UUID // uuid
	CourtID           string    // uuid
	DateTime          time.Time
	DurationInMinutes int
	LastStatus        TypeStatus // created, paid, confirmed, started, completed/cancelled/aborted
	CreatedAt         time.Time  // timestamp
	UpdatedAt         time.Time  // timestamp
}

type BookingHistories struct {
	ID        uuid.UUID
	BookingID uuid.UUID
	Timestamp time.Time
	Status    string
}

type TypeStatus string

func (s TypeStatus) ToString() string {
	return string(s)
}

const (
	StatusCreated   = TypeStatus("created")
	StatusPaid      = TypeStatus("paid")
	StatusConfirmed = TypeStatus("confirmed")
	StatusStarted   = TypeStatus("started")
	StatusCompleted = TypeStatus("completed")
	StatusCanceled  = TypeStatus("cancelled") // normal cancellation (trainee or trainer)
	StatusAborted   = TypeStatus("aborted")   // cancelled by system
)
