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
) (*entry.SpaceUserAttribute, error) {
	spaceUserAttribute, err := n.db.SpaceUserAttributesUpsertSpaceUserAttribute(n.ctx, spaceUserAttributeID, modifyFn)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to upsert space user attribute")
	}

	go func() {
		var value any
		if spaceUserAttribute.AttributePayload != nil {
			value = spaceUserAttribute.AttributePayload.Value
		}
		if err := n.onSpaceUserAttributeChanged(
			universe.ChangedAttributeChangeType, spaceUserAttributeID, value, nil,
		); err != nil {
			n.log.Error(
				errors.WithMessagef(
					err,
					"Node: UpsertSpaceUserAttribute: failed to call onSpaceUserAttributeChanged: %+v",
					spaceUserAttributeID,
				),
			)
		}
	}()

	return spaceUserAttribute, nil
}

func (n *Node) UpdateSpaceUserAttributeValue(
	spaceUserAttributeID entry.SpaceUserAttributeID, modifyFn modify.Fn[entry.AttributeValue],
) (*entry.AttributeValue, error) {
	value, err := n.db.SpaceUserAttributesUpdateSpaceUserAttributeValue(n.ctx, spaceUserAttributeID, modifyFn)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to update space user attribute value")
	}

	go func() {
		changeType := universe.ChangedAttributeChangeType
		if value == nil {
			changeType = universe.RemovedAttributeChangeType
		}
		if err := n.onSpaceUserAttributeChanged(changeType, spaceUserAttributeID, value, nil); err != nil {
			n.log.Error(
				errors.WithMessagef(
					err,
					"Node: UpdateSpaceUserAttributeValue: failed to call onSpaceUserAttributeChanged: %+v",
					spaceUserAttributeID,
				),
			)
		}
	}()

	return value, nil
}

func (n *Node) UpdateSpaceUserAttributeOptions(
	spaceUserAttributeID entry.SpaceUserAttributeID, modifyFn modify.Fn[entry.AttributeOptions],
) (*entry.AttributeOptions, error) {
	options, err := n.db.SpaceUserAttributesUpdateSpaceUserAttributeOptions(n.ctx, spaceUserAttributeID, modifyFn)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to update space user attribute options")
	}

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
		if err := n.onSpaceUserAttributeChanged(
			universe.ChangedAttributeChangeType, spaceUserAttributeID, value, nil,
		); err != nil {
			n.log.Error(
				errors.WithMessagef(
					err,
					"Node: UpdateSpaceUserAttributeOptions: failed to call onSpaceUserAttributeChanged: %+v",
					spaceUserAttributeID,
				),
			)
		}
	}()

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

	go func() {
		if !attributeEffectiveOptionsOK {
			n.log.Error(
				errors.Errorf(
					"Node: RemoveSpaceUserAttribute: failed to get space user attribute effective options",
				),
			)
			return
		}
		if err := n.onSpaceUserAttributeChanged(
			universe.RemovedAttributeChangeType, spaceUserAttributeID, nil, attributeEffectiveOptions,
		); err != nil {
			n.log.Error(
				errors.WithMessagef(
					err,
					"Node: RemoveSpaceUserAttribute: failed to call onSpaceUserAttributeChanged: %+v",
					spaceUserAttributeID,
				),
			)
		}
	}()

	return true, nil
}

func (n *Node) onSpaceUserAttributeChanged(
	changeType universe.AttributeChangeType, spaceUserAttributeID entry.SpaceUserAttributeID, value any,
	effectiveOptions *entry.AttributeOptions,
) error {
	if effectiveOptions == nil {
		options, ok := n.GetSpaceUserAttributeEffectiveOptions(spaceUserAttributeID)
		if !ok {
			return errors.Errorf("failed to get space user attribute effective options: %+v", spaceUserAttributeID)
		}
		effectiveOptions = options
	}

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
