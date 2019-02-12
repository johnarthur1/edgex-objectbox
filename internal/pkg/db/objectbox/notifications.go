package objectbox

// implements core-data service contract

import (
	"github.com/edgexfoundry/edgex-go/internal/pkg/db"
	"github.com/edgexfoundry/edgex-go/internal/pkg/db/objectbox/obx"
	contract "github.com/edgexfoundry/edgex-go/pkg/models"
	"github.com/objectbox/objectbox-go/objectbox"
	"sync"
)

type notificationsClient struct {
	objectBox *objectbox.ObjectBox

	notificationBox *obx.NotificationBox
	subscriptionBox *obx.SubscriptionBox
	transmissionBox *obx.TransmissionBox

	queries notificationsQueries

	cleanupDefaultAge int
}

//region Queries
type notificationsQueries struct {
	notification struct {
		createdB            notificationQuery
		createdGT           notificationQuery
		createdLT           notificationQuery
		labels              notificationQuery
		modifiedLT          notificationQuery
		sender              notificationQuery
		slug                notificationQuery
		status              notificationQuery
		statusAndModifiedLT notificationQuery
		statusAndSeverity   notificationQuery
	}
	subscription struct {
		categories          subscriptionQuery
		categoriesAndLabels subscriptionQuery
		labels              subscriptionQuery
		receiver            subscriptionQuery
		slug                subscriptionQuery
	}
	transmission struct {
		createdBAndResendCount     transmissionQuery
		createdGTAndResendCount    transmissionQuery
		createdLTAndResendCount    transmissionQuery
		notification               transmissionQuery
		notificationAndResendCount transmissionQuery
		statusAndModifiedLT        transmissionQuery
		statusAndResendCount       transmissionQuery
	}
}

type notificationQuery struct {
	*obx.NotificationQuery
	sync.Mutex
}

type subscriptionQuery struct {
	*obx.SubscriptionQuery
	sync.Mutex
}

type transmissionQuery struct {
	*obx.TransmissionQuery
	sync.Mutex
}

//endregion

