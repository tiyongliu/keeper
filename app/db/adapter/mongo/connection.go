package mongo

import (
	"encoding/json"
	"fmt"
	"keeper/app/db"
	"net/url"
	"strings"
)

const connectionScheme = `mongodb`

// ConnectionURL implements a MongoDB connection struct.
type ConnectionURL struct {
	User     string            `json:"username,omitempty"`
	Password string            `json:"password,omitempty"`
	Host     string            `json:"host"`
	Database string            `json:"database"`
	Port     string            `json:"port"`
	Options  map[string]string `json:"options,omitempty"`
}

func (c ConnectionURL) String() (s string) {
	vv := url.Values{}

	// Do we have any options?
	if c.Options == nil {
		c.Options = map[string]string{}
	}

	// Converting options into URL values.
	for k, v := range c.Options {
		vv.Set(k, v)
	}

	// Has user?
	var userInfo *url.Userinfo

	if c.Port == "" {
		c.Port = "27017"
	}

	if c.User != "" {
		if c.Password == "" {
			userInfo = url.User(c.User)
		} else {
			userInfo = url.UserPassword(c.User, c.Password)
		}
	}

	c.Host = strings.Split(c.Host, ":")[0]
	c.Host = fmt.Sprintf("%s:%s", c.Host, c.Port)

	// Building URL.
	u := url.URL{
		Scheme:   connectionScheme,
		Path:     c.Database,
		Host:     c.Host,
		User:     userInfo,
		RawQuery: vv.Encode(),
	}

	return u.String()
}

func ParseSetting(setting map[string]interface{}) (*ConnectionURL, error) {
	if setting == nil {
		return nil, db.ErrInvalidCollection
	}
	marshal, err := json.Marshal(&setting)
	if err != nil {
		return nil, err
	}
	urlDSN := &ConnectionURL{}
	err = json.Unmarshal(marshal, urlDSN)
	if err != nil {
		return nil, err
	}

	if urlDSN.Host == "" {
		return nil, fmt.Errorf("lack of host")
	}
	return urlDSN, nil
}
