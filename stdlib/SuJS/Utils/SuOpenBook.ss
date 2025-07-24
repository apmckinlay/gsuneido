// Copyright (C) 2023 Axon Development Corporation All rights reserved worldwide.
function (book, user, token)
	{
	a = CreateElement('a', SuUI.GetCurrentDocument().body)
	a.href = SuUI.GetCurrentDocument().location.origin $
		'/' $ Url.BuildQuery(Object(preauth: true, :user, :book, :token))
	a.target = "_blank"
	a.Click()
	a.Remove()
	}
