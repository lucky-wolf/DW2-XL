# input:	research-filename start-id end-id adjustment
# output:	stdout
# example:  ResearchProjectDefinitions.xml 1789 1903 -789
#   techs from 1789..1903 are adjusted downwards by 789 (1789 -> 1000)
#   perl .\AdjustReasearchId.pl XL/ResearchProjectDefinitions.xml 1789 1903 -789 > "temp\ResearchProjectDefinitions.xml.txt"

sub ProcessFile
{
	my($sourcefile) = $_[0];
	my($starting) = $_[1] + 0;
	my($ending) = $_[2] + 0;
	my($adjustment) = $_[3] + 0;

	# print "<!-- AdjustResearchId: starting=$starting, ending=$ending, adjustment=$adjustment -->\n";

	open(SF, $sourcefile) or die "$!";

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
				$line = sprintf("$1<ResearchProjectId>%d</ResearchProjectId>$3\n", $ResearchProjectId + $adjustment);
			}
		}

		# copy this line to output
		print $line;
	}

	close(SF);
}

ProcessFile($ARGV[0], $ARGV[1], $ARGV[2], $ARGV[3]);
