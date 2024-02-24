<# pushes this project to the dw2/mods/xl folder #>
param(
	[parameter(Mandatory = $false)][string]$beta="false"
)

if ($beta -eq "true") {
	$name = "XL Beta"
} elseif ($beta -eq "false") {
	$name = "XL"
} else {
	Write-Host "you must specify -beta true | false"
	exit 1
}

$source = "C:\Users\steve\Projects\DW2-XL\XL"
$target = "C:\Steam\steamapps\common\Distant Worlds 2\mods\$name"
Write-Host "removing old files..."
Remove-Item -Force -Recurse -Path "$target\*"
Copy-Item -Force -Recurse -Container -Path "$source\*" -Destination "$target\" -PassThru | ForEach-Object { Write-Host $_.Name. }
if ($? -ne $True) {
	Write-Host "Copy Failure: exiting"
	exit 1
}

if ($beta -eq "true") {
	Move-Item -Force -Path "$target\beta.json" "$target\mod.json"
}

Write-Host "successfully updated $target"


exit 0