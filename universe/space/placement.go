package space

import (
	"github.com/google/uuid"
	"github.com/momentum-xyz/posbus-protocol/posbus"
	"github.com/momentum-xyz/ubercontroller/pkg/cmath"
	"github.com/momentum-xyz/ubercontroller/pkg/position_algo"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/pkg/errors"
	"sort"
)

// TODO: Rewrite

func (s *Space) GetPlacement(placementMap *entry.SpaceChildPlacement) (position_algo.Algo, error) {

	//fmt.Printf("PLSMAP %+v\n", placementMap)

	var par position_algo.Algo
	algo := "circular"
	if placementMap.Algo != nil {
		algo = *placementMap.Algo
	}

	//fmt.Printf("%s | %+v\n", algo, placementMap.Options)
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

func (s *Space) GetPlacements() map[uuid.UUID]position_algo.Algo {
	//fmt.Printf("eopts %+v\n:", s.GetEffectiveOptions().ChildPlacements)
	placements := s.GetEffectiveOptions().ChildPlacements
	//fmt.Println(len(placements))
	pls := make(map[uuid.UUID]position_algo.Algo)
	for sId, placement := range placements {
		if p, err := s.GetPlacement(placement); err != nil {
			s.log.Error(errors.WithMessage(err, "Space: UpdateMetaFromMap: failed to fill placement"))
		} else {
			//fmt.Printf("%+v | %+v\n", sId, p)
			pls[sId] = p
		}
	}
	return pls
}

func (s *Space) SetActualPosition(pos cmath.SpacePosition, theta float64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if (s.theta != theta) || (*s.actualPosition.Load() != pos) {
		s.theta = theta
		s.actualPosition.Store(&pos)

		go s.UpdateSpawnMessage(false)
		s.GetWorld().Send(
			posbus.NewSetStaticObjectPositionMsg(s.id, *(s.actualPosition.Load())).WebsocketMessage(), false,
		)
	}

	return nil
}

func (s *Space) GetPosition() *cmath.SpacePosition {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.position
}

func (s *Space) GetActualPosition() *cmath.SpacePosition {
	return s.actualPosition.Load()
}

func (s *Space) SetPosition(position *cmath.SpacePosition, updateDB bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if updateDB {
		if err := s.db.SpacesUpdateSpacePosition(s.ctx, s.id, position); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	// TODO: unclear how we have to do it, in case if one or another is nil
	s.position = position
	if s.position != nil {
		s.actualPosition.Store(s.position)
		//s.SetActualPosition(s.position)

		go s.UpdateSpawnMessage(false)
		s.GetWorld().Send(
			posbus.NewSetStaticObjectPositionMsg(s.id, *(s.actualPosition.Load())).WebsocketMessage(), false,
		)
	}

	return nil
}

func (s *Space) UpdateChildrenPosition(recursive bool) error {
	//fmt.Println("pls1", s.GetID())
	pls := s.GetPlacements()
	//fmt.Printf("pls1a:%+v : %+v\n", s.GetID(), pls)
	ChildMap := make(map[uuid.UUID][]uuid.UUID)
	for u := range pls {
		ChildMap[u] = make([]uuid.UUID, 0)
	}
	//fmt.Println("pls2", s.GetID())
	s.Children.Mu.RLock()
	defer s.Children.Mu.RUnlock()

	for _, child := range s.Children.Data {
		if child.GetPosition() == nil {
			spaceTypeId := child.GetSpaceType().GetID()
			if _, ok := pls[spaceTypeId]; !ok {
				spaceTypeId = uuid.Nil
			}
			ChildMap[spaceTypeId] = append(ChildMap[spaceTypeId], child.GetID())
		}
	}
	//fmt.Println("pls3", s.GetID(), ChildMap)
	for u := range pls {
		//fmt.Println("pls4", s.GetID(), u)
		lpm := ChildMap[u]
		//fmt.Println("pls4a", s.GetID(), lpm)
		sort.Slice(lpm, func(i, j int) bool { return lpm[i].ClockSequence() < lpm[j].ClockSequence() })
		//fmt.Println("pls4b", s.GetID(), lpm)
		for i, k := range lpm {
			pos, theta := pls[u].CalcPos(s.theta, *s.GetActualPosition(), i, len(lpm))
			//fmt.Printf(" Position: %s |  %+v\n", s.GetID(), pos)

			child, ok := s.Children.Data[k]
			//fmt.Println(ok)
			if !ok {
				s.log.Errorf("Space: UpdatePosition: failed to get space: %s", k)
				continue
			}
			if err := child.SetActualPosition(pos, theta); err != nil {
				s.log.Errorf("Space: UpdatePosition: failed to update position: %s", k)
			}
			if recursive {
				child.UpdateChildrenPosition(recursive)
			}
		}
	}
	//fmt.Println("pls10", s.GetID())
	return nil
}
