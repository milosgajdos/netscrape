package broker

// Marshaler is used for marshaling ingester data.
type Marshaler interface {
	// Marshal marshals data into slice of bytes.
	Marshal(interface{}) ([]byte, error)
}

// Unmarshaler is used for unmarshaling digester data.
type Unmarshaler interface {
	// Unmarshal unmarshals arbitrary bytes into data.
	Unmarshal([]byte, interface{}) error
}

// Encoder encodes data to Message.
type Encoder interface {
	// Returns Message encoded from data.
	Encode(interface{}) (Message, error)
}

// Decode decodes data from Message.
type Decode interface {
	// Decode decodes data from the given Message.
	Decode(Message, interface{}) error
}
