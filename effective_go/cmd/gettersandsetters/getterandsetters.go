package gettersandsetters

import "fmt"

type ParkingLot struct {
	owner string
}

func New(s string) *ParkingLot {
	return &ParkingLot{owner: ""}
}

func (pl *ParkingLot) Owner() string {
	return pl.owner
}

func (pl *ParkingLot) SetOwner(s string) error {
	if len(s) > 3 {
		pl.owner = s
		return nil
	}

	return fmt.Errorf("Name %s is invalid", s)
}
