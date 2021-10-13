package model

import (
	"fmt"
	"reflect"

	"github.com/ovn-org/libovsdb/mapper"
	"github.com/ovn-org/libovsdb/ovsdb"
)

// A DatabaseModel represents libovsdb's metadata about the database.
// It's the result of combining the client's DatabaseModelRequest and the server's Schema
type DatabaseModel struct {
	valid    bool
	request  *DatabaseModelRequest
	schema   *ovsdb.DatabaseSchema
	mapper   *mapper.Mapper
	metadata map[string]*mapper.Metadata
}

// NewDatabaseModel returns a new DatabaseModel
func NewDatabaseModel(schema *ovsdb.DatabaseSchema, request *DatabaseModelRequest) (*DatabaseModel, []error) {
	dbModel := NewPartialDatabaseModel(request)
	errs := dbModel.SetSchema(schema)
	if len(errs) > 0 {
		return nil, errs
	}
	return dbModel, nil
}

// NewPartialDatabaseModel returns a DatabaseModel what does not have a schema yet
func NewPartialDatabaseModel(request *DatabaseModelRequest) *DatabaseModel {
	return &DatabaseModel{
		valid:   false,
		request: request,
	}
}

// Valid returns whether the DatabaseModel is fully functional
func (db *DatabaseModel) Valid() bool {
	return db.valid
}

// SetSchema adds the Schema to the DatabaseModel making it valid if it was not before
func (db *DatabaseModel) SetSchema(schema *ovsdb.DatabaseSchema) []error {
	errors := db.request.validate(schema)
	if len(errors) > 0 {
		return errors
	}
	db.schema = schema
	db.mapper = mapper.NewMapper(schema)
	errs := db.generateModelInfo()
	if len(errs) > 0 {
		return errs
	}
	db.valid = true
	return []error{}
}

// ClearSchema removes the Schema from the DatabaseModel making it not valid
func (db *DatabaseModel) ClearSchema() {
	db.schema = nil
	db.mapper = nil
	db.valid = false
}

// Request returns the DatabaseModel's request
func (db *DatabaseModel) Request() *DatabaseModelRequest {
	return db.request
}

// Schema returns the DatabaseModel's schema
func (db *DatabaseModel) Schema() *ovsdb.DatabaseSchema {
	return db.schema
}

// Mapper returns the DatabaseModel's mapper
func (db *DatabaseModel) Mapper() *mapper.Mapper {
	return db.mapper
}

// NewModel returns a new instance of a model from a specific string
func (db *DatabaseModel) NewModel(table string) (Model, error) {
	mtype, ok := db.request.types[table]
	if !ok {
		return nil, fmt.Errorf("table %s not found in database model", string(table))
	}
	model := reflect.New(mtype.Elem())
	return model.Interface().(Model), nil
}

// Types returns the DatabaseModel Types
// the DatabaseModel types is a map of reflect.Types indexed by string
// The reflect.Type is a pointer to a struct that contains 'ovs' tags
// as described above. Such pointer to struct also implements the Model interface
func (db *DatabaseModel) Types() map[string]reflect.Type {
	return db.request.types
}

// FindTable returns the string associated with a reflect.Type or ""
func (db *DatabaseModel) FindTable(mType reflect.Type) string {
	for table, tType := range db.request.types {
		if tType == mType {
			return table
		}
	}
	return ""
}

// generateModelMetadata creates metadata objects from all models included in the
// database and caches them for future re-use
func (db *DatabaseModel) generateModelInfo() []error {
	errors := []error{}
	metadata := make(map[string]*mapper.Metadata, len(db.request.types))
	for tableName := range db.request.types {
		tableSchema := db.schema.Table(tableName)
		if tableSchema == nil {
			continue
		}
		obj, err := db.NewModel(tableName)
		if err != nil {
			errors = append(errors, err)
			continue
		}
		info, err := mapper.NewInfo(tableName, tableSchema, obj, db.request.compat)
		if err != nil {
			errors = append(errors, err)
			continue
		}
		metadata[tableName] = info.Metadata
	}
	db.metadata = metadata
	return errors
}

// NewModelInfo returns a mapper.Info object based on a provided model
func (db *DatabaseModel) NewModelInfo(obj interface{}) (*mapper.Info, error) {
	meta, ok := db.metadata[db.FindTable(reflect.TypeOf(obj))]
	if !ok {
		return nil, ovsdb.NewErrWrongType("NewModelInfo", "type that is part of the DatabaseModel", obj)
	}
	return &mapper.Info{
		Obj:      obj,
		Metadata: meta,
	}, nil
}

func (db *DatabaseModel) HasColumn(tableName, column string) bool {
	meta, ok := db.metadata[tableName]
	if !ok {
		return false
	}
	_, ok = meta.Fields[column]
	return ok
}
