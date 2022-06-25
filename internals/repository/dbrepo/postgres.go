package dbrepo

import (
	"bookings/internals/models"
	"context"
	"time"
)

func (m *postgressDBRepo) AllUsers() bool {
	return true
}

func (m *postgressDBRepo) InsertReservation(res models.Reservation) (int, error) {
	var id int
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `insert into reservation (first_name, last_name, email, phone,
		start_date, end_date, room_id, created_at, updated_at)
		values($1,$2,$3,$4,$5,$6,$7,$8,$9) returning id`

	err := m.DB.QueryRowContext(ctx, query,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		res.StartDate,
		res.EndDate,
		res.RoomID,
		time.Now(),
		time.Now(),
	).Scan(&id)

	if err != nil {
		return 0, err
	}
	return id, nil
}

func (m *postgressDBRepo) InsertRoomRestriction(res models.RoomRestriction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `insert into room_restriction (start_date, end_date, room_id, reservation_id,
		restriction_id,created_at, updated_at)
		values($1,$2,$3,$4,$5,$6,$7)`

	_, err := m.DB.ExecContext(ctx, query,
		res.StartDate,
		res.EndDate,
		res.RoomID,
		res.ReservationID,
		res.RestrictionID,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return err
	}
	return nil
}

func (m *postgressDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var rooms []models.Room
	query := `select r.id, r.room_name from room r
			  	where r.id not in (select rr.room_id from room_restriction rr
				where $1 < end_date and $2 > start_date);`

	rows, err := m.DB.QueryContext(ctx, query, start, end)
	if err != nil {
		return rooms, err
	}

	for rows.Next() {
		var room models.Room
		err := rows.Scan(&room.ID, &room.RoomName)
		if err != nil {
			return rooms, err
		}
		rooms = append(rooms, room)
	}
	return rooms, nil
}

func (m *postgressDBRepo) SearchAvailabilityByDatesByRoomId(start, end time.Time, roomid int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var count int
	query := `select count(rr.room_id) from room_restriction rr
				where $1 < end_date and $2 > start_date and rr.room_id = $3;`
	err := m.DB.QueryRowContext(ctx, query, start, end, roomid).Scan(&count)
	if err != nil {
		return false, err
	}

	if count == 0 {
		return true, nil
	}
	return false, nil
}

func (m *postgressDBRepo) GetRoomDetailsByID(room_id int) (models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var room models.Room
	query := `select * from room r where r.id = $1;`

	err := m.DB.QueryRowContext(ctx, query, room_id).Scan(&room.ID, &room.RoomName, &room.CreatedAt, &room.UpdatedAt)
	if err != nil {
		return room, err
	}
	return room, nil
}
