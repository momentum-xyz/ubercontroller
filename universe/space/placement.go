package space

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/momentum-xyz/ubercontroller/pkg/cmath"
	"github.com/momentum-xyz/ubercontroller/pkg/position_algo"
	"github.com/momentum-xyz/ubercontroller/utils"
	"github.com/pkg/errors"
	"sort"
)

// TODO: Rewrite

func (s *Space) FillPlacement(placementMap map[string]interface{}) (position_algo.Algo, error) {

	s.log.Debug("PLSMAP", placementMap)

	var par position_algo.Algo
	algo := "circular"
	if v, ok := placementMap["algo"]; ok {
		algo = v.(string)
	}

	switch algo {
	case "circular":
		par = position_algo.NewCircular(placementMap)
	case "helix":
		par = position_algo.NewHelix(placementMap)
	case "sector":
		par = position_algo.NewSector(placementMap)
	case "spiral":
		par = position_algo.NewSpiral(placementMap)
	case "hexaspiral":
		par = position_algo.NewHexaSpiral(placementMap)
	}

	return par, nil
}

func (s *Space) GetPlacements() map[uuid.UUID]position_algo.Algo {
	childPlacements := s.GetEffectiveOptions().ChildPlacements
	jsonData := []byte(utils.GetFromAny(childPlacements, ""))
	var placements map[uuid.UUID]interface{}
	if err := json.Unmarshal(jsonData, &placements); err != nil {
		s.log.Debug(errors.WithMessage(err, "Space: UpdateMetaFromMap: failed to unmarshal child place"))
	}
	pls := make(map[uuid.UUID]position_algo.Algo)
	for sId, placement := range placements {
		if p, err := s.FillPlacement(placement.(map[string]interface{})); err != nil {
			pls[sId] = p
			s.log.Error(errors.WithMessage(err, "Space: UpdateMetaFromMap: failed to fill placement"))
		}
	}
	return pls
}

func (s *Space) SetActualPosition(pos cmath.Vec3, theta float64, force bool) error {
	//s.mu.Lock()
	//defer s.mu.Unlock()
	if (s.theta != theta) || (*s.actualPosition.Load() != pos) || (force) {
		s.theta = theta
		s.actualPosition.Store(&pos)
		s.UpdateSpawnMessage()
	}
	return nil
}

func (s *Space) GetPosition() *cmath.Vec3 {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.position
}

func (s *Space) GetActualPosition() cmath.Vec3 {
	return *s.actualPosition.Load()
}

func (s *Space) SetPosition(position *cmath.Vec3, updateDB bool) error {

	if updateDB {
		if err := s.db.SpacesUpdateSpacePosition(s.ctx, s.id, position); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	s.position = position
	if s.position != nil {
		s.actualPosition.Store(s.position)
		s.UpdateSpawnMessage()
	}

	return nil
}

func (s *Space) UpdateChildrenPosition(recursive bool, force bool) error {
	fmt.Println("pls1", s.GetID())
	pls := s.GetPlacements()

	ChildMap := make(map[uuid.UUID][]uuid.UUID)
	for u := range pls {
		ChildMap[u] = make([]uuid.UUID, 0)
	}
	fmt.Println("pls2", s.GetID())
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
	fmt.Println("pls3", s.GetID())
	for u := range pls {
		lpm := ChildMap[u]

		sort.Slice(lpm, func(i, j int) bool { return lpm[i].ClockSequence() < lpm[j].ClockSequence() })

		for i, k := range lpm {
			pos, theta := pls[u].CalcPos(s.theta, s.GetActualPosition(), i, len(lpm))

			child, ok := s.Children.Load(k)
			if !ok {
				s.log.Errorf("Space: UpdatePosition: failed to get space: %s", k)
				continue
			}
			if err := child.SetActualPosition(pos, theta, force); err != nil {
				s.log.Errorf("Space: UpdatePosition: failed to update position: %s", k)
			}
			if recursive {
				child.UpdateChildrenPosition(recursive, force)

			}
		}
	}
	fmt.Println("pls10", s.GetID())
	return nil
}
