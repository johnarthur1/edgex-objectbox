//
// Copyright (c) 2018 Cavium
//
// SPDX-License-Identifier: Apache-2.0
//

package test

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/edgexfoundry/edgex-go/internal/core/data/interfaces"
	dbp "github.com/edgexfoundry/edgex-go/internal/pkg/db"
	"github.com/edgexfoundry/edgex-go/pkg/models"
	"gopkg.in/mgo.v2/bson"
)

func populateDbReadings(db interfaces.DBClient, count int) (bson.ObjectId, error) {
	var id bson.ObjectId
	for i := 0; i < count; i++ {
		name := fmt.Sprintf("name%d", i)
		r := models.Reading{}
		r.Name = name
		r.Device = name
		r.Value = name
		var err error
		id, err = db.AddReading(r)
		if err != nil {
			return id, err
		}
	}
	return id, nil
}

func populateDbValues(db interfaces.DBClient, count int) (bson.ObjectId, error) {
	var id bson.ObjectId
	for i := 0; i < count; i++ {
		name := fmt.Sprintf("name%d", i)
		v := models.ValueDescriptor{}
		v.Name = name
		v.Description = name
		v.Type = name
		v.UomLabel = name
		v.Labels = []string{name, "LABEL"}
		var err error
		id, err = db.AddValueDescriptor(v)
		if err != nil {
			return id, err
		}
	}
	return id, nil
}

func populateDbEvents(db interfaces.DBClient, count int, pushed int64) (bson.ObjectId, error) {
	var id bson.ObjectId
	for i := 0; i < count; i++ {
		name := fmt.Sprintf("name%d", i)
		e := models.Event{}
		e.Device = name
		e.Event = name
		e.Pushed = pushed
		var err error
		id, err = db.AddEvent(&e)
		if err != nil {
			return id, err
		}
	}
	return id, nil
}

