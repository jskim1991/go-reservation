package handlers

import (
	"database/sql"
	"errors"
	"reservation/internal/config"
	"reservation/internal/models"
	"time"
)

type testDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func (m *testDBRepo) AllUsers() bool {
	return true
}

func (m *testDBRepo) InsertReservation(res models.Reservation) (int, error) {
	if res.RoomID == 99 {
		return 0, errors.New("some error")
	}
	return 1, nil
}

func (m *testDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	if r.RoomID == 999 {
		return errors.New("some error")
	}
	return nil
}

func (m *testDBRepo) SearchAvailabilityRoomIdAndDates(roomId int, start, end time.Time) (bool, error) {
	return false, nil
}

func (m *testDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	var rooms []models.Room
	return rooms, nil
}

func (m *testDBRepo) GetRoomByID(id int) (models.Room, error) {
	var room models.Room
	if id > 2 {
		return room, errors.New("Some error")
	}
	return room, nil
}
