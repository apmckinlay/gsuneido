// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(s, table = false, embed? = false, extra = #())
		{
		s = (OptContribution('HtmlWrap_Pre',
			function (s, table /*unused*/) { return s }))(s, table)

		s = s.Detab()

		if s.Prefix?("<?xml") or s.Prefix?("<html") or s.Prefix?("<!DOCTYPE")
			return s

		s = BookContent.ToHtml(table, s)

		if table is 'suneidoc'
			s = Asup(s)

		prefix = HtmlPrefix
		suffix = HtmlSuffix
		if table isnt false
			{
			extra = Object(runningHttpServer: ServerSuneido.Get('RunningHttpServer')).
				Merge(extra)
			// passing extra to re-evaluare Memoize when extra param changes
			if false isnt p = HtmlWrapPrefix(table, :extra)
				prefix = p
			if false isnt x = Query1Cached(table $
				" where path = '/res' and name = 'HtmlSuffix'")
				suffix = x.text
			}
		s = prefix $ s $ suffix
		return embed? ? .Embed(s) : s
		}


	mimeTypes: ('jpg': 'jpeg', 'svg': 'svg+xml')
	Embed(s)
		{
		return s.Replace('(suneido:/.*?(' $ SuneidoAPP.Images $ '))',
			{|s|
			imgUrl = s[8..] /*= strip "suneido:" */
			mimeType = imgUrl.AfterLast('.')
			mimeType = .mimeTypes.GetDefault(mimeType, mimeType)
			info = SuneidoAPP.GetBookRec(imgUrl)
			String?(info)
				? s
				: ('data:image/' $ mimeType $ ';base64,' $ Base64.Encode(info.text))
			})
		}
	}
