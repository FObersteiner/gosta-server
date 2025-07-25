package postgis

import (
	"database/sql"
	"errors"
	"fmt"

	entities "github.com/FObersteiner/gosta-core"

	gostErrors "github.com/FObersteiner/gosta-server/errors"
	"github.com/FObersteiner/gosta-server/sensorthings/odata"
)

func observedPropertyParamFactory(values map[string]any) (entities.Entity, error) {
	op := &entities.ObservedProperty{}

	for as, value := range values {
		if value == nil {
			continue
		}

		switch as {
		case asMappings[entities.EntityTypeObservedProperty][observedPropertyID]:
			op.ID = value
		case asMappings[entities.EntityTypeObservedProperty][observedPropertyName]:
			op.Name = value.(string)
		case asMappings[entities.EntityTypeObservedProperty][observedPropertyDescription]:
			op.Description = value.(string)
		case asMappings[entities.EntityTypeObservedProperty][observedPropertyDefinition]:
			op.Definition = value.(string)
		}
	}

	return op, nil
}

// GetObservedProperty returns an ObservedProperty by id
func (gdb *GostDatabase) GetObservedProperty(id any, qo *odata.QueryOptions) (*entities.ObservedProperty, error) {
	intID, ok := ToIntID(id)
	if !ok {
		return nil, gostErrors.NewRequestNotFound(errors.New("ObservedProperty does not exist"))
	}

	query, qi := gdb.QueryBuilder.CreateQuery(&entities.ObservedProperty{}, nil, intID, qo)

	observedProperty, err := processObservedProperty(gdb.Db, query, qi)
	if err != nil {
		return nil, err
	}

	return observedProperty, nil
}

// GetObservedPropertyByDatastream returns an ObservedProperty by id
func (gdb *GostDatabase) GetObservedPropertyByDatastream(id any, qo *odata.QueryOptions) (*entities.ObservedProperty, error) {
	intID, ok := ToIntID(id)
	if !ok {
		return nil, gostErrors.NewRequestNotFound(errors.New("Datastream does not exist"))
	}

	query, qi := gdb.QueryBuilder.CreateQuery(&entities.ObservedProperty{}, &entities.Datastream{}, intID, qo)

	observedProperty, err := processObservedProperty(gdb.Db, query, qi)
	if err != nil {
		return nil, err
	}

	return observedProperty, nil
}

// GetObservedProperties returns all bool, observed properties
func (gdb *GostDatabase) GetObservedProperties(qo *odata.QueryOptions) ([]*entities.ObservedProperty, int, bool, error) {
	query, qi := gdb.QueryBuilder.CreateQuery(&entities.ObservedProperty{}, nil, nil, qo)
	countSQL := gdb.QueryBuilder.CreateCountQuery(&entities.ObservedProperty{}, nil, nil, qo)

	return processObservedProperties(gdb.Db, query, qo, qi, countSQL)
}

func processObservedProperty(db *sql.DB, sql string, qi *QueryParseInfo) (*entities.ObservedProperty, error) {
	ops, _, _, err := processObservedProperties(db, sql, nil, qi, "")
	if err != nil {
		return nil, err
	}

	if len(ops) == 0 {
		return nil, gostErrors.NewRequestNotFound(errors.New("ObservedProperty not found"))
	}

	return ops[0], nil
}

func processObservedProperties(db *sql.DB, sql string, qo *odata.QueryOptions, qi *QueryParseInfo, countSQL string) ([]*entities.ObservedProperty, int, bool, error) {
	data, hasNext, err := ExecuteSelect(db, qi, sql, qo)
	if err != nil {
		return nil, 0, hasNext, fmt.Errorf("Error executing query %w", err)
	}

	obs := make([]*entities.ObservedProperty, 0)

	for _, d := range data {
		entity := d.(*entities.ObservedProperty)
		obs = append(obs, entity)
	}

	var count int
	if len(countSQL) > 0 {
		count, err = ExecuteSelectCount(db, countSQL)
		if err != nil {
			return nil, 0, hasNext, fmt.Errorf("Error executing count %w", err)
		}
	}

	return obs, count, hasNext, nil
}

// PostObservedProperty adds an ObservedProperty to the database
func (gdb *GostDatabase) PostObservedProperty(op *entities.ObservedProperty) (*entities.ObservedProperty, error) {
	var opID int

	query := fmt.Sprintf("INSERT INTO %s.observedproperty (name, definition, description) VALUES ($1, $2, $3) RETURNING id", gdb.Schema)

	err := gdb.Db.QueryRow(query, op.Name, op.Definition, op.Description).Scan(&opID)
	if err != nil {
		return nil, err
	}

	op.ID = opID

	return op, nil
}

// PutObservedProperty updates a ObservedProperty in the database
func (gdb *GostDatabase) PutObservedProperty(id any, op *entities.ObservedProperty) (*entities.ObservedProperty, error) {
	return gdb.PatchObservedProperty(id, op)
}

// ObservedPropertyExists checks if a ObservedProperty is present in the database based on a given id.
func (gdb *GostDatabase) ObservedPropertyExists(id any) bool {
	return EntityExists(gdb, id, "observedproperty")
}

// PatchObservedProperty updates a ObservedProperty in the database
func (gdb *GostDatabase) PatchObservedProperty(id any, op *entities.ObservedProperty) (*entities.ObservedProperty, error) {
	var err error

	var ok bool

	var intID int

	updates := make(map[string]any)

	if intID, ok = ToIntID(id); !ok || !gdb.ObservedPropertyExists(intID) {
		return nil, gostErrors.NewRequestNotFound(errors.New("ObservedProperty does not exist"))
	}

	if len(op.Description) > 0 {
		updates["description"] = op.Description
	}

	if len(op.Definition) > 0 {
		updates["definition"] = op.Definition
	}

	if len(op.Name) > 0 {
		updates["name"] = op.Name
	}

	if err = gdb.updateEntityColumns("observedproperty", updates, intID); err != nil {
		return nil, err
	}

	ns, _ := gdb.GetObservedProperty(intID, nil)

	return ns, nil
}

// DeleteObservedProperty tries to delete a ObservedProperty by the given id
func (gdb *GostDatabase) DeleteObservedProperty(id any) error {
	return DeleteEntity(gdb, id, "observedproperty")
}
