package commands

import (
	"os"
	"path"
	"testing"

	"github.com/DanEngelbrecht/golongtail/longtailutils"
	"github.com/alecthomas/assert/v2"
)

func TestDownsync(t *testing.T) {
	testPath, _ := os.MkdirTemp("", "test")
	fsBlobPathPrefix := "fsblob://" + testPath
	createVersionData(t, fsBlobPathPrefix)
	executeCommandLine("upsync", "--source-path", testPath+"/version/v1", "--target-path", fsBlobPathPrefix+"/index/v1.lvi", "--storage-uri", fsBlobPathPrefix+"/storage")
	executeCommandLine("upsync", "--source-path", testPath+"/version/v2", "--target-path", fsBlobPathPrefix+"/index/v2.lvi", "--storage-uri", fsBlobPathPrefix+"/storage")
	executeCommandLine("upsync", "--source-path", testPath+"/version/v3", "--target-path", fsBlobPathPrefix+"/index/v3.lvi", "--storage-uri", fsBlobPathPrefix+"/storage")

	cmd, err := executeCommandLine("downsync", "--source-path", fsBlobPathPrefix+"/index/v1.lvi", "--target-path", testPath+"/version/current", "--storage-uri", fsBlobPathPrefix+"/storage")
	assert.NoError(t, err, cmd)
	validateContent(t, fsBlobPathPrefix, "version/current", v1FilesCreate)
	cmd, err = executeCommandLine("downsync", "--source-path", fsBlobPathPrefix+"/index/v2.lvi", "--target-path", testPath+"/version/current", "--storage-uri", fsBlobPathPrefix+"/storage")
	assert.NoError(t, err, cmd)
	validateContent(t, fsBlobPathPrefix, "version/current", v2FilesCreate)
	cmd, err = executeCommandLine("downsync", "--source-path", fsBlobPathPrefix+"/index/v3.lvi", "--target-path", testPath+"/version/current", "--storage-uri", fsBlobPathPrefix+"/storage")
	assert.NoError(t, err, cmd)
	validateContent(t, fsBlobPathPrefix, "version/current", v3FilesCreate)
}

func TestDownsyncNoTargetPath(t *testing.T) {
	testPath, _ := os.MkdirTemp("", "test")
	fsBlobPathPrefix := "fsblob://" + testPath
	createVersionData(t, fsBlobPathPrefix)
	executeCommandLine("upsync", "--source-path", testPath+"/version/v1", "--target-path", fsBlobPathPrefix+"/index/v1b.lvi", "--storage-uri", fsBlobPathPrefix+"/storage")
	executeCommandLine("upsync", "--source-path", testPath+"/version/v2", "--target-path", fsBlobPathPrefix+"/index/v2b.lvi", "--storage-uri", fsBlobPathPrefix+"/storage")
	executeCommandLine("upsync", "--source-path", testPath+"/version/v3", "--target-path", fsBlobPathPrefix+"/index/v3b.lvi", "--storage-uri", fsBlobPathPrefix+"/storage")

	popd, _ := os.Getwd()
	defer os.Chdir(popd)
	os.Chdir(path.Join(testPath, "version"))
	cmd, err := executeCommandLine("downsync", "--source-path", fsBlobPathPrefix+"/index/v1b.lvi", "--storage-uri", fsBlobPathPrefix+"/storage")
	assert.NoError(t, err, cmd)
	validateContent(t, fsBlobPathPrefix, "version/v1b", v1FilesCreate)
	cmd, err = executeCommandLine("downsync", "--source-path", fsBlobPathPrefix+"/index/v2b.lvi", "--storage-uri", fsBlobPathPrefix+"/storage")
	assert.NoError(t, err, cmd)
	validateContent(t, fsBlobPathPrefix, "version/v2b", v2FilesCreate)
	cmd, err = executeCommandLine("downsync", "--source-path", fsBlobPathPrefix+"/index/v3b.lvi", "--storage-uri", fsBlobPathPrefix+"/storage")
	assert.NoError(t, err, cmd)
	validateContent(t, fsBlobPathPrefix, "version/v3b", v3FilesCreate)
}

