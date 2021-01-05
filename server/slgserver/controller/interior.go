package controller

import (
	"github.com/goinggo/mapstructure"
	"slgserver/constant"
	"slgserver/middleware"
	"slgserver/net"
	"slgserver/server/slgserver/logic/mgr"
	"slgserver/server/slgserver/model"
	"slgserver/server/slgserver/proto"
)

var DefaultInterior = Interior{}

type Interior struct {

}

func (this*Interior) InitRouter(r *net.Router) {
	g := r.Group("interior").Use(middleware.ElapsedTime(),
		middleware.Log(), middleware.CheckRole())
	g.AddRouter("collection", this.collection)
	g.AddRouter("transform", this.transform)
}

func (this*Interior) collection(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	reqObj := &proto.CollectionReq{}
	rspObj := &proto.CollectionRsp{}

	mapstructure.Decode(req.Body.Msg, reqObj)
	rsp.Body.Msg = rspObj
	rsp.Body.Code = constant.OK

	r, _ := req.Conn.GetProperty("role")
	role := r.(*model.Role)
	roleRes, ok:= mgr.RResMgr.Get(role.RId)
	if ok == false {
		rsp.Body.Code = constant.DBError
		return

	}
	roleRes.Gold += mgr.GetYield(roleRes.RId).Gold
	roleRes.SyncExecute()

}


func (this*Interior) transform(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	reqObj := &proto.TransformReq{}
	rspObj := &proto.TransformRsp{}

	mapstructure.Decode(req.Body.Msg, reqObj)
	rsp.Body.Msg = rspObj
	rsp.Body.Code = constant.OK

	r, _ := req.Conn.GetProperty("role")
	role := r.(*model.Role)
	roleRes, ok:= mgr.RResMgr.Get(role.RId)
	if ok == false {
		rsp.Body.Code = constant.DBError
		return
	}

	//main, ok := logic.RCMgr.GetMainCity(role.RId)
	//add := logic.RFMgr.GetAdditions(main.CityId, facility.TypeTax)

	len := 4
	ret := make([]int, len)

	for i := 0 ;i < len; i++{
		//ret[i] = reqObj.To[i] - reqObj.From[i]
		if reqObj.From[i] > 0{
			ret[i] = -reqObj.From[i]
		}

		if reqObj.To[i] > 0{
			ret[i] = reqObj.To[i]
		}
	}


	if roleRes.Wood + ret[0] < 0{
		rsp.Body.Code = constant.InvalidParam
		return
	}

	if roleRes.Iron + ret[1] < 0{
		rsp.Body.Code = constant.InvalidParam
		return
	}

	if roleRes.Stone + ret[2] < 0{
		rsp.Body.Code = constant.InvalidParam
		return
	}

	if roleRes.Grain + ret[3] < 0{
		rsp.Body.Code = constant.InvalidParam
		return
	}

	roleRes.Wood += ret[0]
	roleRes.Iron += ret[1]
	roleRes.Stone += ret[2]
	roleRes.Grain += ret[3]
	roleRes.SyncExecute()

}

