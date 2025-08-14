// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_RemoveTagFromName()
		{
		fn = LibraryTags.RemoveTagFromName
		Assert(fn('') is: '')
		Assert(fn('Abc') is: 'Abc')
		Assert(fn('Abc__foo') is: 'Abc')
		Assert(fn('Abc__bar') is: 'Abc')
		Assert(fn('Abc__bar__foo') is: 'Abc__bar')
		Assert(fn('Abc__foo__bar') is: 'Abc__foo')

		Assert(fn('Abc__protect') is: 'Abc__protect')
		Assert(fn('Abc__protect__foo') is: 'Abc__protect')
		Assert(fn('Abc__foo__protect') is: 'Abc__foo__protect')
		}

	Test_GetTagFromName()
		{
		fn = LibraryTags.GetTagFromName
		Assert(fn('') is: '')
		Assert(fn('Abc') is: '')
		Assert(fn('Abc__foo') is: '__foo')
		Assert(fn('Abc__bar') is: '__bar')
		Assert(fn('Abc__bar__foo') is: '__foo')
		Assert(fn('Abc__foo__bar') is: '__bar')

		Assert(fn('Abc__protect') is: '')
		Assert(fn('Abc__protect__foo') is: '__foo')
		Assert(fn('Abc__foo__protect') is: '')
		}

	Test_buildTags()
		{
		fn = LibraryTags.LibraryTags_buildTags
		Assert(fn(#(), #()) is: #())
		Assert(fn(#(alpha), #()) is: #(alpha))
		Assert(fn(#(alpha_samples, alpha), #()) is: #(alpha_samples, alpha))

		Assert(fn(#(), #(webgui)) is: #(webgui))
		Assert(fn(#(alpha), #(webgui)) is: #(alpha, webgui, webgui_alpha))
		Assert(fn(#(alpha_samples, alpha), #(webgui))
			is: #(alpha_samples, alpha, webgui, webgui_alpha_samples, webgui_alpha))
		}

	Test_ConvertTagInfo()
		{
		fn = LibraryTags.ConvertTagInfo
		Assert(fn(#()) is: #())
		Assert(fn(#('')) is: #())
		Assert(fn(#('', '__foo')) is: #('foo'))
		Assert(fn(#('', '__foo', '__bar')) is: #('foo', 'bar'))
		}

	Test_SetTrialTag()
		{
		fn = LibraryTags.SetTrialTag
		Assert(fn('', '', #()) is: '')
		Assert(fn('Foo', '', #()) is: 'Foo')
		Assert(fn('Foo__webgui', '', #()) is: 'Foo__webgui')

		trials = #(trial, alpha)
		Assert(fn('Foo__trial', '', trials) is: 'Foo')
		Assert(fn('Foo__webgui_trial', '', trials) is: 'Foo__webgui')

		Assert(fn('Foo', 'alpha', trials) is: 'Foo__alpha')
		Assert(fn('Foo__trial', 'alpha', trials) is: 'Foo__alpha')
		Assert(fn('Foo__alpha', 'alpha', trials) is: 'Foo__alpha')
		Assert(fn('Foo__webgui', 'alpha', trials) is: 'Foo__webgui_alpha')
		Assert(fn('Foo__webgui_alpha', 'alpha', trials) is: 'Foo__webgui_alpha')
		Assert(fn('Foo__webgui_trial', 'alpha', trials) is: 'Foo__webgui_alpha')
		}
	}