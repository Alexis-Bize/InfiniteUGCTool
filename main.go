package main

import (
	"errors"
	"fmt"
	"infinite-bookmarker/internal"
	"infinite-bookmarker/internal/services/auth"
	"infinite-bookmarker/internal/shared/modules/helpers/identity"
	"log"
	"net/mail"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
)

const (
	AUTH_CMD = "auth"
	LOGOUT_CMD = "logout"
	BOOKMARK_CMD = "bookmark"
)

func main() {
	// MOVE TO prompt/main.go
	title := fmt.Sprintf("# %s (%s) #\n", internal.GetConfig().Title, internal.GetConfig().Version)
	os.Stdout.WriteString(fmt.Sprintf("%s\n", strings.Repeat("#", len(title) - 1)))
	os.Stdout.WriteString(title)
	os.Stdout.WriteString(fmt.Sprintf("%s\n", strings.Repeat("#", len(title) - 1)))

	currentIdentity, err := identity.GetOrCreateIdentity(identity.Identity{})
	if err != nil {
		log.Panic(err)
	}

	if currentIdentity != (identity.Identity{}) {
		os.Stdout.WriteString(fmt.Sprintf("✅ Welcome back, %s!\n", currentIdentity.XboxNetwork.Gamertag))
		os.Stdout.WriteString("Refreshing your Spartan Token...\n")
		return
	}

	if currentIdentity == (identity.Identity{}) {
		os.Stdout.WriteString("To bookmark any film, map, or mode, please authenticate with your Microsoft credentials.\n")
		prompt := promptui.Prompt{
			Label: "Would you like to proceed",
			IsConfirm: true,
		}
	
		shouldContinue, err := prompt.Run()
		if err != nil {
			fmt.Print(err)
			return
		} else if strings.ToLower(shouldContinue) == "n" {
			fmt.Printf("Good bye!")
			return
		}

		prompt = promptui.Prompt{
			Label: "Email address",
			Validate: func(input string) error {
				_, err := mail.ParseAddress(input)
				if err != nil {
					return errors.New("invalid email address")
				}
		
				return nil
			},
		}

		email, err := prompt.Run()
		if err != nil {
			fmt.Print(err)
			return
		}
	
		prompt = promptui.Prompt{
			Label: "Password",
			Mask: '*',
			Validate: func(input string) error {
				if len(input) == 0 {
					return errors.New("password can not be empty")
				}
		
				return nil
			},
		}
	
		password, err := prompt.Run()
		if err != nil {
			fmt.Print(err)
			return
		}

		os.Stdout.WriteString("Authenticating...\n")
		profile, spartanToken, err := auth.AuthenticateWithCredentials(email, password)
		if err != nil {
			fmt.Print(err)
			return
		}

		os.Stdout.WriteString(fmt.Sprintf("✅ Authenticated as %s", profile.Gamertag))
		identity.SaveIdentity(identity.Identity{
			User: identity.UserCredentials{
				Email: email,
				Password: password,
			},
			SpartanToken: identity.SpartanTokenDetails{
				Value: spartanToken,
				Expiration: "",
			},
			XboxNetwork: identity.XboxNetworkIdentity{
				Xuid: profile.Xuid,
				Gamertag: profile.Gamertag,
			},
		})
	}
}