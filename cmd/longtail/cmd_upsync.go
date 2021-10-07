package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/DanEngelbrecht/golongtail/longtaillib"
	"github.com/DanEngelbrecht/golongtail/longtailstorelib"
	"github.com/DanEngelbrecht/golongtail/longtailutils"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

func upsync(
	numWorkerCount int,
	blobStoreURI string,
	sourceFolderPath string,
	sourceIndexPath string,
	targetFilePath string,
	targetChunkSize uint32,
	targetBlockSize uint32,
	maxChunksPerBlock uint32,
	compressionAlgorithm string,
	hashAlgorithm string,
	includeFilterRegEx string,
	excludeFilterRegEx string,
	minBlockUsagePercent uint32,
	versionLocalStoreIndexPath string,
	getConfigPath string) ([]longtailutils.StoreStat, []longtailutils.TimeStat, error) {

	storeStats := []longtailutils.StoreStat{}
	timeStats := []longtailutils.TimeStat{}

	setupStartTime := time.Now()
	pathFilter, err := longtailutils.MakeRegexPathFilter(includeFilterRegEx, excludeFilterRegEx)
	if err != nil {
		return storeStats, timeStats, err
	}

	fs := longtaillib.CreateFSStorageAPI()
	defer fs.Dispose()

	sourceFolderScanner := asyncFolderScanner{}
	if sourceIndexPath == "" {
		sourceFolderScanner.scan(sourceFolderPath, pathFilter, fs)
	}

	jobs := longtaillib.CreateBikeshedJobAPI(uint32(numWorkerCount), 0)
	defer jobs.Dispose()
	hashRegistry := longtaillib.CreateFullHashRegistry()
	defer hashRegistry.Dispose()

	compressionType, err := getCompressionType(compressionAlgorithm)
	if err != nil {
		return storeStats, timeStats, err
	}
	hashIdentifier, err := getHashIdentifier(hashAlgorithm)
	if err != nil {
		return storeStats, timeStats, err
	}

	setupTime := time.Since(setupStartTime)
	timeStats = append(timeStats, longtailutils.TimeStat{"Setup", setupTime})

	sourceIndexReader := asyncVersionIndexReader{}
	sourceIndexReader.read(sourceFolderPath,
		sourceIndexPath,
		targetChunkSize,
		compressionType,
		hashIdentifier,
		pathFilter,
		fs,
		jobs,
		hashRegistry,
		&sourceFolderScanner)

	remoteStore, err := createBlockStoreForURI(blobStoreURI, "", jobs, numWorkerCount, targetBlockSize, maxChunksPerBlock, longtailstorelib.ReadWrite)
	if err != nil {
		return storeStats, timeStats, err
	}
	defer remoteStore.Dispose()

	creg := longtaillib.CreateFullCompressionRegistry()
	defer creg.Dispose()

	indexStore := longtaillib.CreateCompressBlockStore(remoteStore, creg)
	defer indexStore.Dispose()

	vindex, hash, readSourceIndexTime, err := sourceIndexReader.get()
	if err != nil {
		return storeStats, timeStats, err
	}
	defer vindex.Dispose()
	timeStats = append(timeStats, longtailutils.TimeStat{"Read source index", readSourceIndexTime})

	getMissingContentStartTime := time.Now()
	existingRemoteStoreIndex, errno := longtailutils.GetExistingStoreIndexSync(indexStore, vindex.GetChunkHashes(), minBlockUsagePercent)
	if errno != 0 {
		return storeStats, timeStats, errors.Wrapf(longtaillib.ErrnoToError(errno, longtaillib.ErrEIO), "upsync: longtailutils.GetExistingStoreIndexSync(%s) failed", blobStoreURI)
	}
	defer existingRemoteStoreIndex.Dispose()

	versionMissingStoreIndex, errno := longtaillib.CreateMissingContent(
		hash,
		existingRemoteStoreIndex,
		vindex,
		targetBlockSize,
		maxChunksPerBlock)
	if errno != 0 {
		return storeStats, timeStats, errors.Wrapf(longtaillib.ErrnoToError(errno, longtaillib.ErrEIO), "upsync: longtaillib.CreateMissingContent(%s) failed", sourceFolderPath)
	}
	defer versionMissingStoreIndex.Dispose()

	getMissingContentTime := time.Since(getMissingContentStartTime)
	timeStats = append(timeStats, longtailutils.TimeStat{"Get content index", getMissingContentTime})

	writeContentStartTime := time.Now()
	if versionMissingStoreIndex.GetBlockCount() > 0 {
		writeContentProgress := longtailutils.CreateProgress("Writing content blocks")
		defer writeContentProgress.Dispose()

		errno = longtaillib.WriteContent(
			fs,
			indexStore,
			jobs,
			&writeContentProgress,
			versionMissingStoreIndex,
			vindex,
			normalizePath(sourceFolderPath))
		if errno != 0 {
			return storeStats, timeStats, errors.Wrapf(longtaillib.ErrnoToError(errno, longtaillib.ErrEIO), "upsync: longtaillib.WriteContent(%s) failed", sourceFolderPath)
		}
	}
	writeContentTime := time.Since(writeContentStartTime)
	timeStats = append(timeStats, longtailutils.TimeStat{"Write version content", writeContentTime})

	flushStartTime := time.Now()

	stores := []longtaillib.Longtail_BlockStoreAPI{
		indexStore,
		remoteStore,
	}
	errno = longtailutils.FlushStoresSync(stores)
	if errno != 0 {
		return storeStats, timeStats, errors.Wrapf(longtaillib.ErrnoToError(errno, longtaillib.ErrEIO), "longtailutils.FlushStoresSync: Failed for `%v`", stores)
	}

	flushTime := time.Since(flushStartTime)
	timeStats = append(timeStats, longtailutils.TimeStat{"Flush", flushTime})

	indexStoreStats, errno := indexStore.GetStats()
	if errno == 0 {
		storeStats = append(storeStats, longtailutils.StoreStat{"Compress", indexStoreStats})
	}
	remoteStoreStats, errno := remoteStore.GetStats()
	if errno == 0 {
		storeStats = append(storeStats, longtailutils.StoreStat{"Remote", remoteStoreStats})
	}

	writeVersionIndexStartTime := time.Now()
	vbuffer, errno := longtaillib.WriteVersionIndexToBuffer(vindex)
	if errno != 0 {
		return storeStats, timeStats, errors.Wrapf(longtaillib.ErrnoToError(errno, longtaillib.ErrEIO), "upsync: longtaillib.WriteVersionIndexToBuffer() failed")
	}

	err = longtailstorelib.WriteToURI(targetFilePath, vbuffer)
	if err != nil {
		return storeStats, timeStats, errors.Wrapf(err, "upsync: longtaillib.longtailstorelib.WriteToURL() failed")
	}
	writeVersionIndexTime := time.Since(writeVersionIndexStartTime)
	timeStats = append(timeStats, longtailutils.TimeStat{"Write version index", writeVersionIndexTime})

	if versionLocalStoreIndexPath != "" {
		writeVersionLocalStoreIndexStartTime := time.Now()
		versionLocalStoreIndex, errno := longtaillib.MergeStoreIndex(existingRemoteStoreIndex, versionMissingStoreIndex)
		if errno != 0 {
			return storeStats, timeStats, errors.Wrapf(longtaillib.ErrnoToError(errno, longtaillib.ErrENOMEM), "upsync: longtaillib.MergeStoreIndex() failed")
		}
		defer versionLocalStoreIndex.Dispose()
		versionLocalStoreIndexBuffer, errno := longtaillib.WriteStoreIndexToBuffer(versionLocalStoreIndex)
		if errno != 0 {
			return storeStats, timeStats, errors.Wrapf(longtaillib.ErrnoToError(errno, longtaillib.ErrENOMEM), "upsync: longtaillib.WriteStoreIndexToBuffer() failed")
		}
		err = longtailstorelib.WriteToURI(versionLocalStoreIndexPath, versionLocalStoreIndexBuffer)
		if err != nil {
			return storeStats, timeStats, errors.Wrapf(err, "upsync: longtailstorelib.WriteToURL() failed")
		}
		writeVersionLocalStoreIndexTime := time.Since(writeVersionLocalStoreIndexStartTime)
		timeStats = append(timeStats, longtailutils.TimeStat{"Write version store index", writeVersionLocalStoreIndexTime})
	}

	if getConfigPath != "" {
		writeGetConfigStartTime := time.Now()

		v := viper.New()
		v.SetConfigType("json")
		v.Set("storage-uri", blobStoreURI)
		v.Set("source-path", targetFilePath)
		if versionLocalStoreIndexPath != "" {
			v.Set("version-local-store-index-path", versionLocalStoreIndexPath)
		}
		tmpFile, err := ioutil.TempFile(os.TempDir(), "longtail-")
		if err != nil {
			return storeStats, timeStats, errors.Wrapf(err, "upsync: ioutil.TempFile() failed")
		}
		tmpFilePath := tmpFile.Name()
		tmpFile.Close()
		fmt.Printf("tmp file: %s", tmpFilePath)
		err = v.WriteConfigAs(tmpFilePath)
		if err != nil {
			return storeStats, timeStats, errors.Wrapf(err, "upsync: v.WriteConfigAs() failed")
		}

		bytes, err := ioutil.ReadFile(tmpFilePath)
		if err != nil {
			return storeStats, timeStats, errors.Wrapf(err, "upsync: ioutil.ReadFile(%s) failed", tmpFilePath)
		}
		os.Remove(tmpFilePath)

		err = longtailstorelib.WriteToURI(getConfigPath, bytes)
		if err != nil {
			return storeStats, timeStats, errors.Wrapf(err, "upsync: longtailstorelib.WriteToURI(%s) failed", getConfigPath)
		}

		writeGetConfigTime := time.Since(writeGetConfigStartTime)
		timeStats = append(timeStats, longtailutils.TimeStat{"Write get config", writeGetConfigTime})
	}

	return storeStats, timeStats, nil
}

