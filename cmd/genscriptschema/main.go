// Command genscriptschema generates JSON Schema from Go types.
package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/invopop/jsonschema"

	"github.com/plexusone/vibium-go/script"
)

func main() {
	r := new(jsonschema.Reflector)
	r.ExpandedStruct = true

	schema := r.Reflect(&script.Script{})
	schema.ID = "https://github.com/plexusone/vibium-go/script/vibium-script.schema.json"
	schema.Title = "Vibium Test Script"
	schema.Description = "Schema for Vibium browser automation test scripts"

	data, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling schema: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(string(data))
}
