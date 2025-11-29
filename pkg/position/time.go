package position

import "time"

// J2000 is the Julian Day number for January 1, 2000 at 12:00 TT
const J2000 = 2451545.0

// DayNumber calculates the day number relative to J2000.0 epoch
// This is used as the primary time parameter in orbital calculations
func DayNumber(t time.Time) float64 {
	return JulianDay(t) - J2000
}

// JulianDay converts a time.Time to Julian Day number
// Formula from Meeus, Astronomical Algorithms
func JulianDay(t time.Time) float64 {
	year := t.Year()
	month := int(t.Month())
	day := float64(t.Day()) + float64(t.Hour())/24.0 +
		float64(t.Minute())/1440.0 + float64(t.Second())/86400.0

	if month <= 2 {
		year--
		month += 12
	}

	a := year / 100
	b := 2 - a + a/4

	jd := float64(int(365.25*float64(year+4716))) +
		float64(int(30.6001*float64(month+1))) +
		day + float64(b) - 1524.5

	return jd
}

// JulianDayToTime converts a Julian Day number back to time.Time (UTC)
func JulianDayToTime(jd float64) time.Time {
	z := int(jd + 0.5)
	f := jd + 0.5 - float64(z)

	var a int
	if z < 2299161 {
		a = z
	} else {
		alpha := int((float64(z) - 1867216.25) / 36524.25)
		a = z + 1 + alpha - alpha/4
	}

	b := a + 1524
	c := int((float64(b) - 122.1) / 365.25)
	d := int(365.25 * float64(c))
	e := int(float64(b-d) / 30.6001)

	day := b - d - int(30.6001*float64(e)) + int(f)
	var month int
	if e < 14 {
		month = e - 1
	} else {
		month = e - 13
	}
	var year int
	if month > 2 {
		year = c - 4716
	} else {
		year = c - 4715
	}

	fracDay := f
	hour := int(fracDay * 24)
	fracDay = fracDay*24 - float64(hour)
	minute := int(fracDay * 60)
	fracDay = fracDay*60 - float64(minute)
	second := int(fracDay * 60)

	return time.Date(year, time.Month(month), day, hour, minute, second, 0, time.UTC)
}

// LocalSiderealTime calculates the Local Sidereal Time in degrees
// for a given Julian Day and geographic longitude
func LocalSiderealTime(jd float64, longitude float64) float64 {
	// Greenwich Mean Sidereal Time at 0h UT
	t := (jd - J2000) / 36525.0
	gmst := 280.46061837 + 360.98564736629*(jd-J2000) +
		0.000387933*t*t - t*t*t/38710000.0

	lst := gmst + longitude
	return NormalizeAngle(lst)
}

// NormalizeAngle reduces an angle to the range [0, 360).
func NormalizeAngle(angle float64) float64 {
	for angle < 0 {
		angle += 360
	}
	for angle >= 360 {
		angle -= 360
	}
	return angle
}

// NormalizeMotion reduces an angular motion to the range (-180, 180].
func NormalizeMotion(motion float64) float64 {
	if motion > 180 {
		return motion - 360
	}
	if motion < -180 {
		return motion + 360
	}
	return motion
}

// DegreesToRadians converts degrees to radians
func DegreesToRadians(deg float64) float64 {
	return deg * 0.017453292519943295 // pi/180
}

// RadiansToDegrees converts radians to degrees
func RadiansToDegrees(rad float64) float64 {
	return rad * 57.29577951308232 // 180/pi
}
