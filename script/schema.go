package script

import _ "embed"

//go:embed vibium-script.schema.json
var SchemaJSON []byte

// Schema returns the JSON Schema for Vibium test scripts.
func Schema() []byte {
	return SchemaJSON
}
