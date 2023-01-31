package node

import (
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/logic/common/posbus"
)

func (n *Node) posBusAutoOnUserAttributeChanged(
	changeType universe.AttributeChangeType, userAttributeID entry.UserAttributeID, value *entry.AttributeValue,
	effectiveOptions *entry.AttributeOptions,
) error {
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

func (n *Node) posBusAutoOnObjectUserAttributeChanged(
	changeType universe.AttributeChangeType, objectUserAttributeID entry.ObjectUserAttributeID,
	value *entry.AttributeValue, effectiveOptions *entry.AttributeOptions,
) error {
	autoOption, err := posbus.GetOptionAutoOption(effectiveOptions)
	if err != nil {
		return errors.WithMessagef(err, "failed to get auto option: %+v", objectUserAttributeID)
	}
	autoMessage, err := posbus.GetOptionAutoMessage(autoOption, changeType, objectUserAttributeID.AttributeID, value)
	if err != nil {
		return errors.WithMessagef(err, "failed to get auto message: %+v", objectUserAttributeID)
	}
	if autoMessage == nil {
		return nil
	}

	object, ok := n.GetObjectFromAllObjects(objectUserAttributeID.ObjectID)
	if !ok {
		return errors.Errorf("object not found: %s", objectUserAttributeID.ObjectID)
	}

	var errs *multierror.Error
	for i := range autoOption.Scope {
		switch autoOption.Scope[i] {
		case entry.WorldPosBusAutoScopeAttributeOption:
			world := object.GetWorld()
			if world == nil {
				errs = multierror.Append(
					errs, errors.Errorf("failed to get object world: %s", autoOption.Scope[i]),
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
		case entry.ObjectPosBusAutoScopeAttributeOption:
			if err := object.Send(autoMessage, false); err != nil {
				errs = multierror.Append(
					errs, errors.WithMessagef(
						err, "failed to send message: %s", autoOption.Scope[i],
					),
				)
			}
		case entry.UserPosBusAutoScopeAttributeOption:
			user, ok := object.GetUser(objectUserAttributeID.UserID, false)
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
