// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
// Active Suneido Server Page
// replaces <$ expr $> with expr.Eval()
function (page, overrides = #())
	{
	while page[i = page.Find('<$') + 2] isnt ''
		{
		n = page[i..].Find('$>')
		content = page[i :: n].Trim()
		for override in overrides.Members()
			if content.Prefix?(override)
				content = overrides[override] $ content[override.Size()..]
		page = page.ReplaceSubstr(i - 2, n + 4, content.Eval()) /*= include <$ and $> */
		}
	return page
	}