func TestDownsyncWithVersionLSI(t *testing.T) {
	testPath, _ := os.MkdirTemp("", "test")
	fsBlobPathPrefix := "fsblob://" + testPath
	createVersionData(t, fsBlobPathPrefix)

	executeCommandLine("upsync", "--source-path", testPath+"/version/v1", "--target-path", fsBlobPathPrefix+"/index/v1.lvi", "--storage-uri", fsBlobPathPrefix+"/storage", "--version-local-store-index-path", fsBlobPathPrefix+"/index/v1.lsi")
	executeCommandLine("upsync", "--source-path", testPath+"/version/v2", "--target-path", fsBlobPathPrefix+"/index/v2.lvi", "--storage-uri", fsBlobPathPrefix+"/storage", "--version-local-store-index-path", fsBlobPathPrefix+"/index/v2.lsi")
	executeCommandLine("upsync", "--source-path", testPath+"/version/v3", "--target-path", fsBlobPathPrefix+"/index/v3.lvi", "--storage-uri", fsBlobPathPrefix+"/storage", "--version-local-store-index-path", fsBlobPathPrefix+"/index/v3.lsi")

	cmd, err := executeCommandLine("downsync", "--source-path", fsBlobPathPrefix+"/index/v1.lvi", "--target-path", testPath+"/version/current", "--storage-uri", fsBlobPathPrefix+"/storage", "--version-local-store-index-path", fsBlobPathPrefix+"/index/v1.lsi")
	assert.NoError(t, err, cmd)
	validateContent(t, fsBlobPathPrefix, "version/current", v1FilesCreate)
	cmd, err = executeCommandLine("downsync", "--source-path", fsBlobPathPrefix+"/index/v2.lvi", "--target-path", testPath+"/version/current", "--storage-uri", fsBlobPathPrefix+"/storage", "--version-local-store-index-path", fsBlobPathPrefix+"/index/v2.lsi")
	assert.NoError(t, err, cmd)
	validateContent(t, fsBlobPathPrefix, "version/current", v2FilesCreate)
	cmd, err = executeCommandLine("downsync", "--source-path", fsBlobPathPrefix+"/index/v3.lvi", "--target-path", testPath+"/version/current", "--storage-uri", fsBlobPathPrefix+"/storage", "--version-local-store-index-path", fsBlobPathPrefix+"/index/v3.lsi")
	assert.NoError(t, err, cmd)
	validateContent(t, fsBlobPathPrefix, "version/current", v3FilesCreate)
}

func TestDownsyncWithCache(t *testing.T) {
	testPath, _ := os.MkdirTemp("", "test")
	fsBlobPathPrefix := "fsblob://" + testPath
	createVersionData(t, fsBlobPathPrefix)
	executeCommandLine("upsync", "--source-path", testPath+"/version/v1", "--target-path", fsBlobPathPrefix+"/index/v1.lvi", "--storage-uri", fsBlobPathPrefix+"/storage")
	executeCommandLine("upsync", "--source-path", testPath+"/version/v2", "--target-path", fsBlobPathPrefix+"/index/v2.lvi", "--storage-uri", fsBlobPathPrefix+"/storage")
	executeCommandLine("upsync", "--source-path", testPath+"/version/v3", "--target-path", fsBlobPathPrefix+"/index/v3.lvi", "--storage-uri", fsBlobPathPrefix+"/storage")

	cmd, err := executeCommandLine("downsync", "--source-path", fsBlobPathPrefix+"/index/v1.lvi", "--target-path", testPath+"/version/current", "--storage-uri", fsBlobPathPrefix+"/storage", "--cache-path", testPath+"/cache")
	assert.NoError(t, err, cmd)
	validateContent(t, fsBlobPathPrefix, "version/current", v1FilesCreate)
	cmd, err = executeCommandLine("downsync", "--source-path", fsBlobPathPrefix+"/index/v2.lvi", "--target-path", testPath+"/version/current", "--storage-uri", fsBlobPathPrefix+"/storage", "--cache-path", testPath+"/cache")
	assert.NoError(t, err, cmd)
	validateContent(t, fsBlobPathPrefix, "version/current", v2FilesCreate)
	cmd, err = executeCommandLine("downsync", "--source-path", fsBlobPathPrefix+"/index/v3.lvi", "--target-path", testPath+"/version/current", "--storage-uri", fsBlobPathPrefix+"/storage", "--cache-path", testPath+"/cache")
	assert.NoError(t, err, cmd)
	validateContent(t, fsBlobPathPrefix, "version/current", v3FilesCreate)
}

