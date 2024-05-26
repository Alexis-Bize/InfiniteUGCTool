#!/usr/bin/env bash

config_file="config.txt"

APP_NAME=$(sed -n 's/^APP_NAME="\([^"]*\)".*/\1/p' "$config_file")
APP_VERSION=$(sed -n 's/^APP_VERSION="\([^"]*\)".*/\1/p' "$config_file")

package_name="${APP_NAME}-${APP_VERSION}"
platforms=("windows/amd64" "darwin/amd64")

for platform in "${platforms[@]}"
do
	platform_split=(${platform//\// })
	GOOS=${platform_split[0]}
	GOARCH=${platform_split[1]}
	output_name=$package_name'-'$GOOS'-'$GOARCH
	if [ $GOOS = "windows" ]; then
		output_name+='.exe'
	fi

	env GOOS=$GOOS GOARCH=$GOARCH go build -o "./builds/"$output_name
	if [ $? -ne 0 ]; then
		echo 'An error has occurred! Aborting the script execution...'
		exit 1
	fi
done