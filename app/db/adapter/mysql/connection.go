package mysql

import (
	"encoding/json"
	"fmt"
	"keeper/app/db"
	"net"
	"net/url"
)

// ConnectionURL implements a MSSQL connection struct.
type ConnectionURL struct {
	User     string            `json:"username"`
	Password string            `json:"password"`
	Port     string            `json:"port"`
	Database string            `json:"database"`
	Host     string            `json:"host"`
	Socket   string            `json:"socket,omitempty"`
	Options  map[string]string `json:"options,omitempty"`
}

func (c ConnectionURL) String() (s string) {
	// Adding username.
	if c.User != "" {
		s = s + c.User
		// Adding password.
		if c.Password != "" {
			s = s + ":" + c.Password
		}
		s = s + "@"
	}

	// Adding protocol and address
	if c.Socket != "" {
		s = s + fmt.Sprintf("unix(%s)", c.Socket)
	} else if c.Host != "" {
		host, port, err := net.SplitHostPort(c.Host)
		if err != nil {
			host = c.Host
			port = "3306"
		}
		s = s + fmt.Sprintf("tcp(%s:%s)", host, port)
	}

	// Adding database
	s = s + "/" + c.Database

	// Do we have any options?
	if c.Options == nil {
		c.Options = map[string]string{}
	}

	// Default options.
	if _, ok := c.Options["charset"]; !ok {
		c.Options["charset"] = "utf8"
	}

	if _, ok := c.Options["parseTime"]; !ok {
		c.Options["parseTime"] = "true"
	}

	// Converting options into URL values.
	vv := url.Values{}

	for k, v := range c.Options {
		vv.Set(k, v)
	}

	// Inserting options.
	if p := vv.Encode(); p != "" {
		s = s + "?" + p
	}

	return s
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

	if urlDSN.User == "" {
		return nil, fmt.Errorf("lack of user")
	}
	if urlDSN.Password == "" {
		return nil, fmt.Errorf("lack of password")
	}
	if urlDSN.Host == "" {
		return nil, fmt.Errorf("lack of host")
	}

	return urlDSN, nil
}
