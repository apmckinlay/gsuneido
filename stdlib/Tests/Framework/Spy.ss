// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
class
	{
	New(target)
		{
		.Target = .convertTarget(target)
		.setupInfo()
		.Id = Timestamp()
		.returns = Object()
		.callLogs = Object()
		.registerSpy()
		}

	convertTarget(target)
		{
		if String?(target)
			target = Global(target)

		.origTarget = target
		switch (Type(target))
			{
		case 'Method':
		case 'Function':
		case 'Class':
			target = .getMethod(target, #CallClass)
		case 'Instance':
			target = .getMethod(target, #Call)
		default:
			throw "can't SpyOn " $ Display(target)
			}
//		if Display(target).Has?('stdlib')
//			{
//			Print("Should not spy on stdlib")
//			StackTrace(1, skip: 4)
//			}

		Assert(Function?(target))
		return target
		}

	getMethod(target, method)
		{
		if not target.Method?(method)
			throw 'SpyOn target (' $ Display(target) $ ') needs ' $ method $ ' method'
		target = target[method]
		}

	setupInfo()
		{
		display = Display(.Target)
		name = Name(.Target)
		if name.Blank?()
			throw 'SpyOn cannot get target name from ' $ name

		note = display.AfterFirst(`/* `).BeforeFirst(` */`)
		.Name = name.BeforeFirst('.')
		.Paths = name.Split('.')[1..]
		.Lib = LibraryTags.RemoveTagFromName(note.BeforeLast(' '))
		.Method? = note.Suffix?('method')
		.Params = .Target.Params()
		if Type(.origTarget) is 'Class'
			{
			origName = Name(.origTarget)
			if origName isnt .Name
				{
				.Name = origName
				display = Display(.origTarget)
				note = display.AfterFirst(`/* `).BeforeFirst(` */`)
				.Lib = LibraryTags.RemoveTagFromName(note.BeforeFirst(' '))
				}
			}
		}

	registerSpy()
		{
		SpyManager().Register(this)
		}

	Call(locals)
		{
		.callLogs.Add(locals)
		if false isnt res = .returns.FindOne({ it.when is true or (it.when)(@locals) })
			return .buildReturn(res)
		return Object(action: 'call through')
		}

	buildReturn(returnOb)
		{
		if returnOb.values.Size() is 1
			return Object(action: returnOb.action, value: returnOb.values[0])
		Assert(returnOb.values.Size() greaterThan: returnOb.index)
		return Object(action: returnOb.action, value: returnOb.values[returnOb.index++])
		}

	Return(@args)
		{
		values = args.Values(list:)
		when = args.GetDefault(#when, true)
		.returns.Add(Object(:values, :when, index: 0, action: #return))
		return this
		}

	ClearAndReturn(@args)
		{
		.returns = Object()
		.Return(@args)
		}

	Throw(@args)
		{
		values = args.Values(list:)
		when = args.GetDefault(#when, true)
		.returns.Add(Object(:values, :when, index: 0, action: #throw))
		return this
		}

	CallLogs()
		{
		return .callLogs
		}

	Close()
		{
		SpyManager().RemoveOne(this)
		.Delete(all:)
		}
	}
