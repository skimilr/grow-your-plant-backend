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

## Step 2: Create a New Database

Once inside the PostgreSQL prompt (postgres=#), create your database with the following command:
