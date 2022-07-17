package znet

import "zinx/zinx/ziface"

// -------------------------------------------
// @file          : route.go
// @author        : binshow
// @time          : 2022/7/17 3:23 PM
// @description   : 基础路由类，实现自定义路由时需要基础这个 BaseRouter
// -------------------------------------------


//实现router时，先嵌入这个基类，然后根据需要对这个基类的方法进行重写
type BaseRouter struct {}

//这里之所以BaseRouter的方法都为空，
// 是因为有的Router不希望有PreHandle或PostHandle
// 所以Router全部继承BaseRouter的好处是，不需要实现PreHandle和PostHandle也可以实例化

func (br *BaseRouter)PreHandle(req ziface.IRequest){}

func (br *BaseRouter)Handle(req ziface.IRequest){}

func (br *BaseRouter)PostHandle(req ziface.IRequest){}