// Copyright (c) 2021. Homeland Interactive Technology Ltd. All rights reserved.

package ruixuego

import (
	"crypto/sha1" // nolint:gosec
	"encoding/hex"
	"fmt"
	"hash"
	"strconv"
	"sync"
	"time"

	"git.jiaxianghudong.com/ruixuesdk/ruixuego/httpclient"

	"github.com/google/uuid"
	"github.com/valyala/fasthttp"
)

const (
	headerTraceID   = "ruixue-traceid"
	headerCPID      = "ruixue-cpid"
	headerTimestamp = "ruixue-cpts"
	headerSign      = "ruixue-cpsign"
	headerVersion   = "ruixue-version"
)

const (
	apiSetUserInfo           = "/Social/ServerAPI/SetUserInfo"
	apiSetCustom             = "/Social/ServerAPI/SetCustom"
	apiAddRelation           = "/Social/ServerAPI/AddRelation"
	apiDelRelation           = "/Social/ServerAPI/DeleteRelation"
	apiUpdateRelationRemarks = "/Social/ServerAPI/UpdateRelationRemarks"
	apiRelationList          = "/Social/ServerAPI/RelationList"
	apiHasRelation           = "/Social/ServerAPI/HasRelation"
	apiAddFriend             = "/Social/ServerAPI/AddFriend"
	apiDelFriend             = "/Social/ServerAPI/DelFriend"
	apiUpdateFriendRemarks   = "/Social/ServerAPI/UpdateFriendRemarks"
	apiFriendList            = "/Social/ServerAPI/FriendList"
	apiIsFriend              = "/Social/ServerAPI/IsFriend"
	apiLBSUpdate             = "/Social/ServerAPI/LBSUpdate"
	apiLBSDelete             = "/Social/ServerAPI/LBSDelete"
	apiLBSRadius             = "/Social/ServerAPI/LBSRadius"
)

var defaultClient *Client

func NewClient() *Client {
	c := &Client{}
	c.sha1Pool = &sync.Pool{
		New: func() interface{} {
			return sha1.New() // nolint:gosec
		},
	}
	return c
}

type Client struct {
	sha1Pool *sync.Pool
}

func (c *Client) getSign(traceID, ts string) string {
	h := c.sha1Pool.Get().(hash.Hash)
	_, _ = h.Write([]byte(traceID + ts + config.CPKey))
	ret := hex.EncodeToString(h.Sum(nil))
	h.Reset()
	c.sha1Pool.Put(h)
	return ret
}

func (c *Client) query(
	path string, arg, ret interface{}) error {

	traceID, cpID, ts := uuid.New().String(),
		strconv.FormatUint(uint64(config.CPID), 10),
		strconv.FormatInt(time.Now().Unix(), 10)

	req := httpclient.GetRequest()
	req.Header.Add(headerTraceID, traceID)
	req.Header.Add(headerCPID, cpID)
	req.Header.Add(headerTimestamp, ts)
	req.Header.Add(headerSign, c.getSign(traceID, ts))
	req.Header.Add(headerVersion, Version)

	if arg != nil {
		b, err := MarshalJSON(arg)
		if err != nil {
			return err
		}
		req.Header.SetMethod("POST")
		req.SetBody(b)
	}

	resp, err := httpclient.DefaultClient().DoRequest(config.APIDomain+path, req)
	if err != nil {
		return err
	}
	if resp.StatusCode() != fasthttp.StatusOK {
		return fmt.Errorf("%s", resp.Body())
	}

	if ret != nil {
		err = UnmarshalJSON(resp.Body(), ret)
		httpclient.PutResponse(resp)
		if err != nil {
			return err
		}
	} else {
		httpclient.PutResponse(resp)
	}

	return nil
}

func (c *Client) checkResponse(resp *response) error {
	if resp.Code != 0 {
		return fmt.Errorf("[%d] %s", resp.Code, resp.Msg)
	}
	return nil
}

// SetUserInfo 设置用户信息
func (c *Client) SetUserInfo(appID, openID string, userinfo *UserInfo) error {
	userinfo.AppID = appID
	userinfo.OpenID = openID

	ret := &response{}
	err := c.query(apiSetUserInfo, userinfo, ret)
	if err != nil {
		return err
	}
	err = c.checkResponse(ret)
	if err != nil {
		return err
	}

	return nil
}

// SetCustom 给用户设置社交模块的自定义信息
func (c *Client) SetCustom(appID, openID, custom string) error {
	ret := &response{}
	err := c.query(apiSetCustom, &argCustom{
		AppID:  appID,
		OpenID: openID,
		Custom: custom,
	}, ret)
	if err != nil {
		return err
	}
	err = c.checkResponse(ret)
	if err != nil {
		return err
	}

	return nil
}

// AddRelation 添加自定义关系
func (c *Client) AddRelation(
	types RelationTypes, openID, targetOpenID string, remark ...string) error {

	ret := &response{}
	arg := &argRelation{
		Types:  types,
		OpenID: openID,
		Target: targetOpenID,
	}
	if len(remark) > 0 {
		arg.TargetRemarks = remark[0]
	}
	if len(remark) > 1 {
		arg.TargetRemarks = remark[1]
	}
	err := c.query(apiAddRelation, arg, ret)
	if err != nil {
		return err
	}
	err = c.checkResponse(ret)
	if err != nil {
		return err
	}

	return nil
}

// DelRelation 删除自定义关系
func (c *Client) DelRelation(
	types RelationTypes, openID, targetOpenID string) error {

	ret := &response{}
	err := c.query(apiDelRelation, &argRelation{
		Types:  types,
		OpenID: openID,
		Target: targetOpenID,
	}, ret)
	if err != nil {
		return err
	}
	err = c.checkResponse(ret)
	if err != nil {
		return err
	}

	return nil
}

