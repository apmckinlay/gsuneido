// Copyright (C) 2018 Axon Development Corporation All rights reserved worldwide.
function (selectors, baseEl = false, _document = false)
	{
	if Type(selectors) isnt "String"
		return selectors

	base = baseEl isnt false
		? baseEl
		: document isnt false
			? document
			: SuUI.GetCurrentDocument()

	try
		res = base.QuerySelector(selectors)
	catch
		return false
	return res
	}
