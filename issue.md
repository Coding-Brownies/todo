## Approach
Un approccio consiste ancora nel memorizzare l'operazione, mentre:

- utilizzare keywords diverse: se si tratta di una Edit chiamarla Update e così via come nell'elenco sotto, e usare l'oldStatus per ripristinare la modifica

- nella tabella Change togliere tutti i campi riferiti ai task () e sostituirli con un campo oldStatus codificandolo in json []byte, il campo serve per salvare tutti i dati dei task modificati

- all'inizio della funzione Undo prelevare l'ultimo ActionID e fare la query per trovare e recuperare l'elenco di change con quell'ActionID (che avrà lunghezza 2 solo nel caso dello swap, e nel caso della swap tratterà come 2 Update)

- creare una funzione ausiliaria Do che accetta un task, ne esegue il marshalling e lo salva nelle tabella di registro delle modifiche Change

Ogni azione di modifica sui task (Edit, Check, Uncheck, Swap, Delete e Create) può essere ricondotta ad un'operazione di Update, Delete o Create:

Edit -> Update

Check -> Update

Uncheck -> Update

Swap -> Update

Delete -> Delete

Create -> Create