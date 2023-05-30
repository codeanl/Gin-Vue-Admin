package response

// Response 基础序列化器
type Res struct {
	Code  int         `json:"code"`
	Data  interface{} `json:"data"`
	Msg   string      `json:"msg"`
	Error string      `json:"error"`
}

// TokenData 带有token的Data结构
type TokenData struct {
	User          interface{} `json:"user"`
	Authorization string      `json:"Authorization"`
}
