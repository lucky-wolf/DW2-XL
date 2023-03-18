# input:	research-filename start-row end-row offset
# output:	stdout
# example:  ResearchProjectDefinitions.xml 10 999 2
#   all techs, starting from row 10, are adjusted downwards by 2 rows (e.g. row 10 -> row 12, 11 -> 13, etc...)
#   perl adjustrow.pl XL/ResearchProjectDefinitions.xml 10 9999 2 > "temp\ResearchProjectDefinitions.xml"

sub ProcessFile
{
	my($source) = $_[0];
	my($starting) = $_[1] + 0;
	my($ending) = $_[2] + 0;
	my($offset) = $_[3] + 0;
	my($target) = $_[4];

	print "source = $source\n";
	print "target = $target\n";
	print "rows = $starting..$ending\n";
	print "adjs = $offset\n";

	open(SF, '<', $source) or die "$!";
	open(OF, '>', $target) or die "$!";

	# read in each line of the file
	while ($line = <SF>)
	{
		# is a row property
		if ($line =~ /^([ \t]+)<Row>([0-9.]+)<\/Row>(.*)$/)
		{
			$row = $2 + 0;

			# and is within criteria
			if ($row >= $starting && $row <= $ending)
			{
				# adjust the row
				$line = sprint OFf("$1<Row>%d</Row>$3\n", $row + $offset);
			}
		}

		# copy this line to output
		print OF $line;
	}

	close(SF);
	close(OF);
}

ProcessFile($ARGV[0], $ARGV[1], $ARGV[2], $ARGV[3], $ARGV[4]);
