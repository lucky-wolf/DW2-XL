# grab all research files that we'll adj
param(
	[parameter(Mandatory = $true)][int]$startid,
	[parameter(Mandatory = $true)][int]$endid,
	[parameter(Mandatory = $true)][int]$offset
)

# note: this is global, not just this file...
# Set-PSDebug -Trace 1

# update the components themselves
$files = Get-ChildItem "XL/ComponentDefinitions*.xml"
foreach ($Item in $files) {
	$Item = $Item.Name
	$target = "XL/$Item"
	$source = "temp/$Item"
	Copy $target $source
	perl "AdjustComponentId.pl" $source $startid $endid $offset $target
}

# update research references to those components
# hack: because I'm unskilled with powershell, we'll do the ugly but workable thing
# todo: really need a way to combine the list of files up-front
# todo: or a way to call a sub-funciton to do the work so we're not repeating ourselves!
$files = Get-ChildItem "XL/ResearchProjectDefinitions*.xml"
foreach ($Item in $files) {
	$Item = $Item.Name
	$target = "XL/$Item"
	$source = "temp/$Item"
	Copy $target $source
	perl "AdjustComponentId.pl" $source $startid $endid $offset $target
}
