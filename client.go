// Copyright (c) 2021. Homeland Interactive Technology Ltd. All rights reserved.

package ruixuego

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"ruixuego/bufferpool"

	"github.com/google/uuid"
	"github.com/valyala/fasthttp"
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
	apiSetUserInfo           = "/v1/social/serverapi/setuserinfo"
	apiSetCustom             = "/v1/social/serverapi/setcustom"
	apiAddRelation           = "/v1/social/serverapi/addrelation"
	apiDelRelation           = "/v1/social/serverapi/deleterelation"
	apiUpdateRelationRemarks = "/v1/social/serverapi/updaterelationremarks"
	apiRelationList          = "/v1/social/serverapi/relationlist"
	apiHasRelation           = "/v1/social/serverapi/hasrelation"
	apiAddFriend             = "/v1/social/serverapi/addfriend"
	apiDelFriend             = "/v1/social/serverapi/delfriend"
	apiUpdateFriendRemarks   = "/v1/social/serverapi/updatefriendremarks"
	apiFriendList            = "/v1/social/serverapi/friendlist"
	apiIsFriend              = "/v1/social/serverapi/isfriend"
	apiLBSUpdate             = "/v1/social/serverapi/lbsupdate"
	apiLBSDelete             = "/v1/social/serverapi/lbsdelete"
	apiLBSRadius             = "/v1/social/serverapi/lbsradius"
	apiCreateRank            = "/v1/social/serverapi/createrank"
	apiCloseRank             = "/v1/social/serverapi/closerank"
	apiRankAddScore          = "/v1/social/serverapi/rankaddscore"
	apiRankSetScore          = "/v1/social/serverapi/ranksetscore"
	apiQueryUserRank         = "/v1/social/serverapi/queryuserrank"
	apiGetRankList           = "/v1/social/serverapi/getranklist"
	apiFriendsRank           = "/v1/social/serverapi/friendsrank"

	apiBigDataTrack = "/v1/data/api/track"

	apiIMSLogin                      = "/v1/ims/server/login"
	apiIMSSendMessage                = "/v1/ims/server/sendmessage"
	apiIMSGetHistory                 = "/v1/ims/server/gethistory"
	apiIMSCreateConversation         = "/v1/ims/server/createconversation"
	apiIMSUpdateConversation         = "/v1/ims/server/updateconversation"
	apiIMSDeleteConversation         = "/v1/ims/server/deleteconversation"
	apiIMSGetConversation            = "/v1/ims/server/getconversation"
	apiIMSJoinConversation           = "/v1/ims/server/joinconversation"
	apiIMSLeaveConversation          = "/v1/ims/server/leaveconversation"
	apiIMSUpdateConversationUserData = "/v1/ims/server/updateconversatonuserdata"
	apiIMSConversationUserList       = "/v1/ims/server/conversationuserlist"

	apiPusherPush = "/v1/pusher/push/push"

	apiRiskGreenSyncScan      = "/v1/risk/green/img/syncscan"
	apiRiskGreenAsyncScan     = "/v1/risk/green/img/asyncscan"
	apiRiskGreenGetScanResult = "/v1/risk/green/img/getscanres"
	apiRiskGreenFeedback      = "/v1/risk/green/img/scanfeedback"
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

// Close SDK ????????????????????????????????????????????????, ???????????????????????????
func (c *Client) Close() error {
	if c.producer != nil {
		return c.producer.Close()
	}
	return nil
}

func (c *Client) getRequest(withoutSign ...bool) (string, *fasthttp.Request) {
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
		req.Header.Add(headerSign, GetSign(traceID, ts))
	}

	return traceID, req
}