func newNotificationsClient(objectBox *objectbox.ObjectBox) (*notificationsClient, error) {
	var client = &notificationsClient{objectBox: objectBox}
	var err error

	// TODO check this, it's picked from ServicesTest
	client.cleanupDefaultAge = 86400001

	client.notificationBox = obx.BoxForNotification(client.objectBox)
	client.subscriptionBox = obx.BoxForSubscription(client.objectBox)
	client.transmissionBox = obx.BoxForTransmission(client.objectBox)

	//region Notification
	if err == nil {
		client.queries.notification.createdB.NotificationQuery, err =
			client.notificationBox.QueryOrError(obx.Notification_.Created.Between(0, 0))
	}

	if err == nil {
		client.queries.notification.createdGT.NotificationQuery, err =
			client.notificationBox.QueryOrError(obx.Notification_.Created.GreaterThan(0))
	}

	if err == nil {
		client.queries.notification.createdLT.NotificationQuery, err =
			client.notificationBox.QueryOrError(obx.Notification_.Created.LessThan(0))
	}

	if err == nil {
		client.queries.notification.labels.NotificationQuery, err =
			client.notificationBox.QueryOrError(obx.Notification_.Labels.Contains("", true))
	}

	if err == nil {
		client.queries.notification.modifiedLT.NotificationQuery, err =
			client.notificationBox.QueryOrError(obx.Notification_.Modified.LessThan(0))
	}

	if err == nil {
		client.queries.notification.sender.NotificationQuery, err =
			client.notificationBox.QueryOrError(obx.Notification_.Sender.Equals("", true))
	}

	if err == nil {
		client.queries.notification.slug.NotificationQuery, err =
			client.notificationBox.QueryOrError(obx.Notification_.Slug.Equals("", true))
	}

	if err == nil {
		client.queries.notification.status.NotificationQuery, err =
			client.notificationBox.QueryOrError(obx.Notification_.Status.Equals("", true))
	}

	if err == nil {
		client.queries.notification.statusAndModifiedLT.NotificationQuery, err = client.notificationBox.QueryOrError(
			obx.Notification_.Status.Equals("", true), obx.Notification_.Modified.LessThan(0))
	}

	if err == nil {
		client.queries.notification.statusAndSeverity.NotificationQuery, err = client.notificationBox.QueryOrError(
			obx.Notification_.Status.Equals("", true),
			obx.Notification_.Severity.Equals("", true))
	}
	//endregion

	//region Subscription
	//if err == nil {
	//	client.queries.subscription.categories.SubscriptionQuery, err =
	//		client.subscriptionBox.QueryOrError(obx.Subscription_.SubscribedCategories.Contains("", true))
	//}
	//
	//if err == nil {
	//	client.queries.subscription.categoriesAndLabels.SubscriptionQuery, err =
	//		client.subscriptionBox.QueryOrError(obx.Reading_.Device.Equals("", true), obx.Reading_.Name.Equals("", true))
	//}

	if err == nil {
		client.queries.subscription.labels.SubscriptionQuery, err =
			client.subscriptionBox.QueryOrError(obx.Subscription_.SubscribedLabels.Contains("", true))
	}

	if err == nil {
		client.queries.subscription.slug.SubscriptionQuery, err =
			client.subscriptionBox.QueryOrError(obx.Subscription_.Slug.Equals("", true))
	}

	if err == nil {
		client.queries.subscription.receiver.SubscriptionQuery, err =
			client.subscriptionBox.QueryOrError(obx.Subscription_.Receiver.Equals("", true))
	}
	//endregion

	//region Transmission
	if err == nil {
		client.queries.transmission.createdBAndResendCount.TransmissionQuery, err = client.transmissionBox.QueryOrError(
			obx.Transmission_.ResendCount.LessThan(0), obx.Transmission_.Created.Between(0, 0))
	}

	if err == nil {
		client.queries.transmission.createdGTAndResendCount.TransmissionQuery, err = client.transmissionBox.QueryOrError(
			obx.Transmission_.ResendCount.LessThan(0), obx.Transmission_.Created.GreaterThan(0))
	}

	if err == nil {
		client.queries.transmission.createdLTAndResendCount.TransmissionQuery, err = client.transmissionBox.QueryOrError(
			obx.Transmission_.ResendCount.LessThan(0), obx.Transmission_.Created.GreaterThan(0))
	}

	if err == nil {
		client.queries.transmission.notification.TransmissionQuery, err =
			client.transmissionBox.QueryOrError(obx.Transmission_.Notification.In())
	}

	if err == nil {
		client.queries.transmission.notificationAndResendCount.TransmissionQuery, err =
			client.transmissionBox.QueryOrError(obx.Transmission_.ResendCount.LessThan(0), obx.Transmission_.Notification.In())
	}

	if err == nil {
		client.queries.transmission.statusAndResendCount.TransmissionQuery, err = client.transmissionBox.QueryOrError(
			obx.Transmission_.ResendCount.LessThan(0), obx.Transmission_.Status.Equals("", true))
	}

	if err == nil {
		client.queries.transmission.statusAndModifiedLT.TransmissionQuery, err = client.transmissionBox.QueryOrError(
			obx.Transmission_.Status.Equals("", true), obx.Transmission_.Modified.LessThan(0))
	}
	//endregion

	if err == nil {
		return client, nil
	} else {
		return nil, err
	}
}

func (client *notificationsClient) GetNotifications() ([]contract.Notification, error) {
	return client.notificationBox.GetAll()
}

func (client *notificationsClient) GetNotificationById(id string) (contract.Notification, error) {
	if id, err := obx.IdFromString(id); err != nil {
		return contract.Notification{}, err
	} else if object, err := client.notificationBox.Get(id); err != nil {
		return contract.Notification{}, err
	} else if object == nil {
		return contract.Notification{}, db.ErrNotFound
	} else {
		return *object, nil
	}
}

func (client *notificationsClient) GetNotificationBySlug(slug string) (contract.Notification, error) {
	var query = &client.queries.notification.slug

	query.Lock()
	defer query.Unlock()

	if err := query.SetStringParams(obx.Notification_.Slug, slug); err != nil {
		return contract.Notification{}, err
	}

	if list, err := query.Limit(1).Find(); err != nil {
		return contract.Notification{}, err
	} else if len(list) == 0 {
		return contract.Notification{}, db.ErrNotFound
	} else {
		return list[0], nil
	}
}

func (client *notificationsClient) GetNotificationBySender(sender string, limit int) ([]contract.Notification, error) {
	var query = &client.queries.notification.sender

	query.Lock()
	defer query.Unlock()

	if err := query.SetStringParams(obx.Notification_.Sender, sender); err != nil {
		return nil, err
	}

	return query.Limit(uint64(limit)).Find()
}

