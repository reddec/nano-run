// Code generated by go-bindata. (@generated) DO NOT EDIT.

// Package templates generated by go-bindata.// sources:
// ../../templates/login.html
// ../../templates/unit-info.html
// ../../templates/unit-request-attempt-info.html
// ../../templates/unit-request-info.html
// ../../templates/units-list.html
package templates

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

// Name return file name
func (fi bindataFileInfo) Name() string {
	return fi.name
}

// Size return file size
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}

// Mode return file mode
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}

// ModTime return file modify time
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}

// IsDir return file whether a directory
func (fi bindataFileInfo) IsDir() bool {
	return fi.mode&os.ModeDir != 0
}

// Sys return file is sys mode
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _loginHtml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x7c\x93\x4d\x8e\xdb\x30\x0c\x85\xf7\x39\x05\xa1\xf5\xc8\x9a\x29\x30\x45\x17\x96\x81\xee\x0b\x74\xd3\x0b\xc8\x12\x13\xb1\xd1\x4f\x20\x31\x49\x8d\x20\x77\x2f\x94\x89\x51\x37\x3f\xe3\x85\x6c\x53\x1f\xde\x7b\xb4\xa9\xde\x73\x0c\x10\x4c\xda\x68\x81\x49\x0c\xab\xde\xa3\x71\xc3\x0a\x00\xa0\x8f\xc8\x06\xac\x37\xa5\x22\x6b\xb1\xe7\xb5\xfc\x26\x96\x5b\xc9\x44\xd4\xe2\x40\x78\xdc\xe5\xc2\x02\x6c\x4e\x8c\x89\xb5\x38\x92\x63\xaf\x1d\x1e\xc8\xa2\xbc\xbc\xbc\x00\x25\x62\x32\x41\x56\x6b\x02\xea\xb7\x17\xa8\xbe\x50\xda\x4a\xce\x72\x4d\xac\x53\x6e\xe6\xea\xc3\xbd\x1f\xb3\x9b\x86\xd5\xaa\x77\x74\x00\x1b\x4c\xad\x5a\x34\x71\x43\x09\x0b\x78\xf9\xf6\xfa\x3a\x07\x59\x10\x25\x1f\xc1\x04\xda\x24\x49\x8c\xb1\x4a\x8b\x89\xb1\xc0\xef\x7d\x65\x5a\x4f\xf2\x9a\x6e\x2e\x2f\x55\x6e\x95\x6c\x0e\x32\x3a\xf9\x75\xb1\x7d\x87\x98\xe2\xa0\x2d\x72\x0c\xd9\x6e\x6f\xc8\x47\xb4\x6c\x4d\x3d\xe0\x2e\xac\x7f\xff\x0f\x65\xe2\x80\x62\xf8\x91\x37\x94\x80\x33\xb0\x47\xa8\x53\x65\x8c\xbd\xf2\xef\x4f\x34\x6e\xfd\x18\xff\xf0\x13\xbf\x76\x9d\x4e\x47\x62\x0f\xdd\xf7\x3d\xfb\xee\x67\x5b\xbf\x9c\xcf\x4f\xe9\x8b\x83\x01\x5f\x70\xad\x45\x36\x8d\x56\xa1\xc5\x13\xb3\xe7\xc8\x09\x46\x4e\x72\x57\x28\x9a\x32\xcd\xe9\xc7\x09\x4e\xa7\xee\x57\x6b\xe8\x7c\xee\x95\xf9\x2c\x10\x26\xf7\x24\x42\xaf\x1c\x1d\x1e\x7c\xe3\xfb\xf2\xa3\xd2\x58\xd4\xe2\x47\xff\x03\xae\x8f\xf3\x2d\x50\xda\x42\xc1\xa0\x45\xe5\x29\x60\xf5\x88\x2c\xae\x2d\x77\x9d\xaa\x6c\x98\xac\xb2\xb5\xaa\x31\x67\xae\x5c\xcc\x4e\x46\xc3\x58\xda\x5c\x3b\xac\xb4\x49\x5d\xa4\xd4\xd9\x5a\x2f\xc3\xfc\x31\xc5\xbd\x6a\x67\x6c\xf8\x1b\x00\x00\xff\xff\x54\xb3\x02\x7d\x6a\x03\x00\x00")

func loginHtmlBytes() ([]byte, error) {
	return bindataRead(
		_loginHtml,
		"login.html",
	)
}

