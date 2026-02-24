// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	threads: 4
	Reps: 1000
	Test_object_add()
		{
		ob = Object()
		ob.done = 0
		Assert(Concurrent?(ob) is: false)
		Suneido.ConcurrencyTest = ob
		Assert(Concurrent?(ob))
		for .. .threads
			Thread(function()
				{
				Assert(Concurrent?(Suneido.ConcurrencyTest))
				for .. ConcurrencyTest.Reps
					Suneido.ConcurrencyTest.Add(1)
				++Suneido.ConcurrencyTest.done
				})
		for (sleep = 1; sleep < 1000 and ob.done < .threads; sleep *= 2)
			Thread.Sleep(sleep) // wait for threads to finish
		Assert(ob.done is: .threads)
		Assert(ob.Size() is: 1 + .threads * .Reps)
		}
	Test_object()
		{
		ob = Object()
		ob.n = ob.done = 0
		Assert(Concurrent?(ob) is: false)
		Suneido.ConcurrencyTest = ob
		Assert(Concurrent?(ob))
		for .. .threads
			Thread(function()
				{
				Assert(Concurrent?(Suneido.ConcurrencyTest))
				for .. ConcurrencyTest.Reps
					++Suneido.ConcurrencyTest.n
				++Suneido.ConcurrencyTest.done
				})
		for (sleep = 1; sleep < 1000 and ob.done < .threads; sleep *= 2)
			Thread.Sleep(sleep) // wait for threads to finish
		Assert(ob.done is: .threads)
		Assert(ob.n is: .threads * .Reps)
		}
	Test_instance()
		{
		ob = class{}()
		ob.n = ob.done = 0
		Assert(Concurrent?(ob) is: false)
		Suneido.ConcurrencyTest = ob
		Assert(Concurrent?(ob))
		for .. .threads
			Thread(function()
				{
				Assert(Concurrent?(Suneido.ConcurrencyTest))
				for .. ConcurrencyTest.Reps
					++Suneido.ConcurrencyTest.n
				++Suneido.ConcurrencyTest.done
				})
		for (sleep = 1; sleep < 1000 and ob.done < .threads; sleep *= 2)
			Thread.Sleep(sleep) // wait for threads to finish
		Assert(ob.done is: .threads)
		Assert(ob.n is: .threads * .Reps)
		}
	Test_exception_callstack_bug()
		{
		x = Object()
		Assert(Concurrent?(x) is: false)
		try
			throw 'foo'
		catch (e)
			Suneido.ConcurrencyTest = e // make it concurrent
		Assert(Concurrent?(x))
		}
	Teardown()
		{
		Suneido.Delete(#ConcurrencyTest)
		super.Teardown()
		}
	}
