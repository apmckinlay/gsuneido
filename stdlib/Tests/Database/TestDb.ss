// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(acts)
		{
		ob = new this
		for (i = 0; i < acts.Size(); ++i)
			ob["Act_" $ acts[i]]()
		}
	New()
		{
		.tran = Stack()
		.fldnum = Object()
		.idxnum = Object()
		}
	Act_t() // transaction
		{
		.tran.Push(Transaction(update:))
		Print(.tran.Top())
		}
	Act_T() // complete transaction
		{
		Print(complete: .tran.Top())
		.tran.Pop().Complete()
		}
	Act_a() // aborted transaction
		{
		Print(rollback: .tran.Top())
		.tran.Pop().Rollback()
		}
	tblnum: 0
	Act_C() // create table
		{
		.db("create testdb" $ ++.tblnum $
			" (a,b,c) key(a) index(b)")
		.fldnum[.tblnum] = .idxnum[.tblnum] = 0
		}
	Act_D() // destroy table
		{
		.db("destroy testdb" $ .tblnum--)
		}
	Act_c() // create field
		{
		.db("alter testdb" $ .tblnum $
			" create (fld" $ ++.fldnum[.tblnum] $ ")")
		}
	Act_d() // delete field
		{
		.db("alter testdb" $ .tblnum $
			" delete (fld" $ .fldnum[.tblnum]-- $ ")")
		}
	Act_i() // create index
		{
		.db("alter testdb" $ .tblnum $
			" create index(fld" $ ++.idxnum[.tblnum] $ ")")
		}
	Act_I() // destroy index
		{
		.db("alter testdb" $ .tblnum $
			" delete index(fld" $ .idxnum[.tblnum]-- $ ")")
		}
	Act_o() // output a record
		{
		Print("output to testdb" $ .tblnum)
		.tran.Top().QueryOutput('testdb' $ .tblnum,
			Object(a: Timestamp(), b: Timestamp()))
		}
	Act_u() // update a record
		{
		Print("update testdb" $ .tblnum)
		x = .tran.Top().Query('testdb' $ .tblnum).Prev()
		x.a = Timestamp()
		x.Update()
		}
	Act_e() // exit
		{
		Exit()
		}
	db(s)
		{
		Print(s)
		Database(s)
		}
	}
