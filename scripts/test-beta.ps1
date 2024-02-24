<# pushes this project to the dw2/mods/xl folder and then runs dw2 #>

# run "all"
scripts\xmltree.ps1 All
if ($? -ne $True) {
	exit $status
}

# push to mod folder
scripts\push-to-mods.ps1 -beta true
if ($? -ne $True) {
	exit $status
}

# run the game
# "--new-game"
Start-Process -FilePath "C:\Users\steve\Downloads\DW2 Unstable\DistantWorlds2.exe" -WorkingDirectory "C:\Users\steve\Downloads\DW2 Unstable"
Write-Host "launching experimental game..."
