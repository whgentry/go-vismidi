package midi

type Settings struct {
	Port      int
	KeyOffset int
	KeyCount  int
}

var packageSettings = Settings{
	Port:      0,
	KeyOffset: 21,
	KeyCount:  88,
}

func GetSettings() Settings {
	return packageSettings
}

func ApplySettings(settings Settings) error {
	packageSettings = settings
	return nil
}
