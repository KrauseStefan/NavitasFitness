package AppEngineHelper

import (
	"time"
)

var timezoneOffsetMap = map[string]int{
	"AST":  -4,  // ATLANTIC STANDARD
	"EST":  -5,  // EASTERN STANDARD
	"EDT":  -4,  // EASTERN DAYLIGHT
	"CST":  -6,  // CENTRAL STANDARD
	"CDT":  -5,  // CENTRAL DAYLIGHT
	"MST":  -7,  // MOUNTAIN STANDARD
	"MDT":  -6,  // MOUNTAIN DAYLIGHT
	"PST":  -8,  // PACIFIC STANDARD
	"PDT":  -7,  // PACIFIC DAYLIGHT
	"AKST": -9,  // ALASKA
	"AKDT": -8,  // ALASKA DAYLIGHT
	"HST":  -10, // HAWAII STANDARD
	"HAST": -10, // HAWAII-ALEUTIAN STANDARD
	"HADT": -9,  // HAWAII-ALEUTIAN DAYLIGHT
	"SST":  -11, // SAMOA STANDARD
	"SDT":  -10, // SAMOA DAYLIGHT
	"CHST": +10, // CHAMORRO STANDARD
}

func LoadLocation(timeZone string) (*time.Location, error) {
	location, err := time.LoadLocation(timeZone)
	if err == nil {
		return location, nil
	}

	if utcOffset, ok := timezoneOffsetMap[timeZone]; ok {
		return time.FixedZone(timeZone, utcOffset*60*60), nil
	}

	return nil, err
}
