$APPNAME = "go-svc"

if (-not (Test-Path ".\bin")) {
    New-Item -ItemType Directory -Path ".\bin"
}

if (Test-Path ".\bin\${APPNAME}.exe") {
    Remove-Item ".\bin\${APPNAME}.exe"
}

# rsrc -manifest "./manifest.xml" -o "./rsrc.syso"