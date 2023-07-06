package universe

import "github.com/momentum-xyz/ubercontroller/utils/umid"

// TODO: redesign all this stuff

var ids = struct {
	systemPluginID umid.UMID
	kusamaPluginID umid.UMID
}{}

var CustomisableObjectTypeID = umid.MustParse("4ed3a5bb-53f8-4511-941b-079029111111")

func InitializeIDs(systemPluginID, kusamaPluginID umid.UMID) error {
	ids.systemPluginID = systemPluginID
	ids.kusamaPluginID = kusamaPluginID

	return nil
}

func GetSystemPluginID() umid.UMID {
	return ids.systemPluginID
}

func GetKusamaPluginID() umid.UMID {
	return ids.kusamaPluginID
}
