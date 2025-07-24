// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.

// TODO: rename global (and _Test)
// TODO: extract constant to member
// TODO: inline local
// TODO: inline method
// TODO: convert method to class
// TODO: move method to superclass
// TODO: move method to subclass
// TODO: make method public
// TODO: make method private (dangerous)
// TODO: move method to global function (if static)

// TODO: add option to preview diff of changes

class
	{
	CallClass(source)
		{
		RefactorControl(source, this)
		}

	Name: ''
	Controls: (Static '')
	Desc: ''
	SelectWord: false
	Init(data/*unused*/)
		{ return true }
	Errors(data/*unused*/)
		{ return "" }
	Warnings(data/*unused*/)
		{ return "" }
	Process(data/*unused*/)
		{ return true }
	Edit_Change(source/*unused*/)
		{ }

	Info(msg)
		{
		Alert(msg, .Name, flags: MB.ICONINFORMATION)
		}
	Warn(msg)
		{
		Alert(msg, .Name, flags: MB.ICONWARNING)
		}
	}
