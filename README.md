# Plant Game Backend

Welcome to the backend setup guide for the Plant Game. This document will walk you through creating a PostgreSQL database, configuring the database connection for your Go application, and setting up the initial database schema.

---

## Table of Contents
- [Step 1: Create the PostgreSQL Database](#step-1-create-the-postgresql-database)
- [Step 2: Set Up the Database Connection](#step-2-set-up-the-database-connection)
- [Step 3: Create the Database Schema](#step-3-create-the-database-schema)
- [Summary](#summary)
- [Next Steps](#next-steps)

---

## Step 1: Create the PostgreSQL Database

### 1. Open the PostgreSQL Command Line Tool
To interact with your PostgreSQL server, open your terminal and use the `psql` command-line tool:

```bash
psql -U your_username
```

Replace your_username with your PostgreSQL username. You may be prompted to enter your password.

### 2: Create a New Database

Once inside the PostgreSQL prompt (postgres=#), create your database with the following command:

```bash
CREATE DATABASE your_db_name;
```

Replace your_db_name with your desired database name.

### 3: Verify the Database Creation

To confirm your new database was created, list all databases:
```bash
\l
```
Look for your_db_name in the output.

### 4: Exit the PostgresSQL Command Line Tool

To exit, type:
```bash
\q
```

## Step 2: Set Up the Database Connection

### 1. Modify Your .env File
Make sure your .env file has the correct database connection string. Use the following format, replacing placeholders with your actual database details:

```bash
DATABASE_URL=host=localhost user=your_username dbname=your_db_name password=your_password sslmode=disable
```
- your_username: Your PostgreSQL username.
- your_db_name: The name of the database you created.
- your_password: Your PostgreSQL password.

**Example .env File**
If your username is postgres, your database name is plant_game, and your password is secret, your .env file should look like this:
```bash
DATABASE_URL=host=localhost user=postgres dbname=plant_game password=secret sslmode=disable
```

## Step 3: Create the Database Schema

### 1. Connect to the Database
To create the necessary table, connect to your database:

```bash
psql -U your_username -d your_db_name
```
Replace your_db_name with the name of the database you created.

### 2. Create the plants Table
Run the following SQL command to set up the plants table:

```bash
CREATE TABLE plants (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    growth_stage VARCHAR(50),
    health_level INT,
    last_watered TIMESTAMP,
    last_fed TIMESTAMP
);
```


## Next Steps

Start your Go application:
```bash
go run main.go
```
If everything is configured correctly, you should be able to hit your API endpoints and interact with the database.

