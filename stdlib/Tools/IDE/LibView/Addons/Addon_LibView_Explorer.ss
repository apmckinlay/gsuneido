// Copyright (C) 2019 Axon Development Corporation All rights reserved worldwide.
LibViewAddon
	{
	Commands(cmds)
		{
		cmds.Add(#(Diff_Item_to, "", "Compare with a previous version from history"))
		}

	Ctrl()
		{
		tabMenu = Object()
		.TabMenuOptions(tabMenu)
		return Object(
			order: 5,
			ctrl: [#ExplorerMulti,
				#LibTreeModel,
				#(LibViewView),
				treeArgs: [multi?:],
				extraTabMenu: tabMenu])
		}

	Init()
		{
		.Redir(#On_Inspect)
		}

	Explorer_RestoreTab(tabOb)
		{
		pathOb = tabOb.path.Split('/')
		if pathOb.Empty?()
			return
		pathOb[0] = .toggleUsed(pathOb[0], used?: .Libs().Has?(pathOb[0].Tr("()")))
		.Explorer.GotoPath(pathOb.Join('/'), skipFolder?: not tabOb.group)
		}

	toggleUsed(text, used?)
		{
		return used? ? text.Tr("()") : text =~ "^\(.*\)$" ? text : '(' $ text $ ')'
		}
	}
