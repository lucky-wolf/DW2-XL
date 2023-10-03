# grab all research files that we'll adj
param(
	[parameter(Mandatory = $true)][int]$start,
	[parameter(Mandatory = $true)][int]$end,
	[parameter(Mandatory = $true)][int]$offset
)

# note: this is global, not just this file...
# Set-PSDebug -Trace 1

$files = Get-ChildItem "XL/ResearchProjectDefinitions*.xml"
foreach ($Item in $files) {
	$Item = $Item.Name
	$target = "XL/$Item"
	$source = "temp/$Item"
	Copy-Item $target $source
	perl "AdjustRow.pl" $source $start $end $offset $target
}
