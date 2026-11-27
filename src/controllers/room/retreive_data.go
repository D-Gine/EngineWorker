package room

import (
	"EngineWorker/src/database"
	"EngineWorker/src/debug"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	ece "github.com/D-Gine/ECE"

	"github.com/gin-gonic/gin"
)

type StaticComponent struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Type_id string `json:"type_id"`
	Value   any    `json:"value"`
}

type DynamicComponent struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Type_id string `json:"type_id"`
	Value   any    `json:"value"`
}

type Entity struct {
	Id          string             `json:"id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Sc          []StaticComponent  `json:"static_components"`
	Dc          []DynamicComponent `json:"dynamic_components"`
}

func retreive_entities_rows() (*sql.Rows, error) {
	rows, err := database.DB.Query(`
		SELECT
			e.id, e.name, e.description,
			COALESCE(sc.scs, '[]'),
			COALESCE(dc.dcs, '[]')
		FROM entities as e
			LEFT JOIN LATERAL (
				SELECT json_agg(json_build_object(
			   		'id', sc.id,
					'name', sct.name,
			   		'type_id', sct.type_id,
			   		'value', sc.value
			   		)) AS scs
			   	FROM static_components as sc
			   		LEFT JOIN static_components_templates as sct
			   			on sc.template_id = sct.id
				WHERE sc.entity_id = e.id
			) sc ON true
			LEFT JOIN LATERAL (
				SELECT json_agg(json_build_object(
			   		'id', dc.id,
					'name', dct.name,
			   		'type_id', dct.type_id,
			   		'value', dc.value
			   		)) AS dcs
			   	FROM dynamic_components as dc
			   		LEFT JOIN dynamic_components_templates as dct
			   			on dc.template_id = dct.id
				WHERE dc.entity_id = e.id
			) dc ON true
		WHERE
			e.ruleset_id=$1
	`, database.Room.Rulest)
	return rows, err
}

func retreive_entities(rows *sql.Rows) ([]Entity, error) {
	defer rows.Close()
	var entities []Entity

	for rows.Next() {
		var e Entity
		var scRaw, dcRaw []byte
		if err := rows.Scan(&e.Id, &e.Name, &e.Description, &scRaw, &dcRaw); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(scRaw, &e.Sc); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(dcRaw, &e.Dc); err != nil {
			return nil, err
		}
		entities = append(entities, e)
	}
	return entities, nil
}

func create_new_registry(entities []Entity) (*ece.Registry, error) {
	reg := ece.NewRegistry()

	for idx, e := range entities {
		for _, sc := range e.Sc {
			reg.RegisterComponentByName[any](sc.Type_id)
			err := reg.AddComponentByName[any](idx, sc.Type_id, sc.Value)
			if err != nil {
				debug.ErrorMessage(fmt.Sprintf("could not add component '%s' to entity '%s'(number %d). %s", sc.Id, e.Id, idx, err.Error()))
				return nil, err
			} else {
				debug.Message(fmt.Sprintf("Added component %s(%s) to %s(%d - %s)", sc.Name, sc.Type_id, e.Name, idx, e.Id))
			}
		}
		for _, dc := range e.Dc {
			reg.RegisterComponentByName[any](dc.Type_id)
			err := reg.AddComponentByName(idx, dc.Type_id, dc.Value)
			if err != nil {
				debug.ErrorMessage(fmt.Sprintf("could not add component '%s' to entity '%s'(number %d). %s", dc.Id, e.Id, idx, err.Error()))
				return nil, err
			} else {
				debug.Message(fmt.Sprintf("Added component %s(%s) to %s(%d - %s)", dc.Name, dc.Type_id, e.Name, idx, e.Id))
			}
		}
	}
	debug.Ok("done adding every components to registry")

	return reg, nil
}

func RetreiveData(c *gin.Context) {
	if database.Room.Errored == true {
		c.JSON(http.StatusBadRequest, gin.H{"error": "rooms id are either missing or does not exist"})
		return
	}

	debug.Message("Retreiving row entities from database...")
	rows, err := retreive_entities_rows()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "retreiving raw entities from database... " + err.Error()})
		return
	}
	debug.Ok("row entities retreived")

	debug.Message("Entities data verifications...")
	entities, err := retreive_entities(rows)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "reading raw entities... " + err.Error()})
		return
	}
	debug.Ok("entities info and components verified")

	debug.Message("Creating the new registry")
	new_reg, err := create_new_registry(entities)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "inserting data into registry... " + err.Error()})
		return
	}
	debug.Success("new registry created")

	c.Set("registry", new_reg)
	// TODO change return status ?
	c.JSON(http.StatusAccepted, nil)
}
