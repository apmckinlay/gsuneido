// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_no_photos()
		{
		cl = ImagesList { ImagesList_getImageFiles(unused) { return Object() }}
		Assert(cl(#(query: '')) is: "<p>no photos</p>")
		}

	Test_image_list()
		{
		cl = ImagesList
			{
			ImagesList_getImageFiles(unused)
				{
				return Object(`img_1.jpg`, `img_2.jpg`, `img_3.jpg`)
				}
			}
		expected = `<html><body bgcolor="lightblue"><p align="center" ` $
			`style="margin-top: 0; margin-bottom: 0"><b>Tmp</b></p>` $
			`<p align="center" style="margin-top: 0; margin-bottom: .5em">` $
			`<b>&lt; prev</b>&nbsp;&nbsp;` $
			`<a href="ImagesPage?c:/tmp/img_2.jpg" target="_top">` $
				`<b>next &gt;</b></a>` $
			`</p><font size=2 face=Arial><b>Img_1</b><br />` $
			`<a href="ImagesPage?c:/tmp/img_2.jpg" target="_top">Img_2</a><br />` $
			`<a href="ImagesPage?c:/tmp/img_3.jpg" target="_top">Img_3</a><br />` $
			`</font></body></html>`
		actual = cl(#(query: `c:/tmp/`))
		Assert(actual.Tr('\n') is: expected)
		}
	}