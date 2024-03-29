# input:	source-filename xml-key min max factor target-filename

sub ProcessFile
{
	my($source) = $_[0];
	my($key) = $_[1];
	my($min) = $_[2] + 0.0;
	my($max) = $_[3] + 0.0;
	my($factor) = $_[4] + 0.0;
	my($target) = $_[5];

	# print "source = $source\n";
	# print "target = $target\n";
	# print "key = $key\n";
	# print "min = $min\n";
	# print "max = $max\n";
	# print "factor = $factor\n";

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
		# doesn't handle leading decimal point, nor trailing one - but requires a leading and trailing digit(s)
		if ($line =~ /^([ \t]+)<([^>]+)>(((\d+\.\d+)|(\d+))([eE][+-]\d+)?)<\/([^>]+)>(.*)$/)
		{
			# convert to numeric value
			$value = $3 + 0.0;

			# and is within criteria
			# subtle: ignore zero because that will always result in zero
			if ($key eq $2 && $value != 0.0 && $value >= $min && $value <= $max)
			{
				# compute new value
				$new = $value * $factor;

				# this format lets VSCode give you a jump to file & line number option!
				printf("$target:$ln: $value -> $new (%.5g)\n", $new);

				# adjust the value
				$line = sprintf("$1<${key}>%.5g</${key}>\n", $new);
			}
		}

		# copy this line to output
		print OF $line;
	}

	close(SF);
	close(OF);
}

ProcessFile($ARGV[0], $ARGV[1], $ARGV[2], $ARGV[3], $ARGV[4], $ARGV[5]);
