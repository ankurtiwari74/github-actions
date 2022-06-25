package repository

import (
	"bookings/internals/models"
	"time"
)

type DatabaseRepo interface {
	AllUsers() bool
	InsertReservation(res models.Reservation) (int, error)
	InsertRoomRestriction(res models.RoomRestriction) error
	SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error)
	GetRoomDetailsByID(room_id int) (models.Room, error)
	SearchAvailabilityByDatesByRoomId(start, end time.Time, roomid int) (bool, error)
}
