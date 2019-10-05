[![CircleCI](https://circleci.com/gh/johnnywidth/cql-calculator.svg?style=svg)](https://circleci.com/gh/johnnywidth/cql-calculator) [![GolangCI](https://golangci.com/badges/github.com/golangci/golangci-worker.svg)](https://golangci.com/r/github.com/johnnywidth/cql-calculator) [![codecov](https://codecov.io/gh/johnnywidth/cql-calculator/branch/master/graph/badge.svg)](https://codecov.io/gh/johnnywidth/cql-calculator) [![Go Report Card](https://goreportcard.com/badge/github.com/johnnywidth/cql-calculator)](https://goreportcard.com/report/github.com/johnnywidth/cql-calculator)
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fjohnnywidth%2Fcql-calculator.svg?type=small)](https://app.fossa.io/projects/git%2Bgithub.com%2Fjohnnywidth%2Fcql-calculator?ref=badge_small)

# Calculating the size of a table in Cassandra

### About project
Cassandra is known for splitting data into the partitions and has a few limitations as to the performance, therefore, it's critical to keep the size of partitions small.

To solve that problem partially, you can run this tool which is used to determine the size of partitions in order to anticipate the needed disk space.

To run the program, you'll need to specify the `CREATE TABLE` query and indicate the assumed number of rows that will appear in the table. In the end, you'll get the exact partition size and the amount of cells.

### The benefits
This tool is a quick solution to change the schema design and apply different assumptions to the various data types. Overall, it helps to keep Cassandra productive and avoid any performance bugs or complications in production mode.

## Calculating Partition Size

In order to calculate the size of our partitions, we use the following formula:

**Nv=Nr(Nc−Npk−Ns)+Ns**

The number of values (or cells) in the partition (Nv) is equal to the number of static columns (Ns) plus the product of the number of rows (Nr) and the number of of values per row. The number of values per row is defined as the number of columns (Nc) minus the number of primary key columns (Npk) and static columns (Ns).

In order to determine the size, we use the following formula to determine the size St of a partition:

![](https://github.com/johnnywidth/cql-calculator/raw/master/size-formula.png "Formula")

 - In this formula, ck refers to partition key columns, cs to static columns, cr to regular columns, and cc to clustering columns.
 - The term tavg refers to the average number of bytes of metadata stored per cell, such as timestamps. It is typical to use an estimate of 8 bytes for this value.
 - We recognize the number of rows Nr and number of values Nv from our previous calculations.
 - The sizeOf() function refers to the size in bytes of the CQL data type of each referenced column.

## Install

```
$ go get -u github.com/johnnywidth/cql-calculator/cmd/cql-calculator
```

## Examples

#### Run for CREATE TABLE query
```
$ cql-calculator -query "CREATE TABLE video ( \
    video_id int, email text, name text STATIC, \
    status tinyint, uploaded_at timestamp, \
    PRIMARY KEY (video_id, email))"

# Output
Enter rows count per one partition: 10000
Enter (avarage) size for `email (text)` column: 150
Enter (avarage) size for `name (text)` column: 250
Number of Values:
(10000 * (5 - 2 - 1) + 1) = 20001

Partition Size on Disk:
(4 + 250 + (10000 * 309) + (8 * 20001)) = 3250262 bytes (3.10 Mb)
```

#### Run for CREATE TABLE query and save to file
```
$ cql-calculator -file generated.yaml -query "CREATE TABLE video ( \
    video_id int, email text, name text STATIC, \
    status tinyint, uploaded_at timestamp, \
    PRIMARY KEY (video_id, email))"

# Output
Enter rows count per one partition: 10000
Enter (avarage) size for `email (text)` column: 150
Enter (avarage) size for `name (text)` column: 250
Number of Values:
(10000 * (5 - 2 - 1) + 1) = 20001

Partition Size on Disk:
(4 + 250 + (10000 * 309) + (8 * 20001)) = 3250262 bytes (3.10 Mb)
```

```
$ cql-calculator -file generated.yaml

# Output
Number of Values:
(10000 * (5 - 2 - 1) + 1) = 20001

Partition Size on Disk:
(4 + 250 + (10000 * 309) + (8 * 20001)) = 3250262 bytes (3.10 Mb)
```

## TODO
 - Parsing simple `PRIMARY KEY`: `CREATE TABLE video (video_id int PRIMARY KEY, email text)`

[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fjohnnywidth%2Fcql-calculator.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fjohnnywidth%2Fcql-calculator?ref=badge_large)
