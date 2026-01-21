package okx

// UsersSubaccount 表示子账户信息（Users 模块）。
type UsersSubaccount struct {
	Type       string   `json:"type"`
	Enable     bool     `json:"enable"`
	SubAcct    string   `json:"subAcct"`
	UID        string   `json:"uid"`
	Label      string   `json:"label"`
	Mobile     string   `json:"mobile"`
	GAuth      bool     `json:"gAuth"`
	FrozenFunc []string `json:"frozenFunc"`

	CanTransOut    bool   `json:"canTransOut"`
	TS             string `json:"ts"`
	SubAcctLv      string `json:"subAcctLv"`
	FirstLvSubAcct string `json:"firstLvSubAcct"`
	IfDma          bool   `json:"ifDma"`
}

// UsersSubaccountCreateSubaccountResult 表示创建子账户的返回项。
type UsersSubaccountCreateSubaccountResult struct {
	SubAcct string `json:"subAcct"`
	Label   string `json:"label"`
	UID     string `json:"uid"`
	TS      string `json:"ts"`
}

// UsersSubaccountAPIKeyCreateResult 表示创建子账户 API Key 的返回项（会返回 secretKey）。
type UsersSubaccountAPIKeyCreateResult struct {
	SubAcct    string          `json:"subAcct"`
	Label      string          `json:"label"`
	APIKey     string          `json:"apiKey"`
	SecretKey  SensitiveString `json:"secretKey"`
	Passphrase SensitiveString `json:"passphrase"`
	Perm       string          `json:"perm"`
	IP         string          `json:"ip"`
	TS         string          `json:"ts"`
}

// UsersSubaccountAPIKeyInfo 表示子账户 API Key 信息（查询/重置返回项）。
type UsersSubaccountAPIKeyInfo struct {
	SubAcct string `json:"subAcct,omitempty"`
	Label   string `json:"label"`
	APIKey  string `json:"apiKey"`
	Perm    string `json:"perm"`
	IP      string `json:"ip"`
	TS      string `json:"ts"`
}

// UsersSubaccountDeleteAPIKeyResult 表示删除子账户 API Key 的返回项。
type UsersSubaccountDeleteAPIKeyResult struct {
	SubAcct string `json:"subAcct"`
}

// UsersSubaccountTransferOutPermission 表示子账户主动转出权限返回项。
type UsersSubaccountTransferOutPermission struct {
	SubAcct     string `json:"subAcct"`
	CanTransOut bool   `json:"canTransOut"`
}

// UsersEntrustSubaccount 表示被托管的子账户信息返回项。
type UsersEntrustSubaccount struct {
	SubAcct string `json:"subAcct"`
}
