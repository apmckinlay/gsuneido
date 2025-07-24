// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
// NOTE: too slow and touchy to make a regular test
// requires dd on the path for RunPiped
Test
	{
	Test_one()
		{
		srcFile = .MakeFile("helloworld")
		dstFile = .MakeFile()
		Suneido.TestCopyTo = dstFile // for serverDst
		fileSrc = {|block| File(srcFile, "r", :block) }
		fileDst = {|block| File(dstFile, "w", :block) }

		pipeSrc = {|block| RunPiped("dd if=" $ srcFile, :block) }
		pipeDst = {|block|
			RunPiped("dd of=" $ dstFile, :block)
			Thread.Sleep(100) }

		sockSrc = {|block|
			Thread(.serverSrc, name: "TestCopyTo")
			SocketClient("localhost", 8888, :block) }
		sockDst = {|block|
			Thread(.serverDst, name: "TestCopyTo")
			SocketClient("localhost", 9999, :block)
			Thread.Sleep(100) }

		for si, src in [fileSrc, pipeSrc, sockSrc]
			for di, dst in [fileDst, pipeDst, sockDst]
				for read1 in #(false, true)
					{
					.test(src, dst, read1)
					Assert(GetFile(dstFile) is: "helloworld", msg: si $ "," $ di)
					DeleteFile(dstFile)
					}
		}
	test(src, dst, read1)
		{
		size = 10
		src()
			{|from|
			dst()
				{|to|
				if read1
					{
					Assert(from.Read(1) is: 'h')
					to.Write('h')
					--size
					}
				n = from.CopyTo(to, size)
				Assert(n is: size)
				// to.Write(from.Read(size))
				}
			}
		}
	serverSrc: SocketServer
		{
		Port: 8888
		Killer(.killer)
			{ }
		Run()
			{
			.Write("helloworld")
			.killer.Kill()
			}
		}
	serverDst: SocketServer
		{
		Port: 9999
		Killer(.killer)
			{ }
		Run()
			{
			PutFile(Suneido.TestCopyTo, .Read(10))
			.killer.Kill()
			}
		}
	Teardown()
		{
		Suneido.Delete(#TestCopyTo)
		super.Teardown()
		}
	}