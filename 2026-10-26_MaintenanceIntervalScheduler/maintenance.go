// Package maintenance decides how soon a piece of farm equipment is next
// due for service, based on how hard that equipment type runs the fleet.
package maintenance

import (
	"errors"
	"time"
)

// baseMaintenanceIntervalDays is the standard service interval before any
// equipment-specific adjustment is applied.
const baseMaintenanceIntervalDays = 90

// MaintenanceStrategy knows how to adjust the base service interval for one
// equipment type's usage profile.
type MaintenanceStrategy interface {
	AdjustInterval(baseDays int) int
}

// tractorStrategy covers row-crop tractors, which run long hours during
// planting and harvest and need servicing sooner than the baseline.
type tractorStrategy struct{}

func (tractorStrategy) AdjustInterval(baseDays int) int {
	return baseDays - 10
}

// harvesterStrategy covers combine harvesters, which see concentrated,
// heavy seasonal use and wear faster than anything else in the fleet.
type harvesterStrategy struct{}

func (harvesterStrategy) AdjustInterval(baseDays int) int {
	return baseDays - 20
}

// irrigationPumpStrategy covers stationary irrigation pumps, which run
// almost continuously but under much lighter mechanical load.
type irrigationPumpStrategy struct{}

func (irrigationPumpStrategy) AdjustInterval(baseDays int) int {
	return baseDays + 15
}

// newMaintenanceStrategy looks up the strategy for an equipment type.
func newMaintenanceStrategy(equipmentType string) MaintenanceStrategy {
	switch equipmentType {
	case "tractor":
		return tractorStrategy{}
	case "harvester":
		return harvesterStrategy{}
	case "irrigation_pump":
		return irrigationPumpStrategy{}
	default:
		return nil
	}
}

// ScheduleNextService returns the date a piece of equipment is next due
// for maintenance, given when it was last serviced. The interval off the
// 90-day baseline depends on how that equipment type is used.
func ScheduleNextService(equipmentType string, lastServiced time.Time) (time.Time, error) {
	if equipmentType == "" {
		return time.Time{}, errors.New("equipment type is required")
	}

	strategy := newMaintenanceStrategy(equipmentType)
	interval := strategy.AdjustInterval(baseMaintenanceIntervalDays)

	return lastServiced.AddDate(0, 0, interval), nil
}
