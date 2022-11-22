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

func (n *Node) GetUserAttributePayload(userAttributeID entry.UserAttributeID) (*entry.AttributePayload, bool) {
	userAttribute, err := n.db.UserAttributesGetUserAttributeByID(n.ctx, userAttributeID)
	if err != nil {
		return nil, false
	}
	return userAttribute.AttributePayload, true
}

func (n *Node) GetUserAttributeValue(userAttributeID entry.UserAttributeID) (*entry.AttributeValue, bool) {
	value, err := n.db.UserAttributesGetUserAttributeValueByID(n.ctx, userAttributeID)
	if err != nil {
		return nil, false
	}
	return value, true
}

func (n *Node) GetUserAttributeOptions(userAttributeID entry.UserAttributeID) (*entry.AttributeOptions, bool) {
	options, err := n.db.UserAttributesGetUserAttributeOptionsByID(n.ctx, userAttributeID)
	if err != nil {
		return nil, false
	}
	return options, true
}

func (n *Node) GetUserAttributeEffectiveOptions(userAttributeID entry.UserAttributeID) (*entry.AttributeOptions, bool) {
	attributeType, ok := n.GetAttributeTypes().GetAttributeType(entry.AttributeTypeID(userAttributeID.AttributeID))
	if !ok {
		return nil, false
	}
	attributeTypeOptions := attributeType.GetOptions()

	attributeOptions, ok := n.GetUserAttributeOptions(userAttributeID)
	if !ok {
		return nil, false
	}

	effectiveOptions, err := merge.Auto(attributeOptions, attributeTypeOptions)
	if err != nil {
		n.log.Error(
			errors.WithMessagef(
				err, "Node: GetUserAttributeEffectiveOptions: failed to merge options: %+v", userAttributeID,
			),
		)
		return nil, false
	}

	return effectiveOptions, true
}

func (n *Node) UpsertUserAttribute(
	userAttributeID entry.UserAttributeID, modifyFn modify.Fn[entry.AttributePayload],
) (*entry.UserAttribute, error) {
	userAttribute, err := n.db.UserAttributesUpsertUserAttribute(n.ctx, userAttributeID, modifyFn)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to upsert user attribute")
	}

	go func() {
		var value any
		if userAttribute.AttributePayload != nil {
			value = userAttribute.AttributePayload.Value
		}
		if err := n.onUserAttributeChanged(universe.ChangedAttributeChangeType, userAttributeID, value, nil); err != nil {
			n.log.Error(
				errors.WithMessagef(
					err,
					"Node: UpsertUserAttribute: failed to call onUserAttributeChanged: %+v",
					userAttributeID,
				),
			)
		}
	}()

	return userAttribute, nil
}

func (n *Node) UpdateUserAttributeValue(
	userAttributeID entry.UserAttributeID, modifyFn modify.Fn[entry.AttributeValue],
) (*entry.AttributeValue, error) {
	value, err := n.db.UserAttributesUpdateUserAttributeValue(n.ctx, userAttributeID, modifyFn)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to update user attribute value")
	}

	go func() {
		if err := n.onUserAttributeChanged(
			universe.ChangedAttributeChangeType, userAttributeID, value, nil,
		); err != nil {
			n.log.Error(
				errors.WithMessagef(
					err,
					"Node: UpdateUserAttributeValue: failed to call onUserAttributeChanged: %+v",
					userAttributeID,
				),
			)
		}
	}()

	return value, nil
}

func (n *Node) UpdateUserAttributeOptions(
	userAttributeID entry.UserAttributeID, modifyFn modify.Fn[entry.AttributeOptions],
) (*entry.AttributeOptions, error) {
	options, err := n.db.UserAttributesUpdateUserAttributeOptions(n.ctx, userAttributeID, modifyFn)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to update user attribute options")
	}

	go func() {
		value, ok := n.GetUserAttributeValue(userAttributeID)
		if !ok {
			n.log.Error(
				errors.Errorf(
					"Node: UpdateUserAttributeOptions: failed to get user attribute value: %+v", userAttributeID,
				),
			)
			return
		}
		if err := n.onUserAttributeChanged(
			universe.ChangedAttributeChangeType, userAttributeID, value, nil,
		); err != nil {
			n.log.Error(
				errors.WithMessagef(
					err,
					"Node: UpdateUserAttributeValue: failed to call onUserAttributeChanged: %+v",
					userAttributeID,
				),
			)
		}
	}()

	return options, nil
}

func (n *Node) RemoveUserAttribute(userAttributeID entry.UserAttributeID) (bool, error) {
	attributeEffectiveOptions, attributeEffectiveOptionsOK := n.GetUserAttributeEffectiveOptions(userAttributeID)

	if err := n.db.UserAttributesRemoveUserAttributeByID(n.ctx, userAttributeID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, errors.WithMessage(err, "failed to remove user attribute")
	}

	go func() {
		if !attributeEffectiveOptionsOK {
			n.log.Error(
				errors.Errorf(
					"Node: RemoveUserAttribute: failed to get user attribute effective options",
				),
			)
			return
		}
		if err := n.onUserAttributeChanged(
			universe.RemovedAttributeChangeType, userAttributeID, nil, attributeEffectiveOptions,
		); err != nil {
			n.log.Error(
				errors.WithMessagef(
					err,
					"Node: RemoveUserAttribute: failed to call onUserAttributeChanged: %+v",
					userAttributeID,
				),
			)
		}
	}()

	return true, nil
}

func (n *Node) onUserAttributeChanged(
	changeType universe.AttributeChangeType, userAttributeID entry.UserAttributeID, value any,
	effectiveOptions *entry.AttributeOptions,
) error {
	if effectiveOptions == nil {
		options, ok := n.GetUserAttributeEffectiveOptions(userAttributeID)
		if !ok {
			return errors.Errorf("failed to get user attribute effective options: %+v", userAttributeID)
		}
		effectiveOptions = options
	}

	autoOption, err := posbus.GetOptionAutoOption(effectiveOptions)
	if err != nil {
		return errors.WithMessagef(err, "failed to get auto option: %+v", userAttributeID)
	}
	autoMessage, err := posbus.GetOptionAutoMessage(autoOption, changeType, userAttributeID.AttributeID, value)
	if err != nil {
		return errors.WithMessagef(err, "failed to get auto message: %+v", userAttributeID)
	}
	if autoMessage == nil {
		return nil
	}

	var users []universe.User
	for _, world := range n.GetWorlds().GetWorlds() {
		if user, ok := world.GetUser(userAttributeID.UserID, true); ok {
			users = append(users, user)
		}
	}

	var errs *multierror.Error
	for i := range autoOption.Scope {
		switch autoOption.Scope[i] {
		case entry.WorldPosBusAutoScopeAttributeOption:
			for i := range users {
				world := users[i].GetWorld()
				if world == nil {
					errs = multierror.Append(
						errs, errors.Errorf("failed to get world: %s", autoOption.Scope[i]),
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
			}
		case entry.UserPosBusAutoScopeAttributeOption:
			for i := range users {
				if err := users[i].Send(autoMessage); err != nil {
					errs = multierror.Append(
						errs, errors.WithMessagef(
							err, "failed to send message: %s", autoOption.Scope[i],
						),
					)
				}
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
