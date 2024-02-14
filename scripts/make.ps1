# compile our cli tool
Push-Location
Set-Location ./code
go build -o ../bin/xmltree.exe cmd/xmltree/main.go
$build = $LastExitCode
Pop-Location

if ($build -ne 0) {
	Write-Host "build failure: $build"
	exit $build
}

Write-Host "build succeeded: bin/xmltree.exe"
