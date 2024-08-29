package repository

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"main/internal/v1/core_api/domain/ports"
)

// postgresqlRepo Struct
type postgresqlRepo struct {
	ctx context.Context
	db  *pgxpool.Pool
}

// NewPostgresqlRepository Auth Domain postgresql repository constructor
func NewPostgresqlRepository(ctx context.Context, db *pgxpool.Pool) ports.IPostgresqlRepository {
	return &postgresqlRepo{
		ctx: ctx,
		db:  db,
	}
}

// Example Postgresql repository Function
/*
func (r *postgresqlRepo) GetVehicleTypes() (record int64, data json.RawMessage) {
	var query string
	var errDb error
	record = 1

	query = `SELECT JSON_AGG(GVC) FROM (SELECT type_id,def_val,tr_val,en_val FROM vehicles.veh_types ORDER BY def_val ASC) GVC `
	errDb = r.db.QueryRow(r.ctx, query).Scan(&data)
	if errDb != nil && cm.CheckStringIfContains(errDb.Error(), "no rows in result set") == false {
		record = -1
		fmt.Printf("GetVehicleTypes DB Err: ", errDb.Error())

	} else if errDb != nil && cm.CheckStringIfContains(errDb.Error(), "no rows in result set") == true {
		record = 0
	}

	return
}
*/