func (c *Client) queryAndCheckResponse(
	path string, req interface{}, resp *response, compress ...bool) error {

	if resp == nil {
		resp = &response{}
	}

	traceID, err := c.query(path, req, resp, compress...)
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
	path string, arg, ret interface{}, compress ...bool) (string, error) {
	traceID, req := c.getRequest()
	_, err := c.queryCode(path, req, config.Timeout, arg, ret, compress...)
	return traceID, err
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

// SetUserInfo ??????????????????
func (c *Client) SetUserInfo(appID, openID string, userinfo *UserInfo) error {
	if openID == "" {
		return ErrInvalidOpenID
	}
	if appID == "" {
		return ErrInvalidAppID
	}

	userinfo.AppID = appID
	userinfo.OpenID = openID

	return c.queryAndCheckResponse(apiSetUserInfo, userinfo, nil)
}

// SetCustom ?????????????????????????????????????????????
func (c *Client) SetCustom(appID, openID, custom string) error {
	if openID == "" {
		return ErrInvalidOpenID
	}
	if appID == "" {
		return ErrInvalidAppID
	}

	return c.queryAndCheckResponse(apiSetCustom, &argCustom{
		AppID:  appID,
		OpenID: openID,
		Custom: custom,
	}, nil)
}

// AddRelation ?????????????????????
// remarks[0] openID ????????? targetOpenID ?????????????????????
// remarks[1] targetOpenID ????????? openID ?????????????????????
func (c *Client) AddRelation(
	types RelationTypes, openID, targetOpenID string, remarks ...string) error {
	if openID == "" || targetOpenID == "" {
		return ErrInvalidOpenID
	}
	if len(types) == 0 {
		return ErrInvalidType
	}

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

	return c.queryAndCheckResponse(apiAddRelation, arg, nil)
}

// DelRelation ?????????????????????
func (c *Client) DelRelation(
	types RelationTypes, openID, targetOpenID string) error {
	if openID == "" || targetOpenID == "" {
		return ErrInvalidOpenID
	}
	if len(types) == 0 {
		return ErrInvalidType
	}

	return c.queryAndCheckResponse(apiDelRelation, &argRelation{
		Types:  types,
		OpenID: openID,
		Target: targetOpenID,
	}, nil)
}

// UpdateRelationRemarks ????????????????????????
func (c *Client) UpdateRelationRemarks(
	typ, openID, targetOpenID, remarks string) error {
	if openID == "" || targetOpenID == "" {
		return ErrInvalidOpenID
	}
	if typ == "" {
		return ErrInvalidType
	}

	return c.queryAndCheckResponse(apiUpdateRelationRemarks, &argRelation{
		Type:          typ,
		OpenID:        openID,
		Target:        targetOpenID,
		TargetRemarks: remarks,
	}, nil)
}

// RelationList ????????????????????????
func (c *Client) RelationList(typ, openID string) ([]*RelationUser, error) {
	if openID == "" {
		return nil, ErrInvalidOpenID
	}
	if typ == "" {
		return nil, ErrInvalidType
	}

	ret := make([]*RelationUser, 0)
	resp := &response{Data: &ret}

	err := c.queryAndCheckResponse(apiRelationList, &argRelation{
		Type:   typ,
		OpenID: openID,
	}, resp)

	if err != nil {
		return nil, err
	}

	return ret, nil
}

// HasRelation ?????? Target ????????? User ??????????????????
func (c *Client) HasRelation(typ, openID, targetOpenID string) (bool, error) {
	ret := false
	resp := &response{Data: &ret}

	err := c.queryAndCheckResponse(apiHasRelation, &argRelation{
		Type:   typ,
		OpenID: openID,
		Target: targetOpenID,
	}, resp)

	return ret, err
}

// AddFriend ????????????
// remarks[0] openID ????????? targetOpenID ?????????????????????
// remarks[1] targetOpenID ????????? openID ?????????????????????
func (c *Client) AddFriend(
	openID, targetOpenID string, remarks ...string) error {
	if openID == "" || targetOpenID == "" {
		return ErrInvalidOpenID
	}

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

	return c.queryAndCheckResponse(apiAddFriend, arg, nil)
}

// DelFriend ????????????
func (c *Client) DelFriend(
	openID, targetOpenID string) error {
	if openID == "" || targetOpenID == "" {
		return ErrInvalidOpenID
	}

	return c.queryAndCheckResponse(apiDelFriend, &argRelation{
		OpenID: openID,
		Target: targetOpenID,
	}, nil)
}

// UpdateFriendRemarks ??????????????????
func (c *Client) UpdateFriendRemarks(
	openID, targetOpenID, remarks string) error {
	if openID == "" || targetOpenID == "" {
		return ErrInvalidOpenID
	}

	return c.queryAndCheckResponse(apiUpdateFriendRemarks, &argRelation{
		OpenID:        openID,
		Target:        targetOpenID,
		TargetRemarks: remarks,
	}, nil)
}

// FriendList ??????????????????
func (c *Client) FriendList(openID string) ([]*RelationUser, error) {
	if openID == "" {
		return nil, ErrInvalidOpenID
	}

	ret := make([]*RelationUser, 0)
	resp := &response{Data: &ret}

	err := c.queryAndCheckResponse(apiFriendList, &argRelation{
		OpenID: openID,
	}, resp)

	if err != nil {
		return nil, err
	}

	return ret, nil
}

// IsFriend ?????? Target ????????? User ?????????
func (c *Client) IsFriend(openID, targetOpenID string) (bool, error) {
	if openID == "" || targetOpenID == "" {
		return false, ErrInvalidOpenID
	}

	ret := false
	resp := &response{Data: &ret}

	err := c.queryAndCheckResponse(apiIsFriend, &argRelation{
		OpenID: openID,
		Target: targetOpenID,
	}, resp)

	return ret, err
}

// LBSUpdate ?????? WGS84 ??????
//
//	types ??? CP	?????????????????????, ???????????????????????????????????? all ??? female ???????????????
func (c *Client) LBSUpdate(openID string, types []string, lon, lat float64) error {
	if openID == "" {
		return ErrInvalidOpenID
	}
	if len(types) == 0 {
		return ErrInvalidType
	}

	return c.queryAndCheckResponse(apiLBSUpdate, &argLocation{
		OpenID:    openID,
		Types:     types,
		Longitude: lon,
		Latitude:  lat,
	}, nil)
}

// LBSDelete ?????? WGS84 ??????
func (c *Client) LBSDelete(openID string, types []string) error {
	if openID == "" {
		return ErrInvalidOpenID
	}
	if len(types) == 0 {
		return ErrInvalidType
	}

	return c.queryAndCheckResponse(apiLBSDelete, &argLocation{
		OpenID: openID,
		Types:  types,
	}, nil)
}

// LBSRadius ????????????????????????
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

	err := c.queryAndCheckResponse(apiLBSRadius, arg, resp)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

// Tracks ???????????????????????????
//
//	devicecode (????????????) ???????????????. ????????????????????????????????????
//	distinctID (?????????) ????????????. ??????????????? OpenID
//	opts: ??????????????????
func (c *Client) Tracks(
	devicecode, distinctID string, opts ...BigdataOptions) error {
	return c.producer.Tracks(devicecode, distinctID, opts...)
}

// track ?????????????????????????????????
func (c *Client) track(data []byte, logCount int, compress bool) (int, error) {
	if len(data) == 0 {
		return defaultStatus, nil
	}

	traceID, req := c.getRequest(true)
	ret := &response{}
	req.Header.Add(headerDataCount, Itoa(logCount))
	code, err := c.queryCode(apiBigDataTrack, req, config.TrackTimeout, data, ret, compress)
	if err != nil {
		return code, errWithTraceID(err, traceID)
	}
	err = c.checkResponse(ret)
	if err != nil {
		return code, errWithTraceID(err, traceID)
	}
	return code, nil
}

// CreateRank ???????????????
func (c *Client) CreateRank(rankID string, startTime, destroyTime time.Time) error {
	if rankID == "" {
		return ErrInvalidOpenID
	}

	err := c.queryAndCheckResponse(apiCreateRank, &rankAPIArg{
		RankID:      rankID,
		StartTime:   startTime.Format(time.RFC3339),
		DestroyTime: destroyTime.Format(time.RFC3339),
	}, nil)

	return err
}

// CloseRank ???????????????
func (c *Client) CloseRank(rankID string) error {
	if rankID == "" {
		return ErrInvalidOpenID
	}

	err := c.queryAndCheckResponse(apiCloseRank, &rankAPIArg{
		RankID: rankID,
	}, nil)

	return err
}

// RankAddScore ??????????????????
func (c *Client) RankAddScore(rankID string, openId string, score int64) error {
	if rankID == "" || openId == "" {
		return ErrInvalidOpenID
	}

	err := c.queryAndCheckResponse(apiRankAddScore, &rankAPIArg{
		RankID: rankID,
		OpenID: openId,
		Score:  score,
	}, nil)

	return err
}

// RankSetScore ??????????????????
func (c *Client) RankSetScore(rankID string, openId string, score int64) error {
	if rankID == "" || openId == "" {
		return ErrInvalidOpenID
	}

	err := c.queryAndCheckResponse(apiRankSetScore, &rankAPIArg{
		RankID: rankID,
		OpenID: openId,
		Score:  score,
	}, nil)

	return err
}

// QueryUserRank ????????????????????????
func (c *Client) QueryUserRank(rankID string, openId string) (*RankMember, error) {
	if rankID == "" || openId == "" {
		return nil, ErrInvalidOpenID
	}

	ret := &RankMember{}
	resp := &response{Data: ret}

	err := c.queryAndCheckResponse(apiQueryUserRank, &rankAPIArg{
		RankID: rankID,
		OpenID: openId,
	}, resp)

	return ret, err
}

// GetRankList ???????????????
func (c *Client) GetRankList(rankID string) ([]*RankMember, error) {
	if rankID == "" {
		return nil, ErrInvalidOpenID
	}

	var ret []*RankMember
	resp := &response{Data: &ret}

	err := c.queryAndCheckResponse(apiGetRankList, &rankAPIArg{
		RankID: rankID,
	}, resp)

	return ret, err
}

// GetFriendRankList ?????????????????????
func (c *Client) GetFriendRankList(rankID string, openId string) ([]*RankMember, error) {
	if rankID == "" || openId == "" {
		return nil, ErrInvalidOpenID
	}

	ret := make([]*RankMember, 0)
	resp := &response{Data: ret}

	err := c.queryAndCheckResponse(apiFriendsRank, &rankAPIArg{
		RankID: rankID,
		OpenID: openId,
	}, resp)

	return ret, err
}

func (c *Client) IMSLogin(req *IMSLoginReq) (*IMSLoginResp, error) {
	ret := &IMSLoginResp{}
	resp := &response{Data: ret}
	err := c.queryAndCheckResponse(apiIMSLogin, req, resp)
	return ret, err
}

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
	err := c.queryAndCheckResponse(apiIMSSendMessage, req, resp)
	return ret, err
}

