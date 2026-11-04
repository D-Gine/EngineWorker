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

func verify_ruleset() error {
	var id string
	err := DB.QueryRow(`
		SELECT
			r.id
		FROM rulesets r
		WHERE r.id=$1
	`, Room.Rulest).Scan(&id)
	if err != nil {
		return err
	}
	return nil
}

func verify_campaign() error {
	return nil
}

func RetreiveRoomInformations() {
	debug.PrintMsg("retreiving rooms informations...")
	Room.Rulest = os.Getenv("RULESET_ID")
	Room.Campaign = os.Getenv("CAMPAIGN_ID")
	var exit_err bool = false
	if Room.Rulest == "" {
		debug.PrintErrorMsg("RULESET_ID variable not found in environt")
		exit_err = true
	} else {
		err := verify_ruleset()
		if err != nil {
			debug.PrintError(err)
			exit_err = true
		}
	}
	if Room.Campaign == "" {
		debug.PrintErrorMsg("CAMPAIGN_ID variable not found in environt")
		exit_err = true
	} else {
		err := verify_campaign()
		if err != nil {
			debug.PrintError(err)
			exit_err = true
		}
	}
	if exit_err {
		return
	}
	debug.PrintSuccess("room informations found")
	debug.PrintMsg("running campaign on\n\t- Ruleset: " + Room.Rulest + "\n\t- Campaign: " + Room.Campaign)
}
