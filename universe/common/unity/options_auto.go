package unity

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/momentum-xyz/ubercontroller/config"
	"github.com/momentum-xyz/ubercontroller/types"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/common/api/dto"
	"github.com/momentum-xyz/ubercontroller/utils"

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

const (
	RenderKindNone uint = iota
	RenderKindVideo
	RenderKindText
)

// TODO: requre optimization to not do it every time we touch attribute
func GetOptionAutoOption(
	attributeID entry.AttributeID, options *entry.AttributeOptions,
) (*entry.UnityAutoAttributeOption, error) {
	if options == nil {
		return nil, nil
	}

	autoOptionsValue, ok := (*options)["unity_auto"]
	if !ok {
		return nil, nil
	}

	//fmt.Printf("FFF0: %+v \n", autoOptionsValue)

	var autoOption *entry.UnityAutoAttributeOption
	if err := utils.MapDecode(autoOptionsValue, &autoOption); err != nil {
		return nil, errors.WithMessage(err, "failed to decode auto option")
	}
	//fmt.Printf("FFF: %+v %+v \n", autoOption, autoOptionsValue)

	if autoOption.SlotType == entry.UnitySlotTypeInvalid || autoOption.ContentType == entry.UnityContentTypeInvalid {
		return nil, nil
	}
	if autoOption.SlotName == "" {
		autoOption.SlotName = attributeID.Name
	}
	if autoOption.ValueField == "" {
		if autoOption.ContentType == "image" {
			autoOption.ValueField = "render_hash"
		} else {
			autoOption.ValueField = "value"
		}
	}
	return autoOption, nil
}

func PrerenderAutoValue(
	ctx context.Context, option *entry.UnityAutoAttributeOption, value *entry.AttributeValue,
) (*dto.HashResponse, error) {
	if option == nil || option.SlotType != "texture" || value == nil {
		return nil, nil
	}

	if option.ContentType == "image" {
		return nil, nil
	}

	valueAny, ok := (*value)[option.ValueField]
	if !ok {
		return nil, nil
	}
	var valueString string

	var renderKind uint
	renderKind = RenderKindText
	if option.ContentType == "video" {
		renderKind = RenderKindVideo
	}
	switch option.ContentType {
	case "video", "text", "string":
		valueString, ok = valueAny.(string)
		if !ok {
			errors.New("Can not cast value to string in PrerenderAutoValue")
		}
	case "number":
		valueUint, ok := valueAny.(uint32)
		if !ok {
			errors.New("Can not cast value to uint32 in PrerenderAutoValue")
		}
		valueString = strconv.FormatUint(uint64(valueUint), 10)
	default:
		return nil, nil
	}

	var hash *dto.HashResponse
	var err error

	switch renderKind {
	case RenderKindVideo:
		payload, err := json.Marshal(
			map[string]any{
				"url": valueString,
			},
		)
		if err != nil {
			return nil, errors.WithMessage(err, "Failed to marshal preRenderHash")
		}

		hash, err = renderVideo(ctx, payload)
		if err != nil {
			return nil, err
		}
	case RenderKindText:
		if option.TextRenderTemplate == "" {
			return nil, nil
		}

		payload := []byte(strings.Replace(option.TextRenderTemplate, "%TEXT%", valueString, -1))

		hash, err = renderFrame(ctx, payload)
		if err != nil {
			return nil, err
		}
	}
	return hash, nil
}

func renderFrame(ctx context.Context, textJob []byte) (*dto.HashResponse, error) {
	// need config for the media-manager render URLs
	cfg := utils.GetFromAny(ctx.Value(types.ConfigContextKey), (*config.Config)(nil))
	if cfg == nil {
		return nil, errors.Errorf("failed to get config from context: %T", ctx.Value(types.ConfigContextKey))
	}

	req, err := http.NewRequest("POST", cfg.Common.RenderInternalURL+"/render/addframe", bytes.NewBuffer(textJob))
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

func renderVideo(ctx context.Context, url []byte) (*dto.HashResponse, error) {
	// need config for the media-manager render URLs
	cfg := utils.GetFromAny(ctx.Value(types.ConfigContextKey), (*config.Config)(nil))
	if cfg == nil {
		return nil, errors.Errorf("failed to get config from context: %T", ctx.Value(types.ConfigContextKey))
	}

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
