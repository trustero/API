// This file is subject to the terms and conditions defined in
// file 'LICENSE.txt', which is part of this source code package.
package main

import (
	"github.com/trustero/api/go/examples/gitlab_receptor/receptorPackage"

	"github.com/trustero/api/go/receptor_sdk"
	"github.com/trustero/api/go/receptor_sdk/cmd"
	"github.com/trustero/api/go/receptor_v1"
)

// This struct holds the credentials the receptor needs to authenticate with the
// service provider. A display name and placeholder tag should be provided
// which will be used in the UI when activating the receptor.
// This is what will be returned in the GetCredentialObj call

// Receptor defines the GitLab service credentials required for connecting to the GitLab
// service and gathering necessary evidence to support its use.
type Receptor struct {
	Token   string `trustero:"display:GitLab Access Token;placeholder:token"`
	GroupID string `trustero:"display:GitLab Group ID;placeholder:group id"`
}

// Set the name of the receptor in the const declaration above
// This will let the receptor inform Trustero about itself
func (r *Receptor) GetReceptorType() string {
	return receptorPackage.GetReceptorTypeImpl()
}

// Set the names of the services in the const declaration above
// This will let the receptor inform Trustero about itself
// Feel free to add or remove services as needed
func (r *Receptor) GetKnownServices() []string {
	return receptorPackage.GetKnownServicesImpl()
}

// This will return Receptor struct defined above when the receptor is asked to
// identify itself
func (r *Receptor) GetCredentialObj() (credentialObj interface{}) {
	return r
}

// This function will call into the service provider API with the provided
// credentials and confirm that the credentials are valid. Usually a simple
// API call like GET org name. If the credentials are not valid,
// return a relevant error message
func (r *Receptor) Verify(credentials interface{}, config interface{}) (ok bool, err error) {
	c := credentials.(*Receptor)
	return receptorPackage.VerifyImpl(c.Token, c.GroupID)
}

// The Discover function returns a list of Service Entities. This function
// makes any relevant API calls to the Service Provider to gather information
// about how many Service Entity Instances are in use. If at any point this
// function runs into an error, log that error and continue
func (r *Receptor) Discover(credentials interface{}, config interface{}) (svcs []*receptor_v1.ServiceEntity, err error) {
	c := credentials.(*Receptor)
	return receptorPackage.DiscoverImpl(c.Token, c.GroupID)
}

// Report will often make the same API calls made in the Discover call, but it
// will additionally create evidences with the data returned from the API calls
func (r *Receptor) Report(credentials interface{}, config interface{}) (evidences []*receptor_sdk.Evidence, err error) {
	c := credentials.(*Receptor)
	return receptorPackage.ReportImpl(c.Token, c.GroupID)
}

func (r *Receptor) ReportBatch(credentials interface{}, evidenceChan chan []*receptor_sdk.Evidence) {
	close(evidenceChan)
}

// GetEvidenceInfo returns a list of all the possible evidence created. The return value should not have any actual
// rows or source data, just the Caption and Description
func (r *Receptor) GetEvidenceInfo(credentials interface{}) (evidences []*receptor_sdk.Evidence) {
	return
}

func (r *Receptor) Configure(credentials interface{}) (config *receptor_v1.ReceptorConfiguration, err error) {
	c := credentials.(*Receptor)
	return receptorPackage.ConfigureImpl(c.Token, c.GroupID)
}

// This will return Config json when the receptor is asked for
// config
func (r *Receptor) GetConfigObj() (config interface{}) {
	return nil
}

func (r *Receptor) GetAuthMethods() interface{} {
	return receptorPackage.GetAuthMethodImpl()
}

func (r *Receptor) GetConfigObjDesc() interface{} {
	return receptorPackage.GetConfigObjDescImpl()
}

func (r *Receptor) GetLogo() (string, error) {
	logo, _ := receptorPackage.GetLogoImpl()
	return logo, nil
}
func (r *Receptor) GetInstructions() (string, error) {
	logo, _ := receptorPackage.GetInstructionsImpl()
	return logo, nil
}

func main() {
	cmd.Execute(&Receptor{})
}
