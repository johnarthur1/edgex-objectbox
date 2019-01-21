// Code generated by ObjectBox; DO NOT EDIT.

package obx

import (
	. "github.com/edgexfoundry/edgex-go/pkg/models"
	"github.com/google/flatbuffers/go"
	"github.com/objectbox/objectbox-go/objectbox"
	"github.com/objectbox/objectbox-go/objectbox/fbutils"
)

type reading_EntityInfo struct {
	Id  objectbox.TypeId
	Uid uint64
}

var ReadingBinding = reading_EntityInfo{
	Id:  5,
	Uid: 5720922007709447864,
}

// Reading_ contains type-based Property helpers to facilitate some common operations such as Queries.
var Reading_ = struct {
	Id       *objectbox.PropertyUint64
	Pushed   *objectbox.PropertyInt64
	Created  *objectbox.PropertyInt64
	Origin   *objectbox.PropertyInt64
	Modified *objectbox.PropertyInt64
	Device   *objectbox.PropertyString
	Name     *objectbox.PropertyString
	Value    *objectbox.PropertyString
}{
	Id: &objectbox.PropertyUint64{
		BaseProperty: &objectbox.BaseProperty{
			Id: 1,
			Entity: &objectbox.Entity{
				Id: 5,
			},
		},
	},
	Pushed: &objectbox.PropertyInt64{
		BaseProperty: &objectbox.BaseProperty{
			Id: 2,
			Entity: &objectbox.Entity{
				Id: 5,
			},
		},
	},
	Created: &objectbox.PropertyInt64{
		BaseProperty: &objectbox.BaseProperty{
			Id: 3,
			Entity: &objectbox.Entity{
				Id: 5,
			},
		},
	},
	Origin: &objectbox.PropertyInt64{
		BaseProperty: &objectbox.BaseProperty{
			Id: 4,
			Entity: &objectbox.Entity{
				Id: 5,
			},
		},
	},
	Modified: &objectbox.PropertyInt64{
		BaseProperty: &objectbox.BaseProperty{
			Id: 5,
			Entity: &objectbox.Entity{
				Id: 5,
			},
		},
	},
	Device: &objectbox.PropertyString{
		BaseProperty: &objectbox.BaseProperty{
			Id: 6,
			Entity: &objectbox.Entity{
				Id: 5,
			},
		},
	},
	Name: &objectbox.PropertyString{
		BaseProperty: &objectbox.BaseProperty{
			Id: 7,
			Entity: &objectbox.Entity{
				Id: 5,
			},
		},
	},
	Value: &objectbox.PropertyString{
		BaseProperty: &objectbox.BaseProperty{
			Id: 8,
			Entity: &objectbox.Entity{
				Id: 5,
			},
		},
	},
}

// GeneratorVersion is called by the ObjectBox to verify the compatibility of the generator used to generate this code
func (reading_EntityInfo) GeneratorVersion() int {
	return 1
}

// AddToModel is called by the ObjectBox during model build
func (reading_EntityInfo) AddToModel(model *objectbox.Model) {
	model.Entity("Reading", 5, 5720922007709447864)
	model.Property("Id", objectbox.PropertyType_Long, 1, 506318890259717221)
	model.PropertyFlags(objectbox.PropertyFlags_ID)
	model.Property("Pushed", objectbox.PropertyType_Long, 2, 857645422171573341)
	model.Property("Created", objectbox.PropertyType_Long, 3, 2293804091244003818)
	model.Property("Origin", objectbox.PropertyType_Long, 4, 6649600333394362209)
	model.Property("Modified", objectbox.PropertyType_Long, 5, 3466270577664053413)
	model.Property("Device", objectbox.PropertyType_String, 6, 3447591269737415871)
	model.Property("Name", objectbox.PropertyType_String, 7, 6909240474200120732)
	model.Property("Value", objectbox.PropertyType_String, 8, 7605577018169475074)
	model.EntityLastPropertyId(8, 7605577018169475074)
}

