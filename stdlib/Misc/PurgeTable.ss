// Copyright (C) 2015 Suneido Software Corp. All rights reserved worldwide.
function (table, dateField, date = false, monthsToKeep = 3, outputToFile? = false)
	{
	if outputToFile?
		{
		logPath = 'logs/purgeTable/' $ table $ '/'
		EnsureDirectories('./' $ logPath)
		Save10(logPath $ table $ ".su")
		Database.Dump(table)
		MoveFile(table $ ".su", logPath $ table $ ".su.1")
		}
	cutoff = date isnt false ? date : Date().Plus(months: -monthsToKeep)
	// can't use QueryDo to delete
	// because it may fail from too many writes in one transaction
	QueryApplyMulti(table $ " where " $ dateField $ " < " $ Display(cutoff), update:)
		{ |x|
		x.Delete()
		}
	}
