package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

var auctions []Auction

var separator = "|"

const (
	auctionStatusSold   = "SOLD"
	auctionStatusUnsold = "UNSOLD"
)

// ParseLineCallback enforces synchronous callback aka filepath.Walk from standard library
// Used to avoid hardcoding functionality to parse line within readLineByLine function.
type ParseLineCallback func(line string) (expired []string, err error)

func readLineByLine(file string, parseFn ParseLineCallback) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}

	defer f.Close()

	r := bufio.NewReader(f)

	for {
		line, err := r.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		line = strings.TrimSuffix(line, "\n")
		if err != nil {
			return err
		}

		expired, err := parseFn(line)
		if err != nil {
			return err
		}

		if expired != nil {
			for _, auction := range expired {
				fmt.Println(auction)
			}
		}
	}

	return err
}

func parseLine(line string) (expired []string, err error) {
	s := strings.Split(line, separator)
	length := len(s)

	switch {

	case length == 1:
		timestamp, err := strconv.ParseInt(s[0], 10, 64)
		if err != nil {
			return expired, err
		}

		activity := Heartbeat{timestamp}
		expired = timecheckActivities(activity)

	case length == 5:
		timestamp, err := strconv.ParseInt(s[0], 10, 64)
		if err != nil {
			return expired, err
		}

		user_id := s[1]
		action := s[2]
		item := s[3]

		bid_amount, err := strconv.ParseFloat(s[4], 64)
		if err != nil {
			return expired, err
		}

		bid := Bid{timestamp, user_id, action, item, ToCurrency(bid_amount)}

		// check if bid is still applicable
		expired = timecheckActivities(bid)

		// apply bid to relevant active auction
		applyBid(bid)

	case length == 6:
		timestamp, err := strconv.ParseInt(s[0], 10, 64)
		if err != nil {
			return expired, err
		}

		user_id := s[1]
		action := s[2]
		item := s[3]

		reserve_price, err := strconv.ParseFloat(s[4], 64)
		if err != nil {
			return expired, err
		}

		close_time, err := strconv.ParseInt(s[5], 10, 64)
		if err != nil {
			return expired, err
		}

		activity := Item{timestamp, user_id, action, item, ToCurrency(reserve_price), close_time}

		addAuction(activity)

		expired = timecheckActivities(activity)

	default:
		//errors.New("Minimum match not found")
		fmt.Printf("Match not found")
	}

	return expired, err
}

func addAuction(item Item) {
	auction := Auction{
		CloseTime:     item.CloseTime,
		Item:          item.Item,
		UserID:        "",
		Status:        auctionStatusUnsold,
		PricePaid:     0,
		TotalBidCount: 0,
		HighestBid:    0,
		LowestBid:     0,
		reservePrice:  item.ReservePrice,
	}

	// add new item for sale
	auctions = append(auctions, auction)
}

func applyBid(bid Bid) {
	for i := 0; i < len(auctions); i++ {
		auction := &auctions[i]

		if auction.Item != bid.Item {
			continue
		}

		// TODO: verify isOpen vs auction.CloseTime <= bid.getTimestamp()
		if !auction.isOpen(bid.getTimestamp()) {
			continue
		}

		// increase TotalBidCount
		auction.TotalBidCount = auction.TotalBidCount + 1

		// apply user with highest bid to auction
		if bid.BidAmount > auction.reservePrice {
			if bid.BidAmount > auction.HighestBid {
				auction.UserID = bid.UserID
			}
		}

		// set LowestBid
		if auction.LowestBid == 0 {
			auction.LowestBid = bid.BidAmount
		}

		// set LowestBid
		if auction.LowestBid > bid.BidAmount {
			auction.LowestBid = bid.BidAmount
		}

		// if applicable, set PricePaid to second HighestBid before increasing HighestBid
		if bid.BidAmount > auction.reservePrice {
			// if there is only a single valid bid buyer will pay the reserve price
			if auction.TotalBidCount == 1 {
				auction.PricePaid = auction.reservePrice

				// if previous bids all below reserve price, buyer with winning bid
				// will pay the reserve price, since effectively same case as a single valid bid
			} else if auction.HighestBid < auction.reservePrice {
				auction.PricePaid = auction.reservePrice

				// PricePaid to second HighestBid before increasing HighestBid below
			} else if bid.BidAmount > auction.HighestBid {
				auction.PricePaid = auction.HighestBid
			}
		}

		// increase HighestBid
		if bid.BidAmount > auction.HighestBid {
			auction.HighestBid = bid.BidAmount
		}
	}
}

