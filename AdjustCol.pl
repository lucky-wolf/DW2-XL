# input:	research-filename start-column end-column offset
# output:	stdout
# example:  ResearchProjectDefinitions.xml 10 9999 2 99 3
#   all techs in rows 10..9999 and columns 2..99 are adjusted by +3 columns
#   perl adjustcol.pl data/ResearchProjectDefinitions.xml 10 9999 2 99 3 > "Test.xml"

sub ProcessFile
{
	my($source) = $_[0];
	my($rfirst) = $_[1] + 0;
	my($rlast) = $_[2] + 0;
	my($cfirst) = $_[3] + 0;
	my($clast) = $_[4] + 0;
	my($offset) = $_[5] + 0;

	# print STDERR "rows $rfirst..$rlast columns $cfirst..$clast by $offset\n";
	# return;

	open(SF, $source) or die "$!";

	# read in each line of the file
	while ($lncol = <SF>)
	{
		# it must be a column spec
		if ($lncol =~ /^([ \t]+)<Column>([0-9.]+)<\/Column>(.*)$/)
		{
			$pre = $1;
			$column = $2 + 0;
			$post = $3;

			# must be a next line
			$lnrow = <SF>;

			# see if this is a row spec
			if ($lnrow =~ /^([ \t]+)<Row>([0-9.]+)<\/Row>(.*)$/)
			{
				# extract the row
				$row = $2 + 0;

				# check if within criteria
				if ($cfirst <= $column && $column <= $clast && $rfirst <= $row  && $row <= $rlast)
				{
					# adjust the column
					$lncol = sprintf("$pre<Column>%d</Column>$post\n", $column + $offset);
				}
			}

			print $lncol;
			print $lnrow;
		}
		else
		{
			print $lncol;
		}
	}

	close(SF);
}

ProcessFile($ARGV[0], $ARGV[1], $ARGV[2], $ARGV[3], $ARGV[4], $ARGV[5]);
