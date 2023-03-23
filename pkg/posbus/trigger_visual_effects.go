package posbus

// FIXME: finalize VisualEffect type with all data
type VisualEffect struct {
	Name string `json:"name"`
}

type TriggerVisualEffects struct {
	Effects []VisualEffect `json:"effects"`
}

func init() {
	registerMessage(TriggerVisualEffects{})
	addExtraType(VisualEffect{})
}

func (g *TriggerVisualEffects) GetType() MsgType {
	return 0xD96089C6
}
