// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_openrouterUrl()
		{
		openrouterUrl = AiAgentModelHelp.AiAgentModelHelp_openrouterUrl
		Assert(openrouterUrl("deepseek/deepseek-v4.1-flash")
			is: "https://openrouter.ai/deepseek/deepseek-v4.1-flash")
		// the Addon_url used to display the links must be able to match them
		for m in AiAgentControl.AiAgentControl_models.Members()
			Assert(false isnt Addon_url.MatchUrl(openrouterUrl(m)))
		}

	Test_control()
		{
		// the control is sized to fit the text (in characters and lines)
		control = AiAgentModelHelp.AiAgentModelHelp_control
		modelHelpText = AiAgentModelHelp.AiAgentModelHelp_modelHelpText
		models = #("deepseek/deepseek-v4.1-flash": [context: "1M", in: .15, out: .60])
		ctrl = control(models)
		text = modelHelpText(models)
		Assert(ctrl[1].width is: text.Lines().Map(#Size).Max() + 2 /*= control margins */)
		// header + dashes + one model + blank + note
		Assert(ctrl[1].height is: 4)
		}
	}
