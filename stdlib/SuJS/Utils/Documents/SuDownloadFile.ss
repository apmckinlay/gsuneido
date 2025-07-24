// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
function (target, saveName)
	{
	// opening downloaded files automatically needs to be configured on the user's browser
	a = CreateElement('a', SuUI.GetCurrentDocument().body)
	a.href = 'download' $
		Url.BuildQuery(Object(target, token: SuRender().GetToken(), :saveName))
	a.download = saveName
	a.Click()
	a.Remove()
	}
