package flac

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
)

// File represents a handler of FLAC file
type File struct {
	Meta   []*MetaDataBlock
	Frames FrameData
}

// Marshal encodes all meta tags and returns the content of the resulting whole FLAC file
func (c *File) Marshal() []byte {
	res := bytes.NewBuffer([]byte{})
	res.Write([]byte("fLaC"))
	for i, meta := range c.Meta {
		last := i == len(c.Meta)-1
		res.Write(meta.Marshal(last))
	}
	res.Write(c.Frames)
	return res.Bytes()
}

// Save encapsulates Marshal and save the file to the file system
func (c *File) Save(fn string) error {
	return ioutil.WriteFile(fn, c.Marshal(), 0644)
}

// ParseBytes accepts a reader to a FLAC stream and returns the final file
func ParseBytes(f io.Reader) (*File, error) {
	res := new(File)

	if err := readFLACHead(f); err != nil {
		return nil, err
	}
	if meta, err := readMetadataBlocks(f); err != nil {
		return nil, err
	} else {
		res.Meta = meta
	}

	res.Frames = readFLACStream(f)
	return res, nil
}

// ParseFile parses a FLAC file
func ParseFile(filename string) (*File, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ParseBytes(f)
}
