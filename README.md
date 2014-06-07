Quickbeam
========

A simple in-memory key-based byte store

put\t<key>\t<data>\n  -   put data into the store at key
get\t<key>\t\n        -   get data from store at key. supports * for all.
del\t<key>\t\n        -   delete data from store at key. supports * for all.
con\t<key>\t\n        -   consume (get and delete) data from store at key. supports * for all.
cnt\t\n               -   show store count in log
det\t\n               -   detail entire store in the log
log\t\n               -   toggle logging
