package objectbox

/*
#cgo LDFLAGS: -L ${SRCDIR}/libs -lobjectboxc
#include <stdlib.h>
#include <string.h>
#include "objectbox.h"
*/
import "C"

import (
	"errors"
	"fmt"
	"github.com/google/flatbuffers/go"
	"runtime"
	"strconv"
	"strings"
	"unsafe"
)

const Unavailable = flatbuffers.UOffsetT(0)

//noinspection GoUnusedConst
const (
	PropertyType_Bool       = 1
	PropertyType_Byte       = 2
	PropertyType_Short      = 3
	PropertyType_Char       = 4
	PropertyType_Int        = 5
	PropertyType_Long       = 6
	PropertyType_Float      = 7
	PropertyType_Double     = 8
	PropertyType_String     = 9
	PropertyType_Date       = 10
	PropertyType_Relation   = 11
	PropertyType_ByteVector = 23
)

//noinspection GoUnusedConst
const (
	/// One long property on an entity must be the ID
	PropertyFlags_ID = 1

	/// On languages like Java, a non-primitive type is used (aka wrapper types, allowing null)
	PropertyFlags_NON_PRIMITIVE_TYPE = 2

	/// Unused yet
	PropertyFlags_NOT_NULL = 4
	PropertyFlags_INDEXED  = 8
	PropertyFlags_RESERVED = 16
	/// Unused yet: Unique index
	PropertyFlags_UNIQUE = 32
	/// Unused yet: Use a persisted sequence to enforce ID to rise monotonic (no ID reuse)
	PropertyFlags_ID_MONOTONIC_SEQUENCE = 64
	/// Allow IDs to be assigned by the developer
	PropertyFlags_ID_SELF_ASSIGNABLE = 128
	/// Unused yet
	PropertyFlags_INDEX_PARTIAL_SKIP_NULL = 256
	/// Unused yet, used by References for 1) back-references and 2) to clear references to deleted objects (required for ID reuse)
	PropertyFlags_INDEX_PARTIAL_SKIP_ZERO = 512
	/// Virtual properties may not have a dedicated field in their entity class, e.g. target IDs of to-one relations
	PropertyFlags_VIRTUAL = 1024
	/// Index uses a 32 bit hash instead of the value
	/// (32 bits is shorter on disk, runs well on 32 bit systems, and should be OK even with a few collisions)

	PropertyFlags_INDEX_HASH = 2048
	/// Index uses a 64 bit hash instead of the value
	/// (recommended mostly for 64 bit machines with values longer >200 bytes; small values are faster with a 32 bit hash)
	PropertyFlags_INDEX_HASH64 = 4096
)

//noinspection GoUnusedConst
const (
	DebugFlags_LOG_TRANSACTIONS_READ  = 1
	DebugFlags_LOG_TRANSACTIONS_WRITE = 2
	DebugFlags_LOG_QUERIES            = 4
	DebugFlags_LOG_QUERY_PARAMETERS   = 8
	DebugFlags_LOG_ASYNC_QUEUE        = 16
)

type TypeId uint32

type ObjectBinding interface {
	GetTypeId() TypeId
	GetTypeName() string
	GetId(object interface{}) (id uint64, err error)
	Flatten(object interface{}, fbb *flatbuffers.Builder, id uint64)
	ToObject(bytes []byte) interface{}
	AppendToSlice(slice interface{}, object interface{}) (sliceNew interface{})
}

type Model struct {
	model *C.OB_model
	err   error
}

type ObjectBox struct {
	store          *C.OB_store
	bindingsById   map[TypeId]ObjectBinding
	bindingsByName map[string]ObjectBinding
}

type Transaction struct {
	txn       *C.OB_txn
	objectBox *ObjectBox
}

type Cursor struct {
	cursor  *C.OB_cursor
	binding ObjectBinding
	fbb     *flatbuffers.Builder
}

type Box struct {
	box     *C.OB_box
	binding ObjectBinding
	// FIXME not synchronized:
	fbb *flatbuffers.Builder
}

type TableArray struct {
	tableArray *C.OB_table_array
}

type BytesArray struct {
	bytesArray  [][]byte
	cBytesArray *C.OB_bytes_array
}

type TxnFun func(transaction *Transaction) (err error)
type CursorFun func(cursor *Cursor) (err error)

func NewModel() (model *Model, err error) {
	model = &Model{}
	model.model = C.ob_model_create()
	if model.model == nil {
		model = nil
		err = createError()
	}
	return
}

func (model *Model) LastEntityId(id uint32, uid uint64) {
	if model.err != nil {
		return
	}
	C.ob_model_last_entity_id(model.model, C.uint(id), C.ulong(uid))
}

