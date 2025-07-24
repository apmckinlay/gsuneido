// Copyright (C) 2006 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_GetState()
		{
		// test with most book components missing (not constructed)
		book = BookControl
			{
			LogFont: "BookControl_LogFont"
			BookControl_login: "BookControl_login"
			Title: "Test Book Title"
			}
		Assert(book.GetState() is: Object(login: "BookControl_login",
			title: "Test Book Title",
			help_book: false
			book: false).Add("BookControl_LogFont", at: GetCustomFontSaveName()))

		// test all items in GetState
		book = BookControl
			{
			LogFont: "BookControl_LogFont2"
			BookControl_login: "BookControl_login2"
			Title: "Test Book Title2"
			BookControl_book: "test_book"
			BookControl_tree: class { GetState() { return "tree_state" }}
			BookControl_marks: class { GetState() { return "marks_state" }}
			BookControl_toolbar: class { GetState() { return "toolbar_state" }}
			BookControl_start: "StartPage"
			BookControl_marksplit: class { GetSplit() { return "mark_split" }}
			BookControl_treesplit: class { GetSplit() { return "tree_split" }}
			BookControl_help_book: false
			}
		Assert(book.GetState() is: Object(login: "BookControl_login2",
			title: "Test Book Title2",
			book: "test_book",
			tree: "tree_state",
			marks: "marks_state",
			toolbar: "toolbar_state",
			start: "StartPage",
			marksplit: "mark_split",
			treesplit: "tree_split",
			help_book: false
			).Add("BookControl_LogFont2", at: GetCustomFontSaveName()))
		}

	Test_help_book()
		{
		book = BookControl
			{
			LogFont: "BookControl_LogFont2"
			BookControl_login: "BookControl_login2"
			Title: "Test Book Title2"
			BookControl_book: "test_book"
			BookControl_tree: class { GetState() { return "tree_state" }}
			BookControl_marks: class { GetState() { return "marks_state" }}
			BookControl_toolbar: class { GetState() { return "toolbar_state" }}
			BookControl_start: "StartPage"
			BookControl_marksplit: class { GetSplit() { return "mark_split" }}
			BookControl_treesplit: class { GetSplit() { return "tree_split" }}
			BookControl_help_book: true
			}
		Assert(book.GetState() is: Object(login: "BookControl_login2",
			title: "Test Book Title2",
			book: "test_book",
			tree: "tree_state",
			marks: "marks_state",
			toolbar: "toolbar_state",
			start: "StartPage",
			marksplit: "mark_split",
			treesplit: "tree_split",
			help_book: true
			))
		}

	Test_url_to_name()
		{
		f = BookControl.BookControl_url_to_name
		Assert(f('foo') is: '')
		Assert(f('suneido:/one/two/three') is: '/two/three')
		}

	Test_buildWikiNotesURLName()
		{
		m = BookControl.BookControl_buildWikiNotesURLName

		path = ""
		name = ""
		Assert(m(path, name) is: "")

		path = ""
		name = "Something"
		Assert(m(path, name) is: "Something")

		path = "Location1"
		name = "Item2"
		Assert(m(path, name) is: "Location1Item2")

		path = "Location"
		name = "One And Two"
		Assert(m(path, name) is: "LocationOneAndTwo")

		path = "/Location/Location2/"
		name = "Something &@- Something Else"
		Assert(m(path, name) is: "LocationLocation2SomethingSomethingElse")

		path = "/Base/Here And There/"
		name = "One And Two"
		Assert(m(path, name) is: "BaseHereAndThereOneAndTwo")
		}
	}