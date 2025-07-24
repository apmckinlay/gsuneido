// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		tp = new .mockThreadPool

		// Fresh start
		// No existing thread running
		taskLog = Object()
		task = function (taskLog)
			{
			taskLog.Add('Current task done')
			}
		tp.Submit({ task(taskLog) })
		Assert(tp.ThreadPool_nThreads is: 0)
		Assert(taskLog is: #('Current task done'))

		// One more thread allowed to run
		// Queue is not empty
		taskLog = Object()
		task1 = function (taskLog)
			{
			taskLog.Add('task1 done')
			}
		task2 = function (taskLog)
			{
			taskLog.Add('task2 done')
			}
		task3 = function (taskLog)
			{
			taskLog.Add('task3 done')
			}
		tp.ThreadPool_queue = Object({ task1(taskLog) }, { task2(taskLog) },
			{ task3(taskLog) })
		tp.ThreadPool_nThreads = 3
		tp.Submit({ task(taskLog) })
		Assert(tp.ThreadPool_nThreads is: 3)
		Assert(taskLog is: #('Current task done', 'task1 done', 'task2 done',
			'task3 done'))

		// One more thread allowed to run
		// Newly-started task throws exception
		// Queue is not empty
		taskLog = Object()
		_errLog = Object()
		task = function ()
			{
			throw 'Task Throwing'
			}
		tp.ThreadPool_queue = Object({ task1(taskLog) }, { task2(taskLog) },
			{ task3(taskLog) })
		tp.ThreadPool_nThreads = 3
		tp.Submit({ task() })
		Assert(tp.ThreadPool_nThreads is: 3)
		Assert(taskLog is: #('task1 done', 'task2 done', 'task3 done'))
		Assert(_errLog isSize: 1)
		Assert(_errLog[0] is: 'Task Throwing')

		// One more thread allowed to run
		// Both task and SuneidoLog throw
		taskLog = Object()
		_throwLog = Object()
		tptt = new .threadPoolTestingThrow
		task = function ()
			{
			throw 'Task Throwing'
			}
		tptt.ThreadPool_queue = Object({ task1(taskLog) }, { task2(taskLog) },
			{ task3(taskLog) })
		tptt.ThreadPool_nThreads = 3
		// There is no return value for the Submit method so 'throws' can't be used here
		// Can not test the size of queue after throwing exception since there is no real
		// exception that stops the code from running
		tptt.Submit({ task() })
		Assert(_throwLog isSize: 1)
		Assert(_throwLog[0] is: 'Task Throwing')
		Assert(tptt.ThreadPool_nThreads is: 3)
		}

	mockThreadPool: ThreadPool
		{
		New()
			{
			.ThreadPool_queue = Object()
			}
		Submit(task)
			{
			if .ThreadPool_nThreads < .ThreadPool_maxThreads
				{
				++.ThreadPool_nThreads
				.ThreadPool_worker(task)
				}
			else if .ThreadPool_queue.Size() < .ThreadPool_maxQueue
				.ThreadPool_queue.Add(task)
			else
				task()
			}
		ThreadPool_logException(e, _errLog = false)
			{
			if Object?(errLog)
				errLog.Add(e)
			}
		}

	threadPoolTestingThrow: ThreadPool
		{
		New()
			{
			.ThreadPool_queue = Object()
			}
		Submit(task)
			{
			if .ThreadPool_nThreads < .ThreadPool_maxThreads
				{
				++.ThreadPool_nThreads
				.ThreadPool_worker(task)
				}
			else if .ThreadPool_queue.Size() < .ThreadPool_maxQueue
				.ThreadPool_queue.Add(task)
			else
				task()
			}
		ThreadPool_logException(unused)
			{
			throw 'Testing SuneidoLog Throwing'
			}
		ThreadPool_throwException(e, _throwLog = false)
			{
			if Object?(throwLog)
				throwLog.Add(e)
			}
		}
	}
