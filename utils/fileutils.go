package utils

import "os"

func GetRootFilePath(fileName string) string {
	return os.ExpandEnv("$HOME/" + fileName)
}

func DoesFileExist(filePath string) bool {
	_, err := os.Stat(GetRootFilePath(filePath))
	return !os.IsNotExist(err)
}

func WriteToFile(filePath string, data string) {
	file, err := os.Create(GetRootFilePath(filePath))
	if err != nil {
		panic(err)
	}
	defer file.Close()
	file.WriteString(data)
}

func ReadFile(filePath string) string {
	file, err := os.Open(GetRootFilePath(filePath))
	if err != nil {
		panic(err)
	}
	defer file.Close()
	buf := make([]byte, 1024)
	n, err := file.Read(buf)
	if err != nil {
		panic(err)
	}
	return string(buf[:n])
}