func TestDownsyncWithLSIAndCache(t *testing.T) {
	testPath, _ := os.MkdirTemp("", "test")
	fsBlobPathPrefix := "fsblob://" + testPath
	createVersionData(t, fsBlobPathPrefix)
	executeCommandLine("upsync", "--source-path", testPath+"/version/v1", "--target-path", fsBlobPathPrefix+"/index/v1.lvi", "--storage-uri", fsBlobPathPrefix+"/storage", "--version-local-store-index-path", fsBlobPathPrefix+"/index/v1.lsi")
	executeCommandLine("upsync", "--source-path", testPath+"/version/v2", "--target-path", fsBlobPathPrefix+"/index/v2.lvi", "--storage-uri", fsBlobPathPrefix+"/storage", "--version-local-store-index-path", fsBlobPathPrefix+"/index/v2.lsi")
	executeCommandLine("upsync", "--source-path", testPath+"/version/v3", "--target-path", fsBlobPathPrefix+"/index/v3.lvi", "--storage-uri", fsBlobPathPrefix+"/storage", "--version-local-store-index-path", fsBlobPathPrefix+"/index/v3.lsi")

	cmd, err := executeCommandLine("downsync", "--source-path", fsBlobPathPrefix+"/index/v1.lvi", "--target-path", testPath+"/version/current", "--storage-uri", fsBlobPathPrefix+"/storage", "--version-local-store-index-path", fsBlobPathPrefix+"/index/v1.lsi", "--cache-path", testPath+"/cache")
	assert.NoError(t, err, cmd)
	validateContent(t, fsBlobPathPrefix, "version/current", v1FilesCreate)
	cmd, err = executeCommandLine("downsync", "--source-path", fsBlobPathPrefix+"/index/v2.lvi", "--target-path", testPath+"/version/current", "--storage-uri", fsBlobPathPrefix+"/storage", "--version-local-store-index-path", fsBlobPathPrefix+"/index/v2.lsi", "--cache-path", testPath+"/cache")
	assert.NoError(t, err, cmd)
	validateContent(t, fsBlobPathPrefix, "version/current", v2FilesCreate)
	cmd, err = executeCommandLine("downsync", "--source-path", fsBlobPathPrefix+"/index/v3.lvi", "--target-path", testPath+"/version/current", "--storage-uri", fsBlobPathPrefix+"/storage", "--version-local-store-index-path", fsBlobPathPrefix+"/index/v3.lsi", "--cache-path", testPath+"/cache")
	assert.NoError(t, err, cmd)
	validateContent(t, fsBlobPathPrefix, "version/current", v3FilesCreate)
}

