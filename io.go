package neurgo

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func WriteStringToFile(value string, filepath string) error {
	outputFile, outputError := os.OpenFile(filepath,
		os.O_WRONLY|os.O_CREATE,
		0666)
	if outputError != nil {
		return outputError
	}
	defer outputFile.Close()
	outputWriter := bufio.NewWriter(outputFile)
	outputWriter.WriteString(value)
	outputWriter.Flush()
	return nil
}

func JsonString(v interface{}) string {
	json, err := json.Marshal(v)
	if err != nil {
		log.Fatal(err)
		return err.Error()
	}
	return fmt.Sprintf("%s", json)
}
