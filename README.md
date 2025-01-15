# go-flow-migrate
An opensource Golang Migration tool for gophers! It allows you to easily manage your database schema changes, ensuring smooth migrations and rollbacks.

## 📁 Configuration

To get started, you need to create a configuration file to store your database credentials.
	1.	Create a directory named config in your application’s root.
	2.	Inside config, create a file named database.yml with the following structure:
```yaml
host:
port:
user:
password:
dbname:
```
## 🚀 Installation

Add go-flow-migrate to your project using go get:
```go
go get github.com/einnovationlabs/go-flow-migrate@v1.0.0
```
## Starting the tool
```go
flow.Start()
```
Upon starting, you’ll see the following prompt:
```
Welcome to flow migrate.
1. Create a migration file
2. Run Migrations
3. Rollback Migrations

Select an option to proceed:
```

## 🛠️ Creating a Migration File
To create a new migration file, select option 1 from the prompt or directly use the Create(migration name) method:
```go
flow.Create("create users table")
```
This command will create a migration file containing up and down sections for you to define your migration scripts.
Example Migration File:
```yaml
version: 20250108013715
name: Create User Table
description: Users table
up: |
  CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
  );
down: |
  DROP TABLE users;
```
- up: SQL commands to apply the migration (e.g., creating tables or columns).
- down: SQL commands to reverse the migration (e.g., dropping tables or columns).

## 🔼 Running Migrations

After defining your migration scripts, select option 2 from the prompt or directly run all pending migrations using the MigrateUp method:
```go
flow.MigrateUp()
```
This will execute the up scripts in all migration files that haven’t been applied yet.

## 🔽 Rolling Back Migrations

To rollback migrations, select option 3 from the prompt or use the MigrateDown() method:
```go
flow.MigrateDown(step int)
```
- step: Number of migrations to roll back. For example,
  - step = 1 rolls back the most recent migration
  - step = 3 rolls back the last three migrations.


## 🛡️ Best Practices
- Use descriptive names for your migration files (e.g., add users table).
- Test migrations in a staging environment before running them in production.
- Always keep backups of your database before applying significant changes.

## 💬 Support

For questions or issues, please reach out via GitHub Issues.

## ⚖️ License

Flow is open-source software licensed under the MIT License.

Happy Migrating! 🚀

