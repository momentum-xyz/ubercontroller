package node

import (
	"context"
	"strconv"

	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/momentum-xyz/ubercontroller/harvester"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/logic/tree"
	"github.com/momentum-xyz/ubercontroller/utils"
	"github.com/momentum-xyz/ubercontroller/utils/umid"
	"github.com/pkg/errors"
	"go.uber.org/zap/zapcore"
)

func (n *Node) Listener(bcName string, events []*harvester.UpdateEvent, stakeEvents []*harvester.StakeEvent, nftEvent []*harvester.NftEvent) error {
	if n.log.Level() == zapcore.DebugLevel {
		n.log.Debugln("Table Listener:")
		for k, v := range events {
			n.log.Debugf("%+v %+v %+v %+v \n", k, v.Wallet, v.Contract, v.Amount.String())
		}
	}
	if nftEvent != nil && len(nftEvent) > 0 {
		for _, event := range nftEvent {
			if event.To != (ethCommon.Address{}).Hex() {
				seqID := utils.UMIDToSEQ(event.OdysseyID)

				user, err := n.db.GetUsersDB().GetUserByWallet(n.ctx, event.To)
				if user == nil || err != nil {
					n.log.Infof("NFT %d orphan, user with %s not found yet.", seqID, event.To)
					continue
				}

				world, _ := n.GetObjectFromAllObjects(event.OdysseyID)
				if world != nil {
					n.log.Infof("NFT %d world already exists", seqID)
					continue
				}

				templateValue, _ := n.GetNodeAttributes().GetValue(
					entry.NewAttributeID(universe.GetSystemPluginID(), universe.ReservedAttributes.Node.WorldTemplate.Name),
				)

				var worldTemplate tree.WorldTemplate
				err = utils.MapDecode(*templateValue, &worldTemplate)
				if err != nil {
					return errors.WithMessage(err, "failed to decode template")
				}

				objectName := "Odyssey#" + strconv.FormatUint(seqID, 10)

				worldTemplate.ObjectID = &event.OdysseyID
				worldTemplate.ObjectName = &objectName
				worldTemplate.OwnerID = &user.UserID

				n.log.Debugf("Adding odyssey for: %s...", event.OdysseyID)
				_, err = tree.AddWorldFromTemplate(&worldTemplate, true)
				if err != nil {
					return errors.WithMessage(err, "failed to add world from template")
				}
			}
		}
	}

	return nil
}

// Check if this (new) user has an NFT, create world if it doesn't exist yet.
func (n *Node) checkNFTWorld(ctx context.Context, userID umid.UMID, wallet string) error {
	// TODO:  wallet(s) from entry.User or vice-versa
	n.log.Debugf("check nft worlds for wallet %s", wallet)
	nfts, err := n.db.GetNFTsDB().ListNewByWallet(ctx, wallet)
	if err != nil {
		return err
	}
	for _, nft := range nfts {
		n.log.Debugf("check nft %s", nft.ObjectID.String())
		// TODO: Refactor, extract nft world creation function
		seqID := utils.UMIDToSEQ(nft.ObjectID)
		templateValue, _ := n.GetNodeAttributes().GetValue(
			entry.NewAttributeID(universe.GetSystemPluginID(), universe.ReservedAttributes.Node.WorldTemplate.Name),
		)

		var worldTemplate tree.WorldTemplate
		err = utils.MapDecode(*templateValue, &worldTemplate)
		if err != nil {
			return errors.WithMessage(err, "failed to decode template")
		}

		objectName := "Odyssey#" + strconv.FormatUint(seqID, 10)

		worldTemplate.ObjectID = &nft.ObjectID
		worldTemplate.ObjectName = &objectName
		worldTemplate.OwnerID = &userID

		n.log.Debugf("Adding odyssey for: %s...", nft.ObjectID)
		_, err = tree.AddWorldFromTemplate(&worldTemplate, true)
		if err != nil {
			return errors.WithMessage(err, "failed to add world from template")
		}
	}
	return nil
}
