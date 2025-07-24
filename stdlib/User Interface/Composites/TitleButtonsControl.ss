// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
PassthruController
	{
	New(args)
		{
		super(.controls(args))
		.insertButtons()
		}
	controls(args)
		{
		.args = args
		.border = .args.Member?(#border) ? .args.border : 6
		.postControls = Object()
		controls = ['Horz']
		Plugins().ForeachContribution('TitleButtons', 'button', showErrors:, sort:)
			{|x|
			if not x.Member?('condition') or true is (x.condition)(args)
				{
				if x.Member?('postCondition')
					.postControls.Add(Object(pos: controls.Size() - 1,
						control: x.control, condition: x.postCondition))
				else
					controls.Add(Object('Border', x.control, .border))
				}
			}
		return controls
		}
	Notes_Title()
		{
		return .args.GetDefault('Notes_Title', .args[0])
		}
	insertButtons()
		{
		added = 0
		for ob in .postControls
			if ((ob.condition)(this))
				{
				.Horz.Insert(ob.pos + added, Object('Border', ob.control, .border))
				++added
				}
		}
	}
