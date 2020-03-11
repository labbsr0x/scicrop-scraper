package schema

type Schema struct {
	Key string          `json:"_key"`
	Id string           `json:"_id"`
	Uri string          `json:"uri"`
	Fields []Field 		`json:"fields"`

	//Meta fields: meta data to guide clients operation
	Node string         `json:"-"`
}