package object

import (
	"github.com/hashicorp/go-multierror"
	pb "github.com/momentum-xyz/ubercontroller/pkg/posbus"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe/logic/common/posbus"
	"github.com/pkg/errors"
)

func (o *Object) posBusAutoOnObjectAttributeChanged(
	changeType pb.AttributeChangeType, attributeID entry.AttributeID, value *entry.AttributeValue,
	effectiveOptions *entry.AttributeOptions,
) error {
	autoOption, err := posbus.GetOptionAutoOption(effectiveOptions)
	if err != nil {
		return errors.WithMessagef(err, "failed to get auto option: %+v", attributeID)
	}
	targetID := o.GetID()
	autoMessage, err := posbus.GetOptionAutoMessage(autoOption, changeType, attributeID, targetID, value)
	if err != nil {
		return errors.WithMessagef(err, "failed to get auto message: %+v", attributeID)
	}
	if autoMessage == nil {
		return nil
	}

	var errs *multierror.Error
	for i := range autoOption.Scope {
		switch autoOption.Scope[i] {
		case entry.WorldPosBusAutoScopeAttributeOption:
			world := o.GetWorld()
			if world == nil {
				errs = multierror.Append(
					err, errors.Errorf("failed to get world: %s", autoOption.Scope[i]),
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
			if err := o.Send(autoMessage, false); err != nil {
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
