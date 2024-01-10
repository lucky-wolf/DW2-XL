# require some parameters for this
param(
	[parameter(Mandatory = $true)][string]$filenamepattern,
	[parameter(Mandatory = $true)][string]$algorithm
	# [parameter(Mandatory = $false)][flag]$quiet,
)

# note: this is global, not just this file...
# Set-PSDebug -Trace 1

# ensure our code is compiled
Push-Location
Set-Location ./code
go build -o bin/transform.exe cmd/xmltree/main.go
$build = $LastExitCode
Pop-Location

if ($build -ne 0) {
	Write-Host "build failure: $build"
	exit $build
}

$files = Get-ChildItem "XL/${filenamepattern}.xml"
foreach ($Item in $files) {
	$Item = $Item.Name
	$target = "XL/$Item"
	$source = "temp/$Item"
	Copy-Item $target $source
	./code/bin/transform.exe -source $source -target $target -algorithm $algorithm -quiet
}
