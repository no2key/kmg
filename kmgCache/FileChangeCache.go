package kmgCache

import (
	"github.com/bronze1man/kmg/encoding/kmgGob"
	"github.com/bronze1man/kmg/kmgConfig"
	"github.com/bronze1man/kmg/kmgCrypto"
	"github.com/bronze1man/kmg/kmgFile"
	"os"
	"path/filepath"
)

func getFileChangeCachePath(key string) string {
	return filepath.Join(kmgConfig.DefaultEnv().TmpPath, "FileChangeCache", key)
}

func MustMd5FileChangeCache(key string, pathList []string, f func()) {
	//读取文件修改时间缓存信息
	toChange := false
	cacheInfo := map[string]string{}
	cacheFilePath := getFileChangeCachePath(key)
	err := kmgGob.ReadFile(cacheFilePath, &cacheInfo)
	if err != nil {
		//忽略缓存读取的任何错误
		cacheInfo = map[string]string{}
	}
	for _, path := range pathList {
		statList, err := kmgFile.GetAllFileAndDirectoryStat(path)
		if err != nil {
			if os.IsNotExist(err) {
				toChange = true
				//fmt.Printf("[MustFileChangeCache] path:[%s] not exist\n", path)
				break
			}
			panic(err)
		}
		for _, stat := range statList {
			if stat.Fi.IsDir() {
				continue
			}

			cacheMd5 := cacheInfo[stat.FullPath]
			if kmgCrypto.MustMd5File(stat.FullPath) != cacheMd5 {
				toChange = true
				//fmt.Printf("[MustMd5FileChangeCache] path:[%s] mod md5 not match save[%s] file[%s]\n", stat.FullPath,
				//	cacheMd5, kmgCrypto.MustMd5File(stat.FullPath))
				break
			}
		}
		if toChange {
			break
		}
	}
	if !toChange {
		return
	}
	f()
	cacheInfo = map[string]string{}
	for _, path := range pathList {
		statList, err := kmgFile.GetAllFileAndDirectoryStat(path)
		if err != nil {
			panic(err)
		}
		for _, stat := range statList {
			if stat.Fi.IsDir() {
				continue
			}
			cacheInfo[stat.FullPath] = kmgCrypto.MustMd5File(stat.FullPath)
		}
	}
	kmgFile.MustMkdirForFile(cacheFilePath)
	kmgGob.MustWriteFile(cacheFilePath, cacheInfo)
	//保存文件缓存信息
	return
}
