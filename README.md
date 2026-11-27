# EngineWorker

## Build and Run with
### Env
For this project you need the following environment variables:
```
# ROOM INFORMATIONS
RULESET_ID="aaaaa"
CAMPAIGN_ID="aaaaa"

# DATABASE
DB_HOST=dengine-db-{x.x.x}
DB_PORT=5432
DB_USER=[DB_USER : string]
DB_PWD=[DB_PWD : string]
DB_NAME=dengine

# WORKER-API
VERSION=1.0.0
API_PORT=[API_PORT : int]
GIN_MODE=debug

# NETWORKS
BACK_NETWORK=dengine-back-net
```

## Dev with local version of ECE / Gestionnaire2Type
Packages `github.com/D-Gine/ECE` and `github.com/D-Gine/Gestionnaire2Type` are needed for the project, and their version are fixed inside `go.mod`.
For unpushed modifications for testing purpose, we use a Go workspace (`go.work`, local et non-commited) :

```bash
make dev-setup   # clone ../ECE et ../Gestionnaire2Type if needed, then generate the go.work
```

Note that you can modify the localisation of the ECE and Gestionnaire2Type if you want easier switch between project. You just have to modify the paths in the `./scripts/setup_local_dev.sh` line 6-7.
```bash
...
LOCAL_ECE_PATH=../ECE
LOCAL_GDT_PATH=../Gestionnaire2Type
...
```

Workflow :
1. In `../ECE` or `../Gestionnaire2Type`, create a branch and do your modification. EngineWorker will use them directly (`make run`).
2. Push then merge the branch on ECE.
3. In EngineWorker, update the used version : `go get github.com/D-Gine/ECE@<tag|branche>` or `go get github.com/D-Gine/Gestionnaire2Type@<tag|branche>` , then commit `go.mod` and `go.sum`.
4. **BEFORE PUSHING**, verify that the project work without like in the Docker/CI (without local project) : `go build ./...`
