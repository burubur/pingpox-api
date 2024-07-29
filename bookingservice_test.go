package main

import (
	"context"
	"testing"

	"github.com/burubur/pingpox-api/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestBookingService(t *testing.T) {
	ctrl := gomock.NewController(t)
	repositoryMock := mocks.NewMockRepositoryManager(ctrl)
	eventManagerMock := mocks.NewMockEventManager(ctrl)

	bookingService := NewBookingService(repositoryMock, eventManagerMock)
	_, err := bookingService.CreateBooking(context.TODO(), BookingRequest{})
	assert.Error(t, err)
}
