package position

// OrbitalElements contains the Keplerian orbital elements at J2000.0 epoch
// and their rates of change per century
// Source: Paul Schlyter (stjarnhimlen.se) and JPL approximate positions
type OrbitalElements struct {
	// Longitude of ascending node (degrees)
	N float64
	// Rate of change of N (degrees/day)
	NRate float64
	// Inclination to ecliptic (degrees)
	I float64
	// Rate of change of I (degrees/day)
	IRate float64
	// Argument of perihelion (degrees)
	W float64
	// Rate of change of W (degrees/day)
	WRate float64
	// Semi-major axis (AU, or Earth radii for Moon)
	A float64
	// Rate of change of A (AU/day)
	ARate float64
	// Eccentricity
	E float64
	// Rate of change of E (per day)
	ERate float64
	// Mean anomaly at epoch (degrees)
	M float64
	// Mean daily motion (degrees/day)
	MRate float64
}

// AtDay returns computed orbital elements for a given day number from J2000
func (o OrbitalElements) AtDay(d float64) ComputedElements {
	return ComputedElements{
		N: NormalizeAngle(o.N + o.NRate*d),
		I: o.I + o.IRate*d,
		W: NormalizeAngle(o.W + o.WRate*d),
		A: o.A + o.ARate*d,
		E: o.E + o.ERate*d,
		M: NormalizeAngle(o.M + o.MRate*d),
	}
}

// ComputedElements are the orbital elements computed for a specific date
type ComputedElements struct {
	N float64 // Longitude of ascending node
	I float64 // Inclination
	W float64 // Argument of perihelion
	A float64 // Semi-major axis
	E float64 // Eccentricity
	M float64 // Mean anomaly
}

// PlanetElements contains orbital elements for all celestial bodies
// Elements are for J2000.0 epoch with rates per day
// Source: Paul Schlyter's "Computing planetary positions"
var PlanetElements = map[CelestialBody]OrbitalElements{
	// Sun (actually Earth's orbit seen from geocentric perspective)
	Sun: {
		N: 0.0, NRate: 0.0,
		I: 0.0, IRate: 0.0,
		W: 282.9404, WRate: 4.70935e-5,
		A: 1.000000, ARate: 0.0,
		E: 0.016709, ERate: -1.151e-9,
		M: 356.0470, MRate: 0.9856002585,
	},
	// Moon
	Moon: {
		N: 125.1228, NRate: -0.0529538083,
		I: 5.1454, IRate: 0.0,
		W: 318.0634, WRate: 0.1643573223,
		A: 60.2666, ARate: 0.0, // Earth radii
		E: 0.054900, ERate: 0.0,
		M: 115.3654, MRate: 13.0649929509,
	},
	// Mercury
	Mercury: {
		N: 48.3313, NRate: 3.24587e-5,
		I: 7.0047, IRate: 5.00e-8,
		W: 29.1241, WRate: 1.01444e-5,
		A: 0.387098, ARate: 0.0,
		E: 0.205635, ERate: 5.59e-10,
		M: 168.6562, MRate: 4.0923344368,
	},
	// Venus
	Venus: {
		N: 76.6799, NRate: 2.46590e-5,
		I: 3.3946, IRate: 2.75e-8,
		W: 54.8910, WRate: 1.38374e-5,
		A: 0.723330, ARate: 0.0,
		E: 0.006773, ERate: -1.302e-9,
		M: 48.0052, MRate: 1.6021302244,
	},
	// Mars
	Mars: {
		N: 49.5574, NRate: 2.11081e-5,
		I: 1.8497, IRate: -1.78e-8,
		W: 286.5016, WRate: 2.92961e-5,
		A: 1.523688, ARate: 0.0,
		E: 0.093405, ERate: 2.516e-9,
		M: 18.6021, MRate: 0.5240207766,
	},
	// Jupiter
	Jupiter: {
		N: 100.4542, NRate: 2.76854e-5,
		I: 1.3030, IRate: -1.557e-7,
		W: 273.8777, WRate: 1.64505e-5,
		A: 5.20256, ARate: 0.0,
		E: 0.048498, ERate: 4.469e-9,
		M: 19.8950, MRate: 0.0830853001,
	},
	// Saturn
	Saturn: {
		N: 113.6634, NRate: 2.38980e-5,
		I: 2.4886, IRate: -1.081e-7,
		W: 339.3939, WRate: 2.97661e-5,
		A: 9.55475, ARate: 0.0,
		E: 0.055546, ERate: -9.499e-9,
		M: 316.9670, MRate: 0.0334442282,
	},
	// Uranus
	Uranus: {
		N: 74.0005, NRate: 1.3978e-5,
		I: 0.7733, IRate: 1.9e-8,
		W: 96.6612, WRate: 3.0565e-5,
		A: 19.18171, ARate: -1.55e-8,
		E: 0.047318, ERate: 7.45e-9,
		M: 142.5905, MRate: 0.011725806,
	},
	// Neptune
	Neptune: {
		N: 131.7806, NRate: 3.0173e-5,
		I: 1.7700, IRate: -2.55e-7,
		W: 272.8461, WRate: -6.027e-6,
		A: 30.05826, ARate: 3.313e-8,
		E: 0.008606, ERate: 2.15e-9,
		M: 260.2471, MRate: 0.005995147,
	},
	// Pluto (simplified elements)
	Pluto: {
		N: 110.30347, NRate: 0.0,
		I: 17.14175, IRate: 0.0,
		W: 224.06676, WRate: 0.0,
		A: 39.48168677, ARate: 0.0,
		E: 0.24880766, ERate: 0.0,
		M: 238.92881, MRate: 0.003971354,
	},
	// Chiron
	Chiron: {
		N: 209.3851, NRate: 0.0,
		I: 6.9311, IRate: 0.0,
		W: 339.5574, WRate: 0.0,
		A: 13.6697, ARate: 0.0,
		E: 0.3792, ERate: 0.0,
		M: 72.3100, MRate: 0.01953663,
	},
	// Ceres
	Ceres: {
		N: 80.3932, NRate: 0.0,
		I: 10.5935, IRate: 0.0,
		W: 73.5968, WRate: 0.0,
		A: 2.7658, ARate: 0.0,
		E: 0.0758, ERate: 0.0,
		M: 113.4104, MRate: 0.21408169,
	},
	// Pallas
	Pallas: {
		N: 173.0962, NRate: 0.0,
		I: 34.8413, IRate: 0.0,
		W: 310.0474, WRate: 0.0,
		A: 2.7716, ARate: 0.0,
		E: 0.2313, ERate: 0.0,
		M: 78.2287, MRate: 0.21343011,
	},
	// Juno
	Juno: {
		N: 169.8712, NRate: 0.0,
		I: 12.9717, IRate: 0.0,
		W: 248.4100, WRate: 0.0,
		A: 2.6691, ARate: 0.0,
		E: 0.2562, ERate: 0.0,
		M: 18.2795, MRate: 0.22610627,
	},
	// Vesta
	Vesta: {
		N: 103.8513, NRate: 0.0,
		I: 7.1340, IRate: 0.0,
		W: 151.1983, WRate: 0.0,
		A: 2.3615, ARate: 0.0,
		E: 0.0887, ERate: 0.0,
		M: 169.1467, MRate: 0.27154186,
	},
}

// Obliquity of the ecliptic at J2000.0 (degrees)
const Obliquity = 23.4393
