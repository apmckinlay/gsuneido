// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_imageSizeOb()
		{
		_report = FakeObject(
			GetImageSize: function (unused)	{ return Object(width: 1, height: 1) }
			)
		mock = Mock(ImageFormat)
		mock.When.imageSizeOb([anyArgs:]).CallThrough()
		mock.ImageFormat_data = ''
		mock.ImageFormat_validImageSize? = ''

		// Initial run, no error occurs, image size is valid
		Assert(mock.imageSizeOb('') isnt: false)
		Assert(mock.ImageFormat_validImageSize?)

		// Image previously calculated, was valid, get size again
		mock.ImageFormat_validImageSize? = true
		Assert(mock.imageSizeOb('') isnt: false)

		// Image previously calculated, was invalid, do not get size again
		mock.ImageFormat_validImageSize? = false
		Assert(mock.imageSizeOb('') is: false)

		// Image previously calculated, was invalid, data is different however
		// get the new size. Do not touch the valid flag
		mock.ImageFormat_data = 'origValue'
		Assert(mock.imageSizeOb('diffValue') is: Object(width: 1, height: 1))
		Assert(mock.ImageFormat_validImageSize? is: false)
		}

	Test_dataString()
		{
		c = ImageFormat
			{
			ImageFormat_data: 'test2.file'
			}
		fn = c.ImageFormat_dataString
		data, textOnly? = fn('')
		Assert(data endsWith: 'test2.file')
		Assert(textOnly?)

		data, textOnly? = fn('test.file')
		Assert(data endsWith: 'test.file')
		Assert(textOnly?)

		data, textOnly? = fn('\xfc\xff')
		Assert(data is: '\xfc\xff')
		Assert(textOnly? is: false)

		c = ImageFormat
			{
			ImageFormat_data: ''
			}
		fn = c.ImageFormat_dataString
		data, textOnly? = fn('')
		Assert(data is: '')
		Assert(textOnly?)

		fn = c.ImageFormat_dataString
		data, textOnly? = fn('test.file')
		Assert(data endsWith: 'test.file')
		Assert(textOnly?)

		fn = c.ImageFormat_dataString
		data, textOnly? = fn('test.pdf')
		Assert(data endsWith: 'test.pdf')
		Assert(textOnly?)

		fn = c.ImageFormat_dataString
		data, textOnly? = fn('test.pdf Axon Label: test')
		Assert(data endsWith: 'test.pdf')
		Assert(textOnly?)

		data, textOnly? = fn(GetAppTempPath() $ 'test.jpg')
		Assert(data endsWith: 'test.jpg')
		Assert(textOnly?)

		data, textOnly? = fn(GetAppTempPath() $ 'test.jpg Axon Label: test')
		Assert(data endsWith: 'test.jpg')
		Assert(textOnly?)
		}
	}
