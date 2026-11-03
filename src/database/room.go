package database

import (
	"os"

	"github.com/D-Gine/EngineWorker/src/debug"
)

type RoomIds struct {
	Campaign string
	Rulest   string
}

var Room RoomIds

func RetreiveRoomInformations() {
	debug.PrintMsg("retreiving rooms informations...")
	Room.Rulest = os.Getenv("RULESET_ID")
	Room.Campaign = os.Getenv("CAMPAIGN_ID")
	if Room.Rulest == "" {
		debug.PrintErrorMsg("RULESET_ID variable not found in environt")
	} else if Room.Campaign == "" {
		debug.PrintErrorMsg("CAMPAIGN_ID variable not found in environt")
	} else {
		debug.PrintSuccess("room informations found")
		debug.PrintMsg("running campaign on\n\t- Ruleset: " + Room.Rulest + "\n\t- Campaign: " + Room.Campaign)
	}
}
