package animation

import colorful "github.com/lucasb-eyer/go-colorful"

type Settings struct {
	FlowingNotes      *AnimationSettings
	VelocityBar       *AnimationSettings
	VelocityBarMirror *AnimationSettings
}

type AnimationSettings struct {
	*CommonSettings
}

type CommonSettings struct {
	LowerColor colorful.Color
	UpperColor colorful.Color
}

var DefaultCommonSettings = &CommonSettings{
	LowerColor: colorful.FastLinearRgb(0, 1, 0),
	UpperColor: colorful.FastLinearRgb(1, 0, 0),
}

var packageSettings = Settings{
	FlowingNotes: &AnimationSettings{
		CommonSettings: DefaultCommonSettings,
	},
	VelocityBar: &AnimationSettings{
		CommonSettings: DefaultCommonSettings,
	},
	VelocityBarMirror: &AnimationSettings{
		CommonSettings: DefaultCommonSettings,
	},
}

func GetSettings() Settings {
	return packageSettings
}

func ApplySettings(settings Settings) error {
	packageSettings = settings
	return nil
}
