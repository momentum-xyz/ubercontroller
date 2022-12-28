package node

import (
	"github.com/hashicorp/go-multierror"
	"github.com/jackc/pgx/v4"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/common/posbus"
	"github.com/momentum-xyz/ubercontroller/utils/merge"
	"github.com/momentum-xyz/ubercontroller/utils/modify"
	"github.com/pkg/errors"
)

func (n *Node) GetSpaceUserAttributePayload(spaceUserAttributeID entry.SpaceUserAttributeID) (*entry.AttributePayload, bool) {
	spaceUserAttribute, err := n.db.SpaceUserAttributesGetSpaceUserAttributeByID(n.ctx, spaceUserAttributeID)
	if err != nil {
		return nil, false
	}
	return spaceUserAttribute.AttributePayload, true
}

func (n *Node) GetSpaceUserAttributeValue(spaceUserAttributeID entry.SpaceUserAttributeID) (*entry.AttributeValue, bool) {
	value, err := n.db.SpaceUserAttributesGetSpaceUserAttributeValueByID(n.ctx, spaceUserAttributeID)
	if err != nil {
		return nil, false
	}
	return value, true
}

func (n *Node) GetSpaceUserAttributeOptions(spaceUserAttributeID entry.SpaceUserAttributeID) (*entry.AttributeOptions, bool) {
	options, err := n.db.SpaceUserAttributesGetSpaceUserAttributeOptionsByID(n.ctx, spaceUserAttributeID)
	if err != nil {
		return nil, false
	}
	return options, true
}

func (n *Node) GetSpaceUserAttributeEffectiveOptions(
	spaceUserAttributeID entry.SpaceUserAttributeID,
) (*entry.AttributeOptions, bool) {
	attributeType, ok := n.GetAttributeTypes().GetAttributeType(entry.AttributeTypeID(spaceUserAttributeID.AttributeID))
	if !ok {
		return nil, false
	}
	attributeTypeOptions := attributeType.GetOptions()

	attributeOptions, ok := n.GetSpaceUserAttributeOptions(spaceUserAttributeID)
	if !ok {
		return nil, false
	}

	effectiveOptions, err := merge.Auto(attributeOptions, attributeTypeOptions)
	if err != nil {
		n.log.Error(
			errors.WithMessagef(
				err, "Node: GetSpaceUserAttributeEffectiveOptions: failed to merge options: %+v", spaceUserAttributeID,
			),
		)
		return nil, false
	}

	return effectiveOptions, true
}

func (n *Node) UpsertSpaceUserAttribute(
	spaceUserAttributeID entry.SpaceUserAttributeID, modifyFn modify.Fn[entry.AttributePayload],
) (*entry.AttributePayload, error) {
	payload, err := n.db.SpaceUserAttributesUpsertSpaceUserAttribute(n.ctx, spaceUserAttributeID, modifyFn)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to upsert space user attribute")
	}

	if n.GetEnabled() {
		go func() {
			var value *entry.AttributeValue
			if payload != nil {
				value = payload.Value
			}
			n.onSpaceUserAttributeChanged(universe.ChangedAttributeChangeType, spaceUserAttributeID, value, nil)
		}()
	}

	return payload, nil
}

func (n *Node) UpdateSpaceUserAttributeValue(
	spaceUserAttributeID entry.SpaceUserAttributeID, modifyFn modify.Fn[entry.AttributeValue],
) (*entry.AttributeValue, error) {
	value, err := n.db.SpaceUserAttributesUpdateSpaceUserAttributeValue(n.ctx, spaceUserAttributeID, modifyFn)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to update space user attribute value")
	}

	if n.GetEnabled() {
		go n.onSpaceUserAttributeChanged(universe.ChangedAttributeChangeType, spaceUserAttributeID, value, nil)
	}

	return value, nil
}

func (n *Node) UpdateSpaceUserAttributeOptions(
	spaceUserAttributeID entry.SpaceUserAttributeID, modifyFn modify.Fn[entry.AttributeOptions],
) (*entry.AttributeOptions, error) {
	options, err := n.db.SpaceUserAttributesUpdateSpaceUserAttributeOptions(n.ctx, spaceUserAttributeID, modifyFn)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to update space user attribute options")
	}

	if n.GetEnabled() {
		go func() {
			value, ok := n.GetSpaceUserAttributeValue(spaceUserAttributeID)
			if !ok {
				n.log.Error(
					errors.Errorf(
						"Node: UpdateSpaceUserAttributeOptions: failed to get space user attribute value: %+v",
						spaceUserAttributeID,
					),
				)
				return
			}
			n.onSpaceUserAttributeChanged(universe.ChangedAttributeChangeType, spaceUserAttributeID, value, nil)
		}()
	}

	return options, nil
}

