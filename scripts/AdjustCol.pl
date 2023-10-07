# input:	source-filename rowstart rowend colsstart coleend coffset output-filename
# output:	stdout
# example:  temp/ResearchProjectDefinitions.xml 10 9999 2 8 +1 XL/ResearchProjectDefinitions.xml
#   all techs in rows 10..9999 and columns 2..8 are adjusted by +1 columns
#   perl adjustcol.pl temp/ResearchProjectDefinitions.xml 10 9999 2 8 +1 XL/ResearchProjectDefinitions.xml

sub ProcessFile
{
	my($source) = $_[0];
	my($rfirst) = $_[1] + 0;
	my($rlast) = $_[2] + 0;
	my($cfirst) = $_[3] + 0;
	my($clast) = $_[4] + 0;
	my($offset) = $_[5] + 0;
	my($target) = $_[6];

	print "source = $source\n";
	print "target = $target\n";
	print "rows = $rfirst..$rlast\n";
	print "cols = $cfirst..$clast\n";
	print "adjs = $offset\n";

	open(SF, '<', $source) or die "$!";
	open(OF, '>', $target) or die "$!";

	# die;

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

			print OF $lncol;
			print OF $lnrow;
		}
		else
		{
			print OF $lncol;
		}
	}

	close(SF);
	close(OF);
}

ProcessFile($ARGV[0], $ARGV[1], $ARGV[2], $ARGV[3], $ARGV[4], $ARGV[5], $ARGV[6]);
