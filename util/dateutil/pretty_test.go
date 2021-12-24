// Package dateutil
//
// @author: xwc1125
package dateutil

import (
	"testing"
	"time"
)

func TestPrettyAge_String(t *testing.T) {
	t.Log("ms", time.Now().UnixMilli())
	t.Log("s", time.Now().Unix())
	millisecondToTime := MillisecondToTime(1635050328000)
	t.Log("ms", millisecondToTime.UnixMilli())
	distanceTime := GetDistanceTime(1)
	t.Log(distanceTime)
	t.Log(Format(millisecondToTime, Default))
	age := PrettyAge(millisecondToTime)
	t.Log(age.String())
	duration := PrettyDuration(time.Since(millisecondToTime))
	t.Log(duration.String())
}
