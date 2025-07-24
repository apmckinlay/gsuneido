// Copyright (C) 2013 Axon Development Corporation All rights reserved worldwide.
function (tasks)
	{
	Plugins().ForeachContribution('ScheduledTasks', tasks)
		{ |task|
		try
			task.func.Eval() // needs Eval
		catch (err)
			{
			pre = "ERROR: Run Scheduled Task failed (" $ task.func $ ") - "
			SuneidoLog(pre $ err)
			}
		}
	}