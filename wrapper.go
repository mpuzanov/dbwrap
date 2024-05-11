package dbwrap

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	"go.uber.org/multierr"
)

func sqlErr(err error, query string, args ...interface{}) error {
	return fmt.Errorf(`run query "%s" with args %+v: %w`, query, args, err)
}

func namedQuery(query string, arg interface{}) (nq string, args []interface{}, err error) {
	nq, args, err = sqlx.Named(query, arg)
	if err != nil {
		return "", nil, sqlErr(err, query, args...)
	}
	return nq, args, nil
}

// Exec Выполнение запроса DML
func (d *DBSQL) Exec(query string, args ...interface{}) (int64, error) {

	// ограничим время выполнения запроса
	dur := time.Duration(d.Cfg.TimeoutQuery) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), dur)
	defer cancel()

	result, err := d.DBX.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, sqlErr(err, query, args...)
	}
	count, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return count, nil
}

// NamedExec Выполнение запроса DML
func (d *DBSQL) NamedExec(query string, arg interface{}) (int64, error) {

	nq, args, err := namedQuery(query, arg)
	if err != nil {
		return 0, err
	}

	return d.Exec(d.DBX.Rebind(nq), args...)
}

// Select получаем данные из запроса в слайс структур
//
// var users []User
//
// err := ts.db.Select(&users, "select * from users")
func (d *DBSQL) Select(dest interface{}, query string, args ...interface{}) error {

	// ограничим время выполнения запроса
	dur := time.Duration(d.Cfg.TimeoutQuery) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), dur)
	defer cancel()

	if err := sqlx.SelectContext(ctx, d.DBX, dest, query, args...); err != nil {
		return sqlErr(err, query, args...)
	}

	return nil
}

// NamedSelect получаем данные из запроса в слайс структур
//
// var users []User
//
// err := ts.db.NamedSelect(&users, "select * from users where name=:Name", map[string]interface{}{"Name": "admin"})
func (d *DBSQL) NamedSelect(dest interface{}, query string, arg interface{}) error {

	nq, args, err := namedQuery(query, arg)
	if err != nil {
		return err
	}

	return d.Select(dest, d.DBX.Rebind(nq), args...)
}

// SelectMaps ...
func (d *DBSQL) SelectMaps(query string, args ...interface{}) (ret []map[string]interface{}, err error) {

	dur := time.Duration(d.Cfg.TimeoutQuery) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), dur)
	defer cancel()

	rows, err := d.DBX.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, sqlErr(err, query, args...)
	}

	defer func() {
		err = multierr.Combine(err, rows.Close())
	}()

	ret = []map[string]interface{}{}
	numCols := -1
	for rows.Next() {
		var m map[string]interface{}
		if numCols < 0 {
			m = map[string]interface{}{}
		} else {
			m = make(map[string]interface{}, numCols)
		}

		if err = rows.MapScan(m); err != nil {
			return nil, sqlErr(err, query, args...)
		}

		for key, val := range m {
			switch v := val.(type) {
			case []byte:
				if resFloat, err := strconv.ParseFloat(string(v), 64); err == nil {
					m[key] = resFloat
					continue
				}
				if v, ok := val.([]uint8); ok {
					m[key] = string(v)
				} else {
					m[key] = v
				}
			default:
				m[key] = v
			}
		}

		ret = append(ret, m)
		numCols = len(m)
	}

	if err = rows.Err(); err != nil {
		return nil, sqlErr(err, query, args...)
	}

	return ret, nil
}

// NamedSelectMaps ...
func (d *DBSQL) NamedSelectMaps(query string, arg interface{}) (ret []map[string]interface{}, err error) {
	nq, args, err := namedQuery(query, arg)
	if err != nil {
		return nil, err
	}

	return d.SelectMaps(d.DBX.Rebind(nq), args...)
}

// Get ...
func (d *DBSQL) Get(dest interface{}, query string, args ...interface{}) error {

	// ограничим время выполнения запроса
	dur := time.Duration(d.Cfg.TimeoutQuery) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), dur)
	defer cancel()

	if err := sqlx.GetContext(ctx, d.DBX, dest, query, args...); err != nil {
		return sqlErr(err, query, args...)
	}

	return nil
}

// NamedGet ...
func (d *DBSQL) NamedGet(dest interface{}, query string, arg interface{}) error {
	nq, args, err := namedQuery(query, arg)
	if err != nil {
		return err
	}

	return d.Get(dest, d.DBX.Rebind(nq), args...)
}

// GetMap ...
func (d *DBSQL) GetMap(query string, args ...interface{}) (ret map[string]interface{}, err error) {
	// ограничим время выполнения запроса
	dur := time.Duration(d.Cfg.TimeoutQuery) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), dur)
	defer cancel()

	row := d.DBX.QueryRowxContext(ctx, query, args...)
	if row.Err() != nil {
		return nil, sqlErr(row.Err(), query, args...)
	}

	ret = map[string]interface{}{}
	if err := row.MapScan(ret); err != nil {
		return nil, sqlErr(err, query, args...)
	}

	for key, val := range ret {
		switch v := val.(type) {
		case []byte:
			if resFloat, err := strconv.ParseFloat(string(v), 64); err == nil {
				ret[key] = resFloat
				continue
			}
			if v, ok := val.([]uint8); ok {
				ret[key] = string(v)
			} else {
				ret[key] = v
			}
		default:
			ret[key] = v
		}
	}

	return ret, nil
}

// NamedGetMap ...
func (d *DBSQL) NamedGetMap(query string, arg interface{}) (ret map[string]interface{}, err error) {
	nq, args, err := namedQuery(query, arg)
	if err != nil {
		return nil, err
	}

	return d.GetMap(d.DBX.Rebind(nq), args...)
}
