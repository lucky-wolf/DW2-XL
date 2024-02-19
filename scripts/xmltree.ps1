# require some parameters for this
param(
	[parameter(Mandatory = $true)][string]$algorithm,
	# [parameter(Mandatory = $false)][flag]$quiet,
	[parameter(Mandatory = $false)][double]$scale,
	[parameter(Mandatory = $false)][string]$filter
)

# note: this is global, not just this file...
# Set-PSDebug -Trace 1

# scale algorithms require a scale factor
if ($algorithm -eq "ScalePlanetFrequencies" -And $scale -eq 0) {
	Write-Host "ScalePlanetFrequencies requires a scale factor: $scale"
	exit 1
}

# ensure our code is compiled
scripts\make.ps1
if ($? -ne $True) {
	exit $status
}

if ($algorithm -eq "ScalePlanetFrequencies") {
	bin/xmltree.exe -algorithm $algorithm -quiet -scale $scale -filter $filter
} else {
	bin/xmltree.exe -algorithm $algorithm -quiet
}
