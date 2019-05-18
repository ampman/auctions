package main

import "testing"

func TestAuctions(t *testing.T) {

	var received []string

	input1 := []string{
		"10|1|SELL|toaster_1|10.00|20",
		"12|8|BID|toaster_1|7.50",
		"13|5|BID|toaster_1|12.50",
		"15|8|SELL|tv_1|250.00|20",
		"16",
		"17|8|BID|toaster_1|20.00",
		"18|1|BID|tv_1|150.00",
		"19|3|BID|tv_1|200.00",
		"20",
		"21|3|BID|tv_1|300.00",
	}

	expected1 := []string{
		"20|toaster_1|8|SOLD|12.50|3|20.00|7.50",
		"20|tv_1||UNSOLD|0.00|2|200.00|150.00",
	}


	input2 := []string{
		"10|1|SELL|toaster_1|10.00|20",
		"12|8|BID|toaster_1|7.50",
		"13|5|BID|toaster_1|12.50",
		"15|8|SELL|tv_1|250.00|20",
		"16|6|BID|tv_1|255.00",
		"17|8|BID|toaster_1|20.00",
		"18|1|BID|tv_1|150.00",
		"19|3|BID|tv_1|200.00",
		"20",
		"21|3|BID|tv_1|300.00",
	}

	expected2 := []string{
		"20|toaster_1|8|SOLD|12.50|3|20.00|7.50",
		"20|tv_1|6|SOLD|250.00|3|255.00|150.00",
	}


	input3 := []string{
		"10|1|SELL|toaster_1|10.00|20",
		"11|8|SELL|tv_1|250.00|20",
		"12|8|BID|toaster_1|7.50",
		"13|5|BID|toaster_1|12.50",
		"15|6|BID|tv_1|255.00",
		"16|12|BID|tv_1|275.00",
		"17|8|BID|toaster_1|20.00",
		"18|1|BID|tv_1|150.00",
		"19|3|BID|tv_1|200.00",
		"20",
		"21|3|BID|tv_1|300.00",
	}

	expected3 := []string{
		"20|toaster_1|8|SOLD|12.50|3|20.00|7.50",
		"20|tv_1|12|SOLD|255.00|4|275.00|150.00",
	}


	input4 := []string{
		"10|1|SELL|toaster_1|10.00|20",
		"11|8|SELL|tv_1|250.00|20",
		"12|8|BID|toaster_1|7.50",
		"13|5|BID|toaster_1|12.50",
		"15|6|BID|tv_1|255.00",
		"16|12|BID|tv_1|275.00",
		"17|8|BID|toaster_1|20.00",
		"18|1|BID|tv_1|150.00",
		"19|3|BID|tv_1|200.00",
	}

	expected4 := []string{}


	input5 := []string{
		"10|1|SELL|toaster_1|10.00|20",
		"12|8|BID|toaster_1|7.50",
		"13|5|BID|toaster_1|12.50",
		"15|8|SELL|tv_1|250.00|21",
		"16",
		"17|8|BID|toaster_1|20.00",
		"18|1|BID|tv_1|150.00",
		"19|3|BID|tv_1|255.00",
		"20|4|BID|tv_1|200.00",
		"21|3|BID|tv_1|300.00",
	}

	expected5 := []string{
		"20|toaster_1|8|SOLD|12.50|3|20.00|7.50",
		"21|tv_1|3|SOLD|250.00|3|255.00|150.00",
	}


	tt := []struct {
		name     string
		lines    []string
		expected []string
	}{
		{"input1", input1, expected1},
		{"input2", input2, expected2},
		{"input3", input3, expected3},
		{"input4", input4, expected4},
		{"input5", input5, expected5},
	}

	for _, tc := range tt {
		received = nil
		auctions = nil

		t.Run(tc.name, func(t *testing.T) {

			for _, line := range tc.lines {
				expired, _ := parseLine(line)
				if expired != nil {
					received = append(received, expired...)
				}
			}

			if len(tc.expected) != len(received) {
				t.Fatalf("Invalid auction - %v, expected %v, received: %v", tc.name, tc.expected, received)
			}

			for i, auction := range tc.expected {
				if auction != received[i] {
					t.Fatalf("Invalid auction - %v, expected %v, received: %v", tc.name, auction, received[i])
				}
			}
		})
	}
}
