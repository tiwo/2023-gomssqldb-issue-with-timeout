# QueryContext blocks when the network is down

Issue [microsoft/go-mssqldb#160](https://github.com/microsoft/go-mssqldb/issues/160)

I want to deal with sometimes slow or unreliable network connections, and use QueryContext to limit query time.

After a connection is established (by Ping() or previous queries), and then fails (in my experiments I disconnect the VPN connection to the server), the next QueryContext hangs indefinitely even after the context is cancelled:

    CONNECTIONSTRING='odbc:DRIVER={ODBC Driver 17 for SQL Server};SERVER={10.0.11.13};PORT={11433};DATABASE={ABC};USER ID={johndoe};PASSWORD={secret password};ApplicationIntent=ReadOnly;' \
    go run main.go
    03:10:39 main.go:68: main()
    03:10:39 main.go:31: query_with_timeout("SELECT 'first';")
    03:10:41 main.go:62: query_with_timeout returning "first"
    03:10:41 main.go:91: to reproduce, disconnect from your network now
    03:10:46 main.go:31: query_with_timeout("SELECT 'second';")
    03:10:58 main.go:38: Cancelling query... state of ctx is: context.deadlineExceededError{})
    03:11:29 main.go:84: Closing database
    ^C03:17:57 main.go:75: Received SIGINT; exiting.

Tested with Go 1.21.3, go-mssqldb 1.6.0.
For what it's worth, I am setting db.SetConnMaxLifetime which should be unrelated, and I am making sure the context really is cancelled.

Is the context even used in [Stmt.sendQuery](https://github.com/microsoft/go-mssqldb/blob/e51fa150588f719f052ab5b715cc595f94088bc3/mssql.go#L506), except for logging?
