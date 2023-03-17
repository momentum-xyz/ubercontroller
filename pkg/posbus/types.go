package posbus

const (
	MsgTypeSize      = 4
	MsgArrTypeSize   = 4
	MsgUUIDTypeSize  = 16
	MsgLockStateSize = 4
)

type MsgType uint32

/* can use fmt.Sprintf("%x", int) to display hex */
const (
	TypeNONE               MsgType = 0x00000000
	TypeSetUsersTransforms MsgType = 0x285954B8
	TypeMyTransform        MsgType = 0xF878C4BF
	TypeGenericMessage     MsgType = 0xF508E4A3
	TypeHandShake          MsgType = 0x7C41941A
	TypeSetWorld           MsgType = 0xCCDF2E49

	TypeAddObjects     MsgType = 0x2452A9C1
	TypeRemoveObjects  MsgType = 0x6BF88C24
	TypeObjectPosition MsgType = 0xEA6DA4B4

	TypeSetObjectData MsgType = 0xCACE197C

	TypeAddUsers    MsgType = 0xF51F2AFF
	TypeRemoveUsers MsgType = 0xF5A14BB0
	TypeSetUserData MsgType = 0xF702EF5F

	TypeLockObject       MsgType = 0xA7DE9F59
	TypeObjectLockResult MsgType = 0x0924668C

	TypeTriggerVisualEffects MsgType = 0xD96089C6
	TypeUserAction           MsgType = 0xEF1A2E75

	TypeSignal       MsgType = 0xADC1964D
	TypeNotification MsgType = 0xC1FB41D7

	TypeTeleportRequest MsgType = 0x78DA55D9
)
