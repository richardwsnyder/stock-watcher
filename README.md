# stock-watcher

## Goals
This project was created to perturb my usage of trading platforms like Robinhood. This platform will notify me whenever a stock hits a specified price target or price floor. With that in place, I can know when to sell or buy a stock without constantly being online watching stocks and their daily movements.

## Prerequisites

### Postgresql
This project uses a postgresql database to store your stock entries. Details on how to set up a database both locally and remotely are in the [Database](#db) section.

### finnhub.io

This project utilizes the _free_ finnhub.io api. In order to make requests, you will need a token from the api. Create an account [here](https://finnhub.io), save the token that they give you.

### Environment file
This project makes use of a `.env` file to manage all of the secret values. These values include your database and email password, database host, and finnhub token. The `.env` file should look like this if you're running your database on Heroku (explained more below):
```
TOKEN=finnhub_token
EUSER=email_address (email address that the notification will be sent from)
EPASS=email_password (password of the sending email address)
TO=send_to (recipient email address)
HEROKU=heroku_uri
```

And like this if you're running the instance locally:
```
HOST=host_url (probably localhost)
USER=postgres_username (I'd recommend user postgres)
PASSWORD=user_password
DBNAME=database_name
TOKEN=finnhub_token
EUSER=email_address (email address that the notification will be sent from)
EPASS=email_password (password of the sending email address)
TO=send_to (recipient email address)
```

## <a name="db"></a>Database

### Heroku
An easy method creating a postgres database is using the free `hobby-dev` version of [Heroku Postgres](https://devcenter.heroku.com/articles/heroku-postgresql#using-the-cli). The method `ConnectHeroku` Will utilize the database URI that you created with the heroku app, so place that in your `.env` file.

### Local
I recommend you follow [this](https://medium.com/coding-blocks/creating-user-database-and-adding-access-on-postgresql-8bfcd2f4a91e) tutorial on how to create a database and user profile in postgres to run a local database. The same database that you create with those commands will be the one you put in the `.env` file

### Stocks Table

The database that you will be using for this project will have a single table: `stocks`. The table columns and types are as follows
```
symbol: VARCHAR(50), NOT NULL, Primary Key
name: VARCHAR(50)
price: double precision
pricetarget: double precision
```

## Install
```
$ go get github.com/richardwsnyder/stock-watcher github.com/joho/godotenv github.com/lib/pq github.com/gorilla/mux
```

This will place the project in your `$GOPATH` in the `src/github.com` directory. 

## Usage
Before you are able to run the executable created below, you must create the `.env` file as described above in the root directory of the project.
```
$ cd $GOPATH/src/github.com/richardwsnyder/stock-watcher
$ touch .env # this will create a .env file
```
Once you have finished editing your `.env` file, you can the move on to building the binary of the project
```
$ cd src/app
$ go build main.go
$ ./main <action>
```

`<action>` can either be `watch` or `insert`.

`watch` will watch all of the stocks that you have in your database

`insert` will prompt you to add a new stock. First you will be asked for the stock's symbol. Then you will be asked to provide a price target. 

`update` will prompt you to update the price target of a stock in your database

`remove` will prompt you to remove a stock from your database

### Heroku vs local database
The default will be to attempt to connect to a Heroku database. If you are running a lcoal postgres database, uncomment the line

`db, err := database.Connect(database.Host, database.Port, database.User, database.Password, database.Dbname)`

And comment out 

`// db, err := database.ConnectHeroku(en.GoDotEnvVariable("HEROKU"))`