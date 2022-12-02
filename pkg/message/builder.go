package message

import (
	"encoding/binary"

	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"github.com/momentum-xyz/posbus-protocol/flatbuff/go/api"
	"github.com/momentum-xyz/posbus-protocol/posbus"
	"github.com/momentum-xyz/ubercontroller/logger"
	"github.com/momentum-xyz/ubercontroller/pkg/cmath"
	"github.com/momentum-xyz/ubercontroller/universe/common/api/dto"
	"github.com/momentum-xyz/ubercontroller/utils"
)

type ObjectDefinition struct {
	ObjectID         uuid.UUID
	ParentID         uuid.UUID
	AssetType        uuid.UUID
	AssetFormat      dto.Asset3dType // TODO: Rename AssetType to AssetID, so Type can be used for this.
	Name             string
	Position         cmath.SpacePosition
	TetheredToParent bool
	Minimap          bool
	InfoUI           uuid.UUID
}

type DecorationMetadata struct {
	AssetID  uuid.UUID  `json:"AssetID"`
	Position cmath.Vec3 `json:"Position"`
	rotation cmath.Vec3
}

type Builder struct {
	builders chan *flatbuffers.Builder
}

var msgBuilder *Builder

func GetBuilder() *Builder {
	return msgBuilder
}

var log = logger.L()

func NewSendPosBuffer(id uuid.UUID) []byte {
	buf := make([]byte, posbus.UserPositionsMessageSize)
	copy(buf[:posbus.MsgUUIDTypeSize], utils.BinID(id))
	return buf
}

func (mb *Builder) NewPosBusMessageBuffer(msgId uint32, len int) []byte {
	msg := make([]byte, posbus.MsgTypeSize*2+len)
	binary.LittleEndian.PutUint32(msg, msgId)
	binary.LittleEndian.PutUint32(msg[posbus.MsgTypeSize+len:], ^msgId)
	return msg
}

func InitBuilder(count int, size int) {
	msgBuilder = &Builder{
		builders: make(chan *flatbuffers.Builder, count),
	}

	for i := 0; i < count; i++ {
		msgBuilder.builders <- flatbuffers.NewBuilder(size)
	}
	//return msgBuilder
}

func (mb *Builder) GetBuilder() *flatbuffers.Builder {
	builder := <-mb.builders
	return builder
}

func (mb *Builder) ReleaseBuilder(builder *flatbuffers.Builder) {
	builder.Reset()
	mb.builders <- builder
}

func (mb *Builder) MsgSetWorld(
	worldId uuid.UUID, name string, avatarControllerId, skyboxControllerId uuid.UUID, lodDistances []uint32,
	decorations []DecorationMetadata,
) *websocket.PreparedMessage {
	builder := mb.GetBuilder()
	defer func() {
		mb.ReleaseBuilder(builder)
	}()
	nameObj := builder.CreateString(name)

	api.SetWorldStartLodDistancesVector(builder, len(lodDistances))
	for i := len(lodDistances) - 1; i >= 0; i-- {
		builder.PrependUint32(lodDistances[i])
	}
	LODsOffset := builder.EndVector(len(lodDistances))
	var decorationOffsets []flatbuffers.UOffsetT
	for i := len(decorations) - 1; i >= 0; i-- {
		pos := decorations[i].Position
		rot := decorations[i].rotation
		api.DecorationMetadataStart(builder)
		api.DecorationMetadataAddAssetId(builder, mb.SerializeGUID(builder, decorations[i].AssetID))
		api.DecorationMetadataAddPos(
			builder, api.CreatePosition(builder, pos.X, pos.Y, pos.Z, rot.X, rot.Y, rot.Z, 0, 0, 0),
		)
		offset := api.DecorationMetadataEnd(builder)
		decorationOffsets = append(decorationOffsets, offset)
	}

	api.SetWorldStartDecorationsVector(builder, len(decorationOffsets))
	for _, decorationOffset := range decorationOffsets {
		builder.PrependUOffsetT(decorationOffset)
	}
	decsOffset := builder.EndVector(len(decorationOffsets))

	api.SetWorldStart(builder)
	api.SetWorldAddWorldId(builder, mb.SerializeGUID(builder, worldId))
	api.SetWorldAddAvatarControllerId(builder, mb.SerializeGUID(builder, avatarControllerId))
	api.SetWorldAddSkyboxControllerId(builder, mb.SerializeGUID(builder, skyboxControllerId))
	api.SetWorldAddLodDistances(builder, LODsOffset)
	api.SetWorldAddDecorations(builder, decsOffset)
	api.SetWorldAddName(builder, nameObj)
	msgOffset := api.SetWorldEnd(builder)

	return mb.FinishMessage(builder, api.MsgSetWorld, msgOffset)
}

