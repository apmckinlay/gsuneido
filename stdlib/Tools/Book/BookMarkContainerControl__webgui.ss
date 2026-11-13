// Copyright (C) 2020 Axon Development Corporation All rights reserved worldwide.
/* README:
	While BookControl can be constructed with: "book: false", BookMarkContainerControl
	is only added to the BookControl when "book" is not false. As a result, we can
	assume throughout this record that "book" will always be set to a valid book name.
*/
PassthruController
	{
	Name:          #BookMarkContainer
	ComponentName: #BookMarkContainer
	New(.book)
		{
		.list = .FindControl(#BookMarkList)
		tool = .FindControl(#BookMarkTool)
		.list.Controller = tool.Controller = .Controller
		tool.Parent = this
		}

	Controls()
		{
		return [#Vert,
			#(BookMarkTool, name: BookMarkTool),
			[#BookMarkList, .book]]
		}

	Default(@args)
		{
		.list[args[0]](@+1args)
		}
	}
