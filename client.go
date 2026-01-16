// Copyright (c) 2021. Homeland Interactive Technology Ltd. All rights reserved.

package ruixuego

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	url2 "net/url"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/valyala/fasthttp"
)

const defaultStatus = -1

const (
	headerTraceID     = "ruixue-traceid"     // 请求唯一标识
	headerCPID        = "ruixue-cpid"        // CPID
	headerTimestamp   = "ruixue-cpts"        // 时间戳
	headerSign        = "ruixue-cpsign"      // 签名
	headerVersion     = "ruixue-version"     // SDK 版本
	headerDataCount   = "ruixue-datacount"   // 大数据批量上传时的数据条数
	HeaderProductID   = "ruixue-productid"   // 产品
	HeaderChannelID   = "ruixue-channelid"   // 渠道
	HeaderServiceMark = "ruixue-servicemark" // 用于区分区服信息
	HeaderNameRegion  = "ruixue-region"      // 支持分区
	HeaderLanguage    = "ruixue-language"    // 语言
	headerMethod      = "Method"             // 临时设置请求方法
)

const (
	apiSetCustom                          = "/v1/social/serverapi/setcustom"
	apiAddRelation                        = "/v1/social/serverapi/addrelation"
	apiDelRelation                        = "/v1/social/serverapi/deleterelation"
	apiUpdateRelationRemarks              = "/v1/social/serverapi/updaterelationremarks"
	apiRelationList                       = "/v1/social/serverapi/relationlist"
	apiHasRelation                        = "/v1/social/serverapi/hasrelation"
	apiAddFriend                          = "/v1/social/serverapi/addfriend"
	apiDelFriend                          = "/v1/social/serverapi/delfriend"
	apiUpdateFriendRemarks                = "/v1/social/serverapi/updatefriendremarks"
	apiFriendList                         = "/v1/social/serverapi/friendlist"
	apiIsFriend                           = "/v1/social/serverapi/isfriend"
	apiLBSUpdate                          = "/v1/social/serverapi/lbsupdate"
	apiLBSDelete                          = "/v1/social/serverapi/lbsdelete"
	apiLBSRadius                          = "/v1/social/serverapi/lbsradius"
	apiCreateRank                         = "/v1/social/serverapi/createrank"
	apiCloseRank                          = "/v1/social/serverapi/closerank"
	apiRankAddScore                       = "/v1/social/serverapi/rankaddscore"
	apiRankSetScore                       = "/v1/social/serverapi/ranksetscore"
	apiQueryUserRank                      = "/v1/social/serverapi/queryuserrank"
	apiGetRankList                        = "/v1/social/serverapi/getranklist"
	apiFriendsRank                        = "/v1/social/serverapi/friendsrank"
	apiRankDeleteUser                     = "/v1/social/serverapi/deleteuserscore"
	apiGetRealtionUser                    = "/v1/social/serverapi/getrelationuser"
	apiRankDetail                         = "/v1/social/serverapi/rankdetail"
	apiAllRankIDList                      = "/v1/social/serverapi/getallranklist"
	apiBigDataTrack                       = "/v1/data/api/track"
	apiIMSLogin                           = "/v1/ims/server/login"
	apiIMSSendMessage                     = "/v1/ims/server/sendmessage"
	apiIMSGetHistory                      = "/v1/ims/server/gethistory"
	apiIMSCreateConversation              = "/v1/ims/server/createconversation"
	apiIMSUpdateConversation              = "/v1/ims/server/updateconversation"
	apiIMSDeleteConversation              = "/v1/ims/server/deleteconversation"
	apiIMSGetConversation                 = "/v1/ims/server/getconversation"
	apiIMSJoinConversation                = "/v1/ims/server/joinconversation"
	apiIMSLeaveConversation               = "/v1/ims/server/leaveconversation"
	apiIMSUpdateConversationUserData      = "/v1/ims/server/updateconversatonuserdata"
	apiIMSConversationUserList            = "/v1/ims/server/conversationuserlist"
	apiIMSChannelUsersCount               = "/v1/ims/server/getchanneluserscount"
	apiPusherPush                         = "/v1/pusher/push/push"
	apiRiskTextScan                       = "/v1/risk/content/text/scan"
	apiRiskImageScan                      = "/v1/risk/content/image/scan"
	apiReportCustomAction                 = "/v1/attribution/user/custom_action"
	apiPassportUpdateCPUserID             = "/v1/passport/users/update_cpuserid"
	apiRiskRealAuthCheck                  = "/v1/risk/auth_check"
	apiOperationToolsExtensionExchange    = "/v1/operationtoolsapi/extension/exchange"
	apiOperationToolsExtensionGameDisplay = "/v1/operationtoolsapi/extension/game_display"
	apiOrderInfoByNo                      = "/v1/ke/api/trade_query" // 获取订单信息 --- IGNORE ---
	apiThirdPartySiyu                     = "/v1/thirdparty/service_api/check_user_in_siyu"
	apiReportCPRole                       = "/v1/report/cp/role"
)

var defaultClient *Client

