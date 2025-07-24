// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
// Usage: ThreadPool().Submit(callable)
Singleton
	{
	maxThreads: 4
	maxQueue: 256
	nThreads: 0
	getter_queue()
		{
		return .queue = Object() // once only
		}
	Submit(task)
		{
		// potential race conditions, small chance that threads or queue could exceed max
		if .nThreads < .maxThreads
			{ // haven't reached max, start another worker thread
			++.nThreads
			Thread({ .worker(task) })
			}
		else if .queue.Size() < .maxQueue
			{ // queue not full, add to it
			.queue.Add(task)
			}
		else // queue full, run in caller thread to throttle it
			task()
		}
	worker(task)
		{
		Finally(
			{
			forever
				{
				try
					task()
				catch (e)
					{
					try
						.logException(e)
					catch
						.throwException(e)
					}
				// concurrency depends on PopFirst being atomic (thread safe)
				if Same?(.queue, task = .queue.PopFirst())
					break // queue empty, end worker thread
				}
			},
			{
			--.nThreads
			})
		}
	logException(e) // overridden by test
		{
		SuneidoLog('ERROR: ' $ e)
		}
	throwException(e) // overridden by test
		{
		throw e
		}
	ClearQueue() // note: doesn't stop running tasks
		{
		.queue = Object()
		}
	Reset()
		{
		}
	}
