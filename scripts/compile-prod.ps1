$scriptDir = Split-Path -Path $MyInvocation.MyCommand.Definition -Parent
. "${scriptDir}\__init.ps1"

# go build -ldflags "-s -w -X 'main.manifest=rsrc.syso'" -o "./bin/${APPNAME}.exe" ".\example\service\"
go build -ldflags "-s -w" -o "./bin/${APPNAME}.exe" ".\example\service\"