func (client *notificationsClient) GetNotificationsByLabels(labels []string, limit int) ([]contract.Notification, error) {
	if query, err := client.notificationBox.QueryOrError(
		stringVectorContainsAny(obx.Notification_.Labels, labels, true)); err != nil {
		return nil, err
	} else {
		return query.Limit(uint64(limit)).Find()
	}
}

func (client *notificationsClient) GetNotificationsByStartEnd(start int64, end int64, limit int) ([]contract.Notification, error) {
	var query = &client.queries.notification.createdB

	query.Lock()
	defer query.Unlock()

	if err := query.SetInt64Params(obx.Notification_.Created, start, end); err != nil {
		return nil, err
	}

	return query.Limit(uint64(limit)).Find()
}

func (client *notificationsClient) GetNotificationsByStart(start int64, limit int) ([]contract.Notification, error) {
	var query = &client.queries.notification.createdGT

	query.Lock()
	defer query.Unlock()

	if err := query.SetInt64Params(obx.Notification_.Created, start); err != nil {
		return nil, err
	}

	return query.Limit(uint64(limit)).Find()
}

func (client *notificationsClient) GetNotificationsByEnd(end int64, limit int) ([]contract.Notification, error) {
	var query = &client.queries.notification.createdLT

	query.Lock()
	defer query.Unlock()

	if err := query.SetInt64Params(obx.Notification_.Created, end); err != nil {
		return nil, err
	}
	return query.Limit(uint64(limit)).Find()
}

func (client *notificationsClient) GetNewNotifications(limit int) ([]contract.Notification, error) {
	var query = &client.queries.notification.status

	query.Lock()
	defer query.Unlock()

	if err := query.SetStringParams(obx.Notification_.Status, "NEW"); err != nil {
		return nil, err
	}

	return query.Limit(uint64(limit)).Find()
}

func (client *notificationsClient) GetNewNormalNotifications(limit int) ([]contract.Notification, error) {
	var query = &client.queries.notification.statusAndSeverity

	query.Lock()
	defer query.Unlock()

	if err := query.SetStringParams(obx.Notification_.Status, "NEW"); err != nil {
		return nil, err
	}

	if err := query.SetStringParams(obx.Notification_.Severity, "NORMAL"); err != nil {
		return nil, err
	}

	return query.Limit(uint64(limit)).Find()
}

func (client *notificationsClient) AddNotification(n contract.Notification) (string, error) {
	onCreate(&n.BaseObject)

	id, err := client.notificationBox.Put(&n)
	return obx.IdToString(id), err
}

func (client *notificationsClient) UpdateNotification(n contract.Notification) error {
	onUpdate(&n.BaseObject)

	if id, err := obx.IdFromString(n.ID); err != nil {
		return err
	} else if exists, err := client.notificationBox.Contains(id); err != nil {
		return err
	} else if !exists {
		return db.ErrNotFound
	}

	_, err := client.notificationBox.Put(&n)
	return err
}

func (client *notificationsClient) MarkNotificationProcessed(n contract.Notification) error {
	n.Status = contract.NotificationsStatus(contract.Processed)
	return client.UpdateNotification(n)
}

func (client *notificationsClient) DeleteNotificationById(id string) error {
	if id, err := obx.IdFromString(id); err != nil {
		return err
	} else {
		return client.notificationBox.Box.Remove(id)
	}
}

func (client *notificationsClient) DeleteNotificationBySlug(slug string) error {
	if obj, err := client.GetNotificationBySlug(slug); err != nil {
		return err
	} else {
		return client.notificationBox.Remove(&obj)
	}
}

func (client *notificationsClient) DeleteNotificationsOld(age int) error {
	var query = &client.queries.notification.statusAndModifiedLT

	query.Lock()
	defer query.Unlock()

	if err := query.SetStringParams(obx.Notification_.Status, "PROCESSED"); err != nil {
		return err
	}

	var end = db.MakeTimestamp() - int64(age)
	if err := query.SetInt64Params(obx.Notification_.Modified, end); err != nil {
		return err
	}

	_, err := query.Remove()
	return err
}

func (client *notificationsClient) GetSubscriptions() ([]contract.Subscription, error) {
	return client.subscriptionBox.GetAll()
}

