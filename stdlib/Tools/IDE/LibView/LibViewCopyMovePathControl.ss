// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
LibViewItemBase
	{
	CallClass(names, pathVal, title)
		{
		return ToolDialog(0, Object(this, names.Join('\r\n'), pathVal, title),
			title: title $ ' To')
		}

	New(.itemVal, pathVal, title)
		{
		super(.itemVal, pathVal, .controls(title))
		}

	controls(title)
		{
		return Object('Vert',
			#(Pair, (Static, Item), (Editor, readonly:, name: item)),
			Object('Pair', Object('Static', title $ ' To'),
				#(Field, readonly:, name: path)),
			'Skip',
			#(ExplorerMultiTree, multi?: false, readonly:, ymin: 250),
			#(Skip 5),
			Object('Horz', 'Fill',
				Object('Button', title, command: 'move' xmin: 50),
				'Skip', #(Button, Cancel, xmin: 50)))
		}

	On_move()
		{
		if .Valid?()
			.Window.Result(.Path.Get())
		}
	}
