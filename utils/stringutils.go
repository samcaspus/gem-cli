package utils

func GetMergedStringArgs(args []string) string {
	var mergedString string
	for _, arg := range args {
		mergedString += arg + " "
	}
	return mergedString

}
