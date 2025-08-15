package sqltest_test

import (
	"database/sql"
	"errors"
	"fmt"
	"testing"

	"github.com/aereal/optional"
	_ "github.com/mattn/go-sqlite3"
)

func TestOption_value(t *testing.T) {
	testCases := []struct {
		input   optional.Option[int]
		want    optional.Option[int]
		wantErr error
	}{
		{
			input:   optional.Some(123),
			want:    optional.Some(123),
			wantErr: nil,
		},
		{
			input:   optional.None[int](),
			want:    optional.None[int](),
			wantErr: nil,
		},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("input=%#v", tc.input), func(t *testing.T) {
			db, err := sql.Open("sqlite3", "file:test.db?mode=memory")
			if err != nil {
				t.Fatal(err)
			}
			ctx := t.Context()
			if _, err := db.ExecContext(ctx, `create table ids (id integer primary key autoincrement, n integer)`); err != nil {
				t.Fatal(err)
			}
			ret, insertErr := db.ExecContext(ctx, `insert into ids (n) values (?)`, tc.input)
			if !errors.Is(tc.wantErr, insertErr) {
				t.Errorf("error:\n\twant: (%T) %s\n\t got: (%T) %s", tc.wantErr, tc.wantErr, insertErr, insertErr)
			}
			if insertErr != nil {
				return
			}
			id, err := ret.LastInsertId()
			if err != nil {
				t.Fatal(err)
			}
			row := db.QueryRowContext(ctx, `select n from ids where id = ?`, id)
			if err := row.Err(); err != nil {
				t.Fatal(err)
			}
			var got optional.Option[int]
			if err := row.Scan(&got); err != nil {
				t.Fatal(err)
			}
			if !optional.Equal(got, tc.want) {
				t.Errorf("inserted value:\n\twant: %#v\n\t got: %#v", tc.want, got)
			}
		})
	}
}

func TestOption_scan(t *testing.T) {
	testCases := []struct {
		query     string
		wantValue optional.Option[int]
		wantErr   error
	}{
		{
			query:     `select 123`,
			wantValue: optional.Some(123),
			wantErr:   nil,
		},
		{
			query:     `select 'abc'`,
			wantValue: optional.None[int](),
			wantErr:   literalError(`sql: Scan error on column index 0, name "'abc'": converting driver.Value type string ("abc") to a int: invalid syntax`),
		},
		{
			query:     `select null`,
			wantValue: optional.None[int](),
			wantErr:   nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.query, func(t *testing.T) {
			db, err := sql.Open("sqlite3", "file:test.db?mode=memory")
			if err != nil {
				t.Fatal(err)
			}
			row := db.QueryRowContext(t.Context(), tc.query)
			if err := row.Err(); err != nil {
				t.Fatal(err)
			}
			var got optional.Option[int]
			gotScanErr := row.Scan(&got)
			if !errors.Is(tc.wantErr, gotScanErr) {
				t.Errorf("error:\n\twant: (%T) %s\n\t got: (%T) %s", tc.wantErr, tc.wantErr, gotScanErr, gotScanErr)
			}
			if gotScanErr != nil {
				return
			}
			if !optional.Equal(got, tc.wantValue) {
				t.Errorf("value:\n\twant: %#v\n\t got: %#v", tc.wantValue, got)
			}
		})
	}
}

type literalError string

var _ error = literalError("")

func (e literalError) Error() string { return string(e) }

func (e literalError) Is(other error) bool {
	return other.Error() == string(e)
}
