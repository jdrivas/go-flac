package flac

import (
	"archive/zip"
	"bytes"
	"testing"

	httpclient "github.com/ddliu/go-httpclient"
)

func TestFLACDecode(t *testing.T) {
	zipres, err := httpclient.Begin().Get("http://helpguide.sony.net/high-res/sample1/v1/data/Sample_BeeMoved_96kHz24bit.flac.zip")
	if err != nil {
		t.Errorf("Error while downloading test file: %s", err.Error())
		t.FailNow()
	}
	zipdata, err := zipres.ReadAll()
	if err != nil {
		t.Errorf("Error while downloading test file: %s", err.Error())
		t.FailNow()
	}
	zipfile, err := zip.NewReader(bytes.NewReader(zipdata), int64(len(zipdata)))
	if err != nil {
		t.Errorf("Error while decompressing test file: %s", err.Error())
		t.FailNow()
	}
	if zipfile.File[0].Name != "Sample_BeeMoved_96kHz24bit.flac" {
		t.Errorf("Unexpected test file content: %s", zipfile.File[0].Name)
		t.FailNow()
	}

	flachandle, err := zipfile.File[0].Open()
	if err != nil {
		t.Errorf("Failed to decompress test file: %s", err)
		t.FailNow()
	}

	f, err := ParseBytes(flachandle)
	if err != nil {
		t.Errorf("Failed to parse flac file: %s", err)
		t.Fail()
	}

	metadata := [][]int{
		[]int{0, 34},
		[]int{4, 149},
		[]int{6, 58388},
		[]int{2, 1402},
		[]int{1, 102},
	}

	for i, meta := range f.Meta {
		if BlockType(metadata[i][0]) != meta.Type {
			t.Errorf("Metadata type mismatch: got %d expected %d", meta.Type, metadata[i][0])
			t.Fail()
		}
		if metadata[i][1] != len(meta.Data) {
			t.Errorf("Metadata size mismatch: got %d expected %d", len(meta.Data), metadata[i][1])
			t.Fail()
		}
	}
}