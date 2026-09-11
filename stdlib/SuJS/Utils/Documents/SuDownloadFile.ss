// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
class
	{
	CallClass(target, saveName)
		{
		if saveName.Lower().AfterLast('.') in (#pdf, #csv) and .isSafariMobile?()
			{
			SuRender().Event(false, "DownloadLinkControl.Show", [target, saveName])
			return
			}

		// opening downloaded files automatically needs to be configured on the user's browser
		a = CreateElement('a', SuUI.GetCurrentDocument().body)
		a.href =
			"download" $ Url.BuildQuery([target, token: SuRender().GetToken(), :saveName])
		a.download = saveName
		a.Click()
		a.Remove()
		}

	isSafariMobile?()
		{
		navigator = SuUI.GetCurrentWindow().navigator
		userAgent = navigator.userAgent.Lower()
		if userAgent.Has?(#chrome)
			return false
		return userAgent.Has?(#iphone) or userAgent.Has?(#ipad) or
			userAgent.Has?(#macintosh) and navigator.maxTouchPoints > 1 // ipad
		}
	}
