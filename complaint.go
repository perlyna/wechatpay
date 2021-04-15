package wechatpay

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/perlyna/wechatpay/core"
	"github.com/perlyna/wechatpay/model"
	"github.com/perlyna/wechatpay/util"
)

// ParseComplaintNotify 解析投诉通知回调数据
// 文档链接: https://pay.weixin.qq.com/wiki/doc/apiv3/apis/chapter10_2_16.shtml
func ParseComplaintNotify(body []byte, apiv3Secret string) (model.ComplaintEvent, error) {
	ret := model.ComplaintEvent{}
	if err := json.Unmarshal(body, &ret); err != nil {
		return ret, err
	}
	rawComplaint, err := util.DecryptToByte(apiv3Secret, ret.Resource.AssociatedData,
		ret.Resource.Nonce, ret.Resource.Ciphertext)
	if err != nil {
		return ret, err
	}
	err = json.Unmarshal(rawComplaint, &ret)
	return ret, err
}

const complaintsURL = "https://api.mch.weixin.qq.com/v3/merchant-service/complaints-v2"

// ListComplaints 查询投诉单列表
// 文档链接: https://pay.weixin.qq.com/wiki/doc/apiv3/apis/chapter10_2_11.shtml
// 最新更新时间：2021.04.01
func (p *WechatPay) ListComplaints(ctx context.Context, begin, end time.Time) ([]model.Complaint, error) {
	v := url.Values{}
	v.Set("begin_date", begin.Format("2006-01-02"))
	v.Set("end_date", end.Format("2006-01-02"))
	v.Set("limit", "50")
	complaints := []model.Complaint{}
	for offset := 0; ; offset++ {
		v.Set("offset", strconv.Itoa(offset))
		body, err := core.Get(ctx, p.Client, p.credential, p.validator, complaintsURL+"?"+v.Encode())
		if err != nil {
			return nil, err
		}
		reply := model.ComplaintReply{}
		if err = json.Unmarshal(body, &reply); err != nil {
			return complaints, err
		}
		for _, complaint := range reply.Complaints {
			if complaint.PayerPhone != "" {
				if plaintext, err := util.DecryptOAEP(complaint.PayerPhone, p.privateKey); err == nil {
					complaint.PayerPhone = plaintext
				}
			}
			complaints = append(complaints, complaint)
		}
		if reply.TotalCount == 0 ||
			len(complaints) >= int(reply.TotalCount) ||
			len(reply.Complaints) == 0 {
			break
		}
	}
	return complaints, nil
}

// GetComplaint 查询投诉详情
// 文档链接: https://pay.weixin.qq.com/wiki/doc/apiv3/apis/chapter10_2_13.shtml
// 最新更新时间：2021.04.01
func (p *WechatPay) GetComplaint(ctx context.Context, complaintID string) (complaint model.Complaint, err error) {
	reqURL := "https://api.mch.weixin.qq.com/v3/merchant-service/complaints-v2/" + complaintID
	body, err := core.Get(ctx, p.Client, p.credential, p.validator, reqURL)
	if err != nil {
		return complaint, err
	}

	err = json.Unmarshal(body, &complaint)
	if complaint.PayerPhone != "" {
		if plaintext, err := util.DecryptOAEP(complaint.PayerPhone, p.privateKey); err == nil {
			complaint.PayerPhone = plaintext
		}
	}
	return complaint, err
}

// complaintNotifyURL 投诉通知回调地址API
const complaintNotifyURL = "https://api.mch.weixin.qq.com/v3/merchant-service/complaint-notifications"

// complaintNotifyReq 投诉通知回调地址请求参数
type complaintNotifyReq struct {
	MchID string `json:"mchid,omitempty"`
	URL   string `json:"url,omitempty"`
}

// ComplaintNotifications 创建投诉通知回调地址
// 文档链接: https://pay.weixin.qq.com/wiki/doc/apiv3/apis/chapter10_2_2.shtml
func (p *WechatPay) CreateComplaintNotification(ctx context.Context, notifyURL string) (err error) {
	reqBody := complaintNotifyReq{URL: notifyURL}
	_, err = core.Post(ctx, p.Client, p.credential, p.validator, complaintNotifyURL, reqBody)
	return err
}

// GetComplaintNotification 查询投诉通知回调地址
// 文档链接: https://pay.weixin.qq.com/wiki/doc/apiv3/apis/chapter10_2_3.shtml
func (p *WechatPay) GetComplaintNotification(ctx context.Context) (notifyURL string, err error) {
	body, err := core.Get(ctx, p.Client, p.credential, p.validator, complaintNotifyURL)
	if err != nil {
		return "", err
	}
	reply := complaintNotifyReq{}
	err = json.Unmarshal(body, &reply)
	return reply.URL, err
}

// ComplaintNotifications 创建投诉通知回调地址
// 文档链接: https://pay.weixin.qq.com/wiki/doc/apiv3/apis/chapter10_2_4.shtml
func (p *WechatPay) UpdateComplaintNotification(ctx context.Context, notifyURL string) (err error) {
	body, err := core.Put(ctx, p.Client, p.credential, p.validator, complaintNotifyURL, complaintNotifyReq{URL: notifyURL})
	if err != nil {
		return err
	}
	reply := complaintNotifyReq{}
	if err = json.Unmarshal(body, &reply); err != nil {
		return err
	}
	if reply.URL != notifyURL {
		return fmt.Errorf("投诉通知回调地址没有更新成功")
	}
	return nil
}

// GetComplaintNotification 查询投诉通知回调地址
// 文档链接: https://pay.weixin.qq.com/wiki/doc/apiv3/apis/chapter10_2_5.shtml
func (p *WechatPay) DeleteComplaintNotification(ctx context.Context) (err error) {
	_, err = core.Delete(ctx, p.Client, p.credential, p.validator, complaintNotifyURL, nil)
	return err
}
