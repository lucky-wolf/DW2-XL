<# pushes this project to the dw2/mods/xl folder and then publishes the mod #>
scripts\push-to-mods.ps1
Push-Location
Set-Location "C:\Steam\steamapps\common\Distant Worlds 2"
./distantworlds2 --ugc-publish mods/XL
Pop-Location