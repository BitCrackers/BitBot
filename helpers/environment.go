package helpers

type Settings struct {
	Token string
	Debug bool
	Owner string
	Guild string
}

var s Settings

func SetSettings(t string, d bool, o string, g string) {
	s.Token = t
	s.Debug = d
	s.Owner = o
	s.Guild = g
}

func GetSettings() Settings {
	return s
}