func TestDownsyncWithValidate(t *testing.T) {
	testPath, _ := os.MkdirTemp("", "test")
	fsBlobPathPrefix := "fsblob://" + testPath
	createVersionData(t, fsBlobPathPrefix)
	executeCommandLine("upsync", "--source-path", testPath+"/version/v1", "--target-path", fsBlobPathPrefix+"/index/v1.lvi", "--storage-uri", fsBlobPathPrefix+"/storage")
	executeCommandLine("upsync", "--source-path", testPath+"/version/v2", "--target-path", fsBlobPathPrefix+"/index/v2.lvi", "--storage-uri", fsBlobPathPrefix+"/storage")
	executeCommandLine("upsync", "--source-path", testPath+"/version/v3", "--target-path", fsBlobPathPrefix+"/index/v3.lvi", "--storage-uri", fsBlobPathPrefix+"/storage")

	cmd, err := executeCommandLine("downsync", "--source-path", fsBlobPathPrefix+"/index/v1.lvi", "--target-path", testPath+"/version/current", "--storage-uri", fsBlobPathPrefix+"/storage", "--validate")
	assert.NoError(t, err, cmd)
	validateContent(t, fsBlobPathPrefix, "version/current", v1FilesCreate)
	cmd, err = executeCommandLine("downsync", "--source-path", fsBlobPathPrefix+"/index/v2.lvi", "--target-path", testPath+"/version/current", "--storage-uri", fsBlobPathPrefix+"/storage", "--validate")
	assert.NoError(t, err, cmd)
	validateContent(t, fsBlobPathPrefix, "version/current", v2FilesCreate)
	cmd, err = executeCommandLine("downsync", "--source-path", fsBlobPathPrefix+"/index/v3.lvi", "--target-path", testPath+"/version/current", "--storage-uri", fsBlobPathPrefix+"/storage", "--validate")
	assert.NoError(t, err, cmd)
	validateContent(t, fsBlobPathPrefix, "version/current", v3FilesCreate)
}

func TestDownsyncWithVersionLSIWithValidate(t *testing.T) {
	testPath, _ := os.MkdirTemp("", "test")
	fsBlobPathPrefix := "fsblob://" + testPath
	createVersionData(t, fsBlobPathPrefix)

	executeCommandLine("upsync", "--source-path", testPath+"/version/v1", "--target-path", fsBlobPathPrefix+"/index/v1.lvi", "--storage-uri", fsBlobPathPrefix+"/storage", "--version-local-store-index-path", fsBlobPathPrefix+"/index/v1.lsi")
	executeCommandLine("upsync", "--source-path", testPath+"/version/v2", "--target-path", fsBlobPathPrefix+"/index/v2.lvi", "--storage-uri", fsBlobPathPrefix+"/storage", "--version-local-store-index-path", fsBlobPathPrefix+"/index/v2.lsi")
	executeCommandLine("upsync", "--source-path", testPath+"/version/v3", "--target-path", fsBlobPathPrefix+"/index/v3.lvi", "--storage-uri", fsBlobPathPrefix+"/storage", "--version-local-store-index-path", fsBlobPathPrefix+"/index/v3.lsi")

	cmd, err := executeCommandLine("downsync", "--source-path", fsBlobPathPrefix+"/index/v1.lvi", "--target-path", testPath+"/version/current", "--storage-uri", fsBlobPathPrefix+"/storage", "--version-local-store-index-path", fsBlobPathPrefix+"/index/v1.lsi", "--validate")
	assert.NoError(t, err, cmd)
	validateContent(t, fsBlobPathPrefix, "version/current", v1FilesCreate)
	cmd, err = executeCommandLine("downsync", "--source-path", fsBlobPathPrefix+"/index/v2.lvi", "--target-path", testPath+"/version/current", "--storage-uri", fsBlobPathPrefix+"/storage", "--version-local-store-index-path", fsBlobPathPrefix+"/index/v2.lsi", "--validate")
	assert.NoError(t, err, cmd)
	validateContent(t, fsBlobPathPrefix, "version/current", v2FilesCreate)
	cmd, err = executeCommandLine("downsync", "--source-path", fsBlobPathPrefix+"/index/v3.lvi", "--target-path", testPath+"/version/current", "--storage-uri", fsBlobPathPrefix+"/storage", "--version-local-store-index-path", fsBlobPathPrefix+"/index/v3.lsi", "--validate")
	assert.NoError(t, err, cmd)
	validateContent(t, fsBlobPathPrefix, "version/current", v3FilesCreate)
}

