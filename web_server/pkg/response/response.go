package response

type Body struct {
    Code   int  `json:"code"`
    Msg    string  `json:"msg"`
    Status bool `json:"status"`
    Data   any  `json:"data"`
}

func Success(data any) Body {
    return Body{Code: 0, Msg: "ok", Status: true, Data: data}
}

func Error(code int, msg string) Body {
    return Body{Code: code, Msg: msg, Status: false, Data: nil}
}