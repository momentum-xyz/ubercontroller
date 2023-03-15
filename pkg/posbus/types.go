package posbus

import (
	"encoding/json"
	"github.com/momentum-xyz/ubercontroller/pkg/cmath"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe/logic/api/dto"
	"github.com/momentum-xyz/ubercontroller/utils/mid"
)

type SignalType uint32

const (
	SignalNone SignalType = iota
	SignalDualConnection
	SignalReady
	SignalInvalidToken
	SignalSpawn
	SignalLeaveWorld
	SignalConnectionFailed
	SignalConnected
	SignalConnectionClosed
	SignalWorldDoesNotExist
)

type Trigger uint32

const (
	TriggerNone = iota
	TriggerWow
	TriggerHighFive
	TriggerEnteredObject
	TriggerLeftObject
	TriggerStake
)

type Notification uint32

const (
	NotificationNone     Notification = 0
	NotificationWow      Notification = 1
	NotificationHighFive Notification = 2

	NotificationStageModeAccept        Notification = 10
	NotificationStageModeInvitation    Notification = 11
	NotificationStageModeSet           Notification = 12
	NotificationStageModeStageJoin     Notification = 13
	NotificationStageModeStageRequest  Notification = 14
	NotificationStageModeStageDeclined Notification = 15

	NotificationTextMessage Notification = 500
	NotificationRelay       Notification = 501

	NotificationGeneric Notification = 999
	NotificationLegacy  Notification = 1000
)

const (
	MsgTypeSize      = 4
	MsgArrTypeSize   = 4
	MsgUUIDTypeSize  = 16
	MsgLockStateSize = 4
)

type Destination byte

const (
	DestinationUnity Destination = 0b01
	DestinationReact Destination = 0b10
	DestinationBoth  Destination = 0b11
)

type MsgType uint32

/* can use fmt.Sprintf("%x", int) to display hex */
const (
	NONEType               MsgType = 0x00000000
	TypeSetUsersTransforms MsgType = 0x285954B8
	TypeSendTransform      MsgType = 0xF878C4BF
	TypeGenericMessage     MsgType = 0xF508E4A3
	TypeHandShake          MsgType = 0x7C41941A
	TypeSetWorld           MsgType = 0xCCDF2E49

	TypeAddObjects        MsgType = 0x2452A9C1
	TypeRemoveObjects     MsgType = 0x6BF88C24
	TypeSetObjectPosition MsgType = 0xEA6DA4B4

	TypeSetObjectData MsgType = 0xCACE197C

	TypeAddUsers    MsgType = 0xF51F2AFF
	TypeRemoveUsers MsgType = 0xF5A14BB0
	TypeSetUserData MsgType = 0xF702EF5F

	TypeSetObjectLock    MsgType = 0xA7DE9F59
	TypeObjectLockResult MsgType = 0x0924668C

	TypeTriggerVisualEffects MsgType = 0xD96089C6
	TypeUserAction           MsgType = 0xEF1A2E75

	TypeSignal       MsgType = 0xADC1964D
	TypeNotification MsgType = 0xC1FB41D7

	TypeTeleportRequest MsgType = 0x78DA55D9
)

type HandShake struct {
	HandshakeVersion int    `json:"handshake_version"`
	ProtocolVersion  int    `json:"protocol_version"`
	Token            string `json:"token"`
	UserId           mid.ID `json:"user_id"`
	SessionId        mid.ID `json:"session_id"`
}

type AddObjects struct {
	Objects []ObjectDefinition `json:"objects"`
}

type RemoveObjects struct {
	Objects []mid.ID `json:"objects"`
}

type AddUsers struct {
	Users []UserDefinition `json:"users"`
}

type RemoveUsers struct {
	Users []mid.ID `json:"users"`
}

type Signal struct {
	Value SignalType `json:"value"`
}

type TeleportRequest struct {
	Target mid.ID `json:"target"`
}

type ObjectDefinition struct {
	ID               mid.ID                `json:"mid"`
	ParentID         mid.ID                `json:"parent_id"`
	AssetType        mid.ID                `json:"asset_type"`
	AssetFormat      dto.Asset3dType       `json:"asset_format"` // TODO: Rename AssetType to AssetID, so Type can be used for this.
	Name             string                `json:"name"`
	Transform        cmath.ObjectTransform `json:"transform"`
	IsEditable       bool                  `json:"is_editable"`
	TetheredToParent bool                  `json:"tethered_to_parent"`
	ShowOnMiniMap    bool                  `json:"show_on_minimap"`
	//InfoUI           mid.ID
}

type UserDefinition struct {
	ID        mid.ID              `json:"mid"`
	Name      string              `json:"name"`
	Avatar    mid.ID              `json:"avatar"`
	Transform cmath.UserTransform `json:"transform"`
	IsGuest   bool                `json:"is_guest"`
}

type SetWorldData struct {
	ID              mid.ID `json:"mid"`
	Name            string `json:"name"`
	Avatar          mid.ID `json:"avatar"`
	Owner           mid.ID `json:"owner"`
	Avatar3DAssetID mid.ID `json:"avatar_3d_asset_id"`
}

type ObjectDataIndex struct {
	Kind     entry.UnitySlotType
	SlotName string
}

type ObjectData struct {
	ID      mid.ID
	Entries map[ObjectDataIndex]interface{}
}

type SetObjectLock struct {
	ID    mid.ID `json:"mid"`
	State uint32 `json:"state"`
}

type ObjectLockResultData struct {
	ID        mid.ID `json:"mid"`
	Result    uint32 `json:"result"`
	LockOwner mid.ID `json:"lock_owner"`
}

type ObjectPosition struct {
	ID              mid.ID                `json:"mid"`
	ObjectTransform cmath.ObjectTransform `json:"object_transform"`
}

type Message struct {
	buf     []byte
	msgType MsgType
}

func (o *ObjectData) MarshalJSON() ([]byte, error) {
	q := make(map[string]map[string]interface{})
	for k, v := range o.Entries {
		t, ok := q[string(k.Kind)]
		if !ok {
			t = make(map[string]interface{})
		}
		t[k.SlotName] = v
		q[string(k.Kind)] = t
	}

	return json.Marshal(
		&struct {
			ID      mid.ID                            `json:"mid"`
			Entries map[string]map[string]interface{} `json:"entries"`
		}{
			ID:      o.ID,
			Entries: q,
		},
	)
}
