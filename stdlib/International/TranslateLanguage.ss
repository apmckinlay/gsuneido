// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
// Based on contributions from Roberto Artigas Jr. (rartiga1@midsouth.rr.com)
function (@args)
	// pre: 	arg[0] is a string to translate
	//			and Suneido.Language is an object with a name member
	// post:	returns the translation of the string
	//			or the original from value if no translation is found
	{
	from = args[0]
	if from is ''
		return from

	language = GetLanguage()
	if language.name is 'english' or language.name is ''
		translate = function (from) { return from is 'decimalpoint' ? 'point' : from }
	else
		{
		if Suneido.CacheLanguage isnt language.name or
			not Suneido.Member?('TranslateCache')
			{
			Suneido.TranslateCache = LruCache(GetTranslation, 200) /*= cache size */
			Suneido.CacheLanguage = language.name
			}
		translate = function (from)
			{
			trfrom = from.Trim()
			prefix = from.BeforeFirst(trfrom)
			suffix = from.AfterFirst(trfrom)
			if trfrom.Suffix?('...')
				{
				trfrom = trfrom[..-3] /*= removing '...' */
				suffix = '...' $ suffix
				}
			if ("" is (translation = Suneido.TranslateCache.Get(trfrom.Tr("&"))))
				return from
			if (not trfrom.Has?("&"))
				translation = translation.Tr("&")
			return prefix $ translation $ suffix
			}
		}
	translation = translate(from)
	if args.Size() > 1
		translation = translation.
			Replace("%[0-9]", { |a| translate(args[Number(a[1 ..])]) }).
			Replace("%[a-z]+", { |a| translate(args[a[1 ..]]) })
	return translation
	}
