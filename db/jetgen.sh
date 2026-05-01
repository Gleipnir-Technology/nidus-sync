#! /usr/bin/env nix-shell
#! nix-shell -i bash -p bash
jet -dsn="postgresql://?host=/var/run/postgresql&sslmode=disable&dbname=nidus-sync" -schema=stadia -schema=arcgis -path=./gen
