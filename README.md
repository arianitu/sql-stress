# Notice

This is unstable and is in heavy development. 

# sql-stress
Stress a SQL server by defining tasks using JSON. It currently supports mysql, postgres, and sqlite. 

# Problem statement

You want to know how your queries scale when you get to 10+ million rows. How long do inserts take? How long do selects take? How big is my table (in terms of disk space) in the worse case? How big are my indexes? (in terms of disk space)

# Example output

	Insert 1,000,000 random integers
		Qps: 16970.22 Avg: 583.479µs Worst: 62.636873ms Best: 279.093µs  

		Table: my_test_table
			table size: 37 MB, index size: 22 MB, avg row size: 36 bytes, rows: 1016610 
	Insert 1,000,000 strings
		Qps: 16552.41 Avg: 595.991µs Worst: 49.696598ms Best: 302.531µs 

		Table: my_test_table_2
			table size: 69 MB, index size: 75 MB, avg row size: 69 bytes, rows: 1001345 


Notice that rows is close, but not exact. This is because table stats are pulled from "show table status" (Note: ANALYZE TABLE is run on the table before SHOW TABLE STATUS for accuracy)



# Fixtures

Fixtures are basically .sql files that run in order. This is where you would put your table definitions in. You must have a folder named fixtures in the directory that you're running sql-stress (if you don't want to run fixtures, use the command line option -run-fixtures=0)

The file names inside the fixtures folder should be in the format: name_order. Name must be a string, and order must be an integer.

Example:

    ./fixtures
    ./fixtures/players_1.sql
    ./fixtures/items_2.sql
  
It would run players_1.sql, then items_2.sql. 

Queries inside .sql files must be separated using a semicolon. _All_ queries must have a semicolon, including the last one (or it will not be executed)

# Tasks

Tasks are defined using JSON. Each task has a series of steps that run in order. You can have multiple tasks and each task will run in order.

You must have a folder named tasks in the directory that you're running sql-stress. The file names inside the fixtures folder should be in the format: name_order. Name must be a string, and order must be an integer.

Example:

    ./tasks
    ./tasks/players_1.json
    ./tasks/items_2.json
    
Here is an example of a task:

    {
		"conn": {
			"vendor": "mysql",
			"url": "root:@/sql_stress_test"
		},
    	"steps": [{
    		"tables": ["my_test_table"],
    		"name": "Insert 1,000,000 random integers",
    		"query": "INSERT INTO my_test_table (x) VALUES(?)",
    		"values": ["randIntInclusive(0, 100)"],
    		"iterations": 1000000
    	}, {
    		"tables": ["my_test_table_2"],
    		"name": "Insert 1,000,000 strings",
    		"query": "INSERT INTO my_test_table_2 (x) VALUES(?)",
    		"values": ["randString(10, 50)"],
    		"iterations": 100000
    	}]
    }

# Task Documentation

### Property: skip (Bool)
  Skip this task
  
### Property: conn (Object)

If you don't provide conn, the connection info that is passed via command line arguments is used by default.

#### vendor (String): 
- mysql
- postgres
- sqlite

#### url (String):
- mysql: username:password@localhost/dbname
- postgres: postgres://username:password@localhost/dbname
- sqlite: /some/location/test.db

#### maxOpenConn(Int): 

The maximum connections that can be opened to the sql server

#### workers(Int): 

The maximum workers to spawn. Generally if you want to test lock contention, you want maxOpenConn == workers to get the right amount of connections to SQL

## Steps (Array\<Object\>)

### Property: tables (Array)
  Tables to output metrics for when a step is completed. 
  
### Property: query (String)
  Prepared statement to execute. MySQL tends to use ?, and Postgres tends to use $1,$2..
  
### Property: values (Array)
  Valid values in the array are string, float64, bool and functions*. If you need to use NOW(), do it in the query statement. 
  
##### Functions (String)
  You can supply functions in the values array when you need random data. Functions that currently exist are:
  
    randIntInclusive(min, max)
    randString(minStringLength, maxStringLength)
    incrementingCount(initialCount, increment) Increment can be negative to count downwards,
    count is unique per value in a query.
    
### Property: iterations (Integer)
   The number of times to run the query. If you use functions inside values, they're computed for each iteration. Iterations are run in parallel if possible (sql-stress has a worker command line option)
   
### Property: skip (Bool)
  Skip this step

### Property: delay (Int)
  Time to sleep in miliseconds after this step finishes executing

### Property: predelay (Int)
  Time to sleep in miliseconds before this step starts executing  






