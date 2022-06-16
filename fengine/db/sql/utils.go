package sql

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/jackc/pgx/v4"
)

type RowCache struct {
	Mapper  *Mapper
	unsafe  bool
	started bool
	fields  [][]int
	values  []interface{}
}

func InitScan() {

}

func Columns(r pgx.Rows) []string {
	descriptions := r.FieldDescriptions()
	columns := make([]string, len(descriptions))
	for i, f := range descriptions {
		columns[i] = string(f.Name)
	}
	return columns
}

func StructScan(row pgx.Rows, dest interface{}) error {
	v := reflect.ValueOf(dest)

	if v.Kind() != reflect.Ptr {
		return errors.New("must pass a pointer, not a value, to StructScan destination")
	}

	v = v.Elem()
	var r RowCache
	if !r.started {
		columns, err := r.Columns()
		if err != nil {
			return err
		}
		m := r.Mapper

		r.fields = m.TraversalsByName(v.Type(), columns)
		// if we are not unsafe and are missing fields, return an error
		if f, err := missingFields(r.fields); err != nil && !r.unsafe {
			return fmt.Errorf("missing destination name %s in %T", columns[f], dest)
		}
		r.values = make([]interface{}, len(columns))
		r.started = true
	}

	err := fieldsByTraversal(v, r.fields, r.values, true)
	if err != nil {
		return err
	}
	// scan into the struct field pointers and append to our results
	err = r.Scan(r.values...)
	if err != nil {
		return err
	}
	return r.Err()
}

func (r RowCache) Columns(values ...any) ([]string, error) {
	return nil, nil
}

func (r RowCache) Scan(values ...any) error {
	return nil
}

func (r RowCache) Err(values ...any) error {
	return nil
}

// fieldsByName fills a values interface with fields from the passed value based on the traversals in int. If ptrs is
// true, return addresses instead of values. We write this instead of using FieldsByName to save allocations and map
// lookups when iterating over many rows.  Empty traversals will get an interface pointer. Because of the necessity of
// requesting ptrs or values, it's considered a bit too specialized for inclusion in reflectx itself.
func fieldsByTraversal(v reflect.Value, traversals [][]int, values []interface{}, ptrs bool) error {
	v = reflect.Indirect(v)
	if v.Kind() != reflect.Struct {
		return errors.New("argument not a struct")
	}

	for i, traversal := range traversals {
		if len(traversal) == 0 {
			values[i] = new(interface{})
			continue
		}
		f := FieldByIndexes(v, traversal)
		if ptrs {
			values[i] = f.Addr().Interface()
		} else {
			values[i] = f.Interface()
		}
	}
	return nil
}

func missingFields(transversals [][]int) (field int, err error) {
	for i, t := range transversals {
		if len(t) == 0 {
			return i, errors.New("missing field")
		}
	}
	return 0, nil
}

func MapScan(row pgx.Rows, mp *map[string]interface{}) {

}
