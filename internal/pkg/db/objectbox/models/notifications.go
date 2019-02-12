/*******************************************************************************
 * Copyright 2018 Dell Technologies Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
 *
 *******************************************************************************/

package models

import (
	"github.com/edgexfoundry/edgex-go/pkg/models"
)

type Notification struct {
	models.BaseObject `inline`
	ID                string
	Slug              string `unique`
	Sender            string
	Category          models.NotificationsCategory
	Severity          models.NotificationsSeverity
	Content           string
	Description       string
	Status            models.NotificationsStatus
	Labels            []string
	ContentType       string
}