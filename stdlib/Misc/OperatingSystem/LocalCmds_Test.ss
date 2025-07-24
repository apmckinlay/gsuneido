 // Copyright (C) 2021 Suneido Software Corp. All rights reserved worldwide.
// TAGS: windows
Test
	{
	Test_main()
		{
		Assert(LocalCmds.Nl is: Sys.Windows?() ? "\r\n" : "\n")
		}
	Test_GetLastBootTime_windows()
		{
		cmd = new LocalCmds
		spy = .SpyOn(RunPipedOutput)
		spy.Return("20210426050013.500000-360", "", "LastBootUpTimefalse")
		Assert(cmd.GetLastBootTime() is: #20210426.050013500)
		Assert(cmd.GetLastBootTime() is: '')
		Assert(cmd.GetLastBootTime() is: 'LastBootUpTimefalse')
		spy.Close()
		}

	Test_GetCpuName()
		{
		cmd = new LocalCmds
		spy = .SpyOn(RunPipedOutput)
		spy.Return('Intel(R) Xeon(R) CPU E5-2697 v4 @ 2.30GHz\r\nIntel(R) Xeon(R) ' $
			'CPU E5-2697 v4 @ 2.30GHz\r\nIntel(R) Xeon(R) CPU E5-2697 v4 @ 2.30GHz\r\n' $
			'Intel(R) Xeon(R) CPU E5-2697 v4 @ 2.30GHz')
		Assert(cmd.GetLastBootTime() is: 'Intel(R) Xeon(R) CPU E5-2697 v4 @ 2.30GHz')
		spy.Close()
		}

	Test_GetLastBootTime_linux()
		{
		spy = .SpyOn(RunPipedOutput)
		spy.Return("2023-06-14 14:05:37"

			"uptime -V display"
			" 22:50:15 up 11 days, 16:49,  1 user,  " $
			"load average: 0.00, 0.00, 0.00",

			"uptime -V display"
			" 14:06:52 up 1 min,  1 user,  load average: 0.02, 0.01, 0.00"

			"uptime -V display"
			" 15:58:51 up  1:53,  1 user,  load average: 0.00, 0.00, 0.00")
		linux = new LocalCmds_Linux

		Assert(linux.GetLastBootTime() is: Date("2023-06-14 14:05:37"))

		Assert(linux.GetLastBootTime()
			is: Date().NoTime().Plus(hours: 22, minutes: 50).
				Minus(days: 11, hours: 16, minutes: 49))

		Assert(linux.GetLastBootTime()
			is: Date().NoTime().Plus(hours: 14, minutes: 6).
				Minus(days: 0, hours: 0, minutes: 1))

		Assert(linux.GetLastBootTime()
			is: Date().NoTime().Plus(hours: 15, minutes: 58).
				Minus(days: 0, hours: 1, minutes: 53))
		}

	Test_GetCpuCores_Win()
		{
		spy = .SpyOn(RunPipedOutput)
		spy.Return(
			// normal
			"8",
			"8",
			"8"

			// when cores are split into lines
			"1
			1
			1
			1",
			"1
			1
			1
			1",
			"1
			1
			1
			1" // 4

			// empty
			"",
			"",
			""

			// from some azure hyper v setup
			"NumberOfCores

			4

			",
			"NumberOfEnabledCore

			8

			",
			"Processor\nProcessor\nProcessor\nProcessor\nProcessor\nProcessor\n" $
			"Processor\nProcessor\n" // 8

			// when NumberOfEnabledCore is not available
			"NumberOfCores

			4

			",
			"Node - TEST

			ERROR:

			Description = Invalid query
			"
			"Processor\nProcessor\nProcessor\nProcessor\nProcessor\nProcessor\n" $
			"Processor\nProcessor\n",  // 8

			// when cpu cores is more than logic processors
			"NumberOfCores
			1
			1
			1
			1

			",
			"NumberOfEnabledCore
			1
			1
			1
			1

			",

			"1
			1"

			// when GetLogicalProcessors throw program error
			"NumberOfCores
			1
			1
			1
			1

			",
			"NumberOfEnabledCore
			1
			1
			1
			1

			",

			false // force GetLogicalProcessors to throw program error

			// when GetLogicalProcessors returns 0
			"NumberOfCores
			1
			1
			1
			1

			",
			"NumberOfEnabledCore
			1
			1
			1
			1

			",

			'' // force GetLogicalProcessors to throw program error
			)
		cmd = new LocalCmds
		Assert(cmd.GetCpuCores() is: 8)
		Assert(cmd.GetCpuCores() is: 4)
		Assert(cmd.GetCpuCores() is: 0)
		Assert(cmd.GetCpuCores() is: 8)
		Assert(cmd.GetCpuCores() is: 4)
		Assert(cmd.GetCpuCores() is: 2)
		Assert(cmd.GetCpuCores() is: 4)
		Assert(cmd.GetCpuCores() is: 4)
		spy.Close()

		// use try catch to show more detailed errors,
		// which fails sporadically: "expected a value greater than 11 but it was 11 - "
		try
			{
			// when NumberOfCores is not available
			.SpyOn(RunPipedOutput).Return("testing invalid result")
			cmd.GetCpuCores()
			}
		catch (err)
			{
			Assert(err has: 'testing invalid result',
				msg: FormatCallStack(err.Callstack(), levels: 10))
			return
			}
		Assert(false msg: 'should not reach here')
		}

	Test_GetNetworkSpeed_Win()
		{
		spy = .SpyOn(RunPipedOutput)
		spy.Return(
'Name                                           Speed
----                                           -----
Intel(R) Ethernet Connection (17) I219-LM 1000000000'
			)

		cmd = new LocalCmds
		Assert(cmd.GetNetworkSpeed() is: "Name                                         " $
			"  Speed\r\r\nIntel(R) Ethernet Connection (17) I219-LM 1000000000")
		spy.Close()
		}

	Test_linux_amd64()
		{
		spy = .SpyOn(RunPipedOutput)
		spy.Return(`model name	: Intel(R) Xeon(R) E-2276G CPU @ 3.80GHz
			stepping	: 10
			microcode	: 0xd6
			cpu MHz		: 3792.000
			cache size	: 12288 KB
			physical id	: 0
			siblings	: 2
			core id		: 1
			cpu cores	: 2`,

			`model name	: Intel(R) Xeon(R) E-2276G CPU @ 3.80GHz
			stepping	: 10
			microcode	: 0xd6
			cpu MHz		: 3792.000
			cache size	: 12288 KB
			physical id	: 0
			siblings	: 2
			core id		: 1`,

			`stepping	: 10
			microcode	: 0xd6
			cpu MHz		: 3792.000
			cache size	: 12288 KB
			physical id	: 0
			siblings	: 2
			core id		: 1
			cpu cores	: 4`

			`model name	: Intel(R) Xeon(R) E-2276G CPU @ 3.80GHz
			stepping	: 10
			microcode	: 0xd6
			cpu MHz		: 3792.000
			cache size	: 12288 KB
			physical id	: 0
			siblings	: 2
			core id		: 1`,)
		linux = new LocalCmds_Linux
		Assert(linux.GetCpuCores() is: 2)
		Assert(linux.GetCpuCores() is: 'not available')
		Assert(linux.GetCpuCores() is: 4)

		Assert(linux.GetCpuName() is: 'Intel(R) Xeon(R) E-2276G CPU @ 3.80GHz')
		}

	Test_linxu_arm64()
		{
		spy = .SpyOn(RunPipedOutput)
		spy.Return(`processor	: 0
BogoMIPS	: 243.75
Features	: fp asimd evtstrm aes pmull sha1 sha2 crc32 atomics fphp asimdhp cpuid` $
	` asimdrdm lrcpc dcpop asimddp ssbs
CPU implementer	: 0x41
CPU architecture: 8
CPU variant	: 0x3
CPU part	: 0xd0c
CPU revision	: 1

processor	: 1
BogoMIPS	: 243.75
Features	: fp asimd evtstrm aes pmull sha1 sha2 crc32 atomics fphp asimdhp cpuid` $
	` asimdrdm lrcpc dcpop asimddp ssbs
CPU implementer	: 0x41
CPU architecture: 8
CPU variant	: 0x3
CPU part	: 0xd0c
CPU revision	: 1

processor	: 2
BogoMIPS	: 243.75
Features	: fp asimd evtstrm aes pmull sha1 sha2 crc32 atomics fphp asimdhp cpuid` $
	` asimdrdm lrcpc dcpop asimddp ssbs
CPU implementer	: 0x41
CPU architecture: 8
CPU variant	: 0x3
CPU part	: 0xd0c
CPU revision	: 1`)
		linux = new LocalCmds_Linux
		Assert(linux.GetCpuName() is: '0x41,8,0x3 @ 243.75 BogoMIPS')
		Assert(linux.GetCpuCores() is: 'not available')
		}
	}
