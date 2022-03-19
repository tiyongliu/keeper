package sideQuests

import (
	"keeper/app/pkg/standard"
)

func readVersion(driver standard.SqlStandard) (*standard.VersionMsg, error) {
	version, err := driver.GetVersion()
	if err != nil {
		return nil, err
	}

	return version, nil
}
