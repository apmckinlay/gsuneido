// Copyright (C) 2008 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_base()
		{
		book = .MakeBook()
		svcTable = SvcTable(book)
		svcTable.Output(rec = [name: #folder])
		svcTable.Output([name: #one, path: '/folder'])
		svcTable.Output([name: #two, path: '/folder'])
		svcTable.Output([name: #nested_folder, path: '/folder'])
		Assert(svcTable.Dir(rec.num) like: 'nested_folder\none\ntwo\n')
		}

	Test_newLines()
		{
		rec = [name: 'test_book_item', num: 5,
			text: '\r\ntest one\r\ntest two\r\ntest threee\r\n']
		SvcBook.GetData(rec)
		Assert(SvcBook.CheckRecord(rec, '') is: #())

		rec = [name: 'test_book_item2', num: 6,
			text: '\r\ntest one\r\ntest two\rtest threee\r\n']
		SvcBook.GetData(rec)
		Assert(SvcBook.CheckRecord(rec, '-') is: #())

		SvcBook.GetData(rec)
		Assert(SvcBook.CheckRecord(rec, ' ')
			is: #('Please ensure that this record does not have any invalid newlines.'))

		rec = [name: 'test_book_item3', num: 7,
			text: '\r\ntest one\r\ntest two\r\ntest threee\r<p>']
		SvcBook.GetData(rec)
		Assert(SvcBook.CheckRecord(rec, '+')
			equalsSet: #('unclosed tags: p @ 6',
				'Please ensure that this record does not have any invalid newlines.'))

		rec = [name: 'test_book_item4', num: 8,
			text: '\rtest one\r\ntest two\r\ntest threee\r\n']
		SvcBook.GetData(rec)
		Assert(SvcBook.CheckRecord(rec, '')
			is: #('Please ensure that this record does not have any invalid newlines.'))

		rec = [name: 'test_book_item5', num: 10, path: '/res',
			text: 'test one\r\ntest two\rtest threee\r\n']
		SvcBook.GetData(rec)
		Assert(SvcBook.CheckRecord(rec, '')
			is: #('Please ensure that this record does not have any invalid newlines.'))

		rec = [name: 'test_book_item6.gif', num: 11, path: '/res',
			text: 'test one\rtest two\rtest threee\r\n']
		SvcBook.GetData(rec)
		rec.name = rec.path $ '/' $ rec.name
		Assert(SvcBook.CheckRecord(rec, '') is: #())
		}

	Test_splitText()
		{
		m = SvcBook.SplitText

		// No "Order: " prefix, shouldn't throw errors
		m(rec = [path: '', text: 'no order'])
		Assert(rec.text is: 'no order')

		// path starts with /res and a image name, shouldn't split text
		m(rec = [path: '/res', text: 'Order: 5\n\nHas order', name: 'image.jpg'])
		Assert(rec.text is: 'Order: 5\n\nHas order')

		// Old record, uses \n\n, order is ""
		m(rec = [path: '', text: 'Order: \n\nHas order'])
		Assert(rec.text is: 'Has order')
		Assert(rec.order is: '')

		// Old record, uses \n\n
		m(rec = [path: '', text: 'Order: 5\n\nHas order'])
		Assert(rec.text is: 'Has order')
		Assert(rec.order is: 5)

		// Current record, uses \r\n
		m(rec = [path: '', text: 'Order: 15\n\nHas order'])
		Assert(rec.text is: 'Has order')
		Assert(rec.order is: 15)

		// Current record, uses \r\n, order is ""
		m(rec = [path: '', text: 'Order: \r\n\r\nHas order'])
		Assert(rec.text is: 'Has order')
		Assert(rec.order is: '')

		// Old record, uses \n\n
		m(rec = [path: '', text: 'Order: 115\n\nHas order twice Order: 19, end'])
		Assert(rec.text is: 'Has order twice Order: 19, end')
		Assert(rec.order is: 115)

		// Current record, uses \r\n, "Order: " exists later in text
		m(rec = [path: '', text: 'Order: 1115\n\nHas order twice Order: 19, end'])
		Assert(rec.text is: 'Has order twice Order: 19, end')
		Assert(rec.order is: 1115)

		// Current record, uses \r\n, "Order: " exists later in text
		m(rec = [path: '', text: 'Order: 11115\n\nOrder: 19\n\nHas order twice'])
		Assert(rec.text is: 'Order: 19\n\nHas order twice')
		Assert(rec.order is: 11115)
		}

	Test_FormatText()
		{
		m = SvcBook.FormatText

		// Testing text.Trim for base records
		text = ' \ttext to be trimmed \r\n '
		Assert(m(text, 'Normal') is: 'text to be trimmed')
		Assert(m(text, '/res/Image.GIF') is: text)
		Assert(m(text, '/res/Image.ico') is: text)
		Assert(m(text, '/res/Image.JPG') is: text)
		Assert(m(text, '/res/Javascript.js') is: text)
		Assert(m(text, '/res/Javascript.js.map') is: text)

		// Testing text.Trim for html records, no longer Detabs
		text = '<htmltag>' $ text
		Assert(m(text, 'Normal') is: '<htmltag> \ttext to be trimmed')
		Assert(m(text, '/res/Image.GIF') is: text)
		Assert(m(text, '/res/Image.ico') is: text)
		Assert(m(text, '/res/Image.JPG') is: text)
		Assert(m(text, '/res/Javascript.js') is: text)
		Assert(m(text, '/res/Javascript.js.map') is: text)
		}
	}
