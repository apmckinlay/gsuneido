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
	Test_closure()
		{
		ob = Object(n: 0, done: 0)
		Assert(Concurrent?(ob) is: false)
		for .. .threads
			Thread()
				{
				Assert(Concurrent?(ob))
				for .. .Reps
					++ob.n
				++ob.done
				}
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
	n: 100
	Test_closure_race_bug()
		{
		n = .n // make it a closure (shared variable)
		c = {|x| x + n }
		Suneido.ConcurrencyTest = c // make it concurrent
		Assert(Concurrent?(c))
		Assert(c(1) is: 101)
		Suneido.tmp = 0
		for ..10
			{
			Thread({ Suneido.tmp += c(2) - 102 })
			Thread({ Suneido.tmp += c(3) - 103 })
			Thread({ Suneido.tmp += c(4) - 104 })
			Thread({ Suneido.tmp += c(5) - 105 })
			Thread({ Suneido.tmp += c(6) - 106 })
			Thread({ Suneido.tmp += c(7) - 107 })
			Thread({ Suneido.tmp += c(8) - 108 })
			Thread({ Suneido.tmp += c(9) - 109 })
			}
		// give threads time to run (but not critical if they're all done)
		Thread.Sleep(10)
		Assert(Suneido.tmp is: 0)
		// can't clean up Suneido.tmp because we don't know when threads are done
		}
	Teardown()
		{
		Suneido.Delete(#ConcurrencyTest)
		super.Teardown()
		}
	}
