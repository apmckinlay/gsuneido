// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
// NOTE: this is deliberately not named _Test
// because we don't want it running when you run the tests
Test
	{
	Setup()
		{
		.serverTags = ServerEval('LibraryTags.GetTagsInUse')
		ServerEval('Suneido.LibraryTags')

		.clientTags = LibraryTags.GetTagsInUse()
		Suneido.LibraryTags()
		}

	Test_non_client()
		{
		if Client?()
			return

		Assert(Suneido.Info("library.tags").SafeEval() is: #(""))

		Assert(TestLibTags is: 123)

		Suneido.LibraryTags("tlt1")
		Assert(TestLibTags is: 456)

		Suneido.LibraryTags("tlt2")
		Assert(TestLibTags is: 789)

		Suneido.LibraryTags("tlt1", "tlt2")
		Assert(TestLibTags is: 789)

		Suneido.LibraryTags("tlt2", "tlt1")
		Assert(TestLibTags is: 456)
		}
	Test_client()
		{
		if not Client?()
			return

		Assert(Suneido.Info("library.tags").SafeEval() is: #())
		Assert(ServerEval("Suneido.Info", "library.tags").SafeEval() is: #(""))

		ServerEval("Suneido.LibraryTags", "tlt1", "tlt2")
		Unload("TestLibTags")
		Assert(TestLibTags is: 789)

		Suneido.LibraryTags("tlt1")
		Assert(TestLibTags is: 456)
		}
	Teardown()
		{
		ob = LibraryTags.ConvertTagInfo(.serverTags)
		ob.Add('Suneido.LibraryTags', at: 0)
		ServerEval(@ob)
		Suneido.LibraryTags(@LibraryTags.ConvertTagInfo(.clientTags))
		}
	}