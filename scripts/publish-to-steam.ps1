<# pushes this project to the dw2/mods/xl folder and then publishes the mod #>
param(
	[parameter(Mandatory = $false)][string]$beta
)

if ($beta -eq "true") {
	$name = "XL Beta"
} elseif ($beta -eq "false") {
	$name = "XL"
} else {
	Write-Host "you must specify -beta true | false"
	exit 1
}

scripts\push-to-mods.ps1 -beta $beta
if ($? -ne $True) {
	exit $status
}

Push-Location
Set-Location "C:\Steam\steamapps\common\Distant Worlds 2"
./distantworlds2 --ugc-publish "mods/$name"
Pop-Location
