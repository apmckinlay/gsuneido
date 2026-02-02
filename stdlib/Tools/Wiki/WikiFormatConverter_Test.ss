// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		for x in .data
			Assert(.format(x[0]) is: x[1])
		Assert(.format("[sic]" $ 'X'.Repeat(2000) $ "[/sic]")
			is: '<p>' $ 'X'.Repeat(2000) $ '</p>\n')
		}

	addon: Md_Addon_Wiki
		{
		Md_Addon_Wiki_nameNotExist?(name)
			{
			return name isnt 'FooBar'
			}
		}

	format(s)
		{
		return MarkdownToHtml(WikiFormatConverter(s), noIndent?:,
			addons: [
				[Md_Addon_Table, #(border: 1, cellpadding: 3)],
				Md_Addon_Definition,
				Md_Addon_suneido_style,
				.addon])
		}

	data: (
		// (in, out)
		("", "")
		("&<>", "<p>&amp;&lt;&gt;</p>\n")
		("hello\n\nworld", "<p>hello</p>\n<p>world</p>\n")
		(" preformatted", "<pre><code> preformatted\n</code></pre>\n")
		("|table|", '<table border="1" cellpadding="3">\n' $
			'<tr>\n<th style="text-align: left;">table</th>\n</tr>\n' $
			'</table>\n')
		("|one|two|", '<table border="1" cellpadding="3">\n' $
			'<tr>\n<th style="text-align: left;">one</th>\n' $
				'<th style="text-align: left;">two</th>\n</tr>\n' $
			'</table>\n')
		("|table|\n|row|", '<table border="1" cellpadding="3">\n' $
			'<tr>\n<th style="text-align: left;">table</th>\n</tr>\n' $
			'<tr>\n<td style="text-align: left;">row</td>\n</tr>\n' $
			'</table>\n')
		("|table|\nmore", '<table border="1" cellpadding="3">\n' $
			'<tr>\n<th style="text-align: left;">table</th>\n</tr>\n</table>\n' $
			'<p>more</p>\n')
		("http://ibm.com", '<p><a href="http://ibm.com">http://ibm.com</a></p>\n')
		("https://ibm.com",
			'<p><a href="https://ibm.com">https://ibm.com</a></p>\n')
		("[sic]ThisShouldNotLink[/sic]", "<p>ThisShouldNotLink</p>\n")
		)

	Test_table_after_ul_bug()
		{
		s = .format("
* bullet
* list

|one|two|
|three|four|")

		Assert(s has: 'border="1" cellpadding="3"')
		}

	Test_file_link_replacement()
		{
		webpage = Md_Addon_Wiki.InternalFile("http://File?path/toFile")
		Assert(webpage is: '<a href="File/path/toFile">toFile</a>')

		webpage = Md_Addon_Wiki.InternalFile('http://File?hello.txt')
		Assert(webpage is: '<a href="File/hello.txt">hello.txt</a>')

		webpage = Md_Addon_Wiki.InternalFile('http://File/hello.txt')
		Assert(webpage is: '<a href="File/hello.txt">hello.txt</a>')

		webpage = Md_Addon_Wiki.InternalFile('http://File/path/toFile')
		Assert(webpage is: '<a href="File/path/toFile">toFile</a>')

		webpage = Md_Addon_Wiki.InternalFile('http://File/path/toFile.JPG')
		Assert(webpage is: '<img src="File/path/toFile.JPG">')

		webpage = Md_Addon_Wiki.InternalFile('http://File/path/toFile.jpg')
		Assert(webpage is: '<img src="File/path/toFile.jpg">')

		biggerWebpage = .format("http://www.somepage.com/image.png
http://www.anotherpage.com/
CoolWikiPageLink
http://File?path/toAnotherFile
file://somefile/test.txt
http://File?path/galaxy.jpg
http://File/path/truck.jpg
")

		Assert(biggerWebpage is: '<p><img src="http://www.somepage.com/image.png">\n' $
'<a href="http://www.anotherpage.com/">http://www.anotherpage.com/</a>\n' $
'<a href="Wiki?edit=CoolWikiPageLink"><b>?</b></a>CoolWikiPageLink\n' $
'<a href="File/path/toAnotherFile">toAnotherFile</a>\n' $
'<a href="file://somefile/test.txt">file://somefile/test.txt</a>\n' $
'<img src="File/path/galaxy.jpg">\n' $
'<img src="File/path/truck.jpg"></p>\n')
		}

	Test_headings()
		{
		headings = .format("
!!Headingtwo
!!!Heading3
!!!!Heading4
not a heading")

expected = '<small><a name="0" class="noPrint"><em>Headings:</em></a> ' $
'<a href="#1" class="noPrint">Headingtwo</a>&nbsp;&nbsp;\n' $
'<a href="#2" class="noPrint">Heading3</a>&nbsp;&nbsp;\n' $
'<a href="#3" class="noPrint">Heading4</a>&nbsp;&nbsp;</small>\n' $
'<a name="1">\n<h2>Headingtwo</h2>\n</a>\n' $
'<a name="2">\n<h3>Heading3</h3>\n</a>\n' $
'<a name="3">\n<h4>Heading4</h4>\n</a>\n' $
'<p>not a heading</p>\n'

		Assert(headings is: expected)
		}

	Test_handleLineStart_definitionLine()
		{
		line = ':'
		Assert(.format(line) is: '<p>:</p>\n')

		line = '::'
		Assert(.format(line) is: '')

		line = ':Read: You can use this web site like any other.'
		Assert(.format(line) like: '<dl>\n<dt>Read</dt>\n' $
			'<dd>You can use this web site like any other.</dd>\n</dl>')

		line = ':term : definition'
		Assert(.format(line) is: '<dl>\n<dt>term</dt>\n<dd>definition</dd>\n</dl>\n')

		line = '::term2:nested definition'
		Assert(.format(line)
			is: '<dl>\n<dt>term2</dt>\n<dd>nested definition</dd>\n</dl>\n')

		line = ':Test: This is a test: There are colons: like this: colons'
		Assert(.format(line) is: '<dl>\n' $
			'<dt>Test: This is a test: There are colons: like this</dt>\n' $
			'<dd>colons</dd>\n</dl>\n')

		longLine = ":Next: Create your own page. Edit an existing page and insert the" $
			" title of the new page as two or more capitalized words run together" $
			" without spaces, [sic]LikeThis[/sic]. When you save from the edit page, " $
			"Wiki converts all words run together to links. If it can find the title" $
			" in the database of existing pages, it will put in a link to that page." $
			" If it can't find the title, it will put a question mark link next to" $
			" the term. Click on the question mark link to edit your new page. Why" $
			" not start with your name (e.g. FooBar) and say a little" $
			" something about yourself."
		Assert(.format(longLine) like: "<dl>\n<dt>Next</dt>\n<dd>Create your own page." $
			" Edit an existing page and insert the" $
			" title of the new page as two or more capitalized words run together" $
			" without spaces, LikeThis. When you save from the edit page, " $
			"Wiki converts all words run together to links. If it can find the title" $
			" in the database of existing pages, it will put in a link to that page." $
			" If it can't find the title, it will put a question mark link next to" $
			" the term. Click on the question mark link to edit your new page. Why" $
			" not start with your name (e.g. " $
			"<a href=\"Wiki?FooBar\">FooBar</a>) and say a little" $
			" something about yourself.</dd>\n</dl>")
		}
	}