func (n *Node) RemoveSpaceUserAttribute(spaceUserAttributeID entry.SpaceUserAttributeID) (bool, error) {
	attributeEffectiveOptions, attributeEffectiveOptionsOK := n.GetSpaceUserAttributeEffectiveOptions(spaceUserAttributeID)

	if err := n.db.SpaceUserAttributesRemoveSpaceUserAttributeByID(n.ctx, spaceUserAttributeID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, errors.WithMessage(err, "failed to remove space user attribute")
	}

	if n.GetEnabled() {
		go func() {
			if !attributeEffectiveOptionsOK {
				n.log.Error(
					errors.Errorf(
						"Node: RemoveSpaceUserAttribute: failed to get space user attribute effective options",
					),
				)
				return
			}
			n.onSpaceUserAttributeChanged(
				universe.RemovedAttributeChangeType, spaceUserAttributeID, nil, attributeEffectiveOptions,
			)
		}()
	}

	return true, nil
}

func (n *Node) onSpaceUserAttributeChanged(
	changeType universe.AttributeChangeType, spaceUserAttributeID entry.SpaceUserAttributeID, value *entry.AttributeValue,
	effectiveOptions *entry.AttributeOptions,
) {
	if effectiveOptions == nil {
		options, ok := n.GetSpaceUserAttributeEffectiveOptions(spaceUserAttributeID)
		if !ok {
			n.log.Error(
				errors.Errorf(
					"Node: onSpaceUserAttributeChanged: failed to get space user attribute effective options: %+v",
					spaceUserAttributeID,
				),
			)
			return
		}
		effectiveOptions = options
	}

	go func() {
		if err := n.posBusAutoOnSpaceUserAttributeChanged(changeType, spaceUserAttributeID, value, effectiveOptions); err != nil {
			n.log.Error(
				errors.WithMessagef(
					err, "Node: onSpaceUserAttributeChanged: failed to handle pos bus auto: %+v", spaceUserAttributeID,
				),
			)
		}
	}()
}

func (n *Node) posBusAutoOnSpaceUserAttributeChanged(
	changeType universe.AttributeChangeType, spaceUserAttributeID entry.SpaceUserAttributeID, value *entry.AttributeValue,
	effectiveOptions *entry.AttributeOptions,
) error {
	autoOption, err := posbus.GetOptionAutoOption(effectiveOptions)
	if err != nil {
		return errors.WithMessagef(err, "failed to get auto option: %+v", spaceUserAttributeID)
	}
	autoMessage, err := posbus.GetOptionAutoMessage(autoOption, changeType, spaceUserAttributeID.AttributeID, value)
	if err != nil {
		return errors.WithMessagef(err, "failed to get auto message: %+v", spaceUserAttributeID)
	}
	if autoMessage == nil {
		return nil
	}

	space, ok := n.GetSpaceFromAllSpaces(spaceUserAttributeID.SpaceID)
	if !ok {
		return errors.Errorf("space not found: %s", spaceUserAttributeID.SpaceID)
	}

	var errs *multierror.Error
	for i := range autoOption.Scope {
		switch autoOption.Scope[i] {
		case entry.WorldPosBusAutoScopeAttributeOption:
			world := space.GetWorld()
			if world == nil {
				errs = multierror.Append(
					errs, errors.Errorf("failed to get space world: %s", autoOption.Scope[i]),
				)
				continue
			}
			if err := world.Send(autoMessage, true); err != nil {
				errs = multierror.Append(
					errs, errors.WithMessagef(
						err, "failed to send message: %s", autoOption.Scope[i],
					),
				)
			}
		case entry.SpacePosBusAutoScopeAttributeOption:
			if err := space.Send(autoMessage, false); err != nil {
				errs = multierror.Append(
					errs, errors.WithMessagef(
						err, "failed to send message: %s", autoOption.Scope[i],
					),
				)
			}
		case entry.UserPosBusAutoScopeAttributeOption:
			user, ok := space.GetUser(spaceUserAttributeID.UserID, false)
			if !ok {
				continue
			}
			if err := user.Send(autoMessage); err != nil {
				errs = multierror.Append(
					errs, errors.WithMessagef(
						err, "failed to send message: %s", autoOption.Scope[i],
					),
				)
			}
		default:
			errs = multierror.Append(
				errs, errors.Errorf(
					"scope type in not supported: %s", autoOption.Scope[i],
				),
			)
		}
	}

	return errs.ErrorOrNil()
}
