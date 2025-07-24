// Copyright (C) 2014 Suneido Software Corp. All rights reserved worldwide.
// Intended to be used with System() since running executables
// on windows requires the slashes to be `\`.
// Handles getting the paths and executables in a format specific to the
// operating system, that is compatable with System() calls
// Also handles if the process is running an exe in an outside directory
// (i.e. c:\myServerDir> start "C:\MyOtherDir\suneido.exe")
// Dir will return c:\myServerDir, ExePath will return C:\MyOtherDir
class
	{
	New()
		{
		.path = Paths.ToLocal(ServerEval('ExePath'))
		if Sys.Linux?()
			.path = .path.Replace('.exe', '')

		.exe = Paths.Basename(.path)
		.dir = Paths.ToLocal(ServerEval('GetCurrentDirectory'))
		}

	ExeName()
		{
		return .exe
		}

	// Returns the Working Directory for the server
	// This might not be where the exe is running, just where the server is running
	Dir(subdir = "")
		{
		return .dir $ Paths.ToLocal(Opt("/", subdir))
		}

	// Returns the path to the exe that is running (can be different
	// than Dir if the process is pointing to an exe in a different directory)
	ExePath()
		{
		return .path
		}

	ExeDir()
		{
		return .path.BeforeLast(.exe)[.. -1]
		}
	}
