package universe

import "github.com/google/uuid"

// TODO: redesign all this stuff

var ids = struct {
	systemPluginID uuid.UUID
	kusamaPluginID uuid.UUID
}{}

func InitializeIDs(systemPluginID, kusamaPluginID uuid.UUID) error {
	ids.systemPluginID = systemPluginID
	ids.kusamaPluginID = kusamaPluginID

	return nil
}

func GetSystemPluginID() uuid.UUID {
	return ids.systemPluginID
}

func GetKusamaPluginID() uuid.UUID {
	return ids.kusamaPluginID
}