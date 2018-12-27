package postgres

import (
	"context"

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
func NewStorer() (storage.Storer, error) {
	db, err := sqlx.Open(
		"postgres",
		"user=postgres dbname=postgres sslmode=disable",
	)
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
	err := s.db.QueryRowContext(
		ctx,
		`SELECT id, name, city, country, alias, regions, coordinates, province, timezone, unlocs, code FROM port WHERE id = $1`,
		portID,
	).Scan(
		&port.Id, &port.Name, &port.City, &port.Country, pq.Array(&port.Alias),
		pq.Array(&port.Regions), pq.Array(&port.Coordinates), &port.Province,
		&port.Timezone, pq.Array(&port.Unlocs), &port.Code,
	)
	if err != nil {
		return nil, err
	}

	return &port, nil
}

// UpsertPort upserts port entity in database from given information
func (s *Storer) UpsertPort(
	ctx context.Context, port *portpb.Port,
) error {
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
			"coordinates": pq.Array(port.Coordinates),
			"province":    port.Province,
			"timezone":    port.Timezone,
			"unlocs":      pq.Array(port.Unlocs),
			"code":        port.Code,
		},
	)

	return err
}

// Close closes database connection
func (s *Storer) Close() error {
	return s.db.Close()
}
