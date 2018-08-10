// automatically generated by the FlatBuffers compiler, do not modify

package flatcoredata

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type Event struct {
	_tab flatbuffers.Table
}

func GetRootAsEvent(buf []byte, offset flatbuffers.UOffsetT) *Event {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &Event{}
	x.Init(buf, n+offset)
	return x
}

func (rcv *Event) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *Event) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *Event) Id() uint64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.GetUint64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Event) MutateId(n uint64) bool {
	return rcv._tab.MutateUint64Slot(4, n)
}

func (rcv *Event) Pushed() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Event) MutatePushed(n int64) bool {
	return rcv._tab.MutateInt64Slot(6, n)
}

/// Device identifier (name or id)
func (rcv *Event) Device() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

/// Device identifier (name or id)
func (rcv *Event) Created() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Event) MutateCreated(n int64) bool {
	return rcv._tab.MutateInt64Slot(10, n)
}

func (rcv *Event) Modified() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Event) MutateModified(n int64) bool {
	return rcv._tab.MutateInt64Slot(12, n)
}

func (rcv *Event) Origin() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(14))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Event) MutateOrigin(n int64) bool {
	return rcv._tab.MutateInt64Slot(14, n)
}

/// Schedule event identifier
func (rcv *Event) ScheduleEvent() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(16))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

/// Schedule event identifier
func EventStart(builder *flatbuffers.Builder) {
	builder.StartObject(7)
}
func EventAddId(builder *flatbuffers.Builder, id uint64) {
	builder.PrependUint64Slot(0, id, 0)
}
func EventAddPushed(builder *flatbuffers.Builder, Pushed int64) {
	builder.PrependInt64Slot(1, Pushed, 0)
}
func EventAddDevice(builder *flatbuffers.Builder, Device flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(2, flatbuffers.UOffsetT(Device), 0)
}
func EventAddCreated(builder *flatbuffers.Builder, Created int64) {
	builder.PrependInt64Slot(3, Created, 0)
}
func EventAddModified(builder *flatbuffers.Builder, Modified int64) {
	builder.PrependInt64Slot(4, Modified, 0)
}
func EventAddOrigin(builder *flatbuffers.Builder, Origin int64) {
	builder.PrependInt64Slot(5, Origin, 0)
}
func EventAddScheduleEvent(builder *flatbuffers.Builder, ScheduleEvent flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(6, flatbuffers.UOffsetT(ScheduleEvent), 0)
}
func EventEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
