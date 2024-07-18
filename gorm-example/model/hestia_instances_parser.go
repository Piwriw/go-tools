package model

import (
	"encoding/json"
	"errors"

	"github.com/go-playground/validator/v10"
	icarussdk "yunqutech.gitlab.com/agilex/basis/icarus"
	"yunqutech.gitlab.com/agilex/basis/utils/encryption"
)
import (
	"time"
)

const TableNameHestiaInstance = "hestia_instances"

// HestiaInstance mapped from table <hestia_instances>
type HestiaInstance struct {
	ID             int64     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Category       string    `gorm:"column:category;not null;default:[]" json:"category"`
	Usefor         string    `gorm:"column:usefor;not null;default:[]" json:"usefor"`
	Datum          string    `gorm:"column:data;not null;default:{}" json:"data"`
	Status         bool      `gorm:"column:status;not null;default:true" json:"status"`
	CreatedAt      time.Time `gorm:"column:created_at" json:"created_at"`
	CreatedBy      string    `gorm:"column:created_by;not null" json:"created_by"`
	UpdatedAt      time.Time `gorm:"column:updated_at" json:"updated_at"`
	UpdatedBy      string    `gorm:"column:updated_by;not null" json:"updated_by"`
	DeletedBy      string    `gorm:"column:deleted_by;not null" json:"deleted_by"`
	RequestID      string    `gorm:"column:request_id;not null" json:"request_id"`
	Account        string    `gorm:"column:account;not null;default:{}" json:"account"`
	DisabledBy     string    `gorm:"column:disabled_by;not null" json:"disabled_by"`
	MonitorMethod  string    `gorm:"column:monitor_method;not null" json:"monitor_method"`
	MonitorAccount string    `gorm:"column:monitor_account;not null" json:"monitor_account"`
	MonitorStatus  string    `gorm:"column:monitor_status;not null" json:"monitor_status"`
	MonitorError   string    `gorm:"column:monitor_error;not null" json:"monitor_error"`
	MonitorScrape  time.Time `gorm:"column:monitor_scrape" json:"monitor_scrape"`
	Source         int32     `gorm:"column:source;not null" json:"source"`
	PoolMaxActive  int32     `gorm:"column:pool_max_active;default:10" json:"pool_max_active"`
	QueryTimeout   int32     `gorm:"column:query_timeout;default:300" json:"query_timeout"`
	Network        string    `gorm:"column:network;not null;default:server1" json:"network"`
}

// TableName HestiaInstance's table name
func (*HestiaInstance) TableName() string {
	return TableNameHestiaInstance
}

type ConnectionInfo struct {
	IP            string `validate:"required,ip"`
	DBName        string `validate:"required"`
	DBPort        int    `validate:"required"`
	DBUser        string `validate:"required"`
	DBPassword    string `validate:"required"`
	DBVersion     string `validate:"required"`
	PoolMaxActive int    `validate:"required"`
	QueryTimeout  int    `validate:"required"`
	Encoding      string `validate:""`
}

type InstanceInfo struct {
	ID             int    `validate:"required"`
	Name           string `validate:"required"`
	DBType         string `validate:"required"`
	DBVersion      string `validate:"required"`
	ResourceType   string `validate:"required"`
	Network        string `validate:"required"`
	UseFor         []string
	ConnectionInfo ConnectionInfo
	// internal fields
	dbInstance HestiaInstanceModel
}

func (i *InstanceInfo) FromModel(hestiaInstance *HestiaInstanceModel) error {
	validate := validator.New()
	i.ID = hestiaInstance.ID
	i.DBType = hestiaInstance.Category.InstanceType
	i.ResourceType = hestiaInstance.Category.ResourceType
	i.dbInstance = *hestiaInstance
	i.DBVersion = hestiaInstance.Account.InstanceAccount.Content.DbVersion
	i.Name = hestiaInstance.Data.Name
	i.Network = hestiaInstance.Network
	i.UseFor = []string{hestiaInstance.UseFor.ResourceType, hestiaInstance.UseFor.InstanceType}
	i.ConnectionInfo.IP = hestiaInstance.Data.IP
	i.ConnectionInfo.DBName = hestiaInstance.Account.InstanceAccount.Content.DbName
	i.ConnectionInfo.DBPort = hestiaInstance.Account.InstanceAccount.Content.DbPort
	i.ConnectionInfo.DBUser = hestiaInstance.Account.InstanceAccount.Content.DbUsername
	passwordDecode, err := encryption.Decode(hestiaInstance.Account.InstanceAccount.Content.DbPassword)
	if err != nil {
		passwordDecode = hestiaInstance.Account.InstanceAccount.Content.DbPassword
	}
	i.ConnectionInfo.DBPassword = passwordDecode
	i.ConnectionInfo.PoolMaxActive = hestiaInstance.PoolMaxActive
	i.ConnectionInfo.QueryTimeout = hestiaInstance.QueryTimeout

	i.ConnectionInfo.DBVersion = hestiaInstance.Account.InstanceAccount.Content.DbVersion
	i.ConnectionInfo.Encoding = hestiaInstance.Account.InstanceAccount.Content.DbEncoding
	return validate.Struct(i)
}