func (c *Client) IMSGetHistory(req *IMSHistoryReq) (*IMSHistoryResp, error) {
	ret := &IMSHistoryResp{}
	resp := &response{Data: ret}
	err := c.queryAndCheckResponse(apiIMSGetHistory, req, resp)
	return ret, err
}

func (c *Client) IMSCreateConversation(req *IMSCreateConvReq) error {
	return c.queryAndCheckResponse(apiIMSCreateConversation, req, nil)
}

func (c *Client) IMSUpdateConversation(req *IMSUpdateConvReq) error {
	return c.queryAndCheckResponse(apiIMSUpdateConversation, req, nil)
}

func (c *Client) IMSDeleteConversation(req *IMSConvDeleteReq) error {
	return c.queryAndCheckResponse(apiIMSDeleteConversation, req, nil)
}

func (c *Client) IMSGetConversation(req *IMSGetConversationReq) (*IMSConversation, error) {
	ret := &IMSConversation{}
	resp := &response{Data: ret}
	err := c.queryAndCheckResponse(apiIMSGetConversation, req, resp)
	return ret, err
}

func (c *Client) IMSJoinConversation(req *IMSJoinConversationReq) error {
	return c.queryAndCheckResponse(apiIMSJoinConversation, req, nil)
}

func (c *Client) IMSLeaveConversation(req *IMSLeaveConversationReq) error {
	return c.queryAndCheckResponse(apiIMSLeaveConversation, req, nil)
}

