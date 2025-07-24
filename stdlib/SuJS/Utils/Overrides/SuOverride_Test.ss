// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(quiet = false)
		{
		name = Display(this).BeforeFirst(' ')
		observer = .RunTest(name, :quiet)
		return observer.Result
		}

	RunTest(name, observer = false, quiet = false)
		{
		if observer is false
			observer = new TestObserverString(:quiet)
		SuTestRunner.Run1(name, observer)
		return observer
		}

	Foreach_test_method(block)
		{
		base = .Base()
		privatePrefix = Display(base).BeforeFirst(' ') $ '_'
		for member in .getbaseMembers(base).Sort!()
			.debug_one(member, privatePrefix, base, block)
		}

	getbaseMembers(base)
		{
		members = Object()
		do
			{
			members.MergeUnion(base.Members())
			base = base.Base()
			}
			while base isnt Test

		return members
		}

	debug_one(member, privatePrefix, base, block)
		{
		if member.Prefix?("Test") and
			not member.Prefix?(privatePrefix) and
			Function?(base[member])
			Finally(
				{
				block(member)
				},
				{
				.TeardownAfterEachMethod()
				})
		}

	Setup()
		{
		}

	Teardown()
		{
		}

	TeardownAfterEachMethod()
		{
		}
	}