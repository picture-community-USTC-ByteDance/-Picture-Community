package service

import (
	"errors"
	"mime/multipart"
	"net/http"
	"path"
	"picture_community/dao/post"
	"picture_community/entity/db"
	"picture_community/global"
	"picture_community/response"
	"picture_community/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	StorageLocation = "storage"
	ServerName      = "121.5.1.73"
)

func FileUpload(c *gin.Context, id uint, file *multipart.FileHeader) response.ResStruct {
	utils.PathExists(StorageLocation)

	filesuffix := path.Ext(file.Filename)
	file.Filename = strconv.FormatUint(uint64(id), 10) + strconv.FormatInt(time.Now().Unix(), 10) + utils.RandStr(20) + filesuffix
	dst := path.Join(StorageLocation, file.Filename)
	err := c.SaveUploadedFile(file, dst)
	if err != nil {
		return response.ResStruct{
			HttpStatus: http.StatusGatewayTimeout,
			Code:       response.FailCode,
			Message:    err.Error(),
			Data:       nil,
		}
	}
	dst = "http://" + ServerName + ":8080/upload/pictures/" + file.Filename
	return response.ResStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Message:    "ok",
		Data:       gin.H{"url": dst},
	}
}

func CreatePost(c *gin.Context, id uint, url string, content string) response.ResStruct {
	newPost := db.Post{
		UID:              id,
		TitlePhotoUrl:    url,
		Content:          content,
		PhotoNumber:      1,
		CommentNumber:    0,
		LikeNumber:       0,
		ForwardNumber:    0,
		CollectionNumber: 0,
	}
	postID, err := post.InsertPostByUserID(newPost)
	if err != nil {
		return response.ResStruct{
			HttpStatus: http.StatusGatewayTimeout,
			Code:       response.FailCode,
			Message:    err.Error(),
			Data:       nil,
		}
	}
	return response.ResStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Message:    "ok",
		Data:       gin.H{"post_id": postID},
	}
}

func DeletePost(uid uint, pid uint) response.ResStruct {
	post := db.Post{
		PID: pid,
	}

	err := global.MysqlDB.Where("uid = ?", uid).Delete(&post).Error
	if err != nil {
		return response.ResStruct{
			HttpStatus: http.StatusBadRequest,
			Code:       response.FailCode,
			Message:    err.Error(),
		}
	}
	return response.ResStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Message:    "delete success",
		Data:       nil,
	}
}

func NewForward(uid uint, pid uint, content string) response.ResStruct {
	var forward db.Forward
	var post db.Post
	err := global.MysqlDB.Where("p_id = ?", pid).First(&post).Error
	if err != nil {
		return response.ResStruct{
			HttpStatus: http.StatusOK,
			Code:       response.FailCode,
			Message:    "Target post not exist",
			Data:       err.Error(),
		}
	}

	err = global.MysqlDB.Where("author_user_id = ? AND to_forward_post_id = ?", uid, pid).First(&forward).Error
	if err == nil {
		return response.ResStruct{
			HttpStatus: http.StatusOK,
			Code:       response.FailCode,
			Message:    "Already forwarded",
			Data: gin.H{
				"forward_id": forward.FID,
			},
		}
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		forward = db.Forward{
			AuthorUserID:    uid,
			ToForwardPostID: pid,
			CommentNumber:   0,
			LikeNumber:      0,
			Content:         content,
		}
		err = global.MysqlDB.Create(&forward).Error
	}

	if err != nil {
		return response.ResStruct{
			HttpStatus: http.StatusForbidden,
			Code:       response.FailCode,
			Message:    err.Error(),
			Data:       nil,
		}
	}
	return response.ResStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Message:    "Forward success",
		Data: gin.H{
			"forward_id": forward.FID,
		},
	}
}

func DeleteForward(uid uint, fid uint) response.ResStruct {
	forward := db.Forward{
		FID: fid,
	}

	err := global.MysqlDB.Where("author_user_id = ?", uid).Delete(&forward).Error
	if err != nil {
		return response.ResStruct{
			HttpStatus: http.StatusBadRequest,
			Code:       response.FailCode,
			Message:    err.Error(),
		}
	}
	return response.ResStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Message:    "delete success",
		Data:       nil,
	}
}
