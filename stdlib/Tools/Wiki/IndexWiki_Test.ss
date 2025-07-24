// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_wikiWord()
		{
		ww = IndexWiki.IndexWiki_wikiWord
		Assert("foobar" !~ ww)
		Assert("Foobar" !~ ww)
		Assert("fooBar" !~ ww)
		Assert("IBM" !~ ww)
		Assert("123" !~ ww)
		Assert("Foo123" !~ ww)
		Assert("123Foo" !~ ww)

		Assert("WikiWord" matches: ww)
		Assert("NowIsTheTime" matches: ww)
		Assert("April22Meeting2025" matches: ww)
		Assert("WikiWord" matches: ww)
		}
	Test_HandleWikiWords()
		{
		test = function(text, expected)
			{
			Assert(IndexWiki.HandleWikiWords(text) is: expected)
			}
		test("", "")
		test("x", "x")
		test(s = "now is the time", s)
		test("FooBar", "FooBar Foo Bar")
		test("the FooBar thing", "the FooBar Foo Bar thing")
		}
	}