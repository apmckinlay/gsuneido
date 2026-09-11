// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
// displays the models used by AiAgentControl with links to openrouter.ai
Controller
	{
	CallClass(hwnd, models)
		{
		return Dialog(hwnd, [this, models], closeButton?:, title: #Models)
		}

	New(models)
		{
		super(.control(models))
		// NOTE: set readonly this way (not with a readonly: argument) so that
		// the addon's IdleAfterChange still runs and highlights the urls
		.FindControl(#models).SetReadOnly(true)
		}

	// size the control (in characters and lines) to fit the text
	control(models)
		{
		text = .modelHelpText(models)
		lines = text.Lines()
		return [#Vert,
			[#ScintillaAddons, name: #models, wrap:, xstretch: 0,
				ystretch: 0, width: lines.Map(#Size).Max() + 2, /*= margins */
				height: lines.Size() + 1, Addon_url:, set: text]]
		}

	modelHelpText(models)
		{
		// the Link column is sized to fit the longest url
		w0 = models.Members().Map({ .openrouterUrl(it).Size() }).Max()
		w1 = 7
		w2 = 5
		w3 = 5
		s = "Link".RightFill(w0) $ "Context".LeftFill(w1) $ "In".LeftFill(w2) $
			"Out".LeftFill(w3) $ '\n'
		s $= '-'.Repeat(w0 + w1 + w2 + w3) $ '\n'
		for m in models.Members().Sort!()
			{
			x = models[m]
			s $= .openrouterUrl(m).RightFill(w0) $ x.context.LeftFill(w1) $
				x.in.Format("##.##").LeftFill(w2) $ x.out.Format("##.##").LeftFill(w3) $
				'\n'
			}
		return s
		}

	openrouterUrl(model)
		{
		return "https://openrouter.ai/" $ model
		}
	}
