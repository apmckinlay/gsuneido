// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.

Test
	{
	dupeCodeTest: (`// Copyright (C) 2013 Suneido Software
			Corp. All rights reserved worldwide.
class
	{
	N: 8 // block size (number of lines)
	CallClass(libs)
		{
		if String?(libs)
			libs = [libs]
		(new this).Detect(libs)
		}
	Detect(libs)
		{
		.hashes = Object()
		for lib in libs
			.process(lib)
		.output()
		}
	process(lib)
		{
		QueryApply(lib $ ' where name !~ 'Test$'', group: -1)
			{|x|
			.process1(lib, x)
			}
		}
	// hash .N non-blank lines
	hash(lines, i, j, k)
		{
		hasher = Adler32()
		for (j = i, n = 0; n < .N and j < lines.Size(); ++j)
			if lines[j] isnt '' // skip blank lines
				{
				hasher.Update(lines[j])
				++n
				}
		return hasher.Value()
		}
	output()
		{
		dups = .hashes.Values().Filter({ it.Has?(',') })?/////////////////////////////////
			///////////////////////////
		dups.Sort!().Each { Print(it.Tr(',', '\t')) }
		Print(dups.Size())
		}
	}"`,

	`// Copyright (C) 2013 Suneido Software
			Corp. All rights reserved worldwide.
class
	{
	N: 8 // block size (number of lines)
	CallClass(libs)
		{
		if String?(libs)
			libs = [libs]
		(new this).Detect(libs)
		}
	Detect(libs)
		{
		.hashes = Object()
		for lib in libs
			.process(lib)
		.output()
		}
	process(lib)
		{
		QueryApply(lib $ ' where name !~ 'Test$'', group: -1)
			{|x|
			.process1(lib, x)
			}
		}
	// hash .N non-blank lines
	hash(lines, i, j, k)
		{
		hasher = Adler32()
		for (j = i, n = 0; n < .N and j < lines.Size(); ++j)
			if lines[j] isnt '' // skip blank lines
				{
				hasher.Update(lines[j])
				++n
				}
		return hasher.Value()
		}
	output()
		{
		dups = .hashes.Values().Filter({ it.Has?(',') })?/////////////////////////////////
			///////////////////////////
		dups.Sort!().Each { Print(it.Tr(',', '\t')) }
		Print(dups.Size())
		}
	}"`,

	`function_test = "// Copyright (C) 2017
		Suneido Software Corp. All rights reserved worldwide.
function (x1, x2, x3, x4, x5, x6)
	{
	x = 5
	y = 55
	if (x > y)
		{
		return 55
		}
	return (y > x)
	}
		dups = .hashes.Values().Filter({ it.Has?(',') })?/////////////////////////////////
			///////////////////////////
		dups.Sort!().Each { Print(it.Tr(',', '\t')) }
		Print(dups.Size())dups = .hashes.Values().Filter({ it.Has?(',') })?///////////////
			/////////////////////////////////////////////
		dups.Sort!().Each { Print(it.Tr(',', '\t')) }
		Print(dups.Size())dups = .hashes.Values().Filter({ it.Has?(',') })?///////////////
			/////////////////////////////////////////////
		dups.Sort!().Each { Print(it.Tr(',', '\t')) }
		Print(dups.Size())`,

	`function_test = "// Copyright (C) 2017
		Suneido Software Corp. All rights reserved worldwide.
function (x1, x2, x3, x4, x5, x6)
	{
	x = 5
	y = 55
	if (x > y)
		{
		return 55
		}
	return (y > x)
	}
		dups = .hashes.Values().Filter({ it.Has?(',') })?/////////////////////////////////
			///////////////////////////
		dups.Sort!().Each { Print(it.Tr(',', '\t')) }
		Print(dups.Size())dups = .hashes.Values().Filter({ it.Has?(',') })?///////////////
			/////////////////////////////////////////////
		dups.Sort!().Each { Print(it.Tr(',', '\t')) }
		Print(dups.Size())dups = .hashes.Values().Filter({ it.Has?(',') })?///////////////
			/////////////////////////////////////////////
		dups.Sort!().Each { Print(it.Tr(',', '\t')) }
		Print(dups.Size())`,

	`	dups = .hashes.Values().Filter({ it.Has?(',') })?/////////////////////////////////
			///////////////////////////
		dups.Sort!().Each { Print(it.Tr(',', '\t')) }
		Print(dups.Size())dups = .hashes.Values().Filter({ it.Has?(',') })?///////////////
			/////////////////////////////////////////////
		dups.Sort!().Each { Print(it.Tr(',', '\t')) }
		Print(dups.Size())dups = .hashes.Values().Filter({ it.Has?(',') })?///////////////
			/////////////////////////////////////////////
		dups.Sort!().Each { Print(it.Tr(',', '\t')) }
		Print(dups.Size())
function_test = "// Copyright (C) 2017
		Suneido Software Corp. All rights reserved worldwide.
function (x1, x2, x3, x4, x5, x6)
	{
	x = 5
	y = 55
	if (x > y)
		{
		return 55
		}
	return (y > x)
	}`,

	`function() { // no dupe code here }`
	)

	Test_Continuous_CheckDuplicates()
		{
		//Create test libraries to run Qc_Duplicate_Code fast
		.MakeLibraryRecord([name: "test1", text: .dupeCodeTest[0]])
		.MakeLibraryRecord([name: "test2", text: .dupeCodeTest[1]])
		.MakeLibraryRecord([name: "test3", text: .dupeCodeTest[2]])
		.MakeLibraryRecord([name: "test4", text: .dupeCodeTest[3]])
		.MakeLibraryRecord([name: "test5", text: .dupeCodeTest[4]])
		.MakeLibraryRecord([name: "test6", text: .dupeCodeTest[5]])

		recordData = Record()
		recordData.recordName = "test1"
		recordData.lib = "Test_lib"
		dupCode = Qc_DuplicateCode(recordData)
		recordData.recordName = "test3"
		dupCode2 = Qc_DuplicateCode(recordData)

		Assert(dupCode.warnings isSize: 5)
		Assert(dupCode.warnings
			has: [name: "Duplicate: Test_lib:test1:1, Test_lib:test2:1"])
		Assert(dupCode.warnings
			has: [name: "Duplicate: Test_lib:test1:9, Test_lib:test2:9"])
		Assert(dupCode.warnings
			has: [name: "Duplicate: Test_lib:test1:17, Test_lib:test2:17"])
		Assert(dupCode.warnings
			has: [name: "Duplicate: Test_lib:test1:26, Test_lib:test2:26"])
		Assert(dupCode.warnings
			has: [name: "Duplicate: Test_lib:test1:34, Test_lib:test2:34"])
		Assert(dupCode.desc
			is: "Duplicate code from in this class was found " $
				"-> This check does not affect the rating of code")

		Assert(dupCode2.warnings isSize: 3)
		Assert(dupCode2.warnings
			has: [
				name: "Duplicate: Test_lib:test3:1, Test_lib:test4:1, Test_lib:test5:11"
				])
		Assert(dupCode2.warnings
			has: [name: "Duplicate: Test_lib:test3:13, Test_lib:test5:1"])
		Assert(dupCode2.warnings
			has: [name: "Duplicate: Test_lib:test3:9, Test_lib:test4:9"])
		Assert(dupCode2.desc
			is: "Duplicate code from in this class was found " $
				"-> This check does not affect the rating of code")

		recordData.recordName = "test6"
		noDupeCode = Qc_DuplicateCode(recordData)
		Assert(noDupeCode.desc is: "No duplicates found of code in this class -> " $
			"This check does not affect the rating of code")
		}
	}