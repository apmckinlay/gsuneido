// Copyright (C) 2023 Suneido Software Corp. All rights reserved worldwide.
// These are queries that have failed in the past
// and therefore may be more likely to fail in the future.
#(
`((((cus union cus) join (ivc where i4 < "24")) union (((cus join ivc) union
	(cus join ivc)) where ik is "82")) join (aln where a2 < "32"))`,
`(cus join ((((ivc union ivc) where i3 is "21") join (aln union aln)) union
	((ivc union ivc) join (bln union ((bln where ik is "96") where ik < "93")))))`,
`(ivc union (cus join ivc)) where ik is "82"`,
`(ivc join ((((((bln union aln) remove a2,b3) where b2 is "45")) where ik is "86")
	union aln))`,
`(cus leftjoin (ivc union ((ivc union ivc) join
	((aln where ik is "23") union aln))))`,
`(((ivc union ivc) union (ivc union ((((((ivc extend x0="0")) where x0 is "11"))
	extend x2) union ivc))) join bln)`,
`((ivc join (aln union aln)) union ((ivc join ((bln extend x0) union aln))
	where a3 is "79"))`,
`(((cus union cus) extend x0) union ((cus union cus) union (cus union cus)))
	remove c4,c1`,
`(((cus join ivc) join (aln extend x0="0")) union ((cus join ivc) join aln))
	remove a4,x0`,
`((ivc union ivc) join (bln union (((aln where ak is "40") union
	(aln where ak is "15")) union ((aln extend x2="2") union aln)))) sort ck,b1`,
`((((((((((ivc union ivc) join bln) union ((ivc union ivc) join bln)) union
	(ivc join aln)) union (ivc join aln)) where ik is "32")) extend x1="1"))
	where bk is "99") sort ck,i3`,
`(((((cus join ivc) join aln) union ((cus join ivc) join bln)) union
	(((cus join ivc) join bln) union ((cus join ivc) join aln))) union
	((((cus join ivc) join bln) union ((cus join ivc) join bln)) union
	(((cus join ivc) join (bln where bk is "18")) union
	(((cus join ivc) join bln) where bk is "11")))) sort b2`,
`(((cus union cus) join ivc) join ((bln union bln) union (((bln union aln)
	where bk is "79") union (bln where bk is "71"))) remove a1) sort b4,b3`,
`(((cus rename c4 to r1) leftjoin ((ivc rename i4 to r2) union ivc)) join
	(aln union aln))`,
`(((((cus extend r0) join ivc) union (cus leftjoin ivc)) join
	(bln where ik is "67")) union ((cus join (ivc rename i1 to r2)) join aln))
	sort c4,b4,b3`,
`cus leftjoin ((ivc where ck is "56" union (ivc where ck is "56")) join
	(((aln union aln) union (bln rename b3 to y1 union aln)) union
	(aln union (bln extend r0))))`,
`(((cus join (ivc union ivc)) remove c1 leftjoin (bln where ik is "18"))
	where ik is "83")`,
`ivc join (cus extend i1 = c3) where ck='3'`,
`(((aln union bln) remove b4 union ((bln union bln) remove b4 union bln)))
	where ik is "22"`,
`((((((cus union cus) join ivc) extend x1 = "1")) extend bk = x1) join
	(((aln union bln) where ak is "56") union aln))`
`((((cus join ivc) where c1 is "93") join (bln extend c1 = b4)) union (((cus join ivc)
	union (((cus where c2 is "24") union cus) join ivc)) join aln)) sort a3,c3,c4`,
`(cus join (ivc join (((bln union ((aln where a2 is "56") union aln)) union
	(bln extend c3 = ik)) where a1 is "71"))) sort i2,a2,b1`,
`(((cus join (ivc join ((aln union bln) union bln))) union ((cus extend b2 = c2) leftjoin
	(ivc join ((aln union bln) union aln) remove b4))) where a4 is "23") sort i2,b4`,
`(((cus join ((aln union aln) leftjoin ivc)) remove i1 union ((cus extend a4 = ck) join
	(ivc join ((aln union bln) where bk is "61")))) union (cus join (ivc join
	(bln union aln)))) sort b1,i3`,
`((cus extend b1 = c1) join ((((ivc leftjoin (aln union bln)) where b2 < "78"))
	where ak is "21"))`,
`(cus join ((ivc join (((aln union (aln extend c3 = a1)) union bln) where bk is "16"))
	union (ivc join aln)))`,
`((((bln extend ck = bk) leftjoin (ivc where ck is "19")) union ((ivc where i1 is "57")
	join aln)) union (ivc join aln)) sort i1,ak`,
`((((cus extend ak = c1)) rename c4 to y2) join (((ivc union (ivc union (ivc union ivc)))
	join (aln union bln)) where bk is "76"))`,
