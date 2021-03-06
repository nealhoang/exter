// package user contains business object (BO) and data access object (DAO) implementations for User.
package user

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/btnguyen2k/consu/reddo"

	"github.com/btnguyen2k/henge"

	"main/src/utils"
)

// NewUser is helper function to create new User bo.
func NewUser(appVersion uint64, id string) *User {
	user := &User{
		UniversalBo: henge.NewUniversalBo(id, appVersion),
	}
	user.SetAesKey(utils.RandomString(16))
	return user.sync()
}

// NewUserFromUbo is helper function to create User App bo from a universal bo.
func NewUserFromUbo(ubo *henge.UniversalBo) *User {
	if ubo == nil {
		return nil
	}
	user := User{UniversalBo: &henge.UniversalBo{}}
	if err := json.Unmarshal([]byte(ubo.GetDataJson()), &user); err != nil {
		log.Print(fmt.Sprintf("[WARN] NewUserFromUbo - error unmarshalling JSON data: %s", err))
		// log.Print(err)
		return nil
	}
	user.UniversalBo = ubo.Clone()
	return &user
}

const (
	AttrUser_AesKey      = "aes_key"
	AttrUser_DisplayName = "display_name"
	AttrUser_Ubo         = "_ubo"
)

// User is the business object.
// User inherits unique id from bo.UniversalBo. Email address is used to uniquely identify user (e.g. user-id is email address).
type User struct {
	*henge.UniversalBo `json:"_ubo"`
}

// MarshalJSON implements json.encode.Marshaler.MarshalJSON.
//	TODO: lock for read?
func (user *User) MarshalJSON() ([]byte, error) {
	user.sync()
	m := map[string]interface{}{
		AttrUser_Ubo: user.UniversalBo.Clone(),
	}
	return json.Marshal(m)
}

// UnmarshalJSON implements json.decode.Unmarshaler.UnmarshalJSON.
//	TODO: lock for write?
func (user *User) UnmarshalJSON(data []byte) error {
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	if m[AttrUser_Ubo] != nil {
		js, _ := json.Marshal(m[AttrUser_Ubo])
		if err := json.Unmarshal(js, &user.UniversalBo); err != nil {
			return err
		}
	}
	user.sync()
	return nil
}

// GetAesKey returns value of user's 'aes-key' attribute.
func (user *User) GetAesKey() string {
	v, err := user.GetDataAttrAs(AttrUser_AesKey, reddo.TypeString)
	if err != nil || v == nil {
		return ""
	}
	return v.(string)
}

// SetAesKey sets value of user's 'aes-key' attribute.
func (user *User) SetAesKey(v string) *User {
	user.SetDataAttr(AttrUser_AesKey, strings.TrimSpace(v))
	return user
}

// GetDisplayName returns value of user's 'display-name' attribute.
// available since v0.4.0
func (user *User) GetDisplayName() string {
	v, err := user.GetDataAttrAs(AttrUser_DisplayName, reddo.TypeString)
	if err != nil || v == nil {
		return ""
	}
	return v.(string)
}

// SetDisplayName sets value of user's 'display-name' attribute.
// available since v0.4.0
func (user *User) SetDisplayName(v string) *User {
	user.SetDataAttr(AttrUser_DisplayName, strings.TrimSpace(v))
	return user
}

func (user *User) sync() *User {
	user.UniversalBo.Sync()
	return user
}
