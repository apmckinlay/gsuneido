// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Setup()
		{
		.cl = CheckFileName
			{
			MaxAllowedFileNameChars: 20
			}
		}

	Test_Valid?()
		{
		valid? = .cl.ValidWithPath?
		Assert(valid?('') is: false)
		Assert(valid?(`test.txt`))
		Assert(valid?(`test`))
		Assert(valid?(`.txt`) is: false)
		Assert(valid?(`test | more`) is: false)
		Assert(valid?('hello'))
		Assert(valid?('hel lo'))				// space
		Assert(valid?('hel123lo'))				// number
		Assert(valid?('(hello)'))				// brackets
		Assert(valid?('[hello]'))				// brackets
		Assert(valid?('{hello}'))				// brackets
		Assert(valid?('(hello)'))				// brackets
		Assert(valid?("hello!!"))				// exclamation mark
		Assert(valid?("hello'"))				// single quote
		Assert(valid?('hello?') is: false)		// question mark
		Assert(valid?('<hello>') is: false)		// angle brackets
		Assert(valid?('hello\/') is: false)		// slashes
		Assert(valid?('hello:') is: false)		// colon
		Assert(valid?('hello"') is: false)		// double quote
		Assert(valid?('hello\t') is: false)		// tab
		Assert(valid?('hello	') is: false)	// tab
		Assert(valid?('hel\tlo') is: false)		// tab
		Assert(valid?('hel	lo') is: false)		// tab
		Assert(valid?('hello\n') is: false)		// newline (linux)
		Assert(valid?('hello\r') is: false)		// carriage return
		Assert(valid?('hello\r\n') is: false)	// newline (windows)
		Assert(valid?('hello\n\r') is: false)
		Assert(valid?('hello.world'))
		Assert(valid?('hello.
world') is: false)								// newline
		Assert(valid?(`hello|world.pdf `) is: false)
		Assert(valid?(`hello||||world.pdf`) is: false)
		Assert(valid?(`hello||**world.pdf`) is: false)
		Assert(valid?(`helloworld.|pdf`) is: false)
		// control characters
		Assert(valid?('hel\x1aloworld.pdf') is: false)
		Assert(valid?('hel\x01loworld.pdf') is: false)
		}

	Test_ValidWithPath?()
		{
		if not Sys.Windows?()
			return
		valid? = .cl.ValidWithPath?
		Assert(valid?(`c:\work\test.txt`))
		Assert(valid?(`c:work.txt`))
		Assert(valid?(`c:\work\test:again`) is: false)
		Assert(valid?(`c:\work<>`) is: false)
		Assert(valid?(`c:\work\`) is: false)
		Assert(valid?(`c:\work\NUL`) is: false)
		Assert(valid?(`c:\work\con`) is: false)
		Assert(valid?(`c:\work\lpT5`) is: false)
		Assert(valid?(`c:\work\NUL.txt`) is: false)
		Assert(valid?(`c:\work\nul.txt`) is: false)
		Assert(valid?(`c:\work\NULlified.txt`))
		Assert(valid?(`c:\work\` $ 'a'.Repeat(21)) is: false)
		}

	Test_WithErrorMsg()
		{
		func = .cl.WithErrorMsg
		Assert(func('test.txt') is: '')
		Assert(func('C:/work/test.txt', withPath?:) is: '')
		Assert(func('') is: 'File Name cannot be blank or just an extension')
		Assert(func('.txt') is: 'File Name cannot be blank or just an extension')
		Assert(func('C:/work/.txt', withPath?:)
			is: 'File Name cannot be blank or just an extension')

		Assert(func('bob<>.txt') is: CheckFileName.InvalidCharsDisplay)
		Assert(func('C:/work/bobby<hill>.txt', withPath?:)
			is: CheckFileName.InvalidCharsDisplay)

		Assert(func('bob:.txt') is: CheckFileName.InvalidCharsDisplay)
		Assert(func('C:/work/bobbyhill:.txt', withPath?:)
			is: CheckFileName.InvalidCharsDisplay)

		Assert(func('CON.txt')
			is: 'File Name cannot be a reserved Windows File Name')
		Assert(func('CONnected.txt') is: '')
		Assert(func('C:/work/LPT1.txt', withPath?:)
			is: 'File Name cannot be a reserved Windows File Name')
		Assert(func('C:/work/LPT1234.txt', withPath?:) is: '')

		// we are allowed to prefix folders with period
		Assert(func('.txt', isFolder?:) is: '')
		// file name is too long
		Assert(.cl.WithErrorMsg('a'.Repeat(21)) is: .cl.MaxAllowedCharsMsg)
		}
	}