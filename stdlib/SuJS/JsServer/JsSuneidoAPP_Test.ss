// Copyright (C) 2023 Axon Development Corporation All rights reserved worldwide.
Test
	{
	cases: (
		('', ''),
		('this is a test', 'this is a test'),
		(
			'<img src="suneido:/ETAHelp/res/11.png">',
			'<img src="/suneidoapp/ETAHelp/res/11.png">'),
		(
			'<style>background: url("suneido:/ETAHelp/res/11.png")</style>',
			'<style>background: url("/suneidoapp/ETAHelp/res/11.png")</style>'),
		(
			'<a href="suneido:/test.html">Test link</a>',
			`<a href="javascript:suIframeSend('L3Rlc3QuaHRtbA==');void(0)">Test link</a>`)
		(
			'<a href="https://example.com">External link</a>',
			'<a href="https://example.com" target="_blank" ' $
				'rel="noopener noreferrer">External link</a>'))
	Test_convert()
		{
		fn = JsSuneidoAPP.JsSuneidoAPP_convert

		.cases.Each({ Assert(fn(it[0]) is: it[1]) })
		Assert(fn('<head>' $ .cases.Map({ it[0] }).Join('\r\n') $ '</head>')
			is: '<head>' $ .cases.Map({ it[1] }).Join('\r\n') $ '</head>')
		}
	}