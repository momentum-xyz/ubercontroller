package universe

import "github.com/momentum-xyz/ubercontroller/utils/mid"

// TODO: redesign all this stuff

var ids = struct {
	systemPluginID mid.ID
	kusamaPluginID mid.ID
}{}

func InitializeIDs(systemPluginID, kusamaPluginID mid.ID) error {
	ids.systemPluginID = systemPluginID
	ids.kusamaPluginID = kusamaPluginID

	return nil
}

func GetSystemPluginID() mid.ID {
	return ids.systemPluginID
}

func GetKusamaPluginID() mid.ID {
	return ids.kusamaPluginID
}
