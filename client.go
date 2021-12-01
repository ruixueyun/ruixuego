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

	"github.com/valyala/fasthttp"

	"git.jiaxianghudong.com/ruixuesdk/ruixuego/httpclient"

	"github.com/google/uuid"
)

const (
	headerTraceID   = "ruixue-traceid"
	headerCPID      = "ruixue-cpid"
	headerTimestamp = "ruixue-cpts"
	headerSign      = "ruixue-cpsign"
)

const (
	apiAddRelation           = "/Social/ServerAPI/AddRelation"
	apiDelRelation           = "/Social/ServerAPI/DeleteRelation"
	apiUpdateRelationRemarks = "/Social/ServerAPI/UpdateRelationRemarks"
	apiAddFriend             = "/Social/ServerAPI/AddFriend"
	apiDelFriend             = "/Social/ServerAPI/DelFriend"
	apiUpdateFriendRemarks   = "/Social/ServerAPI/UpdateFriendRemarks"
	apiCustom                = "/Social/ServerAPI/SetCustom"
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

// SetCustom 给用户设置社交自定义信息
func (c *Client) SetCustom(appID, openID, custom string) error {
	ret := &response{}
	err := c.query(apiCustom, &argCustom{
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