func testDBReadings(t *testing.T, db interfaces.DBClient) {
	err := db.ScrubAllEvents()
	if err != nil {
		t.Fatalf("Error removing all readings: %v\n", err)
	}

	readings, err := db.Readings()
	if err != nil {
		t.Fatalf("Error getting readings %v", err)
	}
	if readings == nil {
		t.Fatalf("Should return an empty array")
	}
	if len(readings) != 0 {
		t.Fatalf("There should be 0 readings instead of %d", len(readings))
	}

	beforeTime := dbp.MakeTimestamp()
	id, err := populateDbReadings(db, 100)
	if err != nil {
		t.Fatalf("Error populating db: %v\n", err)
	}

	// To have two readings with the same name
	id, err = populateDbReadings(db, 10)
	if err != nil {
		t.Fatalf("Error populating db: %v\n", err)
	}
	afterTime := dbp.MakeTimestamp()

	count, err := db.ReadingCount()
	if err != nil {
		t.Fatalf("Error getting readings count:  %v", err)
	}
	if count != 110 {
		t.Fatalf("There should be 110 readings instead of %d", count)
	}

	readings, err = db.Readings()
	if err != nil {
		t.Fatalf("Error getting readings %v", err)
	}
	if len(readings) != 110 {
		t.Fatalf("There should be 110 readings instead of %d", len(readings))
	}
	r3, err := db.ReadingById(id.Hex())
	if err != nil {
		t.Fatalf("Error getting reading by id %v", err)
	}
	if r3.Id.Hex() != id.Hex() {
		t.Fatalf("Id does not match %s - %s", r3.Id, id)
	}
	_, err = db.ReadingById("INVALID")
	if err == nil {
		t.Fatalf("Reading should not be found")
	}

	readings, err = db.ReadingsByDeviceAndValueDescriptor("name1", "name1", 10)
	if err != nil {
		t.Fatalf("Error getting ReadingsByDeviceAndValueDescriptor: %v", err)
	}
	if len(readings) != 2 {
		t.Fatalf("There should be 2 readings, not %d", len(readings))
	}
	readings, err = db.ReadingsByDeviceAndValueDescriptor("name1", "name1", 1)
	if err != nil {
		t.Fatalf("Error getting ReadingsByDeviceAndValueDescriptor: %v", err)
	}
	if len(readings) != 1 {
		t.Fatalf("There should be 1 readings, not %d", len(readings))
	}
	readings, err = db.ReadingsByDeviceAndValueDescriptor("name20", "name20", 10)
	if err != nil {
		t.Fatalf("Error getting ReadingsByDeviceAndValueDescriptor: %v", err)
	}
	if len(readings) != 1 {
		t.Fatalf("There should be 1 readings, not %d", len(readings))
	}

	readings, err = db.ReadingsByDevice("name1", 10)
	if err != nil {
		t.Fatalf("Error getting ReadingsByDevice: %v", err)
	}
	if len(readings) != 2 {
		t.Fatalf("There should be 2 readings, not %d", len(readings))
	}
	readings, err = db.ReadingsByDevice("name1", 1)
	if err != nil {
		t.Fatalf("Error getting ReadingsByDevice: %v", err)
	}
	if len(readings) != 1 {
		t.Fatalf("There should be 1 readings, not %d", len(readings))
	}
	readings, err = db.ReadingsByDevice("name20", 10)
	if err != nil {
		t.Fatalf("Error getting ReadingsByDevice: %v", err)
	}
	if len(readings) != 1 {
		t.Fatalf("There should be 1 readings, not %d", len(readings))
	}

	readings, err = db.ReadingsByValueDescriptor("name1", 10)
	if err != nil {
		t.Fatalf("Error getting ReadingsByValueDescriptor: %v", err)
	}
	if len(readings) != 2 {
		t.Fatalf("There should be 2 readings, not %d", len(readings))
	}
	readings, err = db.ReadingsByValueDescriptor("name1", 1)
	if err != nil {
		t.Fatalf("Error getting ReadingsByValueDescriptor: %v", err)
	}
	if len(readings) != 1 {
		t.Fatalf("There should be 1 readings, not %d", len(readings))
	}
	readings, err = db.ReadingsByValueDescriptor("name20", 10)
	if err != nil {
		t.Fatalf("Error getting ReadingsByValueDescriptor: %v", err)
	}
	if len(readings) != 1 {
		t.Fatalf("There should be 1 readings, not %d", len(readings))
	}
	readings, err = db.ReadingsByValueDescriptor("name", 10)
	if err != nil {
		t.Fatalf("Error getting ReadingsByValueDescriptor: %v", err)
	}
	if len(readings) != 0 {
		t.Fatalf("There should be 0 readings, not %d", len(readings))
	}

	readings, err = db.ReadingsByValueDescriptorNames([]string{"name1", "name2"}, 10)
	if err != nil {
		t.Fatalf("Error getting ReadingsByValueDescriptorNames: %v", err)
	}
	if len(readings) != 4 {
		t.Fatalf("There should be 4 readings, not %d", len(readings))
	}
	readings, err = db.ReadingsByValueDescriptorNames([]string{"name1", "name2"}, 1)
	if err != nil {
		t.Fatalf("Error getting ReadingsByValueDescriptorNames: %v", err)
	}
	if len(readings) != 1 {
		t.Fatalf("There should be 1 readings, not %d", len(readings))
	}
	readings, err = db.ReadingsByValueDescriptorNames([]string{"name", "noname"}, 10)
	if err != nil {
		t.Fatalf("Error getting ReadingsByValueDescriptorNames: %v", err)
	}
	if len(readings) != 0 {
		t.Fatalf("There should be 0 readings, not %d", len(readings))
	}
	readings, err = db.ReadingsByValueDescriptorNames([]string{"name20"}, 10)
	if err != nil {
		t.Fatalf("Error getting ReadingsByValueDescriptorNames: %v", err)
	}
	if len(readings) != 1 {
		t.Fatalf("There should be 1 readings, not %d", len(readings))
	}

	readings, err = db.ReadingsByCreationTime(beforeTime, afterTime, 200)
	if err != nil {
		t.Fatalf("Error getting ReadingsByCreationTime: %v", err)
	}
	if len(readings) != 110 {
		t.Fatalf("There should be 110 readings, not %d", len(readings))
	}
	readings, err = db.ReadingsByCreationTime(beforeTime, afterTime, 100)
	if err != nil {
		t.Fatalf("Error getting ReadingsByCreationTime: %v", err)
	}
	if len(readings) != 100 {
		t.Fatalf("There should be 100 readings, not %d", len(readings))
	}

	r := models.Reading{}
	r.Id = id
	r.Name = "name"
	err = db.UpdateReading(r)
	if err != nil {
		t.Fatalf("Error updating reading %v", err)
	}
	r2, err := db.ReadingById(r.Id.Hex())
	if err != nil {
		t.Fatalf("Error getting reading by id %v", err)
	}
	if r2.Name != r.Name {
		t.Fatalf("Did not update reading correctly: %s %s", r.Name, r2.Name)
	}

	err = db.DeleteReadingById("INVALID")
	if err == nil {
		t.Fatalf("Reading should not be deleted")
	}

	err = db.DeleteReadingById(id.Hex())
	if err != nil {
		t.Fatalf("Reading should be deleted: %v", err)
	}

	err = db.UpdateReading(r)
	if err == nil {
		t.Fatalf("Update should return error")
	}
}

