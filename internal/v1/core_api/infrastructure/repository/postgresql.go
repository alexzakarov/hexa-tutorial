package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"main/internal/v1/core_api/domain/entities"
)

// postgresqlRepo Struct
type postgresqlRepo struct {
	ctx context.Context
	db  *pgxpool.Pool
}

// NewPostgresqlRepository Auth Domain postgresql repository constructor
func NewPostgresqlRepository(ctx context.Context, db *pgxpool.Pool) *postgresqlRepo {
	return &postgresqlRepo{
		ctx: ctx,
		db:  db,
	}
}

// create
func (r *postgresqlRepo) CreateUser(req_dat entities.UserReqDto) error {
	query := `INSERT INTO users.users (role_id, email, phone, address, user_name,  password, has_gdpr) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.Exec(r.ctx, query, req_dat.RoleId, req_dat.Email, req_dat.Phone, req_dat.Address, req_dat.UserName, req_dat.Password, req_dat.HasGdpr)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

// read
func (r *postgresqlRepo) GetUserById(user_id int) (err error, data entities.UserResDto) {
	query := `SELECT  user_name, email, password FROM users.users WHERE id = $1`
	err = r.db.QueryRow(r.ctx, query, user_id).Scan(&data.UserName, &data.Email, &data.Password)
	if err != nil {
		fmt.Println(err.Error())
	}
	return
}

// update
func (r *postgresqlRepo) UpdateUser(user_id int, req_dat entities.UserReqDto) error {
	query := `UPDATE users.users SET user_name = $1, email = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $3`
	_, err := r.db.Exec(r.ctx, query, req_dat.UserName, req_dat.Email, user_id)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

// delete
func (r *postgresqlRepo) DeleteUser(id int64) error {
	query := `DELETE FROM users.users WHERE id = $1`
	_, err := r.db.Exec(r.ctx, query, id)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

// Example Postgresql repository Function

//func (r *postgresqlRepo) GetVehicleTypes() (record int64, data json.RawMessage) {
//	var query string
//	var errDb error
//	record = 1
//
//	query = `SELECT JSON_AGG(GVC) FROM (SELECT type_id,def_val,tr_val,en_val FROM vehicles.veh_types ORDER BY def_val ASC) GVC `
//	errDb = r.db.QueryRow(r.ctx, query).Scan(&data)
//	if errDb != nil && cm.CheckStringIfContains(errDb.Error(), "no rows in result set") == false {
//		record = -1
//		fmt.Printf("GetVehicleTypes DB Err: ", errDb.Error())
//
//	} else if errDb != nil && cm.CheckStringIfContains(errDb.Error(), "no rows in result set") == true {
//		record = 0
//	}
//
//	return
//}