// timecheckActivities checks each auction against Activity.Timestamp to see if auction expired
// and if so, adjusts relevant values
func timecheckActivities(a Activity) (expired []string) {

	timestamp := a.getTimestamp()

	for i := 0; i < len(auctions); i++ {
		auction := &auctions[i]

		if auction.isOpen(timestamp) {
			continue
		}

		// verify Status
		if auction.TotalBidCount >= 1 && auction.PricePaid >= auction.reservePrice {
			auction.Status = auctionStatusSold
		}

		if auction.Status == auctionStatusUnsold {
			auction.PricePaid = 0
			auction.UserID = ""
		}

		expired = append(expired, fmt.Sprint(auction))

		// remove auction from open auctions
		auctions = append(auctions[:i], auctions[i+1:]...)
	}

	return expired
}

func main() {
	flag.Parse()
	if len(flag.Args()) == 0 {
		fmt.Printf("usage: auctions <file1> <file2> ...\n")
		return
	}

	for _, file := range flag.Args() {
		// reset auctions
		auctions = nil

		// helpful when using multiple files
		//fmt.Printf("\nreading file: %s\n", file)

		// pass synchronous callback aka filepath.Walk
		err := readLineByLine(file, parseLine)
		if err != nil {
			fmt.Println(err)
		}
	}
}

// Auction / Output - expected format
type Auction struct {
	CloseTime     int64    // `close_time` - should be a unix epoch of the time the auction finished
	Item          string   // `item` - is the unique string item code.
	UserID        string   // `user_id` - is the integer id of the winning user, or blank if the item did not sell.
	Status        string   // `status` - should contain either "SOLD" or "UNSOLD" depending on the auction outcome.
	PricePaid     Currency // `price_paid` - should be the price paid by the auction winner (0.00 if the item is UNSOLD), as a number to two decimal places
	TotalBidCount int64    // `total_bid_count` - should be the number of valid bids received for the item.
	HighestBid    Currency // `highest_bid` - the highest bid received for the item as a number to two decimal places
	LowestBid     Currency // `lowest_bid` - the lowest bid placed on the item as a number to two decimal places
	reservePrice  Currency // `reserve_price` is a decimal representing the item reserve price in the site's local currency.
}

func (a Auction) String() string {
	return fmt.Sprintf("%v|%s|%v|%s|%s|%v|%s|%s", a.CloseTime, a.Item, a.UserID, a.Status, a.PricePaid, a.TotalBidCount, a.HighestBid, a.LowestBid)
}

func (a Auction) isOpen(timestamp int64) bool {
	return a.CloseTime > timestamp
}

// Item represents users listing item for sale.
// This appears in the format:
// timestamp|user_id|action|item|reserve_price|close_time
type Item struct {
	Timestamp    int64    // `timestamp` will be an integer representing a unix epoch time and is the auction start time,
	UserID       string   // `user_id` is an integer user id
	Action       string   // `action` will be the string "SELL"
	Item         string   // `item` is a unique string code for that item.
	ReservePrice Currency // `reserve_price` is a decimal representing the item reserve price in the site's local currency.
	CloseTime    int64    // `close_time` will be an integer representing a unix epoch time
}

func (i Item) getTimestamp() int64 {
	return i.Timestamp
}

// Bid represents bids on item.
// This will appear in the format:
// timestamp|user_id|action|item|bid_amount
type Bid struct {
	Timestamp int64    // `timestamp` will be an integer representing a unix epoch time and is the time of the bid,
	UserID    string   // `user_id` is an integer user id
	Action    string   // `action` will be the string "BID"
	Item      string   // `item` is a unique string code for that item.
	BidAmount Currency // `bid_amount` is a decimal representing a bid in the auction site's local currency.
}

func (b Bid) getTimestamp() int64 {
	return b.Timestamp
}

// Heartbeat representing a unix epoch time
// These messages may appear periodically in the input to ensure that auctions can be closed
// in the absence of bids, they take the format:
// `timestamp` - an integer representing a unix epoch time
type Heartbeat struct {
	Timestamp int64 // `timestamp` will be an integer representing a unix epoch time
}

func (h Heartbeat) getTimestamp() int64 {
	return h.Timestamp
}

// Activity interface is used to unify input values, where each entry provides timestamp value
// necessary to determine Auction.Status. Activity interface guarantees access to timestamp.
type Activity interface {
	getTimestamp() int64
}

// Currency represents amount in terms of money using properly sized integer type,
// normalized to the lowest possible currency amount
type Currency int64

// ToCurrency converts a float64 to Currency
// e.g. 1.345 to 1.35
func ToCurrency(f float64) Currency {
	return Currency((f * 100) + 0.5)
}

// Float64 converts a Currency to float64
func (m Currency) Float64() float64 {
	x := float64(m)
	x = x / 100
	return x
}

// Multiply safely multiplies a Currency value by a float64, rounding to the nearest amount.
func (m Currency) Multiply(f float64) Currency {
	x := (float64(m) * f) + 0.5
	return Currency(x)
}

// String returns a formatted Currency value
func (m Currency) String() string {
	x := float64(m)
	x = x / 100
	return fmt.Sprintf("%.2f", x)
}
