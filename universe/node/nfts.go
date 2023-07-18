package node

import (
	"context"
	"strconv"
	"time"

	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"go.uber.org/zap/zapcore"

	"github.com/momentum-xyz/ubercontroller/harvester"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/logic/tree"
	"github.com/momentum-xyz/ubercontroller/utils"
	"github.com/momentum-xyz/ubercontroller/utils/umid"
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

	return n.AddStakeActivities(stakeEvents)
}

func (n *Node) AddStakeActivities(stakeEvents []*harvester.StakeEvent) error {
	node := universe.GetNode()
	activities := node.GetActivities().GetActivities()

	newStakeEvents := make([]*harvester.StakeEvent, 0)

	for _, s := range stakeEvents {
		exists := false
		for _, a := range activities {
			if a.GetData() != nil && a.GetData().BCTxHash != nil && a.GetData().BCLogIndex != nil {
				if *a.GetData().BCTxHash == s.TxHash && *a.GetData().BCLogIndex == s.LogIndex {
					exists = true
				}
			}
		}

		if !exists {
			newStakeEvents = append(newStakeEvents, s)
		}
	}

	for _, s := range newStakeEvents {
		err := n.AddStakeActivity(s)
		if err != nil {
			return errors.WithMessage(err, "failed to AddStakeActivity")
		}
	}

	return nil
}

func (n *Node) AddStakeActivity(stakeEvent *harvester.StakeEvent) error {
	node := universe.GetNode()

	_, ok := n.GetObject(stakeEvent.OdysseyID, true)
	if !ok {
		// An event for an unknown object, in other words: the 'world' hasn't been created yet.
		// Can easily happen during local development, when when DB and blockchain are out-of-sync.
		// In prod could happen if people bypass our UI/frontend, and directly stake on BC.
		n.log.Errorf("Stake activity %s for unknown object %s", stakeEvent.TxHash, stakeEvent.OdysseyID)
		// Can't really easily recover from this, getting stake activity before 'create' activity would be weird.
		// Creating empty/unclaimed worlds is what we try to avoid by lazily creating worlds.
		// For now, ingoring
		return nil
	}
	aType := entry.ActivityTypeStake
	if stakeEvent.ActivityType == "unstake" { // TODO: change type of StakeEvent.ActivityType
		//aType = entry.ActivityTypeUnstake
		// Ignore unstake events for now
		// wait for design on FE
		n.log.Errorf("Unstake activity %s in %s: TODO", stakeEvent.TxHash, stakeEvent.OdysseyID)
		return nil
	}

	a, err := node.GetActivities().CreateActivity(umid.New())
	if err != nil {
		return errors.WithMessage(err, "failed to CreateActivity")
	}

	if err := a.SetObjectID(stakeEvent.OdysseyID, true); err != nil {
		return errors.WithMessage(err, "failed to set object ID")
	}

	user, err := n.db.GetUsersDB().GetUserByWallet(context.Background(), stakeEvent.Wallet)
	if err != nil {
		return errors.WithMessage(err, "failed to get user by wallet")
	}

	if err := a.SetUserID(user.UserID, true); err != nil {
		return errors.WithMessage(err, "failed to set user ID")
	}

	if err := a.SetCreatedAt(time.Now(), true); err != nil {
		return errors.WithMessage(err, "failed to set created_at")
	}

	if err := a.SetType(&aType, true); err != nil {
		return errors.WithMessage(err, "failed to set activity type")
	}

	modifyFn := func(current *entry.ActivityData) (*entry.ActivityData, error) {
		if current == nil {
			current = &entry.ActivityData{}
		}

		current.BCTxHash = &stakeEvent.TxHash
		current.BCLogIndex = &stakeEvent.LogIndex
		symbol := "MOM"
		if stakeEvent.Kind == 1 {
			symbol = "DAD"
		}
		current.TokenSymbol = &symbol
		amount := stakeEvent.Amount.String()
		current.TokenAmount = &amount

		//current.Position = &position
		//current.Hash = &inBody.Hash
		//current.Description = &inBody.Description

		return current, nil
	}

	_, err = a.SetData(modifyFn, true)
	if err != nil {
		return errors.New("failed to set activity data")
	}

	if err := a.GetActivities().Inject(a); err != nil {
		return errors.New("failed to inject activity")
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
