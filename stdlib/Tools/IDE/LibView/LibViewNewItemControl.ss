// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
LibViewItemBase
	{
	Title: "Create New Definition"
	DefaultButton: "Create"
	CallClass(name, msg, path)
		{ return ToolDialog(0, Object(this, name, path, msg)) }

	New(name, path, .msg)
		{ super(name, path) }

	Controls()
		{
		return Object('Vert'
			Object('Static', .msg),
			'Skip',
			#(Pair, (Static, Item), (Field, name: item)),
			#(Pair, (Static, Path), (Field, readonly:, name: path)),
			'Skip',
			#(ExplorerMultiTree, multi?: false, readonly:, ymin: 250),
			#(Skip, 5),
			#(Horz, Fill, (Button, Create, xmin: 50), Skip, (Button, Cancel, xmin: 50)))
		}

	On_Create()
		{
		if not .Valid?()
			return
		// check that name doesn't already exist in library
		path = .Path.Get()
		lib = path.Split('/')[0]
		item = .Item.Get()
		if false isnt Query1(lib, name: item, group: -1)
			.AlertWarn(.Title, item $ ' already exists in ' $ lib)
		else
			.Window.Result([:item, :path])
		}
	}
