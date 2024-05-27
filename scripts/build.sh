#!/usr/bin/env bash

script_path=$(dirname "$0")
script_path=$(cd "$script_path" && pwd)
root_path="${script_path}/.."
config_file="${root_path}/configs/application.yaml"

name=$(grep "^name:" "$config_file" | cut -d ":" -f 2- | sed 's/^ *//g')
description=$(grep "^description:" "$config_file" | cut -d ":" -f 2- | sed 's/^ *//g')
version=$(grep "^version:" "$config_file" | cut -d ":" -f 2- | sed 's/^ *//g')
author=$(grep "^author:" "$config_file" | cut -d ":" -f 2- | sed 's/^ *//g')
repository=$(grep "^repository:" "$config_file" | cut -d ":" -f 2- | sed 's/^ *//g')

package_name="${name}-${version}"
target_platforms=("windows/amd64" "darwin/amd64")

IFS='.' read -r major minor patch <<< "$version"
build=0

versioninfo='{
	"FixedFileInfo": {
		"FileVersion": {
			"Major": '"$major"',
			"Minor": '"$minor"',
			"Patch": '"$patch"',
			"Build": '"$build"'
		},
		"ProductVersion": {
			"Major": '"$major"',
			"Minor": '"$minor"',
			"Patch": '"$patch"',
			"Build": '"$build"'
		},
		"FileFlagsMask": "3f",
		"FileFlags ": "00",
		"FileOS": "040004",
		"FileType": "01",
		"FileSubType": "00"
	},
	"StringFileInfo": {
		"Comments": "",
		"CompanyName": "'"$author"'",
		"FileDescription": "'"$description"'",
		"FileVersion": "v'"$version.$build"'",
		"InternalName": "'"$name".exe'",
		"LegalCopyright": "Copyright (c) '$(date +"%Y")' '"$author"'",
		"LegalTrademarks": "",
		"OriginalFilename": "main.go",
		"PrivateBuild": "",
		"ProductName": "'"$name"'",
		"ProductVersion": "v'"$version.$build"'",
		"SpecialBuild": ""
	},
	"VarFileInfo": {
		"Translation": {
			"LangID": "0409",
			"CharsetID": "04B0"
		}
	},
	"IconPath": "",
	"ManifestPath": ""
}'

echo "$versioninfo" >| "${root_path}/versioninfo.json"
go generate

for platform in "${target_platforms[@]}"
do
	platform_split=(${platform//\// })

	GOOS=${platform_split[0]}
	GOARCH=${platform_split[1]}

	base_output_name=$package_name'-'$GOOS'-'$GOARCH
	output_name=$base_output_name
	if [ $GOOS = "windows" ]; then
		output_name+='.exe'
	fi

	env GOOS=$GOOS GOARCH=$GOARCH go build -o "./build/bin/"$output_name
	if [ $? -ne 0 ]; then
		echo 'An error has occurred! Aborting the script execution...'
		exit 1
	fi

	mkdir -p "./build/archive/"
	zip -r "./build/archive/"$base_output_name".zip" "./build/bin/"$output_name
done
