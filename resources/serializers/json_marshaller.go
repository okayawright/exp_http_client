package serializers

import (
	"encoding/json"
)

/* (Un)marshaller for JSON */
type jsonMarshaller struct{}

/* Marshaller c'tor */
func NewJsonMarshaller() *jsonMarshaller {
	return &jsonMarshaller{}
}

func (marshaller *jsonMarshaller) Serialize(input interface{}) ([]byte, error) {
	return json.Marshal(input)
}

func (marshaller *jsonMarshaller) SerializationCompatibleMimetype() string {
	return "application/vnd.api+json"
}

func (marshaller *jsonMarshaller) Deserialize(input []byte) (interface{}, error) {
	var output interface{}
	err := json.Unmarshal(input, &output)
	return output, err
}

func (marshaller *jsonMarshaller) DeserializationCompatibleMimetypes() []string {
	return []string{
		"application/vnd.api+json",
		"application/json",
	}
}
