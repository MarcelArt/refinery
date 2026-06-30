package common

import (
	"strings"

	"git.bangmarcel.art/marcel/arrays"
	"github.com/MarcelArt/refinery/internal/configs"
	"github.com/MarcelArt/refinery/internal/enums"
)

func CheckEmails(emails ...string) error {
	if configs.Env.ServerENV == "prod" {
		dummyEmail := arrays.Find(emails, func(email string) bool {
			return strings.Contains(email, "@yopmail.com")
		})
		if dummyEmail != nil {
			return enums.ErrDummyEmailOnProd
		}
	}
	return nil
}
