// Copyright (c) 2021. Homeland Interactive Technology Ltd. All rights reserved.

package ruixuego

import (
	"crypto/sha1" // nolint:gosec
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/valyala/fasthttp"

	"git.jiaxianghudong.com/ruixuesdk/ruixuego/bufferpool"
)

const defaultStatus = -1

const (
	headerTraceID   = "ruixue-traceid"
	headerCPID      = "ruixue-cpid"
	headerTimestamp = "ruixue-cpts"
	headerSign      = "ruixue-cpsign"
	headerVersion   = "ruixue-version"
	headerDataCount = "ruixue-datacount"
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
	apiBigDataTrack          = "/Data/API/Track"
)

var defaultClient *Client

func NewClient() (c *Client, err error) {
	c = &Client{
		sha1Pool: &sync.Pool{
			New: func() interface{} {
				return sha1.New() // nolint:gosec
			},
		},
		httpClient: NewHTTPClient(config.Timeout, config.Concurrency),
	}

	if config.BigData != nil {
		c.producer, err = NewProducer(c, config.BigData)
		if err != nil {
			return nil, err
		}
	}
	return c, nil
}

type Client struct {
	httpClient *HTTPClient
	sha1Pool   *sync.Pool
	producer   *Producer
}

// Close SDK 客户端在关闭时必须显式调用该方法, 已保障数据不会丢失
func (c *Client) Close() error {
	if c.producer != nil {
		return c.producer.Close()
	}
	return nil
}

func (c *Client) getRequest(withoutSign ...bool) *fasthttp.Request {
	traceID, cpID, ts := uuid.New().String(),
		strconv.FormatUint(uint64(config.CPID), 10),
		strconv.FormatInt(time.Now().Unix(), 10)

	req := GetRequest()
	req.Header.Add("user-agent", "ruixue-go-sdk")
	req.Header.Add(headerVersion, Version)
	req.Header.Add(headerTraceID, traceID)
	req.Header.Add(headerCPID, cpID)
	req.Header.Add(headerTimestamp, ts)
	if len(withoutSign) == 0 {
		req.Header.Add(headerSign, c.getSign(traceID, ts))
	}

	return req
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
	path string, arg, ret interface{}, compress ...bool) error {
	_, err := c.queryCode(path, c.getRequest(), config.Timeout, arg, ret, compress...)
	return err
}

func (c *Client) queryCode(
	path string, req *fasthttp.Request, timeout time.Duration, arg, ret interface{}, compress ...bool) (int, error) {

	code := defaultStatus

	if arg != nil {
		req.Header.SetMethod("POST")

		var b []byte
		var ok bool
		var err error
		b, ok = arg.([]byte)
		if !ok {
			b, err = MarshalJSON(arg)
			if err != nil {
				return code, err
			}
		}
		if len(compress) == 1 && compress[0] {
			buf, err := gZIPCompress(b)
			if err != nil {
				bufferpool.Put(buf)
				return code, err
			}
			req.Header.Set("content-encoding", "gzip")
			req.SetBody(buf.Bytes())
			bufferpool.Put(buf)
		} else {
			req.SetBody(b)
		}
	}

	resp, err := c.httpClient.DoRequestWithTimeout(
		config.APIDomain+path, req, timeout)
	if err != nil {
		return code, err
	}
	code = resp.StatusCode()
	if code != fasthttp.StatusOK {
		return code, errors.New(http.StatusText(code))
	}

	if ret != nil {
		err = UnmarshalJSON(resp.Body(), ret)
		PutResponse(resp)
		if err != nil {
			return code, err
		}
	} else {
		PutResponse(resp)
	}

	return code, nil
}

func (c *Client) checkResponse(resp *response) error {
	if resp.Code != 0 {
		return fmt.Errorf("[%d] %s", resp.Code, resp.Msg)
	}
	return nil
}

