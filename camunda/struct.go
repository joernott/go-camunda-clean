package camunda

type ProcessInstanceList []ProcessInstance

type ProcessInstance struct {
	Links          []UnknownData `json:"links"`
	Id             string        `json:"id"`
	DefinitionId   string        `json:"definitionId"`
	BusinessKey    string        `json:"businessKey"`
	CaseInstanceId string        `json:"caseInstanceId"`
	Ended          bool          `json:"ended"`
	Suspended      bool          `json:"suspended"`
	TenantId       string        `json:"tenantId"`
}

type UnknownData map[string]interface{}
