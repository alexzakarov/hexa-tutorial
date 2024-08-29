package base_repo

import (
	"context"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
	"main/pkg/utils/common"
	"strings"
)

type QueryFunc struct {
	Params []map[string]interface{} `json:"params"`
	Field  string                   `json:"field"`
	Func   string                   `json:"func"`
}

type SubQuery struct {
	IsSingle bool   `json:"is_one"`
	Alias    string `json:"alias"`
	Query    string `json:"query"`
}

func ConvertStatus(status int64, status_type string) interface{} {
	var data interface{}
	if status_type == "bool" {
		switch status {
		case 1:
			data = true
		case 2:
			data = false
		}
	} else if status_type == "int" {
		data = status
	}
	return data
}

type IBaseRepo[C, U, G any] interface {
	Create(C, func(id int64), func(record int64), ...QueryFunc)
	Update(int64, U, func(), func(int64))
	GetOne(int64, func(G), func(int64), ...SubQuery)
	DeleteOne(int64, func(), func(int64))
	ChangeStatus(int64, int64, func(), func(int64))
	CheckOrgIds(int64, []int64, func(int64), func(int64))
}

type BaseRepo[C, U, G any] struct {
	ctx             context.Context
	db              *pgxpool.Pool
	PrimaryKey      string   `json:"primary_key"`
	Schema          string   `json:"schema"`
	Table           string   `json:"table"`
	createFields    []string `json:"create_fields"`
	updateFields    []string `json:"update_fields"`
	getFields       []string `json:"get_fields"`
	strCreateFields string   `json:"str_create_fields"`
	strUpdateFields string   `json:"str_update_fields"`
	strGetFields    string   `json:"str_get_fields"`
	createReplacer  string   `json:"create_replacers"`
	updateReplacer  string   `json:"update_replacer"`
	softDeletable   bool     `json:"has_status"`
	statusName      string   `json:"status_name"`
	statusType      string   `json:"status_type"`
}

func NewBaseRepo[C, U, G any](ctx context.Context, db *pgxpool.Pool, schema, table, primary_key string, softDeletable bool, status_name, status_type string) IBaseRepo[C, U, G] {
	var createStruc C
	var updateStruc U
	var getStruc G
	var createJsons []string
	var updateJsons []string
	var getJsons []string
	createReplacer := ""
	updateReplacer := ""
	var errParse error

	createJsons, _, errParse = common.Convert(createStruc)
	if errParse != nil {
		println(errParse.Error())
	}

	updateJsons, _, errParse = common.Convert(updateStruc)
	if errParse != nil {
		println(errParse.Error())
	}

	getJsons, _, errParse = common.Convert(getStruc)
	if errParse != nil {
		println(errParse.Error())
	}

	for i, _ := range createJsons {
		if i < len(createJsons)-1 {
			createReplacer += fmt.Sprintf(`$%d,`, i+1)
		} else {
			createReplacer += fmt.Sprintf(`$%d`, i+1)
		}
	}
	for i, field := range updateJsons {
		if i < len(updateJsons)-1 {
			updateReplacer += fmt.Sprintf(`%s=$%d,`, field, i+1)
		} else {
			updateReplacer += fmt.Sprintf(`%s=$%d`, field, i+1)
		}
	}

	strCreateFields := strings.Join(createJsons, ",")
	strUpdateFields := strings.Join(updateJsons, ",")
	strGetFields := strings.Join(getJsons, ",")
	return &BaseRepo[C, U, G]{
		ctx:             ctx,
		db:              db,
		Schema:          schema,
		Table:           table,
		createFields:    createJsons,
		updateFields:    updateJsons,
		getFields:       getJsons,
		PrimaryKey:      primary_key,
		strCreateFields: strCreateFields,
		strUpdateFields: strUpdateFields,
		strGetFields:    strGetFields,
		createReplacer:  createReplacer,
		updateReplacer:  updateReplacer,
		softDeletable:   softDeletable,
		statusName:      status_name,
		statusType:      status_type,
	}
}

func (b *BaseRepo[C, U, G]) Create(dat C, success func(id int64), failure func(record int64), qfuncs ...QueryFunc) {
	var mapEntity []interface{}
	var errParse error
	var errDb error
	var data int64
	var createStruc C
	var funcStrings string
	var funcFields string
	var entity []interface{}
	for _, qfunc := range qfuncs {
		strParams := ""
		index := 0
		funcFields += fmt.Sprintf(",%s", qfunc.Field)

		for j, _ := range qfunc.Params {
			for param, value := range qfunc.Params[j] {

				if len(qfunc.Params) > 1 && index > 0 {
					strParams += ","
				}
				queryIndex, err := common.GetQueryIndex(createStruc, param)
				if err != nil {
					println("Error occured while getting function fields")
					continue
				}
				if value != 0 && value != "" && value != nil && queryIndex == 0 {
					strParams += fmt.Sprintf("%d", value)

				} else {
					strParams += fmt.Sprintf("$%d", queryIndex)
				}
			}
			index += 1
		}
		funcStrings += fmt.Sprintf(",%s(%s)", qfunc.Func, strParams)
	}

	_, mapEntity, errParse = common.Convert(dat)
	if errParse != nil {
		println(errParse.Error())
		failure(-1)
		return
	}
	entity = append(entity, mapEntity...)

	query := fmt.Sprintf(`INSERT INTO %s.%s ( %s %s) VALUES (%s %s) RETURNING %s`, b.Schema, b.Table, b.strCreateFields, funcFields, b.createReplacer, funcStrings, b.PrimaryKey)
	errDb = b.db.QueryRow(b.ctx, query, entity...).Scan(&data)
	println(entity)
	if errDb != nil && common.CheckStringIfContains(errDb.Error(), "duplicate key value") == false {
		println(errDb.Error())
		failure(-1)
		return
	} else if errDb != nil && common.CheckStringIfContains(errDb.Error(), "duplicate key value") == true {
		failure(-2)
		return
	}
	success(data)
	return
}

