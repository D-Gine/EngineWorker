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
	// TO DO after db
	// var id string
	// err = DB.QueryRow(`
	// 	SELECT
	// 		c.id
	// 	FROM campaign c
	// 	WHERE c.id=$1
	// `, Room.Campaign).Scan(&id)
	// if err != nil {
	// 	return err
	// }
	return nil
}

func RetreiveRoomInformations() {
	debug.Message("retreiving rooms informations...")
	Room.Rulest = os.Getenv("RULESET_ID")
	Room.Campaign = os.Getenv("CAMPAIGN_ID")
	exit_err := false
	if Room.Rulest == "" {
		debug.ErrorMessage("RULESET_ID variable not found in environt")
		exit_err = true
	} else {
		err := verify_ruleset()
		if err != nil {
			debug.Error(err)
			exit_err = true
		}
	}
	if Room.Campaign == "" {
		debug.ErrorMessage("CAMPAIGN_ID variable not found in environt")
		exit_err = true
	} else {
		err := verify_campaign()
		if err != nil {
			debug.Error(err)
			exit_err = true
		}
	}
	if exit_err {
		return
	}
	debug.Success("room informations found")
	debug.Message("running campaign on\n\t- Ruleset: " + Room.Rulest + "\n\t- Campaign: " + Room.Campaign)
}