`((bln union aln) leftjoin ((((cus union cus) join ivc) union ((cus extend ik = c1)
	join ivc)) union ((cus extend i1 = c4) join ivc)))`,
`(((cus leftjoin ivc) union ((cus extend ik = c1) join (ivc where i2 < "54"))) union
	((cus join ivc) union (cus leftjoin (ivc where i4 < "51"))))`,
`((((cus union cus) extend a2 = ck) join ((ivc union ivc) join (bln union aln)))
	where bk is "23")`,
`((((cus extend a3 = c1) union cus) join ((ivc union ivc) join
	((bln union aln) where a3 is "56"))) union ((cus union cus) join
	((ivc union ivc) join (aln union bln))))`,
`((ivc union ((ivc union ivc) extend b1 = i4)) join (((aln union (bln extend r0))
	where ak is "20") union aln))`,
`(cus join (((((ivc join aln) where ak is "59")) where ik is "71") union
	(ivc join bln)))`,
`(((ivc leftjoin cus) join ((bln extend i1 = b4) union (bln union aln))) where ik is "57")
	sort c2,c1`,
`(((((cus join ivc) union (cus leftjoin ivc)) extend x0 = ck) join aln) union
	(((cus join ivc) join aln) union (((ivc where ik is "80") leftjoin
	(cus extend b2 = c1)) join bln)))`,
`((((aln leftjoin (ivc leftjoin cus)) union (aln leftjoin (cus join ivc))) union
	(bln leftjoin (((ivc extend bk = i2) leftjoin cus) where ik is "56")))
	extend x2 = "2")`,
`(((ivc leftjoin cus) join (aln union bln)) union (((ivc where ik is "58") leftjoin cus)
	join (bln union (bln extend i3 = b3))))`,
`((((cus union cus) union ((cus union cus) extend i3 = c1)) union (cus union cus)) join
	(((ivc union ivc) join bln) where i3 is "84"))`,
`(((cus extend i2 = c2) join (((ivc union ivc) join bln) where i2 is "66"))
	where b3 is "26")`,
`(((((cus where c1 is "22") union cus) union cus) extend i4 = c4) join ((((ivc union ivc)
	union ivc) where i4 is "30") join bln))`,
`(((cus extend r0 union cus) join ivc) join aln) union (((ivc where ik is '7'
	project ik,i2,i3,ck leftjoin cus) union (cus join (ivc where ik is '7')))
	join (aln where ik is '7'))`,
`((cus union (cus extend a3 = c4)) leftjoin ((ivc join (bln extend x1 = "1")) union
	((ivc join (aln union bln)) where bk is "69")))`,
`(((cus where c3 isnt "54") join ivc) union (((cus join ivc) where i2 <= "98") union
	((ivc leftjoin cus) where c4 is ""))) sort reverse ik`,
`(((cus join ivc) union (cus join ivc)) union ((ivc leftjoin (cus where ck <= "20"))
	where c2 is ""))`,
`(((ivc union ivc) join ((bln union aln) where ik >= "")) union ((aln leftjoin ivc)
	where i4 is ""))`,
`(((ivc union ((ivc union ivc) extend r2)) join bln) union
	(((((ivc union (ivc union ivc)) join aln) where ak isnt "0")) where ik is ""))`,
`(((((cus join ((bln union bln) leftjoin ivc)) union (cus remove c1 leftjoin
	(ivc join ((aln union bln) union aln)))) where bk isnt "58")) where ik is "")`,
`(cus join ((ivc join bln) union (((ivc leftjoin (bln union aln)) union
	(((ivc union ivc) join (bln where ik <= "")) rename b2 to y1))
	where bk isnt "55"))) sort b4`,
`(((cus join ivc) join (((bln extend c1 = ik)) where bk is "82")) union
	((cus join (ivc union ivc)) join bln)) sort b2,c4,i4`,
`(((((cus union cus) join (((ivc union ivc) join aln) where ak > "1"))
	where ik <= "")) extend r2)`,
`(ivc where ik <= "") join by(ik) (aln where ik <= "" and ak > '1')`,
`(((((ivc union ivc) join ((aln union bln) union aln)) where bk >= "38")) where ik <= "")`,
`cus rename c4 to x, x to y remove y`,
`((((((cus union cus) union cus) union (cus union (cus union cus)))
	remove c2,c3) remove c1,c4) where ck is "79")`,
`(((bln union aln) leftjoin ((ivc where ik is "90") leftjoin (((cus extend ik = c2))
	extend x2 = "2"))) union ((cus join ivc) join aln))`,
`cus union cus union cus extend bk = c3
	join ((ivc where ik is "4") join bln) where ck is "" sort ck`,
