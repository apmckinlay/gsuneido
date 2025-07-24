// Copyright (C) 2014 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_ToWindows()
		{
		Assert(Paths.ToWindows('') is: '')
		Assert(Paths.ToWindows('fred') is: 'fred')
		Assert(Paths.ToWindows('c:/work/eta') is: `c:\work\eta`)
		Assert(Paths.ToWindows(`c:\work\eta/tests`) is: `c:\work\eta\tests`)
		Assert(Paths.ToWindows(`\\server\fred`) is: `\\server\fred`)
		}

	Test_ToUnix()
		{
		Assert(Paths.ToUnix(`c:\work\eta/tests`) is: `c:/work/eta/tests`)
		}

	Test_ToStd()
		{
		Assert(Paths.ToStd('') is: '')
		Assert(Paths.ToStd('fred') is: 'fred')
		Assert(Paths.ToStd('c:/work/eta') is: 'c:/work/eta')
		Assert(Paths.ToStd(`c:\work\eta/tests`) is: `c:/work/eta/tests`)
		Assert(Paths.ToStd(`\\server\fred`) is: `//server/fred`)
		}

	Test_Basename()
		{
		mock = Mock(Paths)
		mock.When.Basename([anyArgs:]).CallThrough()

		// Windows
		mock.When.windows?().Return(true)
		for vol in #("", "c:", "c:/", `\\server\share\`)
			for dir in #("", "one/", "one/two/")
				for file in #("", "foo", "foo.bar")
					Assert(mock.Basename(vol $ dir $ file) is: file)
		Assert(mock.Basename(`\\server\fred/c:foo`) is: 'c:foo')
		Assert(mock.Basename(`\\server\fred`) is: '')
		Assert(mock.Basename(`c:\\server\fred/foo`) is: 'foo')

		// Linux/Unix
		mock.When.windows?().Return(false)
		Assert(mock.Basename(`\\server\fred`) is: 'fred')
		Assert(mock.Basename(`c:\\server\fred/foo`) is: 'foo')
		}

	Test_ToLocal()
		{
		mock = Mock(Paths)
		mock.When.ToLocal([anyArgs:]).CallThrough()
		path = `\test\this/path`

		// Windows
		mock.When.windows?().Return(true)
		Assert(mock.ToLocal(path) is: `\test\this\path`)

		// Linux/Unix
		mock.When.windows?().Return(false)
		Assert(mock.ToLocal(path) is: `/test/this/path`)
		}

	Test_ParentOf()
		{
		Assert(Paths.ParentOf('') is: '.')
		Assert(Paths.ParentOf('fred') is: '.')
		Assert(Paths.ParentOf('C:/work/test.txt') is: 'C:/work')
		Assert(Paths.ParentOf('C:\\work\\test.txt') is: 'C:\\work')
		Assert(Paths.ParentOf('\\server\fred') is: '\\server')
		Assert(Paths.ParentOf('C:\work/test\fred.txt') is: 'C:\work/test')
		Assert(Paths.ParentOf('c:/name.txt') is: 'c:')
		Assert(Paths.ParentOf('c:name.txt') is: 'c:')
		}

	Test_Combine()
		{
		Assert(Paths.Combine('a', 'b') is: 'a/b')
		Assert(Paths.Combine('a', '/b') is: 'a/b')
		Assert(Paths.Combine('a/', 'b') is: 'a/b')
		Assert(Paths.Combine('a/', '/b') is: 'a/b')
		Assert(Paths.Combine(`a\`, `\b`) is: 'a/b')

		Assert(Paths.Combine(`a`, `b`, 'c') is: 'a/b/c')
		Assert(Paths.Combine(`a\`, `b`, '/c') is: 'a/b/c')
		Assert(Paths.Combine(`/a\`, `b`, '/c/') is: '/a/b/c/')
		}

	Test_Absolute()
		{
		Assert(Paths.ToAbsolute('C:/work', '../bob/tests')
			is: Paths.ToLocal('C:/bob/tests'))
		Assert(Paths.ToAbsolute('/srv/work', '../updates')
			is: Paths.ToLocal('/srv/updates'))

		Assert(Paths.ToAbsolute('C:/work', './SamIAm') is:
			Paths.ToLocal('C:/work/SamIAm'))
		Assert(Paths.ToAbsolute('C:\work', `.\SamIAm`) is:
			Paths.ToLocal('C:/work/SamIAm'))

		Assert(Paths.ToAbsolute(`C:/work`, `C:/tests`) is:
			Paths.ToLocal('C:/tests'))
		Assert(Paths.ToAbsolute(`C:\work`, `C:\tests`) is:
			Paths.ToLocal('C:/tests'))
		Assert(Paths.ToAbsolute('/srv/work1', '/srv/work2') is:
			Paths.ToLocal('/srv/work2'))

		Assert(Paths.ToAbsolute(`c:\`, '..\bob')
			is: Paths.ToLocal('c:\bob'))
		}

	Test_ValidFileName?()
		{
		Assert(Paths.ValidFileName?('') is: false)
		Assert(Paths.ValidFileName?('hello'))
		Assert(Paths.ValidFileName?('hel lo'))				// space
		Assert(Paths.ValidFileName?('hel123lo'))			// number
		Assert(Paths.ValidFileName?('(hello)'))				// brackets
		Assert(Paths.ValidFileName?('[hello]'))				// brackets
		Assert(Paths.ValidFileName?('{hello}'))				// brackets
		Assert(Paths.ValidFileName?('(hello)'))				// brackets
		Assert(Paths.ValidFileName?("hello!!"))				// exclamation mark
		Assert(Paths.ValidFileName?("hello'"))				// single quote
		Assert(Paths.ValidFileName?('hello?') is: false)	// brackets
		Assert(Paths.ValidFileName?('<hello>') is: false)	// brackets
		Assert(Paths.ValidFileName?('hello\/') is: false)	// brackets
		Assert(Paths.ValidFileName?('hello:') is: false)	// brackets
		Assert(Paths.ValidFileName?('hello"') is false)		// double quote
		Assert(Paths.ValidFileName?('hello\t') is: false)	// tab
		Assert(Paths.ValidFileName?('hello	') is: false)	// tab
		Assert(Paths.ValidFileName?('hel\tlo') is: false)	// tab
		Assert(Paths.ValidFileName?('hel	lo') is: false)	// tab
		Assert(Paths.ValidFileName?('hello\n') is: false)	// newline (linux)
		Assert(Paths.ValidFileName?('hello\r') is: false)	// carriage return
		Assert(Paths.ValidFileName?('hello\r\n') is: false)	// newline (windows)
		Assert(Paths.ValidFileName?('hello\n\r') is: false)
		Assert(Paths.ValidFileName?('hello.world'))
		Assert(Paths.ValidFileName?('hello.
world') is: false)											// newline
		Assert(Paths.ValidFileName?(`hello|world.pdf `) is: false)
		Assert(Paths.ValidFileName?(`hello||||world.pdf`) is: false)
		Assert(Paths.ValidFileName?(`hello||**world.pdf`) is: false)
		Assert(Paths.ValidFileName?(`helloworld.|pdf`) is: false)
		Assert(Paths.ValidFileName?(`hel\xa1loworld.pdf`) is: false) // control character
		Assert(Paths.ValidFileName?(`hel\x01loworld.pdf`) is: false) // control character
		}

	Test_ParseUNC()
		{
		// Invalid UNC paths
		Assert(Paths.ParseUNC(`/linux/path`) is: false)
		Assert(Paths.ParseUNC(`C:/`) is: false)
		Assert(Paths.ParseUNC(`//`) is: false)
		Assert(Paths.ParseUNC(`///share_name`) is: false)
		Assert(Paths.ParseUNC(`//server_name`) is: false)
		Assert(Paths.ParseUNC(`//server_name/`) is: false)
		Assert(Paths.ParseUNC(`//server_name//`) is: false)
		Assert(Paths.ParseUNC(`//server_name//share_name`) is: false)

		// Valid UNC paths
		uncOb = Paths.ParseUNC(`//server_name/share_name`)
		Assert(uncOb.server is: 'server_name')
		Assert(uncOb.share is: 'share_name')
		Assert(uncOb.file is: '')

		uncOb = Paths.ParseUNC(`\\server_name/share_name/file.txt`)
		Assert(uncOb.server is: 'server_name')
		Assert(uncOb.share is: 'share_name')
		Assert(uncOb.file is: 'file.txt')

		uncOb = Paths.ParseUNC(`//server_name\share_name\sub_dir1\file.txt`)
		Assert(uncOb.server is: 'server_name')
		Assert(uncOb.share is: 'share_name')
		Assert(uncOb.file is: 'sub_dir1/file.txt')

		uncOb = Paths.ParseUNC(`//server_name\share_name/sub_dir1/sub_dir2\file.txt`)
		Assert(uncOb.server is: 'server_name')
		Assert(uncOb.share is: 'share_name')
		Assert(uncOb.file is: 'sub_dir1/sub_dir2/file.txt')
		}

	Test_Equal?()
		{
		mock = Mock(Paths)
		mock.When.Equal?([anyArgs:]).CallThrough()

		// Windows
		mock.When.windows?().Return(true)
		Assert(mock.Equal?(`C:\this\is\a\path`, `C:\this\is\a\path`))
		Assert(mock.Equal?(`C:\this\is\a\path`, `c:\this\is\a\path`))
		Assert(mock.Equal?(`\\this\is\a\path`, `\\this\is\a\path`))
		Assert(mock.Equal?(`\\this\is\A\path`, `\\this\is\a\path`))
		Assert(mock.Equal?(`\\this\is\a\path`, `\\this\is\path`) is: false)
		Assert(mock.Equal?(`\\this\is\a/path`, `\\this\is\a\path`) is: false)
		Assert(mock.Equal?(`\\this\is\a\path`, `\\this\is\a\path\`) is: false)

		// Linux/Unix
		mock.When.windows?().Return(false)
		Assert(mock.Equal?('//this/is/a/path', '//this/is/a/path'))
		Assert(mock.Equal?('//this/is/A/path', '//this/is/a/path') is: false)
		Assert(mock.Equal?('//this/is/a/path', '//this/is/path') is: false)
		Assert(mock.Equal?('//this/is\a/path', '//this/is/a/path') is: false)
		Assert(mock.Equal?('//this/is/a/path', '//this/is/a/path/') is: false)
		}

	Test_Equivalent?()
		{
		mock = Mock(Paths)
		mock.When.Equivalent?([anyArgs:]).CallThrough()

		// Windows
		mock.When.windows?().Return(true)
		Assert(mock.Equivalent?(`C:\this\is\a\path`, `C:\this\is\a\path`))
		Assert(mock.Equivalent?(`C:\this\is\a\path`, `c:\this\is\a\path`))
		Assert(mock.Equivalent?(`\\this\is\a\path`, `\\this\is\a\path`))
		Assert(mock.Equivalent?(`\\this\is\A\path`, `\\this\is\a\path`))
		Assert(mock.Equivalent?(`\\this\is\a\path`, `\\this\is\path`) is: false)
		Assert(mock.Equivalent?(`\\this\is\a/path`, `\\this\is\a\path`))
		Assert(mock.Equivalent?(`\\this\is\a\path`, `\\this\is\a\path\`) is: false)

		// Linux/Unix
		mock.When.windows?().Return(false)
		Assert(mock.Equivalent?('//this/is/a/path', '//this/is/a/path'))
		Assert(mock.Equivalent?('//this/is/A/path', '//this/is/a/path') is: false)
		Assert(mock.Equivalent?('//this/is/a/path', '//this/is/path') is: false)
		Assert(mock.Equivalent?('//this/is\a/path', '//this/is/a/path'))
		Assert(mock.Equivalent?('//this/is/a/path', '//this/is/a/path/') is: false)
		}

	Test_Prefix?()
		{
		mock = Mock(Paths)
		mock.When.Prefix?([anyArgs:]).CallThrough()

		// Windows
		mock.When.windows?().Return(true)
		Assert(mock.Prefix?(`C:\this\is\a\path`, `C:\this\is`))
		Assert(mock.Prefix?(`C:\this\is\a\path`, `c:\this\is\`))
		Assert(mock.Prefix?(`\\this\is\a\path`, `\\this\is\`))
		Assert(mock.Prefix?(`\\this\is\A\path`, `\\this\is/`))
		Assert(mock.Prefix?(`\\this\is\a\path`, `\\this\is\path`) is: false)
		Assert(mock.Prefix?(`\\this\is\a\path`, `\\this\is\a\path`)) // Exact match
		Assert(mock.Prefix?(`\\this\is\a\path`, `\\this\is\a\path\`) is: false)

		// Linux/Unix
		mock.When.windows?().Return(false)
		Assert(mock.Prefix?('/this/is/a/path', '/this/is/a/path'))
		Assert(mock.Prefix?('/this/is/A/path', '/this/is/a/path') is: false)
		Assert(mock.Prefix?('/this/is/a/path', '/this/is/path') is: false)
		Assert(mock.Prefix?('/this/is\a/path', '/this/is/a/') is: false)
		Assert(mock.Prefix?('/this/is/a/path', '/this/is/a/'))
		Assert(mock.Prefix?('/this/is/a/path', '/this/is/a/path/') is: false)
		}

	Test_EnsureTrailingSlash()
		{
		mock = Mock(Paths)
		mock.When.EnsureTrailingSlash([anyArgs:]).CallThrough()

		Assert(mock.EnsureTrailingSlash(``) is: `/`)
		Assert(mock.EnsureTrailingSlash(``, noSlashWhenEmpty:) is: ``)
		Assert(mock.EnsureTrailingSlash(`/this/is/a/path`) is: `/this/is/a/path/`)
		Assert(mock.EnsureTrailingSlash(`/this/is/a/path/`) is: `/this/is/a/path/`)
		}
	}
