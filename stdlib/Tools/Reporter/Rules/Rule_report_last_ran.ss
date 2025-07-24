// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
function ()
	{
	last_ran_by = last_ran_on = false
	paramsName = 'Reporter Report - ' $ .report_name
	reportName = 'Reporter - ' $ .report_name
	recs = QueryAll('params extend lastRan = params["lastRan"]
		where (report is ' $ Display(reportName) $ ' or
			report is ' $ Display(paramsName) $ ') and lastRan > ""
		summarize user, max lastRan sort reverse max_lastRan')
	if not recs.Empty?()
		{
		last_ran_by = recs[0].user isnt ''
			? recs[0].user
			: 'Scheduled Report'
		last_ran_on = recs[0].max_lastRan
		}

	return [:last_ran_by, :last_ran_on]
	}
