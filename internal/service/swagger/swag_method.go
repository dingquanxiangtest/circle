package swagger

import (
	"fmt"
	"strings"
)

// SwagItem SwagItem
type SwagItem interface {
	GetTag() string
	SetField(value interface{}, appID, tableID string)
	GenSwagMethods() (SwagMethods, error)
}

// MethodCommon MethodCommon
type MethodCommon struct {
	appID   string
	tableID string
	schema  interface{}
}

// MGenSwaggerInput MGenSwaggerInput
func (m *MethodCommon) MGenSwaggerInput(properties SwagObjectProperties, required []string) (SwagParameters, error) {
	resp := make(SwagParameters, 0)
	schema := SwagValue{
		"type": "object",
	}
	r := SwagValue{
		"schema": schema,
	}
	r["name"] = "root"
	r["description"] = "body inputs"
	r["in"] = "body"

	schema["properties"] = properties
	schema["required"] = required
	resp = append(resp, r)
	return resp, nil

}

// MGenSwaggerOutput  MGenSwaggerOutput
func (m *MethodCommon) MGenSwaggerOutput(swagValue SwagValue) (SwagResponse, error) {
	resp := SwagResponse{"200": &SwagResponseObject{
		Desc: "",
		Schema: SwagValue{
			"type":        "object",
			"description": "response schema.",
			"properties": SwagObjectProperties{
				"code": SwagValue{
					"type":  "number",
					"title": "0:success, others: error",
				},
				"msg": SwagValue{
					"type":  "string",
					"title": "error message when code is not 0",
				},
				"data": swagValue,
			},
		},
	}}
	return resp, nil
}

// MGenSwagMethods MGenSwagMethods
func (m *MethodCommon) MGenSwagMethods(req SwagParameters, response SwagResponse, operationID string) SwagMethods {
	createAPI := &SwagAPI{
		Responses:    response,
		Parameters:   req,
		OperationID:  operationID,
		EncodingsOut: []string{"application/json"},
		EncodingsIn:  []string{"application/json"},
		OpenRequest:  true,
	}

	createMethod := SwagMethods{
		"post": createAPI,
	}
	return createMethod

}

// Create Create
type Create struct {
	MethodCommon
	schema interface{}
}

// GetTag GetTag
func (c *Create) GetTag() string {
	return "create"
}

// SetField SetField
func (c *Create) SetField(value interface{}, appID, tableID string) {
	c.schema = value
	c.tableID = tableID
	c.appID = appID
}

// GenSwaggerOutput GenSwaggerOutput
func (c *Create) GenSwaggerOutput() (SwagResponse, error) {
	s, ok := c.schema.(Schema)
	if !ok {
		return nil, nil
	}
	schema := SwagValue{
		"type": "object",
	}
	properties := SwagObjectProperties{
		"entity": SwagValue{
			"type":       "object",
			"properties": s,
		},
	}
	schema["properties"] = properties
	output, err := c.MGenSwaggerOutput(schema)
	if err != nil {
		return nil, err
	}
	return output, nil
}

// GenSwaggerInput  GenSwaggerInput
func (c *Create) GenSwaggerInput() (SwagParameters, error) {
	if s, ok := c.schema.(Schema); ok {
		properties := SwagObjectProperties{
			"entity": SwagValue{
				"type":        "object",
				"description": "entity",
				"properties":  s,
			},
		}
		input, err := c.MGenSwaggerInput(properties, []string{"entity"})
		if err != nil {
			return nil, err
		}
		return input, nil
	}
	return nil, nil
}

// GenSwagMethods GenSwagMethods
func (c *Create) GenSwagMethods() (SwagMethods, error) {
	output, err := c.GenSwaggerOutput()
	if err != nil {
		return nil, err
	}
	input, err := c.GenSwaggerInput()
	if err != nil {
		return nil, err
	}
	createMethod := c.MGenSwagMethods(input, output, GetInnerXName(c.tableID, c.GetTag()))
	return createMethod, nil
}

// Update  Update
type Update struct {
	MethodCommon
}

// GetTag GetTag
func (u *Update) GetTag() string {
	return "update"
}

// SetField SetField
func (u *Update) SetField(value interface{}, appID, tableID string) {
	u.schema = value
	u.tableID = tableID
	u.appID = appID
}

// GenSwaggerOutput  GenSwaggerOutput
func (u *Update) GenSwaggerOutput() (SwagResponse, error) {
	if _, ok := u.schema.(Schema); ok {
		schema := SwagValue{
			"type": "object",
		}
		properties := SwagObjectProperties{
			"errorCount": SwagValue{
				"type": "number",
			},
		}
		schema["properties"] = properties
		output, err := u.MGenSwaggerOutput(schema)
		if err != nil {
			return nil, err
		}
		return output, nil
	}
	return nil, nil

}

// GenSwaggerInput GenSwaggerInput
func (u *Update) GenSwaggerInput() (SwagParameters, error) {
	if s, ok := u.schema.(Schema); ok {
		properties := SwagObjectProperties{
			"entity": SwagValue{
				"type":        "object",
				"description": "entity",
				"properties":  s,
			},
			"query": SwagValue{
				"type": "object",
				"properties": SwagValue{
					"term": SwagValue{
						"type": "object",
						"properties": SwagValue{
							"_id": SwagValue{
								"type": "string",
							},
						},
					},
				},
			},
		}
		input, err := u.MGenSwaggerInput(properties, []string{"entity", "query"})
		if err != nil {
			return nil, err
		}
		return input, nil
	}
	return nil, nil
}

