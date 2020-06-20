package swagger

// SwaggerJson swagger api 文档需要加载的 json 结构定义
type SwaggerJson struct {
	Swagger string     `json:"swagger"` // swagger 的版本
	Info    InfoStruct `json:"info"`    // api 文档的描述信息

	Consumes []string `json:"consumes"`
	Produces []string `json:"produces"`

	Paths       interface{}            `json:"paths"`       // 请求方法的具体信息集合
	Definitions map[string]ModelStruct `json:"definitions"` // 各种 model 数据模型的定义（包括方法的出入参结构定义）
}

// InfoStruct swagger api 中对文档头部包含的文档描述信息的结构定义
type InfoStruct struct {
	Title       string `json:"title"`                 // 该文档的标题
	Description string `json:"description,omitempty"` // 该文档的描述
	Version     string `json:"version,omitempty"`     // 该文档的版本
}

// MethodStruct swagger json 中对方法详细信息的结构定义
type MethodStruct struct {
	Summary     string                    `json:"summary"`               // 方法的注释
	OperationId string                    `json:"operationId"`           // 方法名字
	Responses   map[string]ResponseStruct `json:"responses"`             // 方法的出参结构定义
	Parameters  []*ParametersStruct       `json:"parameters"`            // 方法的入参结构定义
	Tags        []string                  `json:"tags"`                  // 该方法所属的 service
	Description string                    `json:"description,omitempty"` // 方法的描述
}

// ResponseStruct swagger json 中对方法的出参信息的结构定义
type ResponseStruct struct {
	Description string       `json:"description"`      // 对于方法返回的描述
	Schema      SchemaStruct `json:"schema,omitempty"` // 方法出参的数据模型的引用，必须
}

// SchemaStruct swagger json 中对于数据模型使用schema 引用的定义
type SchemaStruct struct {
	Ref string `json:"$ref"`
}

// ParametersStruct swagger json 对方法的入参信息的结构定义
type ParametersStruct struct {
	Name        string        `json:"name"`                  // 参数的名称
	In          string        `json:"in"`                    // 参数的名称
	Required    bool          `json:"required"`              // 参数是否为必须
	Type        string        `json:"type"`                  // 参数的类型
	Schema      *SchemaStruct `json:"schema,omitempty"`      // 参数的引用，非必须
	Format      string        `json:"format,omitempty"`      // 参数的格式，非必须
	Description string        `json:"description,omitempty"` // 参数的描述，非必须
}

// PropertiesStruct swagger json 中对数据模型单个 field 值描述的结构定义
type PropertiesStruct struct {
	Type        string `json:"type,omitempty"`        // 参数的类型
	Format      string `json:"format,omitempty"`      // 参数的格式
	Ref         string `json:"$ref,omitempty"`        // 参数的内包含的引用
	Description string `json:"description,omitempty"` // 参数的描述
}

// ModelStruct swagger json 中对整个数据模型的结构定义
type ModelStruct struct {
	Type       string                      `json:"type"`       // 数据模型的类型
	Properties map[string]PropertiesStruct `json:"properties"` // 数据类型的参数
	Title      string                      `json:"title"`      // 数据类型的描述
}
