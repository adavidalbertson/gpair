package store

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

const testingDir = "gpair_testing"

var testingPath string

const existingDir = "existing_dir"

var existingDirPath string

const existingFile = "existing_file.txt"

var existingFilePath string

const existingFileContents = "existing file contents"

const forbiddenFile = "forbidden_file.txt"

var forbiddenFilePath string

const forbiddenFileContents = "forbidden file contents"

const forbiddenDir = "forbidden_dir"

var forbiddenDirPath string

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func setUp() {
	testingPath = filepath.Join(os.TempDir(), testingDir)
	err := os.RemoveAll(testingPath)
	check(err)

	existingDirPath = filepath.Join(testingPath, existingDir)
	err = os.MkdirAll(existingDirPath, 0700)
	check(err)

	existingFilePath = filepath.Join(existingDirPath, existingFile)
	f, err := os.Create(existingFilePath)
	check(err)
	defer f.Close()

	_, err = f.WriteString(existingFileContents)
	check(err)

	forbiddenFilePath = filepath.Join(existingDirPath, forbiddenFile)
	ff, err := os.Create(forbiddenFilePath)
	check(err)
	defer ff.Close()

	_, err = ff.WriteString(forbiddenFileContents)
	check(err)

	err = os.Chmod(forbiddenFilePath, 0000)
	check(err)

	forbiddenDirPath = filepath.Join(testingPath, forbiddenDir)
	err = os.MkdirAll(forbiddenDirPath, 000)
	check(err)
}

func TestNewFileStore(t *testing.T) {

	setUp()

	type args struct {
		filename     string
		startDirType int
		dirPath      []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"new file", args{"new_file.txt", ROOT, []string{testingPath, existingDir}}, false},
		{"new dir", args{"new_file.txt", ROOT, []string{testingPath, "new_dir"}}, false},
		{"existing file", args{existingFile, ROOT, []string{testingPath, existingDir}}, false},
		// {"forbidden file", args{forbiddenFile, ROOT, []string{existingDirPath}}, true},
		// {"forbidden dir", args{"new_file.txt", ROOT, []string{forbiddenDirPath}}, true},
		{"invalid dir type", args{"new_file.txt", -1, []string{existingDirPath}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store, err := NewFileStore(tt.args.filename, tt.args.startDirType, tt.args.dirPath...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewFileStore(%s) error = %v, wantErr %v", store.GetPath(), err, tt.wantErr)
				return
			}
		})
	}
}

func Test_fileStore_fileExists(t *testing.T) {

	setUp()

	tests := []struct {
		name    string
		path    string
		want    bool
		wantErr bool
	}{
		{"existing file", existingFilePath, true, false},
		{"nonexistent file", filepath.Join(existingDirPath, "new_file.txt"), false, false},
		{"nonexistent dir", filepath.Join(testingPath, "new_dir", "new_file.txt"), false, false},
		{"forbidden file", forbiddenFilePath, true, false},
		{"forbidden dir", filepath.Join(forbiddenDirPath, "new_file.txt"), false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := &fileStore{tt.path}

			got, err := fs.fileExists()
			if (err != nil) != tt.wantErr {
				t.Errorf("fileStore.fileExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("fileStore.fileExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fileStore_Read(t *testing.T) {

	setUp()

	type fields struct {
		path string
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{"existing file", fields{existingFilePath}, []byte(existingFileContents), false},
		{"new file", fields{filepath.Join(existingDirPath, "new_file.txt")}, nil, false},
		// {"forbidden file", fields{forbiddenFilePath}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := &fileStore{
				path: tt.fields.path,
			}
			got, err := fs.Read()
			if (err != nil) != tt.wantErr {
				t.Errorf("fileStore.Read(%s) error = %v, wantErr %v", fs.GetPath(), err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fileStore.Read(%s) = %v, want %v", fs.GetPath(), got, tt.want)
			}
		})
	}
}

func Test_fileStore_Write(t *testing.T) {

	setUp()

	type fields struct {
		path string
	}
	type args struct {
		jsonBytes []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"existing file", fields{existingFilePath}, args{[]byte("new file contents")}, false},
		{"new file", fields{filepath.Join(existingDirPath, "new_file.txt")}, args{[]byte("new file contents")}, false},
		{"new dir", fields{filepath.Join(testingPath, "new_dir", "new_file.txt")}, args{[]byte("new file contents")}, true},
		
		// The file is overwritten with 0700 permissions, so no error should be thrown
		{"forbidden file", fields{forbiddenFilePath}, args{[]byte("new file contents")}, false},
		{"forbidden dir", fields{filepath.Join(forbiddenDirPath, "new_file.txt")}, args{[]byte("new file contents")}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := &fileStore{
				path: tt.fields.path,
			}
			if err := fs.Write(tt.args.jsonBytes); (err != nil) != tt.wantErr {
				t.Errorf("fileStore.Write() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
