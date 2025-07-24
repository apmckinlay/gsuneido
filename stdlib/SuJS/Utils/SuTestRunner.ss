// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
class
	{
	Run1(name, observer = false, timeEachMethod? = false)
		{
		instance = Construct(name)
		if observer is false
			observer = TestObserverPrint()
		.wrap(observer)
			{
			.run(name, instance, observer, timeEachMethod?)
			}
		}

	RunList(tests, observer)
		{
		return .wrap(observer)
			{
			for name in tests.Members()
				{
				.run(name, Construct(name), observer, excludeMethods: tests[name])
				}
			}
		}

	wrap(observer, block)
		{
		m = .measure(block)
		return observer.After(m.time, dbgrowth: m.dbgrowth, memory: m.memory)
		}

	measure(block)
		{
		t = Timer(block)
		return Object(
			time: t < 0.1 ? 0 : t,
			dbgrowth: 0,
			memory: 0)
		}

	run(name, instance, observer, timeEachMethod? = false, excludeMethods = #())
		{
		observer.BeforeTest(name)

		m = .measure()
			{
			.run1(instance, observer, timeEachMethod?, excludeMethods)
			}

		observer.AfterTest(name, m.time, dbgrowth: m.dbgrowth, memory: m.memory)
		}

	run1(instance, observer, timeEachMethod?, excludeMethods)
		{
		if .runMethod(observer, instance, 'Setup', timeEachMethod?)
			instance.Foreach_test_method
				{ |method|
				if not excludeMethods.Has?(method)
					.runMethod(observer, instance, method, timeEachMethod?)
				}
		.runMethod(observer, instance, 'Teardown', timeEachMethod?)
		}

	runMethod(observer, instance, method, timeEachMethod?)
		{
		observer.BeforeMethod(method)
		time = false
		if timeEachMethod?
			time = Timer({ passed? = .tryMethod(observer, instance, method) })
		else
			passed? = .tryMethod(observer, instance, method)
		observer.AfterMethod(method, :time)
		return passed?
		}

	tryMethod(observer, instance, method)
		{
		try
			instance[method]()
		catch (e)
			{
			.error(observer, method, e)
			return false
			}
		return true
		}

	error(observer, where, msg)
		{
		observer.Error(where, msg)
		}
	}