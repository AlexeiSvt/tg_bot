package handlers

import (
	"strings"
	"fmt"
	"tgbot/internal/models"
)

func formatSummary(u *models.User) string {
	var summary strings.Builder; summary.WriteString(fmt.Sprintf(SummaryHeader,
		u.FirstName, u.LastName, u.Class))

	for game, gd := range u.Disciplines {
		if game == "Chess" {
			summary .WriteString(fmt.Sprintf("   %s: %s\n", game, gd.Nick))
		} else {
			summary .WriteString(fmt.Sprintf("   %s: %s | %s\n", game, gd.Nick, gd.Tag))
		}
	}

	summary .WriteString(SummaryFooter)
	return summary.String()
}

