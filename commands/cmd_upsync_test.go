package commands

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/alecthomas/assert/v2"
)

func TestUpsync(t *testing.T) {
	testPath, _ := os.MkdirTemp("", "test")
	fsBlobPathPrefix := "fsblob://" + testPath
	createVersionData(t, fsBlobPathPrefix)
	cmd, err := executeCommandLine("upsync", "--source-path", testPath+"/version/v1", "--target-path", fsBlobPathPrefix+"/index/v1.lvi", "--storage-uri", fsBlobPathPrefix+"/storage")
	assert.NoError(t, err, cmd)
	cmd, err = executeCommandLine("upsync", "--source-path", testPath+"/version/v2", "--target-path", fsBlobPathPrefix+"/index/v2.lvi", "--storage-uri", fsBlobPathPrefix+"/storage")
	assert.NoError(t, err, cmd)
	cmd, err = executeCommandLine("upsync", "--source-path", testPath+"/version/v3", "--target-path", fsBlobPathPrefix+"/index/v3.lvi", "--storage-uri", fsBlobPathPrefix+"/storage")
	assert.NoError(t, err, cmd)
}

func TestUpsyncWithLSI(t *testing.T) {
	testPath, _ := os.MkdirTemp("", "test")
	fsBlobPathPrefix := "fsblob://" + testPath
	createVersionData(t, fsBlobPathPrefix)

	cmd, err := executeCommandLine("upsync", "--source-path", testPath+"/version/v1", "--target-path", fsBlobPathPrefix+"/index/v1.lvi", "--storage-uri", fsBlobPathPrefix+"/storage", "--version-local-store-index-path", fsBlobPathPrefix+"/index/v1.lsi")
	assert.NoError(t, err, cmd)
	cmd, err = executeCommandLine("upsync", "--source-path", testPath+"/version/v2", "--target-path", fsBlobPathPrefix+"/index/v2.lvi", "--storage-uri", fsBlobPathPrefix+"/storage", "--version-local-store-index-path", fsBlobPathPrefix+"/index/v2.lsi")
	assert.NoError(t, err, cmd)
	cmd, err = executeCommandLine("upsync", "--source-path", testPath+"/version/v3", "--target-path", fsBlobPathPrefix+"/index/v3.lvi", "--storage-uri", fsBlobPathPrefix+"/storage", "--version-local-store-index-path", fsBlobPathPrefix+"/index/v3.lsi")
	assert.NoError(t, err, cmd)
}

func TestUpsyncWithBrokenLSI(t *testing.T) {
	testPath, _ := os.MkdirTemp("", "test")
	fsBlobPathPrefix := "fsblob://" + testPath
	createVersionData(t, fsBlobPathPrefix)

	data := [11]uint8{8, 21, 141, 3, 1, 4, 124, 213, 1, 23, 123}
	err := os.MkdirAll(fsBlobPathPrefix[9:]+"/storage", 0777)
	assert.NoError(t, err)
	err = ioutil.WriteFile(fsBlobPathPrefix[9:]+"/storage/store.lsi", data[:11], 0644)
	assert.NoError(t, err)

	cmd, err := executeCommandLine("upsync", "--source-path", testPath+"/version/v1", "--target-path", fsBlobPathPrefix+"/index/v1.lvi", "--storage-uri", fsBlobPathPrefix+"/storage", "--version-local-store-index-path", fsBlobPathPrefix+"/index/v1.lsi")
	assert.Error(t, err, cmd)
}

//func TestUpsyncWithGetConfig(t *testing.T) {
//	testPath, _ := os.MkdirTemp("", "test")
//	fsBlobPathPrefix := "fsblob://" + testPath
//	createVersionData(t, fsBlobPathPrefix)
//
//	cmd, err := executeCommandLine("upsync", "--source-path", testPath+"/version/v1", "--target-path", fsBlobPathPrefix+"/index/v1.lvi", "--storage-uri", fsBlobPathPrefix+"/storage", "--get-config-path", fsBlobPathPrefix+"/index/v1.json")
//	if err != nil {
//		t.Errorf("%s: %s", cmd, err)
//	}
//	cmd, err = executeCommandLine("upsync", "--source-path", testPath+"/version/v2", "--target-path", fsBlobPathPrefix+"/index/v2.lvi", "--storage-uri", fsBlobPathPrefix+"/storage", "--get-config-path", fsBlobPathPrefix+"/index/v2.json")
//	if err != nil {
//		t.Errorf("%s: %s", cmd, err)
//	}
//	cmd, err = executeCommandLine("upsync", "--source-path", testPath+"/version/v3", "--target-path", fsBlobPathPrefix+"/index/v3.lvi", "--storage-uri", fsBlobPathPrefix+"/storage", "--get-config-path", fsBlobPathPrefix+"/index/v3.json")
//	if err != nil {
//		t.Errorf("%s: %s", cmd, err)
//	}
//}
//
//func TestUpsyncWithGetConfigAndLSI(t *testing.T) {
//	testPath, _ := os.MkdirTemp("", "test")
//	fsBlobPathPrefix := "fsblob://" + testPath
//	createVersionData(t, fsBlobPathPrefix)
//	cmd, err := executeCommandLine("upsync", "--source-path", testPath+"/version/v1", "--target-path", fsBlobPathPrefix+"/index/v1.lvi", "--storage-uri", fsBlobPathPrefix+"/storage", "--get-config-path", fsBlobPathPrefix+"/index/v1.json", "--version-local-store-index-path", fsBlobPathPrefix+"/index/v1.lsi")
//	if err != nil {
//		t.Errorf("%s: %s", cmd, err)
//	}
//	cmd, err = executeCommandLine("upsync", "--source-path", testPath+"/version/v2", "--target-path", fsBlobPathPrefix+"/index/v2.lvi", "--storage-uri", fsBlobPathPrefix+"/storage", "--get-config-path", fsBlobPathPrefix+"/index/v2.json", "--version-local-store-index-path", fsBlobPathPrefix+"/index/v2.lsi")
//	if err != nil {
//		t.Errorf("%s: %s", cmd, err)
//	}
//	cmd, err = executeCommandLine("upsync", "--source-path", testPath+"/version/v3", "--target-path", fsBlobPathPrefix+"/index/v3.lvi", "--storage-uri", fsBlobPathPrefix+"/storage", "--get-config-path", fsBlobPathPrefix+"/index/v3.json", "--version-local-store-index-path", fsBlobPathPrefix+"/index/v3.lsi")
//	if err != nil {
//		t.Errorf("%s: %s", cmd, err)
//	}
//}
