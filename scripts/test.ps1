<# pushes this project to the dw2/mods/xl folder and then runs dw2 #>
scripts\push-to-mods.ps1
Start-Process -FilePath "C:\Steam\steamapps\common\Distant Worlds 2\DistantWorlds2.exe" "--skip-splash --new-game" -WorkingDirectory "C:\Steam\steamapps\common\Distant Worlds 2"