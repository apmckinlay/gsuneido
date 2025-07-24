// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		EnsureDir('test_deletedir')
		EnsureDir('test_deletedir/one')
		EnsureDir('test_deletedir/one/two')
		EnsureDir('test_deletedir/three')
		PutFile('test_deletedir/file1', 'x')
		PutFile('test_deletedir/file2', 'x')
		PutFile('test_deletedir/one/onefile1', 'x')
		PutFile('test_deletedir/one/onefile2', 'x')
		PutFile('test_deletedir/one/two/twofile1', 'x')
		PutFile('test_deletedir/one/two/twofile2', 'x')
		PutFile('test_deletedir/three/threefile1', 'x')

		Assert(DeleteDir('test_deletedir'))
		}
	}