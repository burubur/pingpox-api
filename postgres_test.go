package main

import (
	"context"
	"testing"
	"time"

	"github.com/burubur/pingpox-api/types"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostgres(t *testing.T) {
	cfg := PostgresConfig{
		Host:         "localhost",
		Port:         5432,
		User:         "postgres",
		Password:     "secret",
		DatabaseName: "db_pingpox",
	}

	postgres, err := NewPostgres(cfg)
	_ = postgres

	require.NoError(t, err)

	timestamp := time.Now()
	bookingData := types.Bookings{
		ID:                uuid.New(),
		CoachID:           uuid.New(),
		UserID:            uuid.New(),
		DateTime:          time.Now().Add(+48 * time.Hour),
		CourtID:           "court-a",
		DurationInMinutes: 90,
		LastStatus:        types.StatusCreated,
		CreatedAt:         time.Now(),
		UpdatedAt:         timestamp,
	}

	result, err := postgres.StoreBookingCreationData(context.Background(), bookingData)
	assert.NoError(t, err)
	assert.NotNil(t, result)
}
