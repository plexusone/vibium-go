// Command genvpatschema generates JSON Schema from VPAT Go types.
package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/agentplexus/vibium-go/vpat"
	"github.com/invopop/jsonschema"
)

func main() {
	r := new(jsonschema.Reflector)
	r.DoNotReference = true

	schema := r.Reflect(&vpat.Report{})
	schema.ID = "https://github.com/agentplexus/vibium-go/vpat/vpat.schema.json"
	schema.Title = "VPAT Report Schema"
	schema.Description = "JSON Schema for Vibium VPAT (Voluntary Product Accessibility Template) reports"

	data, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(string(data))
}
