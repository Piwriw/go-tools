package lark

import "fmt"

const larkTemplate = `<at user_id="%s">Tom</at>`

func getLarkContent(userIDS []string, context string) string {
	res := ""
	for _, userID := range userIDS {
		msg := fmt.Sprintf(larkTemplate, userID)
		res += msg
	}
	if res == "" {
		return context
	}
	return res + "\n" + context
}
