# grab all research files that we'll adj
param(
	[parameter(Mandatory = $true)][int]$rowstart,
	[parameter(Mandatory = $true)][int]$rowend,
	[parameter(Mandatory = $true)][int]$colstart,
	[parameter(Mandatory = $true)][int]$colend,
	[parameter(Mandatory = $true)][int]$offset
)

# note: this is global, not just this file...
# Set-PSDebug -Trace 1

$files = Get-ChildItem "XL/ResearchProjectDefinitions*.xml"
foreach ($Item in $files) {
	$Item = $Item.Name
	$target = "../XL/$Item"
	$source = "../temp/$Item"
	Copy-Item $target $source
	# input:	source-filename rowstart rowend colsstart coleend coffset output-filename
	perl "scripts/AdjustCol.pl" $source $rowstart $rowend $colstart $colend $offset $target
}