func testDBEvents(t *testing.T, db interfaces.DBClient) {
	err := db.ScrubAllEvents()
	if err != nil {
		t.Fatalf("Error removing all events")
	}

	events, err := db.Events()
	if err != nil {
		t.Fatalf("Error getting events %v", err)
	}
	if events == nil {
		t.Fatalf("Should return an empty array")
	}
	if len(events) != 0 {
		t.Fatalf("There should be 0 events instead of %d", len(events))
	}

	beforeTime := dbp.MakeTimestamp()
	id, err := populateDbEvents(db, 100, 0)
	if err != nil {
		t.Fatalf("Error populating db: %v\n", err)
	}

	// To have two events with the same name
	id, err = populateDbEvents(db, 10, 1)
	if err != nil {
		t.Fatalf("Error populating db: %v\n", err)
	}
	afterTime := dbp.MakeTimestamp()

	count, err := db.EventCount()
	if err != nil {
		t.Fatalf("Error getting events count:  %v", err)
	}
	if count != 110 {
		t.Fatalf("There should be 110 events instead of %d", count)
	}

	count, err = db.EventCountByDeviceId("name1")
	if err != nil {
		t.Fatalf("Error getting events count:  %v", err)
	}
	if count != 2 {
		t.Fatalf("There should be 2 events instead of %d", count)
	}

	count, err = db.EventCountByDeviceId("name20")
	if err != nil {
		t.Fatalf("Error getting events count:  %v", err)
	}
	if count != 1 {
		t.Fatalf("There should be 1 events instead of %d", count)
	}

	count, err = db.EventCountByDeviceId("name")
	if err != nil {
		t.Fatalf("Error getting events count:  %v", err)
	}
	if count != 0 {
		t.Fatalf("There should be 0 events instead of %d", count)
	}

	events, err = db.Events()
	if err != nil {
		t.Fatalf("Error getting events %v", err)
	}
	if len(events) != 110 {
		t.Fatalf("There should be 110 events instead of %d", len(events))
	}
	e3, err := db.EventById(id.Hex())
	if err != nil {
		t.Fatalf("Error getting event by id %v", err)
	}
	if e3.ID.Hex() != id.Hex() {
		t.Fatalf("Id does not match %s - %s", e3.ID, id)
	}
	_, err = db.EventById("INVALID")
	if err == nil {
		t.Fatalf("Event should not be found")
	}

	events, err = db.EventsForDeviceLimit("name1", 10)
	if err != nil {
		t.Fatalf("Error getting EventsForDeviceLimit: %v", err)
	}
	if len(events) != 2 {
		t.Fatalf("There should be 2 events, not %d", len(events))
	}
	events, err = db.EventsForDeviceLimit("name1", 1)
	if err != nil {
		t.Fatalf("Error getting EventsForDeviceLimit: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("There should be 1 events, not %d", len(events))
	}
	events, err = db.EventsForDeviceLimit("name20", 10)
	if err != nil {
		t.Fatalf("Error getting EventsForDeviceLimit: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("There should be 1 events, not %d", len(events))
	}
	events, err = db.EventsForDeviceLimit("name", 10)
	if err != nil {
		t.Fatalf("Error getting EventsForDeviceLimit: %v", err)
	}
	if len(events) != 0 {
		t.Fatalf("There should be 0 events, not %d", len(events))
	}

	events, err = db.EventsForDevice("name1")
	if err != nil {
		t.Fatalf("Error getting EventsForDevice: %v", err)
	}
	if len(events) != 2 {
		t.Fatalf("There should be 2 events, not %d", len(events))
	}
	events, err = db.EventsForDevice("name20")
	if err != nil {
		t.Fatalf("Error getting EventsForDevice: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("There should be 1 events, not %d", len(events))
	}
	events, err = db.EventsForDevice("name")
	if err != nil {
		t.Fatalf("Error getting EventsForDevice: %v", err)
	}
	if len(events) != 0 {
		t.Fatalf("There should be 0 events, not %d", len(events))
	}

	events, err = db.EventsByCreationTime(beforeTime, afterTime, 200)
	if err != nil {
		t.Fatalf("Error getting EventsByCreationTime: %v", err)
	}
	if len(events) != 110 {
		t.Fatalf("There should be 110 events, not %d", len(events))
	}
	events, err = db.EventsByCreationTime(beforeTime, afterTime, 100)
	if err != nil {
		t.Fatalf("Error getting EventsByCreationTime: %v", err)
	}
	if len(events) != 100 {
		t.Fatalf("There should be 100 events, not %d", len(events))
	}

	events, err = db.EventsOlderThanAge(0)
	if err != nil {
		t.Fatalf("Error getting EventsOlderThanAge: %v", err)
	}
	if len(events) != 110 {
		t.Fatalf("There should be 110 events, not %d", len(events))
	}
	events, err = db.EventsOlderThanAge(1000000)
	if err != nil {
		t.Fatalf("Error getting EventsOlderThanAge: %v", err)
	}
	if len(events) != 0 {
		t.Fatalf("There should be 0 events, not %d", len(events))
	}

	events, err = db.EventsPushed()
	if err != nil {
		t.Fatalf("Error getting EventsOlderThanAge: %v", err)
	}
	if len(events) != 10 {
		t.Fatalf("There should be 10 events, not %d", len(events))
	}

	e := models.Event{}
	e.ID = id
	e.Device = "name"
	err = db.UpdateEvent(e)
	if err != nil {
		t.Fatalf("Error updating event %v", err)
	}
	e2, err := db.EventById(e.ID.Hex())
	if err != nil {
		t.Fatalf("Error getting event by id %v", err)
	}
	if e2.Device != e.Device {
		t.Fatalf("Did not update event correctly: %s %s", e.Device, e2.Device)
	}

	err = db.DeleteEventById("INVALID")
	if err == nil {
		t.Fatalf("Event should not be deleted")
	}

	err = db.DeleteEventById(id.Hex())
	if err != nil {
		t.Fatalf("Event should be deleted: %v", err)
	}

	err = db.UpdateEvent(e)
	if err == nil {
		t.Fatalf("Update should return error")
	}

	err = db.ScrubAllEvents()
	if err != nil {
		t.Fatalf("Error removing all events")
	}

	events, err = db.Events()
	if err != nil {
		t.Fatalf("Error getting events %v", err)
	}
	if len(events) != 0 {
		t.Fatalf("There should be 0 events instead of %d", len(events))
	}
}

func testDBValueDescriptors(t *testing.T, db interfaces.DBClient) {
	err := db.ScrubAllValueDescriptors()
	if err != nil {
		t.Fatalf("Error removing all value descriptors")
	}

	values, err := db.ValueDescriptors()
	if err != nil {
		t.Fatalf("Error getting events %v", err)
	}
	if values == nil {
		t.Fatalf("Should return an empty array")
	}
	if len(values) != 0 {
		t.Fatalf("There should be 0 values instead of %d", len(values))
	}

	id, err := populateDbValues(db, 110)
	if err != nil {
		t.Fatalf("Error populating db: %v\n", err)
	}

	_, err = populateDbValues(db, 110)
	if err == nil {
		t.Fatalf("Should be an error adding a new ValueDescriptor with the same name\n")
	}

	values, err = db.ValueDescriptors()
	if err != nil {
		t.Fatalf("Error getting Values %v", err)
	}
	if len(values) != 110 {
		t.Fatalf("There should be 110 Values instead of %d", len(values))
	}

	v3, err := db.ValueDescriptorById(id.Hex())
	if err != nil {
		t.Fatalf("Error getting Value by id %v", err)
	}
	if v3.Id.Hex() != id.Hex() {
		t.Fatalf("Id does not match %s - %s", v3.Id, id)
	}
	_, err = db.ValueDescriptorById("INVALID")
	if err == nil {
		t.Fatalf("Value should not be found")
	}

	v3, err = db.ValueDescriptorByName("name1")
	if err != nil {
		t.Fatalf("Error getting Value by id %v", err)
	}
	if v3.Name != "name1" {
		t.Fatalf("Name does not match %s - name1", v3.Name)
	}
	_, err = db.ValueDescriptorByName("INVALID")
	if err == nil {
		t.Fatalf("Value should not be found")
	}

	values, err = db.ValueDescriptorsByName([]string{"name1", "name2"})
	if err != nil {
		t.Fatalf("Error getting ValueDescriptorsByName: %v", err)
	}
	if len(values) != 2 {
		t.Fatalf("There should be 2 Values, not %d", len(values))
	}
	values, err = db.ValueDescriptorsByName([]string{"name1", "name"})
	if err != nil {
		t.Fatalf("Error getting ValueDescriptorsByName: %v", err)
	}
	if len(values) != 1 {
		t.Fatalf("There should be 1 Values, not %d", len(values))
	}
	values, err = db.ValueDescriptorsByName([]string{"name", "INVALID"})
	if err != nil {
		t.Fatalf("Error getting ValueDescriptorsByName: %v", err)
	}
	if len(values) != 0 {
		t.Fatalf("There should be 0 Values, not %d", len(values))
	}

	values, err = db.ValueDescriptorsByUomLabel("name1")
	if err != nil {
		t.Fatalf("Error getting ValueDescriptorsByUomLabel: %v", err)
	}
	if len(values) != 1 {
		t.Fatalf("There should be 1 Values, not %d", len(values))
	}
	values, err = db.ValueDescriptorsByUomLabel("INVALID")
	if err != nil {
		t.Fatalf("Error getting ValueDescriptorsByLabel: %v", err)
	}
	if len(values) != 0 {
		t.Fatalf("There should be 0 Values, not %d", len(values))
	}

	values, err = db.ValueDescriptorsByLabel("name1")
	if err != nil {
		t.Fatalf("Error getting ValueDescriptorsByLabel: %v", err)
	}
	if len(values) != 1 {
		t.Fatalf("There should be 1 Values, not %d", len(values))
	}
	values, err = db.ValueDescriptorsByLabel("INVALID")
	if err != nil {
		t.Fatalf("Error getting ValueDescriptorsByLabel: %v", err)
	}
	if len(values) != 0 {
		t.Fatalf("There should be 0 Values, not %d", len(values))
	}

	values, err = db.ValueDescriptorsByType("name1")
	if err != nil {
		t.Fatalf("Error getting ValueDescriptorsByType: %v", err)
	}
	if len(values) != 1 {
		t.Fatalf("There should be 1 Values, not %d", len(values))
	}
	values, err = db.ValueDescriptorsByType("INVALID")
	if err != nil {
		t.Fatalf("Error getting ValueDescriptorsByType: %v", err)
	}
	if len(values) != 0 {
		t.Fatalf("There should be 0 Values, not %d", len(values))
	}

	v := models.ValueDescriptor{}
	v.Id = id
	v.Name = "name"
	err = db.UpdateValueDescriptor(v)
	if err != nil {
		t.Fatalf("Error updating Value %v", err)
	}
	v2, err := db.ValueDescriptorById(v.Id.Hex())
	if err != nil {
		t.Fatalf("Error getting Value by id %v", err)
	}
	if v2.Name != v.Name {
		t.Fatalf("Did not update Value correctly: %s %s", v.Name, v2.Name)
	}

	err = db.DeleteValueDescriptorById("INVALID")
	if err == nil {
		t.Fatalf("Value should not be deleted")
	}

	err = db.DeleteValueDescriptorById(id.Hex())
	if err != nil {
		t.Fatalf("Value should be deleted: %v", err)
	}

	err = db.UpdateValueDescriptor(v)
	if err == nil {
		t.Fatalf("Update should return error")
	}

	err = db.ScrubAllValueDescriptors()
	if err != nil {
		t.Fatalf("Error removing all value descriptors")
	}
}

func TestDataDB(t *testing.T, db interfaces.DBClient) {
	err := db.Connect()
	if err != nil {
		t.Fatalf("Could not connect: %v", err)
	}

	testDBReadings(t, db)
	testDBEvents(t, db)
	testDBValueDescriptors(t, db)

	db.CloseSession()
	// Calling CloseSession twice to test that there is no panic when closing an
	// already closed db
	db.CloseSession()
}

func BenchmarkDB(b *testing.B, db interfaces.DBClient) {

	benchmarkReadings(b, db)
	benchmarkEvents(b, db)
	db.CloseSession()
}

func benchmarkReadings(b *testing.B, db interfaces.DBClient) {
	// Plain IDs do not require .hex(); must use reflect to avoid import cycle to identify DB
	plainIDs := strings.Contains(reflect.TypeOf(db).Name(), "ObjectBox")

	// Remove previous events and readings
	db.ScrubAllEvents()

	b.Run("AddReading", func(b *testing.B) {
		reading := models.Reading{}
		for i := 0; i < b.N; i++ {
			reading.Name = "test" + strconv.Itoa(i)
			reading.Device = "device" + strconv.Itoa(i/100)
			_, err := db.AddReading(reading)
			if err != nil {
				b.Fatalf("Error add reading: %v", err)
			}
		}
	})

	// Remove previous events and readings
	db.ScrubAllEvents()
	// prepare to benchmark n readings
	n := 1000
	readings := make([]string, n)
	reading := models.Reading{}
	for i := 0; i < n; i++ {
		reading.Name = "test" + strconv.Itoa(i)
		reading.Device = "device" + strconv.Itoa(i/100)
		id, err := db.AddReading(reading)
		if err != nil {
			b.Fatalf("Error add reading: %v", err)
		}
		if plainIDs {
			readings[i] = string(id)
		} else {
			readings[i] = id.Hex()
		}
	}

	b.Run("Readings", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := db.Readings()
			if err != nil {
				b.Fatalf("Error readings: %v", err)
			}
		}
	})

	b.Run("ReadingCount", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := db.ReadingCount()
			if err != nil {
				b.Fatalf("Error reading count: %v", err)
			}
		}
	})

	b.Run("ReadingById", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := db.ReadingById(readings[i%len(readings)])
			if err != nil {
				b.Fatalf("Error reading by ID: %v", err)
			}
		}
	})

	b.Run("ReadingsByDevice", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			device := "device" + strconv.Itoa((i%len(readings))/100)
			_, err := db.ReadingsByDevice(device, 100)
			if err != nil {
				b.Fatalf("Error reading by device: %v", err)
			}
		}
	})
}

