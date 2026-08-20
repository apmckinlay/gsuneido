// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		logo = .TempName()
		letterHeadOb = .letterhead()
		items = DrawControl.BuildItems(letterHeadOb)
		fmt = DrawControl.WrapItems(items)
		ob = DrawControl.ApplyJustify(fmt, letterHeadOb)
		Assert(CanvasToSvgConverter(ob, logo) is: .expect(logo))
		}

	sampleJpeg: "/9j/4AAQSkZJRgABAQAAAQABAAD/2wBDAAgGBgcGBQgHBwcJCQgKDBQNDAsLD" $
		"BkSEw8UHRofHh0aHBwgJC4nICIsIxwcKDcpLDAxNDQ0Hyc5PTgyPC4zNDL/wAALCAABAAEBARE" $
		"A/8QAHwAAAQUBAQEBAQEAAAAAAAAAAAECAwQFBgcICQoL/8QAtRAAAgEDAwIEAwUFBAQAAAF9A" $
		"QIDAAQRBRIhMUEGE1FhByJxFDKBkaEII0KxwRVS0fAkM2JyggkKFhcYGRolJicoKSo0NTY3ODk" $
		"6Q0RFRkdISUpTVFVWV1hZWmNkZWZnaGlqc3R1dnd4eXqDhIWGh4iJipKTlJWWl5iZmqKjpKWmp" $
		"6ipqrKztLW2t7i5usLDxMXGx8jJytLT1NXW19jZ2uHi4+Tl5ufo6erx8vP09fb3+Pn6/9oACAE" $
		"BAAA/APf6/9k="

	letterhead()
		{
		return Object(items: Object(
			#("CanvasEllipse", 152, 407, 441, 497),
			#("CanvasRoundRect", 164, 49, 434, 396, 90, 90),
			#("CanvasRect", 184, 76, 215, 369),
			Object("CanvasImage", .sampleJpeg, 222, 85, 425.25, 356),
			#("CanvasText", "dGVzdCB0ZXh0", 150, 415, 434, 517,
				#(name: "Arial", italic: false, weight:400, size: 15), encoded:,
				justify: "center"),
			#("CanvasLine", 158, 509, 422, 509)),
			justify: "Center", resources: #(), pre_printed: "").DeepCopy()
		}

	expect(logo)
		{
		'<div data_logo="' $ logo $ '" ' $
			'style="display:flex;justify-content:center;">' $
			'<svg style="width:100%;height:auto;max-width:329.8px;" ' $
				'viewBox="0 0 329.8 530.4" ' $
				'xmlns="http://www.w3.org/2000/svg">' $
				'<ellipse cx="166.0333333333334" cy="456.7333333333333" ' $
					'rx="163.7666666666667" ry="51" ' $
					'style="stroke-width: 1;fill: #ffffff;stroke: #000000" />\n' $
				'<rect height="393.2666666666667" rx="51" ry="51" ' $
					'style="stroke-width: 1;fill: #ffffff;stroke: #000000" ' $
					'width="306" x="15.86666666666667" y="0" />\n' $
				'<rect height="332.0666666666667" ' $
					'style="stroke-width: 1;fill: #ffffff;stroke: #000000" ' $
					'width="35.13333333333333" x="38.53333333333333" y="30.6" />\n' $
				'<image height="307.1333333333333" ' $
					'href="data:image/jpeg;base64,' $ .sampleJpeg $ '" ' $
					'preserveAspectRatio="none" width="230.35" x="81.6" y="40.8" />\n' $
				'<text style="font-size: 17pt;font-family: Helvetica;' $
					'font-weight: 400;fill: #000000" ' $
					'x="121.244" y="435.9026666666667">test text</text>\n' $
				'<line style="stroke-width: 1;stroke: #000000" ' $
					'x1="9.066666666666666" x2="308.2666666666667" ' $
					'y1="521.3333333333333" y2="521.3333333333333" />\n</svg></div>'
		}

	Test_justify()
		{
		fakeFmt = Object(t: Timestamp())
		fn = CanvasToSvgConverter.CanvasToSvgConverter_justify
		Assert(fn(DrawControl.ApplyJustify(fakeFmt, #(justify: #Left)), fakeFmt)
			is: 'flex-start')
		Assert(fn(DrawControl.ApplyJustify(fakeFmt, #(justify: #Center)), fakeFmt)
			is: 'center')
		Assert(fn(DrawControl.ApplyJustify(fakeFmt, #(justify: #Right)), fakeFmt)
			is: 'flex-end')
		}
	}
