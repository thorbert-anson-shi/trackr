package handlers

import (
	"crypto/rand"
	"fmt"
	"time"

	"tobtoby/trackr/auth"
	"tobtoby/trackr/database"
	"tobtoby/trackr/generated"
	"tobtoby/trackr/logging"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgtype"
)

// CreateInviteLink generates a new invite link with a 5-minute expiry.
// @Summary      Create an invite link
// @Description  Admin-only endpoint that creates a new invite link with a randomly generated code, valid for 5 minutes.
// @Tags         invites
// @Produce      plain
// @Security     ApiKeyAuth
// @Success      201  {string}  string  "Full invite URL"
// @Failure      403  "Forbidden - only admins can create invite links"
// @Failure      500  "Internal Server Error"
// @Router       /api/v1/invite [post]
func CreateInviteLink(c fiber.Ctx) error {
	queries := generated.New(database.DB)
	user := c.Locals(auth.UserContextKey).(generated.User)

	if !user.IsAdmin {
		logging.GlobalLogger.Printf("Non-admin user (%d) tried to create an invite link", user.ID)
		return c.SendStatus(fiber.StatusForbidden)
	}

	invitation, err := queries.AddInvitation(c, generated.AddInvitationParams{
		Code: GenerateRandomCode(),
		ExpiryDate: pgtype.Timestamp{
			Time:             time.Now().Add(time.Duration(5) * time.Minute),
			InfinityModifier: 0,
			Valid:            true,
		},
	})
	if err != nil {
		logging.GlobalLogger.Printf("An error occurred while creating an invite link")
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	inviteLink := fmt.Sprintf("%s/invite?code=%s", c.BaseURL(), invitation.Code)

	return c.Status(fiber.StatusCreated).SendString(inviteLink)
}

func GenerateRandomCode() string {
	return rand.Text()
}
