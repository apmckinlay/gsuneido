// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_New()
		{
		Assert({ File('') } throws: "open")
		Assert({ File("doesnt exist") } throws: "open")
		fileName = .MakeFile()
		File(fileName, "w") { |f/*unused*/| }
		File(fileName) { |f/*unused*/| }
		File(fileName).Close()
		}
	Test_ReadWrite()
		{
		fileName = .MakeFile()
		File(fileName, "w")
			{ |f|
			f.Write("hello world\nhow are you")
			}
		File(fileName)
			{ |f|
			Assert(f.Read() is: "hello world\nhow are you")
			Assert(f.Read() is: false)
			}
		File(fileName)
			{ |f|
			Assert(f.Read(5) is: "hello")
			Assert(f.Read() is: " world\nhow are you")
			Assert(f.Read() is: false)
			}
		}
	Test_ReadWriteline()
		{
		fileName = .MakeFile()
		File(fileName, "w")
			{ |f|
			f.Writeline("hello world")
			f.Writeline("how are you")
			}
		File(fileName)
			{ |f|
			Assert(f.Readline() is: "hello world")
			Assert(f.Readline() is: "how are you")
			Assert(f.Readline() is: false)
			}
		}
	Test_TellSeek()
		{
		fileName = .MakeFile()
		File(fileName, "w")
			{ |f|
			f.Writeline("hello world")
			f.Writeline("how are you")
			}
		File(fileName)
			{ |f|
			Assert(f.Readline() is: "hello world")
			offset = f.Tell()
			Assert(f.Readline() is: "how are you")
			Assert(f.Readline() is: false)
			f.Seek(offset)
			Assert(f.Readline() is: "how are you")
			}
		}
	Test_Flush()
		{
		fileName = .MakeFile()
		File(fileName, "w")
			{ |f|
			f.Write("hello")
			f.Flush()
			File(fileName)
				{ |f| after = f.Read() }
			Assert(after is: "hello")
			}
		}
	Test_bad_method()
		{
		Assert({ File.BadMethod() } throws: 'method not found')
		fileName = .MakeFile()
		PutFile(fileName, '')
		Assert({ File(fileName) { it.BadMethod() }} throws: 'method not found')
		}
	Test_Readline()
		{
		test = {|data, result|
			fileName = .MakeFile(data)
			Assert(.lines(fileName) is: result)
			}
		test('', #())
		test('one', #(one))
		test('\r\n', #(''))
		test('one\n', #(one))
		test('one\r\n', #(one))
		test('\r', #(''))
		test('\r\r\n', #(''))
		test('one\r', #('one'))
		test('one\r\r\n', #('one'))
		test('one\rtwo', #('one\rtwo'))
		test('one\rONE\ntwo\r\nthree\r\r\nfour',
			#('one\rONE', 'two', 'three', 'four'))
		test('one\n\ntwo', #(one, '', two))
		}
	lines(fileName)
		{
		lines = Object()
		File(fileName)
			{|f|
			while false isnt line = f.Readline()
				lines.Add(line)
			return lines
			}
		}
	Test_Writeline()
		{
		fileName = .MakeFile()
		File(fileName, "w") // create
			{|f|
			Assert(f.Tell() is: 0)
			f.Writeline("hello")
			Assert(f.Tell() is: 7)
			f.Writeline("world")
			Assert(f.Tell() is: 14)
			}
		Assert(FileSize(fileName) is: 14)
		Assert(GetFile(fileName) is: "hello\r\nworld\r\n")

		fileName = .MakeFile()
		File(fileName, "a") // append to new file (= create)
			{|f|
			f.Writeline("hello")
			f.Writeline("world")
			}
		Assert(FileSize(fileName) is: 14)
		Assert(GetFile(fileName) is: "hello\r\nworld\r\n")
		File(fileName, "a") // append to existing file
			{|f|
			f.Writeline("again")
			}
		Assert(FileSize(fileName) is: 21)
		Assert(GetFile(fileName) is: "hello\r\nworld\r\nagain\r\n")
		}
	Test_SeekWrite()
		{
		fileName = .MakeFile()
		File(fileName, "w") // create
			{|f|
			f.Write("hello world")
			f.Seek(0)
			f.Write("HELLO")
			}
		Assert(GetFile(fileName) is: "HELLO world")
		Assert(FileSize(fileName) is: 11)
		}
	Test_FileSize()
		{
		Assert({ FileSize("_non_existent_") } throws: "does not exist")

		fileName = .MakeFile("")
		Assert(FileSize(fileName) is: 0)

		PutFile(fileName, "hello world")
		Assert(FileSize(fileName) is: 11)
		}
	Test_Size()
		{
		file = .MakeFile()
		File(file, "w")
			{|f|
			Assert(f.Size() is: 0)
			f.Write("x")
			Assert(f.Size() is: 1)
			}
		File(file, "a")
			{|f|
			Assert(f.Size() is: 1)
			f.Write("x")
			Assert(f.Size() is: 2)
			}
		}
	}