func TestDownsyncWithCacheWithValidate(t *testing.T) {
	testPath, _ := os.MkdirTemp("", "test")
	fsBlobPathPrefix := "fsblob://" + testPath
	createVersionData(t, fsBlobPathPrefix)
	executeCommandLine("upsync", "--source-path", testPath+"/version/v1", "--target-path", fsBlobPathPrefix+"/index/v1.lvi", "--storage-uri", fsBlobPathPrefix+"/storage")
	executeCommandLine("upsync", "--source-path", testPath+"/version/v2", "--target-path", fsBlobPathPrefix+"/index/v2.lvi", "--storage-uri", fsBlobPathPrefix+"/storage")
	executeCommandLine("upsync", "--source-path", testPath+"/version/v3", "--target-path", fsBlobPathPrefix+"/index/v3.lvi", "--storage-uri", fsBlobPathPrefix+"/storage")

	cmd, err := executeCommandLine("downsync", "--source-path", fsBlobPathPrefix+"/index/v1.lvi", "--target-path", testPath+"/version/current", "--storage-uri", fsBlobPathPrefix+"/storage", "--cache-path", testPath+"/cache", "--validate")
	assert.NoError(t, err, cmd)
	validateContent(t, fsBlobPathPrefix, "version/current", v1FilesCreate)
	cmd, err = executeCommandLine("downsync", "--source-path", fsBlobPathPrefix+"/index/v2.lvi", "--target-path", testPath+"/version/current", "--storage-uri", fsBlobPathPrefix+"/storage", "--cache-path", testPath+"/cache", "--validate")
	assert.NoError(t, err, cmd)
	validateContent(t, fsBlobPathPrefix, "version/current", v2FilesCreate)
	cmd, err = executeCommandLine("downsync", "--source-path", fsBlobPathPrefix+"/index/v3.lvi", "--target-path", testPath+"/version/current", "--storage-uri", fsBlobPathPrefix+"/storage", "--cache-path", testPath+"/cache", "--validate")
	assert.NoError(t, err, cmd)
	validateContent(t, fsBlobPathPrefix, "version/current", v3FilesCreate)
}

func TestDownsyncWithLSIAndCacheWithValidate(t *testing.T) {
	testPath, _ := os.MkdirTemp("", "test")
	fsBlobPathPrefix := "fsblob://" + testPath
	createVersionData(t, fsBlobPathPrefix)
	executeCommandLine("upsync", "--source-path", testPath+"/version/v1", "--target-path", fsBlobPathPrefix+"/index/v1.lvi", "--storage-uri", fsBlobPathPrefix+"/storage", "--version-local-store-index-path", fsBlobPathPrefix+"/index/v1.lsi")
	executeCommandLine("upsync", "--source-path", testPath+"/version/v2", "--target-path", fsBlobPathPrefix+"/index/v2.lvi", "--storage-uri", fsBlobPathPrefix+"/storage", "--version-local-store-index-path", fsBlobPathPrefix+"/index/v2.lsi")
	executeCommandLine("upsync", "--source-path", testPath+"/version/v3", "--target-path", fsBlobPathPrefix+"/index/v3.lvi", "--storage-uri", fsBlobPathPrefix+"/storage", "--version-local-store-index-path", fsBlobPathPrefix+"/index/v3.lsi")

	cmd, err := executeCommandLine("downsync", "--source-path", fsBlobPathPrefix+"/index/v1.lvi", "--target-path", testPath+"/version/current", "--storage-uri", fsBlobPathPrefix+"/storage", "--version-local-store-index-path", fsBlobPathPrefix+"/index/v1.lsi", "--cache-path", testPath+"/cache", "--validate")
	assert.NoError(t, err, cmd)
	validateContent(t, fsBlobPathPrefix, "version/current", v1FilesCreate)
	cmd, err = executeCommandLine("downsync", "--source-path", fsBlobPathPrefix+"/index/v2.lvi", "--target-path", testPath+"/version/current", "--storage-uri", fsBlobPathPrefix+"/storage", "--version-local-store-index-path", fsBlobPathPrefix+"/index/v2.lsi", "--cache-path", testPath+"/cache", "--validate")
	assert.NoError(t, err, cmd)
	validateContent(t, fsBlobPathPrefix, "version/current", v2FilesCreate)
	cmd, err = executeCommandLine("downsync", "--source-path", fsBlobPathPrefix+"/index/v3.lvi", "--target-path", testPath+"/version/current", "--storage-uri", fsBlobPathPrefix+"/storage", "--version-local-store-index-path", fsBlobPathPrefix+"/index/v3.lsi", "--cache-path", testPath+"/cache", "--validate")
	assert.NoError(t, err, cmd)
	validateContent(t, fsBlobPathPrefix, "version/current", v3FilesCreate)
}

