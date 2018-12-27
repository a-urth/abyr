package postgres

import (
	"context"
	"fmt"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	"github.com/a-urth/abyr/pb/portpb"
	"github.com/a-urth/abyr/src/service/port/storage"
)

// Storer implements storer interface using postgres as a storage
type Storer struct {
	db *sqlx.DB
}

// NewStorer creates and returns postgres storer implementation object
// TODO: proper configuration should be used here
func NewStorer(host, port string) (storage.Storer, error) {
	connString := fmt.Sprintf(
		"user=postgres dbname=postgres sslmode=disable host=%s port=%s",
		host, port,
	)

	db, err := sqlx.Connect("postgres", connString)
	if err != nil {
		return nil, err
	}

	storer := Storer{db}
	return &storer, nil
}

// GetPort return port information from database by given id
func (s *Storer) GetPort(
	ctx context.Context, portID string,
) (*portpb.Port, error) {
	port := portpb.Port{}
	coordinates := []string{}
	err := s.db.QueryRowContext(
		ctx,
		`SELECT id, name, city, country, alias, regions, coordinates, province, timezone, unlocs, code FROM port WHERE id = $1`,
		portID,
	).Scan(
		&port.Id, &port.Name, &port.City, &port.Country, pq.Array(&port.Alias),
		pq.Array(&port.Regions), pq.Array(&coordinates), &port.Province,
		&port.Timezone, pq.Array(&port.Unlocs), &port.Code,
	)
	if err != nil {
		return nil, err
	}

	port.Coordinates = make([]float32, len(coordinates))
	for i, coord := range coordinates {
		v, err := strconv.ParseFloat(coord, 32)
		if err != nil {
			return nil, err
		}

		port.Coordinates[i] = float32(v)
	}

	return &port, nil
}

// UpsertPort upserts port entity in database from given information
func (s *Storer) UpsertPort(
	ctx context.Context, port *portpb.Port,
) error {
	// unfortunatelly pq doesn't know how to work with pg array of floats
	coordinates := make([]string, len(port.Coordinates))
	for i, coord := range port.Coordinates {
		coordinates[i] = fmt.Sprintf("%f", coord)
	}

	_, err := s.db.NamedExecContext(
		ctx,
		`INSERT INTO port (id, name, city, country, alias, regions, coordinates, province, timezone, unlocs, code) `+
			`VALUES (:id, :name, :city, :country, :alias, :regions, :coordinates, :province, :timezone, :unlocs, :code) `+
			`ON CONFLICT (id) `+
			`DO UPDATE SET `+
			`name = :name, `+
			`city = :city, `+
			`country = :country, `+
			`alias = :alias, `+
			`regions = :regions, `+
			`coordinates = :coordinates, `+
			`province = :province, `+
			`timezone = :timezone, `+
			`unlocs = :unlocs, `+
			`code = :code, `+
			`updated_at = now();`,
		map[string]interface{}{
			"id":          port.Id,
			"name":        port.Name,
			"city":        port.City,
			"country":     port.Country,
			"alias":       pq.Array(port.Alias),
			"regions":     pq.Array(port.Regions),
			"coordinates": pq.Array(coordinates),
			"province":    port.Province,
			"timezone":    port.Timezone,
			"unlocs":      pq.Array(port.Unlocs),
			"code":        port.Code,
		},
	)

	return err
}

// DeletePort hard deletes port from database
func (s *Storer) DeletePort(ctx context.Context, portID string) error {
	_, err := s.db.ExecContext(
		ctx, "DELETE FROM port WHERE id = $1", portID,
	)
	return err
}

// Close closes database connection
func (s *Storer) Close() error {
	return s.db.Close()
}