func (b *BaseRepo[C, U, G]) Update(entity_id int64, dat U, success func(), failure func(record int64)) {
	var mapEntity []interface{}
	var errParse error
	var errDb error
	var cmd pgconn.CommandTag
	var entity []interface{}

	_, mapEntity, errParse = common.Convert(dat)
	if errParse != nil {
		println(errParse.Error())
		failure(-1)
		return
	}
	entity = append(entity, mapEntity...)
	entity = append(entity, entity_id)

	query := fmt.Sprintf(`UPDATE %s.%s SET %s WHERE %s=$%d`, b.Schema, b.Table, b.updateReplacer, b.PrimaryKey, len(entity))
	cmd, errDb = b.db.Exec(b.ctx, query, entity...)
	if errDb != nil {
		println(errDb.Error())
		failure(-1)
		return
	} else if errDb == nil && cmd.RowsAffected() == 0 {
		failure(0)
		return
	}
	success()
	return
}

func (b *BaseRepo[C, U, G]) GetOne(entity_id int64, success func(data G), failure func(record int64), sub_queries ...SubQuery) {
	var errDb error
	var entity G
	var qsReplaced []string

	for _, q := range sub_queries {
		var injected string
		replaced := strings.ReplaceAll(q.Query, "[primary_key]", fmt.Sprintf("ent.%s", b.PrimaryKey))
		if q.IsSingle {
			injected = fmt.Sprintf(`,(SELECT TO_JSON(ENTITY) FROM (%s) ENTITY) as %s`, replaced, q.Alias)
		} else {
			injected = fmt.Sprintf(`,(SELECT JSON_AGG(ENTITY) FROM (%s) ENTITY) as %s`, replaced, q.Alias)
		}
		qsReplaced = append(qsReplaced, injected)
	}

	subQs := strings.Join(qsReplaced, " ")

	query := fmt.Sprintf(`SELECT TO_JSON(ENTITY) FROM (SELECT %s %s FROM %s.%s ent WHERE %s=$1) ENTITY`, b.strGetFields, subQs, b.Schema, b.Table, b.PrimaryKey)
	errDb = b.db.QueryRow(b.ctx, query, entity_id).Scan(&entity)
	if errDb != nil && common.CheckStringIfContains(errDb.Error(), "no rows in result set") == false {
		println(errDb.Error())
		failure(-1)
		return
	} else if errDb != nil && common.CheckStringIfContains(errDb.Error(), "no rows in result set") == true {
		failure(0)
		return
	}
	success(entity)
	return
}

func (b *BaseRepo[C, U, G]) DeleteOne(entity_id int64, success func(), failure func(record int64)) {
	var query string
	var cmd pgconn.CommandTag
	var errDb error
	var statusClause string
	var statusUpdateClause string

	if b.softDeletable {
		statusClause = fmt.Sprintf(`AND %s=%v`, b.statusName, ConvertStatus(1, b.statusType))
		statusUpdateClause = fmt.Sprintf(`%s=%v`, b.statusName, ConvertStatus(2, b.statusType))
		query = fmt.Sprintf(`UPDATE %s.%s SET %s WHERE %s=$1 %s`, b.Schema, b.Table, statusUpdateClause, b.PrimaryKey, statusClause)
	} else {
		query = fmt.Sprintf(`DELETE FROM %s.%s WHERE %s=$1`, b.Schema, b.Table, b.PrimaryKey)
	}

	cmd, errDb = b.db.Exec(b.ctx, query, entity_id)
	if errDb != nil {
		println(errDb.Error())
		failure(-1)
		return
	} else if errDb == nil && cmd.RowsAffected() == 0 {
		failure(0)
		return
	}
	success()
	return
}

func (b *BaseRepo[C, U, G]) ChangeStatus(entity_id, status int64, success func(), failure func(record int64)) {
	var query string
	var cmd pgconn.CommandTag
	var errDb error

	query = fmt.Sprintf(`UPDATE %s.%s SET %s=$1 WHERE %s=$2`, b.Schema, b.Table, b.statusName, b.PrimaryKey)
	cmd, errDb = b.db.Exec(b.ctx, query, ConvertStatus(status, b.statusType), entity_id)
	if errDb != nil {
		println(errDb.Error())
		failure(-1)
		return
	} else if errDb == nil && cmd.RowsAffected() == 0 {
		failure(0)
		return
	}
	success()
	return
}

func (b *BaseRepo[C, U, G]) CheckOrgIds(entity_id int64, org_ids []int64, success func(id int64), failure func(record int64)) {
	var query string
	var errDb error
	var orgId int64
	var stringArr []string

	for _, id := range org_ids {
		stringArr = append(stringArr, fmt.Sprintf("%d", id))
	}

	query = fmt.Sprintf(`SELECT org_id FROM %s.%s WHERE %s=$1 AND org_id IN (%s)`, b.Schema, b.Table, b.PrimaryKey, strings.Join(stringArr, ","))
	errDb = b.db.QueryRow(b.ctx, query, entity_id).Scan(&orgId)
	if errDb != nil && common.CheckStringIfContains(errDb.Error(), "no rows in result set") == false {
		println(errDb.Error())
		failure(-1)
		return
	} else if errDb != nil && common.CheckStringIfContains(errDb.Error(), "no rows in result set") == true {
		failure(0)
		return
	}
	success(orgId)
	return
}
