package object

import (
	"github.com/momentum-xyz/ubercontroller/pkg/posbus"
	"github.com/momentum-xyz/ubercontroller/utils/umid"
	"sort"

	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/pkg/cmath"
	"github.com/momentum-xyz/ubercontroller/pkg/position_algo"
	"github.com/momentum-xyz/ubercontroller/types/entry"
)

// TODO: Rewrite

func (o *Object) GetPlacement(placementMap *entry.ObjectChildPlacement) (position_algo.Algo, error) {

	//fmt.Printf("PLSMAP %+v\n", placementMap)

	var par position_algo.Algo
	algo := "circular"
	if placementMap.Algo != nil {
		algo = *placementMap.Algo
	}

	//fmt.Printf("%o | %+v\n", algo, placementMap.Options)
	switch algo {
	case "circular":
		par = position_algo.NewCircular(placementMap.Options)
	case "helix":
		par = position_algo.NewHelix(placementMap.Options)
	case "sector":
		par = position_algo.NewSector(placementMap.Options)
	case "spiral":
		par = position_algo.NewSpiral(placementMap.Options)
	case "hexaspiral":
		par = position_algo.NewHexaSpiral(placementMap.Options)
	}
	//fmt.Printf("%+v\n", par)
	return par, nil
}

func (o *Object) GetPlacements() map[umid.UMID]position_algo.Algo {
	//fmt.Printf("eopts %+v\n:", o.GetEffectiveOptions().ChildPlacements)
	placements := o.GetEffectiveOptions().ChildPlacements
	//fmt.Println(len(placements))
	pls := make(map[umid.UMID]position_algo.Algo)
	for sId, placement := range placements {
		if p, err := o.GetPlacement(placement); err != nil {
			o.log.Error(errors.WithMessage(err, "Object: UpdateMetaFromMap: failed to fill placement"))
		} else {
			//fmt.Printf("%+v | %+v\n", sId, p)
			pls[sId] = p
		}
	}
	return pls
}

func (o *Object) SetActualTransform(pos cmath.ObjectTransform, theta float64) error {
	o.Mu.Lock()
	defer o.Mu.Unlock()

	if (o.theta != theta) || (*o.actualPosition.Load() != pos) {
		o.theta = theta
		o.actualPosition.Store(&pos)

		if o.enabled.Load() {
			go func() {
				o.UpdateSpawnMessage()
				world := o.GetWorld()
				if world != nil {
					world.Send(
						posbus.NewMessageFromData(posbus.TypeSetObjectPosition, *(o.GetActualTransform())).WSMessage(),
						true,
					)
				}
			}()
		}
	}

	return nil
}

func (o *Object) GetTransform() *cmath.ObjectTransform {
	o.Mu.RLock()
	defer o.Mu.RUnlock()

	return o.transform
}

func (o *Object) GetActualTransform() *cmath.ObjectTransform {
	return o.actualPosition.Load()
}

func (o *Object) SetTransform(position *cmath.ObjectTransform, updateDB bool) error {
	o.Mu.Lock()
	defer o.Mu.Unlock()

	if updateDB {
		if err := o.db.GetObjectsDB().UpdateObjectPosition(o.ctx, o.GetID(), position); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	// TODO: unclear how we have to do it, in case if one or another is nil
	o.transform = position
	if o.transform != nil {
		o.actualPosition.Store(o.transform)

		if o.enabled.Load() {
			go func() {
				o.UpdateSpawnMessage()
				world := o.GetWorld()
				if world != nil {
					world.Send(
						posbus.NewMessageFromData(posbus.TypeSetObjectPosition, *(o.GetActualTransform())).WSMessage(),
						true,
					)
				}
			}()
		}
	}

	return nil
}

func (o *Object) UpdateChildrenPosition(recursive bool) error {
	//fmt.Println("pls1", o.GetID())
	pls := o.GetPlacements()
	//fmt.Printf("pls1a:%+v : %+v\n", o.GetID(), pls)
	ChildMap := make(map[umid.UMID][]umid.UMID)
	for u := range pls {
		ChildMap[u] = make([]umid.UMID, 0)
	}
	//fmt.Println("pls2", o.GetID())
	o.Children.Mu.RLock()
	defer o.Children.Mu.RUnlock()

	for _, child := range o.Children.Data {
		if child.GetTransform() == nil {
			objectTypeID := child.GetObjectType().GetID()
			if _, ok := pls[objectTypeID]; !ok {
				objectTypeID = umid.Nil
			}
			ChildMap[objectTypeID] = append(ChildMap[objectTypeID], child.GetID())
		}
	}
	//fmt.Println("pls3", o.GetID(), ChildMap)
	for u := range pls {
		//fmt.Println("pls4", o.GetID(), u)
		lpm := ChildMap[u]
		//fmt.Println("pls4a", o.GetID(), lpm)
		sort.Slice(lpm, func(i, j int) bool { return lpm[i].ClockSequence() < lpm[j].ClockSequence() })
		//fmt.Println("pls4b", o.GetID(), lpm)
		for i, k := range lpm {
			pos, theta := pls[u].CalcPos(o.theta, *o.GetActualTransform(), i, len(lpm))
			//fmt.Printf(" Position: %o |  %+v\n", o.GetID(), pos)

			child, ok := o.Children.Data[k]
			//fmt.Println(ok)
			if !ok {
				o.log.Errorf("Object: UpdatePosition: failed to get object: %s", k)
				continue
			}
			if err := child.SetActualTransform(pos, theta); err != nil {
				o.log.Errorf("Object: UpdatePosition: failed to update transform: %s", k)
			}

			if !recursive {
				continue
			}

			child.UpdateChildrenPosition(true)
		}
	}
	//fmt.Println("pls10", o.GetID())
	return nil
}
