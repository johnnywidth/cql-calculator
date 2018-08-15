# Calculating the size of a table in Cassandra

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

## How to run

> go run main.go -file example.yaml
