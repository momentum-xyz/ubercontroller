package slot

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/pkg/media"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe/logic/api/dto"
	"github.com/momentum-xyz/ubercontroller/utils"
)

const (
	RenderKindNone uint = iota
	RenderKindVideo
	RenderKindText
)

// TODO: requre optimization to not do it every time we touch attribute
func GetOptionAutoOption(
	attributeID entry.AttributeID, options *entry.AttributeOptions,
) (*entry.RenderAutoAttributeOption, error) {
	if options == nil {
		return nil, nil
	}

	autoOptionsValue, ok := (*options)["render_auto"]
	if !ok {
		return nil, nil
	}

	var autoOption *entry.RenderAutoAttributeOption
	if err := utils.MapDecode(autoOptionsValue, &autoOption); err != nil {
		return nil, errors.WithMessage(err, "failed to decode auto option")
	}

	if autoOption.SlotType == entry.SlotTypeInvalid || autoOption.ContentType == entry.SlotContentTypeInvalid {
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
	media *media.Media, option *entry.RenderAutoAttributeOption, value *entry.AttributeValue,
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

	var resp *dto.HashResponse

	switch renderKind {
	case RenderKindVideo:
		payload, err := json.Marshal(
			map[string]any{
				"url": valueString,
			},
		)
		if err != nil {
			return nil, errors.WithMessage(err, "failed to marshal preRenderHash")
		}

		hash, err := media.AddTube(payload)
		if err != nil {
			return nil, errors.WithMessage(err, "failed to add tube")
		}

		resp = &dto.HashResponse{
			Hash: hash,
		}
	case RenderKindText:
		if option.TextRenderTemplate == "" {
			return nil, nil
		}

		payload := []byte(strings.Replace(option.TextRenderTemplate, "%TEXT%", valueString, -1))

		hash, err := media.AddFrame(payload)
		if err != nil {
			return nil, errors.WithMessage(err, "failed to add frame")
		}

		resp = &dto.HashResponse{
			Hash: hash,
		}
	}
	return resp, nil
}
