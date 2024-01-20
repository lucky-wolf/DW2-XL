# require some parameters for this
param(
	[parameter(Mandatory = $true)][string]$algorithm,
	# [parameter(Mandatory = $false)][flag]$quiet,
	[parameter(Mandatory = $false)][double]$scale
)

# note: this is global, not just this file...
# Set-PSDebug -Trace 1

# scale algorithms require a scale factor
if ($algorithm -eq "ScalePlanetFrequencies" -And $scale -eq 0) {
	Write-Host "ScalePlanetFrequencies requires a scale factor: $scale"
	exit 1
}

# ensure our code is compiled
Push-Location
Set-Location ./code
go build -o bin/xmltree.exe cmd/xmltree/main.go
$build = $LastExitCode
Pop-Location

if ($build -ne 0) {
	Write-Host "build failure: $build"
	exit $build
}

if ($algorithm -eq "ScalePlanetFrequencies") {
	./code/bin/xmltree.exe -algorithm $algorithm -quiet -scale $scale
} else {
	./code/bin/xmltree.exe -algorithm $algorithm -quiet
}
