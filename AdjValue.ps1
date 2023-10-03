# require some parameters for this
param(
	[parameter(Mandatory = $true)][string]$filenamepattern,
	[parameter(Mandatory = $true)][string]$key,
	[parameter(Mandatory = $true)][int]$min,
	[parameter(Mandatory = $true)][int]$max,
	[parameter(Mandatory = $true)][int]$adj
)

# note: this is global, not just this file...
# Set-PSDebug -Trace 1

$files = Get-ChildItem "XL/${filenamepattern}.xml"
foreach ($Item in $files) {
	$Item = $Item.Name
	$target = "XL/$Item"
	$source = "temp/$Item"
	Copy-Item $target $source
	perl "AdjustNumericValue.pl" $source $key $min $max $adj $target
}
