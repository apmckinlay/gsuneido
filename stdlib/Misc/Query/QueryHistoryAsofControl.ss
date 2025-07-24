// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Title: 'Query History Asof'
	CallClass(hwnd)
		{ OkCancel(this, title: .Title, :hwnd) }

	controls: #(Vert
		(Form,
			(StaticText, '', group: 0, name: compact), nl, nl,
			(StaticText, 'From', group: 0), (ChooseDateTime, name: from, group: 1), nl,
			(StaticText, 'To', group: 0), (ChooseDateTime, name: to, group: 1))
		)
	New()
		{
		super(.controls)
		t = Transaction(read:)
		// Stagger the default times slightly so they encompass the first / last records
		compact = t.Asof(Date.Begin())
		.FindControl(#compact).
			Set('Compact Date: ' $ compact.Format('yyyy/MMMM/dd, hh:mm'))
		.from.Set(compact.Minus(minutes: 1))
		.to.Set(t.Asof(Date.End()).Plus(minutes: 1))
		t.Complete()
		}

	getter_from()
		{ return .from = .FindControl(#from) }

	getter_to()
		{ return .to = .FindControl(#to) }

	OK()
		{
		if '' is error = .valid(from = .from.Get(), to = .to.Get())
			return [:from, :to]
		.AlertError(.Title, error)
		return false
		}

	valid(from, to)
		{
		error = ''
		if from is '' or to is ''
			error = 'From and To are both required'
		else if from > to
			error = 'From must be less than To'
		return error
		}
	}
