<# pushes this project to the dw2/mods/xl folder and then runs dw2 #>
scripts\push-to-mods.ps1
if ($? -ne $True) {
	exit $status
}
# Start-Process -FilePath "C:\Steam\steamapps\common\Distant Worlds 2\DistantWorlds2.exe" "--help" -WorkingDirectory "C:\Steam\steamapps\common\Distant Worlds 2"
#  "--new-game"
#  "--gen-xsd"
Start-Process -FilePath "C:\Steam\steamapps\common\Distant Worlds 2\DistantWorlds2.exe" -WorkingDirectory "C:\Steam\steamapps\common\Distant Worlds 2"
Write-Host "launching steam game..."
