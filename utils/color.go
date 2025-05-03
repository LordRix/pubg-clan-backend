package utils

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
)

func Green(msg string) string {
	return colorGreen + msg + colorReset
}

func Red(msg string) string {
	return colorRed + msg + colorReset
}

func Yellow(msg string) string {
	return colorYellow + msg + colorReset
}

func Blue(msg string) string {
	return colorBlue + msg + colorReset
}
