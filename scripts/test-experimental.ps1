<# pushes this project to the dw2/mods/xl folder and then runs dw2 #>
scripts\push-to-mods.ps1
# "--new-game"
Start-Process -FilePath "C:\Users\steve\Downloads\DW2 Unstable\DistantWorlds2.exe" -WorkingDirectory "C:\Users\steve\Downloads\DW2 Unstable"
Write-Host "launching experimental game..."
