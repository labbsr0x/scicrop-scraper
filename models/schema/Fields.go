package schema

type Field struct {
	Name string 							`json:"name"`
	Type string 							`json:"type"`
	Description string						`json:"description"`
	UrlExample  string 						`json:"url_example"`
	From 	    string 						`json:"-"`
	Properties []map[string]interface{}   	`json:"properties"`
}