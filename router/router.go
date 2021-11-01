package router

import (
	"log"
	"wikipay-admin/apis/card"
	"wikipay-admin/apis/exchange"
	"wikipay-admin/apis/financial"
	"wikipay-admin/apis/home"
	"wikipay-admin/apis/im"
	log2 "wikipay-admin/apis/log"
	"wikipay-admin/apis/mch"
	"wikipay-admin/apis/monitor"
	"wikipay-admin/apis/operate"
	"wikipay-admin/apis/system"
	"wikipay-admin/apis/system/dict"
	. "wikipay-admin/apis/tools"
	"wikipay-admin/apis/topupcard"
	"wikipay-admin/apis/user"
	"wikipay-admin/apis/wallet"

	_ "wikipay-admin/docs"
	"wikipay-admin/handler"
	"wikipay-admin/handler/sd"
	"wikipay-admin/middleware"
	_ "wikipay-admin/pkg/jwtauth"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(middleware.LoggerToFile())
	r.Use(middleware.CustomError)
	r.Use(middleware.NoCache)
	r.Use(middleware.Options)
	r.Use(middleware.Secure)
	r.Use(middleware.RequestId())

	r.GET("/", system.HelloWorld)
	r.Static("/static", "./static")
	r.GET("/info", handler.Ping)

	// 监控信息
	svcd := r.Group("/sd")
	{
		svcd.GET("/health", sd.HealthCheck)
		svcd.GET("/disk", sd.DiskCheck)
		svcd.GET("/cpu", sd.CPUCheck)
		svcd.GET("/ram", sd.RAMCheck)
		svcd.GET("/os", sd.OSCheck)
	}

	// the jwt middleware
	authMiddleware, err := middleware.AuthInit()

	if err != nil {
		log.Fatalln("JWT Error", err.Error())
	}

	r.POST("/login", authMiddleware.LoginHandler)

	// Refresh time can be longer than token timeout
	r.GET("/refresh_token", authMiddleware.RefreshHandler)
	r.GET("/routes", Dashboard)

	apiv1 := r.Group("/api/v1")
	{

		apiv1.GET("/monitor/server", monitor.ServerInfo)

		apiv1.GET("/getCaptcha", system.GenerateCaptchaHandler)
		apiv1.GET("/db/tables/page", GetDBTableList)
		apiv1.GET("/db/columns/page", GetDBColumnList)
		apiv1.GET("/sys/tables/page", GetSysTableList)
		apiv1.POST("/sys/tables/info", InsertSysTable)
		apiv1.PUT("/sys/tables/info", UpdateSysTable)
		apiv1.DELETE("/sys/tables/info/:tableId", DeleteSysTables)
		apiv1.GET("/sys/tables/info/:tableId", GetSysTables)
		apiv1.GET("/gen/preview/:tableId", Preview)
		apiv1.GET("/menuTreeselect", system.GetMenuTreeelect)
		apiv1.GET("/dict/databytype/:dictType", dict.GetDictDataByDictType)
		//通用上传
		apiv1.POST("/upload", UploadFile)

	}

	auth := r.Group("/api/v1")
	auth.Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{

		auth.GET("/deptList", system.GetDeptList)
		auth.GET("/deptTree", system.GetDeptTree)
		auth.GET("/dept/:deptId", system.GetDept)
		auth.POST("/dept", system.InsertDept)
		auth.PUT("/dept", system.UpdateDept)
		auth.DELETE("/dept/:id", system.DeleteDept)

		auth.GET("/dict/datalist", dict.GetDictDataList)
		auth.GET("/dict/data/:dictCode", dict.GetDictData)
		auth.POST("/dict/data", dict.InsertDictData)
		auth.PUT("/dict/data/", dict.UpdateDictData)
		auth.DELETE("/dict/data/:dictCode", dict.DeleteDictData)

		auth.GET("/dict/typelist", dict.GetDictTypeList)
		auth.GET("/dict/type/:dictId", dict.GetDictType)
		auth.POST("/dict/type", dict.InsertDictType)
		auth.PUT("/dict/type", dict.UpdateDictType)
		auth.DELETE("/dict/type/:dictId", dict.DeleteDictType)

		auth.GET("/dict/typeoptionselect", dict.GetDictTypeOptionSelect)

		auth.GET("/sysUserList", system.GetSysUserList)
		auth.GET("/sysUser/:userId", system.GetSysUser)
		auth.GET("/sysUser/", system.GetSysUserInit)
		auth.POST("/sysUser", system.InsertSysUser)
		auth.PUT("/sysUser", system.UpdateSysUser)
		auth.DELETE("/sysUser/:userId", system.DeleteSysUser)

		auth.GET("/rolelist", system.GetRoleList)
		auth.GET("/role/:roleId", system.GetRole)
		auth.POST("/role", system.InsertRole)
		auth.PUT("/role", system.UpdateRole)
		auth.PUT("/roledatascope", system.UpdateRoleDataScope)
		auth.DELETE("/role/:roleId", system.DeleteRole)

		auth.GET("/configList", system.GetConfigList)
		auth.GET("/config/:configId", system.GetConfig)
		auth.POST("/config", system.InsertConfig)
		auth.PUT("/config", system.UpdateConfig)
		auth.DELETE("/config/:configId", system.DeleteConfig)

		auth.GET("/roleMenuTreeselect/:roleId", system.GetMenuTreeRoleselect)
		auth.GET("/roleDeptTreeselect/:roleId", system.GetDeptTreeRoleselect)

		auth.GET("/getinfo", system.GetInfo)
		auth.GET("/user/profile", system.GetSysUserProfile)
		auth.POST("/user/avatar", system.InsetSysUserAvatar)
		auth.PUT("/user/pwd", system.SysUserUpdatePwd)

		auth.GET("/postlist", system.GetPostList)
		auth.GET("/post/:postId", system.GetPost)
		auth.POST("/post", system.InsertPost)
		auth.PUT("/post", system.UpdatePost)
		auth.DELETE("/post/:postId", system.DeletePost)

		auth.GET("/menulist", system.GetMenuList)
		auth.GET("/menu/:id", system.GetMenu)
		auth.POST("/menu", system.InsertMenu)
		auth.PUT("/menu", system.UpdateMenu)
		auth.DELETE("/menu/:id", system.DeleteMenu)
		auth.GET("/menurole", system.GetMenuRole)

		auth.GET("/menuids", system.GetMenuIDS)

		auth.GET("/loginloglist", log2.GetLoginLogList)
		auth.GET("/loginlog/:infoId", log2.GetLoginLog)
		auth.POST("/loginlog", log2.InsertLoginLog)
		auth.PUT("/loginlog", log2.UpdateLoginLog)
		auth.DELETE("/loginlog/:infoId", log2.DeleteLoginLog)

		auth.GET("/operloglist", log2.GetOperLogList)
		auth.GET("/operlog/:operId", log2.GetOperLog)
		auth.DELETE("/operlog/:operId", log2.DeleteOperLog)

		auth.GET("/configKey/:configKey", system.GetConfigByConfigKey)

		auth.POST("/logout", handler.LogOut)

		///广告
		auth.GET("/user/advertList", user.GetAdvertList)
		auth.GET("/user/advert/:id", user.GetAdvert)
		auth.POST("/user/advert", user.InsertAdvert)
		auth.PUT("/user/advert", user.UpdateAdvert)
		auth.DELETE("/user/advert/:id", user.DeleteAdvert)

		//
		auth.GET("/user/contentList", user.GetContentList)
		auth.GET("/user/content/:id", user.GetContent)
		auth.POST("/user/content", user.InsertContent)
		auth.PUT("/user/content", user.UpdateContent)
		auth.DELETE("/user/content/:id", user.DeleteContent)

		//问题反馈
		auth.GET("/user/feedbackTypeList", user.GetFeedbackTypeList)
		auth.GET("/user/feedbackType/:code", user.GetFeedbackType)
		auth.POST("/user/feedbackType", user.InsertFeedbackType)
		auth.PUT("/user/feedbackType", user.UpdateFeedbackType)
		auth.DELETE("/user/feedbackType/:code", user.DeleteFeedbackType)

		//问题类型反馈
		auth.GET("/user/feedbackList", user.GetFeedbackList)
		auth.GET("/user/feedback/:id", user.GetFeedback)
		auth.POST("/user/feedback", user.InsertFeedback)
		auth.PUT("/user/feedback", user.UpdateFeedback)
		auth.DELETE("/user/feedback/:id", user.DeleteFeedback)
		//错误码
		auth.GET("/user/errorList", user.GetErrorList)
		auth.GET("/user/error/:code", user.GetError)
		auth.POST("/user/error", user.InsertError)
		auth.PUT("/user/error", user.UpdateError)
		auth.DELETE("/user/error/:code", user.DeleteError)

		//错误码
		auth.GET("/user/userList", user.GetUserList)
		auth.GET("/user/user/:id", user.GetUser)
		auth.POST("/user/user", user.InsertUser)
		auth.PUT("/user/user", user.UpdateUser)
		auth.DELETE("/user/user/:id", user.DeleteUser)

		//设置
		auth.GET("/user/settingList", user.GetSettingList)
		auth.GET("/user/setting/:id", user.GetSetting)
		auth.POST("/user/setting", user.InsertSetting)
		auth.PUT("/user/setting", user.UpdateSetting)
		auth.DELETE("/user/setting/:id", user.DeleteSetting)

		//用户风控
		auth.GET("/user/controlList", user.GetUserControlInfos)
		auth.POST("/user/status", user.UpdateUserStatus)
		auth.POST("/user/reset", user.ResetUser)
		//版本
		auth.GET("/user/versionList", user.GetVersionList)
		auth.GET("/user/version/:id", user.GetVersion)
		auth.POST("/user/version", user.InsertVersion)
		auth.PUT("/user/version", user.UpdateVersion)
		auth.DELETE("/user/version/:id", user.DeleteVersion)

		//登录历史
		auth.GET("/user/loginHistoriesList", user.GetLoginHistoriesList)
		auth.GET("/user/loginHistories/:id", user.GetLoginHistories)
		auth.POST("/user/loginHistories", user.InsertLoginHistories)
		auth.PUT("/user/loginHistories", user.UpdateLoginHistories)
		auth.DELETE("/user/loginHistories/:id", user.DeleteVersion)

		//资金账户
		auth.GET("/wallet/accountList", wallet.GetAccountList)
		auth.GET("/wallet/account/:id", wallet.GetAccount)
		auth.POST("/wallet/account", wallet.InsertAccount)
		auth.PUT("/wallet/account", wallet.UpdateAccount)
		auth.DELETE("/wallet/account/:id", wallet.DeleteAccount)

		//地址列表
		auth.GET("/wallet/addressList", wallet.GetAddressList)
		auth.GET("/wallet/address/:id", wallet.GetAddress)
		auth.POST("/wallet/address", wallet.InsertAddress)
		auth.PUT("/wallet/address", wallet.UpdateAddress)
		auth.DELETE("/wallet/address/:id", wallet.DeleteAddress)

		//汇率
		auth.GET("/wallet/interestList", wallet.GetInterestList)
		auth.GET("/wallet/interest/:id", wallet.GetInterest)
		auth.POST("/wallet/interest", wallet.InsertInterest)
		auth.PUT("/wallet/interest", wallet.UpdateInterest)
		auth.DELETE("/wallet/interest/:id", wallet.DeleteInterest)

		//向表
		auth.GET("/wallet/coinList", wallet.GetCoinList)
		auth.GET("/wallet/coin/:id", wallet.GetCoin)
		auth.POST("/wallet/coin", wallet.InsertCoin)
		auth.PUT("/wallet/coin", wallet.UpdateCoin)
		auth.DELETE("/wallet/coin/:id", wallet.DeleteCoin)

		//转账
		auth.GET("/wallet/transferList", wallet.GetTransferList)
		auth.GET("/wallet/transfer/:id", wallet.GetTransfer)
		auth.POST("/wallet/transfer", wallet.InsertTransfer)
		auth.PUT("/wallet/transfer", wallet.UpdateTransfer)
		auth.DELETE("/wallet/transfer/:id", wallet.DeleteTransfer)

		//预处理转账
		auth.GET("/wallet/transferTempList", wallet.GetTransferTempList)
		auth.GET("/wallet/transferTemp/:id", wallet.GetTransferTemp)
		auth.POST("/wallet/transferTemp", wallet.InsertTransferTemp)
		auth.PUT("/wallet/transferTemp", wallet.UpdateTransferTemp)
		auth.DELETE("/wallet/transferTemp/:id", wallet.DeleteTransferTemp)

		//账单
		auth.GET("/wallet/transactionList", wallet.GetTransactionList)
		auth.GET("/wallet/transaction/:id", wallet.GetTransaction)
		auth.POST("/wallet/transaction", wallet.InsertTransaction)
		auth.PUT("/wallet/transaction", wallet.UpdateTransaction)
		auth.DELETE("/wallet/transaction/:id", wallet.DeleteTransaction)

		//充提现记录
		auth.GET("/wallet/withdrawDepositList", wallet.GetWithdrawDepositList)
		auth.GET("/wallet/withdrawDeposit/:id", wallet.GetWithdrawDeposit)
		auth.POST("/wallet/withdrawDeposit", wallet.InsertWithdrawDeposit)
		auth.PUT("/wallet/withdrawDeposit", wallet.UpdateWithdrawDeposit)
		auth.PUT("/wallet/withdrawDeposit/status", wallet.UpdateWithdrawDepositStatus)
		auth.DELETE("/wallet/withdrawDeposit/:id", wallet.DeleteWithdrawDeposit)

		//提现预付订单
		auth.GET("/wallet/withdrawtempList", wallet.GetWithdrawTempList)
		auth.GET("/wallet/withdrawtemp/:id", wallet.GetWithdrawTemp)
		auth.POST("/wallet/withdrawtemp", wallet.InsertWithdrawTemp)
		auth.PUT("/wallet/withdrawtemp/status", wallet.UpdateWithdrawTemp)
		auth.DELETE("/wallet/withdrawtemp/:id", wallet.DeleteWithdrawTemp)

		//收益
		//充提现记录
		auth.GET("/wallet/profitList", wallet.GetProfitList)
		auth.GET("/wallet/profit/:id", wallet.GetProfit)
		auth.POST("/wallet/profit", wallet.InsertProfit)
		auth.PUT("/wallet/profit", wallet.UpdateProfit)
		auth.DELETE("/wallet/profit/:id", wallet.DeleteProfit)

		//兑换汇率
		auth.GET("/ex/coinpriceList", exchange.GetCoinPriceList)
		auth.GET("/exchange/coinprice/:id", exchange.GetCoinPrice)
		auth.POST("/exchange/coinprice", exchange.InsertCoinPrice)
		auth.PUT("/exchange/coinprice", exchange.UpdateCoinPrice)
		auth.DELETE("/exchange/coinprice/:id", exchange.DeleteCoinPrice)

		//兑换
		auth.GET("/ex/coinbidaskList", exchange.GetCoinBidAskList)
		auth.GET("/exchange/coinbidask/:id", exchange.GetCoinBidAsk)
		auth.POST("/exchange/coinbidask", exchange.InsertCoinBidAsk)
		auth.PUT("/exchange/coinbidask", exchange.UpdateCoinBidAsk)
		auth.DELETE("/exchange/coinbidask/:id", exchange.DeleteCoinBidAsk)

		//卡
		auth.GET("/card/cardList", card.GetCardList)
		auth.GET("/card/card/:id", card.GetCard)
		auth.POST("/card/card", card.InsertCard)
		auth.PUT("/card/card", card.UpdateCard)
		auth.DELETE("/card/card/:id", card.DeleteCard)

		//地址监控
		auth.GET("/monitor/addressList", monitor.GetMonitorAddressList)
		auth.GET("/monitor/address/:id", monitor.GetMonitorAddress)
		auth.POST("/monitor/address", monitor.InsertMonitorAddress)
		auth.PUT("/monitor/address", monitor.UpdateMonitorAddress)
		auth.DELETE("/monitor/address/:id", monitor.DeleteMonitorAddress)
		//同步
		auth.POST("/monitor/address/sync", monitor.SyncMonitorAddress)

		//获取余额
		auth.GET("/monitor/balance", monitor.GetBalance)
		auth.GET("/monitor/networkinfo", monitor.GetNetworkInfo)
		auth.GET("/monitor/balance/satistical", monitor.GetSatistical)
		auth.GET("/monitor/balance/assetComparison", monitor.AssetComparison)

		//用户地址监控
		auth.GET("/monitor/user/addressList", monitor.GetMonitorUserAddressList)
		auth.POST("/monitor/user/balances/sync", monitor.SyncAccountAndBalance)

		//归集设置
		auth.POST("/monitor/setting", monitor.InsertMonitorSetting)
		auth.PUT("/monitor/setting", monitor.UpdateMonitorSetting)
		auth.GET("/monitor/setting", monitor.GetMonitorSetting)

		//获取当日充值、提现统计
		auth.GET("/monitor/withdraw", monitor.GetWithdrawWithNow)

		//更新归集比率
		auth.POST("/monitor/ratio/update", monitor.UpdateMonitorRatio)
		auth.GET("/monitor/historyList", monitor.GetMonitorHistoryList)
		auth.GET("/monitor/history/detail", monitor.GetMonitorHistoryDetail)

		auth.GET("/monitor/gasprice", monitor.GetGasPrice)
		auth.GET("/monitor/assets", monitor.GetAssetsGroup)
		auth.GET("/monitor/collect/assets", monitor.GetCollectAssetsGroup)
		//auth.GET("/monitor/assets/total", monitor.GetAssetsTotal)                //获取gasprice
		auth.POST("/monitor/setting/manual", monitor.InsertMonitorManualSetting) //主动归集设置

		//卡交易
		auth.GET("/card/cardtransferList", card.GetCardTransferList)
		auth.GET("/card/cardtransfer/:id", card.GetCardTransfer)
		auth.POST("/card/cardtransfer", card.InsertCardTransfer)
		auth.PUT("/card/cardtransfer", card.UpdateCardTransfer)
		auth.DELETE("/card/cardtransfer/:id", card.DeleteCardTransfer)
		auth.POST("/card/cardtransfer/audit", card.AuditCardTransfer)
		//联系人
		auth.GET("/im/contactList", im.GetContactList)
		auth.GET("/im/contact/:id", im.GetContact)
		auth.POST("/im/contact", im.InsertContact)
		auth.PUT("/im/contact", im.UpdateContact)
		auth.DELETE("/im/contact/:id", im.DeleteContact)

		//系统消息
		auth.GET("/im/messageList", im.GetMessageList)
		auth.GET("/im/message/:id", im.GetMessage)
		auth.POST("/im/message", im.InsertMessage)
		auth.PUT("/im/message", im.UpdateMessage)
		auth.DELETE("/im/message/:id", im.DeleteMessage)

		//联系历史
		auth.GET("/im/contacthistoryList", im.GetContactHistoryList)
		auth.GET("/im/contacthistory/:id", im.GetContactHistory)
		auth.POST("/im/contacthistory", im.InsertContactHistory)
		auth.PUT("/im/contacthistory", im.UpdateContactHistory)
		auth.DELETE("/im/contacthistory/:id", im.DeleteContactHistory)

		//ab卡信息
		auth.GET("/topupcard/infoList", topupcard.GetTopupCardList)
		auth.GET("/topupcard/info/:id", topupcard.GetTopupCard)
		auth.POST("/topupcard/info", topupcard.InsertTopupCard)
		auth.PUT("/topupcard/info", topupcard.UpdateTopupCard)
		auth.DELETE("/topupcard/info/:id", topupcard.DeleteTopupCard)

		//ab卡交易
		auth.GET("/topupcard/tradeList", topupcard.GetTopupCardTradeList)
		auth.GET("/topupcard/trade/:id", topupcard.GetTopupCardTrade)
		auth.POST("/topupcard/trade", topupcard.InsertTopupCardTrade)
		auth.PUT("/topupcard/trade", topupcard.UpdateTopupCardTrade)
		auth.DELETE("/topupcard/trade/:id", topupcard.DeleteTopupCardTrade)

		//商户信息
		auth.GET("/mch/mchinfoList", mch.GetMchInfoList)
		auth.GET("/mch/mchinfo/:id", mch.GetMchInfo)
		auth.POST("/mch/mchinfo", mch.InsertMchInfo)
		auth.PUT("/mch/mchinfo", mch.UpdateMchInfo)
		auth.DELETE("/mch/mchinfo/:id", mch.DeleteMchInfo)

		//密钥
		auth.GET("/mch/mchkeyList", mch.GetMchKeyList)
		auth.GET("/mch/mchkey/:id", mch.GetMchKey)
		auth.POST("/mch/mchkey", mch.InsertMchKey)
		auth.PUT("/mch/mchkey", mch.UpdateMchKey)
		auth.DELETE("/mch/mchkey/:id", mch.DeleteMchKey)

		//流水订单
		auth.GET("/mch/mchorderList", mch.GetMchOrderList)
		auth.GET("/mch/mchorder/:id", mch.GetMchOrder)
		auth.POST("/mch/mchorder", mch.InsertMchOrder)
		auth.PUT("/mch/mchorder", mch.UpdateMchOrder)
		auth.DELETE("/mch/mchorder/:id", mch.DeleteMchOrder)

		//流水交易
		auth.GET("/mch/mchtransactionList", mch.GetMchTransactionList)
		auth.GET("/mch/mchtransaction/:id", mch.GetMchTransaction)
		auth.POST("/mch/mchtransaction", mch.InsertMchTransaction)
		auth.PUT("/mch/mchtransaction", mch.UpdateMchTransaction)
		auth.DELETE("/mch/mchtransaction/:id", mch.DeleteMchTransaction)

		//首页
		auth.GET("/home/total/info", home.GetTotalInfo)
		auth.GET("/home/profit/info", home.GetProfitInfo)
		auth.GET("/home/trade/info", home.GetTradeInfo)
		auth.GET("/home/financial/info", home.GetFinancialInfo)

		//提现审核记录
		auth.GET("/financial/withdraw/audit/info", financial.GetWithdrawAuditInfo)
		auth.GET("/financial/info", financial.GetFinancialInfoWithDateTime)
		auth.GET("/financial/transaction/info", financial.GetTransactionFinancialInfo)
		auth.POST("/financial/withdraw/audit", financial.WithdrawAudit)
		auth.GET("/financial/bidask/list", financial.ListBidAsk)
		auth.GET("/financial/fee/list", financial.ListFee)

		//
		auth.GET("/identity/list", user.GetIdentityList)
		auth.POST("/identity/audit", user.AuditIdentity)
		auth.GET("/identity/info/:id", user.GetIdentityFromId)

		//运营管理
		auth.GET("/operate/list", operate.GetOperateList)
		auth.POST("/operate", operate.InsertOperate)
		auth.PUT("/operate", operate.UpdateOperate)
		auth.GET("/operate/info/:id", operate.GetOperate)
		auth.DELETE("/operate/:id", operate.DeleteOperate)
	}
	//r.NoRoute(authMiddleware.MiddlewareFunc(), NoFound)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	log.Println("路由加载成功！")
	return r
}

