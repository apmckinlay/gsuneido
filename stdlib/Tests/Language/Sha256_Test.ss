// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		Assert(Sha256("").ToHex()
			is: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855")
		Assert(Sha256("\x00\xff").ToHex()
			is: "06eb7d6a69ee19e5fbdf749018d3d2abfa04bcbd1365db312eb86dc7169389b8")

		cksum = "936a185caaa266bb9cbe981e9e05cb78cd732b0b3280eb944412bb6f8f8f07af"
		Assert(Sha256("helloworld").ToHex() is: cksum)
		Assert(Sha256("hello", "world").ToHex() is: cksum)
		Assert(Sha256().Update("hello").Update("world").Value().ToHex() is: cksum)
		}
	}