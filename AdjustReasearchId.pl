# input:	research-filename start-id end-id offset
# output:	stdout
# example:  perl temp\ResearchProjectDefinitions.xml 1789 1903 -789 XL\ResearchProjectDefinitions.xml
#   techs from 1789..1903 are adjusted downwards by 789 (1789 -> 1000)

sub ProcessFile
{
	my($source) = $_[0];
	my($starting) = $_[1] + 0;
	my($ending) = $_[2] + 0;
	my($offset) = $_[3] + 0;
	my($target) = $_[4];

	print "source = $source\n";
	print "target = $target\n";
	print "IDs = $starting..$ending\n";
	print "adj = $offset\n";

	open(SF, '<', $source) or die "$!";
	open(OF, '>', $target) or die "$!";

	# read in each line of the file
	while ($line = <SF>)
	{
		# is a ResearchProjectId property
		if ($line =~ /^([ \t]+)<ResearchProjectId>([0-9.]+)<\/ResearchProjectId>(.*)$/)
		{
			# convert to a number
			$ResearchProjectId = $2 + 0;

			# and is within criteria
			if ($ResearchProjectId >= $starting && $ResearchProjectId <= $ending)
			{
				# adjust the ResearchProjectId
				$line = sprintf("$1<ResearchProjectId>%d</ResearchProjectId>$3\n", $ResearchProjectId + $offset);
			}
		}

		# copy this line to output
		print OF $line;
	}

	close(SF);
	close(OF);
}

ProcessFile($ARGV[0], $ARGV[1], $ARGV[2], $ARGV[3], $ARGV[4]);
