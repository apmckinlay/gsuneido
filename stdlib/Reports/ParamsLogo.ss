// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
function (report)
	{
	logoFunc = OptContribution('ParamsLogo', function(@unused){ return 'Skip' })
	return logoFunc(report)
	}