func (mb *Builder) FinishMessage(
	builder *flatbuffers.Builder, msgType api.Msg, msgOffset flatbuffers.UOffsetT,
) *websocket.PreparedMessage {
	api.FlatBuffMsgStart(builder)
	api.FlatBuffMsgAddMsgType(builder, msgType)
	api.FlatBuffMsgAddMsg(builder, msgOffset)
	flatBuffMsgOffset := api.FlatBuffMsgEnd(builder)
	builder.Finish(flatBuffMsgOffset)
	rawBytes := builder.FinishedBytes()

	return mb.WrapMessage(rawBytes)
}

// FinishMessageBytes is used when we want to skip gorilla/websocket.PreparedMessage type
func (mb *Builder) FinishMessageBytes(
	builder *flatbuffers.Builder, msgType api.Msg, msgOffset flatbuffers.UOffsetT,
) (msg, buf []byte) {
	api.FlatBuffMsgStart(builder)
	api.FlatBuffMsgAddMsgType(builder, msgType)
	api.FlatBuffMsgAddMsg(builder, msgOffset)
	flatBuffMsgOffset := api.FlatBuffMsgEnd(builder)
	builder.Finish(flatBuffMsgOffset)
	rawBytes := builder.FinishedBytes()

	m := posbus.NewMessage(posbus.MsgTypeFlatBufferMessage, len(rawBytes))
	copy(m.Msg(), rawBytes)
	return m.Msg(), m.Buf()
}

func (mb *Builder) SerializeGUID(builder *flatbuffers.Builder, uuid uuid.UUID) flatbuffers.UOffsetT {
	swappedBytes := unityUUID(uuid[:])
	least := binary.LittleEndian.Uint64(swappedBytes[0:8])
	most := binary.LittleEndian.Uint64(swappedBytes[8:])
	return api.CreateID(builder, least, most)
}

func (mb *Builder) MsgObjectDefinition(obj ObjectDefinition) *websocket.PreparedMessage {
	builder := mb.GetBuilder()
	defer func() {
		mb.ReleaseBuilder(builder)
		log.Debug("Builder: MsgObjectDefinition")
	}()

	objName := builder.CreateString(obj.Name)

	api.ObjectDefinitionStart(builder)
	api.ObjectDefinitionAddObjectId(builder, mb.SerializeGUID(builder, obj.ObjectID))
	api.ObjectDefinitionAddName(builder, objName)
	api.ObjectDefinitionAddPosition(
		builder, api.CreatePosition(
			builder, obj.Position.Location.X, obj.Position.Location.Y, obj.Position.Location.Z, obj.Position.Rotation.X,
			obj.Position.Rotation.Y, obj.Position.Rotation.Z, obj.Position.Scale.X, obj.Position.Scale.Y,
			obj.Position.Scale.Z,
		),
	)
	api.ObjectDefinitionAddParentId(builder, mb.SerializeGUID(builder, obj.ParentID))
	api.ObjectDefinitionAddAssetType(builder, mb.SerializeGUID(builder, obj.AssetType))
	api.ObjectDefinitionAddTetheredToParent(builder, obj.TetheredToParent)
	api.ObjectDefinitionAddMinimap(builder, obj.Minimap)
	api.ObjectDefinitionAddInfouiType(builder, mb.SerializeGUID(builder, obj.InfoUI))
	api.ObjectDefinitionAddIsGltf(builder, obj.AssetFormat == dto.GLTFAsset3dType)
	msgOffset := api.ObjectDefinitionEnd(builder)

	return mb.FinishMessage(builder, api.MsgObjectDefinition, msgOffset)
}

