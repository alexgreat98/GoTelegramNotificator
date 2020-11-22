package telegram

import "fmt"

func ListAll(list []string) string {
	var result string
	for _, item := range list {
		result = fmt.Sprintf("%s\n%s", result, item)
	}
	return result
}
