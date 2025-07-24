// Copyright (C) 2008 Suneido Software Corp. All rights reserved worldwide.
// filter can be a list of filters that are AND'ed together
// [0] is company
// [1] is timestamp
// [2] is heap
// [3] is user (session id)
// [4] is message (in quotes)
// in the results [0] is message, [1] is count
class
	{
	CallClass(file, filter = function (unused) { true }, outfile = false)
		{
		if outfile is false
			return .process(file, filter, Print)
		else
			File(outfile, "w")
				{|outf|
				return .process(file, filter, { outf.Writeline(it) })
				}
		}
	process(file, filter, print)
		{
		c = Object().Set_default(0)
		.Foreach(file, filter)
			{|x|
			++c[x[4]]
			if c.Size() > 10000
				throw "too many"
			}
		f = Object()
		total = 0
		for m in c.Members().Sort!()
			{
			print(c[m].Pad(6, ' ') $ ' ' $ m)
			f.Add([m, c[m]])
			total += c[m]
			}
		print('='.Repeat(60))
		f.Sort!({|x, y| x[1] > y[1] }) // reverse
		for x in f
			print(x[1].Pad(6, ' ') $ ' ' $ x[0])
		print("")
		print('TOTAL: ' $ total)
		return f
		}
	Penny(file, outfile = false)
		{
		//axon = function (x) { x[0].Has?('axon') }
		not_axon = function (x) { not x[0].Has?('axon') }
		help = function (x) { x[4] =~ `^.(ETAHelp/)|(How Do I)|(Book Search)` }
		BookLogProcessing(file, [not_axon, help], outfile)
		}
	Xp(file)
		{
		c = Object().Set_default(#())
		BookLogProcessing.Foreach(file, {|x| x[4].Has?('Login WinXP') })
			{|x|
			cus = x[0].Replace(' \d$', '') // combine multiple companies
			ses = x[3]
			c[cus][ses] = true
			}
		for cus in c.Members()
			Print(c[cus].Size() $ '\t' $ cus)
		}
	Users(file, filter = function (unused) { true })
		{
		companies = Object()
		users = Object()
		.Foreach(file, filter)
			{|x|
			companies[x[0]] = true
			users[x[0] $ ':' $ x[2].BeforeFirst('@')] = true
			}
		Print(companies.Size(), "companies", users.Size(), "users")
		return [companies.Members().Sort!(), users.Members().Sort!()]
		}
	BadLines(file, filter = function (unused) { true })
		{
		.Foreach(file, filter)
			{|x, line|
			if Date(x[1]) is false
				Print(line $ ": " $ Display(x))
			}
		}
	Foreach(file, filter, block)
		{
		if not Object?(filter)
			filter = [filter]
		line = 0
		File(file)
			{ |f|
			while false isnt s = f.Readline()
				{
				++line
				try
					{
					y = s.SplitOnFirst(', ')
					z = y[1].Split('\t')
					x = Object(y[0]).Add(@z)
					include? = true
					for filt in filter
						include? = include? and filt(x)
					if include?
						block(x, line: line)
					}
				catch (err)
					{
					Print("ERROR", err)
					Print("line", line $ ':', s)
					}
				}
			}
		Print(lines: line)
		}
	}