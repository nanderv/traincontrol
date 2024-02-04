package subcodes

const (
	Announce    = 1
	AckAnnounce = 2

	EEPROM_READ         = 3
	EEPROM_READ_RETURN  = 4
	EEPROM_WRITE        = 5
	EEPROM_WRITE_RETURN = 6

	LightOn     = 16
	AckLightOn  = 17
	LightOff    = 18
	AckLightOff = 19
)

const (
	StatusNew = 1

	StatusUpdating = 128
	StatusRunning  = 254
)
