// Copyright (C) 2004 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(helpbook)
		{
		index = helpbook $ "HowToIndex"
		try Database("destroy " $ index)
		Database("create " $ index $ " (name, howtos) key(name)")
		howtos = .BuildHowTos(helpbook)
		bad = #()
		if helpbook.Suffix?('Help')
			bad = .checknames(howtos, helpbook[.. -4]) /* = size of the word help*/
		for name in howtos.Members()
			QueryOutput(index, [:name, howtos: howtos[name]])
		return bad
		}

	BuildHowTos(helpbook)
		{
		howtos = Object().Set_default(Object())
		QueryApply(helpbook)
			{ |x|
			.addone(howtos, x)
			}
		return howtos
		}

	addone(howtos, x)
		{
		startIdx = 13
		endIdx  =  -4
		x.text.Replace('<!-- option: .+? -->')
			{ |s|
			if x.name in ('How do I use a form screen?', 'How do I use a list screen?',
				'How do I use a list+form screen?')
				howtos[s[startIdx .. endIdx]].Add(x.path $ '/' $ x.name, at: 0)
			else
				howtos[s[startIdx .. endIdx]].Add(x.path $ '/' $ x.name)
			}
		}

	checknames(howtos, book)
		{
		bad = Object()
		for name in howtos.Members()
			{
			if name.Has?('/')
				{
				path = name.BeforeLast('/')
				item = name.AfterLast('/')
				if 1 isnt QueryCount(book $ ' where name = ' $ Display(item) $
					' and path is ' $ Display(path))
					bad.Add(name)
				}
			else if 1 isnt QueryCount(book $ ' where name = ' $ Display(name) $
				' and not path.Has?("Reporter Reports")' $
				' and not path.Has?("Reporter Forms")')
				bad.Add(name)
			}
		return bad
		}
	}
