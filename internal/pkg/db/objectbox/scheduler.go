package objectbox

// implements export-client service contract

import (
	contract "github.com/edgexfoundry/go-mod-core-contracts/models"
	"github.com/objectbox/edgex-objectbox/internal/pkg/db"
	"github.com/objectbox/edgex-objectbox/internal/pkg/db/objectbox/obx"
	"github.com/objectbox/objectbox-go/objectbox"
	"sync"
)

type schedulerClient struct {
	objectBox *objectbox.ObjectBox

	intervalBox       *obx.IntervalBox       // no async - has unique and requires insert/update to fail
	intervalActionBox *obx.IntervalActionBox // no async - has unique and requires insert/update to fail

	queries schedulerQueries
}

//region Queries
type schedulerQueries struct {
	interval struct {
		all  intervalQuery
		name intervalQuery
	}
	intervalAction struct {
		all      intervalActionQuery
		interval intervalActionQuery
		name     intervalActionQuery
		target   intervalActionQuery
	}
}

type intervalQuery struct {
	*obx.IntervalQuery
	sync.Mutex
}

type intervalActionQuery struct {
	*obx.IntervalActionQuery
	sync.Mutex
}

//endregion

func newSchedulerClient(objectBox *objectbox.ObjectBox) (*schedulerClient, error) {
	var client = &schedulerClient{objectBox: objectBox}
	var err error

	client.intervalBox = obx.BoxForInterval(objectBox)
	client.intervalActionBox = obx.BoxForIntervalAction(objectBox)

	//region Interval
	if err == nil {
		client.queries.interval.all.IntervalQuery, err = client.intervalBox.QueryOrError()
	}
	if err == nil {
		client.queries.interval.name.IntervalQuery, err =
			client.intervalBox.QueryOrError(obx.Interval_.Name.Equals("", true))
	}
	//endregion

	//region IntervalAction
	if err == nil {
		client.queries.intervalAction.all.IntervalActionQuery, err = client.intervalActionBox.QueryOrError()
	}
	if err == nil {
		client.queries.intervalAction.interval.IntervalActionQuery, err =
			client.intervalActionBox.QueryOrError(obx.IntervalAction_.Interval.Equals("", true))
	}
	if err == nil {
		client.queries.intervalAction.name.IntervalActionQuery, err =
			client.intervalActionBox.QueryOrError(obx.IntervalAction_.Name.Equals("", true))
	}
	if err == nil {
		client.queries.intervalAction.target.IntervalActionQuery, err =
			client.intervalActionBox.QueryOrError(obx.IntervalAction_.Target.Equals("", true))
	}
	//endregion

	if err == nil {
		return client, nil
	} else {
		return nil, mapError(err)
	}
}

func (client *schedulerClient) Intervals() ([]contract.Interval, error) {
	result, err := client.intervalBox.GetAll()
	return result, mapError(err)
}

func (client *schedulerClient) IntervalsWithLimit(limit int) ([]contract.Interval, error) {
	var query = &client.queries.interval.all

	query.Lock()
	defer query.Unlock()

	result, err := query.Limit(uint64(limit)).Find()
	return result, mapError(err)
}

func (client *schedulerClient) IntervalByName(name string) (contract.Interval, error) {
	var query = &client.queries.interval.name

	query.Lock()
	defer query.Unlock()

	if err := query.SetStringParams(obx.Interval_.Name, name); err != nil {
		return contract.Interval{}, mapError(err)
	}

	if list, err := query.Limit(1).Find(); err != nil {
		return contract.Interval{}, mapError(err)
	} else if len(list) == 0 {
		return contract.Interval{}, mapError(db.ErrNotFound)
	} else {
		return list[0], nil
	}
}

func (client *schedulerClient) IntervalById(id string) (contract.Interval, error) {
	if id, err := obx.IdFromString(id); err != nil {
		return contract.Interval{}, mapError(err)
	} else if object, err := client.intervalBox.Get(id); err != nil {
		return contract.Interval{}, mapError(err)
	} else if object == nil {
		return contract.Interval{}, mapError(db.ErrNotFound)
	} else {
		return *object, nil
	}
}

func (client *schedulerClient) AddInterval(interval contract.Interval) (string, error) {
	onCreate(&interval.Timestamps)

	id, err := client.intervalBox.Put(&interval)
	return obx.IdToString(id), mapError(err)
}