func TestDownsyncMissingChunks(t *testing.T) {
	testPath, _ := os.MkdirTemp("", "test")
	fsBlobPathPrefix := "fsblob://" + testPath
	createVersionData(t, fsBlobPathPrefix)
	executeCommandLine("upsync", "--source-path", testPath+"/version/v1", "--target-path", fsBlobPathPrefix+"/index/v1.lvi", "--storage-uri", fsBlobPathPrefix+"/storage")
	executeCommandLine("upsync", "--source-path", testPath+"/version/v2", "--target-path", fsBlobPathPrefix+"/index/v2.lvi", "--storage-uri", fsBlobPathPrefix+"/storage")
	executeCommandLine("upsync", "--source-path", testPath+"/version/v3", "--target-path", fsBlobPathPrefix+"/index/v3.lvi", "--storage-uri", fsBlobPathPrefix+"/storage")

	os.RemoveAll(path.Join(testPath, "storage/chunks"))

	cmd, err := executeCommandLine("downsync", "--source-path", fsBlobPathPrefix+"/index/v1.lvi", "--target-path", testPath+"/version/current", "--storage-uri", fsBlobPathPrefix+"/storage")
	assert.Error(t, err, cmd)
	cmd, err = executeCommandLine("downsync", "--source-path", fsBlobPathPrefix+"/index/v2.lvi", "--target-path", testPath+"/version/current", "--storage-uri", fsBlobPathPrefix+"/storage")
	assert.Error(t, err, cmd)
	cmd, err = executeCommandLine("downsync", "--source-path", fsBlobPathPrefix+"/index/v3.lvi", "--target-path", testPath+"/version/current", "--storage-uri", fsBlobPathPrefix+"/storage")
	assert.Error(t, err, cmd)
}

