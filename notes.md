- use List bubble
- session state enum in main model like this:
    ```golang
    type sessionState int
    const (
       "/" = iota,
       "/paths",
       "/paths/{id}"
    )

    type MainModel struct{
       state sessionState 
       view1 View1Model
       view2 View2Model
       view3 View3Model
    }
    ```
- main model always gets messages first. it delegates downwards

---

# libopenapi
base (all except DynamicValue and SchemaProxy have a *ChangesStruct)
    - Contact
    - Discriminator
    - Example
    - ExternalDoc
    - Info
    - License
    - Schema
    - SecurityRequirement
    - Tag
    - XML
- DynamicValue
- SchemaProxy

model
    - base <=> ContactChanges
    - base <=> DiscriminatorChanges
    - base <=> ExampleChanges
    - base <=> ExternalDocChanges
    - base <=> InfoChanges
    - base <=> LicenseChanges
    - base <=> SchemaChanges
    - base <=> SecurityRequirementChanges
    - base <=> TagChanges
    - base <=> XMLChanges
  - v3 <=> CallbackChanges
  - v3 <=> ComponentsChanges
  - v3 <=> DocumentChanges
  - v3 <=> EncodingChanges
  - v3 <=> HeaderChanges
  - v3 <=> LinkChanges
  - v3 <=> MediaTypeChanges
  - v3 <=> OAuthFlowChanges
  - v3 <=> OAuthFlowsChanges
  - v3 <=> OperationChanges
  - v3 <=> ParameterChanges
  - v3 <=> PathItemChanges
  - v3 <=> PathsChanges
  - v3 <=> RequestBodyChanges
  - v3 <=> ResponseChanges
  - v3 <=> ResponsesChanges
  - v3 <=> SecuritySchemeChanges
  - v3 <=> ServerChanges
  - v3 <=> ServerVariableChanges
- Change
- ChangeContext
- DocumentChangesFlat
- ExamplesChanges
- ExtensionChanges
- ItemsChanges
- PropertyChanges
- PropertyCheck
- ScopesChanges
- WhatChanged

v3 (all have a *Changes struct)
  - Callback
  - Components
  - Document
  - Encoding
  - Header
  - Link
  - MediaType
  - OAuthFlow
  - OAuthFlows
  - Operation
  - Parameter
  - PathItem
  - Paths
  - RequestBody
  - Response
  - Responses
  - SecurityScheme
  - Server
  - ServerVariable

---

Recursively parsing the abstract syntax tree of the v3 package yields the following coordinates, describing how all the structs I care about are composed:

```
- MediaType
	- Schema	*base.SchemaProxy
	- Example	*yaml.Node
	- Examples	*[string, *base.Example]
	- Encoding	*[string, *Encoding]
	- Extensions	*[string, *yaml.Node]

- Operation
	- Tags	[]string
	- Summary	string
	- Description	string
	- ExternalDocs	*base.ExternalDoc
	- OperationId	string
	- Parameters	[]*Parameter
	- RequestBody	*RequestBody
	- Responses	*Responses
	- Callbacks	*[string, *Callback]
	- Deprecated	*bool
	- Security	[]*base.SecurityRequirement
	- Servers	[]*Server
	- Extensions	*[string, *yaml.Node]

- SecurityScheme
	- Type	string
	- Description	string
	- Name	string
	- In	string
	- Scheme	string
	- BearerFormat	string
	- Flows	*OAuthFlows
	- OpenIdConnectUrl	string
	- Extensions	*[string, *yaml.Node]

- RequestBody
	- Description	string
	- Content	*[string, *MediaType]
	- Required	*bool
	- Extensions	*[string, *yaml.Node]

- Responses
	- Codes	*[string, *Response]
	- Default	*Response
	- Extensions	*[string, *yaml.Node]

- ServerVariable
	- Enum	[]string
	- Default	string
	- Description	string
	- Extensions	*[string, *yaml.Node]

- Document
	- Version	string
	- Info	*base.Info
	- Servers	[]*Server
	- Paths	*Paths
	- Components	*Components
	- Security	[]*base.SecurityRequirement
	- Tags	[]*base.Tag
	- ExternalDocs	*base.ExternalDoc
	- Extensions	*[string, *yaml.Node]
	- JsonSchemaDialect	string
	- Webhooks	*[string, *PathItem]
	- Index	*index.SpecIndex
	- Rolodex	*index.Rolodex

- Encoding
	- ContentType	string
	- Headers	*[string, *Header]
	- Style	string
	- Explode	*bool
	- AllowReserved	bool

- OAuthFlow
	- AuthorizationUrl	string
	- TokenUrl	string
	- RefreshUrl	string
	- Scopes	*[string, string]
	- Extensions	*[string, *yaml.Node]

- OAuthFlows
	- Implicit	*OAuthFlow
	- Password	*OAuthFlow
	- ClientCredentials	*OAuthFlow
	- AuthorizationCode	*OAuthFlow
	- Extensions	*[string, *yaml.Node]

- Paths
	- PathItems	*[string, *PathItem]
	- Extensions	*[string, *yaml.Node]

- Callback
	- Expression	*[string, *PathItem]
	- Extensions	*[string, *yaml.Node]

- Components
	- Schemas	*[string, *highbase.SchemaProxy]
	- Responses	*[string, *Response]
	- Parameters	*[string, *Parameter]
	- Examples	*[string, *highbase.Example]
	- RequestBodies	*[string, *RequestBody]
	- Headers	*[string, *Header]
	- SecuritySchemes	*[string, *SecurityScheme]
	- Links	*[string, *Link]
	- Callbacks	*[string, *Callback]
	- PathItems	*[string, *PathItem]
	- Extensions	*[string, *yaml.Node]

- Header
	- Description	string
	- Required	bool
	- Deprecated	bool
	- AllowEmptyValue	bool
	- Style	string
	- Explode	bool
	- AllowReserved	bool
	- Schema	*highbase.SchemaProxy
	- Example	*yaml.Node
	- Examples	*[string, *highbase.Example]
	- Content	*[string, *MediaType]
	- Extensions	*[string, *yaml.Node]

- Link
	- OperationRef	string
	- OperationId	string
	- Parameters	*[string, string]
	- RequestBody	string
	- Description	string
	- Server	*Server
	- Extensions	*[string, *yaml.Node]

- Parameter
	- Name	string
	- In	string
	- Description	string
	- Required	*bool
	- Deprecated	bool
	- AllowEmptyValue	bool
	- Style	string
	- Explode	*bool
	- AllowReserved	bool
	- Schema	*base.SchemaProxy
	- Example	*yaml.Node
	- Examples	*[string, *base.Example]
	- Content	*[string, *MediaType]
	- Extensions	*[string, *yaml.Node]

- PathItem
	- Description	string
	- Summary	string
	- Get	*Operation
	- Put	*Operation
	- Post	*Operation
	- Delete	*Operation
	- Options	*Operation
	- Head	*Operation
	- Patch	*Operation
	- Trace	*Operation
	- Servers	[]*Server
	- Parameters	[]*Parameter
	- Extensions	*[string, *yaml.Node]

- Response
	- Description	string
	- Headers	*[string, *Header]
	- Content	*[string, *MediaType]
	- Links	*[string, *Link]
	- Extensions	*[string, *yaml.Node]

- Server
	- URL	string
	- Description	string
	- Variables	*[string, *ServerVariable]
	- Extensions	*[string, *yaml.Node]
```