func (mb *Builder) MsgAddStaticObjects(objects []ObjectDefinition) *websocket.PreparedMessage {
	builder := mb.GetBuilder()
	defer func() {
		mb.ReleaseBuilder(builder)
		log.Debug("MSG: msgAddStaticObjects")
	}()

	var objectOffsets []flatbuffers.UOffsetT
	for i := len(objects) - 1; i >= 0; i-- {
		obj := objects[i]
		nameObj := builder.CreateString(obj.Name)

		api.ObjectDefinitionStart(builder)
		api.ObjectDefinitionAddObjectId(builder, mb.SerializeGUID(builder, obj.ObjectID))
		api.ObjectDefinitionAddName(builder, nameObj)
		api.ObjectDefinitionAddPosition(
			builder,
			api.CreatePosition(
				builder, obj.Position.Location.X, obj.Position.Location.Y, obj.Position.Location.Z,
				obj.Position.Rotation.X,
				obj.Position.Rotation.Y, obj.Position.Rotation.Z, obj.Position.Scale.X, obj.Position.Scale.Y,
				obj.Position.Scale.Z,
			),
		)
		api.ObjectDefinitionAddParentId(builder, mb.SerializeGUID(builder, obj.ParentID))
		api.ObjectDefinitionAddAssetType(builder, mb.SerializeGUID(builder, obj.AssetType))
		api.ObjectDefinitionAddTetheredToParent(builder, obj.TetheredToParent)
		api.ObjectDefinitionAddMinimap(builder, obj.Minimap)
		api.ObjectDefinitionAddInfouiType(builder, mb.SerializeGUID(builder, obj.InfoUI))
		offsetObj := api.ObjectDefinitionEnd(builder)
		objectOffsets = append(objectOffsets, offsetObj)
	}

	api.AddStaticObjectsStartObjectsVector(builder, len(objectOffsets))
	for _, objOffset := range objectOffsets {
		builder.PrependUOffsetT(objOffset)
	}
	offsetObjects := builder.EndVector(len(objectOffsets))

	api.AddStaticObjectsStart(builder)
	api.AddStaticObjectsAddObjects(builder, offsetObjects)
	msgOffset := api.AddStaticObjectsEnd(builder)

	return mb.FinishMessage(builder, api.MsgAddStaticObjects, msgOffset)
}

func (mb *Builder) SetObjectTextures(id uuid.UUID, textures map[string]string) *websocket.PreparedMessage {
	builder := mb.GetBuilder()
	defer func() {
		mb.ReleaseBuilder(builder)
		log.Debug("MSG: makeSetObjectTextures")
	}()

	var objectOffsets []flatbuffers.UOffsetT
	for k, v := range textures {
		objLabel := builder.CreateString(k)
		objData := builder.CreateString(v)
		api.TextureDefinitionStart(builder)
		api.TextureDefinitionAddLabel(builder, objLabel)
		api.TextureDefinitionAddData(builder, objData)
		offsetObj := api.TextureDefinitionEnd(builder)
		objectOffsets = append(objectOffsets, offsetObj)
	}

	api.SetObjectTexturesStartObjectsVector(builder, len(objectOffsets))
	for _, objOffset := range objectOffsets {
		builder.PrependUOffsetT(objOffset)
	}
	offsetObjects := builder.EndVector(len(objectOffsets))

	api.SetObjectTexturesStart(builder)
	api.SetObjectTexturesAddObjectId(builder, mb.SerializeGUID(builder, id))
	api.SetObjectTexturesAddObjects(builder, offsetObjects)

	msgOffset := api.SetObjectTexturesEnd(builder)

	return mb.FinishMessage(builder, api.MsgSetObjectTextures, msgOffset)
}