// GetId is called by the ObjectBox during Put operations to check for existing ID on an object
func (reading_EntityInfo) GetId(object interface{}) (uint64, error) {
	if obj, ok := object.(*Reading); ok {
		return objectbox.StringIdConvertToDatabaseValue(obj.Id), nil
	} else {
		return objectbox.StringIdConvertToDatabaseValue(object.(Reading).Id), nil
	}
}

// SetId is called by the ObjectBox during Put to update an ID on an object that has just been inserted
func (reading_EntityInfo) SetId(object interface{}, id uint64) error {
	if obj, ok := object.(*Reading); ok {
		obj.Id = objectbox.StringIdConvertToEntityProperty(id)
	} else {
		// NOTE while this can't update, it will at least behave consistently (panic in case of a wrong type)
		_ = object.(Reading).Id
	}
	return nil
}

// Flatten is called by the ObjectBox to transform an object to a FlatBuffer
func (reading_EntityInfo) Flatten(object interface{}, fbb *flatbuffers.Builder, id uint64) {
	obj := object.(*Reading)
	var offsetDevice = fbutils.CreateStringOffset(fbb, obj.Device)
	var offsetName = fbutils.CreateStringOffset(fbb, obj.Name)
	var offsetValue = fbutils.CreateStringOffset(fbb, obj.Value)

	// build the FlatBuffers object
	fbb.StartObject(8)
	fbutils.SetUint64Slot(fbb, 0, id)
	fbutils.SetInt64Slot(fbb, 1, obj.Pushed)
	fbutils.SetInt64Slot(fbb, 2, obj.Created)
	fbutils.SetInt64Slot(fbb, 3, obj.Origin)
	fbutils.SetInt64Slot(fbb, 4, obj.Modified)
	fbutils.SetUOffsetTSlot(fbb, 5, offsetDevice)
	fbutils.SetUOffsetTSlot(fbb, 6, offsetName)
	fbutils.SetUOffsetTSlot(fbb, 7, offsetValue)
}

// ToObject is called by the ObjectBox to load an object from a FlatBuffer
func (reading_EntityInfo) ToObject(bytes []byte) interface{} {
	table := &flatbuffers.Table{
		Bytes: bytes,
		Pos:   flatbuffers.GetUOffsetT(bytes),
	}

	return &Reading{
		Id:       objectbox.StringIdConvertToEntityProperty(table.GetUint64Slot(4, 0)),
		Pushed:   table.GetInt64Slot(6, 0),
		Created:  table.GetInt64Slot(8, 0),
		Origin:   table.GetInt64Slot(10, 0),
		Modified: table.GetInt64Slot(12, 0),
		Device:   fbutils.GetStringSlot(table, 14),
		Name:     fbutils.GetStringSlot(table, 16),
		Value:    fbutils.GetStringSlot(table, 18),
	}
}

// MakeSlice is called by the ObjectBox to construct a new slice to hold the read objects
func (reading_EntityInfo) MakeSlice(capacity int) interface{} {
	return make([]Reading, 0, capacity)
}

// AppendToSlice is called by the ObjectBox to fill the slice of the read objects
func (reading_EntityInfo) AppendToSlice(slice interface{}, object interface{}) interface{} {
	return append(slice.([]Reading), *object.(*Reading))
}

// Box provides CRUD access to Reading objects
type ReadingBox struct {
	*objectbox.Box
}

// BoxForReading opens a box of Reading objects
func BoxForReading(ob *objectbox.ObjectBox) *ReadingBox {
	return &ReadingBox{
		Box: ob.InternalBox(5),
	}
}

// Put synchronously inserts/updates a single object.
// In case the Id is not specified, it would be assigned automatically (auto-increment).
// When inserting, the Reading.Id property on the passed object will be assigned the new ID as well.
func (box *ReadingBox) Put(object *Reading) (uint64, error) {
	return box.Box.Put(object)
}

