<# erases source bundle, builds them again, erases project bundle, and copies newly created ones here #>
$project = "C:\Users\steve\Projects\XL"
$source = "$project\Bin\Windows\Debug\win-x64\data\db\bundles"
$target = "C:\Users\steve\Projects\DW2-XL\XL"

Write-Host "removing old source bundle files..."
Remove-Item -Force -Path "$source\*.bundle"

Write-Host "building source bundle files..."
Push-Location
Set-Location $project
dotnet.exe build "XL.sln"
$build = $LastExitCode
Pop-Location

if ($build -ne 0) {
	Write-Host "build failure: $build"
	exit $build
}

Write-Host "removing old bundle files..."
Remove-Item -Force -Path "$target\*.bundle"

Write-Host "copying new bundle files..."
Copy-Item -Container -Force -Path "$source\XL*.bundle" -Destination "$target\" -PassThru | ForEach-Object { Write-Host $_.Name. }

Write-Host "bundle updated successfully"
