// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	table_created?: false
	Setup()
		{
		// if translation doesn't exist, or is empty, then create it and output
		// some test records
		if TableExists?('translatelanguage')
			.delete_test_data()
		else
			.table_created? = true
		.addLines()
		.language = Suneido.Language
		Suneido.Language = #(name: "translatelanguagetest", charset: "DEFAULT")
		}
	values: #(xxx: 'y&yy', eee: '&fff', hhh: 'iii')
	addLines()
		{
		Database("ensure translatelanguage (trlang_from, trlang_translatelanguagetest)
			key(trlang_from)")
		for val in .values.Members()
			QueryOutput("translatelanguage",
				Record(trlang_from: val, trlang_translatelanguagetest: .values[val]))
		}
	Test_main()
		{
		Assert(TranslateLanguage("xxx") is: "yyy")
		Assert(TranslateLanguage("eee") is: "fff")
		Assert(TranslateLanguage("hhh") is: "iii")
		// test handling of  "..."
		Assert(TranslateLanguage("xxx...") is: "yyy...")
		// test stripping off ampersand
		Assert(TranslateLanguage("e&e&e") is: "&fff")
		// test legitimate ampersand
		Assert(TranslateLanguage("&hhh") is: "iii")
		// test non-existent
		Assert(TranslateLanguage("XXXXXXXXXXXXXXXXXX") is: "XXXXXXXXXXXXXXXXXX")
		// test empty string
		Assert(TranslateLanguage("") is: "")
		}
	Test_args()
		{
		Assert(TranslateLanguage("Can't open %1 for input", "temp"),
			is: "Can't open temp for input")
		Assert(TranslateLanguage("Can't open %file for input", file: "temp"),
			is: "Can't open temp for input")
		Assert(TranslateLanguage("Can't %1", "xxx"),
			is: "Can't yyy")
		Assert(TranslateLanguage("Can't %file", file: "xxx"),
			is: "Can't yyy")
		}
	Test_nonexistent_language()
		{
		language = Suneido.Language
		Suneido.Language = #(name: "translatelanguagetest2", charset: "DEFAULT")
		Assert(TranslateLanguage("xxx") is: "xxx")
		Assert(TranslateLanguage("eee") is: "eee")
		Assert(TranslateLanguage("hhh") is: "hhh")
		Suneido.Language = language
		}
	Test_trim()
		{
		Assert(TranslateLanguage("  xxx") is: "  yyy")
		Assert(TranslateLanguage("xxx  ") is: "yyy  ")
		Assert(TranslateLanguage("xxx...  ") is: "yyy...  ")
		Assert(TranslateLanguage("   xxx  ") is: "   yyy  ")
		}
	Teardown()
		{
		Suneido.Language = .language
		Database('alter translatelanguage drop (trlang_translatelanguagetest)')
		.delete_test_data()
		if .table_created?
			Database("destroy translatelanguage")
		}
	delete_test_data()
		{
		for val in .values.Members()
			QueryDo('delete translatelanguage where trlang_from is ' $ Display(val))
		}
	}
