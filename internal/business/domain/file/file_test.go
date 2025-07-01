package file_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/machilan1/cruise/internal/business/domain/file"
	"github.com/machilan1/cruise/internal/business/domain/file/stores/filedb"
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
	file *file.Core
}

func newTestSuite(t *testing.T) *testSuite {
	t.Helper()

	log := testhelper.TestLogger(t)
	testDB, _ := testDatabaseInstance.NewDatabase(t, log)

	return &testSuite{
		file: file.NewCore(filedb.NewStore(testDB)),
	}
}

func TestFile_Lifecycle(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ts := newTestSuite(t)

	want := file.File{
		ID:        1,
		Path:      file.MustParsePath("dir/test.jpg"),
		Original:  "test.jpg",
		MimeType:  "image/jpeg",
		Size:      999,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Query by unknown ID, should return not found error
	{
		if _, err := ts.file.QueryByID(ctx, -999); !errors.Is(err, file.ErrNotFound) {
			t.Fatalf("want not found error, got %v", err)
		}
	}

	// Query by unknown path, should return not found error
	{
		if _, err := ts.file.QueryByPath(ctx, file.MustParsePath("unknown/file.jpg")); !errors.Is(err, file.ErrNotFound) {
			t.Fatalf("want not found error, got %v", err)
		}
	}

	// Create a file
	{
		nFl := file.NewFile{
			Path:     want.Path,
			Original: want.Original,
			Size:     want.Size,
		}
		got, err := ts.file.Create(ctx, nFl)
		if err != nil {
			t.Fatal(err)
		}
		checkFile(t, got, want)
	}

	// Read back, should return the saved file
	{
		got, err := ts.file.QueryByID(ctx, want.ID)
		if err != nil {
			t.Fatal(err)
		}
		checkFile(t, got, want)
	}

	// Read back using path, should return the saved file
	{
		got, err := ts.file.QueryByPath(ctx, want.Path)
		if err != nil {
			t.Fatal(err)
		}
		checkFile(t, got, want)
	}
}

func checkFile(t *testing.T, got, want file.File) {
	t.Helper()

	if diff := cmp.Diff(got, want, sqldb.ApproxTime); diff != "" {
		t.Fatalf("mismatch (-got +want):\n%s", diff)
	}
}
