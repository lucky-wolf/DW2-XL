# input:	source-filename xml-key min max factor target-filename

sub ProcessFile
{
	my($source) = $_[0];
	my($key) = $_[1];
	my($min) = $_[2] + 0.0;
	my($max) = $_[3] + 0.0;
	my($factor) = $_[4] + 0.0;
	my($target) = $_[5];

	print "source = $source\n";
	print "target = $target\n";
	print "key = $key\n";
	print "min = $min\n";
	print "max = $max\n";
	print "factor = $factor\n";

	if ($factor == 0.0 && $min <= 0.0 && $max >= 0.0) {
		print "error: cannot scale zero by zero";
		return
	}

	open(SF, '<', $source) or die "$!";
	open(OF, '>', $target) or die "$!";

	# read in each line of the file
	$ln = 0;
	while ($line = <SF>)
	{
		# keep track of 1 based line number
		++$ln;

		# matches xml key with floating point number (including sci notation)
		if ($line =~ /^([ \t]+)<([^>]+)>((0\.|[1-9]\.?)\d*(e(\+|-)\d+)?)<\/([^>]+)>(.*)$/)
		{
			# convert to numeric value
			# a non-number's value will presumably be 0, and we explicitly skip zeros, so we're safe enough here (I think)
			$value = $3 + 0.0;

			# and is within criteria
			if ($key eq $2)
			{
				if ($value == 0.0)
				{
					# skip it
				}
				elsif ($value < $min)
				{
					# print "$ln: X $value < $min\n";
				}
				elsif ($value > $max)
				{
					# print "$ln: X $value > $max\n";
				}
				else
				{
					$new = $value * $factor;
					print "$ln: matched $2 for $value -> $new\n";
					# adjust the value
					$line = sprintf("$1<${key}>%g</${key}>$5\n", $new);
				}
			}
		}

		# copy this line to output
		print OF $line;
	}

	close(SF);
	close(OF);
}

ProcessFile($ARGV[0], $ARGV[1], $ARGV[2], $ARGV[3], $ARGV[4], $ARGV[5]);
