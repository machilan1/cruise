package auction_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/machilan1/cruise/internal/business/domain/auction"
	"github.com/machilan1/cruise/internal/business/domain/auction/stores/auctiondb"
	"github.com/machilan1/cruise/internal/business/sdk/sqldb"
	"github.com/machilan1/cruise/internal/business/sdk/testhelper"
)

var testDatabaseInstance *sqldb.TestInstance

func TestMain(m *testing.M) {
	testDatabaseInstance = sqldb.MustTestInstance()
	defer testDatabaseInstance.MustClose()
	m.Run()
}

type testSuite struct {
	auction *auction.Core
}

func newTestSuite(t *testing.T) *testSuite {
	t.Helper()

	log := testhelper.TestLogger(t)
	testDB, _ := testDatabaseInstance.NewDatabase(t, log)

	aucCore := auction.NewCore(auctiondb.NewStore(testDB))

	return &testSuite{
		auction: aucCore,
	}
}

func TestAuction_Lifecycle(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ts := newTestSuite(t)

	want := auction.Auction{
		ID:        1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Create a new auction
	{
		nAuc := auction.NewAuction{}
		got, err := ts.auction.Create(ctx, nAuc)
		if err != nil {
			t.Fatalf("failed to create auction: %v", err)
		}
		checkAuction(t, got, want)
	}

	var target auction.Auction
	// Read back, should return the saved auction
	{
		got, err := ts.auction.QueryByID(ctx, want.ID)
		if err != nil {
			t.Fatalf("failed to get auction: %v", err)
		}
		checkAuction(t, got, want)
		target = got
	}

	// Update the auction
	{
		uAuc := auction.UpdateAuction{}
		want.UpdatedAt = time.Now()

		got, err := ts.auction.Update(ctx, target, uAuc)
		if err != nil {
			t.Fatalf("failed to update auction: %v", err)
		}
		checkAuction(t, got, want)
	}

	// Read back, should return the updated auction
	{
		got, err := ts.auction.QueryByID(ctx, want.ID)
		if err != nil {
			t.Fatalf("failed to get auction: %v", err)
		}
		checkAuction(t, got, want)
		target = got
	}

	// Delete the auction
	{
		if err := ts.auction.Delete(ctx, target); err != nil {
			t.Fatalf("failed to delete auction: %v", err)
		}
	}

	// Read back, should return an error
	{
		if _, err := ts.auction.QueryByID(ctx, want.ID); !errors.Is(err, auction.ErrNotFound) {
			t.Fatalf("got err: %v, want: %v", err, auction.ErrNotFound)
		}
	}
}

func checkAuction(t *testing.T, got, want auction.Auction) {
	t.Helper()

	if diff := cmp.Diff(got, want, sqldb.ApproxTime); diff != "" {
		t.Errorf("mismatch (-got +want):\n%s", diff)
	}
}
