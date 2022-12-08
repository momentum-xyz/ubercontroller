package node

import (
	"github.com/hashicorp/go-multierror"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/common/posbus"
	"github.com/momentum-xyz/ubercontroller/utils/merge"
	"github.com/momentum-xyz/ubercontroller/utils/modify"
)

func (n *Node) GetUserUserAttributePayload(userUserAttributeID entry.UserUserAttributeID) (*entry.AttributePayload, bool) {
	userUserAttribute, err := n.db.UserUserAttributesGetUserUserAttributeByID(n.ctx, userUserAttributeID)
	if err != nil {
		return nil, false
	}
	return userUserAttribute.AttributePayload, true
}

func (n *Node) GetUserUserAttributeValue(userUserAttributeID entry.UserUserAttributeID) (*entry.AttributeValue, bool) {
	value, err := n.db.UserUserAttributesGetUserUserAttributeValueByID(n.ctx, userUserAttributeID)
	if err != nil {
		return nil, false
	}
	return value, true
}

func (n *Node) GetUserUserAttributeOptions(userUserAttributeID entry.UserUserAttributeID) (*entry.AttributeOptions, bool) {
	options, err := n.db.UserUserAttributesGetUserUserAttributeOptionsByID(n.ctx, userUserAttributeID)
	if err != nil {
		return nil, false
	}
	return options, true
}

func (n *Node) GetUserUserAttributeEffectiveOptions(userUserAttributeID entry.UserUserAttributeID) (*entry.AttributeOptions, bool) {
	attributeType, ok := n.GetAttributeTypes().GetAttributeType(entry.AttributeTypeID(userUserAttributeID.AttributeID))
	if !ok {
		return nil, false
	}
	attributeTypeOptions := attributeType.GetOptions()

	attributeOptions, ok := n.GetUserUserAttributeOptions(userUserAttributeID)
	if !ok {
		return nil, false
	}

	effectiveOptions, err := merge.Auto(attributeOptions, attributeTypeOptions)
	if err != nil {
		n.log.Error(
			errors.WithMessagef(
				err, "Node: GetUserUserAttributeEffectiveOptions: failed to merge options: %+v", userUserAttributeID,
			),
		)
		return nil, false
	}

	return effectiveOptions, true
}

func (n *Node) UpsertUserUserAttribute(
	userUserAttributeID entry.UserUserAttributeID, modifyFn modify.Fn[entry.AttributePayload],
) (*entry.UserUserAttribute, error) {
	userUserAttribute, err := n.db.UserUserAttributesUpsertUserUserAttribute(n.ctx, userUserAttributeID, modifyFn)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to upsert user user attribute")
	}

	if n.GetEnabled() {
		go func() {
			var value *entry.AttributeValue
			if userUserAttribute.AttributePayload != nil {
				value = userUserAttribute.AttributePayload.Value
			}
			n.onUserUserAttributeChanged(universe.ChangedAttributeChangeType, userUserAttributeID, value, nil)
		}()
	}

	return userUserAttribute, nil
}

func (n *Node) onUserUserAttributeChanged(
	changeType universe.AttributeChangeType, userUserAttributeID entry.UserUserAttributeID, value *entry.AttributeValue,
	effectiveOptions *entry.AttributeOptions,
) {
	if effectiveOptions == nil {
		options, ok := n.GetUserUserAttributeEffectiveOptions(userUserAttributeID)
		if !ok {
			n.log.Error(
				errors.Errorf(
					"Node: onUserUserAttributeChanged: failed to get user user attribute effective options: %+v",
					userUserAttributeID,
				),
			)
			return
		}
		effectiveOptions = options
	}

	go func() {
		if err := n.posBusAutoOnUserUserAttributeChanged(changeType, userUserAttributeID, value, effectiveOptions); err != nil {
			n.log.Error(
				errors.WithMessagef(
					err, "Node: onUserUserAttributeChanged: failed to handle pos bus auto: %+v", userUserAttributeID,
				),
			)
		}
	}()
}

func (n *Node) posBusAutoOnUserUserAttributeChanged(
	changeType universe.AttributeChangeType, userUserAttributeID entry.UserUserAttributeID, value *entry.AttributeValue,
	effectiveOptions *entry.AttributeOptions,
) error {
	autoOption, err := posbus.GetOptionAutoOption(effectiveOptions)
	if err != nil {
		return errors.WithMessagef(err, "failed to get auto option: %+v", userUserAttributeID)
	}
	autoMessage, err := posbus.GetOptionAutoMessage(autoOption, changeType, userUserAttributeID.AttributeID, value)
	if err != nil {
		return errors.WithMessagef(err, "failed to get auto message: %+v", userUserAttributeID)
	}
	if autoMessage == nil {
		return nil
	}

	var users []universe.User
	for _, world := range n.GetWorlds().GetWorlds() {
		if user, ok := world.GetUser(userUserAttributeID.SourceUserID, true); ok {
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

func (n *Node) RemoveUserUserAttribute(userUserAttributeID entry.UserUserAttributeID) (bool, error) {
	attributeEffectiveOptions, attributeEffectiveOptionsOK := n.GetUserUserAttributeEffectiveOptions(userUserAttributeID)

	if err := n.db.UserUserAttributesRemoveUserUserAttributeByID(n.ctx, userUserAttributeID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, errors.WithMessage(err, "failed to remove user attribute")
	}

	if n.GetEnabled() {
		go func() {
			if !attributeEffectiveOptionsOK {
				n.log.Error(
					errors.Errorf(
						"Node: RemoveUserAttribute: failed to get user attribute effective options",
					),
				)
				return
			}
			n.onUserUserAttributeChanged(universe.RemovedAttributeChangeType, userUserAttributeID, nil, attributeEffectiveOptions)
		}()
	}

	return true, nil
}

func (n *Node) UpdateUserUserAttributeOptions(
	userUserAttributeID entry.UserUserAttributeID, modifyFn modify.Fn[entry.AttributeOptions],
) (*entry.AttributeOptions, error) {
	options, err := n.db.UserUserAttributesUpdateUserUserAttributeOptions(n.ctx, userUserAttributeID, modifyFn)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to update user attribute options")
	}

	if n.GetEnabled() {
		go func() {
			value, ok := n.GetUserUserAttributeValue(userUserAttributeID)
			if !ok {
				n.log.Error(
					errors.Errorf(
						"Node: UpdateUserAttributeOptions: failed to get user attribute value: %+v", userUserAttributeID,
					),
				)
				return
			}
			n.onUserUserAttributeChanged(universe.ChangedAttributeChangeType, userUserAttributeID, value, nil)
		}()
	}

	return options, nil
}

func (n *Node) UpdateUserUserAttributeValue(
	userUserAttributeID entry.UserUserAttributeID, modifyFn modify.Fn[entry.AttributeValue],
) (*entry.AttributeValue, error) {
	value, err := n.db.UserUserAttributesUpdateUserUserAttributeValue(n.ctx, userUserAttributeID, modifyFn)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to update user attribute value")
	}

	if n.GetEnabled() {
		go n.onUserUserAttributeChanged(universe.ChangedAttributeChangeType, userUserAttributeID, value, nil)
	}

	return value, nil
}
