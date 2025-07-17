package mcp

// Based on the MCP specification: https://modelcontextprotocol.io/spec/

// ParameterDetail describes a single parameter for an operation.
type ParameterDetail struct {
	Name string `json:"name"`
	In   string `json:"in"` // Location (query, header, path, cookie)
	// Add other details if needed, e.g., required, type
}

// OperationDetail holds the necessary information to execute a specific API operation.
type OperationDetail struct {
	Method     string            `json:"method"`
	Path       string            `json:"path"` // Path template (e.g., /users/{id})
	BaseURL    string            `json:"baseUrl"`
	Parameters []ParameterDetail `json:"parameters,omitempty"`
	// Add RequestBody schema if needed
}

// ToolSet represents the collection of tools provided by an MCP server.
type ToolSet struct {
	MCPVersion  string `json:"mcp_version"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	// Auth        *AuthInfo `json:"auth,omitempty"` // Removed authentication info
	Tools []Tool `json:"tools"`

	// Operations maps Tool.Name (operationId) to its execution details.
	// This is internal to the server and not part of the standard MCP JSON response.
	Operations map[string]OperationDetail `json:"-"` // Use json:"-" to exclude from JSON

	// Internal fields for server-side auth handling (not exposed in JSON)
	apiKeyName string // e.g., "key", "X-API-Key"
	apiKeyIn   string // e.g., "query", "header"
}

// SetAPIKeyDetails allows the parser to set internal API key info.
func (ts *ToolSet) SetAPIKeyDetails(name, in string) {
	ts.apiKeyName = name
	ts.apiKeyIn = in
}

// GetAPIKeyDetails allows the server to retrieve internal API key info.
// We might need this later when making the request.
func (ts *ToolSet) GetAPIKeyDetails() (name, in string) {
	return ts.apiKeyName, ts.apiKeyIn
}

// Tool represents a single function or capability exposed via MCP.
type Tool struct {
	Name         string           `json:"name"` // Corresponds to OpenAPI operationId or generated name
	Description  string           `json:"description,omitempty"`
	InputSchema  Schema           `json:"inputSchema"` // Renamed from Parameters, consolidate parameters/body here
	OutputSchema Schema           `json:"outputSchema,omitempty"`
	Annotations  *ToolAnnotations `json:"annotations,omitempty"`
	// Entrypoint  string      `json:"entrypoint"`             // Removed for simplicity, schema should contain enough info?
	// RequestBody RequestBody `json:"request_body,omitempty"` // Removed, info should be part of InputSchema
	// HTTPMethod  string      `json:"http_method"`            // Removed for simplicity
	// TODO: Add Response handling if needed by spec/client
}

// ToolAnnotations captures optional MCP annotation hints for a tool.
type ToolAnnotations struct {
	Title           string `json:"title,omitempty"`
	ReadOnlyHint    *bool  `json:"readOnlyHint,omitempty"`
	DestructiveHint *bool  `json:"destructiveHint,omitempty"`
	IdempotentHint  *bool  `json:"idempotentHint,omitempty"`
	OpenWorldHint   *bool  `json:"openWorldHint,omitempty"`
}

// RequestBody describes the expected request body for a tool.
// This might become redundant if all info is in InputSchema.
// Keeping it for now as the parser might still use it internally.
type RequestBody struct {
	Description string            `json:"description,omitempty"`
	Required    bool              `json:"required,omitempty"`
	Content     map[string]Schema `json:"content"` // Keyed by media type (e.g., "application/json")
}

// Schema defines the structure and constraints of data (parameters or request/response bodies).
// This mirrors a subset of JSON Schema properties.
type Schema struct {
	Type        string            `json:"type,omitempty"` // e.g., "object", "string", "integer", "array"
	Description string            `json:"description,omitempty"`
	Properties  map[string]Schema `json:"properties,omitempty"` // For type "object"
	Required    []string          `json:"required,omitempty"`   // For type "object"
	Items       *Schema           `json:"items,omitempty"`      // For type "array"
	Format      string            `json:"format,omitempty"`     // e.g., "int32", "date-time"
	Enum        []interface{}     `json:"enum,omitempty"`
	OneOf       []Schema          `json:"oneOf,omitempty"`
	Nullable    bool              `json:"nullable,omitempty"`
	// Add other relevant JSON Schema fields as needed (e.g., minimum, maximum, pattern)
}
