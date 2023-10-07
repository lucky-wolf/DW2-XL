<# Makes a backup of xml and txt files given the version number you're backing up#>

# must supply a version string (e.g. 1.1.2.4)
param(
	[ValidateLength(4, 10)]
	[parameter(Mandatory = $true)]
	[String]
	$version
)

$target = "data - " + $version
$policy = $target + "\policy"

Push-Location
Set-Location "C:\Steam\steamapps\common\Distant Worlds 2"
New-Item -ItemType Directory -Force -Path $target
New-Item -ItemType Directory -Force -Path $policy
Copy-Item -Path "data\*" -Filter *.xml -Container -Destination $target
Copy-Item -Path "data\policy\*" -Filter *.xml -Container -Destination $policy
Pop-Location

$final = "C:\Steam\steamapps\common\Distant Worlds 2\" + $target
Get-ChildItem -Path $final -Recurse
