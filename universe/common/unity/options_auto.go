package unity

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/momentum-xyz/posbus-protocol/posbus"
	"github.com/momentum-xyz/ubercontroller/config"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/common/api/dto"
	"github.com/momentum-xyz/ubercontroller/utils"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

type AttributeValueChangedMessage struct {
	Type universe.AttributeChangeType     `json:"type"`
	Data AttributeValueChangedMessageData `json:"data"`
}

type AttributeValueChangedMessageData struct {
	AttributeName string `json:"attribute_name"`
	Value         any    `json:"value"`
}

type AttributeSlotKind string

const (
	AttributeSlotKindInvalid AttributeSlotKind = ""
	AttributeSlotKindTexture AttributeSlotKind = "texture"
	AttributeSlotKindVideo   AttributeSlotKind = "video"
	AttributeSlotKindNumber  AttributeSlotKind = "number"
	AttributeSlotKindString  AttributeSlotKind = "string"
	AttributeSlotKindImage   AttributeSlotKind = "image"
)

func GetOptionAutoOption(options *entry.AttributeOptions) (*entry.UnityAutoAttributeOption, error) {
	if options == nil {
		return nil, nil
	}

	autoOptionsValue, ok := (*options)["unity_auto"]
	if !ok {
		return nil, nil
	}

	var autoOption *entry.UnityAutoAttributeOption
	if err := utils.MapDecode(autoOptionsValue, autoOption); err != nil {
		return nil, errors.WithMessage(err, "failed to decode auto option")
	}

	return autoOption, nil
}

func GetOptionAutoMessage(option *entry.UnityAutoAttributeOption, changeType universe.AttributeChangeType,
	attributeID entry.AttributeID, value *entry.AttributeValue) (*websocket.PreparedMessage, error) {

	// do checks if obligatory fields are present
	if option == nil {
		return nil, nil
	}
	if option.SlotKind == "" || option.ContentType == "" {
		return nil, nil
	}

	data := &AttributeValueChangedMessage{
		Type: changeType,
		Data: AttributeValueChangedMessageData{
			AttributeName: attributeID.Name,
			Value:         value,
		},
	}

	// do some checks depending on slot kind
	// for textures we need to render them depending on the
	// content type: currently there are 2 options:
	// video and text/number/string
	switch option.SlotKind {
	case "texture":
		if option.ContentType == "video" {
			payload, err := json.Marshal(map[string]interface{}{
				"url": option.ValueField,
			})
			if err != nil {
				return nil, errors.WithMessage(err, "Failed to marshal preRenderHash")
			}

			hash, err := renderVideo(payload)
			if err != nil {
				return nil, err
			}
			if hash != nil {
				data.Data.Value = hash.Hash
			}
		}

		// if we need to render a string or a number we extract the
		// "string" field of the ValueField
		start := strings.Index(option.ValueField, "string")
		end := strings.Index(option.ValueField[start:], ",")

		preRenderHash := option.ValueField[start : start+end]

		payload, err := json.Marshal(map[string]interface{}{
			"hash": preRenderHash,
		})
		if err != nil {
			return nil, errors.WithMessage(err, "Failed to marshal preRenderHash")
		}

		hash, err := renderFrame(payload)
		if err != nil {
			return nil, err
		}
		if hash != nil {
			data.Data.Value = hash.Hash
		}
	}

	sendData, err := json.Marshal(data)
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to marshal message payload")
	}

	topic := string(option.ContentType)
	if topic == "" {
		topic = attributeID.PluginID.String()
	}

	return posbus.NewRelayToUnityMsg(topic, sendData).WebsocketMessage(), errors.Errorf("got invalid slot kind: %s", option.SlotKind)
}

func renderFrame(preRenderHash []byte) (*dto.HashResponse, error) {
	// need config for the media-manager render URLs

	cfg := config.GetConfig()
	req, err := http.NewRequest("POST", cfg.Common.RenderInternalURL+"/render/addframe", bytes.NewBuffer(preRenderHash))
	if err != nil {
		return nil, errors.WithMessage(err, "Common: renderFrame: failed to create post request")
	}

	req.Header.Set("Content-Type", "image/png")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.WithMessage(err, "Common: renderFrame: failed to post data to media-manager")
	}

	defer resp.Body.Close()

	response := &dto.HashResponse{}

	errs := json.NewDecoder(resp.Body).Decode(response)
	if errs != nil {
		return nil, errors.WithMessage(err, "Common: renderFrame: failed to decode json into response")
	}

	return response, nil
}

func renderVideo(url []byte) (*dto.HashResponse, error) {
	// need config for the media-manager render URLs
	cfg := config.GetConfig()
	req, err := http.NewRequest("POST", cfg.Common.RenderInternalURL+"/render/addtube", bytes.NewBuffer(url))
	if err != nil {
		return nil, errors.WithMessage(err, "Common: renderFrame: failed to create post request")
	}

	req.Header.Set("Content-Type", "image/png")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.WithMessage(err, "Common: renderFrame: failed to post data to media-manager")
	}
	defer resp.Body.Close()

	response := &dto.HashResponse{}

	if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
		return nil, errors.WithMessage(err, "Common: renderFrame: failed to decode json into response")
	}

	return response, nil
}
