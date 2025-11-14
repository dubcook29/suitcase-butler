# BUTLER Service

BUTLER Service is developed based on Golang (the database is MongoDB and use GORM framework, web framework use Gin).

``` txt
       ____  _   _ _____ _     _____ ____        
      | __ )| | | |_   _| |   | ____|  _ \       
      |  _ \| | | | | | | |   |  _| | |_) |      
      | |_) | |_| | | | | |___| |___|  _ <       
      |____/ \___/  |_| |_____|_____|_| \_\      
                                                 
       github.com/dubcook29/suitcase-butler      
                                                 
           ------  0.0.1(Alpha)  ------          
                                     
```

## Project Setup

```sh
Usage of /suitcase/butler:
  -dbs string
        connection address for mongodb (default "localhost:27017")
  -h string
        host / ipaddress for butler api (default "localhost")
  -p int
        port number for butler api (default 8080)
  -pass string
        password formongodb connection
  -user string
        username for mongodb connection
```

- Start the main program:

```sh
go run .
```

> because it is in Alpha and does not have the intervention of external WMP and internal WMP, it cannot be used normally after startup.

- If you want to use the internal wmp which is still in testing:

```sh
go test -timeout 0 -run ^TestButlerAPIServiceStarter$ github.com/suitcase/butler/api/v1
```

> please modify the relevant database configuration during testing