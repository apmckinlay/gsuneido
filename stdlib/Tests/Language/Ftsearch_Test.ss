// Copyright (C) 2023 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		b = Ftsearch.Create()
		b.Add(1, "Big Trees", "douglas fir")
		b.Add(2, "Small Trees", "apple pear not big")
		b.Add(3, "Pretty Flowers", "spring time")
		data = b.Pack()

		ix = Ftsearch.Load(data)
		Assert(Display(ix) like: "Ftsearch{ndocs: 3, unique terms: 11, ntermsTotal: 13,
			avgTermsPerDoc: 4, avgDocsPerTerm: 1, maxDocsPerTerm: 2}")
		Assert(ix.WordInfo("big flowers") like:
			"big count 4 in 2 documents
			flower count 3 in 1 documents")

		Assert(ix.Search("nada") is: #())
		Assert(ix.Search("fir") is: #(1))
		Assert(ix.Search("big") is: #(1, 2))
		Assert(ix.Search("flower") is: #(3))
		Assert(ix.Search("big fir") is: #(1, 2))

		ix.Update(4, "", "", "New One", "nothing special apple")	// add
		Assert(ix.Search("special new") is: #(4))

		ix.Update(2, "Small Trees", "apple pear not big", 			// update
			"Small Shrubs", "Douglas street")
		Assert(ix.Search("shrub") is: #(2))

		ix.Update(3, "Pretty Flowers", "spring time", "", "")	// delete
		Assert(ix.Search("flower") is: #())

		expected = "Ftsearch{ndocs: 3, unique terms: 12, ntermsTotal: 13,
			avgTermsPerDoc: 4, avgDocsPerTerm: 1, maxDocsPerTerm: 2}"
		Assert(Display(ix) like: expected)

		b = Ftsearch.Create()
		b.Add(1, "Big Trees", "douglas fir")
		b.Add(2, "Small Shrubs", "Douglas street")
		b.Add(4, "New One", "nothing special apple")
		ix = b.Index()
		Assert(Display(ix) like: expected)
		}
	}