package mcp

import (
	"encoding/json"
	"reflect"
	"testing"
)

type TestParam struct {
	name        string
	mcpRequest  *MCPRequest
	schema      *SchemaMapping
	wantQuery   string
	wantHeaders map[string]string
	wantErr     bool
}

func TestProcessQueryParameters(t *testing.T) {
	tests := []TestParam{
		{
			name: "all required present",
			mcpRequest: &MCPRequest{
				Arguments: `{"foo":"bar","baz":"qux"}`,
			},
			schema: &SchemaMapping{
				QueryParameters: []Param{
					{Name: "foo", Required: true},
					{Name: "baz", Required: false},
				},
			},
			wantQuery: "?foo=bar&baz=qux",
			wantErr:   false,
		},
		{
			name: "missing required param",
			mcpRequest: &MCPRequest{
				Arguments: `{"foo":"bar"}`,
			},
			schema: &SchemaMapping{
				QueryParameters: []Param{
					{Name: "foo", Required: true},
					{Name: "baz", Required: true},
				},
			},
			wantQuery: "",
			wantErr:   true,
		},
		{
			name: "optional param missing",
			mcpRequest: &MCPRequest{
				Arguments: `{"foo":"bar"}`,
			},
			schema: &SchemaMapping{
				QueryParameters: []Param{
					{Name: "foo", Required: true},
					{Name: "baz", Required: false},
				},
			},
			wantQuery: "?foo=bar",
			wantErr:   false,
		},
		{
			name: "no query params",
			mcpRequest: &MCPRequest{
				Arguments: `{}`,
			},
			schema: &SchemaMapping{
				QueryParameters: []Param{},
			},
			wantQuery: "",
			wantErr:   false,
		},
		{
			name: "url encoded params",
			mcpRequest: &MCPRequest{
				Arguments: `{"foo":"bar baz","baz":"qux!"}`,
			},
			schema: &SchemaMapping{
				QueryParameters: []Param{
					{Name: "foo", Required: true},
					{Name: "baz", Required: false},
				},
			},
			wantQuery: "?foo=bar%20baz&baz=qux%21",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args, err := parseArgs(tt.mcpRequest)
			if err != nil {
				t.Errorf("parseArgs() error = %v", err)
				return
			}
			got, err := processQueryParameters((args), tt.schema)
			if (err != nil) != tt.wantErr {
				t.Errorf("processQueryParameters() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.wantQuery {
				t.Errorf("processQueryParameters() = %v, want %v", got, tt.wantQuery)
			}
		})
	}
}

func TestProcessHeaderParameters(t *testing.T) {
	tests := []TestParam{
		{
			name: "all required present",
			mcpRequest: &MCPRequest{
				Arguments: `{"header1":"value1","header2":"value2"}`,
				API:       APIInfo{},
			},
			schema: &SchemaMapping{
				HeaderParameters: []Param{
					{Name: "header1", Required: true},
					{Name: "header2", Required: false},
				},
				ContentType: "application/json",
			},
			wantHeaders: map[string]string{
				"header1":      "value1",
				"header2":      "value2",
				"Content-Type": "application/json",
			},
			wantErr: false,
		},
		{
			name: "missing required header",
			mcpRequest: &MCPRequest{
				Arguments: `{"header1":"value1"}`,
				API:       APIInfo{},
			},
			schema: &SchemaMapping{
				HeaderParameters: []Param{
					{Name: "header1", Required: true},
					{Name: "header2", Required: true},
				},
				ContentType: "application/json",
			},
			wantHeaders: nil,
			wantErr:     true,
		},
		{
			name: "optional header missing",
			mcpRequest: &MCPRequest{
				Arguments: `{"header1":"value1"}`,
				API:       APIInfo{},
			},
			schema: &SchemaMapping{
				HeaderParameters: []Param{
					{Name: "header1", Required: true},
					{Name: "header2", Required: false},
				},
				ContentType: "application/json",
			},
			wantHeaders: map[string]string{
				"header1":      "value1",
				"Content-Type": "application/json",
			},
			wantErr: false,
		},
		{
			name: "auth header present",
			mcpRequest: &MCPRequest{
				Arguments: `{"header1":"value1"}`,
				API:       APIInfo{Auth: "Authorization: Bearer token"},
			},
			schema: &SchemaMapping{
				HeaderParameters: []Param{
					{Name: "header1", Required: true},
				},
				ContentType: "application/json",
			},
			wantHeaders: map[string]string{
				"header1":       "value1",
				"Authorization": "Bearer token",
				"Content-Type":  "application/json",
			},
			wantErr: false,
		},
		{
			name: "default content type",
			mcpRequest: &MCPRequest{
				Arguments: `{"header1":"value1"}`,
				API:       APIInfo{},
			},
			schema: &SchemaMapping{
				HeaderParameters: []Param{
					{Name: "header1", Required: true},
				},
				ContentType: "",
			},
			wantHeaders: map[string]string{
				"header1":      "value1",
				"Content-Type": ContentTypeJSON,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := processHeaderParameters(tt.mcpRequest, tt.schema)
			if (err != nil) != tt.wantErr {
				t.Errorf("processHeaderParameters() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.wantHeaders) {
				gotJSON, _ := json.Marshal(got)
				wantJSON, _ := json.Marshal(tt.wantHeaders)
				t.Errorf("processHeaderParameters() = %s, want %s", gotJSON, wantJSON)
			}
		})
	}
}