func (model *Model) Entity(name string, id uint32, uid uint64) (err error) {
	if model.err != nil {
		return model.err
	}
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	rc := C.ob_model_entity(model.model, cname, C.uint(id), C.ulong(uid))
	if rc != 0 {
		err = createError()
	}
	return
}

func (model *Model) EntityLastPropertyId(id uint32, uid uint64) (err error) {
	if model.err != nil {
		return model.err
	}
	rc := C.ob_model_entity_last_property_id(model.model, C.uint(id), C.ulong(uid))
	if rc != 0 {
		err = createError()
	}
	return
}

func (model *Model) Property(name string, propertyType int, id uint32, uid uint64) (err error) {
	if model.err != nil {
		return model.err
	}
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	rc := C.ob_model_property(model.model, cname, C.OBPropertyType(propertyType), C.uint(id), C.ulong(uid))
	if rc != 0 {
		err = createError()
	}
	return
}

func (model *Model) PropertyFlags(propertyFlags int) (err error) {
	if model.err != nil {
		return model.err
	}
	rc := C.ob_model_property_flags(model.model, C.OBPropertyFlags(propertyFlags))
	if rc != 0 {
		err = createError()
	}
	return
}

func NewObjectBox(model *Model, name string) (objectBox *ObjectBox, err error) {
	fmt.Println("Ignoring name %v", name)
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	objectBox = &ObjectBox{}
	objectBox.store = C.ob_store_open(model.model, nil)
	if objectBox.store == nil {
		objectBox = nil
		err = createError()
	}
	if err == nil {
		objectBox.bindingsById = make(map[TypeId]ObjectBinding)
		objectBox.bindingsByName = make(map[string]ObjectBinding)
	}
	return
}

func (ob *ObjectBox) RegisterBinding(binding ObjectBinding) {
	id := binding.GetTypeId()
	name := binding.GetTypeName()
	existingBinding := ob.bindingsById[id]
	if existingBinding != nil {
		panic("Already registered a binding for ID " + strconv.Itoa(int(id)) + ": " + binding.GetTypeName())
	}
	existingBinding = ob.bindingsByName[name]
	if existingBinding != nil {
		panic("Already registered a binding for name " + name + ": ID " + strconv.Itoa(int(binding.GetTypeId())))
	}
	ob.bindingsById[id] = binding
	ob.bindingsByName[name] = binding
}

func (ob *ObjectBox) BeginTxn() (txn *Transaction, err error) {
	var ctxn = C.ob_txn_begin(ob.store)
	if ctxn == nil {
		return nil, createError()
	}
	return &Transaction{ctxn, ob}, nil
}

func (ob *ObjectBox) BeginTxnRead() (txn *Transaction, err error) {
	var ctxn = C.ob_txn_begin_read(ob.store)
	if ctxn == nil {
		return nil, createError()
	}
	return &Transaction{ctxn, ob}, nil
}

func (ob *ObjectBox) RunInTxn(readOnly bool, txnFun TxnFun) (err error) {
	runtime.LockOSThread()
	var txn *Transaction
	if readOnly {
		txn, err = ob.BeginTxnRead()
	} else {
		txn, err = ob.BeginTxn()
	}
	if err != nil {
		return
	}

	//fmt.Println(">>> START TX")
	//os.Stdout.Sync()

	err = txnFun(txn)

	//fmt.Println("<<< END TX")
	//os.Stdout.Sync()

	if !readOnly && err == nil {
		err = txn.Commit()
	}
	err2 := txn.Destroy()
	if err == nil {
		err = err2
	}
	runtime.UnlockOSThread()

	//fmt.Println("<<< END TX Destroy")
	//os.Stdout.Sync()

	return
}

func (ob ObjectBox) getBindingById(typeId TypeId) ObjectBinding {
	binding := ob.bindingsById[typeId]
	if binding == nil {
		// Configuration error by the dev, OK to panic
		panic("Configuration error; no binding registered for type ID " + strconv.Itoa(int(typeId)))
	}
	return binding
}

func (ob ObjectBox) getBindingByName(typeName string) ObjectBinding {
	binding := ob.bindingsByName[strings.ToLower(typeName)]
	if binding == nil {
		// Configuration error by the dev, OK to panic
		panic("Configuration error; no binding registered for type name " + typeName)
	}
	return binding
}

func (ob *ObjectBox) RunWithCursor(typeId TypeId, readOnly bool, cursorFun CursorFun) (err error) {
	binding := ob.getBindingById(typeId)
	return ob.RunInTxn(readOnly, func(txn *Transaction) (err error) {
		cursor, err := txn.Cursor(binding)
		if err != nil {
			return
		}
		//fmt.Println(">>> START C")
		//os.Stdout.Sync()

		err = cursorFun(cursor)

		//fmt.Println("<<< END C")
		//os.Stdout.Sync()

		err2 := cursor.Destroy()
		if err == nil {
			err = err2
		}
		return
	})
}

