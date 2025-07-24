// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
class
	{
	New(.table, .changes, .local_list, .settings)
		{
		}

	CallClass(table, local_list, settings)
		{
		if .stop(local_list, table)
			return

		ThreadPool().ClearQueue() // WARNING: potentially clear other's tasks
		if Thread.List().Any?({ it.Has?('SvcGetMaster') })
			{
			Suneido.SvcCommit_Warnings.newGetMaster = true
			Suneido.SvcCommit_Warnings.table = table
			Suneido.SvcCommit_Warnings.local_list = local_list
			return
			}

		if not Suneido.Member?('SvcCommit_Warnings')
			SvcCommitChecker.ClearPreCheck()
		if not Suneido.Member?('SvcPreCheck_ForceStop')
			Suneido.SvcPreCheck_ForceStop = Object()
		Suneido.Svc_Check_Batch = Timestamp()
		for thread in Thread.List().Filter({
			it.Has?('__') and it.AfterFirst('__') < Suneido.Svc_Check_Batch})
			Suneido.SvcPreCheck_ForceStop[thread.AfterFirst('SvcPreCheck')] = true

		changes = .getLocalRecs(local_list)
		ThreadPool().Submit(new this(table, changes, local_list, settings))
		}

	stop(local_list, table)
		{
		return .runningTestRunner() or local_list.Get().Empty?() or table is ''
		}

	Call()
		{
		try
			.run()
		catch (e)
			SuneidoLog('ERROR SvcGetMaster: ' $ e)
		}

	runningTestRunner()
		{
		return Suneido.GetDefault('TestRunner', false)
		}

	getLocalRecs(local_list)
		{
		return local_list.Get().
			Map({Object(name: it.svc_name, lib: it.svc_lib, type: it.svc_type )}).
			Instantiate()
		}

	run()
		{
		Thread.Name('SvcGetMaster' $ Display(Timestamp()))
		.svc = Svc(server: .settings.svc_server, local?: .settings.svc_local?)
		index = 0
		for change in .changes
			{
			if not .local_list.Member?('Hwnd') or .forceStop()
				return false
			masterRec = .getMasterRec(change.lib, change.name)
			SvcRunChecks(.table, .local_list, change, masterRec, index)
			index++
			}
		}

	forceStop()
		{
		if .needNew() is true
			{
			.startNewCheck()
			return true
			}
		return false
		}

	needNew()
		{
		return Suneido.SvcCommit_Warnings.newGetMaster
		}

	startNewCheck()
		{
		Suneido.SvcCommit_Warnings.newGetMaster = false
		Thread.Name('SvcStopGetMaster' $ Display(Timestamp()))
		SvcGetMaster(Suneido.SvcCommit_Warnings.table,
			Suneido.SvcCommit_Warnings.local_list, .settings)
		}

	getMasterRec(lib, name)
		{
		return [master: .svc.Get(lib, name), missingTestOld: .svc.MissingTest?(lib, name)]
		}
	}
