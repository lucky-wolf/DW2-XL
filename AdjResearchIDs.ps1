# grab all research files that we'll adj
param(
	[parameter(Mandatory = $true)][int]$startid,
	[parameter(Mandatory = $true)][int]$endid,
	[parameter(Mandatory = $true)][int]$offset
)

# note: this is global, not just this file...
# Set-PSDebug -Trace 1

$files = Get-ChildItem "XL/ResearchProjectDefinitions*.xml"
foreach ($Item in $files) {
	$Item = $Item.Name
	$target = "XL/$Item"
	$source = "temp/$Item"
	Copy $target $source
	perl "AdjustResearchId.pl" $source $startid $endid $offset $target
}
