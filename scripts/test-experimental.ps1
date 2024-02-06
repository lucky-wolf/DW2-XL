<# pushes this project to the dw2/mods/xl folder and then runs dw2 #>
scripts\push-to-mods.ps1
if ($? -ne $True) {
	exit $status
}
# "--new-game"
Start-Process -FilePath "C:\Users\steve\Downloads\DW2 Unstable\DistantWorlds2.exe" "--new-game" -WorkingDirectory "C:\Users\steve\Downloads\DW2 Unstable"
Write-Host "launching experimental game..."
