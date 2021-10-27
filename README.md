# Simple REST client library in Go

## Context
The goal was to make a minimal library that could be easily reused and expanded in other projects. It doesn't rely on any third-party modules except for [mapstructure](https://github.com/mitchellh/mapstructure) which is used in the test implementation. Neither does it feature more complex features that are expected to be found in mature clients (e.g. HATEOS support, metrics, Swagger support, authentication, etc).

## Usage

### Resources package
The generic REST API client library.

1. In order to call an API you need to first instantiate a **resource** with *NewResource()*
    ```
    res := resources.NewResource("http://localhost:8080/v1/membership/users")
    ```
    Specify the API endpoint corresponding to this resource as a parameter. By default, you can read and write JSON requests and responses.

    You can also use the URL as a template with named parameters enclosed within curly braces and that will be resolved later on with *Request()* (see below.)
    ```
    res := resources.NewResource("http://localhost:8080/v1/membership/users/{user_id}")
    ```

    You can change the behaviour of this **resource** with chainable methods:

    - *WithMarshaller()* lets you change the default request and response serializer/deserializer by specifying a new **marshaller** midlleware. Right now only a JSON-specific (de)serializer is implemented.
        ```
        res.WithMarshaller(serializers.NewJsonMarshaller())
        ```
    - *WithClient()* lets you override the default HTTP client engine if needed.
        ```
        res.WithClient(http.DefaultClient)
        ```
    - *WithTimeout()* lets you define a custom request timeout, in seconds, that globally applies to all verbs on this resource. By default it is set to 30s.
        ```
        res.WithTimeout(60)
        ```
    - *WithRetrier()* lets you select an alternative call retry strategy by specifying a new **retrier** middleware. Right now only an exponential backoff is implemented.
        ```
        res.WithRetrier(retriers.NewExponentialRetrier())
        ```
        The default exponential backoff retrying middleware will actually try to call the API again in case it captures a client timeout and the HTTP error codes 429, 500, 503, and 504. Each retry might not be triggered perfectly on time and deviate a bit from their planned scheduled (jittering). In any case, it will then stop retrying after 3 attempts.

        You can change the default behaviour of **exponentialRetrier** with the following chainable methods: *WithRetryableCodes()*, *WithJitter()*, *WithMaxTries()*.
2. On this **resource** you can then define a set of actions that corresponds to a specific combination of an HTTP verb and inputs. An action is setup using the *Request()* method.
    ```
    call, cancel, err := res.Request("GET", &map[string]string{
		"user_id": id,
	}, nil)
    ```
    ```
    call, cancel, err := res.Request("POST", nil, save)
    ```
    The first parameter is the case-insensitive HTTP verb to use for this request. The second one is an optional map of string keys and values representing the named parameters and their corresponding values to replace in the template URL. The third parameter is the optional struct body to send as well, if needed.

    It returns a **CallFunc** and a **CancelFunc** (see below), and potential errors.
3. The **CallFunc** function will let you make the actual HTTP request, that can be programmatically cancelled by executing the corresponding **CancelFunc** function. You can execute **CallFunc** multiple times in a row, or in parallel.
    ```
    body, code, err := call()
    ```
    ```
    go func() {
		body, code, err := call()
        ...
	}()
    ...
	cancel()
    ```
    **CallFunc** returns the structured body of the response if available, as a *map[string]interface{}*, the HTTP status code, and potential errors.

### Example of use
As a test implementation for this library, there's an example package *example_user* that provides standard `Create`, `Fetch`, and `Delete` operations on an imaginary `user` resource.
In order to keep it simple I didn't expose the cancel function in this version.
#### Create
Create an `user` resource 
```
body, status, err := sample_membership.Create("http://sampleapi:8080/", &user)
```
Set the protocol and host of the API and a `example_user.Data` struct that represents the user to persist.
Gets back a `example_user.Data` struct corresponding to the response, the HTTP status, and potential errors.
#### Fetch
Retch a specific `user` resource 
```
body, status, err := example_user.Fetch("http://sampleapi:8080", userId)
```
Set the protocol and host of the API and the user identifier of the resource you want to read.
Gets back a `example_user.Data` struct corresponding to theresponse, the HTTP status, and potential errors.
#### Delete
Delete a specific `user` resource 
```
body, status, err := example_user.Delete("http://sampleapi:8080/", userId, version)
```
Set the protocol and host of the API and the user identifier and version of the resource you want to remove.
Gets back a `example_user.Data` struct corresponding to the response, the HTTP status, and potential errors.