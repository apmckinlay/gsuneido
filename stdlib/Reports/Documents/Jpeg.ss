// Copyright (C) 2014 Suneido Software Corp. All rights reserved worldwide.
class
	{
	New(.data)
		{
		if not String?(data)
			throw .errPrefix $ 'invalid jpg file'
		.info = .getImageInfo(data)
		}
	errPrefix: 'Image: Jpeg - ' // so ImageFormat can treat it as caught error

	// http://www.anttikupila.com/flash/
	//        getting-jpg-dimensions-with-as3-without-loading-the-entire-file/
	// Reference: http://www.garykessler.net/library/file_sigs.html
	getImageInfo(data)
		{
		if not data.Prefix?('\xff\xd8')
			throw .errPrefix $ 'invalid jpg file format'

		index = false
		width = 0
		height = 0
		pos = 0
		while false isnt cur = data.Match('\xff[^\x00]', pos)
			{
			curIndex = cur[0][0]
			marker = data[curIndex+1]
			if '\xc0\xc1\xc2\xc9'.Has?(marker)
				{
				if width < newWidth = .extractWidth(data, curIndex)
					{
					width = newWidth
					height = .extractHeight(data, curIndex)
					index = curIndex
					}
				}
			pos = curIndex + 2 +
				(.variableSize?(marker) ? .extractSegmentLength(data, curIndex) : 0)
			}

		if index is false
			throw .errPrefix $ 'cannot find SOI (Start Of Image)'
		return Object(:index, :width, :height)
		}

	variableSize?(marker)
		{
		return '\xc4\xdb\xdd\xda\xfe\xc0\xc1\xc2\xc9'.Has?(marker) or
			(marker >= '\xe0' and marker <= '\xef')
		}

	extractWidth(data, index)
		{
		return data[index + 7].Asc() * 256 + data[index + 8/*=offset*/].Asc()
		}

	extractHeight(data, index)
		{
		return data[index + 5].Asc() * 256 + data[index + 6/*=offset*/].Asc()
		}

	extractSegmentLength(data, index)
		{
		return data[index + 2].Asc() * 256 + data[index + 3/*=offset*/].Asc()
		}

	flagMap: (1: 'DeviceGray', 3: 'DeviceRGB', 4: 'DeviceCMYK')
	GetColorSpace()
		{
		if false is colorSpace = .flagMap.GetDefault(.getColorSpaceFlag(), false)
			throw .errPrefix $ 'unsupported file format'
		return colorSpace
		}

	getColorSpaceFlag()
		{
		return .data[.info.index + 9].Asc() /* = the offset of number of components*/
		}

	colorSpaceFlags: #(1, 3, 4)
	ColorSpaceValid?()
		{
		return .colorSpaceFlags.Has?(.getColorSpaceFlag())
		}

	GetWidth()
		{
		return .info.width
		}

	GetHeight()
		{
		return .info.height
		}

	RunWithCatch(block)
		{
		try
			{
			block()
			return true
			}
		catch (unused, 'Image: Jpeg - ')
			{
			return false
			}
		}

	InvalidExtension: 'Image: invalid extension used'
	ValidExtension?(filename)
		{
		return filename =~ `(?i)\.jpe?g\Z`
		}
	}