// SetUserInfo 设置用户信息
func (c *Client) SetUserInfo(appID, openID string, userinfo *UserInfo) error {
	if openID == "" {
		return ErrInvalidOpenID
	}
	if appID == "" {
		return ErrInvalidAppID
	}

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
	if openID == "" {
		return ErrInvalidOpenID
	}
	if appID == "" {
		return ErrInvalidAppID
	}

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
// remarks[0] openID 用户给 targetOpenID 用户设置的备注
// remarks[1] targetOpenID 用户给 openID 用户设置的备注
func (c *Client) AddRelation(
	types RelationTypes, openID, targetOpenID string, remarks ...string) error {
	if openID == "" || targetOpenID == "" {
		return ErrInvalidOpenID
	}
	if len(types) == 0 {
		return ErrInvalidType
	}

	ret := &response{}
	arg := &argRelation{
		Types:  types,
		OpenID: openID,
		Target: targetOpenID,
	}
	if len(remarks) > 0 {
		arg.TargetRemarks = remarks[0]
	}
	if len(remarks) > 1 {
		arg.UserRemarks = remarks[1]
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
	if openID == "" || targetOpenID == "" {
		return ErrInvalidOpenID
	}
	if len(types) == 0 {
		return ErrInvalidType
	}

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
	if openID == "" || targetOpenID == "" {
		return ErrInvalidOpenID
	}
	if typ == "" {
		return ErrInvalidType
	}

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
	if openID == "" {
		return nil, ErrInvalidOpenID
	}
	if typ == "" {
		return nil, ErrInvalidType
	}

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
// remarks[0] openID 用户给 targetOpenID 用户设置的备注
// remarks[1] targetOpenID 用户给 openID 用户设置的备注
func (c *Client) AddFriend(
	openID, targetOpenID string, remarks ...string) error {
	if openID == "" || targetOpenID == "" {
		return ErrInvalidOpenID
	}

	ret := &response{}
	arg := &argRelation{
		OpenID: openID,
		Target: targetOpenID,
	}
	if len(remarks) > 0 {
		arg.TargetRemarks = remarks[0]
	}
	if len(remarks) > 1 {
		arg.UserRemarks = remarks[1]
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
	if openID == "" || targetOpenID == "" {
		return ErrInvalidOpenID
	}

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
	if openID == "" || targetOpenID == "" {
		return ErrInvalidOpenID
	}

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
	if openID == "" {
		return nil, ErrInvalidOpenID
	}

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
	if openID == "" || targetOpenID == "" {
		return false, ErrInvalidOpenID
	}

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
	if openID == "" {
		return ErrInvalidOpenID
	}
	if len(types) == 0 {
		return ErrInvalidType
	}

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
	if openID == "" {
		return ErrInvalidOpenID
	}
	if len(types) == 0 {
		return ErrInvalidType
	}

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

	if openID == "" {
		return nil, ErrInvalidOpenID
	}
	if typ == "" {
		return nil, ErrInvalidType
	}

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

// Track 大数据埋点事件上报
// 		distinctID 用户标识, 通常为瑞雪 OpenID
// 		event 事件名, 由 CP 自行指定, 后续应与大数据平台创建的埋点名一致
//		properties 自定义事件属性
// 		isLogined 用以标记 distinctID 是否为登录后的用户标识
func (c *Client) Track(
	devicecode, distinctID, event string, properties map[string]interface{}) error {
	return c.producer.Track(devicecode, distinctID, event, properties)
}

// track 大数据埋点记录
func (c *Client) track(data []byte, logCount int, compress bool) (int, error) {
	if len(data) == 0 {
		return defaultStatus, nil
	}

	req, ret := c.getRequest(true), &response{}
	req.Header.Add(headerDataCount, Itoa(logCount))
	code, err := c.queryCode(apiBigDataTrack, req, config.TrackTimeout, data, ret, compress)
	if err != nil {
		return code, err
	}
	err = c.checkResponse(ret)
	if err != nil {
		return code, err
	}
	return code, nil
}
