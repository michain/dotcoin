package packet

type JsonRequest struct{
	Version string
	Command string
	Message interface{}
}

type JsonResult struct{
	RetCode int
	RetMsg string
	Message interface{}
}



