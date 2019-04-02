package follow

//var getBaseLocker = new(sync.Mutex)
//
//// GetBase 获取 base 脚本
//func GetBase() (err error) {
//	getBaseLocker.Lock()
//	defer getBaseLocker.Unlock()
//	baseListURL := fmt.Sprintf(viper.GetString("API.FollowScriptBaseList"),
//		viper.GetString("BaseURL"), viper.GetString("MacAddress"))
//
//	var scripts []tables.FollowBaseScript
//	if scripts, err = new(tables.FollowBaseScript).GetAll(); err != nil {
//		return
//	}
//
//	resp, body, errs := gorequest.New().
//		Retry(3, 5*time.Second, http.StatusBadRequest, http.StatusInternalServerError).
//		Timeout(time.Second * 10).Get(baseListURL).End()
//	if len(errs) != 0 {
//		return errs[len(errs)-1]
//	}
//	if resp.StatusCode != 200 {
//		return fmt.Errorf("response code: %d", resp.StatusCode)
//	}
//	var newScripts map[string]string
//	if err = json.Unmarshal([]byte(body), &newScripts); err != nil {
//		return
//	}
//	if len(scripts) != 0 && len(newScripts) == len(scripts) {
//		var catchScript bool // 检测当前的版本是否一致
//		for scriptID, version := range newScripts {
//			var temp bool
//			for _, script := range scripts {
//				if script.ScriptID == scriptID && version == script.Version {
//					temp = true
//				}
//			}
//			if !temp {
//				catchScript = true
//				break
//			}
//		}
//		if !catchScript {
//			return
//		}
//	}
//
//	var hasError bool
//	var manifests = map[string][]byte{}
//	for scriptID := range newScripts {
//		var manifest []byte
//		if manifest, err = getBaseSpecial(fmt.Sprintf("%s/api/devices/%s/book-followreadding-base/%s",
//			viper.GetString("BaseURL"), viper.GetString("MacAddress"), scriptID)); err != nil {
//			hasError = true
//			logstash.Error(err.Error())
//			break
//		} else {
//			var manifestData GlobalSetting
//			if err = json.Unmarshal(manifest, &manifestData); err != nil {
//				return
//			}
//			if err = downloadUncompressPkg(manifestData.BasePkg.File); err != nil {
//				return
//			}
//			manifests[scriptID] = manifest
//		}
//	}
//	if hasError { // 有一个请求错误，直接返回
//		return fmt.Errorf("something error occurred during request follow-base-script")
//	}
//	// store it to database
//	tx := store.Major.Begin()
//	if err = tx.Delete(&tables.FollowBaseScript{}).Error; err != nil {
//		tx.Rollback()
//		return
//	}
//	for scriptID, manifest := range manifests {
//		if err = tx.Create(&tables.FollowBaseScript{
//			ScriptID: scriptID,
//			Manifest: manifest,
//			Version:  newScripts[scriptID],
//		}).Error; err != nil {
//			tx.Rollback()
//			return
//		}
//	}
//	if err = tx.Commit().Error; err != nil {
//		tx.Rollback()
//	}
//	return
//}
//
//func getBaseSpecial(url string) (manifest []byte, err error) {
//	resp, body, errs := gorequest.New().Timeout(time.Second * 10).Get(url).End()
//
//	if len(errs) != 0 {
//		return manifest, errs[len(errs)-1]
//	}
//	if resp.StatusCode != 200 {
//		return manifest, fmt.Errorf("response code: %d", resp.StatusCode)
//	}
//	if resp.Header.Get("Result-Code") != "0" {
//		return manifest, fmt.Errorf("response Result-Code: %s", resp.Header.Get("Result-Code"))
//	}
//	manifest = []byte(body)
//	return
//}
//
//const serviceName = "followBase"
//
//// downloadUncompressPkg 同步的下载解压 pkg 包
//func downloadUncompressPkg(address string) (err error) {
//	var downloadDestination = path.Join(viper.GetString("TaskDir"), "download", util.UUID()+filepath.Ext(address))
//	if err = alioss.Download(address, downloadDestination); err != nil {
//		return
//	}
//	var uncompressDestination = path.Join(viper.GetString("TaskDir"), "uncompress", util.UUID())
//	if err = archiver.TarGz.Unarchive(downloadDestination, uncompressDestination); err != nil {
//		return
//	}
//	if !com.IsFile(path.Join(uncompressDestination, "files.json")) {
//		return fmt.Errorf("cannot find files.json: %s", path.Join(uncompressDestination, "files.json"))
//	}
//	var getCon []byte
//	if getCon, err = ioutil.ReadFile(path.Join(uncompressDestination, "files.json")); err != nil {
//		return
//	}
//	tx := store.Major.Begin()
//	var files map[string]string
//	if err = json.Unmarshal(getCon, &files); err != nil {
//		logstash.Error(err.Error())
//	}
//	for hash, name := range files {
//		filename := path.Join(uncompressDestination, name)
//
//		// 文件复制之前验证文件 hash
//		var verify bool
//		if verify, err = util.MD5Verify(filename, hash); err != nil {
//			return
//		} else if !verify {
//			continue
//		}
//
//		// 复制文件
//		f := path.Join(viper.GetString("TipsDir"), util.UUID()+filepath.Ext(name))
//		if err = com.Copy(path.Join(uncompressDestination, name), f); err != nil {
//			tx.Rollback()
//			return
//		}
//
//		// 文件复制之后验证文件 hash
//		if verify, err = util.MD5Verify(f, hash); err != nil {
//			return
//		} else if !verify {
//			continue
//		}
//
//		var tips = tables.Tip{}
//		if err = tx.Find(&tips, &tables.Tip{Hash: hash}).Error; err != nil && err == gorm.ErrRecordNotFound {
//			if err = tx.Create(&tables.Tip{Hash: hash, Path: f}).Error; err != nil {
//				tx.Rollback()
//				return
//			}
//		} else if err != nil {
//			tx.Rollback()
//			return
//		} else {
//			logstash.Info(fmt.Sprintf("Audio hash %s is already exist.", hash))
//		}
//	}
//	if err = tx.Commit().Error; err != nil {
//		tx.Rollback()
//		return
//	}
//	logstash.Info("Success download base pkg and compress.")
//	return
//}
