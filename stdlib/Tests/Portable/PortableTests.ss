// Copyright (C) 2014 Suneido Software Corp. All rights reserved worldwide.
/*
Runs test cases listed in a text file.
This is so the same test cases can be
shared between gSuneido, suneido.js, and from within Suneido.
*/
class
	{
	New(file = false, .dir = '../suneido_tests/')
		{
		.Run(file)
		}
	Run(file)
		{
		if file isnt false
			.runFile(file)
		else
			for file in Dir(.testdir $ '*.test')
				.runFile(file)
		}
	getter_testdir()
		{
		if .dir isnt false
			return .dir
		path = 'ptestdir.txt'
		for (i = 0; ; ++i)
			{
			if false isnt dir = GetFile(path)
				break
			if i > 8
				throw "can't find ptestdir.txt"
			path = '../' $ path
			}
		return .testdir = dir.Trim() // cache
		}
	runFile(file)
		{
		if false is s = GetFile(.testdir $ file)
			throw "can't get " $ .testdir $ file
		.file = file
		.scnr = Scanner(s)
		.eof = false
		.next(true)
		.run()
		}
	run()
		{
		while not .eof
			.run1()
		}
	run1()
		{
		.match('@')
		name = .scnr.Text()
		.match('IDENTIFIER', skip:)
		Print(.file $ ":", name $ ":", .comment)
		n = 0
		try
			test = Global('PT_' $ name)
		catch (unused, "can't find")
			{
			Print('\tMISSING TEST FIXTURE')
			test = false
			}
		ok = true
		while not .eof and .scnr.Text() isnt '@'
			{
			args = Object()
			str? = Object()
			do
				{
				if .scnr.Text() is '-'
					text = '-' $ .scnr.Next()
				else
					text = .scnr.Value()
				args.Add(text)
				str?.Add(.scnr.Type() is 'STRING')
				.next()
				if .scnr.Text() is ','
					.next(skip:)
				}
				while not .eof and .scnr.Type() isnt 'NEWLINE'
			args.str? = str?
			// here's the actual test
			if false is .test(test, args)
				{
				ok = false
				Print("\tFAILED: ", args)
				}
			else
				++n
			.next(skip:)
			}
		if test isnt false
			Print("\t" $ n $ " passed")
		return ok
		}
	test(test, args)
		{
		if test is false
			return true
		try
			return test(@args)
		catch (e)
			{
			Print('\t' $ e)
			return false
			}
		}
	match(expected, skip = false)
		{
		if .scnr.Text() isnt expected and .scnr.Type() isnt expected
			throw "expecting " $ Display(expected) $ " got " $ .scnr.Text()
		.next(skip)
		}
	eof: false
	next(skip = false)
		{
		.comment = ''
		nl = false
		do
			{
			if .scnr is .scnr.Next2()
				{
				.eof = true
				break
				}
			if .scnr.Type() is 'NEWLINE'
				if skip
					nl = true
				else
					break
			if .scnr.Type() is 'COMMENT' and not nl
				.comment = .scnr.Text()
			} while .scnr.Type() in ('WHITESPACE', 'COMMENT', 'NEWLINE')
		}
	}