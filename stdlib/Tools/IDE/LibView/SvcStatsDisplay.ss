// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
class
	{
	Title: 'SvcStatistics'

	CallClass(name, lib)
		{
		if name is ""
			{
			Beep()
			return
			}
		Alert(.getRecordStats(name, lib), 'Main Contributors')
		}
	getRecordStats(name, lib)
		{
		if false is settings = SvcSettings()
			return "Svc Unavailable"
		svc = Svc(server: settings.svc_server, local?: settings.svc_local?)
		stats = SvcStatistics(svc, lib, name)
		numContributors = stats.GetNumOfContributors()
		contributions = stats.WeighContributions()[::3/*=result limit*/]
		return "Number of Contributors: " $ Display(numContributors) $ '\n\n' $
			contributions.Map({ it[0] $ ': ' $ it[1].Round(0) }).Join('\n')
		}
	}
