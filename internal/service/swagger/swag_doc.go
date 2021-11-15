package swagger

import "encoding/json"

// SwagDoc is the top swagger structure
type SwagDoc struct {
	Host         string                 `json:"host"`
	Version      string                 `json:"swagger"`
	Info         SwagInfo               `json:"info"`
	Tags         []SwagTag              `json:"tags"`
	Schemes      []string               `json:"schemes"`
	BasePath     string                 `json:"basePath"`
	EncodingsIn  []string               `json:"consumes"`
	EncodingsOut []string               `json:"produces"`
	Paths        map[string]SwagMethods `json:"paths"` // path -> methods
	Auth         SwagValue              `json:"x-auth"`
}

// SwagMethods is the method map
type SwagMethods map[string]*SwagAPI // method -> api

// SwagAPI is the api specific
type SwagAPI struct {
	OpenRequest  bool            `json:"x-open-request"`
	Tags         []string        `json:"tags"`
	Parameters   SwagParameters  `json:"parameters"`
	Responses    SwagResponse    `json:"responses"`
	EncodingsIn  []string        `json:"consumes"`
	EncodingsOut []string        `json:"produces"`
	Summary      string          `json:"summary"`
	Desc         string          `json:"description"`
	Deprecated   bool            `json:"deprecated"`
	Security     json.RawMessage `json:"security"`
	OperationID  string          `json:"operationId"`
}

// SwagInfo is the information of swagger
type SwagInfo struct {
	Title   string      `json:"title"`
	Version string      `json:"version"`
	Desc    string      `json:"description"`
	Contact SwagContact `json:"contact"`
}

// SwagContact is the contact info of this swag
type SwagContact struct {
	Name  string `json:"name"`
	URL   string `json:"url"`
	Email string `json:"email"`
}

// SwagTag is the swag tag
type SwagTag struct {
	Name string `json:"name"`
	Desc string `json:"description"`
}

// SwagObjectProperties SwagObjectProperties
type SwagObjectProperties map[string]SwagValue

// SwagResponse represents response in swagger
type SwagResponse map[string]*SwagResponseObject

// SwagResponseObject represents response object in swagger
type SwagResponseObject struct {
	Desc   string    `json:"description"`
	Schema SwagValue `json:"schema"`
}

// SwagValue represents common value structure in swagger
type SwagValue map[string]interface{}

// SwagParameters represents input from swagger
type SwagParameters []interface{}
