package dbwrap

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	"go.uber.org/multierr"
)

func sqlErr(err error, query string, args ...any) error {
	return fmt.Errorf(`run query "%s" with args %+v: %w`, query, args, err)
}

func namedQuery(query string, arg any) (nq string, args []any, err error) {
	nq, args, err = sqlx.Named(query, arg)
	if err != nil {
		return "", nil, sqlErr(err, query, args...)
	}
	return nq, args, nil
}

// ExecContext Выполнение запроса DML.
func (d *DBSQL) ExecContext(ctx context.Context, query string, args ...any) (int64, error) {

	// ограничим время выполнения запроса по умолчанию
	dur := time.Duration(d.Cfg.TimeoutQuery) * time.Second
	ctx, cancel := context.WithTimeout(ctx, dur)
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

// NamedExecContext Выполнение запроса DML.
func (d *DBSQL) NamedExecContext(ctx context.Context, query string, arg any) (int64, error) {

	nq, args, err := namedQuery(query, arg)
	if err != nil {
		return 0, err
	}

	return d.ExecContext(ctx, d.DBX.Rebind(nq), args...)
}

// SelectContext получаем данные из запроса в слайс структур.
//
// var users []User
//
// err := ts.db.Select(ctx, &users, "select * from users")
func (d *DBSQL) SelectContext(ctx context.Context, dest any, query string, args ...any) error {

	// ограничим время выполнения запроса по умолчанию
	dur := time.Duration(d.Cfg.TimeoutQuery) * time.Second
	ctx, cancel := context.WithTimeout(ctx, dur)
	defer cancel()

	if err := sqlx.SelectContext(ctx, d.DBX, dest, query, args...); err != nil {
		return sqlErr(err, query, args...)
	}

	return nil
}

// NamedSelectContext получаем данные из запроса в слайс структур
//
// var users []User
//
// err := ts.db.NamedSelectContext(ctx, &users, "select * from users where name=:Name", map[string]any{"Name": "admin"})
func (d *DBSQL) NamedSelectContext(ctx context.Context, dest any, query string, arg any) error {

	nq, args, err := namedQuery(query, arg)
	if err != nil {
		return err
	}

	return d.SelectContext(ctx, dest, d.DBX.Rebind(nq), args...)
}

// SelectMapsContext ...
func (d *DBSQL) SelectMapsContext(ctx context.Context, query string, args ...any) (ret []map[string]any, err error) {

	dur := time.Duration(d.Cfg.TimeoutQuery) * time.Second
	ctx, cancel := context.WithTimeout(ctx, dur)
	defer cancel()

	rows, err := d.DBX.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, sqlErr(err, query, args...)
	}

	defer func() {
		err = multierr.Combine(err, rows.Close())
	}()

	ret = []map[string]any{}
	numCols := -1
	for rows.Next() {
		var m map[string]any
		if numCols < 0 {
			m = map[string]any{}
		} else {
			m = make(map[string]any, numCols)
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

// NamedSelectMapsContext ...
func (d *DBSQL) NamedSelectMapsContext(ctx context.Context, query string, arg any) (ret []map[string]any, err error) {
	nq, args, err := namedQuery(query, arg)
	if err != nil {
		return nil, err
	}

	return d.SelectMapsContext(ctx, d.DBX.Rebind(nq), args...)
}

// GetContext ...
func (d *DBSQL) GetContext(ctx context.Context, dest any, query string, args ...any) error {

	// ограничим время выполнения запроса по умолчанию
	dur := time.Duration(d.Cfg.TimeoutQuery) * time.Second
	ctx, cancel := context.WithTimeout(ctx, dur)
	defer cancel()

	if err := sqlx.GetContext(ctx, d.DBX, dest, query, args...); err != nil {
		return sqlErr(err, query, args...)
	}

	return nil
}

// NamedGetContext ...
func (d *DBSQL) NamedGetContext(ctx context.Context, dest any, query string, arg any) error {
	nq, args, err := namedQuery(query, arg)
	if err != nil {
		return err
	}

	return d.GetContext(ctx, dest, d.DBX.Rebind(nq), args...)
}

// GetMapContext ...
func (d *DBSQL) GetMapContext(ctx context.Context, query string, args ...any) (ret map[string]any, err error) {
	// ограничим время выполнения запроса по умолчанию
	dur := time.Duration(d.Cfg.TimeoutQuery) * time.Second
	ctx, cancel := context.WithTimeout(ctx, dur)
	defer cancel()

	row := d.DBX.QueryRowxContext(ctx, query, args...)
	if row.Err() != nil {
		return nil, sqlErr(row.Err(), query, args...)
	}

	ret = map[string]any{}
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

// NamedGetMapContext ...
func (d *DBSQL) NamedGetMapContext(ctx context.Context, query string, arg any) (ret map[string]any, err error) {
	nq, args, err := namedQuery(query, arg)
	if err != nil {
		return nil, err
	}

	return d.GetMapContext(ctx, d.DBX.Rebind(nq), args...)
}