func TestDownsyncMissingIndex(t *testing.T) {
	testPath, _ := os.MkdirTemp("", "test")
	fsBlobPathPrefix := "fsblob://" + testPath
	createVersionData(t, fsBlobPathPrefix)
	executeCommandLine("upsync", "--source-path", testPath+"/version/v1", "--target-path", fsBlobPathPrefix+"/index/v1.lvi", "--storage-uri", fsBlobPathPrefix+"/storage")
	executeCommandLine("upsync", "--source-path", testPath+"/version/v2", "--target-path", fsBlobPathPrefix+"/index/v2.lvi", "--storage-uri", fsBlobPathPrefix+"/storage")
	executeCommandLine("upsync", "--source-path", testPath+"/version/v3", "--target-path", fsBlobPathPrefix+"/index/v3.lvi", "--storage-uri", fsBlobPathPrefix+"/storage")

	longtailutils.DeleteByURI(fsBlobPathPrefix + "/storage/store.lsi")

	cmd, err := executeCommandLine("downsync", "--source-path", fsBlobPathPrefix+"/index/v1.lvi", "--target-path", testPath+"/version/current", "--storage-uri", fsBlobPathPrefix+"/storage")
	assert.Error(t, err, cmd)
	cmd, err = executeCommandLine("downsync", "--source-path", fsBlobPathPrefix+"/index/v2.lvi", "--target-path", testPath+"/version/current", "--storage-uri", fsBlobPathPrefix+"/storage")
	assert.Error(t, err, cmd)
	cmd, err = executeCommandLine("downsync", "--source-path", fsBlobPathPrefix+"/index/v3.lvi", "--target-path", testPath+"/version/current", "--storage-uri", fsBlobPathPrefix+"/storage")
	assert.Error(t, err, cmd)

	cmd, err = executeCommandLine("init-remote-store", "--storage-uri", fsBlobPathPrefix+"/storage", "--worker-count", "1")
	assert.NoError(t, err, cmd)

	cmd, err = executeCommandLine("downsync", "--source-path", fsBlobPathPrefix+"/index/v1.lvi", "--target-path", testPath+"/version/current", "--storage-uri", fsBlobPathPrefix+"/storage", "--version-local-store-index-path", fsBlobPathPrefix+"/index/v1.lsi", "--cache-path", testPath+"/cache", "--validate")
	assert.NoError(t, err, cmd)
	validateContent(t, fsBlobPathPrefix, "version/current", v1FilesCreate)
	cmd, err = executeCommandLine("downsync", "--source-path", fsBlobPathPrefix+"/index/v2.lvi", "--target-path", testPath+"/version/current", "--storage-uri", fsBlobPathPrefix+"/storage", "--version-local-store-index-path", fsBlobPathPrefix+"/index/v2.lsi", "--cache-path", testPath+"/cache", "--validate")
	assert.NoError(t, err, cmd)
	validateContent(t, fsBlobPathPrefix, "version/current", v2FilesCreate)
	cmd, err = executeCommandLine("downsync", "--source-path", fsBlobPathPrefix+"/index/v3.lvi", "--target-path", testPath+"/version/current", "--storage-uri", fsBlobPathPrefix+"/storage", "--version-local-store-index-path", fsBlobPathPrefix+"/index/v3.lsi", "--cache-path", testPath+"/cache", "--validate")
	assert.NoError(t, err, cmd)
	validateContent(t, fsBlobPathPrefix, "version/current", v3FilesCreate)
}

func TestMultiVersionDownsync(t *testing.T) {
	testPath, _ := os.MkdirTemp("", "test")
	fsBlobPathPrefix := "fsblob://" + testPath
	createLayeredData(t, fsBlobPathPrefix)
	executeCommandLine("upsync", "--exclude-filter-regex", ".*layer2$**.*layer3$", "--source-path", testPath+"/source", "--target-path", fsBlobPathPrefix+"/index/base.lvi", "--storage-uri", fsBlobPathPrefix+"/storage", "--version-local-store-index-path", fsBlobPathPrefix+"/index/base.lsi")
	executeCommandLine("upsync", "--include-filter-regex", ".*/$**.*\\.layer2$", "--source-path", testPath+"/source", "--target-path", fsBlobPathPrefix+"/index/layer2.lvi", "--storage-uri", fsBlobPathPrefix+"/storage", "--version-local-store-index-path", fsBlobPathPrefix+"/index/layer2.lsi")
	executeCommandLine("upsync", "--include-filter-regex", ".*/$**.*\\.layer3$", "--source-path", testPath+"/source", "--target-path", fsBlobPathPrefix+"/index/layer3.lvi", "--storage-uri", fsBlobPathPrefix+"/storage", "--version-local-store-index-path", fsBlobPathPrefix+"/index/layer3.lsi")

	cmd, err := executeCommandLine("downsync", "--source-paths", fsBlobPathPrefix+"/index/base.lvi|"+fsBlobPathPrefix+"/index/layer2.lvi|"+fsBlobPathPrefix+"/index/layer3.lvi", "--target-path", testPath+"/target", "--storage-uri", fsBlobPathPrefix+"/storage", "--version-local-store-index-path", fsBlobPathPrefix+"/index/base.lsi|"+fsBlobPathPrefix+"/index/layer2.lsi|"+fsBlobPathPrefix+"/index/layer3.lsi")
	assert.NoError(t, err, cmd)
	validateContent(t, fsBlobPathPrefix, "target", layerData)
}