// PutAsync asynchronously inserts/updates a single object.
// When inserting, the Reading.Id property on the passed object will be assigned the new ID as well.
//
// It's executed on a separate internal thread for better performance.
//
// There are two main use cases:
//
// 1) "Put & Forget:" you gain faster puts as you don't have to wait for the transaction to finish.
//
// 2) Many small transactions: if your write load is typically a lot of individual puts that happen in parallel,
// this will merge small transactions into bigger ones. This results in a significant gain in overall throughput.
//
//
// In situations with (extremely) high async load, this method may be throttled (~1ms) or delayed (<1s).
// In the unlikely event that the object could not be enqueued after delaying, an error will be returned.
//
// Note that this method does not give you hard durability guarantees like the synchronous Put provides.
// There is a small time window (typically 3 ms) in which the data may not have been committed durably yet.
func (box *ReadingBox) PutAsync(object *Reading) (uint64, error) {
	return box.Box.PutAsync(object)
}

// PutAll inserts multiple objects in single transaction.
// In case Ids are not set on the objects, they would be assigned automatically (auto-increment).
//
// Returns: IDs of the put objects (in the same order).
// When inserting, the Reading.Id property on the objects in the slice will be assigned the new IDs as well.
//
// Note: In case an error occurs during the transaction, some of the objects may already have the Reading.Id assigned
// even though the transaction has been rolled back and the objects are not stored under those IDs.
//
// Note: The slice may be empty or even nil; in both cases, an empty IDs slice and no error is returned.
func (box *ReadingBox) PutAll(objects []Reading) ([]uint64, error) {
	return box.Box.PutAll(objects)
}

// Get reads a single object.
//
// Returns nil (and no error) in case the object with the given ID doesn't exist.
func (box *ReadingBox) Get(id uint64) (*Reading, error) {
	object, err := box.Box.Get(id)
	if err != nil {
		return nil, err
	} else if object == nil {
		return nil, nil
	}
	return object.(*Reading), nil
}

// Get reads all stored objects
func (box *ReadingBox) GetAll() ([]Reading, error) {
	objects, err := box.Box.GetAll()
	if err != nil {
		return nil, err
	}
	return objects.([]Reading), nil
}

// Remove deletes a single object
func (box *ReadingBox) Remove(object *Reading) (err error) {
	return box.Box.Remove(objectbox.StringIdConvertToDatabaseValue(object.Id))
}

// Creates a query with the given conditions. Use the fields of the Reading_ struct to create conditions.
// Keep the *ReadingQuery if you intend to execute the query multiple times.
// Note: this function panics if you try to create illegal queries; e.g. use properties of an alien type.
// This is typically a programming error. Use QueryOrError instead if you want the explicit error check.
func (box *ReadingBox) Query(conditions ...objectbox.Condition) *ReadingQuery {
	return &ReadingQuery{
		box.Box.Query(conditions...),
	}
}

// Creates a query with the given conditions. Use the fields of the Reading_ struct to create conditions.
// Keep the *ReadingQuery if you intend to execute the query multiple times.
func (box *ReadingBox) QueryOrError(conditions ...objectbox.Condition) (*ReadingQuery, error) {
	if query, err := box.Box.QueryOrError(conditions...); err != nil {
		return nil, err
	} else {
		return &ReadingQuery{query}, nil
	}
}

// Query provides a way to search stored objects
//
// For example, you can find all Reading which Id is either 42 or 47:
// 		box.Query(Reading_.Id.In(42, 47)).Find()
type ReadingQuery struct {
	*objectbox.Query
}

// Find returns all objects matching the query
func (query *ReadingQuery) Find() ([]Reading, error) {
	objects, err := query.Query.Find()
	if err != nil {
		return nil, err
	}
	return objects.([]Reading), nil
}

// Offset defines the index of the first object to process (how many objects to skip)
func (query *ReadingQuery) Offset(offset uint64) *ReadingQuery {
	query.Query.Offset(offset)
	return query
}

// Limit sets the number of elements to process by the query
func (query *ReadingQuery) Limit(limit uint64) *ReadingQuery {
	query.Query.Limit(limit)
	return query
}