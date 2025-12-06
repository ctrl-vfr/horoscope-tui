package i18n

var messages = map[Lang]map[string]string{
	EN: {
		// Form
		"FormBirthDate":            "Birth date",
		"FormBirthDateDesc":        "Format: DD/MM/YYYY",
		"FormBirthDatePlaceholder": "21/03/1990",
		"FormTransitDate":          "Transit date",
		"FormTransitDateDesc":      "For predictions (DD/MM/YYYY)",
		"FormTransitDatePlaceholder": "01/01/2025",
		"FormQuestion":             "Ask the oracle your question",
		"FormQuestionDesc":         "Ex: Should I accept this job? Is it the right time to...?",
		"FormQuestionPlaceholder":  "What's on your mind?",

		// Validation
		"ValidationRequired":      "date required",
		"ValidationInvalidFormat": "invalid format (DD/MM/YYYY)",

		// Status
		"StatusLoading":           "Loading...",
		"StatusGeocoding":         "Looking up coordinates...",
		"StatusGeocodingError":    "Geocoding error: ",
		"StatusCalculating":       "Calculating chart...",
		"StatusGeneratingWheel":   "Generating wheel...",
		"StatusTransmittingImage": "Transmitting image...",
		"StatusWaitingData":       "Waiting for data...",
		"StatusWaitingNatal":      "Waiting for natal chart...",
		"StatusError":             "Error: ",

		// Missing city
		"MissingCityError": "HOROSCOPE_CITY variable missing",
		"MissingCityHint":  "Set the environment variable:",

		// Header
		"HeaderTitle": "MY ORACLE",

		// Interpretation
		"InterpTitle":   "Interpretation",
		"InterpLoading": "(loading...)",
		"InterpError":   "[Error]",

		// Positions
		"PositionPlanet":   "Planet",
		"PositionNatal":    "Natal",
		"PositionTransit":  "Transit",
		"PositionPosition": "Position",
		"PositionTransits": "Transits",
		"PositionBoth":     "Natal / Transits",

		// Wheel
		"WheelPlaceholder": "[ Zodiac wheel ]\n(Kitty/resvg required)",

		// Navigation
		"NavNavigate":    " navigate",
		"NavScroll":      " scroll",
		"NavNewQuestion": " new question",
		"NavQuit":        " quit",

		// Elements
		"ElementFire":  "Fire",
		"ElementEarth": "Earth",
		"ElementAir":   "Air",
		"ElementWater": "Water",
		"ElementCount": "%d planets",

		// Weekdays
		"WeekdaySunday":    "Sunday",
		"WeekdayMonday":    "Monday",
		"WeekdayTuesday":   "Tuesday",
		"WeekdayWednesday": "Wednesday",
		"WeekdayThursday":  "Thursday",
		"WeekdayFriday":    "Friday",
		"WeekdaySaturday":  "Saturday",

		// Prompt builder
		"PromptQuestion":        "QUESTION",
		"PromptDefaultQuestion": "Give me a general cosmic reading for today based on my chart.",
		"PromptToday":           "Today",
		"PromptTransitsTitle":   "TODAY'S TRANSITS (current positions)",
		"PromptNatalTitle":      "NATAL CHART (birth positions)",
		"PromptBirthDate":       "Birth date",
		"PromptLocation":        "Location",
		"PromptPlanetPositions": "Planetary positions",
		"PromptMajorAspects":    "Major aspects",
		"PromptElementDist":     "Element distribution",
		"PromptRetrograde":      " (RETROGRADE)",
		"PromptOrb":             "orb",
	},

	FR: {
		// Form
		"FormBirthDate":            "Date de naissance",
		"FormBirthDateDesc":        "Format: JJ/MM/AAAA",
		"FormBirthDatePlaceholder": "21/03/1990",
		"FormTransitDate":          "Date de transit",
		"FormTransitDateDesc":      "Pour les prédictions (JJ/MM/AAAA)",
		"FormTransitDatePlaceholder": "01/01/2025",
		"FormQuestion":             "Pose ta question à l'oracle",
		"FormQuestionDesc":         "Ex: Dois-je accepter ce job? C'est le bon moment pour...?",
		"FormQuestionPlaceholder":  "Qu'est-ce qui te tracasse?",

		// Validation
		"ValidationRequired":      "date requise",
		"ValidationInvalidFormat": "format invalide (JJ/MM/AAAA)",

		// Status
		"StatusLoading":           "Chargement...",
		"StatusGeocoding":         "Recherche des coordonnées...",
		"StatusGeocodingError":    "Erreur géocodage: ",
		"StatusCalculating":       "Calcul du thème...",
		"StatusGeneratingWheel":   "Génération de la roue...",
		"StatusTransmittingImage": "Transmission de l'image...",
		"StatusWaitingData":       "En attente des données...",
		"StatusWaitingNatal":      "En attente du thème natal...",
		"StatusError":             "Erreur: ",

		// Missing city
		"MissingCityError": "Variable HOROSCOPE_CITY manquante",
		"MissingCityHint":  "Définissez la variable d'environnement:",

		// Header
		"HeaderTitle": "MON ORACLE",

		// Interpretation
		"InterpTitle":   "Interprétation",
		"InterpLoading": "(chargement...)",
		"InterpError":   "[Erreur]",

		// Positions
		"PositionPlanet":   "Planète",
		"PositionNatal":    "Natal",
		"PositionTransit":  "Transit",
		"PositionPosition": "Position",
		"PositionTransits": "Transits",
		"PositionBoth":     "Natal / Transits",

		// Wheel
		"WheelPlaceholder": "[ Roue zodiacale ]\n(Kitty/resvg requis)",

		// Navigation
		"NavNavigate":    " naviguer",
		"NavScroll":      " défiler",
		"NavNewQuestion": " nouvelle question",
		"NavQuit":        " quitter",

		// Elements
		"ElementFire":  "Feu",
		"ElementEarth": "Terre",
		"ElementAir":   "Air",
		"ElementWater": "Eau",
		"ElementCount": "%d planètes",

		// Weekdays
		"WeekdaySunday":    "Dimanche",
		"WeekdayMonday":    "Lundi",
		"WeekdayTuesday":   "Mardi",
		"WeekdayWednesday": "Mercredi",
		"WeekdayThursday":  "Jeudi",
		"WeekdayFriday":    "Vendredi",
		"WeekdaySaturday":  "Samedi",

		// Prompt builder
		"PromptQuestion":        "QUESTION",
		"PromptDefaultQuestion": "Donne-moi une lecture cosmique générale pour aujourd'hui basée sur mon thème.",
		"PromptToday":           "Aujourd'hui",
		"PromptTransitsTitle":   "TRANSITS DU JOUR (positions actuelles)",
		"PromptNatalTitle":      "THÈME NATAL (positions à la naissance)",
		"PromptBirthDate":       "Date de naissance",
		"PromptLocation":        "Lieu",
		"PromptPlanetPositions": "Positions planétaires",
		"PromptMajorAspects":    "Aspects majeurs",
		"PromptElementDist":     "Répartition des éléments",
		"PromptRetrograde":      " (RÉTROGRADE)",
		"PromptOrb":             "orbe",
	},

	ES: {
		// Form
		"FormBirthDate":            "Fecha de nacimiento",
		"FormBirthDateDesc":        "Formato: DD/MM/AAAA",
		"FormBirthDatePlaceholder": "21/03/1990",
		"FormTransitDate":          "Fecha de tránsito",
		"FormTransitDateDesc":      "Para predicciones (DD/MM/AAAA)",
		"FormTransitDatePlaceholder": "01/01/2025",
		"FormQuestion":             "Hazle tu pregunta al oráculo",
		"FormQuestionDesc":         "Ej: ¿Debo aceptar este trabajo? ¿Es el momento adecuado para...?",
		"FormQuestionPlaceholder":  "¿Qué te preocupa?",

		// Validation
		"ValidationRequired":      "fecha requerida",
		"ValidationInvalidFormat": "formato inválido (DD/MM/AAAA)",

		// Status
		"StatusLoading":           "Cargando...",
		"StatusGeocoding":         "Buscando coordenadas...",
		"StatusGeocodingError":    "Error de geocodificación: ",
		"StatusCalculating":       "Calculando carta...",
		"StatusGeneratingWheel":   "Generando rueda...",
		"StatusTransmittingImage": "Transmitiendo imagen...",
		"StatusWaitingData":       "Esperando datos...",
		"StatusWaitingNatal":      "Esperando carta natal...",
		"StatusError":             "Error: ",

		// Missing city
		"MissingCityError": "Variable HOROSCOPE_CITY faltante",
		"MissingCityHint":  "Configure la variable de entorno:",

		// Header
		"HeaderTitle": "MI ORÁCULO",

		// Interpretation
		"InterpTitle":   "Interpretación",
		"InterpLoading": "(cargando...)",
		"InterpError":   "[Error]",

		// Positions
		"PositionPlanet":   "Planeta",
		"PositionNatal":    "Natal",
		"PositionTransit":  "Tránsito",
		"PositionPosition": "Posición",
		"PositionTransits": "Tránsitos",
		"PositionBoth":     "Natal / Tránsitos",

		// Wheel
		"WheelPlaceholder": "[ Rueda zodiacal ]\n(Kitty/resvg requerido)",

		// Navigation
		"NavNavigate":    " navegar",
		"NavScroll":      " desplazar",
		"NavNewQuestion": " nueva pregunta",
		"NavQuit":        " salir",

		// Elements
		"ElementFire":  "Fuego",
		"ElementEarth": "Tierra",
		"ElementAir":   "Aire",
		"ElementWater": "Agua",
		"ElementCount": "%d planetas",

		// Weekdays
		"WeekdaySunday":    "Domingo",
		"WeekdayMonday":    "Lunes",
		"WeekdayTuesday":   "Martes",
		"WeekdayWednesday": "Miércoles",
		"WeekdayThursday":  "Jueves",
		"WeekdayFriday":    "Viernes",
		"WeekdaySaturday":  "Sábado",

		// Prompt builder
		"PromptQuestion":        "PREGUNTA",
		"PromptDefaultQuestion": "Dame una lectura cósmica general para hoy basada en mi carta.",
		"PromptToday":           "Hoy",
		"PromptTransitsTitle":   "TRÁNSITOS DE HOY (posiciones actuales)",
		"PromptNatalTitle":      "CARTA NATAL (posiciones al nacer)",
		"PromptBirthDate":       "Fecha de nacimiento",
		"PromptLocation":        "Lugar",
		"PromptPlanetPositions": "Posiciones planetarias",
		"PromptMajorAspects":    "Aspectos mayores",
		"PromptElementDist":     "Distribución de elementos",
		"PromptRetrograde":      " (RETRÓGRADO)",
		"PromptOrb":             "orbe",
	},

	DE: {
		// Form
		"FormBirthDate":            "Geburtsdatum",
		"FormBirthDateDesc":        "Format: TT/MM/JJJJ",
		"FormBirthDatePlaceholder": "21/03/1990",
		"FormTransitDate":          "Transitdatum",
		"FormTransitDateDesc":      "Für Vorhersagen (TT/MM/JJJJ)",
		"FormTransitDatePlaceholder": "01/01/2025",
		"FormQuestion":             "Stelle dem Orakel deine Frage",
		"FormQuestionDesc":         "Z.B.: Soll ich diesen Job annehmen? Ist es der richtige Zeitpunkt für...?",
		"FormQuestionPlaceholder":  "Was beschäftigt dich?",

		// Validation
		"ValidationRequired":      "Datum erforderlich",
		"ValidationInvalidFormat": "Ungültiges Format (TT/MM/JJJJ)",

		// Status
		"StatusLoading":           "Laden...",
		"StatusGeocoding":         "Suche Koordinaten...",
		"StatusGeocodingError":    "Geokodierungsfehler: ",
		"StatusCalculating":       "Berechne Horoskop...",
		"StatusGeneratingWheel":   "Erzeuge Rad...",
		"StatusTransmittingImage": "Übertrage Bild...",
		"StatusWaitingData":       "Warte auf Daten...",
		"StatusWaitingNatal":      "Warte auf Geburtshoroskop...",
		"StatusError":             "Fehler: ",

		// Missing city
		"MissingCityError": "Variable HOROSCOPE_CITY fehlt",
		"MissingCityHint":  "Setzen Sie die Umgebungsvariable:",

		// Header
		"HeaderTitle": "MEIN ORAKEL",

		// Interpretation
		"InterpTitle":   "Deutung",
		"InterpLoading": "(lädt...)",
		"InterpError":   "[Fehler]",

		// Positions
		"PositionPlanet":   "Planet",
		"PositionNatal":    "Natal",
		"PositionTransit":  "Transit",
		"PositionPosition": "Position",
		"PositionTransits": "Transite",
		"PositionBoth":     "Natal / Transite",

		// Wheel
		"WheelPlaceholder": "[ Tierkreisrad ]\n(Kitty/resvg erforderlich)",

		// Navigation
		"NavNavigate":    " navigieren",
		"NavScroll":      " scrollen",
		"NavNewQuestion": " neue Frage",
		"NavQuit":        " beenden",

		// Elements
		"ElementFire":  "Feuer",
		"ElementEarth": "Erde",
		"ElementAir":   "Luft",
		"ElementWater": "Wasser",
		"ElementCount": "%d Planeten",

		// Weekdays
		"WeekdaySunday":    "Sonntag",
		"WeekdayMonday":    "Montag",
		"WeekdayTuesday":   "Dienstag",
		"WeekdayWednesday": "Mittwoch",
		"WeekdayThursday":  "Donnerstag",
		"WeekdayFriday":    "Freitag",
		"WeekdaySaturday":  "Samstag",

		// Prompt builder
		"PromptQuestion":        "FRAGE",
		"PromptDefaultQuestion": "Gib mir eine allgemeine kosmische Deutung für heute basierend auf meinem Horoskop.",
		"PromptToday":           "Heute",
		"PromptTransitsTitle":   "HEUTIGE TRANSITE (aktuelle Positionen)",
		"PromptNatalTitle":      "GEBURTSHOROSKOP (Geburtspositionen)",
		"PromptBirthDate":       "Geburtsdatum",
		"PromptLocation":        "Ort",
		"PromptPlanetPositions": "Planetenpositionen",
		"PromptMajorAspects":    "Hauptaspekte",
		"PromptElementDist":     "Elementverteilung",
		"PromptRetrograde":      " (RÜCKLÄUFIG)",
		"PromptOrb":             "Orbis",
	},
}
