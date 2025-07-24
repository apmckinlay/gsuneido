// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
function (name, book = 'imagebook', path = '/res')
	{
	if false isnt code = IconFont().MapToCharCode(name)
		return ImageFont(code.char, code.font)

	if false isnt image = Query1Cached(book, :path, :name)
		return Image(image.text)

	// fallback
	SuneidoLog('ERROR: (CAUGHT) Nonexistent Image Resource', calls:,
		params: Object(:name, :book, :path),
		caughtMsg: 'fall back to use triangle-warning.emf')
	return Image(
		Query1Cached('imagebook', path: '/res', name: 'triangle-warning.emf').text)
	}