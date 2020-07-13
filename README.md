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
An easy method to create an always-available postgres database is using the free `hobby-dev` version of [Heroku Postgres](https://devcenter.heroku.com/articles/heroku-postgresql#using-the-cli). The method `ConnectHeroku` will utilize the database URI that you created with the heroku app, so place that in your `.env` file. 

I prefer this method because I can connect to the database from mulitple dev machines that this project is located on.

### Local
I recommend you follow [this](https://medium.com/coding-blocks/creating-user-database-and-adding-access-on-postgresql-8bfcd2f4a91e) tutorial on how to create a database and user profile in postgres to run a local database. The same database that you create with those commands will be the one you put in the `.env` file.

Again, I prefer using Heroku, but if you desire to have more control over your database, use a local instance.

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
$ ./main <mode>
```

`<mode>` can be either `local` or `heroku`, which will define whether or not you're connecting to a local or heroku postgres database.

You will then be prompted to perform an action. The options are as follows

`w` will watch all of the stocks that you have in your database.

`i` will prompt you to add a new stock. First you will be asked for the stock's symbol. Then you will be asked to provide a price target. 

`u` will prompt you to update the price target of a stock in your database.

`r` will prompt you to remove a stock from your database.

`s` will start an http server on port `8080`. This options spawns a goroutine, so you can continue with other options after starting the server.

`q` to quit the program.