package translateShell

import (
	"Multimedia_Processing_Pipeline/constant"
	"log"
	"os/exec"
	"strings"
	"sync"
)

func TransByBing(proxy, language, src string, ans chan Result, c *constant.Count, once *sync.Once) {
	cmd := exec.Command("trans", "-brief", "-engine", "bing", "-proxy", proxy, language, src)
	output, err := cmd.CombinedOutput()
	if err != nil || strings.Contains(string(output), "u001b") || strings.Contains(string(output), "Didyoumean") || strings.Contains(string(output), "Connectiontimedout") {
		log.Printf("bing查询命令执行出错\t命令原文:%v\t错误原文:%v\n", cmd.String(), err.Error())
		return
	}
	r := Result{
		From: "Bing",
		Dst:  string(output),
	}
	once.Do(func() {
		ans <- r
		c.SetBing()
	})
}