func benchmarkEvents(b *testing.B, db interfaces.DBClient) {

	// Remove previous events and readings
	db.ScrubAllEvents()

	b.Run("AddEvent", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			device := fmt.Sprintf("device" + strconv.Itoa(i/100))
			e := models.Event{
				Device: device,
			}
			for j := 0; j < 5; j++ {
				r := models.Reading{
					Device: device,
					Name:   fmt.Sprintf("name%d", j),
				}
				e.Readings = append(e.Readings, r)
			}
			_, err := db.AddEvent(&e)
			if err != nil {
				b.Fatalf("Error add event: %v", err)
			}
		}
	})

	// Remove previous events and readings
	db.ScrubAllEvents()
	// prepare to benchmark n events (5 readings each)
	n := 1000
	events := make([]string, n)
	for i := 0; i < n; i++ {
		device := fmt.Sprintf("device" + strconv.Itoa(i/100))
		e := models.Event{
			Device: device,
		}
		for j := 0; j < 5; j++ {
			r := models.Reading{
				Device: device,
				Name:   fmt.Sprintf("name%d", j),
			}
			e.Readings = append(e.Readings, r)
		}
		id, err := db.AddEvent(&e)
		if err != nil {
			b.Fatalf("Error add event: %v", err)
		}
		events[i] = id.Hex()
	}

	b.Run("Events", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := db.Events()
			if err != nil {
				b.Fatalf("Error events: %v", err)
			}
		}
	})

	b.Run("EventCount", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := db.EventCount()
			if err != nil {
				b.Fatalf("Error event count: %v", err)
			}
		}
	})

	b.Run("EventById", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := db.EventById(events[i%len(events)])
			if err != nil {
				b.Fatalf("Error event by ID: %v", err)
			}
		}
	})

	b.Run("EventsForDevice", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			device := "device" + strconv.Itoa(i%len(events)/100)
			_, err := db.EventsForDevice(device)
			if err != nil {
				b.Fatalf("Error events for device: %v", err)
			}
		}
	})
}