func NewClient() (c *Client, err error) {
	c = &Client{
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
	producer   *Producer
}

func (c *Client) GetProductID() string {
	return config.ProductID
}

func (c *Client) GetChannelID() string {
	return config.ChannelID
}

func (c *Client) GetRegion() string {
	return config.Region
}

func (c *Client) GetLanguage() string {
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
func (c *Client) Close() error {
	if c.producer != nil {
		return c.producer.Close()
	}
	return nil
}

func (c *Client) getRequest(header *ReqHeader, withoutSign ...bool) (string, *fasthttp.Request) {
	traceID, cpID, ts := uuid.New().String(),
		strconv.FormatUint(uint64(config.CPID), 10),
		strconv.FormatInt(time.Now().Unix(), 10)

	req := GetRequest()
	req.Header.Add("user-agent", "ruixue-go-sdk")
	req.Header.Add(headerVersion, Version)
	req.Header.Add(headerTraceID, traceID)
	req.Header.Add(headerCPID, cpID)
	req.Header.Add(HeaderProductID, c.GetProductID())
	req.Header.Add(HeaderChannelID, c.GetChannelID())
	req.Header.Add(HeaderServiceMark, config.ServiceMark)
	req.Header.Add(headerTimestamp, ts)
	req.Header.Add(HeaderNameRegion, c.GetRegion())
	if c.GetLanguage() != "" {
		req.Header.Add(HeaderLanguage, c.GetLanguage())
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

func (c *Client) getAndCheckResponse(url string, args map[string]string,
	header *ReqHeader, resp *response, compress ...bool) error {

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

func (c *Client) queryAndCheckResponse(
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

func (c *Client) query(
	path string, header *ReqHeader, arg, ret interface{}, compress ...bool) (string, error) {
	traceID, req := c.getRequest(header)
	_, err := c.queryCode(path, header, req, config.Timeout, arg, ret, compress...)
	return traceID, err
}

func (c *Client) queryCode(
	path string, header *ReqHeader, req *fasthttp.Request, timeout time.Duration,
	arg, ret interface{}, compress ...bool) (int, error) {

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
			buf, err := GzipCompressV2(b)
			if err != nil {
				return code, err
			}
			req.Header.Set("content-encoding", "gzip")
			req.SetBody(buf)
		} else {
			req.SetBody(b)
		}
	}

	// 设置请求方式
	if header != nil {
		if method, ok := header.Header[headerMethod]; ok {
			req.Header.SetMethod(method)
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
func (c *Client) queryAndCheckResponseWithProductIDAndChannelID(
	path string, header *ReqHeader, req interface{}, resp *response,
	productID, channelID string, compress ...bool) error {

	if resp == nil {
		resp = &response{}
	}

	traceID, err := c.queryWithProductIDAndChannelID(path, header,
		req, resp, productID, channelID, compress...)
	if err != nil {
		return errWithTraceID(err, traceID)
	}

	err = c.checkResponse(resp)
	if err != nil {
		return errWithTraceID(err, traceID)
	}

	return nil
}

func (c *Client) queryWithProductIDAndChannelID(
	path string, header *ReqHeader, arg, ret interface{},
	productID, channelID string, compress ...bool) (string, error) {

	traceID, req := c.getRequest(header)
	c.queryAddProductIDAndChannelID(req, productID, channelID)
	_, err := c.queryCode(path, header, req, config.Timeout, arg, ret, compress...)
	return traceID, err
}

func (c *Client) queryAddProductIDAndChannelID(
	req *fasthttp.Request, productID, channelID string) {
	req.Header.Add(HeaderProductID, productID)
	req.Header.Add(HeaderChannelID, channelID)
}

// SetCustom 给用户设置社交模块的自定义信息
// openID 为 瑞雪opendid
// cpUserID cp侧用户id
func (c *Client) SetCustom(productID, openID, cpUserID, custom string) error {
	if openID == "" && cpUserID == "" {
		return ErrInvalidOpenID
	}
	if productID == "" {
		return ErrInvalidProductID
	}

	return c.queryAndCheckResponse(apiSetCustom, &ReqHeader{}, &argCustom{
		ProductID: productID,
		OpenID:    openID,
		CPUserID:  cpUserID,
		Custom:    custom,
	}, nil)
}

// SetCustomV2 给用户设置社交模块的自定义信息
// OpenID 为 瑞雪opendid
// CpUserID cp侧用户id
func (c *Client) SetCustomV2(req *ReqSetCustom) error {
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

// AddRelation 添加自定义关系
// remarks[0] openID 用户给 targetOpenID 用户设置的备注
// remarks[1] targetOpenID 用户给 openID 用户设置的备注
func (c *Client) AddRelation(
	types RelationTypes, openID, userID, targetOpenID, targetUserID string, remarks ...string) error {
	if (userID == "" && openID == "") || (targetUserID == "" && targetOpenID == "") {
		return ErrInvalidOpenID
	}
	if len(types) == 0 {
		return ErrInvalidType
	}

	arg := &argRelation{
		Types:          types,
		OpenID:         openID,
		CPUserID:       userID,
		Target:         targetOpenID,
		TargetCPUserID: targetUserID,
	}
	if len(remarks) > 0 {
		arg.TargetRemarks = remarks[0]
	}
	if len(remarks) > 1 {
		arg.UserRemarks = remarks[1]
	}

	return c.queryAndCheckResponse(apiAddRelation, &ReqHeader{}, arg, nil)
}

// AddRelationV2 添加自定义关系
// remarks[0] OpenID 用户给 targetOpenID 用户设置的备注
// remarks[1] targetOpenID 用户给 OpenID 用户设置的备注
func (c *Client) AddRelationV2(req *ReqAddRelation) error {
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

// DelRelation 删除自定义关系
func (c *Client) DelRelation(
	types RelationTypes, openID, userID, targetOpenID, targetUserID string) error {
	if (userID == "" && openID == "") || (targetUserID == "" && targetOpenID == "") {
		return ErrInvalidOpenID
	}
	if len(types) == 0 {
		return ErrInvalidType
	}

	return c.queryAndCheckResponse(apiDelRelation, &ReqHeader{}, &argRelation{
		Types:          types,
		OpenID:         openID,
		CPUserID:       userID,
		Target:         targetOpenID,
		TargetCPUserID: targetUserID,
	}, nil)
}

// DelRelationV2 删除自定义关系
func (c *Client) DelRelationV2(req *ReqDelRelation) error {
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

// UpdateRelationRemarks 更新自定关系备注
func (c *Client) UpdateRelationRemarks(
	typ, openID, userID, targetOpenID, targetUserID string, remarks string) error {
	if (userID == "" && openID == "") || (targetUserID == "" && targetOpenID == "") {
		return ErrInvalidOpenID
	}
	if typ == "" {
		return ErrInvalidType
	}

	return c.queryAndCheckResponse(apiUpdateRelationRemarks, &ReqHeader{}, &argRelation{
		Type:           typ,
		OpenID:         openID,
		CPUserID:       userID,
		Target:         targetOpenID,
		TargetCPUserID: targetUserID,
		TargetRemarks:  remarks,
	}, nil)
}

// UpdateRelationRemarksV2 更新自定关系备注
func (c *Client) UpdateRelationRemarksV2(req *ReqUpdateRelationRemarks) error {
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

// RelationList 获取自定关系列表
func (c *Client) RelationList(typ, openID, userID string) ([]*RelationUser, error) {
	if openID == "" && userID == "" {
		return nil, ErrInvalidOpenID
	}
	if typ == "" {
		return nil, ErrInvalidType
	}

	ret := make([]*RelationUser, 0)
	resp := &response{Data: &ret}

	err := c.queryAndCheckResponse(apiRelationList, &ReqHeader{}, &argRelation{
		Type:     typ,
		OpenID:   openID,
		CPUserID: userID,
	}, resp)

	if err != nil {
		return nil, err
	}

	return ret, nil
}

// RelationListV2 获取自定关系列表
func (c *Client) RelationListV2(req *ReqRelationList) ([]*RelationUser, error) {
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

// HasRelation 判断 Target 是否与 User 存在指定关系
func (c *Client) HasRelation(typ, openID, userID, targetOpenID, targetUserID string) (bool, error) {
	ret := false
	resp := &response{Data: &ret}

	err := c.queryAndCheckResponse(apiHasRelation, &ReqHeader{}, &argRelation{
		Type:           typ,
		OpenID:         openID,
		CPUserID:       userID,
		Target:         targetOpenID,
		TargetCPUserID: targetUserID,
	}, resp)

	return ret, err
}

// HasRelationV2 判断 Target 是否与 User 存在指定关系
func (c *Client) HasRelationV2(req *ReqHasRelation) (bool, error) {
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

// AddFriend 添加好友
// remarks[0] openID 用户给 targetOpenID 用户设置的备注
// remarks[1] targetOpenID 用户给 openID 用户设置的备注
func (c *Client) AddFriend(
	openID, userID, targetOpenID, targetUserID string, remarks ...string) error {
	if (userID == "" && openID == "") || (targetUserID == "" && targetOpenID == "") {
		return ErrInvalidOpenID
	}

	arg := &argRelation{
		OpenID:         openID,
		CPUserID:       userID,
		Target:         targetOpenID,
		TargetCPUserID: targetUserID,
	}
	if len(remarks) > 0 {
		arg.TargetRemarks = remarks[0]
	}
	if len(remarks) > 1 {
		arg.UserRemarks = remarks[1]
	}

	return c.queryAndCheckResponse(apiAddFriend, &ReqHeader{}, arg, nil)
}

// AddFriendV2 添加好友
// remarks[0] OpenID 用户给 targetOpenID 用户设置的备注
// remarks[1] targetOpenID 用户给 OpenID 用户设置的备注
func (c *Client) AddFriendV2(req *ReqAddFriend) error {
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

// DelFriend 删除好友
func (c *Client) DelFriend(
	openID, userID, targetOpenID, targetUserID string) error {
	if (userID == "" && openID == "") || (targetUserID == "" && targetOpenID == "") {
		return ErrInvalidOpenID
	}
	return c.queryAndCheckResponse(apiDelFriend, &ReqHeader{}, &argRelation{
		OpenID:         openID,
		CPUserID:       userID,
		Target:         targetOpenID,
		TargetCPUserID: targetUserID,
	}, nil)
}

// DelFriendV2 删除好友
func (c *Client) DelFriendV2(req *ReqDelFriend) error {
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

// UpdateFriendRemarks 更新好友备注
func (c *Client) UpdateFriendRemarks(
	openID, userID, targetOpenID, targetUserID string, remarks string) error {
	if (userID == "" && openID == "") || (targetUserID == "" && targetOpenID == "") {
		return ErrInvalidOpenID
	}
	return c.queryAndCheckResponse(apiUpdateFriendRemarks, &ReqHeader{}, &argRelation{
		OpenID:         openID,
		CPUserID:       userID,
		Target:         targetOpenID,
		TargetCPUserID: targetUserID,
		TargetRemarks:  remarks,
	}, nil)
}

// UpdateFriendRemarksV2 更新好友备注
func (c *Client) UpdateFriendRemarksV2(req *ReqUpdateFriendRemarks) error {
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

// FriendList 获取好友列表
func (c *Client) FriendList(openID, userID string) ([]*RelationUser, error) {
	if userID == "" && openID == "" {
		return nil, ErrInvalidOpenID
	}

	ret := make([]*RelationUser, 0)
	resp := &response{Data: &ret}

	err := c.queryAndCheckResponse(apiFriendList, &ReqHeader{}, &argRelation{
		OpenID:   openID,
		CPUserID: userID,
	}, resp)

	if err != nil {
		return nil, err
	}

	return ret, nil
}

// FriendListV2 获取好友列表
func (c *Client) FriendListV2(req *ReqFriendList) ([]*RelationUser, error) {
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

// GetRelationUser 查询好友信息
func (c *Client) GetRelationUser(typ, openID, userID,
	targetOpenID, targetUserID string) (*RelationUser, error) {

	if (userID == "" && openID == "") || (targetUserID == "" && targetOpenID == "") {
		return nil, ErrInvalidOpenID
	}

	if typ == "" {
		return nil, ErrInvalidType
	}

	ret := &RelationUser{}
	resp := &response{Data: ret}

	err := c.queryAndCheckResponse(apiGetRealtionUser, &ReqHeader{}, &argRelation{
		OpenID:         openID,
		CPUserID:       userID,
		Target:         targetOpenID,
		TargetCPUserID: targetUserID,
		Type:           typ,
	}, resp)

	return ret, err
}

// GetRelationUserV2 查询好友信息
func (c *Client) GetRelationUserV2(req *ReqGetRelationUser) (*RelationUser, error) {
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

// IsFriend 判断 Target 是否为 User 的好友
func (c *Client) IsFriend(openID, userID, targetOpenID, targetUserID string) (bool, error) {
	if (userID == "" && openID == "") || (targetUserID == "" && targetOpenID == "") {
		return false, ErrInvalidOpenID
	}
	ret := false
	resp := &response{Data: &ret}

	err := c.queryAndCheckResponse(apiIsFriend, &ReqHeader{}, &argRelation{
		OpenID:         openID,
		CPUserID:       userID,
		Target:         targetOpenID,
		TargetCPUserID: targetUserID,
	}, resp)

	return ret, err
}

// IsFriendV2 判断 Target 是否为 User 的好友
func (c *Client) IsFriendV2(req *ReqIsFriend) (bool, error) {
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

// LBSUpdate 更新 WGS84 坐标
//
//	types 为 CP	自定义坐标分组, 比如可以同时将用户加入到 all 和 female 两个列表中
func (c *Client) LBSUpdate(openID, userID string, types []string, lon, lat float64) error {
	if userID == "" && openID == "" {
		return ErrInvalidOpenID
	}

	if len(types) == 0 {
		return ErrInvalidType
	}

	return c.queryAndCheckResponse(apiLBSUpdate, &ReqHeader{}, &argLocation{
		OpenID:    openID,
		CPUserID:  userID,
		Types:     types,
		Longitude: lon,
		Latitude:  lat,
	}, nil)
}

// LBSUpdateV2 更新 WGS84 坐标
//
//	types 为 CP	自定义坐标分组, 比如可以同时将用户加入到 all 和 female 两个列表中
func (c *Client) LBSUpdateV2(req *ReqLBSUpdate) error {
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

// LBSDelete 删除 WGS84 坐标
func (c *Client) LBSDelete(openID, userID string, types []string) error {
	if userID == "" && openID == "" {
		return ErrInvalidOpenID
	}
	if len(types) == 0 {
		return ErrInvalidType
	}

	return c.queryAndCheckResponse(apiLBSDelete, &ReqHeader{}, &argLocation{
		OpenID:   openID,
		CPUserID: userID,
		Types:    types,
	}, nil)
}

// LBSDeleteV2 删除 WGS84 坐标
func (c *Client) LBSDeleteV2(req *ReqLBSDelete) error {
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

// LBSRadius 获取附近的人列表
func (c *Client) LBSRadius(
	openID, userID, typ string,
	lon, lat, radius float64,
	page, pageSize int,
	count ...int) ([]*RelationUser, error) {

	if userID == "" && openID == "" {
		return nil, ErrInvalidOpenID
	}
	if typ == "" {
		return nil, ErrInvalidType
	}

	ret := make([]*RelationUser, 0)
	resp := &response{Data: &ret}
	arg := &argLocation{
		OpenID:    openID,
		CPUserID:  userID,
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

	err := c.queryAndCheckResponse(apiLBSRadius, &ReqHeader{}, arg, resp)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

// LBSRadiusV2 获取附近的人列表
func (c *Client) LBSRadiusV2(req *ReqLBSRadius) ([]*RelationUser, error) {

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
func (c *Client) Tracks(
	devicecode, distinctID string, opts ...BigdataOptions) error {
	return c.producer.Tracks(devicecode, distinctID, opts...)
}

// Track 将埋点数据上报给瑞雪云
func (c *Client) Track(track *ReqTrack) (int, error) {
	if len(track.Data) == 0 {
		return defaultStatus, nil
	}

	traceID, req := c.getRequest(&track.ReqHeader, true)
	ret := &response{}
	req.Header.Add(headerDataCount, Itoa(track.LogCount))
	code, err := c.queryCode(apiBigDataTrack, &track.ReqHeader, req, config.TrackTimeout, track.Data, ret, track.Compress)
	if err != nil {
		return code, errWithTraceID(err, traceID)
	}
	err = c.checkResponse(ret)
	if err != nil {
		return code, errWithTraceID(err, traceID)
	}
	return code, nil
}

// SyncTrack 同步接口 直接将埋点数据上报给瑞雪云
// 前提要设置好 config
func (c *Client) SyncTrack(devicecode, distinctID string, opts ...BigdataOptions) error {
	if devicecode == "" && distinctID == "" {
		return ErrInvalidDevicecode
	}

	logData := &BigDataLog{
		DistinctID: distinctID,
		Devicecode: devicecode,
	}
	for _, opt := range opts {
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

	traceID, req := c.getRequest(&ReqHeader{}, true)
	ret := &response{}
	req.Header.Add(headerDataCount, Itoa(1))
	_, err = c.queryCode(apiBigDataTrack, &ReqHeader{}, req, config.TrackTimeout, b, ret, !config.BigData.DisableCompress)
	if err != nil {
		return errWithTraceID(err, traceID)
	}
	err = c.checkResponse(ret)
	if err != nil {
		return errWithTraceID(err, traceID)
	}
	return nil
}

// SyncTrackV2 同步接口 直接将埋点数据上报给瑞雪云
// 前提要设置好 config
func (c *Client) SyncTrackV2(track *ReqSyncTrack) error {
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
	_, err = c.queryCode(apiBigDataTrack, &track.ReqHeader, req, config.TrackTimeout,
		b, ret, !config.BigData.DisableCompress)

	if err != nil {
		return errWithTraceID(err, traceID)
	}
	err = c.checkResponse(ret)
	if err != nil {
		return errWithTraceID(err, traceID)
	}
	return nil
}

// CreateRank 创建排行榜
func (c *Client) CreateRank(rankID string, startTime, destroyTime time.Time) error {
	if rankID == "" {
		return ErrInvalidOpenID
	}

	err := c.queryAndCheckResponse(apiCreateRank, &ReqHeader{}, &rankAPIArg{
		RankID:      rankID,
		StartTime:   startTime.Format(time.RFC3339),
		DestroyTime: destroyTime.Format(time.RFC3339),
	}, nil)

	return err
}

// CreateRankV2 创建排行榜
func (c *Client) CreateRankV2(req *ReqCreateRank) error {
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

// CloseRank 关闭排行榜
func (c *Client) CloseRank(rankID string) error {
	if rankID == "" {
		return ErrInvalidOpenID
	}

	err := c.queryAndCheckResponse(apiCloseRank, &ReqHeader{}, &rankAPIArg{
		RankID: rankID,
	}, nil)

	return err
}

// CloseRankV2 关闭排行榜
func (c *Client) CloseRankV2(req *ReqCloseRank) error {
	if req.RankID == "" {
		return ErrInvalidOpenID
	}

	err := c.queryAndCheckResponse(apiCloseRank, &req.ReqHeader, &rankAPIArg{
		RankID: req.RankID,
	}, nil)

	return err
}

// RankAddScore 用户添加分数
func (c *Client) RankAddScore(rankID string, openId, cpUserID string, score int64) error {
	if rankID == "" || (openId == "" && cpUserID == "") {
		return ErrInvalidOpenID
	}

	err := c.queryAndCheckResponse(apiRankAddScore, &ReqHeader{}, &rankAPIArg{
		RankID:   rankID,
		OpenID:   openId,
		CPUserID: cpUserID,
		Score:    score,
	}, nil)

	return err
}

// RankAddScoreV2 用户添加分数
func (c *Client) RankAddScoreV2(req *ReqRankAddScore) error {
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

// RankSetScore 用户设置分数
func (c *Client) RankSetScore(rankID string, openId, cpUserID string, score int64) error {
	if rankID == "" || (openId == "" && cpUserID == "") {
		return ErrInvalidOpenID
	}

	err := c.queryAndCheckResponse(apiRankSetScore, &ReqHeader{}, &rankAPIArg{
		RankID:   rankID,
		OpenID:   openId,
		CPUserID: cpUserID,
		Score:    score,
	}, nil)

	return err
}

// RankSetScoreV2 用户设置分数
func (c *Client) RankSetScoreV2(req *ReqRankSetScore) error {
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

// DeleteRankUser 删除排行榜用户
func (c *Client) DeleteRankUser(rankID string, openId, cpUserID string) error {
	if rankID == "" || (openId == "" && cpUserID == "") {
		return ErrInvalidOpenID
	}
	err := c.queryAndCheckResponse(apiRankDeleteUser, &ReqHeader{}, &rankAPIArg{
		RankID:   rankID,
		OpenID:   openId,
		CPUserID: cpUserID,
	}, nil)

	return err
}

// DeleteRankUserV2 删除排行榜用户
func (c *Client) DeleteRankUserV2(req *ReqDeleteRankUser) error {
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

// QueryUserRank 查询用户排行情况
func (c *Client) QueryUserRank(rankID string, openId, cpUserID string) (*RankMember, error) {
	if rankID == "" || (openId == "" && cpUserID == "") {
		return nil, ErrInvalidOpenID
	}
	ret := &RankMember{}
	resp := &response{Data: ret}

	err := c.queryAndCheckResponse(apiQueryUserRank, &ReqHeader{}, &rankAPIArg{
		RankID:   rankID,
		OpenID:   openId,
		CPUserID: cpUserID,
	}, resp)

	return ret, err
}

// QueryUserRankV2 查询用户排行情况
func (c *Client) QueryUserRankV2(req *ReqQueryUserRank) (*RankMember, error) {
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

// GetRankList 查询排行榜
func (c *Client) GetRankList(rankID string, start, end int32) ([]*RankMember, error) {
	if rankID == "" {
		return nil, ErrInvalidOpenID
	}

	var ret []*RankMember
	resp := &response{Data: &ret}

	err := c.queryAndCheckResponse(apiGetRankList, &ReqHeader{}, &rankAPIArg{
		RankID:    rankID,
		StartRank: start,
		EndRank:   end,
	}, resp)

	return ret, err
}

// GetRankListV2 查询排行榜
func (c *Client) GetRankListV2(req *ReqGetRankList) ([]*RankMember, error) {
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

// GetFriendRankList 查询好友排行榜
func (c *Client) GetFriendRankList(rankID string, openId, cpUserID string) ([]*RankMember, error) {
	if rankID == "" || (openId == "" && cpUserID == "") {
		return nil, ErrInvalidOpenID
	}

	var ret []*RankMember
	resp := &response{Data: &ret}

	err := c.queryAndCheckResponse(apiFriendsRank, &ReqHeader{}, &rankAPIArg{
		RankID:   rankID,
		OpenID:   openId,
		CPUserID: cpUserID,
	}, resp)

	return ret, err
}

// GetFriendRankListV2 查询好友排行榜
func (c *Client) GetFriendRankListV2(req *ReqGetFriendRankList) ([]*RankMember, error) {
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

// GetRankDetail 查询排行榜详情
func (c *Client) GetRankDetail(rankID string) (*RespRankDetail, error) {
	if rankID == "" {
		return nil, ErrInvalidOpenID
	}

	ret := &RespRankDetail{}
	resp := &response{Data: ret}

	err := c.queryAndCheckResponse(apiRankDetail, &ReqHeader{}, &rankAPIArg{
		RankID: rankID,
	}, resp)

	return ret, err
}

// GetRankDetailV2 查询排行榜详情
func (c *Client) GetRankDetailV2(req *ReqGetRankDetail) (*RespRankDetail, error) {
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
func (c *Client) GetAllRankIDList() (*RespAllRankID, error) {

	ret := &RespAllRankID{}
	resp := &response{Data: ret}

	err := c.queryAndCheckResponse(apiAllRankIDList, &ReqHeader{}, nil, resp)

	return ret, err
}

// GetAllRankIDListV2 查询所有排行ID
func (c *Client) GetAllRankIDListV2(req *ReqHeader) (*RespAllRankID, error) {
	ret := &RespAllRankID{}
	resp := &response{Data: ret}
	err := c.queryAndCheckResponse(apiAllRankIDList, req, nil, resp)
	return ret, err
}

// IMSLogin ims 登陆接口
func (c *Client) IMSLogin(req *IMSLoginReq) (*IMSLoginResp, error) {
	ret := &IMSLoginResp{}
	resp := &response{Data: ret}
	err := c.queryAndCheckResponse(apiIMSLogin, &req.ReqHeader, req, resp)
	return ret, err
}

// IMSSendMessage 发送消息
func (c *Client) IMSSendMessage(req *IMSMessage) (*IMSMessageAck, error) {
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
func (c *Client) IMSGetHistory(req *IMSHistoryReq) (*IMSHistoryResp, error) {
	ret := &IMSHistoryResp{}
	resp := &response{Data: ret}
	err := c.queryAndCheckResponse(apiIMSGetHistory, &req.ReqHeader, req, resp)
	return ret, err
}

// IMSCreateConversation 创建会话
func (c *Client) IMSCreateConversation(req *IMSCreateConvReq) error {
	return c.queryAndCheckResponse(apiIMSCreateConversation, &req.ReqHeader, req, nil)
}

// IMSUpdateConversation 更新会话信息
func (c *Client) IMSUpdateConversation(req *IMSUpdateConvReq) error {
	return c.queryAndCheckResponse(apiIMSUpdateConversation, &req.ReqHeader, req, nil)
}

// IMSDeleteConversation 删除会话
func (c *Client) IMSDeleteConversation(req *IMSConvDeleteReq) error {
	return c.queryAndCheckResponse(apiIMSDeleteConversation, &req.ReqHeader, req, nil)
}

// IMSGetConversation 获取会话信息
func (c *Client) IMSGetConversation(req *IMSGetConversationReq) (*IMSConversation, error) {
	ret := &IMSConversation{}
	resp := &response{Data: ret}
	err := c.queryAndCheckResponse(apiIMSGetConversation, &req.ReqHeader, req, resp)
	return ret, err
}

// IMSJoinConversation 加入会话
func (c *Client) IMSJoinConversation(req *IMSJoinConversationReq) error {
	return c.queryAndCheckResponse(apiIMSJoinConversation, &req.ReqHeader, req, nil)
}

// IMSLeaveConversation 离开会话
func (c *Client) IMSLeaveConversation(req *IMSLeaveConversationReq) error {
	return c.queryAndCheckResponse(apiIMSLeaveConversation, &req.ReqHeader, req, nil)
}

// IMSUpdateConversationUserData 更新会话内用户信息
func (c *Client) IMSUpdateConversationUserData(req *IMSUpdateConvUserDataReq) error {
	return c.queryAndCheckResponse(apiIMSUpdateConversationUserData, &req.ReqHeader, req, nil)
}

// IMSConversationUserList 获取会话中成员列表
func (c *Client) IMSConversationUserList(req *IMSConversationUserListReq) ([]*IMSConversation, error) {
	ret := make([]*IMSConversation, 0)
	resp := &response{Data: &ret}
	err := c.queryAndCheckResponse(apiIMSConversationUserList, &req.ReqHeader, req, resp)
	return ret, err
}

// IMSChannelUsersCount 获取频道会话中玩家数量
func (c *Client) IMSChannelUsersCount(channelConvIds []string) (map[string]int64, error) {
	ret := make(map[string]int64)
	resp := &response{Data: &ret}

	req := &IMSChannelUsesCountReq{
		ConversationIDs: channelConvIds,
	}

	err := c.queryAndCheckResponse(apiIMSChannelUsersCount, &ReqHeader{}, req, resp)
	return ret, err
}

// IMSChannelUsersCountV2 获取频道会话中玩家数量
func (c *Client) IMSChannelUsersCountV2(count *ReqIMSChannelUsersCount) (map[string]int64, error) {
	ret := make(map[string]int64)
	resp := &response{Data: &ret}

	req := &IMSChannelUsesCountReq{
		ConversationIDs: count.ChannelConvIds,
	}

	err := c.queryAndCheckResponse(apiIMSChannelUsersCount, &count.ReqHeader, req, resp)
	return ret, err
}

// PusherPush 推送信息
func (c *Client) PusherPush(req *PusherPushReq, productID, channelID string) (*PusherPushRes, error) {
	ret := &PusherPushRes{}
	resp := &response{Data: ret}
	err := c.queryAndCheckResponseWithProductIDAndChannelID(apiPusherPush, &ReqHeader{}, req, resp, productID, channelID)
	return ret, err
}

// PusherPushV2 推送信息
func (c *Client) PusherPushV2(req *ReqPusher) (*PusherPushRes, error) {
	ret := &PusherPushRes{}
	resp := &response{Data: ret}
	err := c.queryAndCheckResponseWithProductIDAndChannelID(apiPusherPush, &req.ReqHeader, req.Req, resp, req.ProductID, req.ChannelID)
	return ret, err
}

// RiskContentTextScan 内容安全文本检查（增强）
func (c *Client) RiskContentTextScan(req *RiskContentTextScanReq) (*RiskContentTextScanResp, error) {
	if req == nil || len(req.Scene) <= 0 || len(req.Content) <= 0 {
		return nil, ErrInvalidParam
	}
	ret := &RiskContentTextScanResp{}
	resp := &response{Data: ret}
	err := c.queryAndCheckResponse(apiRiskTextScan, &req.ReqHeader, req, resp)
	return ret, err
}

// RiskContentImageScan 内容安全图片检查（增强）
func (c *Client) RiskContentImageScan(url string) (*RiskContentImageScanResp, error) {
	if len(url) <= 0 {
		return nil, ErrInvalidParam
	}
	ret := &RiskContentImageScanResp{}
	resp := &response{Data: ret}
	err := c.queryAndCheckResponse(apiRiskImageScan, &ReqHeader{}, &RiskContentImageScanReq{
		URL: url,
	}, resp)
	return ret, err
}

// RiskContentImageScanV2 内容安全图片检查（增强）
func (c *Client) RiskContentImageScanV2(req *RiskContentImageScanReq) (*RiskContentImageScanResp, error) {
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
func (c *Client) ReportCustomAction(openID, action string) error {
	if openID == "" || action == "" {
		return ErrInvalidOpenID
	}
	arg := &ReportCustomAction{
		OpenID: openID,
		Action: action,
	}
	resp := &response{}
	err := c.queryAndCheckResponse(apiReportCustomAction, &ReqHeader{}, arg, resp)
	if err != nil {
		return err
	}
	if resp.Code != 0 {
		return fmt.Errorf(resp.Msg)
	}
	return nil
}

// ReportCustomActionV2 投放归因上报自定义action， 例如广点通小游戏创角， action传 CREATE_ROLE
// https://developers.e.qq.com/docs/guide/conversion/new_version/Mini_Game_api
func (c *Client) ReportCustomActionV2(req *ReportCustomAction) error {
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
func (c *Client) UpdateCPuserID(openID, cpUserID string) error {
	if len(cpUserID) == 0 || len(openID) == 0 {
		return ErrInvalidCPuserID
	}
	arg := &UpdateCPUserIDRequest{}
	arg.OpenID = openID
	arg.CPUserID = cpUserID
	resp := &response{}
	err := c.queryAndCheckResponse(apiPassportUpdateCPUserID, &ReqHeader{}, arg, resp)
	if err != nil {
		return err
	}
	if resp.Code != 0 {
		return fmt.Errorf(resp.Msg)
	}
	return nil
}

// UpdateCPuserIDV2 用于瑞雪侧更新cp侧的user_id
func (c *Client) UpdateCPuserIDV2(req *UpdateCPUserIDRequest) error {
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
func (c *Client) RealAuth(productID, idCard, realName string) (*RealAuthResponse, error) {
	if len(productID) == 0 || len(idCard) == 0 || len(realName) == 0 {
		return nil, ErrInvalidParam
	}
	arg := &RealAuthReq{}
	arg.IDCard = idCard
	arg.RealName = realName
	arg.ProductID = productID
	arg.CPID = config.CPID
	data := &RealAuthResponse{}
	resp := &response{
		Data: data,
	}
	err := c.queryAndCheckResponse(apiRiskRealAuthCheck, &ReqHeader{}, arg, resp)
	if err != nil {
		return nil, err
	}
	log.Println(resp.Code, resp.Msg)
	if resp.Code != 0 {
		return nil, fmt.Errorf(resp.Msg)
	}
	return data, nil
}

// RealAuthV2 实名检测
func (c *Client) RealAuthV2(arg *RealAuthReq) (*RealAuthResponse, error) {
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
func (c *Client) ExtensionExchange(arg *ExtensionExchangeReq) ([]*ExtensionProp, error) {
	if len(arg.CdKey) == 0 || len(arg.CpUserID) == 0 {
		return nil, ErrInvalidParam
	}
	data := []*ExtensionProp{}
	resp := &response{
		Data: &data,
	}
	err := c.queryAndCheckResponse(apiOperationToolsExtensionExchange, &arg.ReqHeader, arg, resp)
	if err != nil {
		return nil, err
	}
	if resp.Code != 0 {
		return nil, fmt.Errorf(resp.Msg)
	}
	return data, nil
}

// ExtensionGameDisplay 主播获取游戏内显示码
func (c *Client) ExtensionGameDisplay(gameID string) (*GameDisplayWelfareCodeInfoExp, error) {
	if gameID == "" {
		return nil, ErrInvalidParam
	}
	dataValue := make(url2.Values)
	dataValue.Add("game_id", gameID)
	url := apiOperationToolsExtensionGameDisplay
	uri := url + "?" + dataValue.Encode()
	data := &GameDisplayWelfareCodeInfoExp{}
	resp := &response{
		Data: data,
	}
	err := c.queryAndCheckResponse(uri, &ReqHeader{}, nil, resp)
	if err != nil {
		return nil, err
	}
	log.Println(resp.Code, resp.Msg)
	if resp.Code != 0 {
		return nil, fmt.Errorf(resp.Msg)
	}
	return data, nil
}

// ExtensionGameDisplayV2 主播获取游戏内显示码
func (c *Client) ExtensionGameDisplayV2(req *ReqExtensionGameDisplay) (*GameDisplayWelfareCodeInfoExp, error) {
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

// TradeOrderStatusByNo 平台订单状态
func (c *Client) TradeOrderStatusByNo(orderNo string) (*OrderStatusRes, error) {
	if orderNo == "" || !strings.HasSuffix(orderNo, "v1") {
		return nil, ErrInvalidParam
	}
	dataValue := make(url2.Values)
	dataValue.Add("order_no", orderNo)
	url := apiOrderInfoByNo
	uri := url + "?" + dataValue.Encode()
	data := &OrderStatusRes{}
	resp := &response{
		Data: data,
	}
	err := c.queryAndCheckResponse(uri, &ReqHeader{}, nil, resp)
	if err != nil {
		return nil, err
	}
	log.Println(resp.Code, resp.Msg)
	if resp.Code != 0 {
		return nil, fmt.Errorf(resp.Msg)
	}
	return data, nil
}

// TradeOrderStatusByNoV2 平台订单状态
func (c *Client) TradeOrderStatusByNoV2(req *ReqTradeOrderStatusByNo) (*OrderStatusRes, error) {
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

// CheckUserInSiyu 检查用户是否在私域用户池
func (c *Client) CheckUserInSiyu(rxOpenID, cpUserID string) (*RespUserInSiyu, error) {
	if cpUserID == "" && rxOpenID == "" {
		return nil, ErrInvalidOpenID
	}

	ret := &RespUserInSiyu{}
	resp := &response{Data: ret}

	err := c.queryAndCheckResponse(apiThirdPartySiyu, &ReqHeader{}, &ArgsUserInSiyu{
		RxOpenID: rxOpenID,
		CPUserID: cpUserID,
	}, resp)

	if err != nil {
		return nil, err
	}
	return ret, nil
}

// CheckUserInSiyuV2 检查用户是否在私域用户池
func (c *Client) CheckUserInSiyuV2(req *ReqCheckUserInSiyu) (*RespUserInSiyu, error) {
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

// CPRoleAdd 角色上报
func (c *Client) CPRoleAdd(args *CPRoleInfo) error {
	if args == nil || args.RxOpenID == "" || args.RegionTag == "" || args.CPRoleID == "" {
		return ErrInvalidOpenID
	}
	resp := &response{}
	err := c.queryAndCheckResponse(apiReportCPRole, &args.ReqHeader, args, resp)

	if err != nil {
		return err
	}
	if resp.Code != 0 {
		return fmt.Errorf("code:%v,msg:%v", resp.Code, resp.Msg)
	}
	return nil
}

// CPRoleUpdate 角色更新
func (c *Client) CPRoleUpdate(args *CPRoleInfo) error {
	if args == nil || args.RxOpenID == "" || args.RegionTag == "" || args.CPRoleID == "" {
		return ErrInvalidOpenID
	}
	resp := &response{}
	args.ReqHeader.Set(headerMethod, "PUT")
	err := c.queryAndCheckResponse(apiReportCPRole, &args.ReqHeader, args, resp)

	if err != nil {
		return err
	}
	if resp.Code != 0 {
		return fmt.Errorf("code:%v,msg:%v", resp.Code, resp.Msg)
	}
	return nil
}

// CPRoleDel 角色删除
func (c *Client) CPRoleDel(args *CPRoleInfoDel) error {
	if args == nil || args.RxOpenID == "" || args.RegionTag == "" || args.CPRoleID == "" {
		return ErrInvalidOpenID
	}
	resp := &response{}
	args.ReqHeader.Set(headerMethod, "DELETE")
	err := c.queryAndCheckResponse(apiReportCPRole, &args.ReqHeader, args, resp)

	if err != nil {
		return err
	}
	if resp.Code != 0 {
		return fmt.Errorf("code:%v,msg:%v", resp.Code, resp.Msg)
	}
	return nil
}

// CPRoleListByOpenID 通过openid查询角色列表
func (c *Client) CPRoleListByOpenID(args *CPRoleList) ([]*CPRoleRes, error) {
	if args == nil || args.RxOpenID == "" {
		return nil, ErrInvalidOpenID
	}

	type CPRoleListRes struct {
		List []*CPRoleRes `json:"list"`
	}

	ret := &CPRoleListRes{}
	resp := &response{Data: ret}
	args.ReqHeader.Set(headerMethod, "GET")
	link := apiReportCPRole + "/list?rx_openid=" + args.RxOpenID + "&extension=" + args.Extension
	err := c.queryAndCheckResponse(link, &args.ReqHeader, args, resp)

	if err != nil {
		return nil, err
	}
	if resp.Code != 0 {
		return nil, fmt.Errorf("code:%v,msg:%v", resp.Code, resp.Msg)
	}
	return ret.List, nil
}
