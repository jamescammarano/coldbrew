package utils

import (
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/spaceweasel/promptui"
)

func Base64() string {
	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprint(time.Now().UTC().Nanosecond())))
}

func MergeMaps(mapArr []map[string]string) map[string]string {
	merged := map[string]string{}

	for _, t := range mapArr {
		for k, v := range t {
			merged[k] = v

		}
	}

	return merged
}

func PromptInput(label string) string {
	validator := func(input string) error {
		if input == "" {
			return errors.New("no input")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    label,
		Validate: validator,
	}

	res, err := prompt.Run()

	for err != nil {
		res, err = prompt.Run()
	}

	return res
}
