<# pushes this project to the dw2/mods/xl folder #>
$source = "C:\Users\steve\Projects\DW2-XL\XL"
$target = "C:\Steam\steamapps\common\Distant Worlds 2\mods\XL"
Remove-Item -Force -Recurse -Path "$target\*"
Copy-Item -Force -Recurse -Container -Path "$source\*" -Destination "$target\" -PassThru | ForEach-Object { Write-Host $_.Name. }
