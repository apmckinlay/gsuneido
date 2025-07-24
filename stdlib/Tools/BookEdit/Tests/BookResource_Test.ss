// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	// Should detect: gif, jpg, png, bmp, emf, wmf, ico, cur, ttf, svg, js, map, css, afm
	Test_allBookResources()
		{ .tests() }

	// Should detect: gif, jpg, png, bmp, emf, wmf, ico, cur
	Test_imagesOnlyBookResources()
		{ .tests(imagesOnly?:) }

	// Should detect: gif, jpg, png, bmp, emf, wmf, ico, cur, ttf, svg, map, afm
	Test_readOnlyBookResources()
		{ .tests(readOnly?:) }

	// Should detect: gif, jpg, png, bmp, emf, wmf, ico, cur
	Test_readOnlyAndImagesOnlyBookResources()
		{ .tests(imagesOnly?:, readOnly?:) }

	tests(imagesOnly? = false, readOnly? = false)
		{
		_imagesOnly? = imagesOnly?
		_readOnly? = readOnly?

		// Non-book resources
		Assert(.fn('') is: false)
		Assert(.fn('file') is: false)
		Assert(.fn('/res') is: false)
		Assert(.fn('/res/gif') is: false)
		Assert(.fn('/res/file.txt') is: false)
		Assert(.fn('/res/file.gif.txt') is: false)

		// Book image resources
		Assert(.fn('/res/file.gif'))
		Assert(.fn('/res/file.JPG'))
		Assert(.fn('/res/file.png'))
		Assert(.fn('/res/file.Bmp'))
		Assert(.fn('/res/file.emf'))
		Assert(.fn('/res/file.wmf'))
		Assert(.fn('/res/file.ico'))
		Assert(.fn('/file.gif') is: false) // not in res
		Assert(.fn('/cores/result/file.gif') is: false) // not in res
		Assert(.fn('/res/file.gif.gif'))
		Assert(.fn('/res/nestedPath/file.cur'))

		// Other book resources (fonts, javascript, css, etc.)
		// Non-editable
		expectedNonEditable? = not imagesOnly?
		Assert(.fn('/res/file.ttf') is: expectedNonEditable?)
		Assert(.fn('/res/file.svg') is: expectedNonEditable?)
		Assert(.fn('/res/file.map') is: expectedNonEditable?)

		// Editable
		expectedEditable = expectedNonEditable? and not readOnly?
		Assert(.fn('/res/file.js') is: expectedEditable)
		Assert(.fn('/res/file.css') is: expectedEditable)
		}

	fn(name, _imagesOnly? = false, _readOnly? = false)
		{ return BookResource?(name, :imagesOnly?, :readOnly?) }
	}
