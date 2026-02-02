// Copyright (C) 2024 Axon Development Corporation All rights reserved worldwide.
Test
	{
	Test_one()
		{
		cl = DoTaskWithPause
			{
			DoTaskWithPause_createTask(msg)
				{
				return new this(msg)
				}
			DoTaskWithPause_setup(msg/*unused*/) { }
			DoTaskWithPause_cleanup() { }
			DoTaskWithPause_now(_env)
				{
				return env.now
				}
			DoTaskWithPause_pause(_env)
				{
				env.log.Add(#Pause)
				.DoTaskWithPause_cancel? = env.GetDefault(#cancelAt, false) is env.count
				}
			Finish(_env)
				{
				env.log.Add(#Finish)
				}
			}

		// 1, 2, 3, 4, #Finish
		_env = Object(log: Object(), now: Date(), count: 0)
		Assert(cl('Test', { _env.now = _env.now.Plus(seconds: 1); _env.count++ < 4 }))
		Assert(_env.log is: #(#Finish))

		// 2, 4, 6, #Pause, 8, #Finish
		_env = Object(log: Object(), now: Date(), count: 0)
		Assert(cl('Test', { _env.now = _env.now.Plus(seconds: 2); _env.count++ < 4 }))
		Assert(_env.log is: #(#Pause, #Finish))

		// 4, 8, #Pause, 12, 16, #Pause, #Finish
		_env = Object(log: Object(), now: Date(), count: 0)
		Assert(cl('Test', { _env.now = _env.now.Plus(seconds: 4); _env.count++ < 4 }))
		Assert(_env.log is: #(#Pause, #Pause, #Finish))

		// 5, #Pause, 10, #Pause, 15, #Pause, 20, #Pause, #Finish
		_env = Object(log: Object(), now: Date(), count: 0)
		Assert(cl('Test', { _env.now = _env.now.Plus(seconds: 5); _env.count++ < 4 }))
		Assert(_env.log is: #(#Pause, #Pause, #Pause, #Pause, #Finish))

		// nested
		// 2, 4, 6, #Pause, 8, 10, #Finish, 15, #Pause,
		// 17, 19, 21, #Pause, 23, 25, #Finish, 30, #Pause,
		// 32, 34, 36, #Pause, 38, 40, #Finish, 45, #Pause,
		// 47, 49, 51, #Pause, 53, 55, #Finish, 60, #Pause,
		// 62, 64, 66, #Pause, 68, 70, #Finish, 60, #Finish
		_env = Object(log: Object(), now: Date(), count: 0, count2: 0)
		Assert(cl('Test', {
			_env.count2 = 0
			cl('Test', {
				_env.now = _env.now.Plus(seconds: 2)
				_env.count2++ < 5
				})
			_env.now = _env.now.Plus(seconds: 5)
			_env.count++ < 4 }))
		Assert(_env.log
			is: #(#Pause, #Finish, #Pause, #Pause, #Finish, #Pause, #Pause, #Finish,
				#Pause, #Pause, #Finish, #Pause, #Pause, #Finish, #Finish))

		// throw after 2 calls
		// 5, #Pause, 10, #Pause, #Finish
		_env = Object(log: Object(), now: Date(), count: 0)
		Assert({ cl('Test', { _env.now = _env.now.Plus(seconds: 5)
			if _env.count is 2
				throw 'test'
			_env.count++ < 4 })} throws: 'test')
		Assert(_env.log is: #(#Pause, #Pause, #Finish))

		// cancel after 3 calls
		// 5, #Pause, 10, #Pause, 15, #Pause, #Finish
		_env = Object(log: Object(), now: Date(), count: 0, cancelAt: 3)
		Assert(cl('Test', { _env.now = _env.now.Plus(seconds: 5); _env.count++ < 4 })
			is: false)
		Assert(_env.log is: #(#Pause, #Pause, #Pause, #Finish))
		}
	}