// UpdateRelationRemarks 更新自定关系备注
func (c *Client) UpdateRelationRemarks(
	typ, openID, targetOpenID, remarks string) error {

	ret := &response{}
	err := c.query(apiUpdateRelationRemarks, &argRelation{
		Type:          typ,
		OpenID:        openID,
		Target:        targetOpenID,
		TargetRemarks: remarks,
	}, ret)
	if err != nil {
		return err
	}
	err = c.checkResponse(ret)
	if err != nil {
		return err
	}

	return nil
}

// RelationList 获取自定关系列表
func (c *Client) RelationList(typ, openID string) ([]*RelationUser, error) {
	ret := make([]*RelationUser, 0)
	resp := &response{Data: &ret}
	err := c.query(apiRelationList, &argRelation{
		Type:   typ,
		OpenID: openID,
	}, resp)
	if err != nil {
		return nil, err
	}
	err = c.checkResponse(resp)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// HasRelation 判断 Target 是否与 User 存在指定关系
func (c *Client) HasRelation(typ, openID, targetOpenID string) (bool, error) {
	ret := false
	resp := &response{Data: &ret}
	err := c.query(apiHasRelation, &argRelation{
		Type:   typ,
		OpenID: openID,
		Target: targetOpenID,
	}, resp)
	if err != nil {
		return false, err
	}
	err = c.checkResponse(resp)
	if err != nil {
		return false, err
	}
	return ret, nil
}

// AddFriend 添加好友
func (c *Client) AddFriend(
	openID, targetOpenID string, remark ...string) error {

	ret := &response{}
	arg := &argRelation{
		OpenID: openID,
		Target: targetOpenID,
	}
	if len(remark) > 0 {
		arg.TargetRemarks = remark[0]
	}
	if len(remark) > 1 {
		arg.TargetRemarks = remark[1]
	}
	err := c.query(apiAddFriend, arg, ret)
	if err != nil {
		return err
	}
	err = c.checkResponse(ret)
	if err != nil {
		return err
	}

	return nil
}

// DelFriend 删除好友
func (c *Client) DelFriend(
	openID, targetOpenID string) error {

	ret := &response{}
	err := c.query(apiDelFriend, &argRelation{
		OpenID: openID,
		Target: targetOpenID,
	}, ret)
	if err != nil {
		return err
	}
	err = c.checkResponse(ret)
	if err != nil {
		return err
	}

	return nil
}

// UpdateFriendRemarks 更新好友备注
func (c *Client) UpdateFriendRemarks(
	openID, targetOpenID, remarks string) error {

	ret := &response{}
	err := c.query(apiUpdateFriendRemarks, &argRelation{
		OpenID:        openID,
		Target:        targetOpenID,
		TargetRemarks: remarks,
	}, ret)
	if err != nil {
		return err
	}
	err = c.checkResponse(ret)
	if err != nil {
		return err
	}

	return nil
}

// FriendList 获取好友列表
func (c *Client) FriendList(openID string) ([]*RelationUser, error) {
	ret := make([]*RelationUser, 0)
	resp := &response{Data: &ret}
	err := c.query(apiFriendList, &argRelation{
		OpenID: openID,
	}, resp)
	if err != nil {
		return nil, err
	}
	err = c.checkResponse(resp)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// IsFriend 判断 Target 是否为 User 的好友
func (c *Client) IsFriend(openID, targetOpenID string) (bool, error) {
	ret := false
	resp := &response{Data: &ret}
	err := c.query(apiIsFriend, &argRelation{
		OpenID: openID,
		Target: targetOpenID,
	}, resp)
	if err != nil {
		return false, err
	}
	err = c.checkResponse(resp)
	if err != nil {
		return false, err
	}
	return ret, nil
}

// LBSUpdate 更新 WGS84 坐标
// 		types 为 CP	自定义坐标分组, 比如可以同时将用户加入到 all 和 female 两个列表中
func (c *Client) LBSUpdate(openID string, types []string, lon, lat float64) error {
	ret := &response{}
	err := c.query(apiLBSUpdate, &argLocation{
		OpenID:    openID,
		Types:     types,
		Longitude: lon,
		Latitude:  lat,
	}, ret)
	if err != nil {
		return err
	}
	err = c.checkResponse(ret)
	if err != nil {
		return err
	}
	return nil
}

// LBSDelete 删除 WGS84 坐标
func (c *Client) LBSDelete(openID string, types []string) error {
	ret := &response{}
	err := c.query(apiLBSDelete, &argLocation{
		OpenID: openID,
		Types:  types,
	}, ret)
	if err != nil {
		return err
	}
	err = c.checkResponse(ret)
	if err != nil {
		return err
	}
	return nil
}

// LBSRadius 获取附近的人列表
func (c *Client) LBSRadius(
	openID, typ string,
	lon, lat, radius float64,
	page, pageSize int,
	count ...int) ([]*RelationUser, error) {

	ret := make([]*RelationUser, 0)
	resp := &response{Data: &ret}
	arg := &argLocation{
		OpenID:    openID,
		Type:      typ,
		Longitude: lon,
		Latitude:  lat,
		Radius:    radius,
		Page:      page,
		PageSize:  pageSize,
	}
	if len(count) == 1 {
		arg.Count = count[0]
	}
	err := c.query(apiLBSRadius, arg, resp)
	if err != nil {
		return nil, err
	}
	err = c.checkResponse(resp)
	if err != nil {
		return nil, err
	}
	return ret, nil
}