func (client *notificationsClient) GetSubscriptionById(id string) (contract.Subscription, error) {
	if id, err := obx.IdFromString(id); err != nil {
		return contract.Subscription{}, err
	} else if object, err := client.subscriptionBox.Get(id); err != nil {
		return contract.Subscription{}, err
	} else if object == nil {
		return contract.Subscription{}, db.ErrNotFound
	} else {
		return *object, nil
	}
}

func (client *notificationsClient) GetSubscriptionBySlug(slug string) (contract.Subscription, error) {
	var query = &client.queries.subscription.slug

	query.Lock()
	defer query.Unlock()

	if err := query.SetStringParams(obx.Subscription_.Slug, slug); err != nil {
		return contract.Subscription{}, err
	}

	if list, err := query.Limit(1).Find(); err != nil {
		return contract.Subscription{}, err
	} else if len(list) == 0 {
		return contract.Subscription{}, db.ErrNotFound
	} else {
		return list[0], nil
	}
}

func (client *notificationsClient) GetSubscriptionByReceiver(receiver string) ([]contract.Subscription, error) {
	var query = &client.queries.subscription.receiver

	query.Lock()
	defer query.Unlock()

	if err := query.SetStringParams(obx.Subscription_.Receiver, receiver); err != nil {
		return nil, err
	}

	return query.Find()
}

func (client *notificationsClient) GetSubscriptionByCategories(categories []string) ([]contract.Subscription, error) {
	if query, err := client.subscriptionBox.QueryOrError(
		stringVectorContainsAny(obx.Subscription_.SubscribedCategories, categories, true)); err != nil {
		return nil, err
	} else {
		return query.Find()
	}
}

func (client *notificationsClient) GetSubscriptionByLabels(labels []string) ([]contract.Subscription, error) {
	if query, err := client.subscriptionBox.QueryOrError(
		stringVectorContainsAny(obx.Subscription_.SubscribedLabels, labels, true)); err != nil {
		return nil, err
	} else {
		return query.Find()
	}
}

func (client *notificationsClient) GetSubscriptionByCategoriesLabels(categories []string, labels []string) ([]contract.Subscription, error) {
	if query, err := client.subscriptionBox.QueryOrError(
		stringVectorContainsAny(obx.Subscription_.SubscribedCategories, categories, true),
		stringVectorContainsAny(obx.Subscription_.SubscribedLabels, labels, true)); err != nil {
		return nil, err
	} else {
		return query.Find()
	}
}

func (client *notificationsClient) AddSubscription(s contract.Subscription) (string, error) {
	onCreate(&s.BaseObject)

	id, err := client.subscriptionBox.Put(&s)
	return obx.IdToString(id), err
}

func (client *notificationsClient) UpdateSubscription(s contract.Subscription) error {
	onUpdate(&s.BaseObject)

	if id, err := obx.IdFromString(s.ID); err != nil {
		return err
	} else if exists, err := client.subscriptionBox.Contains(id); err != nil {
		return err
	} else if !exists {
		return db.ErrNotFound
	}

	_, err := client.subscriptionBox.Put(&s)
	return err
}

func (client *notificationsClient) DeleteSubscriptionBySlug(slug string) error {
	if obj, err := client.GetSubscriptionBySlug(slug); err != nil {
		return err
	} else {
		return client.subscriptionBox.Remove(&obj)
	}
}

func (client *notificationsClient) GetTransmissionsByNotificationSlug(slug string, resendLimit int) ([]contract.Transmission, error) {
	var query = &client.queries.notification.slug

	query.Lock()
	defer query.Unlock()

	if err := query.SetStringParams(obx.Notification_.Slug, slug); err != nil {
		return nil, err
	}

	if ids, err := query.FindIds(); err != nil {
		return nil, err
	} else {
		return client.getTransmissionsByNotificationIds(resendLimit, ids)
	}
}

func (client *notificationsClient) getTransmissionsByNotificationIds(resendLimit int, ids []uint64) ([]contract.Transmission, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	var query = &client.queries.transmission.notificationAndResendCount

	query.Lock()
	defer query.Unlock()

	intIds := make([]int64, len(ids))
	for k, v := range ids {
		intIds[k] = int64(v)
	}

	if err := query.SetInt64ParamsIn(obx.Transmission_.Notification, intIds...); err != nil {
		return nil, err
	}

	if err := query.SetInt64Params(obx.Transmission_.ResendCount, int64(resendLimit)); err != nil {
		return nil, err
	}

	return query.Find()
}

