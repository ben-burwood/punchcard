set shell := ["powershell.exe", "-c"]

default:
    just --list

format:
    gofmt -w .
    Set-Location frontend; npx prettier . --write

backend *args:
    go run . {{ args }}

frontend *args:
    Set-Location frontend; npm run {{ args }}
