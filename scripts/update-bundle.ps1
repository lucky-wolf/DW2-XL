<# pushes this project to the dw2/mods/xl folder #>
$source = "C:\Users\steve\Projects\XL\Bin\Windows\Debug\win-x64\data\db\bundles"
$target = "C:\Users\steve\Projects\DW2-XL\XL"
Write-Host "removing old bundle files..."
Remove-Item -Force -Path "$target\*.bundle"
Write-Host "copying new bundle files..."
Copy-Item -Container -Force -Path "$source\XL*.bundle" -Destination "$target\" -PassThru | ForEach-Object { Write-Host $_.Name. }
Write-Host "bundle updated successfully"
