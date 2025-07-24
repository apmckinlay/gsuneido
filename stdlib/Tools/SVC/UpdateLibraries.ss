// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
class
	{
	errorPrefix: 'ERROR: UpdateLibraries: '
	CallClass(tables, preGetFn = false)
		{
		result = .updateLibraries(tables, preGetFn)
		SvcSocketClient().Close()
		return result
		}

	updateLibraries(tables, preGetFn)
		{
		try
			{
			if String?(svc = .svc())
				throw svc
			status = svc.CheckSvcStatus()
			result = status isnt ''
				? .errorPrefix $ 'Svc status: ' $ status
				: .getChanges(tables, svc, preGetFn)
			}
		catch (e)
			result = .errorPrefix $ e
		return result
		}

	svc()
		{
		if false is svcSetting = SvcSettings()
			return 'Invalid Version Control settings'
		return Svc(server: svcSetting.svc_server)
		}

	getChanges(tables, svc, preGetFn)
		{
		nChanges = 0
		result = ''
		for table in tables
			{
			if not TableExists?(table)
				{
				.log('WARNING: ' $ table $ ' does not exist', '')
				continue
				}
			changes = svc.GetChanges(table)
			if not changes.conflicts.Empty?()
				{
				.log(.errorPrefix $ table $ ' has conflicts', .params(changes.conflicts))
				result $= 'FAILURES: Version Control Conflicts in ' $ table $ '\n'
				continue
				}
			.checkLocalChange(changes, table)
			.preGet(QueryColumns(table).Has?('group'), preGetFn, changes, svc)
			nChanges += svc.UpdateLibrary(changes.master_changes)
			}
		return result isnt '' ? result : nChanges
		}

	log(msg, params)
		{ SuneidoLog(msg, :params) }

	params(changes)
		{ return changes.Map({ it.name }).Sort!() }

	checkLocalChange(changes, table)
		{
		if not changes.local_changes.Empty?()
			.log('WARNING: ' $ table $ ' has local changes',
				.params(changes.local_changes))
		}

	preGet(lib?, preGetFn, changes, svc)
		{
		if lib? and preGetFn isnt false
			preGetFn(changes.master_changes, svc)
		}
	}