// GenSwagMethods GenSwagMethods
func (u *Update) GenSwagMethods() (SwagMethods, error) {
	output, err := u.GenSwaggerOutput()
	if err != nil {
		return nil, err
	}
	input, err := u.GenSwaggerInput()
	if err != nil {
		return nil, err
	}
	createMethod := u.MGenSwagMethods(input, output, GetInnerXName(u.tableID, u.GetTag()))
	return createMethod, nil
}

// Delete Delete
type Delete struct {
	MethodCommon
}

// GetTag GetTag
func (d *Delete) GetTag() string {
	return "delete"
}

// SetField SetField
func (d *Delete) SetField(value interface{}, appID, tableID string) {
	d.schema = value
	d.tableID = tableID
	d.appID = appID
}

// GenSwaggerOutput GenSwaggerOutput
func (d *Delete) GenSwaggerOutput() (SwagResponse, error) {
	if _, ok := d.schema.(Schema); ok {
		schema := SwagValue{
			"type": "object",
		}
		properties := SwagObjectProperties{
			"errorCount": SwagValue{
				"type": "number",
			},
		}
		schema["properties"] = properties
		output, err := d.MGenSwaggerOutput(schema)
		if err != nil {
			return nil, err
		}
		return output, nil
	}
	return nil, nil

}

// GenSwaggerInput GenSwaggerInput
func (d *Delete) GenSwaggerInput() (SwagParameters, error) {
	if _, ok := d.schema.(Schema); ok {
		properties := SwagObjectProperties{
			"query": SwagValue{
				"type": "object",
				"properties": SwagValue{
					"terms": SwagValue{
						"type": "object",
						"properties": SwagValue{
							"_id": SwagValue{
								"type": "array",
								"items": SwagValue{
									"type": "string",
								},
							},
						},
					},
				},
			},
		}
		input, err := d.MGenSwaggerInput(properties, []string{"query"})
		if err != nil {
			return nil, err
		}
		return input, nil
	}
	return nil, nil
}

// GenSwagMethods  GenSwagMethods
func (d *Delete) GenSwagMethods() (SwagMethods, error) {
	output, err := d.GenSwaggerOutput()
	if err != nil {
		return nil, err
	}
	input, err := d.GenSwaggerInput()
	if err != nil {
		return nil, err
	}
	createMethod := d.MGenSwagMethods(input, output, GetInnerXName(d.tableID, d.GetTag()))
	return createMethod, nil
}

// Search Search
type Search struct {
	MethodCommon
}

// GetTag GetTag
func (s *Search) GetTag() string {
	return "search"
}

// SetField SetField
func (s *Search) SetField(value interface{}, appID, tableID string) {
	s.schema = value
	s.tableID = tableID
	s.appID = appID
}

// GenSwaggerOutput GenSwaggerOutput
func (s *Search) GenSwaggerOutput() (SwagResponse, error) {
	if s1, ok := s.schema.(Schema); ok {
		schema := SwagValue{
			"type": "object",
		}
		properties := SwagObjectProperties{
			"entities": SwagValue{
				"type":        "array",
				"description": "entity",
				"items": SwagValue{
					"type":       "object",
					"properties": s1,
				},
			},
			"total": SwagValue{
				"type": "number",
			},
		}
		schema["properties"] = properties
		output, err := s.MGenSwaggerOutput(schema)
		if err != nil {
			return nil, err
		}
		return output, nil
	}
	return nil, nil

}

// GenSwaggerInput GenSwaggerInput
func (s *Search) GenSwaggerInput() (SwagParameters, error) {
	types := SwagValue{
		"type": "object",
	}
	if _, ok := s.schema.(Schema); ok {
		properties := SwagObjectProperties{
			"query": SwagValue{
				"type": "object",
				"properties": SwagValue{
					"term":  types,
					"match": types,
					"range": types,
					"terms": types,
					"bool": SwagValue{
						"type": "object",
						"properties": SwagValue{
							"must": SwagValue{
								"type":  "array",
								"items": types,
							},
							"should": SwagValue{
								"type":  "array",
								"items": types,
							},
							"must_not": SwagValue{
								"type":  "array",
								"items": types,
							},
						},
					},
				},
			},
			"page": SwagValue{
				"type": "number",
			},
			"size": SwagValue{
				"type": "number",
			},
		}
		input, err := s.MGenSwaggerInput(properties, []string{})
		if err != nil {
			return nil, err
		}
		return input, nil
	}
	return nil, nil
}

// GenSwagMethods GenSwagMethods
func (s *Search) GenSwagMethods() (SwagMethods, error) {
	output, err := s.GenSwaggerOutput()
	if err != nil {
		return nil, err
	}
	input, err := s.GenSwaggerInput()
	if err != nil {
		return nil, err
	}
	createMethod := s.MGenSwagMethods(input, output, GetInnerXName(s.tableID, s.GetTag()))
	return createMethod, nil
}

// GenXName GenXName
func GenXName(appID, tableID, tag, content string) string {

	return fmt.Sprintf("/system/app/%s/structor/%s/%s", appID, content, GetInnerXName(tableID, tag))
}

// GetInnerXName GetInnerXName
func GetInnerXName(tableID, tag string) string {
	tableIDs := strings.Split(tableID, "_")
	return fmt.Sprintf("%s_%s", tableIDs[len(tableIDs)-1], tag)
}
