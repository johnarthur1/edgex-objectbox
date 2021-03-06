/*******************************************************************************
 * Copyright 2017 Dell Inc.
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
 *******************************************************************************/

package defs

import (
	"github.com/edgexfoundry/go-mod-core-contracts/models"
)

// DeviceProfile is a template for a device
type DeviceProfile struct {
	models.DescribedObject `objectbox:"inline"`
	Id                     string
	Name                   string `objectbox:"unique"`
	Manufacturer           string
	Model                  string
	Labels                 []string
	DeviceResources        []models.DeviceResource  `objectbox:"type:[]byte converter:deviceResourcesJson"`
	DeviceCommands         []models.ProfileResource `objectbox:"type:[]byte converter:profileResourcesJson"`
	CoreCommands           []models.Command         `objectbox:"type:[]byte converter:commandsJson"`
}
