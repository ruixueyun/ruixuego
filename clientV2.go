// Copyright (c) 2021. Homeland Interactive Technology Ltd. All rights reserved.

package ruixuego

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/ruixueyun/ruixuego/bufferpool"
	"github.com/valyala/fasthttp"
	"log"
	"net/http"
	url2 "net/url"
	"strconv"
	"strings"
	"time"
)

var defaultClientV2 *ClientV2

func NewClientV2() (c *ClientV2, err error) {
	c = &ClientV2{
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

type ClientV2 struct {
	httpClient *HTTPClient
	producer   *Producer
}

func (c *ClientV2) GetProductID() string {
	return config.ProductID
}

func (c *ClientV2) GetChannelID() string {
	return config.ChannelID
}

func (c *ClientV2) GetRegion() string {
	return config.Region
}

func (c *ClientV2) GetLanguage() string {
	return config.Language
}

type ReqHeader struct {
	Header map[string]string `json:"header"`
}

func (h *ReqHeader) Set(key string, value string) {
	if h.Header == nil {
		h.Header = make(map[string]string)
	}
	h.Header[key] = value
}

// Close SDK 客户端在关闭时必须显式调用该方法, 已保障数据不会丢失
func (c *ClientV2) Close() error {
	if c.producer != nil {
		return c.producer.Close()
	}
	return nil
}

func (c *ClientV2) getRequest(header *ReqHeader, withoutSign ...bool) (string, *fasthttp.Request) {
	traceID, cpID, ts := uuid.New().String(),
		strconv.FormatUint(uint64(config.CPID), 10),
		strconv.FormatInt(time.Now().Unix(), 10)

	req := GetRequest()
	req.Header.Add("user-agent", "ruixue-go-sdk")
	req.Header.Add(headerVersion, Version)
	req.Header.Add(headerTraceID, traceID)
	req.Header.Add(headerCPID, cpID)
	req.Header.Add(headerProductID, c.GetProductID())
	req.Header.Add(headerChannelID, c.GetChannelID())
	req.Header.Add(headerServiceMark, config.ServiceMark)
	req.Header.Add(headerTimestamp, ts)
	req.Header.Add(headerNameRegion, c.GetRegion())
	if c.GetLanguage() != "" {
		req.Header.Add(headerLanguage, c.GetLanguage())
	}
	if len(withoutSign) == 0 {
		req.Header.Add(headerSign, GetSign(config.CPKey, traceID, ts))
	}

	if header != nil {
		for k, v := range header.Header {
			req.Header.Add(k, v)
		}
	}

	return traceID, req
}

func (c *ClientV2) getAndCheckResponse(url string, args map[string]string, header *ReqHeader, resp *response, compress ...bool) error {

	if resp == nil {
		resp = &response{}
	}

	dataValue := make(url2.Values)
	for k, v := range args {
		dataValue.Add(k, v)
	}

	uri := url + "?" + dataValue.Encode()

	traceID, err := c.query(uri, header, nil, resp, compress...)
	if err != nil {
		return errWithTraceID(err, traceID)
	}

	err = c.checkResponse(resp)
	if err != nil {
		return errWithTraceID(err, traceID)
	}

	return nil
}

func (c *ClientV2) queryAndCheckResponse(
	path string, header *ReqHeader, req interface{}, resp *response, compress ...bool) error {

	if resp == nil {
		resp = &response{}
	}

	traceID, err := c.query(path, header, req, resp, compress...)
	if err != nil {
		return errWithTraceID(err, traceID)
	}

	err = c.checkResponse(resp)
	if err != nil {
		return errWithTraceID(err, traceID)
	}

	return nil
}

func (c *ClientV2) query(
	path string, header *ReqHeader, arg, ret interface{}, compress ...bool) (string, error) {
	traceID, req := c.getRequest(header)
	_, err := c.queryCode(path, req, config.Timeout, arg, ret, compress...)
	return traceID, err
}

func (c *ClientV2) queryCode(
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

func (c *ClientV2) checkResponse(resp *response) error {
	if resp.Code != 0 {
		return fmt.Errorf("[%d] %s", resp.Code, resp.Msg)
	}
	return nil
}
func (c *ClientV2) queryAndCheckResponseWithProductIDAndChannelID(
	path string, header *ReqHeader, req interface{}, resp *response, productID, channelID string, compress ...bool) error {

	if resp == nil {
		resp = &response{}
	}

	traceID, err := c.queryWithProductIDAndChannelID(path, header, req, resp, productID, channelID, compress...)
	if err != nil {
		return errWithTraceID(err, traceID)
	}

	err = c.checkResponse(resp)
	if err != nil {
		return errWithTraceID(err, traceID)
	}

	return nil
}
func (c *ClientV2) queryWithProductIDAndChannelID(
	path string, header *ReqHeader, arg, ret interface{}, productID, channelID string, compress ...bool) (string, error) {
	traceID, req := c.getRequest(header)
	c.queryAddProductIDAndChannelID(req, productID, channelID)
	_, err := c.queryCode(path, req, config.Timeout, arg, ret, compress...)
	return traceID, err
}
func (c *ClientV2) queryAddProductIDAndChannelID(
	req *fasthttp.Request, productID, channelID string) {
	req.Header.Add(headerProductID, productID)
	req.Header.Add(headerChannelID, channelID)

}

type ReqSetCustom struct {
	ReqHeader
	ProductID string
	OpenID    string
	CpUserID  string
	Custom    string
}

// SetCustom 给用户设置社交模块的自定义信息
// OpenID 为 瑞雪opendid
// CpUserID cp侧用户id
func (c *ClientV2) SetCustom(req *ReqSetCustom) error {
	if req.OpenID == "" && req.CpUserID == "" {
		return ErrInvalidOpenID
	}
	if req.ProductID == "" {
		return ErrInvalidProductID
	}

	return c.queryAndCheckResponse(apiSetCustom, &req.ReqHeader, &argCustom{
		ProductID: req.ProductID,
		OpenID:    req.OpenID,
		CPUserID:  req.CpUserID,
		Custom:    req.Custom,
	}, nil)
}

type ReqAddRelation struct {
	ReqHeader
	Types        RelationTypes
	OpenID       string
	UserID       string
	TargetOpenID string
	TargetUserID string
	Remarks      []string
}

// AddRelation 添加自定义关系
// remarks[0] OpenID 用户给 targetOpenID 用户设置的备注
// remarks[1] targetOpenID 用户给 OpenID 用户设置的备注
func (c *ClientV2) AddRelation(req *ReqAddRelation) error {
	if (req.UserID == "" && req.OpenID == "") || (req.TargetUserID == "" && req.TargetOpenID == "") {
		return ErrInvalidOpenID
	}
	if len(req.Types) == 0 {
		return ErrInvalidType
	}

	arg := &argRelation{
		Types:          req.Types,
		OpenID:         req.OpenID,
		CPUserID:       req.UserID,
		Target:         req.TargetOpenID,
		TargetCPUserID: req.TargetUserID,
	}
	if len(req.Remarks) > 0 {
		arg.TargetRemarks = req.Remarks[0]
	}
	if len(req.Remarks) > 1 {
		arg.UserRemarks = req.Remarks[1]
	}

	return c.queryAndCheckResponse(apiAddRelation, &req.ReqHeader, arg, nil)
}

type ReqDelRelation struct {
	ReqHeader
	Types        RelationTypes
	OpenID       string
	UserID       string
	TargetOpenID string
	TargetUserID string
}

// DelRelation 删除自定义关系
func (c *ClientV2) DelRelation(req *ReqDelRelation) error {
	if (req.UserID == "" && req.OpenID == "") || (req.TargetUserID == "" && req.TargetOpenID == "") {
		return ErrInvalidOpenID
	}
	if len(req.Types) == 0 {
		return ErrInvalidType
	}

	return c.queryAndCheckResponse(apiDelRelation, &req.ReqHeader, &argRelation{
		Types:          req.Types,
		OpenID:         req.OpenID,
		CPUserID:       req.UserID,
		Target:         req.TargetOpenID,
		TargetCPUserID: req.TargetUserID,
	}, nil)
}

type ReqUpdateRelationRemarks struct {
	ReqHeader
	Type         string
	OpenID       string
	UserID       string
	TargetOpenID string
	TargetUserID string
	Remarks      string
}

// UpdateRelationRemarks 更新自定关系备注
func (c *ClientV2) UpdateRelationRemarks(req *ReqUpdateRelationRemarks) error {
	if (req.UserID == "" && req.OpenID == "") || (req.TargetUserID == "" && req.TargetOpenID == "") {
		return ErrInvalidOpenID
	}
	if len(req.Type) == 0 {
		return ErrInvalidType
	}

	return c.queryAndCheckResponse(apiUpdateRelationRemarks, &req.ReqHeader, &argRelation{
		Type:           req.Type,
		OpenID:         req.OpenID,
		CPUserID:       req.UserID,
		Target:         req.TargetOpenID,
		TargetCPUserID: req.TargetUserID,
	}, nil)
}

type ReqRelationList struct {
	ReqHeader
	Type   string
	OpenID string
	UserID string
}

// RelationList 获取自定关系列表
func (c *ClientV2) RelationList(req *ReqRelationList) ([]*RelationUser, error) {
	if req.UserID == "" && req.OpenID == "" {
		return nil, ErrInvalidOpenID
	}
	if len(req.Type) == 0 {
		return nil, ErrInvalidType
	}
	ret := make([]*RelationUser, 0)
	resp := &response{Data: &ret}

	err := c.queryAndCheckResponse(apiRelationList, &req.ReqHeader, &argRelation{
		Type:     req.Type,
		OpenID:   req.OpenID,
		CPUserID: req.UserID,
	}, resp)

	if err != nil {
		return nil, err
	}

	return ret, nil
}

type ReqHasRelation struct {
	ReqHeader
	Type         string
	OpenID       string
	UserID       string
	TargetOpenID string
	TargetUserID string
}

// HasRelation 判断 Target 是否与 User 存在指定关系
func (c *ClientV2) HasRelation(req *ReqHasRelation) (bool, error) {
	ret := false
	resp := &response{Data: &ret}

	err := c.queryAndCheckResponse(apiHasRelation, &req.ReqHeader, &argRelation{
		Type:           req.Type,
		OpenID:         req.OpenID,
		CPUserID:       req.UserID,
		Target:         req.TargetOpenID,
		TargetCPUserID: req.TargetUserID,
	}, resp)

	return ret, err
}

type ReqAddFriend struct {
	ReqHeader
	OpenID       string
	UserID       string
	TargetOpenID string
	TargetUserID string
	Remarks      []string
}

// AddFriend 添加好友
// remarks[0] OpenID 用户给 targetOpenID 用户设置的备注
// remarks[1] targetOpenID 用户给 OpenID 用户设置的备注
func (c *ClientV2) AddFriend(req *ReqAddFriend) error {
	if (req.UserID == "" && req.OpenID == "") || (req.TargetUserID == "" && req.TargetOpenID == "") {
		return ErrInvalidOpenID
	}

	arg := &argRelation{
		OpenID:         req.OpenID,
		CPUserID:       req.UserID,
		Target:         req.TargetOpenID,
		TargetCPUserID: req.TargetUserID,
	}
	if len(req.Remarks) > 0 {
		arg.TargetRemarks = req.Remarks[0]
	}
	if len(req.Remarks) > 1 {
		arg.UserRemarks = req.Remarks[1]
	}

	return c.queryAndCheckResponse(apiAddFriend, &req.ReqHeader, arg, nil)
}

type ReqDelFriend struct {
	ReqHeader
	OpenID       string
	UserID       string
	TargetOpenID string
	TargetUserID string
}

// DelFriend 删除好友
func (c *ClientV2) DelFriend(req *ReqDelFriend) error {
	if (req.UserID == "" && req.OpenID == "") || (req.TargetUserID == "" && req.TargetOpenID == "") {
		return ErrInvalidOpenID
	}
	return c.queryAndCheckResponse(apiDelFriend, &req.ReqHeader, &argRelation{
		OpenID:         req.OpenID,
		CPUserID:       req.UserID,
		Target:         req.TargetOpenID,
		TargetCPUserID: req.TargetUserID,
	}, nil)
}

type ReqUpdateFriendRemarks struct {
	ReqHeader
	OpenID       string
	UserID       string
	TargetOpenID string
	TargetUserID string
	Remarks      string
}

// UpdateFriendRemarks 更新好友备注
func (c *ClientV2) UpdateFriendRemarks(req *ReqUpdateFriendRemarks) error {
	if (req.UserID == "" && req.OpenID == "") || (req.TargetUserID == "" && req.TargetOpenID == "") {
		return ErrInvalidOpenID
	}
	return c.queryAndCheckResponse(apiUpdateFriendRemarks, &req.ReqHeader, &argRelation{
		OpenID:         req.OpenID,
		CPUserID:       req.UserID,
		Target:         req.TargetOpenID,
		TargetCPUserID: req.TargetUserID,
		TargetRemarks:  req.Remarks,
	}, nil)
}

type ReqFriendList struct {
	ReqHeader
	OpenID string
	UserID string
}

// FriendList 获取好友列表
func (c *ClientV2) FriendList(req *ReqFriendList) ([]*RelationUser, error) {
	if req.UserID == "" && req.OpenID == "" {
		return nil, ErrInvalidOpenID
	}

	ret := make([]*RelationUser, 0)
	resp := &response{Data: &ret}

	err := c.queryAndCheckResponse(apiFriendList, &req.ReqHeader, &argRelation{
		OpenID:   req.OpenID,
		CPUserID: req.UserID,
	}, resp)

	if err != nil {
		return nil, err
	}

	return ret, nil
}

type ReqGetRelationUser struct {
	ReqHeader
	Type         string
	OpenID       string
	UserID       string
	TargetOpenID string
	TargetUserID string
}

// GetRelationUser 查询好友信息
func (c *ClientV2) GetRelationUser(req *ReqGetRelationUser) (*RelationUser, error) {
	if (req.UserID == "" && req.OpenID == "") || (req.TargetUserID == "" && req.TargetOpenID == "") {
		return nil, ErrInvalidOpenID
	}

	if req.Type == "" {
		return nil, ErrInvalidType
	}

	ret := &RelationUser{}
	resp := &response{Data: ret}

	err := c.queryAndCheckResponse(apiGetRealtionUser, &req.ReqHeader, &argRelation{
		OpenID:         req.OpenID,
		CPUserID:       req.UserID,
		Target:         req.TargetOpenID,
		TargetCPUserID: req.TargetUserID,
		Type:           req.Type,
	}, resp)

	return ret, err
}

type ReqIsFriend struct {
	ReqHeader
	OpenID       string
	UserID       string
	TargetOpenID string
	TargetUserID string
}

// IsFriend 判断 Target 是否为 User 的好友
func (c *ClientV2) IsFriend(req *ReqIsFriend) (bool, error) {
	if (req.UserID == "" && req.OpenID == "") || (req.TargetUserID == "" && req.TargetOpenID == "") {
		return false, ErrInvalidOpenID
	}
	ret := false
	resp := &response{Data: &ret}

	err := c.queryAndCheckResponse(apiIsFriend, &req.ReqHeader, &argRelation{
		OpenID:         req.OpenID,
		CPUserID:       req.UserID,
		Target:         req.TargetOpenID,
		TargetCPUserID: req.TargetUserID,
	}, resp)

	return ret, err
}

type ReqLBSUpdate struct {
	ReqHeader
	OpenID string
	UserID string
	Types  []string
	Lon    float64
	Lat    float64
}

// LBSUpdate 更新 WGS84 坐标
//
//	types 为 CP	自定义坐标分组, 比如可以同时将用户加入到 all 和 female 两个列表中
func (c *ClientV2) LBSUpdate(req *ReqLBSUpdate) error {
	if req.UserID == "" && req.OpenID == "" {
		return ErrInvalidOpenID
	}

	if len(req.Types) == 0 {
		return ErrInvalidType
	}

	return c.queryAndCheckResponse(apiLBSUpdate, &req.ReqHeader, &argLocation{
		OpenID:    req.OpenID,
		CPUserID:  req.UserID,
		Types:     req.Types,
		Longitude: req.Lon,
		Latitude:  req.Lat,
	}, nil)
}

type ReqLBSDelete struct {
	ReqHeader
	OpenID string
	UserID string
	Types  []string
}

// LBSDelete 删除 WGS84 坐标
func (c *ClientV2) LBSDelete(req *ReqLBSDelete) error {
	if req.UserID == "" && req.OpenID == "" {
		return ErrInvalidOpenID
	}
	if len(req.Types) == 0 {
		return ErrInvalidType
	}

	return c.queryAndCheckResponse(apiLBSDelete, &req.ReqHeader, &argLocation{
		OpenID:   req.OpenID,
		CPUserID: req.UserID,
		Types:    req.Types,
	}, nil)
}

type ReqLBSRadius struct {
	ReqHeader
	OpenID   string
	UserID   string
	Typ      string
	Lon      float64
	Lat      float64
	Radius   float64
	Page     int
	PageSize int
	Count    []int
}

// LBSRadius 获取附近的人列表
func (c *ClientV2) LBSRadius(req *ReqLBSRadius) ([]*RelationUser, error) {

	if req.UserID == "" && req.OpenID == "" {
		return nil, ErrInvalidOpenID
	}
	if req.Typ == "" {
		return nil, ErrInvalidType
	}

	ret := make([]*RelationUser, 0)
	resp := &response{Data: &ret}
	arg := &argLocation{
		OpenID:    req.OpenID,
		CPUserID:  req.UserID,
		Type:      req.Typ,
		Longitude: req.Lon,
		Latitude:  req.Lat,
		Radius:    req.Radius,
		Page:      req.Page,
		PageSize:  req.PageSize,
	}
	if len(req.Count) == 1 {
		arg.Count = req.Count[0]
	}

	err := c.queryAndCheckResponse(apiLBSRadius, &req.ReqHeader, arg, resp)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

// Tracks 大数据埋点事件上报
//
//	devicecode (不能为空) 用户设备码. 用户使用设备的唯一识别码
//	distinctID (可为空) 用户标识. 通常为瑞雪 OpenID
//	opts: 动态参数设置
func (c *ClientV2) Tracks(
	devicecode, distinctID string, opts ...BigdataOptions) error {
	return c.producer.Tracks(devicecode, distinctID, opts...)
}

type ReqTrack struct {
	ReqHeader
	Data     []byte
	LogCount int
	Compress bool
}

// track 将埋点数据上报给瑞雪云
func (c *ClientV2) Track(track *ReqTrack) (int, error) {
	if len(track.Data) == 0 {
		return defaultStatus, nil
	}

	traceID, req := c.getRequest(&track.ReqHeader, true)
	ret := &response{}
	req.Header.Add(headerDataCount, Itoa(track.LogCount))
	code, err := c.queryCode(apiBigDataTrack, req, config.TrackTimeout, track.Data, ret, track.Compress)
	if err != nil {
		return code, errWithTraceID(err, traceID)
	}
	err = c.checkResponse(ret)
	if err != nil {
		return code, errWithTraceID(err, traceID)
	}
	return code, nil
}

type ReqSyncTrack struct {
	ReqHeader
	DeviceCode string
	DistinctID string
	Opts       []BigdataOptions
}

// SyncTrack 同步接口 直接将埋点数据上报给瑞雪云
// 前提要设置好 config
func (c *ClientV2) SyncTrack(track *ReqSyncTrack) error {
	if track.DeviceCode == "" && track.DistinctID == "" {
		return ErrInvalidDevicecode
	}

	logData := &BigDataLog{
		DistinctID: track.DistinctID,
		Devicecode: track.DeviceCode,
	}
	for _, opt := range track.Opts {
		err := opt(logData)
		if err != nil {
			return err
		}
	}
	if logData.Type == "" {
		return ErrInvalidType
	}
	if logData.CPID == 0 {
		if config.CPID == 0 {
			return ErrInvalidCPID
		}
		logData.CPID = config.CPID
	}
	if logData.PlatformID <= 0 {
		logData.PlatformID = 10
	}
	if logData.UUID == "" {
		logData.UUID = uuid.New().String()
	}
	if logData.Time == "" {
		logData.Time = time.Now().Format(time.RFC3339Nano)
	}

	data := []*BigDataLog{logData}

	b, err := MarshalJSON(data)
	if err != nil {
		return err
	}

	traceID, req := c.getRequest(&track.ReqHeader, true)
	ret := &response{}
	req.Header.Add(headerDataCount, Itoa(1))
	_, err = c.queryCode(apiBigDataTrack, req, config.TrackTimeout, b, ret, !config.BigData.DisableCompress)
	if err != nil {
		return errWithTraceID(err, traceID)
	}
	err = c.checkResponse(ret)
	if err != nil {
		return errWithTraceID(err, traceID)
	}
	return nil
}

type ReqCreateRank struct {
	ReqHeader
	RankID      string
	StartTime   time.Time
	DestroyTime time.Time
}

// CreateRank 创建排行榜
func (c *ClientV2) CreateRank(req *ReqCreateRank) error {
	if req.RankID == "" {
		return ErrInvalidOpenID
	}

	err := c.queryAndCheckResponse(apiCreateRank, &req.ReqHeader, &rankAPIArg{
		RankID:      req.RankID,
		StartTime:   req.StartTime.Format(time.RFC3339),
		DestroyTime: req.DestroyTime.Format(time.RFC3339),
	}, nil)

	return err
}

type ReqCloseRank struct {
	ReqHeader
	RankID string
}

// CloseRank 关闭排行榜
func (c *ClientV2) CloseRank(req *ReqCloseRank) error {
	if req.RankID == "" {
		return ErrInvalidOpenID
	}

	err := c.queryAndCheckResponse(apiCloseRank, &req.ReqHeader, &rankAPIArg{
		RankID: req.RankID,
	}, nil)

	return err
}

type ReqRankAddScore struct {
	ReqHeader
	RankID   string
	OpenId   string
	CpUserID string
	Score    int64
}

// RankAddScore 用户添加分数
func (c *ClientV2) RankAddScore(req *ReqRankAddScore) error {
	if req.RankID == "" || (req.OpenId == "" && req.CpUserID == "") {
		return ErrInvalidOpenID
	}

	err := c.queryAndCheckResponse(apiRankAddScore, &req.ReqHeader, &rankAPIArg{
		RankID:   req.RankID,
		OpenID:   req.OpenId,
		CPUserID: req.CpUserID,
		Score:    req.Score,
	}, nil)

	return err
}

type ReqRankSetScore struct {
	ReqHeader
	RankID   string
	OpenId   string
	CpUserID string
	Score    int64
}

// RankSetScore 用户设置分数
func (c *ClientV2) RankSetScore(req *ReqRankSetScore) error {
	if req.RankID == "" || (req.OpenId == "" && req.CpUserID == "") {
		return ErrInvalidOpenID
	}

	err := c.queryAndCheckResponse(apiRankSetScore, &req.ReqHeader, &rankAPIArg{
		RankID:   req.RankID,
		OpenID:   req.OpenId,
		CPUserID: req.CpUserID,
		Score:    req.Score,
	}, nil)

	return err
}

type ReqDeleteRankUser struct {
	ReqHeader
	RankID   string
	OpenId   string
	CpUserID string
}

// DeleteRankUser 删除排行榜用户
func (c *ClientV2) DeleteRankUser(req *ReqDeleteRankUser) error {
	if req.RankID == "" || (req.OpenId == "" && req.CpUserID == "") {
		return ErrInvalidOpenID
	}
	err := c.queryAndCheckResponse(apiRankDeleteUser, &req.ReqHeader, &rankAPIArg{
		RankID:   req.RankID,
		OpenID:   req.OpenId,
		CPUserID: req.CpUserID,
	}, nil)

	return err
}

type ReqQueryUserRank struct {
	ReqHeader
	RankID   string
	OpenId   string
	CpUserID string
}

// QueryUserRank 查询用户排行情况
func (c *ClientV2) QueryUserRank(req *ReqQueryUserRank) (*RankMember, error) {
	if req.RankID == "" || (req.OpenId == "" && req.CpUserID == "") {
		return nil, ErrInvalidOpenID
	}
	ret := &RankMember{}
	resp := &response{Data: ret}

	err := c.queryAndCheckResponse(apiQueryUserRank, &req.ReqHeader, &rankAPIArg{
		RankID:   req.RankID,
		OpenID:   req.OpenId,
		CPUserID: req.CpUserID,
	}, resp)

	return ret, err
}

type ReqGetRankList struct {
	ReqHeader
	RankID string
	Start  int32
	End    int32
}

// GetRankList 查询排行榜
func (c *ClientV2) GetRankList(req *ReqGetRankList) ([]*RankMember, error) {
	if req.RankID == "" {
		return nil, ErrInvalidOpenID
	}

	var ret []*RankMember
	resp := &response{Data: &ret}

	err := c.queryAndCheckResponse(apiGetRankList, &req.ReqHeader, &rankAPIArg{
		RankID:    req.RankID,
		StartRank: req.Start,
		EndRank:   req.End,
	}, resp)

	return ret, err
}

type ReqGetFriendRankList struct {
	ReqHeader
	RankID   string
	OpenId   string
	CpUserID string
}

// GetFriendRankList 查询好友排行榜
func (c *ClientV2) GetFriendRankList(req *ReqGetFriendRankList) ([]*RankMember, error) {
	if req.RankID == "" || (req.OpenId == "" && req.CpUserID == "") {
		return nil, ErrInvalidOpenID
	}

	var ret []*RankMember
	resp := &response{Data: &ret}

	err := c.queryAndCheckResponse(apiFriendsRank, &req.ReqHeader, &rankAPIArg{
		RankID:   req.RankID,
		OpenID:   req.OpenId,
		CPUserID: req.CpUserID,
	}, resp)

	return ret, err
}

type ReqGetRankDetail struct {
	ReqHeader
	RankID string
}

// GetRankDetail 查询排行榜详情
func (c *ClientV2) GetRankDetail(req *ReqGetRankDetail) (*RespRankDetail, error) {
	if req.RankID == "" {
		return nil, ErrInvalidOpenID
	}

	ret := &RespRankDetail{}
	resp := &response{Data: ret}

	err := c.queryAndCheckResponse(apiRankDetail, &req.ReqHeader, &rankAPIArg{
		RankID: req.RankID,
	}, resp)

	return ret, err
}

// GetAllRankIDList 查询所有排行ID
func (c *ClientV2) GetAllRankIDList(req *ReqHeader) (*RespAllRankID, error) {

	ret := &RespAllRankID{}
	resp := &response{Data: ret}

	err := c.queryAndCheckResponse(apiAllRankIDList, req, nil, resp)

	return ret, err
}

// IMSLogin ims 登陆接口
func (c *ClientV2) IMSLogin(req *IMSLoginReq) (*IMSLoginResp, error) {
	ret := &IMSLoginResp{}
	resp := &response{Data: ret}
	err := c.queryAndCheckResponse(apiIMSLogin, &req.ReqHeader, req, resp)
	return ret, err
}

// IMSSendMessage 发送消息
func (c *ClientV2) IMSSendMessage(req *IMSMessage) (*IMSMessageAck, error) {
	ret := &IMSMessageAck{}
	resp := &response{Data: ret}
	if req.MilliTS == 0 {
		req.MilliTS = time.Now().UnixMilli()
	}
	if req.UUID == "" {
		req.UUID = uuid.New().String()
	}
	convType, _, ok := IMSParseConversationID(req.ConversationID)
	if !ok {
		return nil, ErrInvalidIMSConversationID
	}
	if req.ConvType == 0 {
		req.ConvType = convType
	}
	req.CPID = config.CPID
	err := c.queryAndCheckResponse(apiIMSSendMessage, &req.ReqHeader, req, resp)
	return ret, err
}

// IMSGetHistory 获取历史记录
func (c *ClientV2) IMSGetHistory(req *IMSHistoryReq) (*IMSHistoryResp, error) {
	ret := &IMSHistoryResp{}
	resp := &response{Data: ret}
	err := c.queryAndCheckResponse(apiIMSGetHistory, &req.ReqHeader, req, resp)
	return ret, err
}

// IMSCreateConversation 创建会话
func (c *ClientV2) IMSCreateConversation(req *IMSCreateConvReq) error {
	return c.queryAndCheckResponse(apiIMSCreateConversation, &req.ReqHeader, req, nil)
}

// IMSUpdateConversation 更新会话信息
func (c *ClientV2) IMSUpdateConversation(req *IMSUpdateConvReq) error {
	return c.queryAndCheckResponse(apiIMSUpdateConversation, &req.ReqHeader, req, nil)
}

// IMSDeleteConversation 删除会话
func (c *ClientV2) IMSDeleteConversation(req *IMSConvDeleteReq) error {
	return c.queryAndCheckResponse(apiIMSDeleteConversation, &req.ReqHeader, req, nil)
}

// IMSGetConversation 获取会话信息
func (c *ClientV2) IMSGetConversation(req *IMSGetConversationReq) (*IMSConversation, error) {
	ret := &IMSConversation{}
	resp := &response{Data: ret}
	err := c.queryAndCheckResponse(apiIMSGetConversation, &req.ReqHeader, req, resp)
	return ret, err
}

// IMSJoinConversation 加入会话
func (c *ClientV2) IMSJoinConversation(req *IMSJoinConversationReq) error {
	return c.queryAndCheckResponse(apiIMSJoinConversation, &req.ReqHeader, req, nil)
}

// IMSLeaveConversation 离开会话
func (c *ClientV2) IMSLeaveConversation(req *IMSLeaveConversationReq) error {
	return c.queryAndCheckResponse(apiIMSLeaveConversation, &req.ReqHeader, req, nil)
}

// IMSUpdateConversationUserData 更新会话内用户信息
func (c *ClientV2) IMSUpdateConversationUserData(req *IMSUpdateConvUserDataReq) error {
	return c.queryAndCheckResponse(apiIMSUpdateConversationUserData, &req.ReqHeader, req, nil)
}

// IMSConversationUserList 获取会话中成员列表
func (c *ClientV2) IMSConversationUserList(req *IMSConversationUserListReq) ([]*IMSConversation, error) {
	ret := make([]*IMSConversation, 0)
	resp := &response{Data: &ret}
	err := c.queryAndCheckResponse(apiIMSConversationUserList, &req.ReqHeader, req, resp)
	return ret, err
}

// IMSChannelUsersCount 获取频道会话中玩家数量
func (c *ClientV2) IMSChannelUsersCount(count *ReqIMSChannelUsersCount) (map[string]int64, error) {
	ret := make(map[string]int64)
	resp := &response{Data: &ret}

	req := &IMSChannelUsesCountReq{
		ConversationIDs: count.ChannelConvIds,
	}

	err := c.queryAndCheckResponse(apiIMSChannelUsersCount, &count.ReqHeader, req, resp)
	return ret, err
}

// PusherPush 推送信息
func (c *ClientV2) PusherPush(req *ReqPusher) (*PusherPushRes, error) {
	ret := &PusherPushRes{}
	resp := &response{Data: ret}
	err := c.queryAndCheckResponseWithProductIDAndChannelID(apiPusherPush, &req.ReqHeader, req.Req, resp, req.ProductID, req.ChannelID)
	return ret, err
}

// RiskContentTextScan 内容安全文本检查（增强）
func (c *ClientV2) RiskContentTextScan(req *RiskContentTextScanReq) (*RiskContentTextScanResp, error) {
	if req == nil || len(req.Scene) <= 0 || len(req.Content) <= 0 {
		return nil, ErrInvalidParam
	}
	ret := &RiskContentTextScanResp{}
	resp := &response{Data: ret}
	err := c.queryAndCheckResponse(apiRiskTextScan, &req.ReqHeader, req, resp)
	return ret, err
}

// RiskContentImageScan 内容安全图片检查（增强）
func (c *ClientV2) RiskContentImageScan(req *RiskContentImageScanReq) (*RiskContentImageScanResp, error) {
	if len(req.URL) <= 0 {
		return nil, ErrInvalidParam
	}
	ret := &RiskContentImageScanResp{}
	resp := &response{Data: ret}
	err := c.queryAndCheckResponse(apiRiskImageScan, &req.ReqHeader, &RiskContentImageScanReq{
		URL: req.URL,
	}, resp)
	return ret, err
}

// ReportCustomAction 投放归因上报自定义action， 例如广点通小游戏创角， action传 CREATE_ROLE
// https://developers.e.qq.com/docs/guide/conversion/new_version/Mini_Game_api
func (c *ClientV2) ReportCustomAction(req *ReportCustomAction) error {
	if req.OpenID == "" || req.Action == "" {
		return ErrInvalidOpenID
	}
	arg := &ReportCustomAction{
		OpenID: req.OpenID,
		Action: req.Action,
	}
	resp := &response{}
	err := c.queryAndCheckResponse(apiReportCustomAction, &req.ReqHeader, arg, resp)
	if err != nil {
		return err
	}
	if resp.Code != 0 {
		return fmt.Errorf(resp.Msg)
	}
	return nil
}

// UpdateCPuserID 用于瑞雪侧更新cp侧的user_id
func (c *ClientV2) UpdateCPuserID(req *UpdateCPUserIDRequest) error {
	if len(req.CPUserID) == 0 || len(req.OpenID) == 0 {
		return ErrInvalidCPuserID
	}

	resp := &response{}
	err := c.queryAndCheckResponse(apiPassportUpdateCPUserID, &req.ReqHeader, req, resp)
	if err != nil {
		return err
	}
	if resp.Code != 0 {
		return fmt.Errorf(resp.Msg)
	}
	return nil
}

// RealAuth 实名检测
func (c *ClientV2) RealAuth(arg *RealAuthReq) (*RealAuthResponse, error) {
	if len(arg.ProductID) == 0 || len(arg.IDCard) == 0 || len(arg.RealName) == 0 {
		return nil, ErrInvalidParam
	}

	arg.CPID = config.CPID

	data := &RealAuthResponse{}
	resp := &response{
		Data: data,
	}
	err := c.queryAndCheckResponse(apiRiskRealAuthCheck, &arg.ReqHeader, arg, resp)
	if err != nil {
		return nil, err
	}
	log.Println(resp.Code, resp.Msg)
	if resp.Code != 0 {
		return nil, fmt.Errorf(resp.Msg)
	}
	return data, nil
}

// ExtensionExchange 兑换福利码
func (c *ClientV2) ExtensionExchange(arg *ExtensionExchangeReq) ([]*ExtensionProp, error) {
	if len(arg.CdKey) == 0 || len(arg.CpUserID) == 0 {
		return nil, ErrInvalidParam
	}
	data := []*ExtensionProp{}
	resp := &response{
		Data: data,
	}
	err := c.queryAndCheckResponse(apiOperationToolsExtensionExchange, &arg.ReqHeader, arg, resp)
	if err != nil {
		return nil, err
	}
	log.Println(resp.Code, resp.Msg)
	if resp.Code != 0 {
		return nil, fmt.Errorf(resp.Msg)
	}
	if resp.Data == nil {
		return nil, nil
	}
	return resp.Data.([]*ExtensionProp), nil
}

type ReqExtensionGameDisplay struct {
	ReqHeader
	GameID string `json:"game_id"`
}

// ExtensionGameDisplay 主播获取游戏内显示码
func (c *ClientV2) ExtensionGameDisplay(req *ReqExtensionGameDisplay) (*GameDisplayWelfareCodeInfoExp, error) {
	if req.GameID == "" {
		return nil, ErrInvalidParam
	}
	dataValue := make(url2.Values)
	dataValue.Add("game_id", req.GameID)
	url := apiOperationToolsExtensionGameDisplay
	uri := url + "?" + dataValue.Encode()
	data := &GameDisplayWelfareCodeInfoExp{}
	resp := &response{
		Data: data,
	}
	err := c.queryAndCheckResponse(uri, &req.ReqHeader, nil, resp)
	if err != nil {
		return nil, err
	}
	log.Println(resp.Code, resp.Msg)
	if resp.Code != 0 {
		return nil, fmt.Errorf(resp.Msg)
	}
	return data, nil
}

type ReqTradeOrderStatusByNo struct {
	ReqHeader
	OrderNo string `json:"order_no"`
}

// TradeOrderStatusByNo 平台订单状态
func (c *ClientV2) TradeOrderStatusByNo(req *ReqTradeOrderStatusByNo) (*OrderStatusRes, error) {
	if req.OrderNo == "" || !strings.HasSuffix(req.OrderNo, "v1") {
		return nil, ErrInvalidParam
	}
	dataValue := make(url2.Values)
	dataValue.Add("order_no", req.OrderNo)
	url := apiOrderInfoByNo
	uri := url + "?" + dataValue.Encode()
	data := &OrderStatusRes{}
	resp := &response{
		Data: data,
	}
	err := c.queryAndCheckResponse(uri, &req.ReqHeader, nil, resp)
	if err != nil {
		return nil, err
	}
	log.Println(resp.Code, resp.Msg)
	if resp.Code != 0 {
		return nil, fmt.Errorf(resp.Msg)
	}
	return data, nil
}

type ReqCheckUserInSiyu struct {
	ReqHeader
	RxOpenID string `json:"rx_open_id"`
	CpUserID string `json:"cp_user_id"`
}

// CheckUserInSiyu 检查用户是否在私域用户池
func (c *ClientV2) CheckUserInSiyu(req *ReqCheckUserInSiyu) (*RespUserInSiyu, error) {
	if req.CpUserID == "" && req.RxOpenID == "" {
		return nil, ErrInvalidOpenID
	}

	ret := &RespUserInSiyu{}
	resp := &response{Data: ret}

	err := c.queryAndCheckResponse(apiThirdPartySiyu, &req.ReqHeader, &ArgsUserInSiyu{
		RxOpenID: req.RxOpenID,
		CPUserID: req.CpUserID,
	}, resp)

	if err != nil {
		return nil, err
	}
	return ret, nil
}
