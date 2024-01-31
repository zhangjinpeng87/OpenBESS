// rules.go
// There are several physical rules in the battery management system we should consider.
// 1) Efficiency: when charging, we hope all the energy can be stored in the battery, and we don't want 
// to waste any energy when discharging. We don't want the energy lost in the form of heat.
// 2) Safety: when charging, we don't want the battery to be overcharged, and run into thermal issue.
// 3) Longevity: we want the battery to last as long as possible. We don't want the battery to be overcharged
// or overdischarged, which may shorten the battery's life.

package balancing

import (
)
