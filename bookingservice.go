package main

import (
	"context"
	"errors"
	"time"

	"github.com/burubur/pingpox-api/types"
	"github.com/google/uuid"
)

type RepositoryManager interface {
	StoreBookingCreationData(context.Context, types.Bookings) (uuid.UUID, error)
	FetchBookingData(context.Context, uuid.UUID) (types.Bookings, error)
	UpdateBookingStatus(context.Context, types.TypeStatus) error
}

type EventManager interface {
	PublishBookingCreationEvent(context.Context, types.Bookings) error
	PublishBookingAcceptanceEvent(context.Context, types.Bookings) error
}

type BookingService struct {
	repository   RepositoryManager
	eventManager EventManager
}

type BookingRequest struct {
	VoBooking
}

type VoBooking struct {
	UserID   uuid.UUID
	CoachID  uuid.UUID
	CourtID  string
	DateTime time.Time
	Duration time.Duration
}

type BookingResult struct {
	ID uuid.UUID
	VoBooking
	CreatedAt  time.Time
	LastStatus string
}

func NewBookingService() BookingService {
	return BookingService{}
}

// performed by user/trainee
func (b *BookingService) CreateBooking(ctx context.Context, request BookingRequest) (bookingResult *BookingResult, err error) {
	// 1. run validation
	// - double validation, validate the coach's availablity

	// 2. b.postgreSQL.insertBookingData()
	generatedBookingID := uuid.New()
	bookingData := types.Bookings{
		ID:         generatedBookingID,
		UserID:     request.UserID,
		CoachID:    request.CoachID,
		CourtID:    request.CourtID,
		DateTime:   request.DateTime,
		Duration:   request.Duration,
		LastStatus: types.StatusCreated,
	}

	bookingID, bookingErr := b.repository.StoreBookingCreationData(ctx, bookingData)
	if bookingErr != nil {
		return nil, bookingErr
	}

	bookingDetail, bookingDetailErr := b.repository.FetchBookingData(ctx, bookingID)
	if bookingDetailErr != nil {
		return &BookingResult{
			ID: bookingID,
		}, bookingDetailErr
	}

	bookingResult = &BookingResult{
		ID: bookingDetail.ID,
		VoBooking: VoBooking{
			UserID:   bookingDetail.UserID,
			CoachID:  bookingDetail.CoachID,
			CourtID:  bookingData.CourtID,
			DateTime: bookingData.DateTime,
		},
		CreatedAt:  bookingData.CreatedAt,
		LastStatus: bookingData.LastStatus.ToString(),
	}

	// 3. send `bookingCreated` to event driver handler/manager
	go func() {
		b.eventManager.PublishBookingCreationEvent(ctx, bookingDetail)
	}()

	return bookingResult, nil
}

// performed by coach/trainer
// confirming might accept or reject the booking appointment
func (b *BookingService) ConfirmBooking(ctx context.Context, bookingID uuid.UUID, coachID uuid.UUID) (lastStatus types.TypeStatus, err error) {
	if bookingID.String() == "" {
		return "", errors.New("invalid booking ID")
	}

	if coachID.String() == "" {
		return "", errors.New("invalid trainer ID")
	}

	// fetch from DB by the given bookingID
	bookingDetail, bookingDetailErr := b.repository.FetchBookingData(ctx, bookingID)
	if bookingDetailErr != nil {
		return "", errors.New("booking detail not found")
	}

	if bookingDetail.LastStatus != "paid" {
		return bookingDetail.LastStatus, errors.New("payment incomplete")
	}

	if !bookingDetail.DateTime.After(time.Now()) {
		return bookingDetail.LastStatus, errors.New("invalid request schedule")
	}

	if bookingDetail.DateTime.Sub(time.Now()) >= 3*time.Hour {
		return bookingDetail.LastStatus, errors.New("the requested scheduled is too tight")
	}

	updatingErr := b.repository.UpdateBookingStatus(ctx, types.StatusConfirmed)
	if updatingErr != nil {
		return bookingDetail.LastStatus, errors.New("failed to update booking status on DB")
	}

	go func() {
		b.eventManager.PublishBookingAcceptanceEvent(ctx, bookingDetail)
	}()

	return types.StatusConfirmed, nil
}
