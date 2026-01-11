package okx

// SystemStatus 表示系统维护状态（system/status 与 WS status 频道）。
type SystemStatus struct {
	Title        string    `json:"title"`
	State        string    `json:"state"`
	Begin        UnixMilli `json:"begin"`
	End          UnixMilli `json:"end"`
	PreOpenBegin UnixMilli `json:"preOpenBegin"`
	Href         string    `json:"href"`
	ServiceType  string    `json:"serviceType"`
	System       string    `json:"system"`
	ScheDesc     string    `json:"scheDesc"`
	MaintType    string    `json:"maintType"`
	Env          string    `json:"env"`
	TS           UnixMilli `json:"ts"`
}
