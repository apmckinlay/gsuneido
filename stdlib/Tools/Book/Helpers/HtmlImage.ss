// Copyright (C) 2013 Axon Development Corporation All rights reserved worldwide.
function (image, clas = 'ss', id = '1')
	{
	return '<img id="ss' $ id $ '" class="' $ clas $ '"
		title="Click to expand/shrink"
		src="suneido:/' $ _table $ '/res/' $ image $ '"
		onClick="screenshot(' $ id $ ')" />'
	}
