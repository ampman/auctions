### CHALLENGE:

The program should produce the following expected output, with each line representing the
outcome of a completed auction. This should be written to stdout or a file and be pipe
delimited with the following format:
close_time|item|user_id|status|price_paid|total_bid_count|highest_bid|lowest_bid

##### EXAMPLE:

Input:

    10|1|SELL|toaster_1|10.00|20
    12|8|BID|toaster_1|7.50
    13|5|BID|toaster_1|12.50
    15|8|SELL|tv_1|250.00|20
    16
    17|8|BID|toaster_1|20.00
    18|1|BID|tv_1|150.00
    19|3|BID|tv_1|200.00
    20
    21|3|BID|tv_1|300.00


Output:

    20|toaster_1|8|SOLD|12.50|3|20.00|7.50
    20|tv_1||UNSOLD|0.00|2|200.00|150.00


## SOLUTION:

##### Usage:

Quick run, based on supplied input.txt and expected output:

    $ go run main.go input.txt

Program can handle multiple file inputs:

    $ go run main.go input.txt input2.txt input3.txt input4.txt input5.txt

Alternatively build & run:
    
    $ go build -o auctions
    $ ./auctions <file1> <file2> ...


##### Tests:

Run all tests

    $ go test -v

Run specific test

    $ go test -v -run TestAuctions/TestAuctions/input1