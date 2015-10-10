package namedtuple

import (
	"errors"
	"math"
)

type TupleBuilder struct {
	fields    map[string]Field
	offsets   map[string]int
	tupleType TupleType
	buffer    []byte
	pos       int
}

func NewBuilder(t TupleType, buffer []byte) TupleBuilder {

	// init instance variables
	fields := make(map[string]Field)
	offsets := make(map[string]int)

	// populate instance fields for builder
	for _, version := range t.Versions() {
		for _, field := range version.Fields {
			fields[field.Name] = field
			offsets[field.Name] = 0
		}
	}

	// create new builder
	return TupleBuilder{fields: fields, offsets: offsets, tupleType: t, buffer: buffer, pos: 0}
}

func (b TupleBuilder) available() int {
	return len(b.buffer) - b.pos
}

func (b TupleBuilder) reset() {
	b.pos = 0
}

func (t *TupleBuilder) typeCheck(fieldName string, fieldType FieldType) error {
	field, exists := t.fields[fieldName]
	if !exists {
		return errors.New("Field does not exist: " + fieldName)
	}

	if field.Type != fieldType {
		return errors.New("Incorrect field type: " + fieldName)
	}

	return nil
}

func (b *TupleBuilder) Build() (Tuple, error) {
	defer b.reset()
	header, err := b.newTupleHeader()
	if err != nil {
		return NIL, err
	}
	return Tuple{data: b.buffer[:b.pos], Header: header}, nil
}

func (b *TupleBuilder) newTupleHeader() (TupleHeader, error) {

	// validation of required fields
	var tupleVersion uint8
	var missingField string
	var fieldSize uint8
	var fieldCount int

	totalFieldCount := uint32(len(b.tupleType.fields))
	offsets := make([]uint64, totalFieldCount)

	// iterate over all the versions
	for _, version := range b.tupleType.Versions() {
	OUTER:

		// iterate over all the fields for the current version
		for _, field := range version.Fields {

			// get offset for field
			offset, exists := b.offsets[field.Name]

			// if the field is required, determine if it has been added to the builder
			if field.Required {

				// if the field has not been written
				// exit the loop and save the missing field name
				if !exists {
					missingField = field.Name
					break OUTER
				}

				// set byte offset of field in tuple data
				offsets[fieldCount] = uint64(offset)
			} else {

				// if the optional fields was not written, encode a maximum offset
				if !exists {

					// set byte offset of field in tuple data
					offsets[fieldCount] = uint64(math.MaxUint64)

				} else {
					// if the optional field does exist
					// set byte offset of field in tuple data
					offsets[fieldCount] = uint64(offset)
				}
			}
			fieldCount++
		}

		// increment the version number after all required fields have been satisfied
		tupleVersion++
	}

	// If the first version is missing a field, return an error
	// At least one version must contain all the required fields.
	// The version number will increment for each version which
	// contains all the required fields.
	if tupleVersion < 1 {
		return TupleHeader{}, errors.New("Missing required field: " + missingField)
	}

	// TODO: Add Field level validation

	// Calculate minimum offset for accessing all fields in data
	// If the total data size is < 256 bytes, all field offsets
	if b.pos < math.MaxUint8-1 {
		fieldSize = 1
	} else if b.pos < math.MaxUint16 {
		fieldSize = 2
	} else if b.pos < math.MaxUint32 {
		fieldSize = 4
	} else {
		fieldSize = 8
	}

	return TupleHeader{
		ProtocolVersion: 1,
		TupleVersion:    tupleVersion,
		NamespaceHash:   b.tupleType.NamespaceHash,
		Hash:            b.tupleType.Hash,
		FieldCount:      totalFieldCount,
		FieldSize:       fieldSize,
		ContentLength:   uint64(b.pos),
		Offsets:         offsets,
		Type:            b.tupleType,
	}, nil
}
