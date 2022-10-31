package bridge

import "keeper/app/pkg/serializer"

type Configs struct {
}

type Settings struct {
	UseNativeMenu bool `json:"useNativeMenu"`
}

func (cfg *Configs) GetSettings() *serializer.Response {
	return serializer.SuccessData(serializer.SUCCESS, map[string]interface{}{
		"app": &Settings{},
	})
}

func loadSettings() {

}

func fillMissingSettings() {

}
