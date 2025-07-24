// Copyright (C) 2021 Suneido Software Corp. All rights reserved worldwide.
class
	{
	title: 'Svc Sync Tables'
	CallClass(table = false)
		{
		outOfSync = .discrepancies(.tables(table), svc = .svc())
		if outOfSync.Empty?()
			Alert('No discrepancies to correct', .title, flags: MB.ICONINFORMATION)
		else
			{
			outOfSync.Members().Sort!().Any?({ .Sync(it, svc, outOfSync[it]) })
			Unload()
			}
		}

	tables(table)
		{
		return table isnt false
			? Object(table)
			: LibraryTables().
				Append(BookTables()).
				Remove(@SvcControl.SvcExcludeLibraries)
		}

	svc()
		{
		settings = SvcSettings.Get()
		return Svc(settings.svc_server, settings.svc_local?)
		}

	discrepancies(tables, svc)
		{
		outOfSync = Object()
		tables.Each()
			{
			discrepancies = svc.Compare(SvcTable(it))
			if discrepancies.NotEmpty?()
				outOfSync[it] = discrepancies
			}
		return outOfSync
		}

	Sync(table, svc, discrepancies)
		{
		try
			{
			changes = svc.GetChanges(table)
			svc.UpdateLibrary(changes.master_changes)
			if changes.local_changes.Empty?() or
				.continue?(table $ ' has local changes.\r\n\r\nOverwrite?')
				.overwrite(table, svc, discrepancies)
			}
		catch (e)
			return not .continue?('An error occurred while syncing: ' $ table $ '\r\n' $
				'Error: ' $ e $ '\r\n\r\nContinue?')
		return false
		}

	continue?(msg)
		{
		return YesNo(msg, .title, flags: MB.ICONWARNING)
		}

	overwrite(table, svc, discrepancies)
		{
		svc.Overwrite(discrepancies.Map!({ [:table, name: it.AfterFirst(' ').Trim()] }))
		}

	Discrepancies(table = false, svc = false)
		{
		if svc is false
			svc = .svc()
		return .discrepancies(.tables(table), svc)
		}
	}
