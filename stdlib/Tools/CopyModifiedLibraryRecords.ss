// Copyright (C) 2012 Suneido Software Corp. All rights reserved worldwide.
// e.g. CopyModifiedLibraryRecords(Libraries(), "demobookoptions", "Changes")
class
	{
	CallClass(srclibs, dstlib, dstfolder)
		{
		srclibs = srclibs.Difference(['configlib', dstlib])
		num = .createDstFolder(dstlib, dstfolder)
		Database("ensure " $ dstlib $ " (lib_modified, copymodlib)")
		for lib in srclibs
			.copyLibrary(lib, dstlib, num)
		}
	createDstFolder(dstlib, dstfolder)
		{
		if false is x = Query1(dstlib, name: dstfolder, group: 0, parent: 0)
			QueryOutput(dstlib,
				x = [num: .nextNum(dstlib), name: dstfolder, group: 0, parent: 0])
		return x.num
		}
	copyLibrary(lib, dstlib, dstfoldernum)
		{
		QueryApply(lib $ " where lib_modified > ''", group: -1)
			{|x|
			.copyRecord(lib, x, dstlib, dstfoldernum)
			}
		}
	copyRecord(lib, srcrec, dstlib, dstfoldernum)
		{
		.checkIfOverloaded(lib, srcrec.name, dstlib)
		Transaction(update:)
			{|t|
			if false is x = t.Query1(dstlib, name: srcrec.name, group:-1)
				{
				Print("output", lib, srcrec.name)
				QueryOutput(dstlib, [copymodlib: lib, num: .nextNum(dstlib),
					name: srcrec.name, text: srcrec.text,
					lib_invalid_text: srcrec.lib_invalid_text
					parent: dstfoldernum, group: -1, lib_modified: srcrec.lib_modified])
				}
			else
				{
				Assert(x.copymodlib is: lib)
				old = x.Copy()
				x.lib_modified = srcrec.lib_modified
				x.text = srcrec.text
				x.lib_invalid_text = srcrec.lib_invalid_text
				if x isnt old
					{
					Print("update", lib, x.name)
					x.Update()
					}
				}
			}
		}
	checkIfOverloaded(lib, name, dstlib)
		{
		libs = Libraries().Difference([dstlib])
		libi = libs.Find(lib)
		Assert(libs.Member?(libi))
		for (i = libi + 1; i < libs.Size(); ++i)
			if false isnt Query1(libs[i], :name, group: -1)
				Print("WARNING: " $ name $ " is overloaded in " $ libs[i])
		}
	nextNum(lib)
		{
		return QueryMax(lib, "num", 0) + 1
		}
	}