func (ob *ObjectBox) SetDebugFlags(flags uint) (err error) {
	rc := C.ob_store_debug_flags(ob.store, C.uint32_t(flags))
	if rc != 0 {
		err = createError()
	}
	return
}

func (ob *ObjectBox) Box(entitySchemaId TypeId) (*Box, error) {
	binding := ob.getBindingById(entitySchemaId)
	cbox := C.ob_box_create(ob.store, C.uint(entitySchemaId))
	if cbox == nil {
		return nil, createError()
	}
	return &Box{cbox, binding, flatbuffers.NewBuilder(512)}, nil
}

func (ob *ObjectBox) Strict() *ObjectBox {
	if C.ob_store_await_async_completion(ob.store) != 0 {
		fmt.Println(createError())
	}
	return ob
}

func (txn *Transaction) Destroy() (err error) {
	rc := C.ob_txn_destroy(txn.txn)
	txn.txn = nil
	if rc != 0 {
		err = createError()
	}
	return
}

func (txn *Transaction) Abort() (err error) {
	rc := C.ob_txn_abort(txn.txn)
	if rc != 0 {
		err = createError()
	}
	return
}

func (txn *Transaction) Commit() (err error) {
	rc := C.ob_txn_commit(txn.txn)
	if rc != 0 {
		err = createError()
	}
	return
}

func (txn *Transaction) Cursor(binding ObjectBinding) (*Cursor, error) {
	ccursor := C.ob_cursor_create(txn.txn, C.uint(binding.GetTypeId()))
	if ccursor == nil {
		return nil, createError()
	}
	return &Cursor{ccursor, binding, flatbuffers.NewBuilder(512)}, nil
}

func (txn *Transaction) CursorForName(entitySchemaName string) (*Cursor, error) {
	binding := txn.objectBox.getBindingByName(entitySchemaName)
	cname := C.CString(entitySchemaName)
	defer C.free(unsafe.Pointer(cname))

	ccursor := C.ob_cursor_create2(txn.txn, cname)
	if ccursor == nil {
		return nil, createError()
	}
	return &Cursor{ccursor, binding, flatbuffers.NewBuilder(512)}, nil
}

func (cursor *Cursor) Destroy() (err error) {
	rc := C.ob_cursor_destroy(cursor.cursor)
	cursor.cursor = nil
	if rc != 0 {
		err = createError()
	}
	return
}

func (cursor *Cursor) Get(id uint64) (object interface{}, err error) {
	bytes, err := cursor.GetBytes(id)
	if bytes == nil || err != nil {
		return
	}
	return cursor.binding.ToObject(bytes), nil
}

func (cursor *Cursor) GetAll() (slice interface{}, err error) {
	var bytes []byte
	binding := cursor.binding
	slice = nil

	for bytes, err = cursor.First(); bytes != nil; bytes, err = cursor.Next() {
		if err != nil || bytes == nil {
			slice = nil
			return
		}
		object := binding.ToObject(bytes)
		slice = binding.AppendToSlice(slice, object)
	}
	return slice, nil
}

func (cursor *Cursor) GetBytes(id uint64) (bytes []byte, err error) {
	var data *C.void
	var dataSize C.size_t
	dataPtr := unsafe.Pointer(data) // Need ptr to an unsafe ptr here
	rc := C.ob_cursor_get(cursor.cursor, C.uint64_t(id), &dataPtr, &dataSize)
	if rc != 0 {
		if rc != 404 {
			err = createError()
		}
		return
	}
	bytes = C.GoBytes(dataPtr, C.int(dataSize))
	return
}

func (cursor *Cursor) First() (bytes []byte, err error) {
	var data *C.void
	var dataSize C.size_t
	dataPtr := unsafe.Pointer(data) // Need ptr to an unsafe ptr here
	rc := C.ob_cursor_first(cursor.cursor, &dataPtr, &dataSize)
	if rc != 0 {
		if rc != 404 {
			err = createError()
		}
		return
	}
	bytes = C.GoBytes(dataPtr, C.int(dataSize))
	return
}

func (cursor *Cursor) Next() (bytes []byte, err error) {
	var data *C.void
	var dataSize C.size_t
	dataPtr := unsafe.Pointer(data) // Need ptr to an unsafe ptr here
	rc := C.ob_cursor_next(cursor.cursor, &dataPtr, &dataSize)
	if rc != 0 {
		if rc != 404 {
			err = createError()
		}
		return
	}
	bytes = C.GoBytes(dataPtr, C.int(dataSize))
	return
}

