# input:	source-filename xml-key min max adj target-filename

sub ProcessFile
{
	my($source) = $_[0];
	my($key) = $_[1];
	my($min) = $_[2] + 0.0;
	my($max) = $_[3] + 0.0;
	my($adj) = $_[4] + 0.0;
	my($target) = $_[5];

	print "source = $source\n";
	print "target = $target\n";
	print "key = $key\n";
	print "min = $min\n";
	print "max = $max\n";
	print "adj = $adj\n";

	open(SF, '<', $source) or die "$!";
	open(OF, '>', $target) or die "$!";

	# read in each line of the file
	$ln = 0;
	while ($line = <SF>)
	{
		# keep track of 1 based line number
		++$ln;

		# matches xml key with floating point number (including sci notation)
		if ($line =~ /^([ \t]+)<([^>]+)>((0\.|[1-9]\.?)\d*((e|E)(\+|-)\d+)?)<\/([^>]+)>(.*)$/)
		{
			$value = $3 + 0.0;

			# and is within criteria
			if ($key eq $2 && $value >= $min && $value <= $max)
			{
				$new = $value + $adj;
				print "$ln: matched $2 for $value -> $new\n";
				# adjust the value
				$line = sprintf("$1<${key}>%g</${key}>\n", $new);
			}
		}

		# copy this line to output
		print OF $line;
	}

	close(SF);
	close(OF);
}

ProcessFile($ARGV[0], $ARGV[1], $ARGV[2], $ARGV[3], $ARGV[4], $ARGV[5]);
