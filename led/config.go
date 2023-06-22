package led

type Settings struct {
	FrameRate int
}

var packageSettings = Settings{
	FrameRate: 60,
}

func GetSettings() Settings {
	return packageSettings
}

func ApplySettings(settings Settings) error {
	packageSettings = settings
	return nil
}
