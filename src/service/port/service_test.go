package port

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/a-urth/abyr/pb/portpb"
	"github.com/a-urth/abyr/src/service/port/storage/postgres"
)

func TestPort(t *testing.T) {
	ctx := context.TODO()

	store, err := postgres.NewStorer()
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

	service := Service{store}

	p, err := service.GetPort(ctx, &portpb.PortID{Id: "AEAJM"})
	// since there are no ports we expect error and nil result
	if !(assert.EqualError(t, err, sql.ErrNoRows.Error()) && assert.Nil(t, p)) {
		return
	}

	port := portpb.Port{
		Id:      "AEAJM",
		Name:    "Ajman",
		City:    "Ajman",
		Country: "United Arab Emirates",
		Alias:   []string{},
		Regions: []string{},
		Coordinates: []string{
			"55.5136433",
			"25.4052165",
		},
		Province: "Ajman",
		Timezone: "Asia/Dubai",
		Unlocs:   []string{"AEAJM"},
		Code:     "52000",
	}
	resp, err := service.UpsertPort(ctx, &port)
	if !(assert.NoError(t, err) && assert.NotNil(t, resp)) {
		return
	}

	dbPort, err := service.GetPort(ctx, &portpb.PortID{Id: "AEAJM"})
	if !(assert.NoError(t, err) && assert.NotNil(t, dbPort)) {
		return
	}

	assert.Equal(t, port, *dbPort)
}
