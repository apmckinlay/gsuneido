// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_encode()
		{
		fn = SuBookSendEmail.SuBookSendEmail_encode
		Assert(fn('') is: '')
		Assert(fn('hello') is: 'hello')
		Assert(fn('Hello !@#$') is: 'Hello !@#$')

		Assert(fn('Hello \xe0') is: '=?UTF-8?B?SGVsbG8gw6A?=')

		// includes all french and mexico letters
		Assert(fn(
			"ABCDEFGHIJKLMNOPQRSTUVWXYZ" $
			"\xe0\xe2\xe6\xe7\xe9\xe8\xea\xeb\xee\xef\xf4\x9c\xf9\xfb" $
			"\xfc\xff\xe1\xed\xf3\xfa\xf1\xc1\xc9\xcd\xd3\xda\xdc")
			is: "=?UTF-8?B?QUJDREVGR0hJSktMTU5PUFFSU1RVVldYWVrDoMOiw6b" $
				"Dp8Opw6jDqsOrw67Dr8O0xZPDucO7w7zDv8Ohw63Ds8O6w7HDgcOJw43Dk8Oaw5w?=")
		}
	}
