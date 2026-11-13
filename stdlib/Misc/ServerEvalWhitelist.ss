class
	{
	CallClass()
		{
		contribs = GetContributions(#ServerEvalWhitelist)
		return contribs.ide.Union(contribs.customer)
		}

	GetIdeWhitelist()
		{
		contribs = GetContributions(#ServerEvalWhitelist)
		return contribs.ide
		}

	GetCustomerWhitelist()
		{
		contribs = GetContributions(#ServerEvalWhitelist)
		return contribs.customer
		}
	}
