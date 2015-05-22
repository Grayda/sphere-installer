package main

import (
	"encoding/json"
	"fmt"

	"github.com/ninjasphere/go-ninja/model"
	"github.com/ninjasphere/go-ninja/suit"
)

// This file contains most of the code for the UI (i.e. what appears in the Labs)

type configService struct {
	app *InstallerConfig
}

// This function is common across all UIs, and is called by the Sphere. Shows our menu option on the main Labs screen
// The "c" bit at the start means this func is an extension of the configService struct (like prototyping, I think?)
func (c *configService) GetActions(request *model.ConfigurationRequest) (*[]suit.ReplyAction, error) {
	// What we're going to show
	var screen []suit.ReplyAction
	// Loop through all Orvibo devices. We do this so we can find an AllOne
	screen = append(screen, suit.ReplyAction{
		Name:        "",
		Label:       "Install Drivers",
		DisplayIcon: "plus",
	})

	// Return our screen to the sphere-ui for rendering
	return &screen, nil
}

// When you click on a ReplyAction button (e.g. the "Configure AllOne" button defined above), Configure is called. requests.Action == the "Name" portion of the ReplyAction
func (c *configService) Configure(request *model.ConfigurationRequest) (*suit.ConfigurationScreen, error) {
	fmt.Sprintf("Incoming configuration request. Action:%s Data:%s", request.Action, string(request.Data))

	switch request.Action {

	case "add": // Very similar to savegroup, but saves an IR code instead
		var vals map[string]string
		err := json.Unmarshal(request.Data, &vals)
		if err != nil {
			return c.error(fmt.Sprintf("Failed to unmarshal save config request %s: %s", request.Data, err))
		}

	case "": // Coming in from the main menu
		return c.menu()

	default: // Everything else

		// return c.list()
		return c.error(fmt.Sprintf("Unknown action: %s", request.Action))
	}

	// If this code runs, then we done fucked up, because default: didn't catch. When this code runs, the universe melts into a gigantic heap. But
	// removing this violates Apple guidelines and ensures the downfall of humanity (probably) so I don't want to risk it.
	// Then again, I could be making all this up. Do you want to remove it and try? ( ͡° ͜ʖ ͡°)
	return nil, nil
}

// So this function (which is an extension of the configService struct that suit (or Sphere-UI) requires) creates a box with a single "Okay" button and puts in a title and text
func (c *configService) confirm(title string, description string) (*suit.ConfigurationScreen, error) {
	// We create a new suit.ConfigurationScreen which is a whole page within the UI
	screen := suit.ConfigurationScreen{
		Title: title,
		Sections: []suit.Section{ // The UI lets us create sections for separating options. This line creates an array of sections
			suit.Section{ // And within that array of sections, a single section
				Contents: []suit.Typed{ // The contents of that section. I don't know what suit.Typed is. It's an interface, but asides from that, I don't know much else just yet
					suit.StaticText{ // Create some static text
						Title: "About this screen",
						Value: description,
					},
				},
			},
		},
		Actions: []suit.Typed{ // This configuration screen can show actionable buttons at the bottom. ReplyAction, as shown above, calls Configure. There is also CloseAction for cancel buttons
			suit.ReplyAction{
				Label:        "Okay",
				Name:         "list",
				DisplayClass: "success", // These are bootstrap classes (or rather, font-awesome classes). They are basically btn-*, where * is DisplayClass (e.g. btn-success)
				DisplayIcon:  "ok",      // Same as above. If you want to show fa-open-folder, you'd set DisplayIcon to "open-folder"
			},
		},
	}

	return &screen, nil
}

// Error! Same as above. It's a function that is added on to configService and displays an error message
func (c *configService) error(message string) (*suit.ConfigurationScreen, error) {

	return &suit.ConfigurationScreen{
		Sections: []suit.Section{
			suit.Section{
				Contents: []suit.Typed{
					suit.Alert{
						Title:        "Error",
						Subtitle:     message,
						DisplayClass: "danger",
					},
				},
			},
		},
		Actions: []suit.Typed{
			suit.ReplyAction{ // Shows a button we can click on. Takes us back to c.Configuration (reply.Action will be "list")
				Label:        "Cancel",
				Name:         "list",
				DisplayClass: "success",
				DisplayIcon:  "ok",
			},
		},
	}, nil
}

// Shows the UI to learn a new IR code
func (c *configService) menu() (*suit.ConfigurationScreen, error) {

	title := "Install New Driver" // Up here for readability

	screen := suit.ConfigurationScreen{
		Title: title,
		Sections: []suit.Section{ // New array of sections
			suit.Section{ // New section
				Contents: []suit.Typed{
					suit.StaticText{ // Some introductory text
						Title: "About this screen",
						Value: "This application lets you install a new driver or application from a .deb file. Please read all the instructions, as failure to do so may cause your Sphere to stop working, security risks on your network, or other undesired effects",
					},
					suit.InputText{ // Textbox
						Name:        "name",
						Before:      "Location of the .deb file",
						Placeholder: "http://example.com/my-driver.deb", // Placeholder is the faded text that appears inside a textbox, giving you a hint as to what to type in
						Value:       "",
					},
				},
			},
		},
		Actions: []suit.Typed{
			suit.CloseAction{ // This is not a CloseAction, because we want to go back to the list of IR codes, not back to the main menu. Hence why we use a ReplyAction with "list"
				Label: "Cancel",
			},
			suit.ReplyAction{
				Label:        "Begin installation",
				Name:         "add",
				DisplayClass: "success",
				DisplayIcon:  "star",
			},
		},
	}

	return &screen, nil
}

// Aye-aye, captain.
// Not actually needed (?)
func i(i int) *int {
	return &i
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