func (client *schedulerClient) UpdateInterval(interval contract.Interval) error {
	onUpdate(&interval.Timestamps)

	if id, err := obx.IdFromString(interval.ID); err != nil {
		return mapError(err)
	} else if exists, err := client.intervalBox.Contains(id); err != nil {
		return mapError(err)
	} else if !exists {
		return mapError(db.ErrNotFound)
	}

	_, err := client.intervalBox.Put(&interval)
	return mapError(err)
}

func (client *schedulerClient) DeleteIntervalById(id string) error {
	if id, err := obx.IdFromString(id); err != nil {
		return mapError(err)
	} else {
		return mapError(client.intervalBox.RemoveId(id))
	}
}

func (client *schedulerClient) IntervalActions() ([]contract.IntervalAction, error) {
	result, err := client.intervalActionBox.GetAll()
	return result, mapError(err)
}

func (client *schedulerClient) IntervalActionsWithLimit(limit int) ([]contract.IntervalAction, error) {
	var query = &client.queries.intervalAction.all

	query.Lock()
	defer query.Unlock()

	result, err := query.Limit(uint64(limit)).Find()
	return result, mapError(err)
}

func (client *schedulerClient) IntervalActionsByIntervalName(name string) ([]contract.IntervalAction, error) {
	var query = &client.queries.intervalAction.interval

	query.Lock()
	defer query.Unlock()

	if err := query.SetStringParams(obx.IntervalAction_.Interval, name); err != nil {
		return nil, mapError(err)
	}

	result, err := query.Limit(0).Find()
	return result, mapError(err)
}

func (client *schedulerClient) IntervalActionsByTarget(name string) ([]contract.IntervalAction, error) {
	var query = &client.queries.intervalAction.target

	query.Lock()
	defer query.Unlock()

	if err := query.SetStringParams(obx.IntervalAction_.Target, name); err != nil {
		return nil, mapError(err)
	}

	result, err := query.Limit(0).Find()
	return result, mapError(err)
}

func (client *schedulerClient) IntervalActionById(id string) (contract.IntervalAction, error) {
	if id, err := obx.IdFromString(id); err != nil {
		return contract.IntervalAction{}, mapError(err)
	} else if object, err := client.intervalActionBox.Get(id); err != nil {
		return contract.IntervalAction{}, mapError(err)
	} else if object == nil {
		return contract.IntervalAction{}, mapError(db.ErrNotFound)
	} else {
		return *object, nil
	}
}

func (client *schedulerClient) IntervalActionByName(name string) (contract.IntervalAction, error) {
	var query = &client.queries.intervalAction.name

	query.Lock()
	defer query.Unlock()

	if err := query.SetStringParams(obx.IntervalAction_.Name, name); err != nil {
		return contract.IntervalAction{}, mapError(err)
	}

	if list, err := query.Limit(1).Find(); err != nil {
		return contract.IntervalAction{}, mapError(err)
	} else if len(list) == 0 {
		return contract.IntervalAction{}, mapError(db.ErrNotFound)
	} else {
		return list[0], nil
	}
}

func (client *schedulerClient) AddIntervalAction(intervalAction contract.IntervalAction) (string, error) {
	// NOTE this is done instead of onCreate because there is no reg.Timestamps
	if intervalAction.Created == 0 {
		intervalAction.Created = db.MakeTimestamp()
	}

	id, err := client.intervalActionBox.Put(&intervalAction)
	return obx.IdToString(id), mapError(err)
}

func (client *schedulerClient) UpdateIntervalAction(intervalAction contract.IntervalAction) error {
	// NOTE this is done instead of onUpdate because there is no reg.Timestamps
	intervalAction.Modified = db.MakeTimestamp()

	if id, err := obx.IdFromString(intervalAction.ID); err != nil {
		return mapError(err)
	} else if exists, err := client.intervalActionBox.Contains(id); err != nil {
		return mapError(err)
	} else if !exists {
		return mapError(db.ErrNotFound)
	}

	_, err := client.intervalActionBox.Put(&intervalAction)
	return mapError(err)
}

func (client *schedulerClient) DeleteIntervalActionById(id string) error {
	if id, err := obx.IdFromString(id); err != nil {
		return mapError(err)
	} else {
		return mapError(client.intervalActionBox.RemoveId(id))
	}
}

func (client *schedulerClient) ScrubAllIntervalActions() (int, error) {
	var query = &client.queries.intervalAction.all

	query.Lock()
	defer query.Unlock()

	count, err := query.Limit(0).Remove()
	return int(count), mapError(err)
}

func (client *schedulerClient) ScrubAllIntervals() (int, error) {
	var query = &client.queries.interval.all

	query.Lock()
	defer query.Unlock()

	count, err := query.Limit(0).Remove()
	return int(count), mapError(err)
}
