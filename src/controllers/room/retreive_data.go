package room

import (
	"EngineWorker/src/database"
	"database/sql"
	"net/http"

	ece "github.com/D-Gine/ECE"
	"github.com/gin-gonic/gin"
)

func retreive_entities() (*sql.Rows, error) {
	rows, err := database.DB.Query(`
		SELECT
			e.id, e.name, e.description,
			COALESCE(sc.scs, '[]'),
			COALESCE(dc.dcs, '[]')
		FROM entities as e
			LEFT JOIN LATERAL (
				SELECT json_agg(json_build_object(
			   		'id', sc.id,
			   		'template_id', sc.template_id,
			   		'value', sc.value
			   		)) AS scs
			   	FROM static_components as sc
				WHERE sc.entity_id = e.id
			) sc ON true
			LEFT JOIN LATERAL (
				SELECT json_agg(json_build_object(
			   		'id', dc.id,
			   		'template_id', dc.template_id,
			   		'value', dc.value
			   		)) AS dcs
			   	FROM dynamic_components as dc
				WHERE dc.entity_id = e.id
			) dc ON true
		WHERE
			e.ruleset_id=$1
	`, database.Room.Rulest)
	return rows, err
}

func create_new_registry(*sql.Rows) (*ece.Registry, error) {
	reg := ece.NewRegistry()
	return reg, nil
}

func RetreiveData(c *gin.Context) {
	if database.Room.Errored == true {
		c.JSON(http.StatusBadRequest, gin.H{"error": "rooms id are either missing or does not exist"})
		return
	}

	rows, err := retreive_entities()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "retreiving entities from database, " + err.Error()})
		return
	}

	new_reg, err := create_new_registry(rows)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "inserting data into registry, " + err.Error()})
	}

	reg := c.MustGet("registry").(*ece.Registry)
	reg = new_reg

	// TODO change return status ?
	c.JSON(http.StatusAccepted, nil)
}
