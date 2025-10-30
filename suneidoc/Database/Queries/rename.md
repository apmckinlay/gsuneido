### rename

*query* **rename** *oldname* to *newname* [ , ... ]

The result of rename is simply the input
with the specified column(s) renamed.

rename can be used to control which fields to join on.

rename may also be used to control rules since rules are applied based on column names.

<div style="display: flex; justify-content: space-around; align-items: center;" class="table-style table-full-width">

<div style="flex-shrink: 0;flex-grow: 1;">

| table | column | 
| :---: | :---: |
| 17 | name | 
| 17 | city | 
| 19 | item | 
| 19 | cost |

</div>
<div style="flex-shrink: 0;text-align: center; padding-left: 1em; padding-right: 1em;">

<code>rename table to tblnum,<br />
      column to colname</code>

</div>
<div style="flex-shrink: 0;text-align: center; padding-left: 1em; padding-right: 1em;">

=

</div>
<div style="flex-shrink: 0;flex-grow: 1;">

| tblnum | colname | 
| :---: | :---: |
| 17 | name | 
| 17 | city | 
| 19 | item | 
| 19 | cost |

</div>
</div>