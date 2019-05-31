# Truora Challenge - Backend 

Truora Challenge API Rest with the following technologies: 

- CockroachDB
- GORM 
- Go lang 

First, you must install CockroachDB and setup an insecure (just for quick testing) local cluster 

https://www.cockroachlabs.com/docs/stable/install-cockroachdb-windows.html 

Second, install Go lang 

https://golang.org/doc/install 

With these properly installed, run these commands to install the required dependencies for GORM
```
go get -u github.com/lib/pq # dependency
go get -u github.com/jinzhu/gorm
``` 
Start the built-in CockroachDB
```
cockroach start --insecure
```  
Create the *abc11* user and the *db_troura* database 
```
cockroach sql --insecure
CREATE USER IF NOT EXISTS abc11;
CREATE DATABASE db_troura;
GRANT ALL ON DATABASE db_troura TO abc11;
\q
```
Install the following dependencies before running the project 
```
go get -u github.com/go-chi/chi
go get -u github.com/go-chi/cors
go get -u github.com/go-chi/cors
go get -u github.com/PuerkitoBio/goquery
go get -u github.com/google/go-cmp/cmp
go get -u github.com/likexian/whois-go
```  
Run the project
```
cd main
go run main.go
```