func (c *Client) IMSUpdateConversationUserData(req *IMSUpdateConvUserDataReq) error {
	return c.queryAndCheckResponse(apiIMSUpdateConversationUserData, req, nil)
}

func (c *Client) IMSConversationUserList(req *IMSConversationUserListReq) ([]*IMSConversation, error) {
	ret := make([]*IMSConversation, 0)
	resp := &response{Data: &ret}
	err := c.queryAndCheckResponse(apiIMSConversationUserList, req, resp)
	return ret, err
}

// PusherPush ????????????
func (c *Client) PusherPush(req *PusherPushReq) error {

	return c.queryAndCheckResponse(apiPusherPush, req, nil)
}

func (c *Client) RiskGreenSyncScan(scenes []string, tasks []*GreenRequestTask, extend string) (*GreenUsercaseResult, error) {
	if len(scenes) <= 0 || len(tasks) <= 0 {
		return nil, ErrInvalidOpenID
	}
	ret := &GreenUsercaseResult{}
	resp := &response{Data: ret}
	err := c.queryAndCheckResponse(apiRiskGreenSyncScan, &GreenRequest{
		Scenes: scenes,
		Tasks:  tasks,
		Extend: extend,
	}, resp)
	return ret, err
}

func (c *Client) RiskGreenAsyncScan(scenes []string, tasks []*GreenRequestTask, extend string, callback string) (*GreenUsercaseResult, error) {
	if len(scenes) <= 0 || len(tasks) <= 0 {
		return nil, ErrInvalidOpenID
	}
	ret := &GreenUsercaseResult{}
	resp := &response{Data: ret}
	err := c.queryAndCheckResponse(apiRiskGreenAsyncScan, &GreenRequest{
		Scenes:     scenes,
		Tasks:      tasks,
		Extend:     extend,
		CPCallback: callback,
	}, resp)
	return ret, err
}

func (c *Client) RiskGreenGetScanRes(taskID []string) (*GreenUsercaseResult, error) {
	if len(taskID) <= 0 {
		return nil, ErrInvalidOpenID
	}
	ret := &GreenUsercaseResult{}
	resp := &response{Data: ret}
	err := c.queryAndCheckResponse(apiRiskGreenGetScanResult, &GreenRequest{
		TaskID: taskID,
	}, resp)
	return ret, err
}

func (c *Client) RiskGreenFeedback(taskID, url string, results map[string]string) error {
	if len(url) <= 0 || len(results) <= 0 {
		return ErrInvalidOpenID
	}
	err := c.queryAndCheckResponse(apiRiskGreenFeedback, &GreenFeedbackRequest{
		TaskID:  taskID,
		URL:     url,
		Results: results,
	}, nil)
	return err
}
