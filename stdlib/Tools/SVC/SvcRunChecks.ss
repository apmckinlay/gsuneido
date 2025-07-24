// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
class
	{
	svcCheckThreadName: #SvcPreCheck

	CallClass(table, local_list, change, masterRec, index)
		{
		if masterRec is false or masterRec is #()
			return

		ThreadPool().Submit({ .threadFn(table, local_list, change, masterRec, index) })
		}

	threadFn(table, local_list, change, masterRec, index)
		{
		id = table $ '__' $ Display(Timestamp())
		Thread.Name(.svcCheckThreadName $ id)
		.buildWarnings(id, local_list, change, masterRec, index)
		}

	buildWarnings(id, local_list, change, masterRec, index)
		{
		SvcCommitChecker.PreCommitChecks(local_list, id, change, masterRec, index)
		}

	GetPreCheckResults(table, changes = false)
		{
		for thread in Thread.List()
			if thread.Has?('SvcPreCheck') and
				thread.AfterFirst('SvcPreCheck').Has?(table)
				return 'Thread_Running'

		return SvcCommitChecker.ProcessWarnings(table, changes)
		}
	}
