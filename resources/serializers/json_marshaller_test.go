package serializers

import (
	"reflect"
	"testing"

	"github.com/mitchellh/mapstructure"
	"github.com/okayawright/exp_http_client/resources/misc"
)

/* Nominal case, marshalling and unmashalling JSON data*/
func TestJsonMarshallerNominalMarshallingUnmarshalling(t *testing.T) {
	type Movement struct {
		Label string
		Date  string
		Price float32
	}
	type AccountBalance struct {
		Movements []Movement
		Total     float32
	}
	date1 := "2021-11-29T11:45:26.371Z"
	date2 := "2021-12-01T07:52:21.002Z"
	date3 := "2021-12-02T19:57:02.584Z"
	data := map[string]AccountBalance{
		"Simon": {
			Movements: []Movement{
				{
					Label: "Supermarket",
					Date:  date1,
					Price: 10.52,
				},
				{
					Label: "Gas station",
					Date:  date2,
					Price: 60.1,
				},
			},
			Total: 70.62,
		},
		"Dominique": {
			Movements: []Movement{
				{
					Label: "Theater",
					Date:  date1,
					Price: 9,
				},
				{
					Label: "Theater",
					Date:  date2,
					Price: 9,
				},
				{
					Label: "Theater",
					Date:  date3,
					Price: 9,
				},
			},
			Total: 27,
		},
	}

	serializer := NewJsonMarshaller()
	serializedData, err := serializer.Serialize(data)
	if err != nil {
		t.Fatalf("Serialize() unexpected error %v", err)
	}
	unserializedData, err := serializer.Deserialize(serializedData)
	if err != nil {
		t.Fatalf("Deserialize() unexpected error %v", err)
	}
	//Convert the interface{} into a proper map for reflect.DeepEqual() to properly work
	var unserializedDataStruct map[string]AccountBalance
	mapstructure.Decode(unserializedData, &unserializedDataStruct)
	if err != nil {
		t.Fatalf("Deserialize() unexpected error %v", err)
	}
	if !reflect.DeepEqual(unserializedDataStruct, data) {
		t.Errorf("Deserialize() = %v, want %v", unserializedDataStruct, data)
	}

}

/* Error case, empty JSON data to unmarshal*/
func TestJsonMarshallerErrorEmptyUnmarshalling(t *testing.T) {
	serializer := NewJsonMarshaller()
	empty := new([]byte)
	_, err := serializer.Deserialize(*empty)
	if err == nil {
		t.Errorf("Deserialize() unexpected success")
	}
}

/* Nominal case, compatible JSON mime type for serialization*/
func TestJsonMarshallerNominalCompatibleSerializationMimeType(t *testing.T) {
	expected := "application/vnd.api+json"
	serializer := NewJsonMarshaller()
	if serializer.SerializationCompatibleMimetype() != expected {
		t.Errorf("SerializationCompatibleMimetype() = %v, want %v", serializer.SerializationCompatibleMimetype(), expected)
	}
}

/* Nominal case, not compatible JSON mime type for serialization*/
func TestJsonMarshallerNominalNotCompatibleSerializationMimeType(t *testing.T) {
	expected := "application/soap+xml"
	serializer := NewJsonMarshaller()
	if serializer.SerializationCompatibleMimetype() == expected {
		t.Errorf("SerializationCompatibleMimetype() = %v, want %v", serializer.SerializationCompatibleMimetype(), expected)
	}
}

/* Nominal case, compatible JSON mime type for deserialization */
func TestJsonMarshallerNominalCompatibleDeserializationMimeType(t *testing.T) {
	probe1 := "application/vnd.api+json"
	serializer := NewJsonMarshaller()
	if misc.Find(serializer.DeserializationCompatibleMimetypes(), probe1, false) == -1 {
		t.Errorf("DeserializationCompatibleMimetypes() = %v, want %v", serializer.DeserializationCompatibleMimetypes(), probe1)
	}
	probe2 := "application/json"
	if misc.Find(serializer.DeserializationCompatibleMimetypes(), probe2, false) == -1 {
		t.Errorf("DeserializationCompatibleMimetypes() = %v, want %v", serializer.DeserializationCompatibleMimetypes(), probe2)
	}
}
