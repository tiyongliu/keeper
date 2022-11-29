package sideQuests

import (
	"keeper/app/db"
	"keeper/app/db/standard/modules"
)

func readVersion(driver db.Session) (*modules.Version, error) {
	version, err := driver.Version()
	if err != nil {
		return nil, err
	}

	return version, nil
}
