package clog

import (
	"fmt"
	"log"
	"net"
	"reflect"
	"strings"

	"github.com/paulraysmile/utils"
)

var (
	Level     int
	Mode      int
	cate_dbg  string
	cate_war  string
	cate_err  string
	cate_info string
	cate_busi string
)

// 请赋值成自己的获取master addr的函数
var AddrFunc = func() (string, error) {
	return "127.0.0.1:28702", nil
}

// 请赋值成自己的split content的函数
var SplitFunc = func(content string) (contents []string) {
	for bpos, epos, seg, l := 0, 0, 65000, len(content); bpos < l; bpos += seg {
		epos = bpos + seg
		if epos > l {
			epos = l
		}
		contents = append(contents, content[bpos:epos])
	}

	return
}

func sendAgent(tube, content string) {
	addr, err := AddrFunc()
	if err != nil {
		return
	}
	conn, err := net.Dial("udp", addr)
	if err != nil {
		return
	}
	defer conn.Close()

	for _, c := range SplitFunc(content) {
		out := tube + "," + c
		conn.Write([]byte(out))
	}
}

func iprint(params []interface{}) {
	for pos, param := range params {
		switch param.(type) {
		case nil, string, []byte:
			continue
		}

		typ := reflect.TypeOf(param)
		if kind := typ.Kind(); kind <= reflect.Complex128 {
			continue
		}

		v := reflect.ValueOf(param)
		switch typ.Kind() {
		case reflect.Ptr, reflect.Map, reflect.Slice:
			if v.IsNil() {
				continue
			}
		}

		if typ.Implements(reflect.TypeOf((*error)(nil)).Elem()) || typ.Implements(reflect.TypeOf((*fmt.Stringer)(nil)).Elem()) {
			params[pos] = fmt.Sprintf("%v", param)
			continue
		}

		params[pos] = utils.Iprint(param)
	}
}

func Init(module, subcate string, level int, mode int) {
	if strings.Contains(module, ",") || strings.Contains(subcate, ",") {
		panic("clog Init error, module or subcate contains ','")
	}

	cate_dbg = strings.Join([]string{module, "logdbg", utils.LocalIp, subcate}, ",")
	cate_war = strings.Join([]string{module, "logwar", utils.LocalIp, subcate}, ",")
	cate_err = strings.Join([]string{module, "logerr", utils.LocalIp, subcate}, ",")
	cate_info = strings.Join([]string{module, "loginfo", utils.LocalIp, subcate}, ",")
	cate_busi = strings.Join([]string{module, "logbusi_%s", utils.LocalIp, subcate}, ",")

	Level = level
	Mode = mode
}

func Debug(format string, params ...interface{}) {
	if Level&1 != 0 {
		iprint(params)
		content := fmt.Sprintf(format, params...)
		if Mode&1 != 0 {
			log.Println("[DEBUG]", content)
		}
		if Mode&2 != 0 {
			sendAgent(cate_dbg, content)
		}
	}
}

func Warn(format string, params ...interface{}) {
	if Level&2 != 0 {
		iprint(params)
		content := fmt.Sprintf(format, params...)
		if Mode&1 != 0 {
			log.Println("[WARN]", content)
		}
		if Mode&2 != 0 {
			sendAgent(cate_war, content)
		}
	}
}

func Error(format string, params ...interface{}) {
	if Level&4 != 0 {
		iprint(params)
		content := fmt.Sprintf(format, params...)
		if Mode&1 != 0 {
			log.Println("[ERROR]", content)
		}
		if Mode&2 != 0 {
			sendAgent(cate_err, content)
		}
	}
}

func Info(format string, params ...interface{}) {
	if Level&8 != 0 {
		iprint(params)
		content := fmt.Sprintf(format, params...)
		if Mode&1 != 0 {
			log.Println("[INFO]", content)
		}
		if Mode&2 != 0 {
			sendAgent(cate_info, content)
		}
	}
}

func Busi(sub string, format string, params ...interface{}) {
	iprint(params)
	content := fmt.Sprintf(format, params...)
	if Mode&1 != 0 {
		log.Println("[BUSI]", sub, content)
	}
	if Mode&2 != 0 {
		sendAgent(fmt.Sprintf(cate_busi, sub), content)
	}
}
