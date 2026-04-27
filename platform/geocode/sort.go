package geocode

import (
	"github.com/Gleipnir-Technology/nidus-sync/platform/types"
	"github.com/Gleipnir-Technology/nidus-sync/stadia"
)

type SortAddressByDistance struct {
	Addresses []types.Address
	Location  types.Location
}

func (s SortAddressByDistance) Len() int { return len(s.Addresses) }
func (s SortAddressByDistance) Swap(i, j int) {
	s.Addresses[i], s.Addresses[j] = s.Addresses[j], s.Addresses[i]
}
func (s SortAddressByDistance) Less(i, j int) bool {
	ai := s.Addresses[i]
	aj := s.Addresses[j]
	if ai.Location == nil || (ai.Location.Latitude == 0 && ai.Location.Longitude == 0) {
		if aj.Location == nil || (aj.Location.Latitude == 0 && aj.Location.Longitude == 0) {
			return ai.Raw > aj.Raw
		}
		return false
	} else if aj.Location == nil || (aj.Location.Latitude == 0 && aj.Location.Longitude == 0) {
		return true
	}
	di := types.LocationDistance(s.Location, *ai.Location)
	dj := types.LocationDistance(s.Location, *ai.Location)
	return di < dj
}

type SortFeaturesByDistance []stadia.GeocodeFeature

func (s SortFeaturesByDistance) Len() int      { return len(s) }
func (s SortFeaturesByDistance) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s SortFeaturesByDistance) Less(i, j int) bool {
	fi := s[i].Properties.Distance
	fj := s[j].Properties.Distance
	if fi == nil {
		if fj == nil {
			return s[i].Properties.GID < s[j].Properties.GID
		}
		return false
	} else if fj == nil {
		return true
	}
	return *fi < *fj
}
