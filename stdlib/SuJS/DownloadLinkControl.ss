// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Show(target, saveName)
		{
		Dialog(0, [this, target, saveName], title: #Download, closeButton?:)
		}

	New(target, saveName)
		{
		super(.layout(target, saveName))
		}

	layout(target, saveName)
		{
		return [#Vert,
			#(Static, "Please click the link below to start the download",
				weight: bold, size: "+2"),
			#Skip,
			[#Html_ahref_, saveName,
				href: #unused,
				hrefOnBrowser: "download" $
					Url.BuildQuery([target, token: SuRenderBackend().Token, :saveName])],
			#Skip]
		}

	Goto(unused)
		{
		.Window.Result(false)
		}
	}
