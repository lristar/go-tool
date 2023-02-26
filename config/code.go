package config

type Code int

//周边错误码
const (
	SUCCESS            Code = 200    //成功
	INVALID            Code = 101001 //缺少必填入参，或入参格式错误
	DATA_ERROR         Code = 101002 //提交日期格式错误
	NOT_EXIST          Code = 101003 //客户账户不存在
	HS_ERROR           Code = 101004 //hs接口错误
	F_NOT_EXIST        Code = 101005 //交易账户不存在
	PROCESS_ERROR      Code = 101006 //流程处理错误
	CURRENCY           Code = 101008 //通用错误码
	C_MODIFY_INVALID   Code = 302001 //入参错误
	C_MODIFY_DB        Code = 302002 //DB错误
	C_CLIENT_NOT_EXIST Code = 302003 //客户不存在
	C_PROCESS_LOCK     Code = 302004 //流程被锁定
	C_PROCESS_END      Code = 302005 //流程已结束
	C_MODIFY_ERROR     Code = 302008 //通用错误
	F_STORE_INVALID    Code = 501001
	F_MT_NOT_EXIST     Code = 501002 //币种类型不存在
	F_STORE_DB         Code = 501003 //资金存DB错误
	F_TAKE_INVALID     Code = 502001 //资金取入参错误
	F_TAKE_DB          Code = 502003 //资金取DB错误
	F_TRANSFER_INVALID Code = 503001 //资金调拨入参错误
	F_TRANSFER_DB      Code = 503003 //资金调拨DB错误
	C_STATUS_ERROR     Code = 503005 //客户状态异常
)
