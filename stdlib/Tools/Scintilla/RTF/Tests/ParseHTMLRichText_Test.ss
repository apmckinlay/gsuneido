// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Setup()
		{
		.report = class
			{
			Construct(item)
				{
				size = item[1].Size() * item.font.size
				itm = FakeObject(GetSize: Object(w: size))
				return itm
				}
			GetFont()
				{
				return #(name: 'Arial', size: 10, weight: 'regular')
				}
			}
		.font = #(name: 'Arial', size: 10, weight: 'regular')
		}
	Test_findSplitPoint()
		{
		_report = Mock()
		_report.When.PlainText?().Return(false)
		findTextFit = ParseHTMLRichText.ParseHTMLRichText_findTextFit
		str = "A B C D E"
		w = 10800
		currentSize = 0
		Assert(findTextFit(str, currentSize, w, .font, .report) is: "A B C D E")

		currentSize = 10750
		Assert(findTextFit(str, currentSize, w, .font, .report) is: "A B C ")

		currentSize = 10799
		Assert(findTextFit(str, currentSize, w, .font, .report) is: "")

		currentSize = 6496
		str = ' Fuel Download be replaced with the '
		Assert(findTextFit(str, currentSize, w, .font, .report)
			is: ' Fuel Download be replaced with the ')

		currentSize = 10802
		str = 'Testing'
		Assert(findTextFit(str, currentSize, w, .font, .report) is: '')

		currentSize = 10750
		Assert(findTextFit(str, currentSize, w, .font, .report) is: '')

		str = 'Testing again'
		Assert(findTextFit(str, currentSize, w, .font, .report) is: '')

		str = 'bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb' $
			'bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb' $
			'bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb' $
			'bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb' $
			'bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb' $
			'bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb'
		w = 3500
		Assert(findTextFit(str, currentSize, w, .font, .report) is: '')

		currentSize = 0
		Assert(findTextFit(str, currentSize, w, .font, .report)
			is: 'bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb' $
			'bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb' $
			'bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb' $
			'bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb' $
			'bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb' $
			'bbbbbbbbbbbbbbbbbbbbbb')
		}

	testFormats: #(
		#(text: '', result: #()),
		#(text: '<span style="font-weight:normal;font-style:normal"><br /></span>',
			result: #(())),
		#(text: '<span style="font-weight:normal;font-style:normal"><br /><br /></span>',
			result: #((), ())),
		#(text: '<span style="font-weight:normal;font-style:normal">This is text</span>',
			result: #(#(
				#('Text', 'This is text',
					font: #(name: 'Arial', size: 10, weight: 'regular'))))),
		#(text: '<span style="font-weight:normal;font-style:normal">' $
			'This is normal text</span>' $
			'<span style="font-weight:bold;font-style:normal">This is bold text</span>',
			result: #(#(
				#('Text', 'This is normal text',
					font: #(name: "Arial", size: 10, weight: 'regular')),
				#('Text', 'This is bold text',
					font: #(name: "Arial", size: 10, weight: 'bold'))))),
		#(text: '<span style="font-weight:normal;font-style:normal">' $
			'This is normal text that is designed to go off the end of the page. ' $
			'It will just be long enough that the very end of it will be on a ' $
			'brand new page this text should be on its own line</span>',
			result: #(
				#(#('Text',
					'This is normal text that is designed to go off the end of ' $
					'the page. It will just be long enough that the very end of ' $
					'it will be on a brand new page ',
						font: #(name: "Arial", size: 10, weight: 'regular'))),
				#(#('Text', 'this text should be on its own line',
					font: #(name: "Arial", size: 10, weight: 'regular')))
				)),
		#(text: '<span style="font-weight:normal;font-style:normal">' $
			'This is the first line<br />This is the second line</span>',
			result: #(
				#(#('Text', 'This is the first line',
					font: #(name: "Arial", size: 10, weight: 'regular'))),
				#(#('Text', 'This is the second line',
					font: #(name: "Arial", size: 10, weight: 'regular')))
				)),
		#(text: '<span style="font-weight:normal;font-style:normal">' $
			'This is the first line<br /></span>' $
			'<span style="font-weight:bold;font-style:normal">' $
			'This is the second line<br /></span>'
			result: #(
				#(#('Text', 'This is the first line',
					font: #(name: "Arial", size: 10, weight: 'regular'))),
				#(#('Text', 'This is the second line',
					font: #(name: "Arial", size: 10, weight: 'bold')))
				)),
		#(text:'This is normal text',
			result: #(#(#('Text', 'This is normal text',
				font: #(name: "Arial", size: 10, weight: 'regular')))))
		#(text: 'This is a mix of normal text ' $
			'<span style="font-weight:bold;font-style:normal">and formatted text</span>',
			result: #(
				#(#('Text', 'This is a mix of normal text ',
					font: #(name: "Arial", size: 10, weight: 'regular'))
				#('Text', 'and formatted text',
					font: #(name: 'Arial', size: 10, weight: 'bold')))
				)),
		#(text: '<span style="font-weight:normal;font-style:normal">' $
			'This is normal text that is designed to go off the end of the page. ' $
			'It will just be long enough that the very end of it will be on a ' $
			'brand new page this text should be on its own line<br /><br />Testing a ' $
			'glitch with multiple lines that go over a page. Need to ensure that ' $
			'this text is going to be fairly large</span>',
			result: #(
				#(#('Text', 'This is normal text that is designed to go off the ' $
					'end of the page. It will just be long enough that the very ' $
					'end of it will be on a brand new page ',
						font: #(name: "Arial", size: 10, weight: 'regular'))),
				#(#('Text', 'this text should be on its own line',
					font: #(name: "Arial", size: 10, weight: 'regular'))),
				#(),
				#(#('Text', 'Testing a glitch with multiple lines that go over a ' $
					'page. Need to ensure that this text is going to be fairly ' $
					'large',
						font: #(name: "Arial", size: 10, weight: 'regular')))
				))
		#(text: '<span style="font-weight:normal;font-style:normal">' $
			'&amp;amp;<br />' $
			'&amp;lt;<br /></span>',
			result: #(
				#(#('Text', '&amp;',
					font: #(name: "Arial", size: 10, weight: 'regular')))
				#(#('Text', '&lt;',
					font: #(name: "Arial", size: 10, weight: 'regular')))
				))
		#(text: '<span style="">first\nsecond\n\nthird\n\n</span>',
			result: #(
				((Text, 'first', font: #(name: "Arial", size: 10, weight: 'regular')))
				((Text, 'second', font: #(name: "Arial", size: 10, weight: 'regular')))
				()
				((Text, 'third', font: #(name: "Arial", size: 10, weight: 'regular')))
				()))
		)
	Test_GetFormats()
		{
		_report = Mock()
		_report.When.PlainText?().Return(false)
		getFormats = ParseHTMLRichText.GetFormats
		w = 1500
		fmt = Mock()
		fmt.GetSize()

		for format in .testFormats
			Assert(getFormats(format.text, w, .report) is: format.result)
		}
	Test_lineLimit()
		{
		text = '<span style="font-weight:normal;font-style:normal">' $
			'line1\r\nline2\r\nline3\r\nline4\r\nline5\r\nline6\r\nline7\r\nline8\r\n' $
			'line9\r\nline10</span>'
		_report = Mock()
		_report.When.PlainText?().Return(false)
		w = 1500

		Assert(ParseHTMLRichText.GetFormats(text, w, .report, lineLimit: 5) is:
			#(#(#('Text', 'line1', font: #(name: "Arial", size: 10, weight: 'regular'))),
			#(#('Text', 'line2', font: #(name: "Arial", size: 10, weight: 'regular'))),
			#(#('Text', 'line3', font: #(name: "Arial", size: 10, weight: 'regular'))),
			#(#('Text', 'line4', font: #(name: "Arial", size: 10, weight: 'regular'))),
			#(#('Text', 'line5', font: #(name: "Arial", size: 10, weight: 'regular'))),
			#(#('Text', '...', font: #(name: "Arial", size: 10, weight: 'regular')))))

		}
	Test_detab()
		{
		_report = Mock()
		_report.When.PlainText?().Return(false)
		getFormats = ParseHTMLRichText.GetFormats

		// Testing Detab
		w = 10800
		str = '<span style="font-weight:normal;font-style:normal">' $
			'Transport		4500	75<br />' $
			'FT				1450	48<br /><br />' $
			'PO&apos;s		1900	32<br /><br />' $
			'Sat Inter		1450	48<br /><br />' $
			'Tickets		1500	50<br /><br />' $
			'Pings			825		100<br /><br />' $
			'Total			11625	353</span>'
		result = getFormats(str, w, .report)
		Assert(result is:
			#(
				#(#('Text', 'Transport       4500    75',
					font: #(name: "Arial", size: 10, weight: 'regular'))),
				#(#('Text', 'FT              1450    48',
					font: #(name: "Arial", size: 10, weight: 'regular'))),
				#(),
				#(#('Text', "PO's        1900    32",
					font: #(name: "Arial", size: 10, weight: 'regular'))),
				#(),
				#(#('Text', 'Sat Inter       1450    48',
					font: #(name: "Arial", size: 10, weight: 'regular'))),
				#(),
				#(#('Text', 'Tickets     1500    50',
					font: #(name: "Arial", size: 10, weight: 'regular'))),
				#(),
				#(#('Text', 'Pings           825     100',
					font: #(name: "Arial", size: 10, weight: 'regular'))),
				#(),
				#(#('Text', 'Total           11625   353',
					font: #(name: "Arial", size: 10, weight: 'regular')))
			))
		}
	}