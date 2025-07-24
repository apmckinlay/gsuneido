// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		book = .MakeBook()

		s = BookMenuPage('/' $ book $ '/nonexistent', 3)
		Assert(s like:
			'<h1>nonexistent</h1>

			<table width="100%">

			</table>
			')

		s = BookMenuPage('/' $ book $ '/nonexistent', 3,
			before: '{', after: '}', headingLevel: 2)
		Assert(s like:
			'<h2>nonexistent</h2>
			{
			<table width="100%">

			</table>
			}')
		QueryOutput(book, #{num: 4, path: '', name: 'section' })
		QueryOutput(book, #{num: 1, path: '', name: 'section2' })
		QueryOutput(book, #{num: 5, path: '/section', name: 'one' })
		QueryOutput(book, #{num: 2, path: '/section', name: 'two' })
		QueryOutput(book, #{num: 3, path: '/section2', name: 'three', })
		BookModel.ClearCache(book)
		s = BookMenuPage('/' $ book $ '/section', 3)
		Assert(s like:
			'<h1>section</h1>

			<table width="100%">

			<tr>
				<td><a href="/' $ book $ '/section/one">one</a></td>
				<td><a href="/' $ book $ '/section/two">two</a></td>

			</tr>

			</table>
			')
		}
	}
