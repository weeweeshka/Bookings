package storage

import (
	"bookings/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"time"
)

func (pos *Postgres) CreateVisitor(hotelId int, hotelRoom int, firstName string, lastName string, age int) error {
	const op = "storage.postgres.CreateVisitor"
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := pos.conn.Exec(ctx, "create_visitor", hotelId, hotelRoom, firstName, lastName, age)
	if err != nil {
		return fmt.Errorf("%s: exec failed: %w", op, err)
	}
	return nil
}

func (pos *Postgres) GetAllVisitors() (string, error) {
	const op = "storage.postgres.CreateVisitor"
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := pos.conn.Query(ctx, "get_all_visitors")
	if err != nil {
		return "", fmt.Errorf("%s: query failed: %w", op, err)
	}

	var visitors []models.Visitor
	for rows.Next() {
		var vis models.Visitor

		rows.Scan(&vis.Id, &vis.HotelId, &vis.HotelRoom, &vis.FirstName, &vis.LastName, &vis.Age)
		visitors = append(visitors, vis)

	}

	jsonVis, _ := json.Marshal(visitors)
	return string(jsonVis), nil

}

func (pos *Postgres) GetVisitor(id int) (string, error) {
	const op = "storage.postgres.GetVisitor"
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var vis models.Visitor

	err := pos.conn.QueryRow(ctx, "get_visitor", id).Scan(&vis.Id, &vis.HotelId, &vis.HotelRoom, &vis.FirstName, &vis.LastName, &vis.Age)
	if err != nil {
		return "", fmt.Errorf("%s: query failed: %w", op, err)
	}

	jsonVis, _ := json.Marshal(vis)

	return string(jsonVis), nil

}

func (pos *Postgres) DeleteVisitor(id int) error {
	const op = "storage.postgres.DeleteVisitor"
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := pos.conn.Exec(ctx, "delete_visitor", id)
	if err != nil {
		return fmt.Errorf("%s: exec failed: %w", op, err)
	}

	return nil
}

func (pos *Postgres) UpdateVisitor(id int, hotelId int, hotelRoom int, firstName string, lastName string, age int) (string, error) {
	const op = "storage.postgres.UpdateVisitor"
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := pos.conn.Exec(ctx, "update_visitor", hotelId, hotelRoom, firstName, lastName, age, id)
	if err != nil {
		return "", fmt.Errorf("%s: exec failed: %w", op, err)
	}

	getUV, _ := pos.GetVisitor(id)

	return getUV, nil
}