`(cus where ck is "" extend r1, i3 = c4) join by(ck,i3) ivc`,
`((cus join (((ivc union ivc) join (bln union bln)) union (ivc join (bln union bln))))
	union (((cus extend b1 = c1) leftjoin (ivc join ((bln union aln) extend x1 = "1")))
	where ak > ""))`,
`((cus extend a4 = c3) join (ivc join (((((aln union bln) where b4 is ""))
	where b3 is "66") union bln)))`,
`(cus join ((ivc union ivc) join ((((((aln union aln) where a4 < "")) extend c4 = ik)
	union bln) union bln))) sort a4`,
`((cus extend ak = c1) join (ivc join (((((bln union bln) union (aln union aln)) union
	((aln union bln) union (bln union bln))) where b2 >= "43")) remove a2)) sort c4`,
`((cus extend ak = c3) join (ivc join ((bln union (aln rename a3 to y0))
	where b3 >= "77")))`,
`(cus join ((ivc join (((aln union aln) extend c3 = a4) union (bln union
	(bln rename b2 to y0)))) where bk > "")) sort y0,a1,b1`,
`(((cus union cus) extend x1 = "1") leftjoin
	(((bln where ik is "95") leftjoin ivc) where bk is "25"))`,
`(((bln union aln) leftjoin ((ivc where ik is "12") leftjoin (((cus extend ik = c2))
	extend x2 = "3"))) union ((cus join ivc) join aln))`,
`((cus union (cus union cus)) join ((ivc leftjoin (aln where ak is ""))
	where ik is "44"))`,
`((bln leftjoin ((cus union (cus extend x1 = "1")) leftjoin ivc)) union
	(((ivc where ik is "96") leftjoin (cus extend ik = c3)) join aln))`,
`(cus join ((ivc where i2 > "") join (((((bln union (bln union bln)) extend ck = b1)
	union (aln union bln)) where ck is "") union bln)))`,
`(cus join (ivc join (((bln extend ck = b1) where ck is "") union bln)))`,
`bln join by(ik,b2) ((ivc where ik is "") leftjoin by(ck) (cus extend b2 = c1))`,
`(((((cus union cus) join ivc remove i1,i4) union (cus leftjoin ivc)) union
	((cus extend ik = c4) join ivc)) join (aln union bln)) sort ik,b4,i3`,
`(((cus extend ik = ck) join ivc) union ((cus extend x2 = c1) join
	(ivc union (ivc extend x1 = "12"))))`,
`(cus leftjoin ((((((ivc union ivc) union (ivc union ivc)) join (bln union
	(aln extend ck = a4))) where b2 > "")) where b2 is "12")) sort reverse c4,i3`,
`cus join (ivc join ((bln union aln) where b3 is "3" extend ck = a4))`,
`((cus remove c4 leftjoin (ivc where ik is "3")) join ((aln union (bln union bln))
	where ik is "3"))`,
`((cus leftjoin (ivc where ik is "3")) join ((aln union bln) where ik is "3"))`,
`(((cus union cus) leftjoin ((((ivc union ivc) where ik is "12")) where i1 < ""))
	join bln) sort i1`,
`((((cus join ivc) union (((cus join ivc) extend x0 = "12") union (cus leftjoin ivc)))
	extend ak = x0) join (aln where ak is "12"))`,
`((((cus where c4 is "")) extend b3 = c4) leftjoin (ivc leftjoin (aln union bln)))`,
`(((cus join (ivc where ik is "")) join (((bln extend r0) union bln) where ik is "36"))
	union ((cus join ivc) join (aln union aln))) sort reverse c3`,
`((((ivc extend a3 = ik)) extend b4 = ik) join ((bln union aln) where b4 is "35"))`,
`((((ivc extend a2 = i2)) where i2 is "67") join (((bln union bln) union (aln union aln))
	where b3 >= "69"))`,
`((((cus union (cus union ((cus extend r0) union cus))) extend a4 = r0)) where a4 < "3")`
`(((((ivc union ivc) leftjoin ((cus extend r0) union cus)) extend x1 = r0)) where x1 is "")`
`(ivc join (bln where b1 is "")) leftjoin ((cus union (cus extend r1)) extend b1 = r1)`
`(((((cus extend r0)) extend a3 = c4) union (cus union cus)) where r0 is "") sort reverse c3,ck`
`(cus extend r1, a2 = r1) join by(ck,a2) (((ivc join by(ik) aln) where a2 is '3') union ivc)`
`cus extend r1, a2 = r1 where a2 is ""`
`cus extend r1, a2 = r1 where a2 in ("", "3")`
`cus extend r1, a2 = r1 where a2 >= "" and a2 < "5" // InRange`
`cus extend r1, a2 = r1 where String?(a2)`
)