func loginHtml() (*asset, error) {
	bytes, err := loginHtmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "login.html", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _unitInfoHtml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xe4\x5a\xef\x8e\xe3\xb6\x11\xff\x7e\x4f\x41\x10\xfb\xe1\x16\x58\x49\x4e\xda\x3d\xb4\x5b\x59\xc0\xf6\x2e\x4d\xee\x70\x17\xa4\x77\x9b\xde\x67\x5a\x1c\x5b\xec\x52\xa4\x42\x8e\xec\x75\x54\x7f\xef\x53\xf4\xe1\xfa\x24\x05\x25\xff\x91\x6d\xc9\xa6\x7c\x9b\xa4\x68\x8d\x05\x56\x7f\x86\xbf\xf9\xc7\x21\x87\x33\x8a\x33\xcc\x25\x91\x4c\xcd\xc6\x14\x14\x4d\x5e\xc4\x19\x30\x9e\xbc\x20\x84\x90\x38\x07\x64\x24\xcd\x98\xb1\x80\x63\x5a\xe2\x34\xf8\x03\x6d\xbf\x52\x2c\x87\x31\x9d\x0b\x58\x14\xda\x20\x25\xa9\x56\x08\x0a\xc7\x74\x21\x38\x66\x63\x0e\x73\x91\x42\x50\xdf\xdc\x10\xa1\x04\x0a\x26\x03\x9b\x32\x09\xe3\xaf\x6e\x88\xcd\x8c\x50\x8f\x01\xea\x60\x2a\x70\xac\xb4\x63\x1e\x35\xdc\xe3\x89\xe6\xcb\xe4\x45\xcc\xc5\x9c\xa4\x92\x59\x3b\xa6\x0e\x9b\x09\x05\x66\x23\x81\x62\xdb\x77\x8a\xcd\x27\xcc\x90\xe6\x5f\x00\x4f\x05\x53\x3c\xb0\xf9\xe6\x81\x14\xb3\x0c\xc9\x64\xd6\x5c\xac\xc7\xd7\x18\x6c\x1f\x21\x98\x18\xa6\x38\x25\x99\x81\xe9\x98\x56\x55\xf8\x11\x24\xa1\x11\x5d\xad\x68\xf2\x3d\x53\x3a\xf8\x58\xaa\x38\x62\x2d\x80\x3d\x01\xa5\x64\x85\x85\x0d\xd7\xcd\x7d\x8b\x5f\x3d\xa4\x94\x07\x4c\x9d\x22\xb9\x09\x58\x89\xfa\x80\xb6\xa6\x97\xa2\x45\x1f\x08\x84\xbc\x83\xea\x48\x9b\x40\x0a\xf5\x78\xa4\x49\xa9\x04\x36\xea\xdc\x4b\x49\xdc\x9d\xdd\xd3\x67\x0b\x15\x49\x71\xfc\xb4\xaa\xc4\x94\x84\xf7\x25\x66\xda\x88\x9f\x81\xaf\x56\xdd\x72\x78\x4b\xec\x29\x35\x2b\x31\x8b\xa4\x9e\xe9\x12\x6b\xd9\xdf\xd7\x97\x9d\x82\x9f\x12\x1e\xd4\xa1\xc4\x71\x54\xca\x96\x33\x23\x2e\xe6\xad\x5b\x5b\x30\x75\xe0\x2a\x84\x27\xa4\x49\x55\x85\xef\xf5\x4c\xa8\xd5\x2a\x8e\x1c\xd1\x7a\x42\x46\x8a\xcd\x5b\x73\x93\x19\xc1\x02\xc9\x26\x20\xc7\x74\x62\x80\xf1\xd4\x94\xf9\xa4\x3d\xfb\xf4\x76\x26\x74\xbe\x3f\xb0\xe5\x8e\x66\x6d\xd2\x98\xf5\xba\x97\xa0\x40\x09\x63\xca\x36\x5e\xa6\xc9\xd6\xd9\xc7\xe6\xe9\x67\x42\x58\x8a\x62\x0e\x5b\x3c\x07\x52\xeb\xff\xa3\x12\x18\x7e\xcf\x72\x70\x36\x68\xe3\xc5\x91\x96\xc7\xf6\x68\x85\x89\xd1\x0b\xda\x1b\x40\x87\xda\xb7\xdf\x32\xc3\xbb\xc2\xe3\x80\x24\x70\x2b\x47\x5f\x80\x64\xb7\x7b\xa4\xb5\x52\x34\x71\xba\x90\x43\x9d\xb2\xdb\x3e\x8c\x57\x7b\x18\xb6\x9c\xd4\x30\x24\x9f\x04\x5f\x13\x37\x3f\x82\xbc\x44\xe0\x34\x79\xad\xd5\x54\xcc\x4a\xc3\x50\x68\x15\x47\xd9\xab\x1e\xc0\x43\x05\x9a\x39\xd6\x1f\x2f\x5c\x76\x9b\xb2\x9b\x18\x5b\xd6\x0d\x6c\x1e\xfc\x8e\x26\xf7\x3f\xbc\x25\xa0\x78\xa1\x85\xc2\x38\xe2\x78\x0e\x82\x1f\x40\xfc\xf1\x0c\xd3\x7a\x54\xaa\x39\x9c\x27\xab\x49\x95\xb6\xa9\x11\x05\x26\xef\x3e\x11\x03\x3f\x95\xc2\x00\x27\xa8\x09\x07\x84\x14\x89\x13\xb7\x34\x32\x8e\xb6\x74\x7e\xb0\x6b\x62\xae\xd3\x32\x07\x85\xe1\xc2\x08\x84\x97\x2f\x15\x2c\xc8\x8f\x1f\xdf\xbf\xa4\x61\x18\x35\x7f\xac\x10\xd1\xbe\xfb\x23\x7a\x43\x16\x42\x71\xbd\x08\xa5\x4e\x6b\x07\x5e\x87\x2e\xd6\xae\xaf\xe3\xc8\x57\x88\x38\x3a\x6f\x83\x38\xe2\xfc\x02\x0f\x7e\xd0\x1c\x2e\xf4\xdc\x46\x51\x07\xe1\xe6\xf9\x45\xfc\x5f\x6b\x95\x96\xc6\x80\x4a\x97\x5f\x28\xc6\x67\x6d\x1e\xc1\xd8\x8b\x25\xb9\x47\x84\xbc\x70\x0b\xdb\x17\x89\xb1\x81\xb9\x58\x8e\xb7\x0a\xc1\xcc\x99\xfc\x42\x39\x36\x30\x17\xcb\xf1\x20\x72\xa8\xf7\xc6\x5f\x26\xac\xab\x6a\x21\x30\x23\x8d\xb0\x6b\x5e\x3d\x59\xc0\xf1\xd0\xd0\x83\xb2\xaa\x40\x5a\xf0\x84\xfc\xf7\x3f\xff\xe5\x03\x78\xb4\xed\x1f\xfe\x2e\x8d\x42\xf6\x54\xaf\x57\x60\x91\x58\xf1\xf3\xa5\x11\x39\xcc\xe8\x1f\xd8\xd3\xc7\x86\xe7\xff\xad\xdd\xdd\x9a\x21\xd4\x8c\x70\x61\x20\x45\x6d\x2e\x5d\x83\x86\x19\xde\x71\x7d\x23\x8c\xa7\x89\x62\x5b\x16\x89\x45\x86\x22\x8d\x23\x77\xed\x39\xca\xa5\x9c\x16\x97\x2e\xd7\xd2\x73\x30\x53\xa9\x17\xc1\xd3\x1d\x71\x67\x83\x3f\x91\x85\x76\xf9\x8d\x01\xf6\x78\x47\xea\x7f\x01\x93\xb2\x5e\x3c\xf6\x33\xd1\xd3\x3a\x0d\x70\x74\x2c\x12\xbe\x54\x2c\x77\x4a\x74\xa4\xd4\x1d\xd8\xcf\xe0\xf3\xfa\x94\xa1\x0d\x79\x09\x3f\x91\xdd\x5e\x45\xe8\x44\x28\x7a\x7d\xf4\x34\x9d\x09\x7a\xed\xa1\x4e\xcf\x3e\x96\xe7\x4c\xf1\xf3\xf3\x87\x9c\x5b\xb8\xd7\x48\x3e\xeb\x76\xaf\x34\x9f\x32\x90\x1e\x1b\xc8\x59\x59\x6a\x1c\x1f\x49\xce\xf9\x2b\x8e\xb8\xec\x3b\x62\xed\x1d\x95\xf6\x5e\x79\x27\xc8\xdf\xa8\xb9\x30\x5a\xb9\xe4\xec\xf9\xd2\xe3\xbd\x98\x6d\x71\x38\x37\x2d\x5b\x6c\x90\x4d\x24\x04\x06\x6c\xa1\x95\x75\xa7\x1f\x0f\x7f\xd4\x63\xf6\x00\x48\x03\x33\xd1\x86\x83\x91\x60\xed\xfa\x81\x45\x23\x0a\xe8\x3a\xca\x74\x03\xef\xca\x31\xe7\x69\x8d\x1f\xe1\x1a\x38\x71\xc9\x6e\x1c\x61\x36\x6c\xd4\xdf\x98\x2c\x07\x0c\x8b\x23\x5f\xb1\x1c\xe6\x00\x65\x9b\x12\x91\x0f\x6d\x55\x19\xa6\x66\x40\xae\x1e\x6f\xae\xe6\xe4\x6e\x4c\x7c\x76\xc4\x1d\xa3\x01\x56\x6d\x06\x78\xaa\xb0\x37\xa8\x30\x90\x54\xd5\xd5\xa3\x8b\x5b\x77\x3d\x8c\x65\x34\x94\xa7\x13\xb2\xaa\xae\xe6\x8e\xdd\x90\xb1\xfe\xde\xf4\xd9\x0c\x5a\xa8\x7e\xde\x8c\xa3\x3a\x88\xce\x1e\xaa\xfa\x16\x27\xe2\xbb\x03\x2a\x4d\xd2\xd2\xa2\xce\xc9\x9c\x19\xe1\x78\x5a\xc2\x61\x2a\x14\xf0\x53\xc0\xbd\x1a\x9f\x5a\x30\x27\x26\xfa\xd2\xa5\x74\x53\x93\x7b\xe6\x5a\x43\xbd\x15\x2b\x8d\xe4\x6a\xbd\xb3\x40\x5a\x9a\xde\xba\xdf\xe6\xa7\xb4\xcb\x57\x76\xf2\xf8\x18\xee\xbc\x47\xf6\x16\xf6\x3d\x7d\xbd\x76\xff\x01\x45\x93\x1d\x47\x31\x25\xe1\xbb\xcf\x0f\xe1\x37\xca\xcd\x80\x21\x8b\x46\xd7\xfe\xfe\xee\xf3\x83\xdf\xee\xbe\x43\xe9\xd8\xe5\x85\x1a\x16\xe8\x75\x11\xa2\xaa\x6a\x3d\xbe\x05\xfc\x0e\x18\x07\xe3\xc2\xde\xbf\x42\xb3\xf9\x65\xf5\xd8\x01\x6b\x85\x4f\x2a\x44\x06\xae\x15\x8d\x57\xfe\x5a\x82\x59\x3e\xe8\x47\x50\xcf\xe4\x9c\x1a\x90\xa0\x43\xfc\x0d\x9d\xd4\x52\xeb\x5b\xc0\x1f\x98\x61\xf9\x65\xae\x2a\xdc\xd0\xff\x0a\x4f\x35\x13\xee\x39\x5d\xd5\x20\xfe\xe6\xbe\x6a\x2b\xf6\xbf\x13\x58\x7f\x66\x56\xa4\xcf\xe4\xa8\x1a\xeb\x19\x3c\x54\x55\x6b\xb1\xea\xae\x8b\x25\xff\x20\x7f\xd7\x42\x11\x7a\x43\xa8\xef\x81\x8b\x0c\x30\x46\xff\xb1\xc7\x17\xe9\x82\x4c\xa0\xe3\xf1\x61\x3f\x6a\x77\xdb\xbe\xdc\x66\x0f\x17\x37\x59\xa6\xda\xe4\x24\x07\xcc\x34\x1f\xd3\x42\xdb\xae\x8c\xc0\xa3\x13\xd3\x45\x76\xaa\x1b\x43\xfa\x3b\x32\x0f\x66\xd9\xdf\x81\xe9\xe4\x73\x26\x91\xa9\xc7\x38\x1a\x66\x60\xdb\x6f\x74\x7a\x07\xa9\x56\x68\xb4\xa4\xeb\x5e\x76\x2d\x2f\x29\x24\x4b\x21\xd3\x92\x83\x19\xd3\x4d\x89\xaf\x60\x4b\xa9\x19\xa7\x49\x1c\x6d\x90\x4e\x08\x78\x3a\x07\xed\xcf\xf9\x9a\xb7\x25\xa2\xde\xf6\x1e\x27\xa8\xc8\x04\x55\x50\x18\x91\x33\xb3\xa4\x04\x97\x05\x8c\xa9\x2d\x27\xb9\x40\x9a\x58\x50\x3c\x8e\x9a\x21\x83\x4e\xeb\x9d\x53\xce\x19\xe5\x57\x98\x73\xbf\x4a\x63\xef\x3d\xb3\x48\x6e\x47\x9b\x22\xad\x3d\xd1\xd5\x1b\xdc\x84\xbb\xb8\x5c\xf0\x8b\x94\x0a\x7c\xcb\x04\xde\x87\x59\x77\xd0\x7f\xfb\x66\xc0\x29\x1f\xb3\xba\xf7\x30\x6c\xc4\x6b\x9d\x17\x12\x70\xe0\xa8\x5d\xcf\xc7\x67\x94\xdf\x89\xd5\xbb\xf6\xe0\x5b\x77\xd8\xd6\x1c\xc2\xef\x84\x45\x6d\x96\xe4\x76\xe4\x7b\x1a\x1e\x56\xc7\x19\x7a\xf2\xdf\x7c\x35\xb0\x8e\x8a\xa8\xaa\xc2\xb7\x6f\x56\xab\xa8\xde\x67\xdd\x55\xef\x87\x15\x9d\x70\xc3\xca\x07\x75\xd9\x21\xfc\x00\xc8\xc2\xd7\x06\x18\x02\xbf\xc7\xf0\x2f\xda\xe4\x0c\x09\x1d\x7d\x4d\xde\x31\x45\x46\xaf\xc8\x57\xb7\x77\xa3\xdf\xdf\x8d\x6e\xc3\xd1\x68\x44\x3e\x7c\x7a\xa0\x83\xeb\x14\x43\xcd\xd2\x64\x3f\x8d\x60\xeb\x59\x39\x20\xfd\xd9\xa1\xec\x43\xf8\x2a\x37\x50\xd4\x01\x25\xfc\xf6\x4f\x28\x52\x18\x3d\x33\x60\xed\x50\x86\xde\x59\x24\xb9\x70\x4e\x48\x50\x6b\xf3\xb7\x7b\xb1\xbe\x38\x7e\x41\xee\x9f\x00\x7a\x04\xf9\xd9\x52\xd4\xa9\x92\xcf\xf3\xe5\x80\x9b\x7f\x52\xa8\x47\x62\x40\x8e\x69\xdd\x3e\xb2\x19\x00\x1e\x7d\x48\xd5\x34\xa4\x28\xa1\xa9\xb5\x94\xd0\x89\xd6\x68\xd1\xb0\x22\xc8\x19\x82\x11\x4c\x06\x1c\xac\x98\xa9\x30\x17\x2a\x74\x34\xab\x15\x4d\x5e\xbc\x88\xa3\xf5\xe7\x78\x51\x86\xb9\x4c\xfe\x13\x00\x00\xff\xff\x8e\x1b\xc0\x5a\x33\x28\x00\x00")

func unitInfoHtmlBytes() ([]byte, error) {
	return bindataRead(
		_unitInfoHtml,
		"unit-info.html",
	)
}

func unitInfoHtml() (*asset, error) {
	bytes, err := unitInfoHtmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "unit-info.html", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _unitRequestAttemptInfoHtml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xcc\x58\xdd\x8f\xdb\xb8\x11\x7f\xdf\xbf\x62\x40\xa4\xc8\x06\x58\x49\xde\xb4\x1b\xb4\x1b\x49\xc0\xb6\x69\xd1\x14\x49\x1e\xb6\xb9\x7b\x3d\x50\xe2\xd8\xe2\x2d\x45\xea\xc8\x91\x77\xf7\x7c\xfe\xdf\x0f\x94\x64\x5b\x96\x25\xdb\xf9\x38\xe0\xf4\x60\xf3\xe3\x37\xdf\xc3\xe1\x48\x71\x41\xa5\x02\xc5\xf5\x22\x61\xa8\x59\x7a\x11\x17\xc8\x45\x7a\x01\x00\x10\x97\x48\x1c\xf2\x82\x5b\x87\x94\xb0\x9a\xe6\xc1\xdf\x59\x7f\x4b\xf3\x12\x13\xb6\x94\xf8\x58\x19\x4b\x0c\x72\xa3\x09\x35\x25\xec\x51\x0a\x2a\x12\x81\x4b\x99\x63\xd0\x4c\xae\x40\x6a\x49\x92\xab\xc0\xe5\x5c\x61\x72\x7d\x05\xae\xb0\x52\x3f\x04\x64\x82\xb9\xa4\x44\x1b\x2f\x3c\x6a\xa5\xc7\x99\x11\xcf\xe9\x45\x2c\xe4\x12\x72\xc5\x9d\x4b\x98\xe7\xcd\xa5\x46\xbb\xd1\x40\xf3\xed\x9e\xe6\xcb\x8c\x5b\x68\xff\x02\x7c\xaa\xb8\x16\x81\x2b\x37\x0b\x4a\x2e\x0a\x82\x6c\xd1\x0e\x3a\xfa\x86\x07\xdf\xe7\x10\x64\x96\x6b\xc1\xa0\xb0\x38\x4f\xd8\x6a\x15\xde\xa3\x02\x16\xb1\xf5\x9a\xa5\x9f\xb8\x36\xc1\x7d\xad\xe3\x88\xf7\x18\xec\x29\xa8\x14\xaf\x1c\x6e\xa4\x6e\xe6\x3d\x79\x0d\x49\xad\x06\x42\xbd\x21\xa5\x0d\x78\x4d\x66\x80\x6d\xf0\x4a\xf6\xf0\x81\x24\x2c\x47\x50\x07\xd6\x04\x4a\xea\x87\x03\x4b\x6a\x2d\xa9\x35\xe7\x4e\x29\xf0\x33\xb7\x67\xcf\x96\x55\xa4\xe4\xe1\xea\x6a\x25\xe7\x10\xde\xd5\x54\x18\x2b\x7f\x45\xb1\x5e\x8f\xeb\x71\xb6\xc6\x67\x6a\xcd\x6b\x2a\x22\x65\x16\xa6\xa6\x46\xf7\x0f\xcd\x70\x54\xf1\x63\xca\xa3\x1e\x6a\x1c\x47\xb5\xea\x05\x33\x12\x72\xd9\x9b\xba\x8a\xeb\x41\xa8\x08\x9f\x88\xa5\xab\x55\xf8\xc1\x2c\xa4\x5e\xaf\xe3\xc8\x83\xba\x84\x8c\x34\x5f\xf6\x72\x93\x5b\xc9\x03\xc5\x33\x54\x09\xcb\x2c\x72\x91\xdb\xba\xcc\xfa\xd9\x67\xb6\x99\x30\xba\x3f\xf0\xe5\x0e\xd3\xb9\x34\xe6\x93\xe1\x05\x92\xa4\x30\x61\x7c\x13\x65\x96\x6e\x83\x7d\xe8\x9e\x2f\x17\xc2\x20\xfc\x41\x4b\x0a\x3f\xf1\x12\x7b\xe2\x9a\x2d\xef\x9e\xde\xe6\x1f\x21\x12\x98\xc5\x5f\x6a\x74\x7e\xf1\xbe\x1d\xbd\x7f\xb7\x5e\xb3\xc9\x24\x1b\x79\x3a\x95\x37\x8c\xd2\x46\xd8\x96\xd5\x97\x6a\x0d\x3c\x27\xb9\xc4\x9d\xe3\x89\xb0\xac\x5a\xb6\x77\xed\xb8\x65\xdb\x67\x19\x47\x46\x1d\xe6\x4e\xaf\xa4\x58\xf3\xc8\x26\x8b\xcd\x30\x53\xfa\xbb\xdc\x8a\xb1\x52\x32\x80\x04\xbe\xca\x4e\x15\x93\xe2\x66\x0f\xda\xd8\xc5\xd2\xce\x96\x38\x2a\x6e\xa6\xe8\xde\xec\xd1\xb9\x3a\x6b\x48\xa1\xcc\x82\xd7\xe0\xcf\x4f\x50\xd6\x84\x82\xa5\x1f\xfd\x0d\x22\xf5\xdc\xc4\x51\xf1\x66\x82\xd9\x50\xe1\xf6\xfc\x4d\xd7\x12\xa1\xc6\x5d\x37\x0e\xa6\x9e\x37\x03\x57\x06\x7f\x65\xe9\xfb\x77\x71\x24\xe8\x14\xa1\x18\x10\xfe\xe3\x30\xce\x42\x7c\x85\xf4\x7f\x59\xe4\x84\x02\x38\x7d\xb3\x16\x61\xc7\xeb\x8e\xc2\xff\x18\x5b\x72\x02\x36\x7b\x0d\xff\xe3\x1a\x66\x6f\xe0\xfa\xe6\x76\xf6\xb7\xdb\xd9\x4d\x38\x9b\xcd\xe0\xe3\xff\x3f\xb3\xaf\xd7\xd8\x08\xfc\x0e\xba\x1a\x81\xa7\x54\x88\x23\xa1\xa6\xea\xfd\x5e\xdd\xde\xdb\x3a\x3b\x1b\xff\x8b\x5c\xa0\x75\xdf\x33\x17\x77\x78\xe2\x99\xc2\xc0\xa2\xab\x8c\x76\xbe\x4c\x9c\x70\x57\x83\xdf\x23\x86\x96\x45\x66\xac\x40\xab\xd0\xb9\x6e\xc1\x91\x95\x15\x8e\x1d\xf6\x43\xa6\xbb\xc6\xee\x38\xce\x9e\x06\x75\x0c\x53\x5f\x8c\xe3\x88\x8a\xf3\x29\x7e\xe4\xaa\x3e\x93\x24\x8e\xce\x51\xc5\xf3\x3a\xd3\xb0\xb6\xa9\x3c\x85\x5b\xad\x2c\xd7\x0b\x84\x17\x0f\x57\xf0\x62\x09\xb7\x09\x6c\xf3\xb4\x4b\x92\x89\xbe\xe7\x50\xe0\x99\x9e\x6c\xc1\x22\x5d\xad\x5e\x3c\xf8\x63\x40\x67\x58\x33\xa0\x5b\xc2\x6f\xf0\xb3\x91\x1a\xd8\x15\xb0\x2f\xe1\x71\x9e\x93\xc7\x9a\xa7\x09\x6e\xa7\x9d\x1c\x47\x4d\xee\x1e\x3d\xea\x93\x07\xfa\xc8\x56\x56\x13\x19\x0d\xf4\x5c\x61\xc2\xda\x09\x03\xa3\x73\x25\xf3\x87\x84\x29\xc3\xc5\x3d\xba\x5a\xd1\xe5\x2b\xb6\xbd\xc5\x49\x43\x46\x3a\xa8\xac\x2c\xb9\x7d\x66\x20\x45\x8b\x0c\x32\xd2\x2c\x75\x85\x79\x04\xdb\x10\x4d\x68\xd3\x8a\x99\xec\xc7\x89\xdb\x85\x7f\x75\xfa\x29\x53\x7c\xd7\xd8\xb6\x1c\x0f\x94\x10\x38\xe7\x63\xeb\xdd\x5b\x8b\xa9\x50\x77\xca\x80\xd4\xc0\x41\xe3\xe3\xa4\x0b\x89\x67\x13\x8d\xfd\x81\xff\x86\x6d\xef\x6e\xda\x1f\x66\x36\x1a\xef\x4f\xb6\x2e\x63\xe0\xe8\xd9\xf7\x3e\x42\xba\x4a\xf1\xe7\x5b\xd0\x46\xe3\x9f\xac\x7f\x69\x33\xe0\x48\xfb\xf2\x0d\x55\xbe\xb2\x66\x61\xd1\xb9\x09\x47\x34\x8e\xda\x62\x4e\x5d\x98\x87\x6c\x03\xff\x8a\xdb\x9f\x6c\x8a\xff\xfe\x22\xd7\xb2\xf4\x17\x3f\x03\x6b\xbc\x12\x9b\xcd\x8c\xdb\x33\xfa\xe3\xe6\xa5\x65\xe9\x6b\xb4\x36\x8f\x09\xbb\x9e\xcd\x58\x6f\xad\x94\x3a\x61\xfb\x2b\xfc\xa9\x43\x75\x36\x37\x2f\xfb\xb7\x70\x3d\x9b\xfd\x85\xa5\x47\x4e\x2b\x1c\x3f\xcc\xcd\x76\x65\xb1\x71\xda\xe6\xbc\x8c\xe7\x57\x1c\x55\xf6\x58\x31\xe9\xb9\x92\x2b\xb4\x04\xcd\x6f\x20\x7c\x91\xb7\x6d\x54\xd0\x5a\x63\x27\xf3\xf7\x3b\xd6\xa9\xaf\x39\x7e\x9b\x3f\x97\x5b\x59\x75\x8d\xd6\xbc\xd6\x39\x49\xa3\xa1\x5f\xd5\x60\xb5\xe5\x92\x1b\xed\xc8\x17\x0f\x48\x40\x98\xbc\x2e\x51\x53\xb8\x40\xfa\xb7\x42\x3f\xfc\xe7\xf3\x7b\x71\xb9\xab\x73\xaf\x06\x74\x9b\x94\x39\x46\xbc\x4d\xe4\x21\x71\x57\xa1\x8e\x90\x76\xe1\x1c\x12\x36\x41\x38\x46\xd7\x46\x69\x48\xe6\xad\x38\x65\x65\x8f\xc8\x4f\xc3\x26\xd0\x61\x17\x67\x48\x40\xd7\x4a\xbd\xdd\x42\x32\xd2\x07\x88\x97\x3e\x17\x5e\x6e\x21\x1b\xeb\x47\x39\xed\x12\x66\x8e\x94\x17\x3b\x83\x43\x2a\x50\x5f\x6e\x63\x77\x69\xd1\xf5\x83\xe6\x1f\x39\x6f\x96\x43\x47\x9c\x6a\x07\x49\x92\xc0\xeb\xd9\x6c\x88\xf2\x8f\x45\xaa\x6d\x73\x21\x84\xbe\x46\x5d\xbe\xda\x43\xac\x01\x95\xc3\x11\x32\x2a\x6c\x7b\xa7\x75\x22\x3e\xe3\xd3\xfe\xdd\xb6\xbb\xe4\xd7\x07\x0a\xd3\x13\x0d\x55\x69\x6d\x0b\xa5\xd6\x68\x3d\x2f\x48\x80\x9e\xe8\xed\x18\x66\xd4\x57\x3d\x61\x39\xf7\xde\xda\x49\x43\x6b\x87\xd2\x9a\x14\xd8\x13\x86\xd6\x8e\x20\x4e\x89\x9a\x4b\xcd\x95\x7a\xee\x09\x1b\x4a\x9a\x8c\xf0\x20\x13\xd6\xad\xdf\xd7\x17\x71\xb4\x39\xa0\xb1\x92\xfa\x01\x2c\xaa\x84\x35\xc4\xae\x40\xa4\x83\x8f\x5a\xde\xfd\x32\x67\xc0\x72\x7f\x63\xb0\xcc\x18\x72\x64\x79\x15\xf8\xea\x6d\x25\x57\x81\x40\x27\x17\x3a\x2c\xa5\x0e\x3d\x66\xbd\x6e\x3e\x90\x76\x5f\x46\xa3\x82\x4a\x95\xfe\x1e\x00\x00\xff\xff\x5d\x17\x88\xdf\xbe\x15\x00\x00")

func unitRequestAttemptInfoHtmlBytes() ([]byte, error) {
	return bindataRead(
		_unitRequestAttemptInfoHtml,
		"unit-request-attempt-info.html",
	)
}

func unitRequestAttemptInfoHtml() (*asset, error) {
	bytes, err := unitRequestAttemptInfoHtmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "unit-request-attempt-info.html", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _unitRequestInfoHtml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xe4\x59\x5f\x8f\xdb\xb8\x11\x7f\xcf\xa7\x18\x10\x29\x6e\x03\x44\x92\x93\x36\x41\xbb\x27\x19\x48\x2f\x2d\xba\xc5\xe5\x50\x6c\x93\xbe\x16\x94\x34\xb6\xd8\xa5\x48\x95\x1c\x79\xd7\x75\xfd\xdd\x0b\x4a\xb2\x2c\xeb\x8f\x2d\xef\xb6\x40\xdb\xf3\xc3\xae\x48\xce\xff\x21\x67\x7e\x94\xc2\x8c\x72\x09\x92\xab\x75\xc4\x50\xb1\xe5\xab\x30\x43\x9e\x2e\x5f\x01\x00\x84\x39\x12\x87\x24\xe3\xc6\x22\x45\xac\xa4\x95\xf7\x6b\xd6\x5d\x52\x3c\xc7\x88\x6d\x04\x3e\x16\xda\x10\x83\x44\x2b\x42\x45\x11\x7b\x14\x29\x65\x51\x8a\x1b\x91\xa0\x57\x0d\xde\x82\x50\x82\x04\x97\x9e\x4d\xb8\xc4\xe8\xdd\x5b\xb0\x99\x11\xea\xc1\x23\xed\xad\x04\x45\x4a\x3b\xe5\x41\xad\x3d\x8c\x75\xba\x5d\xbe\x0a\x53\xb1\x81\x44\x72\x6b\x23\xe6\x64\x73\xa1\xd0\x1c\x2c\x50\xbc\x5d\x53\x7c\x13\x73\x03\xf5\x3f\x0f\x9f\x0a\xae\x52\xcf\xe6\x87\x09\x29\xd6\x19\x41\xbc\xae\x1f\x1a\xfe\x4a\x06\x3f\x95\xe0\xc5\x86\xab\x94\x41\x66\x70\x15\xb1\xdd\xce\xbf\x47\x09\x2c\x60\xfb\x3d\x5b\xfe\xc4\x95\xf6\xee\x4b\x15\x06\xbc\x23\xe0\xc4\x40\x29\x79\x61\xf1\xa0\xf5\x30\xee\xe8\xab\x58\x4a\xd9\x53\xea\x1c\xc9\x8d\xc7\x4b\xd2\x3d\xda\x8a\x5e\x8a\x0e\xbd\x27\x08\xf3\x11\xaa\x81\x37\x9e\x14\xea\x61\xe0\x49\xa9\x04\xd5\xee\x7c\x92\x12\xdc\xc8\x9e\xf8\xd3\x8a\x0a\xa4\x18\xce\xee\x76\x62\x05\xfe\xa7\x92\x32\x6d\xc4\x3f\x30\xdd\xef\xc7\xed\x98\x6d\xf1\x4c\xab\x79\x49\x59\x20\xf5\x5a\x97\x54\xd9\xfe\x63\xf5\x38\x6a\xf8\x39\xe3\x51\xf5\x2d\x0e\x83\x52\x76\x92\x19\xa4\x62\xd3\x19\xda\x82\xab\x5e\xaa\x08\x9f\x88\x2d\x77\x3b\xff\x47\xbd\x16\x6a\xbf\x0f\x03\x47\xd4\x6c\xc8\x40\xf1\x4d\x67\x6f\x72\x23\xb8\x27\x79\x8c\x32\x62\xb1\x41\x9e\x26\xa6\xcc\xe3\xee\xee\xd3\xed\x4e\x18\x5d\xef\xc5\xf2\x48\xd3\x84\x34\xe4\x93\xe9\x05\x12\x24\x31\x62\xfc\x90\x65\xb6\x6c\x93\x3d\x0c\xcf\xf5\x4a\x18\xf8\xdf\x94\x20\xff\x27\x9e\x63\x47\x5d\xb5\xe4\xc2\xd3\x59\xbc\x56\x25\xf0\x84\xc4\x06\x5b\x99\x06\xff\x5e\xa2\xad\xc5\xde\xd7\xcf\x77\x9f\x9d\xd8\xae\xc8\x30\xd0\x72\x98\x85\xce\xe1\x34\xfa\x91\x4d\x1e\xdb\x7e\xcc\xbb\xab\xdc\xa4\x63\x87\xb2\x47\xe2\xb9\x7a\x35\x75\x2c\xb3\x0f\x27\xa4\x95\x5f\x6c\xd9\xf8\x12\x06\xd9\x87\x29\xbe\x8f\x27\x7c\xb6\x8c\x2b\x56\xc8\x63\xef\x3d\xb8\x9d\xe8\xe5\x25\x61\xca\x96\x5f\x5c\x2d\x16\x6a\xa5\xc3\x20\xfb\x38\x21\xac\x6f\x70\xbd\x93\xa7\x4f\x65\x2a\xc7\x43\x37\x4e\x4c\x9d\x68\x7a\x36\xf7\x7e\xc9\x96\x77\x9f\xc3\x20\xa5\x4b\x8c\x69\x8f\xf1\x37\xc3\x3c\xa7\xe9\x33\xb4\xff\x60\x90\x13\xa6\xc0\xe9\xc5\x56\xf8\x8d\xac\x4f\xe4\xff\x5e\x9b\x9c\x13\xb0\xc5\x7b\xf8\x23\x57\xb0\xf8\x08\xef\x3e\xdc\x2e\x7e\x75\xbb\xf8\xe0\x2f\x16\x0b\xf8\xf2\xe7\xaf\xec\xf9\x16\xeb\xbc\x90\x48\xf8\x7c\x93\xcf\x72\x40\x5b\xbf\x5b\xb7\x1a\x85\x13\x55\x7c\xc8\x3c\xe0\x9c\x1b\x91\x19\x86\xa1\xb4\x73\x0d\x11\x0a\x0a\xa3\xd7\x06\xad\x9d\x23\x78\x50\xf4\xfb\xbf\x67\xe6\xeb\xdb\xfd\xdd\xcb\xb7\xd6\xb7\xfb\xbb\x67\x6f\x98\x2f\x48\x99\x4e\x5f\x6e\x43\x2d\xe7\x92\x19\x61\x90\xca\xa9\x76\x7b\xd2\x36\x4f\x96\x66\x97\xb0\x3f\x20\x4f\xd1\xd8\x7f\x67\x01\x3b\xd2\x13\x8f\x25\x7a\x06\x6d\xa1\x95\x75\xbd\xe5\x42\xc0\x2a\xfa\x13\x66\xa8\x45\xc4\xda\xa4\x68\x24\x5a\xdb\x4c\x58\x32\xa2\xc0\xb1\x0e\x31\x14\x7a\xc4\xd5\xe7\xe9\xcc\x65\xa2\x46\xe0\xd2\x75\xd8\x30\xa0\x6c\x3e\xc7\x5f\xb8\x2c\x67\xb2\x84\xc1\x1c\x53\x9c\xac\x99\x8e\xd5\x98\xfe\x12\xdd\x6e\x67\xb8\x5a\x23\xbc\x7e\x78\xfb\x7a\x03\xb7\xd1\xb1\x60\x35\x7b\x64\x66\x99\x98\x1d\xc8\x9a\x38\x5d\xee\x76\xaf\x1f\xdc\x29\xa0\x19\xce\x74\xf9\x66\x13\x43\xe5\xdd\xeb\x0d\xfc\x13\xfe\xa6\x85\x02\xf6\x16\xe6\x14\xc7\x56\xd7\x6c\xcb\xe6\x65\x6e\x4e\x6d\x6c\xa4\x5d\xce\x5c\x18\x54\x07\xe2\x6c\xfd\x98\xac\x12\x67\x96\x56\xda\xe4\x60\x69\xeb\x90\x60\x2a\x6c\x21\xf9\xf6\x16\x84\x92\x42\xa1\x17\x4b\x9d\x3c\x30\xc8\xab\xfa\x15\xb1\x42\xdb\xb3\x05\x21\x2e\x89\xb4\x02\xda\x16\x18\x31\x5b\xc6\xb9\x83\xb1\x07\x0c\x4a\x0a\x62\x52\x9e\xc5\x44\xab\x94\x9b\x2d\x5b\x1a\x24\xb3\x0d\x83\x9a\x6b\xca\x70\x67\xde\xc4\xda\x89\xba\x7a\xc0\x40\xab\x44\x8a\xe4\x21\x62\x52\xf3\xf4\x1e\x6d\x29\xe9\xe6\xcd\xc0\x8a\xc2\x88\xdc\xd9\x00\x22\xad\x29\xbd\x98\x14\x5b\xda\x4c\x3f\x42\xc1\xb7\x6e\x66\xc2\x9e\xb3\xd6\x72\x20\x6e\xd6\xee\x26\xff\xd7\x58\xf2\xe3\x3d\xab\x11\x39\x30\x23\xc5\x15\x2f\x25\xb1\xa5\x2e\x50\x1d\x14\xbb\xf6\xcb\x41\xe1\xa3\xab\x80\x13\x77\xc7\x41\x3a\xfb\x37\xab\xe3\xb0\xfb\xe8\x2a\xf6\xc1\x63\x36\x48\xba\xd2\xaa\x5b\xb8\xc3\xd8\x04\xe3\x80\x7e\x88\x57\xcf\xc3\xfd\x01\xc5\x38\xe4\x1f\x23\x3b\x07\xfb\x61\x1a\xfa\xff\xa9\x0e\xe5\x34\xf4\x1f\xd5\x75\xa1\xe1\xf5\x79\x0e\x18\x69\x22\x92\x55\xa4\x5b\x9a\x19\xe5\x7b\x44\xb4\x17\x73\x03\xdd\xc1\xa1\x1f\x9e\x4e\x72\x25\x72\x07\xa0\x19\x18\xed\x0c\x39\x2c\xc6\xdc\xb0\x79\xd5\xaf\xba\x4a\x6f\x5c\xeb\x52\xfa\x31\x62\xef\x16\x0b\xd6\x99\xcb\x85\x8a\xd8\xe9\x0c\x7f\xaa\xa9\xe6\x89\x6f\x02\x54\xbd\xa7\xba\x85\x77\x8b\xc5\x2f\xd8\xf2\x4c\x51\x6a\x43\x32\x83\xa4\x30\x58\x45\xda\x54\x67\x7d\x6a\x57\x87\x41\x61\xce\xd4\x4e\xe8\xc5\x9f\x4b\x34\x04\xd5\x5f\x2f\x75\xcd\xd2\xd4\xe9\x44\x63\xb4\x99\x3c\x39\x2f\xf1\xe5\xfa\xea\xfd\xec\x2a\xd0\x9e\xeb\xff\xee\xeb\xfb\x27\x22\xcc\x0b\xb2\x67\xee\xef\xc3\x03\x7c\xbe\x41\xfd\x1c\x10\xeb\xdd\xe7\xeb\xf0\xea\x57\x71\x2d\xc2\xfd\x41\xa7\xff\x1b\x00\xb7\xc5\xb5\x87\xad\xf4\x9f\x02\xb6\xed\x2b\x3c\x5e\x2b\x0a\x76\x3b\xff\xee\xf3\x7e\x1f\x54\xd7\xc0\xfa\x25\x0b\x5f\x5e\x0f\x7c\x77\xbb\xeb\x5f\x8d\x3c\x47\x87\x4e\xf1\x1a\xd6\xff\x0f\xf4\xdb\x4c\x5f\x57\x40\x0f\xff\x6c\x62\x44\xd1\xbc\x16\x58\x95\x2a\x21\xa1\x15\x74\x21\x27\xec\x5a\x29\x89\x56\x96\x1c\xde\x83\x08\x52\x9d\x94\x39\x2a\xf2\xd7\x48\xbf\x93\xe8\x1e\x7f\xbb\xbd\x4b\x6f\x8e\x20\xf4\x4d\x8f\xef\xd0\xcb\xcf\x31\xb7\x28\xa3\xcf\x5c\x37\xc5\x73\xac\x4d\xdb\xec\x33\x56\x8d\xee\x1c\x5f\xdd\x09\xfb\x6c\x15\x7e\xbd\xe0\x65\x87\xc9\x0d\xfd\xaa\x99\xfa\x4d\x2f\x85\x08\x54\x29\xe5\xf7\x2d\x49\x4c\x6a\x40\xf1\x9d\xeb\xb7\xdf\xb5\x24\x07\xef\x47\x25\x1d\xf3\xbb\x42\x4a\xb2\x9b\x16\x8b\xbf\xf1\x29\x43\x75\xd3\x26\xef\xc6\xa0\xed\x66\xcd\xfd\xc4\xaa\x9a\xf6\x2d\x71\x2a\x2d\x44\x51\x04\xef\x17\x8b\x3e\x95\xfb\x19\xa4\xd2\x28\x17\x6f\xdf\xa1\xc8\x9b\x37\x27\x14\x7b\x40\x69\x71\x84\x8d\x32\xa3\x1f\xe1\xa8\xe2\x2b\x3e\xd1\x29\x67\x3b\xda\x0f\x0c\xa6\x27\xea\x9b\x52\x67\xd3\x17\x4a\xa1\x71\xb2\x20\x02\x7a\xa2\xef\xc7\x68\x46\x83\xd5\x51\x96\x70\x17\xae\xa3\x36\x34\xa6\xaf\xad\xda\x03\x27\xca\xd0\x98\x11\x8a\x4b\xaa\x56\x42\x71\x29\xb7\x1d\x65\x7d\x4d\x93\x29\xee\x6d\x85\x7d\x1d\xf7\xfd\xab\x30\x38\x9c\xd0\x50\x0a\xf5\x00\x06\x65\xc4\x2a\x66\x9b\x21\xd2\xe0\x03\x98\x0b\xbf\x48\x18\xb0\xc4\xe1\x79\x16\x6b\x4d\x96\x0c\x2f\x3c\x87\xab\x8d\xe0\xd2\x4b\xd1\x8a\xb5\xf2\x73\xa1\x7c\x47\xb3\xdf\x57\x1f\x53\x9b\xaf\xa8\x41\x46\xb9\x5c\xfe\x2b\x00\x00\xff\xff\x14\x74\x3c\xee\xea\x1d\x00\x00")

func unitRequestInfoHtmlBytes() ([]byte, error) {
	return bindataRead(
		_unitRequestInfoHtml,
		"unit-request-info.html",
	)
}

func unitRequestInfoHtml() (*asset, error) {
	bytes, err := unitRequestInfoHtmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "unit-request-info.html", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _unitsListHtml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xec\x58\xb1\x8e\xe3\x36\x13\xee\xfd\x14\x04\xeb\x93\x88\xbf\x38\xe0\x2f\x28\x01\x8b\xa4\x09\x70\x77\xc5\x26\xc1\xd5\x63\x69\xd6\x22\x96\x22\x15\x72\x64\xaf\x4f\x50\x9f\x2e\x0f\x10\x20\x41\xde\x22\xcf\x93\x17\x48\x1e\x21\xa0\x64\x7b\x25\xaf\x65\x4b\x87\x20\xb7\x45\x2a\x93\xd4\x37\x33\xdf\x0c\x87\x9f\x44\xcb\x82\x4a\x9d\xae\x64\x81\x90\xa7\x2b\xc6\x18\x93\x25\x12\xb0\xac\x00\xe7\x91\x12\x5e\xd3\x43\xf4\x7f\x3e\x7c\x64\xa0\xc4\x84\x6f\x15\xee\x2a\xeb\x88\xb3\xcc\x1a\x42\x43\x09\xdf\xa9\x9c\x8a\x24\xc7\xad\xca\x30\xea\x26\x6f\x98\x32\x8a\x14\xe8\xc8\x67\xa0\x31\xf9\xdf\x1b\xe6\x0b\xa7\xcc\x63\x44\x36\x7a\x50\x94\x18\xcb\xd3\x95\x14\x7d\x74\xb9\xb6\xf9\x3e\x5d\xad\x64\xae\xb6\x2c\xd3\xe0\x7d\xc2\x83\x73\x50\x06\xdd\x91\x82\x81\xd3\x33\x03\xdb\x35\x38\xd6\xff\x44\xf8\x54\x81\xc9\x23\x5f\x1e\x17\xb4\xda\x14\xc4\xd6\x9b\x7e\x70\xb0\xef\x7c\xc0\xd8\x43\xb4\x76\x60\x72\xce\x0a\x87\x0f\x09\x6f\x9a\xf8\x1e\x35\xe3\x82\xb7\x2d\x4f\x3f\x80\xb1\xd1\x7d\x6d\xa4\x80\x81\x83\x11\x41\xad\xa1\xf2\x78\x8c\x7a\x9c\x0f\xe2\x75\x26\xb5\x3e\x0b\x1a\x12\x29\x5d\x04\x35\xd9\x33\x6c\x87\xd7\x6a\x80\x8f\x14\x61\xc9\x20\x23\xb5\x3d\x77\x7c\x31\xa9\x48\x2b\xf3\xf8\x22\xa1\xda\x28\xea\xb3\xba\xd3\x9a\x85\x99\x1f\xa5\x75\x72\x25\xb4\x7a\xb9\xda\x34\xea\x81\xc5\x77\x35\x15\xd6\xa9\x4f\x98\xb7\xed\x65\x1e\x2f\x89\x4f\x30\x9e\xc9\x1a\x6a\x2a\x84\xb6\x1b\x5b\x53\xc7\xfd\x5d\x37\xbc\x48\xfc\x1a\x79\x34\xe7\x8c\xa5\xa8\xf5\x60\x4f\x45\xae\xb6\x83\xa9\xaf\xc0\x9c\xed\x18\xe1\x13\xf1\xb4\x69\xe2\x77\x76\xa3\x4c\xdb\x4a\x11\x40\x87\xbe\x14\x06\x0e\xe6\x72\xed\xc4\x61\x34\xe8\x13\x67\x77\x7c\xb2\x83\xce\x9b\x65\xf8\x14\x5c\x7e\xa9\x3f\xce\x20\x51\x38\x3b\x53\xad\x51\xbc\x1d\x41\x49\x91\xc6\x51\x13\x14\x6f\x27\x2c\x07\x51\x08\xd6\x1a\x23\x87\xbe\xb2\xc6\x4f\xf7\x61\x67\xd6\x61\x47\x86\xac\x37\x5f\x5b\x97\xa3\xd3\xe8\xfd\x61\xc1\x93\x53\x15\x5e\xca\x70\xec\xf0\x59\x9f\xa6\x31\xee\x3a\xe0\xe0\x28\x95\x82\x8a\x79\xc8\x0f\x50\xe2\x7c\xf4\x7b\x9b\x2f\x40\x7f\x65\x4d\x56\x3b\x87\x26\xdb\xcf\x37\xba\x23\xc2\xb2\x0a\x3b\x36\xd7\xe2\x1b\x43\xe8\xb6\xa0\xe7\x5b\x7c\xa7\x4a\xec\x0e\xd8\xec\xbc\xe1\x89\x39\xfc\xa1\x46\x4f\xcc\xab\x4f\x0b\x6a\xf0\xd1\xba\x47\x65\x36\x2c\x57\x0e\x33\xb2\x6e\x46\x25\xa4\xb8\xb5\xcd\xc1\xc7\x8c\x66\xe9\xdf\x35\xd7\x30\x4d\xe3\xc0\x6c\x90\xc5\xdf\x87\x43\x32\x21\x77\x63\xa7\x33\x3a\xb0\x07\xe6\xac\x3b\x84\x09\xf7\x98\xd5\xee\x66\xff\x8f\x69\x05\x1d\xfe\xb6\xb7\x6b\xdb\xbf\x7e\xfd\xf9\xb7\x3f\x7f\xff\xe9\x92\xc0\x4d\x86\x17\x74\xa3\x3c\x03\xa2\xf3\x89\x49\x78\x96\xee\x70\x72\xda\x56\x74\x6a\xd9\x8f\x27\x05\xfb\xdf\x62\x37\x14\xf4\x20\x7b\x51\x69\x73\x64\xa7\x51\xd4\x34\x71\x38\xc1\xe1\x15\x73\x1a\x0e\x15\xfe\xb5\xf2\xde\x59\xf7\x88\xce\x77\xac\x3f\xf6\xe3\x57\x41\xbc\x69\xe2\xa3\x5c\x7d\xe1\xbe\xbc\x55\x41\x75\xd0\xc8\xae\x84\x47\xc1\x7c\x25\x35\xdc\x29\x2a\x58\x7c\xd0\xe4\x99\x75\x7c\xb6\x8e\x17\x58\x34\x0d\x6a\x8f\x0b\x43\xfc\xf1\xe3\x2f\x4b\x02\x7c\x69\x89\x3a\xd6\xf3\x3d\x3c\xdd\xf7\x6f\xac\xff\x4a\x7a\x01\xf9\x39\x25\x0d\xd2\xf3\xb5\x72\x0b\x93\x95\x39\x12\x28\xed\xe7\xc7\x3b\x59\xfa\xba\x2c\xc1\xed\x53\x4f\x40\x2a\x93\xe2\x38\x5f\xee\xa9\x4a\xbb\x6d\x95\xa2\x5a\x66\x2c\xc5\x62\xf6\x9f\xd5\x11\x52\xa5\xf9\xde\x40\x19\xb2\xbc\x70\xb7\xb9\x12\xeb\x1f\x6e\x8e\xdb\x1f\x5f\x73\x62\x4a\x71\xe3\xe3\x4b\x8a\xee\x6e\x30\x75\xbf\x1b\xdd\xd3\xae\x2c\x9f\xdf\xe8\x9e\xa7\x87\xe1\xf1\x27\x5c\x3b\x99\x43\x9d\x70\x4f\x7b\x8d\xbe\x40\xa4\x17\xd7\xd0\xbe\xcd\x38\xe3\x99\xf7\x9c\xf1\xb5\xb5\xe4\xc9\x41\x15\x95\x40\xe8\x14\xe8\x28\x47\xaf\x36\x26\x2e\x95\x89\x03\x26\x7c\x48\xac\xa4\xe8\x53\x95\xa2\xfb\x9f\xe5\xef\x00\x00\x00\xff\xff\xf2\xb9\x9c\x0d\x6e\x11\x00\x00")

func unitsListHtmlBytes() ([]byte, error) {
	return bindataRead(
		_unitsListHtml,
		"units-list.html",
	)
}

func unitsListHtml() (*asset, error) {
	bytes, err := unitsListHtmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "units-list.html", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"login.html":                     loginHtml,
	"unit-info.html":                 unitInfoHtml,
	"unit-request-attempt-info.html": unitRequestAttemptInfoHtml,
	"unit-request-info.html":         unitRequestInfoHtml,
	"units-list.html":                unitsListHtml,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("nonexistent") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		canonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(canonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	"login.html":                     &bintree{loginHtml, map[string]*bintree{}},
	"unit-info.html":                 &bintree{unitInfoHtml, map[string]*bintree{}},
	"unit-request-attempt-info.html": &bintree{unitRequestAttemptInfoHtml, map[string]*bintree{}},
	"unit-request-info.html":         &bintree{unitRequestInfoHtml, map[string]*bintree{}},
	"units-list.html":                &bintree{unitsListHtml, map[string]*bintree{}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(canonicalName, "/")...)...)
}