func (cursor *Cursor) Count() (count uint64, err error) {
	var cCount C.uint64_t
	rc := C.ob_cursor_count(cursor.cursor, &cCount)
	if rc != 0 {
		err = createError()
		return
	}
	return uint64(cCount), nil
}

func (cursor *Cursor) Put(object interface{}) (id uint64, err error) {
	idFromObject, err := cursor.binding.GetId(object)
	if err != nil {
		return
	}
	checkForPreviousValue := idFromObject != 0
	id, err = cursor.IdForPut(idFromObject)
	if err != nil {
		return
	}
	cursor.binding.Flatten(object, cursor.fbb, id)
	return id, cursor.finishInternalFbbAndPut(id, checkForPreviousValue)
}

func (cursor *Cursor) finishInternalFbbAndPut(id uint64, checkForPreviousObject bool) (err error) {
	fbb := cursor.fbb
	fbb.Finish(fbb.EndObject())
	bytes := fbb.FinishedBytes()

	cCheckPrevious := 0
	if checkForPreviousObject {
		cCheckPrevious = 1
	}
	rc := C.ob_cursor_put(cursor.cursor, C.uint64_t(id), unsafe.Pointer(&bytes[0]), C.size_t(len(bytes)),
		C.int(cCheckPrevious))
	if rc != 0 {
		err = createError()
	}

	// Reset to have a clear state for the next caller
	fbb.Reset()

	return
}

func (cursor *Cursor) IdForPut(idCandidate uint64) (id uint64, err error) {
	id = uint64(C.ob_cursor_id_for_put(cursor.cursor, C.uint64_t(idCandidate)))
	if id == 0 {
		err = createError()
	}
	return
}

func (cursor *Cursor) RemoveAll() (err error) {
	rc := C.ob_cursor_remove_all(cursor.cursor)
	if rc != 0 {
		err = createError()
	}
	return
}

func (cursor *Cursor) FindByString(propertyId uint, value string) (bytesArray *BytesArray, err error) {
	cvalue := C.CString(value)
	defer C.free(unsafe.Pointer(cvalue))

	cBytesArray := C.ob_query_by_string(cursor.cursor, C.uint32_t(propertyId), cvalue)
	if cBytesArray == nil {
		err = createError()
		return
	}
	size := int(cBytesArray.size)
	plainBytesArray := make([][]byte, size)
	if size > 0 {
		goBytesArray := (*[1 << 30]C.OB_bytes)(unsafe.Pointer(cBytesArray.bytes))[:size:size]
		for i := 0; i < size; i++ {
			cBytes := goBytesArray[i]
			dataBytes := C.GoBytes(cBytes.data, C.int(cBytes.size))
			plainBytesArray[i] = dataBytes
		}
	}

	return &BytesArray{plainBytesArray, cBytesArray}, nil
}

func (bytesArray *BytesArray) Destroy() {
	cBytesArray := bytesArray.cBytesArray
	if cBytesArray != nil {
		bytesArray.cBytesArray = nil
		C.ob_bytes_array_destroy(cBytesArray)
	}
	bytesArray.bytesArray = nil
}

func (box *Box) Destroy() (err error) {
	rc := C.ob_box_destroy(box.box)
	box.box = nil
	if rc != 0 {
		err = createError()
	}
	return
}

func (box *Box) IdForPut(idCandidate uint64) (id uint64, err error) {
	id = uint64(C.ob_box_id_for_put(box.box, C.uint64_t(idCandidate)))
	if id == 0 {
		err = createError()
	}
	return
}

func (box *Box) PutAsync(object interface{}) (id uint64, err error) {
	idFromObject, err := box.binding.GetId(object)
	if err != nil {
		return
	}
	checkForPreviousValue := idFromObject != 0
	id, err = box.IdForPut(idFromObject)
	if err != nil {
		return
	}
	box.binding.Flatten(object, box.fbb, id)
	return id, box.finishInternalFbbAndPutAsync(id, checkForPreviousValue)
}

func (box *Box) finishInternalFbbAndPutAsync(id uint64, checkForPreviousObject bool) (err error) {
	fbb := box.fbb
	fbb.Finish(fbb.EndObject())
	bytes := fbb.FinishedBytes()

	cCheckPrevious := 0
	if checkForPreviousObject {
		cCheckPrevious = 1
	}
	rc := C.ob_box_put_async(box.box, C.uint64_t(id), unsafe.Pointer(&bytes[0]), C.size_t(len(bytes)),
		C.int(cCheckPrevious))
	if rc != 0 {
		err = createError()
	}

	// Reset to have a clear state for the next caller
	fbb.Reset()

	return
}

func createError() error {
	msg := C.ob_last_error_message()
	return errors.New(C.GoString(msg))
}