func (client *notificationsClient) GetTransmissionsByStartEnd(start int64, end int64, resendLimit int) ([]contract.Transmission, error) {
	var query = &client.queries.transmission.createdBAndResendCount

	query.Lock()
	defer query.Unlock()

	if err := query.SetInt64Params(obx.Transmission_.Created, start, end); err != nil {
		return nil, err
	}

	if err := query.SetInt64Params(obx.Transmission_.ResendCount, int64(resendLimit)); err != nil {
		return nil, err
	}

	return query.Find()
}

func (client *notificationsClient) GetTransmissionsByStart(start int64, resendLimit int) ([]contract.Transmission, error) {
	var query = &client.queries.transmission.createdGTAndResendCount

	query.Lock()
	defer query.Unlock()

	if err := query.SetInt64Params(obx.Transmission_.Created, start); err != nil {
		return nil, err
	}

	if err := query.SetInt64Params(obx.Transmission_.ResendCount, int64(resendLimit)); err != nil {
		return nil, err
	}

	return query.Find()
}

func (client *notificationsClient) GetTransmissionsByEnd(end int64, resendLimit int) ([]contract.Transmission, error) {
	var query = &client.queries.transmission.createdLTAndResendCount

	query.Lock()
	defer query.Unlock()

	if err := query.SetInt64Params(obx.Transmission_.Created, end); err != nil {
		return nil, err
	}

	if err := query.SetInt64Params(obx.Transmission_.ResendCount, int64(resendLimit)); err != nil {
		return nil, err
	}

	return query.Find()
}

func (client *notificationsClient) GetTransmissionsByStatus(resendLimit int, status contract.TransmissionStatus) ([]contract.Transmission, error) {
	var query = &client.queries.transmission.statusAndResendCount

	query.Lock()
	defer query.Unlock()

	if err := query.SetStringParams(obx.Transmission_.Status, string(status)); err != nil {
		return nil, err
	}

	if err := query.SetInt64Params(obx.Transmission_.ResendCount, int64(resendLimit)); err != nil {
		return nil, err
	}

	return query.Find()
}

func (client *notificationsClient) AddTransmission(t contract.Transmission) (string, error) {
	onCreate(&t.BaseObject)

	id, err := client.transmissionBox.Put(&t)
	return obx.IdToString(id), err
}

func (client *notificationsClient) UpdateTransmission(t contract.Transmission) error {
	onUpdate(&t.BaseObject)

	if id, err := obx.IdFromString(t.ID); err != nil {
		return err
	} else if exists, err := client.transmissionBox.Contains(id); err != nil {
		return err
	} else if !exists {
		return db.ErrNotFound
	}

	_, err := client.transmissionBox.Put(&t)
	return err
}

func (client *notificationsClient) DeleteTransmission(age int64, status contract.TransmissionStatus) error {
	var query = &client.queries.transmission.statusAndModifiedLT

	query.Lock()
	defer query.Unlock()

	if err := query.SetStringParams(obx.Transmission_.Status, string(status)); err != nil {
		return err
	}

	var end = db.MakeTimestamp() - age
	if err := query.SetInt64Params(obx.Transmission_.Modified, end); err != nil {
		return err
	}

	_, err := query.Remove()
	return err
}

func (client *notificationsClient) deleteTransmissionsByNotificationIds(ids []uint64) error {
	if len(ids) == 0 {
		return nil
	}

	var query = &client.queries.transmission.notification

	query.Lock()
	defer query.Unlock()

	intIds := make([]int64, len(ids))
	for k, v := range ids {
		intIds[k] = int64(v)
	}

	if err := query.SetInt64ParamsIn(obx.Transmission_.Notification, intIds...); err != nil {
		return err
	}

	_, err := query.Remove()
	return err
}

func (client *notificationsClient) Cleanup() error { return client.CleanupOld(client.cleanupDefaultAge) }

func (client *notificationsClient) CleanupOld(age int) error {
	var query = &client.queries.notification.modifiedLT

	query.Lock()
	defer query.Unlock()

	currentTime := db.MakeTimestamp()
	end := int(currentTime) - age
	if err := query.SetInt64Params(obx.Notification_.Modified, int64(end)); err != nil {
		return err
	}

	// first remove all notifications (this sets related transmission.NotificationId = 0)
	if count, err := query.Remove(); err != nil {
		return err
	} else if count == 0 {
		return nil // nothing deleted, no need to delete transmissions
	} else {
		return client.deleteTransmissionsByNotificationIds([]uint64{0})
	}
}