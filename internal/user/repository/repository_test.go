package repository

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"io"
	"reflect"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

var testSQLDriverID int64

type testSQLStep struct {
	kind string
	sql  string
	args []any

	rows   *testSQLRows
	result testSQLResult
	err    error
}

type testSQLScript struct {
	t     *testing.T
	mu    sync.Mutex
	steps []testSQLStep
}

func newTestDB(t *testing.T, steps ...testSQLStep) *sql.DB {
	t.Helper()

	driverName := fmt.Sprintf("user_repository_test_%d", atomic.AddInt64(&testSQLDriverID, 1))
	script := &testSQLScript{t: t, steps: steps}
	sql.Register(driverName, testSQLDriver{script: script})

	db, err := sql.Open(driverName, "")
	if err != nil {
		t.Fatalf("open test db: %v", err)
	}

	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Fatalf("close test db: %v", err)
		}
		script.assertDone()
	})

	return db
}

func expectQuery(query string, args []any, rows *testSQLRows) testSQLStep {
	return testSQLStep{kind: "query", sql: query, args: args, rows: rows}
}

func expectQueryError(query string, args []any, err error) testSQLStep {
	return testSQLStep{kind: "query", sql: query, args: args, err: err}
}

func expectExec(query string, args []any, rowsAffected int64) testSQLStep {
	return testSQLStep{
		kind:   "exec",
		sql:    query,
		args:   args,
		result: testSQLResult{rowsAffected: rowsAffected},
	}
}

func expectExecError(query string, args []any, err error) testSQLStep {
	return testSQLStep{kind: "exec", sql: query, args: args, err: err}
}

func rows(columns []string, values ...[]driver.Value) *testSQLRows {
	return &testSQLRows{columns: columns, values: values}
}

type testSQLDriver struct {
	script *testSQLScript
}

func (d testSQLDriver) Open(string) (driver.Conn, error) {
	return &testSQLConn{script: d.script}, nil
}

type testSQLConn struct {
	script *testSQLScript
}

func (c *testSQLConn) Prepare(string) (driver.Stmt, error) {
	return nil, errors.New("prepare is not supported by test driver")
}

func (c *testSQLConn) Close() error {
	return nil
}

func (c *testSQLConn) Begin() (driver.Tx, error) {
	return &testSQLTx{}, nil
}

func (c *testSQLConn) BeginTx(context.Context, driver.TxOptions) (driver.Tx, error) {
	return &testSQLTx{}, nil
}

func (c *testSQLConn) QueryContext(_ context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	step := c.script.next("query", query, args)
	if step.err != nil {
		return nil, step.err
	}
	return step.rows, nil
}

func (c *testSQLConn) ExecContext(_ context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	step := c.script.next("exec", query, args)
	if step.err != nil {
		return nil, step.err
	}
	return step.result, nil
}

type testSQLTx struct{}

func (testSQLTx) Commit() error {
	return nil
}

func (testSQLTx) Rollback() error {
	return nil
}

func (s *testSQLScript) next(kind string, query string, args []driver.NamedValue) testSQLStep {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.steps) == 0 {
		s.t.Fatalf("unexpected %s query: %s", kind, query)
	}

	step := s.steps[0]
	s.steps = s.steps[1:]

	if step.kind != kind {
		s.t.Fatalf("operation kind = %q, want %q for query %s", kind, step.kind, query)
	}
	if step.sql != query {
		s.t.Fatalf("query = %q, want %q", query, step.sql)
	}
	assertSQLArgs(s.t, args, step.args)

	return step
}

func (s *testSQLScript) assertDone() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.steps) != 0 {
		s.t.Fatalf("remaining sql expectations = %d", len(s.steps))
	}
}

func assertSQLArgs(t *testing.T, got []driver.NamedValue, want []any) {
	t.Helper()

	if len(got) != len(want) {
		t.Fatalf("arg count = %d, want %d; args = %#v", len(got), len(want), got)
	}
	for i := range got {
		if !equalSQLValue(got[i].Value, want[i]) {
			t.Fatalf("arg[%d] = %#v (%T), want %#v (%T)", i, got[i].Value, got[i].Value, want[i], want[i])
		}
	}
}

func equalSQLValue(got any, want any) bool {
	gotTime, gotIsTime := got.(time.Time)
	wantTime, wantIsTime := want.(time.Time)
	if gotIsTime && wantIsTime {
		return gotTime.Equal(wantTime)
	}
	return reflect.DeepEqual(got, want) || fmt.Sprint(got) == fmt.Sprint(want)
}

type testSQLResult struct {
	rowsAffected    int64
	rowsAffectedErr error
}

func (r testSQLResult) LastInsertId() (int64, error) {
	return 0, errors.New("last insert id is not supported")
}

func (r testSQLResult) RowsAffected() (int64, error) {
	if r.rowsAffectedErr != nil {
		return 0, r.rowsAffectedErr
	}
	return r.rowsAffected, nil
}

type testSQLRows struct {
	columns []string
	values  [][]driver.Value
	index   int
	err     error
}

func (r *testSQLRows) Columns() []string {
	return r.columns
}

func (r *testSQLRows) Close() error {
	return nil
}

func (r *testSQLRows) Next(dest []driver.Value) error {
	if r.err != nil {
		return r.err
	}
	if r.index >= len(r.values) {
		return io.EOF
	}
	copy(dest, r.values[r.index])
	r.index++
	return nil
}
