Quickbeam
========

A simple in-memory key-based byte store


$quickbeam -addr=":12345" &  
$nc localhost 12345  
put	awd	This is my data  
get	awd	  
This is my data  
   
Commands

put\tkey\tdata\n  -   put data into the store at key  
get\tkey\t\n      -   get data from store at key. supports * for all.  
del\tkey\t\n      -   delete data from store at key. supports * for all.  
con\tkey\t\n      -   consume (get and delete) data from store at key. supports * for all.  
cnt\t\n           -   show store count in log  
det\t\n           -   detail entire store in the log  
log\t\n           -   toggle logging  

