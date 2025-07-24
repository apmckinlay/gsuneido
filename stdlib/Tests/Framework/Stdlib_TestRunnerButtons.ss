// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
#(
	('Run QC+All',
		function (source)
			{
			Qc_Main.RunQCWith(source.Controller.GetLibs(), source.Controller.On_Run_All)
			}, front:)
	('Check Quality',
		function (source)
			{
			Qc_Main.RunQC(source.Controller.GetLibs())
			})
	('Check Libraries',
		function () {
			Working('Checking libraries ...')
				{
				s = CheckLibraries(Libraries().Remove('demobookoptions'))
				if s is ""
					s = "no errors found"
				Print("Check Libraries")
				Print(s)
				}
			Alert(s, 'Check Libraries')
			})
	('Check Views',
		function () {
			Working('Checking views ...')
				{
				s = CheckViews()
				if s is ""
					s = "no errors found"
				Print("Check Views")
				Print(s)
				}
			Alert(s, 'Check Views')
			})
)
