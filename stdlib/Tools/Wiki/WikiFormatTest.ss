// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Setup()
		{
		.TearDownIfTablesNotExist('wiki')
		WikiEnsure()
		}
	Test_main()
		{
		for x in .data
			Assert(WikiFormat(x[0]) is: x[1])
		Assert(WikiFormat("[sic]" $ 'X'.Repeat(2000) $ "[/sic]")
			is: 'X'.Repeat(2000) $ '\n')
		}
	data: (
		// (in, out)
		("", "")
		("&<>", "&amp;&lt;&gt;\n")
		("hello\n\nworld", "hello\n<p>\nworld\n")
		(" preformatted", "\n<pre> preformatted\n</pre>\n")
		("|table|", '\n<table border="1" cellpadding="3">' $
			'<tr><td>table</td></tr>\n</table>\n')
		("|one|two|", '\n<table border="1" cellpadding="3">' $
			'<tr><td>one</td><td>two</td></tr>\n</table>\n')
		("|table|\n|row|", '\n<table border="1" cellpadding="3">' $
			'<tr><td>table</td></tr>\n<tr><td>row</td></tr>\n</table>\n')
		("|table|\nmore", '\n<table border="1" cellpadding="3">' $
			'<tr><td>table</td></tr>\n</table>\nmore\n')
		("http://ibm.com", '<a href="http://ibm.com">http://ibm.com</a>\n')
		("https://ibm.com", '<a href="https://ibm.com">https://ibm.com</a>\n')
		("[sic]ThisShouldNotLink[/sic]", "ThisShouldNotLink\n")
		)
	Test_table_after_ul_bug()
		{
		s = WikiFormat("
* bullet
* list

|one|two|
|three|four|")

		Assert(s has: 'border="1" cellpadding="3"')
		}
	Test_file_link_replacement()
		{
		webpage = WikiFormat.InternalFile("http://File?path/toFile")
		Assert(webpage is: '<a href="File/path/toFile">toFile</a>')

		webpage = WikiFormat.InternalFile('http://File?hello.txt')
		Assert(webpage is: '<a href="File/hello.txt">hello.txt</a>')

		webpage = WikiFormat.InternalFile('http://File/hello.txt')
		Assert(webpage is: '<a href="File/hello.txt">hello.txt</a>')

		webpage = WikiFormat.InternalFile('http://File/path/toFile')
		Assert(webpage is: '<a href="File/path/toFile">toFile</a>')

		webpage = WikiFormat.InternalFile('http://File/path/toFile.JPG')
		Assert(webpage is: '<img src="File/path/toFile.JPG">')

		webpage = WikiFormat.InternalFile('http://File/path/toFile.jpg')
		Assert(webpage is: '<img src="File/path/toFile.jpg">')

		biggerWebpage = WikiFormat("http://www.somepage.com/image.png
http://www.anotherpage.com/
CoolWikiPageLink
http://File?path/toAnotherFile
file://somefile/test.txt
http://File?path/galaxy.jpg
http://File/path/truck.jpg
")

		Assert(biggerWebpage is: '<img src="http://www.somepage.com/image.png">\n' $
'<a href="http://www.anotherpage.com/">http://www.anotherpage.com/</a>\n' $
'<a href="Wiki?edit=CoolWikiPageLink"><b>?</b></a>CoolWikiPageLink\n' $
'<a href="File/path/toAnotherFile">toAnotherFile</a>\n' $
'<a href="file://somefile/test.txt">file://somefile/test.txt</a>\n' $
'<img src="File/path/galaxy.jpg">\n' $
'<img src="File/path/truck.jpg">\n' $
'')
		}

	Test_headings()
		{
		headings = WikiFormat("
!!Headingtwo
!!!Heading3
!!!!Heading4
not a heading")

expected = '<small><a name="0" class="noPrint"><em>Headings:</em></a> ' $
'<a href="#1" class="noPrint">Headingtwo</a>&nbsp;&nbsp;' $
'<a href="#2" class="noPrint">Heading3</a>&nbsp;&nbsp;' $
'<a href="#3" class="noPrint">Heading4</a>&nbsp;&nbsp;</small><p>\n' $
'<p>\n\n' $
'<h2><a name="1">Headingtwo</a>\n' $
'</h2>\n' $
'<h3>\n<a name="2">Heading3</a>\n' $
'</h3>\n' $
'<h4>\n' $
'<a name="3">Heading4</a>\n' $
'</h4>\n' $
'not a heading\n'

//DiffControl('main', expected.Lines(), headings.Lines(), 'expected', 'headings')

		Assert(headings is: expected)
		}

	Test_text_modifiers()
		{
		textChanges = WikiFormat("
------------------------------
''emphasis''
'''I is strong'''
==monospaced and monotone==
-------------------------------")
		expected = "<p>\n" $
"<hr>\n" $
"<em>emphasis</em>\n" $
"<strong>I is strong</strong>\n" $
"<tt>monospaced and monotone</tt>\n" $
"<hr>\n"
		Assert(textChanges is: expected)
		}

	Test_assertCodeOkayToEval()
		{
		m = WikiFormat.WikiFormat_assertCodeOkayToEval

		// no code, should not throw anything
		code = ""
		m(code)

		// multiple function calls not allowed, should throw
		code = "Call1();Call2()"
		Assert({ m(code) } throws: "can not have multiple calls")

		// whitelisted, should not throw anything
		code = "SampleTestFunctionThatIsAllowed(1,5)"
		m(code)

		// not whitelisted, should throw
		code = "QueryDo('something')"
		Assert({ m(code) } throws: "function is not in the Wiki function whitelist")
		}

	Test_handleLineStart_definitionLine()
		{
		cl = WikiFormat { WikiFormat_emit(@unused) { return '' } }
		fn = cl.WikiFormat_handleLineStart

		line = ':'
		Assert(fn(line) is: ':')

		line = '::'
		Assert(fn(line) is: '<dt><dd>')

		line = ':Read: You can use this web site like any other.'
		Assert(fn(line) is: '<dt>Read<dd>You can use this web site like any other.')

		line = ':term : definition'
		Assert(fn(line) is: '<dt>term <dd>definition')

		line = '::term2:nested definition'
		Assert(fn(line) is: '<dt>term2<dd>nested definition')

		line = ':Test: This is a test: There are colons: like this: colons'
		Assert(fn(line) is: '<dt>Test: This is a test: There are colons: like' $
			' this<dd>colons')

		longLine = ":Next: Create your own page. Edit an existing page and insert the" $
			" title of the new page as two or more capitalized words run together" $
			" without spaces, [sic]LikeThis[/sic]. When you save from the edit page, " $
			"Wiki converts all words run together to links. If it can find the title" $
			" in the database of existing pages, it will put in a link to that page." $
			" If it can't find the title, it will put a question mark link next to" $
			" the term. Click on the question mark link to edit your new page. Why" $
			" not start with your name (e.g. AndrewMcKinlay) and say a little" $
			" something about yourself."
		Assert(fn(longLine) is: "<dt>Next<dd>Create your own page." $
			" Edit an existing page and insert the" $
			" title of the new page as two or more capitalized words run together" $
			" without spaces, [sic]LikeThis[/sic]. When you save from the edit page, " $
			"Wiki converts all words run together to links. If it can find the title" $
			" in the database of existing pages, it will put in a link to that page." $
			" If it can't find the title, it will put a question mark link next to" $
			" the term. Click on the question mark link to edit your new page. Why" $
			" not start with your name (e.g. AndrewMcKinlay) and say a little" $
			" something about yourself.")
		}

	Test_replace()
		{
		fn = WikiFormat.WikiFormat_replace
		block = { |s| s.Upper() }
		Assert(fn('', '<a>', '</a>', block) is: '')
		Assert(fn('aaa', '<a>', '</a>', block) is: 'aaa')
		Assert(fn('<a><a>', '<a>', '</a>', block) is: '<a><a>')
		Assert(fn('</a></a>', '<a>', '</a>', block) is: '</a></a>')
		Assert(fn('a<<a>a</a>a', '<a>', '</a>', block) is: 'a<<A>A</A>a')
		Assert(fn('a<<a>a<a></a>a', '<a>', '</a>', block) is: 'a<<A>A<A></A>a')
		Assert(fn('a<<a>a</a></a>a', '<a>', '</a>', block) is: 'a<<A>A</A></a>a')
		Assert(fn('a<<a>a</a><a>a', '<a>', '</a>', block) is: 'a<<A>A</A><a>a')
		}
	}