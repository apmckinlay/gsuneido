// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	passphrase: 'hello world'
	Test_symmetric()
		{
		original = "now is the time for all good people"
		encrypted = OpenPGP.SymmetricEncrypt(.passphrase, original)
		decrypted = OpenPGP.SymmetricDecrypt(.passphrase, encrypted)
		Assert(decrypted is: original)
		}
	Test_asymmetric()
		{
		original = "now is the time for all good people"
		k = OpenPGP.KeyGen("name", "email", .passphrase)
		encrypted = OpenPGP.PublicEncrypt(k.public, original)
		decrypted = OpenPGP.PrivateDecrypt(k.private, .passphrase, encrypted)
		Assert(decrypted is: original)
		}
	}