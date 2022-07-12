package sideQuests

import (
	"keeper/app/modules"
	"keeper/app/pkg/standard"
)

func readVersion(ch chan *modules.EchoMessage, pool standard.SqlStandard) error {
	version, err := pool.GetVersion()
	if err != nil {
		return err
	}

	ch <- &modules.EchoMessage{
		Payload: version,
		MsgType: "version",
	}

	return nil
}
