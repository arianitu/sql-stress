# Notice

This is unstable and is in heavy development. 

# sql-stress
Stress a SQL server by defining tasks using JSON. It currently supports mysql, postgres, and sqlite. 

# Problem statement

You want to know how your queries scale when you get to 10+ million rows. How long do inserts take? How long do selects take? How big is my table (in terms of disk space) in the worse case? How big are my indexes? (in terms of disk space)

# Example output
    Processing task: task_1.json
	Insert 300,000 random integers
		Avg: 1.308625ms Worst: 140.21674ms Best: 511.712µs Total: 6m32.587658685s 

		Table: my_test_table
			table size: 9 MB, index size: 0 MB, avg row size: 33 bytes, rows: 299730 
			
	Insert 300,000 strings
		Avg: 1.496862ms Worst: 35.539345ms Best: 629.338µs Total: 7m29.058658422s 

		Table: my_test_table_2
			table size: 19 MB, index size: 19 MB, avg row size: 65 bytes, rows: 298730 

Notice that rows is close, but not exact. This is because these stats are pulled from "show table status" to get speedy results.



# Fixtures

Fixtures are basically .sql files that run in order. This is where you would put your table definitions in. You must have a folder named fixtures in the directory that you're running sql-stress.

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
    
### Property: iterations (Integer)
   The number of times to run the query. If you use functions inside values, they're computed for each iteration. Iterations are run in parallel if possible (sql-stress has a worker command line option)
   
### Property: ignore (Bool)
  Skip this task/step