func (mb *Builder) SetObjectAttributes(id uuid.UUID, attributes map[string]int32) *websocket.PreparedMessage {
	builder := mb.GetBuilder()
	defer func() {
		mb.ReleaseBuilder(builder)
		log.Debug("MSG: makeSetObjectAttributes")
	}()

	var objectOffsets []flatbuffers.UOffsetT
	for k, v := range attributes {
		objLabel := builder.CreateString(k)
		// objData := builder.CreateIn(v)
		api.AttributeDefinitionStart(builder)
		api.AttributeDefinitionAddLabel(builder, objLabel)
		api.AttributeDefinitionAddAttribute(builder, v)
		offsetObj := api.AttributeDefinitionEnd(builder)
		objectOffsets = append(objectOffsets, offsetObj)
	}

	api.SetObjectAttributesStartObjectsVector(builder, len(objectOffsets))
	for _, objOffset := range objectOffsets {
		builder.PrependUOffsetT(objOffset)
	}
	offsetObjects := builder.EndVector(len(objectOffsets))

	api.SetObjectAttributesStart(builder)
	api.SetObjectAttributesAddSpaceId(builder, mb.SerializeGUID(builder, id))
	api.SetObjectAttributesAddObjects(builder, offsetObjects)

	msgOffset := api.SetObjectAttributesEnd(builder)

	return mb.FinishMessage(builder, api.MsgSetObjectAttributes, msgOffset)
}

func (mb *Builder) SetObjectStrings(id uuid.UUID, strings map[string]string) *websocket.PreparedMessage {
	builder := mb.GetBuilder()
	defer func() {
		mb.ReleaseBuilder(builder)
		log.Debug("MSG: makeSetObjectStrings")
	}()

	var objectOffsets []flatbuffers.UOffsetT
	for k, v := range strings {
		objLabel := builder.CreateString(k)
		objData := builder.CreateString(v)
		api.StringDefinitionStart(builder)
		api.StringDefinitionAddLabel(builder, objLabel)
		api.StringDefinitionAddData(builder, objData)
		offsetObj := api.StringDefinitionEnd(builder)
		objectOffsets = append(objectOffsets, offsetObj)
	}

	api.SetObjectStringsStartObjectsVector(builder, len(objectOffsets))
	for _, objOffset := range objectOffsets {
		builder.PrependUOffsetT(objOffset)
	}
	offsetObjects := builder.EndVector(len(objectOffsets))

	api.SetObjectStringsStart(builder)
	api.SetObjectStringsAddObjectId(builder, mb.SerializeGUID(builder, id))
	api.SetObjectStringsAddObjects(builder, offsetObjects)

	msgOffset := api.SetObjectStringsEnd(builder)

	return mb.FinishMessage(builder, api.MsgSetObjectStrings, msgOffset)
}

func (mb *Builder) WrapMessage(data []byte) *websocket.PreparedMessage {
	msg := posbus.NewMessage(posbus.MsgTypeFlatBufferMessage, len(data))
	copy(msg.Msg(), data)
	return msg.WebsocketMessage()
}

func unityUUID(uuid []byte) []byte {
	uuid1 := make([]byte, 16)
	copy(uuid1, uuid)

	uuid1[3], uuid1[0] = uuid1[0], uuid1[3]
	uuid1[2], uuid1[1] = uuid1[1], uuid1[2]
	uuid1[4], uuid1[5] = uuid1[5], uuid1[4]
	uuid1[6], uuid1[7] = uuid1[7], uuid1[6]
	return uuid1
}

func DeserializeGUID(id *api.ID) uuid.UUID {
	rawId := make([]byte, 16)
	binary.LittleEndian.PutUint64(rawId, id.L())
	binary.LittleEndian.PutUint64(rawId[8:], id.M())
	uuidRes, _ := uuid.FromBytes(unityUUID(rawId))
	return uuidRes
}
