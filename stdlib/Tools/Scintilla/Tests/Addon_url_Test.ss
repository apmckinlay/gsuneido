// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_ProcessChunk()
		{
		mock = Mock(Addon_url)
		mock.When.IndicatorIdx([anyArgs:]).Return('indicator')
		mock.When.MatchUrl([anyArgs:]).CallThrough()
		mock.Eval(Addon_url.Init)
		mock.Eval(Addon_url.ProcessChunk, '', 0)
		mock.Verify.Never().mark_word([anyArgs:])

		mock.Eval(Addon_url.ProcessChunk, 'hell world', 0)
		mock.Verify.Never().mark_word([anyArgs:])

		text = 'test http://suneido.com test'
		mock.Eval(Addon_url.ProcessChunk, text, 0)
		mock.Verify.mark_word(5, 18)

		text = 'testwww.g' $ 'o'.Repeat(500) $ 'gle test'
		mock.Eval(Addon_url.ProcessChunk, text, 0)
		mock.Verify.mark_word(4, 508)

		text = 'addr1 https://test.com
ftp://test.com abc http://test.com/test#3?test=true end'
		mock.Eval(Addon_url.ProcessChunk, text, 0)
		mock.Verify.mark_word(6, 16)
		mock.Verify.mark_word(24, 14)
		mock.Verify.mark_word(43, 32)

		text = 'sftp://test.com'
		mock.Eval(Addon_url.ProcessChunk, text, 0)
		mock.Verify.mark_word(0, 15)

		text = 'start sftp://test.com end'
		mock.Eval(Addon_url.ProcessChunk, text, 0)
		mock.Verify.mark_word(6, 15)

		text = 'start new ftps://test.com end'
		mock.Eval(Addon_url.ProcessChunk, text, 0)
		mock.Verify.mark_word(10, 15)

		text = (" ".Repeat(1000) $ "https://" $ "x".Repeat(2000) $ " ".Repeat(1000)).
			Repeat(5)
		mock.Eval(Addon_url.ProcessChunk, text, 0)
		mock.Verify.mark_word(1000, 2008)
		mock.Verify.mark_word(5008, 2008)
		mock.Verify.mark_word(9016, 2008)
		mock.Verify.mark_word(13024, 2008)
		mock.Verify.mark_word(17032, 2008)
		}
	}
