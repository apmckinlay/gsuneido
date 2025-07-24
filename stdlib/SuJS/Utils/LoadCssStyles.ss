// Copyright (C) 2018 Axon Development Corporation All rights reserved worldwide.
function (cssName, styles, override? = false, media = 'screen')
	{
	loadedCss = Suneido.GetInit(#LoadedCss, Object())
	if loadedCss.Member?(cssName)
		{
		if not override?
			return
		loadedCss[cssName].Remove()
		}

	styleEl = CreateElement('style', SuUI.GetCurrentDocument().head)
	styleEl.type = 'text/css'
	if media is 'both'
		styleEl.innerHTML = styles
	else
		styleEl.innerHTML = '@media ' $ media $ ' {' $ styles $ '}'

	loadedCss[cssName] = styleEl
	}
