#!/usr/bin/env bash
# Prépare l'environnement pour modifier ECE / Gestionnaire2Type en local via go.work.
# Les dépôts sont clonés à côté d'EngineWorker (../ECE, ../Gestionnaire2Type).

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
LOCAL_ECE_PATH=../ECE
LOCAL_GDT_PATH=../Gestionnaire2Type


clone_if_missing() {
    if [ -d "$1/.git" ]; then
        echo "[OK] $1 allready here"
    else
        echo "[..] Clonning $1"
        git clone "$2" "$1"
    fi
}

clone_if_missing $LOCAL_ECE_PATH    https://github.com/D-Gine/ECE.git
clone_if_missing $LOCAL_GDT_PATH    https://github.com/D-Gine/Gestionnaire2Type.git

cd "$ROOT"
[ -f go.work ] || go work init .
go work use . $LOCAL_ECE_PATH $LOCAL_GDT_PATH
go work sync

echo "[OK] go.work ready :"
echo "     Before pushing, verify the build as the CI/Docker : go build ./..."
