// Copyright (C) 2023 Axon Development Corporation All rights reserved worldwide.
class
	{
	CallClass(table)
		{
		if Sys.Client?()
			return ServerEval('ReIndexHelp', table)
		if Contributions('HelpBook').Has?(table)
			Thread({ .checkAndGenerate(table) })
		}

	checkAndGenerate(table)
		{
		if not TableExists?(table)
			return

		indexPath = Paths.Combine(ExeDir(), 'index_' $ table)
		// index dont exist, create one
		if not FileExists?(indexPath)
			return .generate('Generating', table)

		latest = QueryMax(table, #lib_committed)
		current = Dir(indexPath, details:)[0].date
		if current < latest
			.generate('Re-Generating', table)
		}

	generate(prefix, table)
		{
		title = prefix $ ' FtSearch for ' $ table
		if '' isnt msg = ServerEval('IndexHelp', table)
			SuneidoLog('ERROR: Process failed while ' $ title, params: [:msg])
		}
	}