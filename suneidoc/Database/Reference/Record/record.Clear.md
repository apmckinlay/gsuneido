<div style="float:right"><span class="builtin">Builtin</span></div>

#### record.Clear

Clears all the record data so the record is as if it was just created with Record().

**Note:** This also clears the internal information that allows a record from a table
to be updated or deleted.   

i.e. New? will return true and Update and Delete will fail after Clear.

**Note**: This is different from [Delete](<../../../Language/Reference/Object/object.Delete.md>)(all:) which keeps the internal information.