type HestiaInstanceModel struct {
	ID            int             `json:"id" gorm:"primary_key;not null"`
	UseFor        UseFor          `json:"usefor" gorm:"serializer:json;type:json;not null;" validate:"required"`
	Category      Category        `json:"category" gorm:"column:category;serializer:json;type:json;not null;" validate:"required"`
	Account       InstanceAccount `json:"account" gorm:"serializer:json;type:json;not null;"`
	Data          InstanceData    `json:"data" gorm:"serializer:json;type:json;not null;"`
	Status        bool            `json:"status" gorm:"type:boolean;not null;default:true"`
	Network       string          `json:"network" gorm:"type:varchar(255);default:''"`
	PoolMaxActive int             `json:"pool_max_active" gorm:"type:int" validate:"required,min=5,max=15"`  // 连接池最大连接数
	QueryTimeout  int             `json:"query_timeout" gorm:"type:int" validate:"required,min=120,max=300"` // 查询超时时间

}

func (*HestiaInstanceModel) TableName() string {
	return TableNameHestiaInstance
}

type UseFor struct {
	ResourceType string `json:"UseForResourceType"`
	InstanceType string `json:"UseForInstanceType"`
}

func (u *UseFor) UnmarshalJSON(data []byte) error {
	var v []interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	if len(v) != 2 {
		return errors.New("invalid Category json")
	}
	u.ResourceType = v[0].(string)
	u.InstanceType = v[1].(string)
	return nil
}

func (u *UseFor) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	// Unmarshal the json.RawMessage into an InstanceData struct
	switch v := value.(type) {
	case []byte:
		if err := json.Unmarshal(v, u); err != nil {
			return err
		}
	case string:
		if err := json.Unmarshal([]byte(v), u); err != nil {
			return err
		}
	default:
		if err := json.Unmarshal(v.([]byte), u); err != nil {
			return err
		}
	}
	return nil
}

type Category struct {
	ResourceType string `json:"ResourceType"`
	InstanceType string `json:"InstanceType"`
}

func (c *Category) UnmarshalJSON(data []byte) error {
	var v []interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	if len(v) != 2 {
		return errors.New("invalid Category json")
	}
	c.ResourceType = v[0].(string)
	c.InstanceType = v[1].(string)
	return nil
}

func (c *Category) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	// Unmarshal the json.RawMessage into an InstanceData struct
	switch v := value.(type) {
	case []byte:
		if err := json.Unmarshal(v, c); err != nil {
			return err
		}
	case string:
		if err := json.Unmarshal([]byte(v), c); err != nil {
			return err
		}
	default:
		if err := json.Unmarshal(v.([]byte), c); err != nil {
			return err
		}
	}
	return nil
}

type InstanceData struct {
	IP            string          `json:"ip"`
	Name          string          `json:"name"`
	Inline        bool            `json:"inline"`
	Method        string          `json:"method"`
	Business      json.RawMessage `json:"business"`
	MonitorStatus bool            `json:"monitor_status"`
}

func (i *InstanceData) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	// Unmarshal the json.RawMessage into an InstanceData struct
	switch v := value.(type) {
	case []byte:
		if err := json.Unmarshal(v, i); err != nil {
			return err
		}
	case string:
		if err := json.Unmarshal([]byte(v), i); err != nil {
			return err
		}
	default:
		if err := json.Unmarshal(v.([]byte), i); err != nil {
			return err
		}
	}
	return nil
}

type InstanceAccount struct {
	InstanceType    string `json:"instanceType"`
	InstanceAccount *icarussdk.Account
}

func (a *InstanceAccount) UnmarshalJSON(data []byte) error {
	var v []interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	a.InstanceType = v[0].(string)
	// Unmarshal the account json.RawMessage into an InstanceAccount struct
	accountBytes, err := json.Marshal(v[1])
	if err != nil {
		return err
	}
	if err := json.Unmarshal(accountBytes, &a.InstanceAccount); err != nil {
		return err
	}

	// Decrypt the password
	var decryptPassword = func(password string) (string, error) {
		if password, err := encryption.Decode(password); err != nil {
			return password, err
		}
		return password, nil
	}

	content := a.InstanceAccount.Content
	if content.DbPassword, err = decryptPassword(content.DbPassword); err != nil {
		return err
	}

	if content.SshPassword, err = decryptPassword(content.SshPassword); err != nil {
		return err
	}
	return nil
}

func (a *InstanceAccount) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	// Unmarshal the json.RawMessage into an InstanceData struct
	switch v := value.(type) {
	case []byte:
		if err := json.Unmarshal(v, a); err != nil {
			return err
		}
	case string:
		if err := json.Unmarshal([]byte(v), a); err != nil {
			return err
		}
	default:
		if err := json.Unmarshal(v.([]byte), a); err != nil {
			return err
		}
	}
	return nil
}

func (a *InstanceAccount) Value() (interface{}, error) {
	return json.Marshal(a)
}
