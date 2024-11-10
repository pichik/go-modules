package request

import (
	"sync"

	"github.com/pichik/go-modules/tool"
)

type Request struct {
	UtilData *tool.UtilData
}

type RequestFlow struct {
	UtilData *tool.UtilData
}

type Repeater struct {
	UtilData *tool.UtilData
}

type IUtil interface {
	SetupFlags() []tool.UtilData
	SetupData()
}

type ITool interface {
	SetupFlags()
	SetupInput(urls []string)
}

type IFlowTool interface {
	Results(requestData RequestData, m *sync.Mutex)
}
