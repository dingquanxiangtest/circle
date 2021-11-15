package restful

import (
	"git.internal.yunify.com/qxp/misc/logger"
	"git.internal.yunify.com/qxp/misc/resp"
	"git.internal.yunify.com/qxp/molecule/internal/eval"
	"git.internal.yunify.com/qxp/molecule/internal/service"
	"git.internal.yunify.com/qxp/molecule/pkg/misc/config"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Formula formula
type Formula struct {
	name string
}

// NewFormula 初始化
func NewFormula(conf *config.Config, opt ...service.Options) (*Formula, error) {
	return &Formula{
		name: "table formula",
	}, nil
}

// FormulaReq FormulaReq
type FormulaReq struct {
	Expression string                   `json:"expression"`
	Parameter map[string]interface{}    `json:"parameter"`
}

// Calculation Calculation
func (f *Formula) Calculation(c *gin.Context) {
	req := &FormulaReq{}
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(eval.Handler(logger.CTXTransfer(c), req.Expression,req.Parameter)).Context(c)
}