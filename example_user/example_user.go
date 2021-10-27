package example_user

import (
	"errors"
	"net/url"
	"strconv"

	"github.com/mitchellh/mapstructure"
	"github.com/okayawright/exp_http_client/resources"
)

type Data struct {
	Data         *User   `json:"data,omitempty"`
	ErrorCode    *string `json:"error_code,omitempty"`
	ErrorMessage *string `json:"error_message,omitempty"`
}

type User struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name,omitempty"`
}

/* Convert the polymorphic JSON-produced map of data into a proper User struct */
func convertResponse(data interface{}) (*Data, error) {
	var output Data
	var err error
	if data != nil {
		dataUserStruct, ok := data.(map[string]interface{})
		if !ok {
			//We were not expecting a totally new structure in the response !
			return nil, errors.New("The response format is invalid")
		}

		//Unmarshalling

		cfg := &mapstructure.DecoderConfig{
			Metadata: nil,
			Result:   &output,
			//Use the JSON tags as hints
			TagName: "json",
		}
		decoder, _ := mapstructure.NewDecoder(cfg)
		err = decoder.Decode(dataUserStruct)
	}
	return &output, err
}

/* Merge the user-defined protocol and host part of an URL into a pre-configured one */
func mergeUrlHost(userDefinedHost string, defaultUrl string) *url.URL {
	actualUrl, err := url.Parse(defaultUrl)
	if err != nil {
		panic("Invalid default URL")
	}
	customUrl, err := url.Parse(userDefinedHost)
	if err != nil {
		panic("Invalid hostname")
	}
	if len(customUrl.Scheme) > 0 {
		actualUrl.Scheme = customUrl.Scheme
	}
	if customUrl.User != nil {
		actualUrl.User = customUrl.User
	}
	if len(customUrl.Host) > 0 {
		actualUrl.Host = customUrl.Host
	}
	return actualUrl
}

/* Create a new User resource defined by the save data, on the given host */
func Create(host string, save *Data) (*Data, int, error) {

	//Configure the resource
	res := resources.NewResource(mergeUrlHost(host, "http://localhost:8080/v1/membership/users"))

	//We do not need to cancel the request here so we can go straight to execute call() after Prepare()
	call, _, err := res.Request("POST", nil, save)
	if err != nil {
		return nil, 0, err
	}
	responseStruct, statusCode, err := call()
	if err != nil {
		return nil, statusCode, err
	}

	//Convert the polymorphic result body into a struct
	dataUserStruct, err := convertResponse(responseStruct)
	return dataUserStruct, statusCode, err
}

/* Get an existing User resource identified by the specified identifier, on the given host */
func Fetch(host string, id string) (*Data, int, error) {

	//Configure the resource
	res := resources.NewResource(mergeUrlHost(host, "http://localhost:8080/v1/membership/users/{user_id}"))

	//We do not need to cancel the request here so we can go straight to execute call() after Prepare()
	call, _, err := res.Request("GET", &map[string]string{
		"user_id": id,
	}, nil)
	if err != nil {
		return nil, 0, err
	}
	responseStruct, statusCode, err := call()
	if err != nil {
		return nil, statusCode, err
	}

	//Convert the polymorphic result body into a struct
	dataUserStruct, err := convertResponse(responseStruct)
	return dataUserStruct, statusCode, err
}

/* Delete an existing User resource identified by the specified identifier and version, on the given host */
func Delete(host string, id string, version int) (*Data, int, error) {

	//Configure the resource
	res := resources.NewResource(mergeUrlHost(host, "http://localhost:8080/v1/membership/users/{user_id}?version={version}"))

	//We do not need to cancel the request here so we can go straight to execute call() after Prepare()
	call, _, err := res.Request("DELETE", &map[string]string{
		"user_id": id,
		"version": strconv.Itoa(version),
	}, nil)
	if err != nil {
		return nil, 0, err
	}
	responseStruct, statusCode, err := call()
	if err != nil {
		return nil, statusCode, err
	}

	//Convert the polymorphic result body into a struct
	dataUserStruct, err := convertResponse(responseStruct)
	return dataUserStruct, statusCode, err
}
