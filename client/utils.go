package client

import "gitlab.sandmanbb.com/perfil-digital-agro/agro-ws/build/scicrop-scraper/models/schema"

var CANONICAL_FIELDS = []string{"dateTime"};

func MapUsingSchema(input map[string]string, schema schema.Schema) map[string]string {
	output := make(map[string]string)
	for _, field := range CANONICAL_FIELDS {
		if val, ok := input[field]; ok {
			output[field] = val
		}
	}

	for _, field := range schema.Fields {
		to := field.Name
		from := field.From

		if len(from) == 0 {
			from = to
		}

		if val, ok := input[from]; ok {
			output[to] = val
		}
	}
	return output
}