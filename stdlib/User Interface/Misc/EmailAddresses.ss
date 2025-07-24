// Copyright (C) 2008 Suneido Software Corp. All rights reserved worldwide.
class
	{
	Ensure()
		{
		Database('ensure email_addresses (word, email) index(word) key(email, word)')
		}

	GetAddrs(prefix, limit)
		{
		prefix = .preparePrefix(prefix)
		list = QueryAll('email_addresses
			where word >= ' $ Display(prefix) $
			' and word < ' $ Display(prefix $ '~'),
			:limit)
		list.Map!({ .stripInternalDesc(it.email) })
		return .splitMultipleAddresses(list).UniqueValues().Sort!()
		}

	stripInternalDesc(email)
		{
		if not email.Has?("<") or not email.Has?("*")
			return email

		name = email.BeforeLast("<").BeforeLast("(")
		strippedName = name.BeforeFirst('*')
		return strippedName.Trim() $ " " $ email[name.Size()..]
		}

	preparePrefix(prefix)
		{
		/*
		1. 500 chars should be lots for the prefix of an email address (probably could be
		significantly less than 500). This limit is to prevent index sizes too large
		when searching the email addresses table for matches.
		2. Since the word index entries are all saved as lower values, we should be
			also lowering the input prefixes for the query
		*/
		Assert(String?(prefix), 'unexpected data type for EmailAddresses.preparePrefix')

		prefixLimit = 500
		return prefix[..prefixLimit].Lower()
		}

	DeleteAddr(addr, t)
		{
		DoWithTran(t, update:)
			{|t|
			t.QueryDo('delete email_addresses where email is ' $ Display(addr))
			}
		}

	MaxSavedAddressSize: 200
	OutputAddr(addr, t = false)
		{
		if addr.Size() > .MaxSavedAddressSize
			{
			SuneidoLog('INFO: attempted to save over-sized email address', calls:,
				params: Object(:addr))
			return
			}
		word = email = addr
		if addr.Has?('<')
			{
			addresses = addr.BeforeFirst('<').Replace(' \((.*?)\) $', ' \1').
				Trim().Tr('* ', ' ')
			for (s = addresses; s > ""; s = s.AfterFirst(' '))
				.add1(s, addr, t)
			word = addr.AfterFirst('<').BeforeLast('>')
			}
		else
			{
			x = QueryFirst('email_addresses where word is ' $ Display(email) $
				' sort email')
			if x isnt false and x.email.Has?('<')
				return // don't add without name if it exists with name
			}
		.add1(word, email, t)
		}
	add1(word, email, t)
		{
		DoWithTran(t, update:)
			{|t|
			try
				t.QueryOutput('email_addresses', [word: word.Lower(), :email])
			catch (unused, 'duplicate key')
				; // ignore
			}
		}

	splitMultipleAddresses(list)
		{
		newList = Object()
		for email in list
			{
			if false is emailStart = email.FindLast('<')
				{
				newList.Add(email)
				continue
				}
			addresses = email[emailStart  ..]
			if '' isnt addresses and addresses.Has1of?(';,')
				{
				addresses = addresses.Tr(',', ';').Tr('<> ')
				name = email[.. emailStart]
				for address in addresses.Split(';')
					newList.Add(name $ '<' $ address $ '>')
				}
			else
				newList.Add(email)
			}
		return newList
		}
	}
