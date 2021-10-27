package serializers

/*
Marshaller and unmarshaller for HTTP request and response
*/
type Marshaller interface {
	//Marshall the data
	Serialize(interface{}) ([]byte, error)
	//Preferred mimetype for the output of the serializer
	SerializationCompatibleMimetype() string
	//Unmarshall the raw stream into the provided data
	Deserialize([]byte) (interface{}, error)
	//Preferred mimetypes for the input of the deserializer,
	//TODO quality values ;q= and * asterisks are not accepted
	DeserializationCompatibleMimetypes() []string
}
