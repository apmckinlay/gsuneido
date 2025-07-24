// Copyright (C) 2004 Suneido Software Corp. All rights reserved worldwide.
function (current, globalExclude = #())
	{
	global_options = Object('Save...', 'Restore...','Reporter...', 'Summarize...',
		'CrossTable...', 'Export...').Difference(globalExclude)
	return Object('Horz'
		#(Button '&New', xmin: 80)
		#(Skip)
		Object('MenuButton', 'Current', current, xmin: 80)
		#(Skip)
		Object('MenuButton', 'Global', global_options, xmin: 80))
	}