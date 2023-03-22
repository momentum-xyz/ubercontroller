package posbus

type NotificationType uint32

type Notification struct {
	Value NotificationType
}

const (
	NotificationNone     NotificationType = 0
	NotificationWow      NotificationType = 1
	NotificationHighFive NotificationType = 2

	NotificationStageModeAccept        NotificationType = 10
	NotificationStageModeInvitation    NotificationType = 11
	NotificationStageModeSet           NotificationType = 12
	NotificationStageModeStageJoin     NotificationType = 13
	NotificationStageModeStageRequest  NotificationType = 14
	NotificationStageModeStageDeclined NotificationType = 15

	NotificationTextMessage NotificationType = 500
	NotificationRelay       NotificationType = 501

	NotificationGeneric NotificationType = 999
	NotificationLegacy  NotificationType = 1000
)

func init() {
	registerMessage(Notification{})
	addExtraType(NotificationType(0))
}

func (g *Notification) Type() MsgType {
	return 0xC1FB41D7
}