//
func Dashboard(c *gin.Context) {
	var user = make(map[string]interface{})
	user["login_name"] = "admin"
	user["user_id"] = 1
	user["user_name"] = "管理员"
	user["dept_id"] = 1

	var cmenuList = make(map[string]interface{})
	cmenuList["children"] = nil
	cmenuList["parent_id"] = 1
	cmenuList["title"] = "用户管理"
	cmenuList["name"] = "Sysuser"
	cmenuList["icon"] = "user"
	cmenuList["order_num"] = 1
	cmenuList["id"] = 4
	cmenuList["path"] = "sysuser"
	cmenuList["component"] = "sysuser/index"

	var lista = make([]interface{}, 1)
	lista[0] = cmenuList

	var menuList = make(map[string]interface{})
	menuList["children"] = lista
	menuList["parent_id"] = 1
	menuList["name"] = "Upms"
	menuList["title"] = "权限管理"
	menuList["icon"] = "example"
	menuList["order_num"] = 1
	menuList["id"] = 4
	menuList["path"] = "/upms"
	menuList["component"] = "Layout"

	var list = make([]interface{}, 1)
	list[0] = menuList
	var data = make(map[string]interface{})
	data["user"] = user
	data["menuList"] = list

	var r = make(map[string]interface{})
	r["code"] = 200
	r["data"] = data

	c.JSON(200, r)
}
