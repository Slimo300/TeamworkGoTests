
# TeamworkGoTests

Package customerimporter reads from the given csv file and returns a sorted (data structure of your choice) of email domains along with the number of customers with e-mail addresses for each domain.

### Executing program
```
    go run cmd/domaincounter/domaincounter.go <PATH_TO_FILE>
```

By default program will look for column named "email" to read emails from it. If column name is different it can be changed with '-column' flag
```
    go run cmd/domaincounter/domaincounter.go -column=<EMAIL_COLUMN> <PATH_TO_FILE> 
``` 

If you want to save program's results to a file instead of printing them to standard output (default option) use -'output-file' flag

```
    go run cmd/domaincounter/domaincounter.go -output-file=<FILE_PATH> <PATH_TO_FILE> 
```

Program currently supports 2 formats in which data can be displayed - yaml and json. YAML is default option. Use format flag if you want to change it.

```
    go run cmd/domaincounter/domaincounter.go -format=json <PATH_TO_FILE>
    go run cmd/domaincounter/domaincounter.go -format=yaml <PATH_TO_FILE>
```

### Run tests

```
    go test ./...
```



