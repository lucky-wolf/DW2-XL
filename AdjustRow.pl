# input:	research-filename start-row end-row adjustment
# output:	stdout
# example:  ResearchProjectDefinitions.xml 10 999 2
#   all techs, starting from row 10, are adjusted downwards by 2 rows (e.g. row 10 -> row 12, 11 -> 13, etc...)
#   perl adjustrow.pl XL/ResearchProjectDefinitions.xml 10 9999 2 > "temp\ResearchProjectDefinitions.xml"

sub ProcessFile
{
	my($sourcefile) = $_[0];
	my($starting) = $_[1] + 0;
	my($ending) = $_[2] + 0;
	my($adjustment) = $_[3] + 0;

	# print "<!-- AdjustRow: starting=$starting, ending=$ending, adjustment=$adjustment -->\n";

	open(SF, $sourcefile) or die "$!";

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
				$line = sprintf("$1<Row>%d</Row>$3\n", $row + $adjustment);
			}
		}

		# copy this line to output
		print $line;
	}

	close(SF);
}

ProcessFile($ARGV[0], $ARGV[1], $ARGV[2], $ARGV[3]);
