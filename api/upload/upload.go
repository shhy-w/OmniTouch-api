package upload

// func transformer(u *model.Upload) any {
// 	url, err := minio.PresignedGetObject(u.Path)
// 	if err != nil {
// 		logger.Error(err)
// 	}
// 	u.Path = url
// 	if u.Thumbnail != "" {
// 		url, err := minio.PresignedGetObject(u.Thumbnail)
// 		if err != nil {
// 			logger.Error(err)
// 		}
// 		u.Thumbnail = url
// 	}
// 	return u
// }

// func Get(c *gin.Context) {
// 	core := cosy.Core[model.Upload](c)

// 	core.SetTransformer(transformer).Get()
// }

// func GetList(c *gin.Context) {
// 	core := cosy.Core[model.Upload](c).
// 		SetPreloads("User").
// 		SetFussy("name", "mime").
// 		SetEqual("to")

// 	lastId := c.Query("last_id")
// 	if lastId != "" && lastId != "0" {
// 		core.GormScope(func(tx *gorm.DB) *gorm.DB {
// 			return tx.Where("id < ?", cast.ToInt(lastId))
// 		})
// 	}

// 	core.SetTransformer(transformer).PagingList()
// }

// func Modify(c *gin.Context) {
// 	core := cosy.Core[model.Upload](c).
// 		SetValidRules(gin.H{
// 			"name": "required",
// 		})

// 	core.Modify()
// }

// func DestroyMedia(c *gin.Context) {
// 	core := cosy.Core[model.Upload](c)

// 	core.ExecutedHook(func(ctx *cosy.Ctx[model.Upload]) {
// 		err := minio.Remove(ctx.OriginModel.Path)
// 		if err != nil {
// 			logger.Error(err)
// 		}
// 		err = minio.Remove(ctx.OriginModel.Thumbnail)
// 		if err != nil {
// 			logger.Error(err)
// 		}
// 	}).PermanentlyDelete()
// }

// func Upload(c *gin.Context) {
// 	u, err := upload.Put(c)
// 	if err != nil {
// 		cosy.ErrHandler(c, err)
// 		return
// 	}

// 	url, err := minio.PresignedGetObject(u.Path)
// 	if err != nil {
// 		cosy.ErrHandler(c, err)
// 		return
// 	}

// 	u.Path = url

// 	if u.Thumbnail != "" {
// 		url, err = minio.PresignedGetObject(u.Thumbnail)
// 		if err == nil {
// 			u.Thumbnail = url
// 		}
// 	}

// 	c.JSON(http.StatusOK, u)
// }
