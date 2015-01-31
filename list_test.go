package secman

import (
	"reflect"
	"testing"
	"time"

	"github.com/awslabs/aws-sdk-go/aws"
	"github.com/awslabs/aws-sdk-go/gen/s3"
)

func TestListPattern(t *testing.T) {
	utc1 := time.FixedZone("X", -3600)

	fakeS3 := &FakeS3{
		ListResponses: []s3.ListObjectsOutput{
			{
				Contents: []s3.Object{
					{
						Key:          aws.String("secrets/one.aes"),
						ETag:         aws.String(`"etag1"`),
						Size:         aws.Long(1004),
						LastModified: time.Date(2006, 1, 2, 15, 4, 5, 0, utc1),
					},
					{
						Key:          aws.String("secrets/two.aes"),
						ETag:         aws.String(`"etag2"`),
						Size:         aws.Long(1005),
						LastModified: time.Date(2007, 1, 2, 15, 4, 5, 0, utc1),
					},
					{
						Key:          aws.String("secrets/winkle.aes"),
						ETag:         aws.String(`"etag3"`),
						Size:         aws.Long(1006),
						LastModified: time.Date(2008, 1, 2, 15, 4, 5, 0, utc1),
					},
				},
			},
		},
	}

	man := Manager{
		Objects: fakeS3,
		Bucket:  "bucket",
		Prefix:  "secrets/",
	}

	actual, err := man.List("three,one*,two*")
	if err != nil {
		t.Fatal(err)
	}

	expected := []File{
		File{
			Path:         "one",
			LastModified: time.Date(2006, 1, 2, 16, 4, 5, 0, time.UTC),
			Size:         1004,
			ETag:         "etag1",
		},
		File{
			Path:         "two",
			LastModified: time.Date(2007, 1, 2, 16, 4, 5, 0, time.UTC),
			Size:         1005,
			ETag:         "etag2",
		},
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Was %#v\n but expected \n%#v", actual, expected)
	}

	req := fakeS3.ListRequests[0]

	if v, want := *req.Bucket, "bucket"; v != want {
		t.Errorf("Bucket was %q but expected %q", v, want)
	}

	if v, want := *req.Prefix, "secrets/"; v != want {
		t.Errorf("Prefix was %q but expected %q", v, want)
	}
}

func TestListNoPattern(t *testing.T) {
	utc1 := time.FixedZone("X", -3600)

	fakeS3 := &FakeS3{
		ListResponses: []s3.ListObjectsOutput{
			{
				Contents: []s3.Object{
					{
						Key:          aws.String("secrets/one.aes"),
						ETag:         aws.String(`"etag1"`),
						Size:         aws.Long(1004),
						LastModified: time.Date(2006, 1, 2, 15, 4, 5, 0, utc1),
					},
					{
						Key:          aws.String("secrets/two.aes"),
						ETag:         aws.String(`"etag2"`),
						Size:         aws.Long(1005),
						LastModified: time.Date(2007, 1, 2, 15, 4, 5, 0, utc1),
					},
					{
						Key:          aws.String("secrets/winkle.aes"),
						ETag:         aws.String(`"etag3"`),
						Size:         aws.Long(1006),
						LastModified: time.Date(2008, 1, 2, 15, 4, 5, 0, utc1),
					},
				},
			},
		},
	}

	man := Manager{
		Objects: fakeS3,
		Bucket:  "bucket",
		Prefix:  "secrets/",
	}

	actual, err := man.List("")
	if err != nil {
		t.Fatal(err)
	}

	expected := []File{
		File{
			Path:         "one",
			LastModified: time.Date(2006, 1, 2, 16, 4, 5, 0, time.UTC),
			Size:         1004,
			ETag:         "etag1",
		},
		File{
			Path:         "two",
			LastModified: time.Date(2007, 1, 2, 16, 4, 5, 0, time.UTC),
			Size:         1005,
			ETag:         "etag2",
		},
		File{
			Path:         "winkle",
			LastModified: time.Date(2008, 1, 2, 16, 4, 5, 0, time.UTC),
			Size:         1006,
			ETag:         "etag3",
		},
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Was %#v\n but expected \n%#v", actual, expected)
	}

	req := fakeS3.ListRequests[0]

	if v, want := *req.Bucket, "bucket"; v != want {
		t.Errorf("Bucket was %q but expected %q", v, want)
	}

	if v, want := *req.Prefix, "secrets/"; v != want {
		t.Errorf("Prefix was %q but expected %q", v, want)
	}
}