# require some parameters for this
param(
	[parameter(Mandatory = $true)][string]$filenamepattern,
	[parameter(Mandatory = $true)][string]$key,
	[parameter(Mandatory = $true)][float]$min,
	[parameter(Mandatory = $true)][float]$max,
	[parameter(Mandatory = $true)][float]$scale
)

# note: this is global, not just this file...
# Set-PSDebug -Trace 1

$files = Get-ChildItem "XL/${filenamepattern}.xml"
foreach ($Item in $files) {
	$Item = $Item.Name
	$target = "XL/$Item"
	$source = "temp/$Item"
	Copy-Item $target $source
	perl "ScaleNumericValue.pl" $source $key $min $max $scale $target
}
