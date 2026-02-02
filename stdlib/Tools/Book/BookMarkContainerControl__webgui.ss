// Copyright (C) 2020 Axon Development Corporation All rights reserved worldwide.
PassthruController
	{
	Name: "BookMarkContainer"
	ComponentName: "BookMarkContainer"
	New()
		{
		.list = .FindControl('BookMarkList')
		tool = .FindControl('BookMarkTool')
		.list.Controller = tool.Controller = .Controller
		tool.Parent = this
		}

	Controls: #(Vert
		(BookMarkTool, name: 'BookMarkTool'),
		BookMarkList)

	SetState(@args)
		{
		.list.SetState(@args)
		}

	GetState()
		{
		return .list.GetState()
		}

	Default(@args)
		{
		.list[args[0]](@+1 args)
		}
	}