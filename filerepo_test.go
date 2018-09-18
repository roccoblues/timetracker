package main

import (
	"os"
	"reflect"
	"testing"

	"github.com/roccoblues/tt/test"
)

var testContent = []byte("foo bar")

func Test_newFileRepo(t *testing.T) {
	want := &fileRepo{path: "foo/bar"}
	if got := newFileRepo("foo/bar"); !reflect.DeepEqual(got, want) {
		t.Errorf("NewFileRepo() = %v, want %v", got, want)
	}
}

func Test_fileRepo_Read(t *testing.T) {
	existingFile := test.NewFile(t, testContent)
	defer os.Remove(existingFile)

	tests := []struct {
		name    string
		path    string
		want    []byte
		wantErr bool
	}{
		{
			name:    "returns empty for non-existing file",
			path:    test.NonExistingFile(t),
			want:    []byte{},
			wantErr: false,
		},
		{
			name:    "returns data from existing file",
			path:    existingFile,
			want:    testContent,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := newFileRepo(tt.path)
			got, err := f.Read()
			if (err != nil) != tt.wantErr {
				t.Errorf("fileStorage.Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fileStorage.Read() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fileRepo_Write(t *testing.T) {
	existingFile := test.NewFile(t, []byte("blubb blubb"))
	defer os.Remove(existingFile)

	tests := []struct {
		name    string
		path    string
		data    []byte
		wantErr bool
	}{
		{
			name:    "creates new file",
			path:    test.NonExistingFile(t),
			data:    testContent,
			wantErr: false,
		},
		{
			name:    "overwrite existing file",
			path:    existingFile,
			data:    testContent,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := newFileRepo(tt.path)
			if err := f.Write(tt.data); (err != nil) != tt.wantErr {
				t.Errorf("fileStorage.Write() error = %v, wantErr %v", err, tt.wantErr)
			}
			data := test.ReadFile(t, f.path)
			if !reflect.DeepEqual(data, tt.data) {
				t.Errorf("fileStorage.Write() = %v, want %v", data, tt.data)
			}
		})
	}
}