func BenchmarkDBFixedN(db interfaces.DBClient, verify bool) {
	defer db.CloseSession()
	durable := true
	benchmarkReadingsN(db, verify, durable)
}

func benchmarkReadingsN(db interfaces.DBClient, verify bool, durable bool) {
	// Plain IDs do not require .hex(); must use reflect to avoid import cycle to identify DB
	dbType := reflect.TypeOf(db).String()
	println("\nBenchmarking " + dbType)
	println("---------------------------------------------")
	plainIDs := strings.Contains(dbType, "ObjectBox")

	// Remove any events and readings before and after test
	_ = db.ScrubAllEvents()
	defer db.ScrubAllEvents()

	count := 100000
	countPostfix := "[" + strconv.Itoa(count) + "]"
	readings := make([]string, count)
	RunBenchmarkN(db, "AddReading", count, func(ctx *BenchmarkContext) error {
		reading := models.Reading{}
		reading.Name = "test" + strconv.Itoa(ctx.I)
		reading.Device = "device" + strconv.Itoa(ctx.I/100)
		ctx.StartClock()
		id, err := db.AddReading(reading)
		if durable && ctx.I == count-1 {
			// Last one; ensure DBs actually made data durable
			db.EnsureAllDurable()
		}
		ctx.StopClock()
		if plainIDs {
			readings[ctx.I] = string(id)
		} else {
			readings[ctx.I] = id.Hex()
		}
		return err
	})

	RunBenchmarkN(db, "Readings"+countPostfix, 10, func(ctx *BenchmarkContext) error {
		readings, err := db.Readings()
		ctx.StopClock()
		size := len(readings)
		if verify && size != count {
			panic("Unexpected size: " + strconv.Itoa(size))
		}
		return err
	})

	RunBenchmarkN(db, "ReadingCount"+countPostfix, 100, func(ctx *BenchmarkContext) error {
		size, err := db.ReadingCount()
		ctx.StopClock()
		if verify && size != count {
			panic("Unexpected size: " + strconv.Itoa(size))
		}
		return err
	})

	RunBenchmarkN(db, "ReadingById", count, func(ctx *BenchmarkContext) error {
		id := readings[ctx.I]
		ctx.StartClock()
		reading, err := db.ReadingById(id)
		ctx.StopClock()

		if verify && ((plainIDs && string(reading.Id) != id) || (!plainIDs && reading.Id.Hex() != id)) {
			println(reading.String())
			panic("Expected ID " + id + " but got " + string(reading.Id))
		}

		return err
	})

	RunBenchmarkN(db, "ReadingsByDevice", 100, func(ctx *BenchmarkContext) error {
		device := "device" + strconv.Itoa(ctx.I)
		ctx.StartClock()
		slice, err := db.ReadingsByDevice(device, 100)
		ctx.StopClock()

		if verify {
			if len(slice) != 100 {
				panic("Unexpected slice size: " + strconv.Itoa(len(slice)))
			}

			for idx, reading := range slice {
				if reading.Device != device {
					println("[" + strconv.Itoa(idx) + "] " + reading.String())
					panic("Expected device " + device + " but got " + reading.Device)
				}
			}
		}
		return err
	})

}
