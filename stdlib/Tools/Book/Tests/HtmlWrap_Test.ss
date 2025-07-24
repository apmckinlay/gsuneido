// Copyright (C) 2019 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_embed()
		{
		embed = HtmlWrap.Embed
		Assert(embed('') is: '')
		Assert(embed('<html><head></head><body></body></html>')
				 is: '<html><head></head><body></body></html>')

		book = .MakeBook()
		rec = .MakeBookRecord(book, '1 1', path: '/res')
		QueryDo('update ' $ book $ ' where name is ' $ Display(rec.name) $
			' set name = "test.png"')
		Assert(embed('<html><img src="suneido:/' $ book $ '/res/test.png" /></html>')
				 is: '<html><img src="data:image/png;base64,MSAx" /></html>')

		QueryDo('update ' $ book $ ' where name is "test.png" ' $
			' set name = "test.jpg"')
		Assert(embed('<html><img src="suneido:/' $ book $ '/res/test.jpg" /><div>' $
			'x'.Repeat(200) $ '</div></html>')
				 is: '<html><img src="data:image/jpeg;base64,MSAx" /><div>' $
					'x'.Repeat(200) $ '</div></html>')

		QueryDo('update ' $ book $ ' where name is "test.jpg" ' $
			' set name = "test.xxx"')
		Assert(embed('<html><img src="suneido:/' $ book $ '/res/test.xxx" /></html>')
				 is: '<html><img src="suneido:/' $ book $ '/res/test.xxx" /></html>')

		}
	}
