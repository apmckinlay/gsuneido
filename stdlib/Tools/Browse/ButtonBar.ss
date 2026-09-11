// Copyright (C) 2004 Suneido Software Corp. All rights reserved worldwide.
function(current, globalExclude = #())
	{
	global_options = ["Save...", "Restore...", "Reporter...", "Reporter Forms...",
			"Summarize...", "CrossTable...", "Export..."].
		Difference(globalExclude)
	return [#Horz,
		#(Button, "&New", xmin: 80),
		#(Skip),
		[#MenuButton, #Current, current, xmin: 80],
		#(Skip),
		[#MenuButton, #Global, global_options, xmin: 80]]
	}
