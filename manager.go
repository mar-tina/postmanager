package manager

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/mar-tina/postmanager/schema"
)

type Manager interface {
	CollectionFromID()
	AddToAncestry()
	MigrateAncestor()
	CreateEndpoint()
	Use(collection string, route string)
	StatusCheck()
	NewCollection(name string, body interface{})
}

type ApiMgr struct {
	endpoints map[string][]string
	resource  string
	username  string
	password  string
}

//Use specifies the collection being used, passing in the endpoint and the name of the resource to use.
func (mgr *ApiMgr) Use(collectionName string, resourceName string) error {

	//find collection in database . fetch single collection from postman API.
	//make sure resource exists and then proceed to add the resource to the endpoints list. Does not check to make sure
	//endpoint is live or responsive.
	payload := schema.ResourceCheckPayload{}
	payload.CollectionName = collectionName
	payload.ResourceName = resourceName

	req, err := mgr.prepareRequest("POST", "/resource/exists", payload)
	if err != nil {
		return err
	}

	client := http.Client{}
	_, err = client.Do(req)
	if err != nil {
		return err
	}

	coll := mgr.endpoints[collectionName]
	coll = append(coll, resourceName)
	mgr.endpoints[collectionName] = coll
	return nil
}

func (mg *ApiMgr) New(resource string) *ApiMgr {
	mg.resource = resource
	return mg
}

func (mg *ApiMgr) WithBasicAuth(username, password string) *ApiMgr {
	mg.username = username
	mg.password = password
	return mg
}

// func (mgr *ApiMgr) StatusCheck() {
// 	for key, val := range mgr.endpoints {
// 		go func(collection string, route string) {
// 			//postman.call(coll, route)
// 		}(key, val)
// 	}
// }

func (mgr *ApiMgr) prepareRequest(method, route string, payload interface{}) (*http.Request, error) {
	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, errors.New("failed to prepare request")
	}

	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", mgr.resource, route), bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to prepare request %s", err)
	}

	req.SetBasicAuth(mgr.username, mgr.password)
	return req, nil
}
