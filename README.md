[![CircleCI](https://circleci.com/gh/johnnywidth/cql-calculator.svg?style=svg)](https://circleci.com/gh/johnnywidth/cql-calculator) [![GolangCI](https://golangci.com/badges/github.com/golangci/golangci-worker.svg)](https://golangci.com/r/github.com/johnnywidth/cql-calculator) [![codecov](https://codecov.io/gh/johnnywidth/cql-calculator/branch/master/graph/badge.svg)](https://codecov.io/gh/johnnywidth/cql-calculator) [![Go Report Card](https://goreportcard.com/badge/github.com/johnnywidth/cql-calculator)](https://goreportcard.com/report/github.com/johnnywidth/cql-calculator)
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fjohnnywidth%2Fcql-calculator.svg?type=small)](https://app.fossa.io/projects/git%2Bgithub.com%2Fjohnnywidth%2Fcql-calculator?ref=badge_small)

# Calculating the size of a table in Cassandra

### Calculating Partition Size

In order to calculate the size of our partitions, we use the following formula:

**Nv=Nr(Nc−Npk−Ns)+Ns**

The number of values (or cells) in the partition (Nv) is equal to the number of static columns (Ns) plus the product of the number of rows (Nr) and the number of of values per row. The number of values per row is defined as the number of columns (Nc) minus the number of primary key columns (Npk) and static columns (Ns).

In order to determine the size, we use the following formula to determine the size St of a partition:

![](https://github.com/johnnywidth/cql-calculator/raw/master/size-formula.png "Formula")

 - In this formula, ck refers to partition key columns, cs to static columns, cr to regular columns, and cc to clustering columns.
 - The term tavg refers to the average number of bytes of metadata stored per cell, such as timestamps. It is typical to use an estimate of 8 bytes for this value.
 - We recognize the number of rows Nr and number of values Nv from our previous calculations.
 - The sizeOf() function refers to the size in bytes of the CQL data type of each referenced column.

### Install

```
> go get -u github.com/johnnywidth/cql-calculator/cmd/cql-calculator
```

### Example

```
> cql-calculator
> cql-calculator -file $GOPATH/src/github.com/johnnywidth/cql-calculator/cmd/cql-calculator/example.yaml
```

```
> cql-calculator -file generated.yaml -query "CREATE TABLE video (video_id int, email text, name text STATIC, status tinyint, uploaded_at timestamp, PRIMARY KEY (video_id, email))"
> cql-calculator -file generated.yaml
```

### Does not support
 - Parse simple `PRIMARY KEY`

[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fjohnnywidth%2Fcql-calculator.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fjohnnywidth%2Fcql-calculator?ref=badge_large)
