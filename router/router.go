package router

import (
	"picture_community/controller"
	"picture_community/global"
)

func SetRouter() {
	r := global.GinEngine

	r.GET("/search", controller.Search)
	g := r.Group("/user")
	{
		g.POST("/login", controller.LoginController)
	}
	p := r.Group("/post")
	{
		p.POST("/create", controller.CreatePostController)
	}
}
