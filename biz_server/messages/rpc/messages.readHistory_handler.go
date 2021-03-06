/*
 *  Copyright (c) 2018, https://github.com/airwide-code/airwide.datacenter
 *  All rights reserved.
 *
 *
 *
 */

package rpc

import (
	"github.com/golang/glog"
	"github.com/airwide-code/airwide.datacenter/baselib/logger"
	"github.com/airwide-code/airwide.datacenter/baselib/grpc_util"
	"github.com/airwide-code/airwide.datacenter/mtproto"
	"golang.org/x/net/context"
	"github.com/airwide-code/airwide.datacenter/biz/base"
	"github.com/airwide-code/airwide.datacenter/biz/dal/dao"
	"github.com/airwide-code/airwide.datacenter/biz_server/sync_client"
)

// messages.readHistory#e306d3a peer:InputPeer max_id:int = messages.AffectedMessages;
func (s *MessagesServiceImpl) MessagesReadHistory(ctx context.Context, request *mtproto.TLMessagesReadHistory) (*mtproto.Messages_AffectedMessages, error) {
	md := grpc_util.RpcMetadataFromIncoming(ctx)
	glog.Infof("messages.readHistory#e306d3a - metadata: %s, request: %s", logger.JsonDebugData(md), logger.JsonDebugData(request))

	peer := base.FromInputPeer(request.GetPeer())
	if peer.PeerType == base.PEER_SELF {
		// TODO(@benqi): 太土！
		peer.PeerType = base.PEER_USER
		peer.PeerId = md.UserId
	}

	// 消息已读逻辑
	// 1. inbox，设置unread_count为0以及read_inbox_max_id
	dao.GetUserDialogsDAO(dao.DB_MASTER).UpdateUnreadByPeer(request.GetMaxId(), md.UserId, int8(peer.PeerType), peer.PeerId)

	updateReadHistoryInbox := mtproto.NewTLUpdateReadHistoryInbox()
	updateReadHistoryInbox.SetPeer(peer.ToPeer())
	updateReadHistoryInbox.SetMaxId(request.MaxId)

	_, err := sync_client.GetSyncClient().SyncOneUpdateData2(md.ServerId, md.AuthId, md.SessionId, md.UserId, md.ClientMsgId, updateReadHistoryInbox.To_Update())
	if err != nil {
		return nil, err
	}

	//affected := mtproto.NewTLMessagesAffectedMessages()
	//affected.SetPts(int32(state.Pts))
	//affected.SetPtsCount(state.PtsCount)

	// 2. outbox, 设置read_outbox_max_id
	if peer.PeerType == base.PEER_USER {
		outboxDO := dao.GetUserDialogsDAO(dao.DB_SLAVE).SelectByPeer(peer.PeerId, int8(peer.PeerType), md.UserId)
		dao.GetUserDialogsDAO(dao.DB_MASTER).UpdateReadOutboxMaxIdByPeer(outboxDO.TopMessage, peer.PeerId, int8(peer.PeerType), md.UserId)

		updateReadHistoryOutbox := mtproto.NewTLUpdateReadHistoryOutbox()
		outboxPeer := &mtproto.TLPeerUser{Data2: &mtproto.Peer_Data{
			UserId: md.UserId,
		}}
		updateReadHistoryOutbox.SetPeer(outboxPeer.To_Peer())
		updateReadHistoryOutbox.SetMaxId(outboxDO.TopMessage)

		sync_client.GetSyncClient().PushToUserOneUpdateData(peer.PeerId, updateReadHistoryOutbox.To_Update())
	} else {
		doList := dao.GetChatParticipantsDAO(dao.DB_SLAVE).SelectByChatId(peer.PeerId)
		for _, do := range doList {
			if do.UserId == md.UserId {
				continue
			}
			outboxDO := dao.GetUserDialogsDAO(dao.DB_SLAVE).SelectByPeer(do.UserId, int8(peer.PeerType), peer.PeerId)
			dao.GetUserDialogsDAO(dao.DB_MASTER).UpdateReadOutboxMaxIdByPeer(outboxDO.TopMessage, do.UserId, int8(peer.PeerType), peer.PeerId)

			updateReadHistoryOutbox := mtproto.NewTLUpdateReadHistoryOutbox()
			outboxPeer := &mtproto.TLPeerChat{Data2: &mtproto.Peer_Data{
				ChatId: peer.PeerId,
			}}
			updateReadHistoryOutbox.SetPeer(outboxPeer.To_Peer())
			updateReadHistoryOutbox.SetMaxId(outboxDO.TopMessage)

			sync_client.GetSyncClient().PushToUserOneUpdateData(do.UserId, updateReadHistoryOutbox.To_Update())
		}
	}

	err = mtproto.NewRpcError2(mtproto.TLRpcErrorCodes_NOTRETURN_CLIENT)
	glog.Infof("messages.readHistory#e306d3a - reply: {%v}", err)
	return nil, err
	// affected.To_Messages_AffectedMessages(), nil
}