type UpsyncCmd struct {
	SourcePath                 string `name:"source-path" help:"Source folder path" required:""`
	SourceIndexPath            string `name:"source-index-path" help:"Optional pre-computed index of source-path"`
	TargetPath                 string `name:"target-path" help:"Target file uri" required:""`
	VersionLocalStoreIndexPath string `name:"version-local-store-index-path" help:"Target file uri for a store index optimized for this particular version"`
	GetConfigPath              string `name:"get-config-path" help:"Target file uri for json formatted get-config file"`
	TargetChunkSizeOption
	MaxChunksPerBlockOption
	TargetBlockSizeOption
	MinBlockUsagePercentOption
	StorageURIOption
	CompressionOption
	HashingOption
	UpsyncIncludeRegExOption
	UpsyncExcludeRegExOption
}

func (r *UpsyncCmd) Run(ctx *Context) error {
	storeStats, timeStats, err := upsync(
		ctx.NumWorkerCount,
		r.StorageURI,
		r.SourcePath,
		r.SourceIndexPath,
		r.TargetPath,
		r.TargetChunkSize,
		r.TargetBlockSize,
		r.MaxChunksPerBlock,
		r.Compression,
		r.Hashing,
		r.IncludeFilterRegEx,
		r.ExcludeFilterRegEx,
		r.MinBlockUsagePercent,
		r.VersionLocalStoreIndexPath,
		r.GetConfigPath)
	ctx.StoreStats = append(ctx.StoreStats, storeStats...)
	ctx.TimeStats = append(ctx.TimeStats, timeStats...)
	return err
}
