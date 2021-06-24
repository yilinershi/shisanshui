package tradition_func

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"
)

type pokerInfo struct {
	Pokers    string
	Desc      string
	IsTest    bool
	TestPoker []*Poker
}

func GenTestPokers(Pokers string) []*Poker {
	testPokers := make([]*Poker, 0)
	pokerStrList := strings.Split(Pokers, ",")
	for _, str := range pokerStrList {
		huaDesc := str[0:3]
		hua := DescHua(huaDesc).ToPokerHua()
		pointDesc := str[3:]
		point := DescPoint(pointDesc).ToPokerPoint()
		poker := NewPoker(point, hua)
		testPokers = append(testPokers, poker)
	}
	return testPokers
}

func TestApi(t *testing.T) {

	//读取json文件中的测试表
	jsonFile, err := os.Open("../TestCase.json")
	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var infos []pokerInfo

	if err := json.Unmarshal(byteValue, &infos); err != nil {
		fmt.Println("Unmarshal json error, err=", err)
		return
	}

	for _, info := range infos {
		info.TestPoker = GenTestPokers(info.Pokers)
		if info.IsTest {
			if len(info.TestPoker) == 13 {
				startTime := time.Now().Nanosecond()
				fmt.Printf("测试组合=%s,牌型={%s},开始时间=%d\n", info.Desc, info.Pokers, startTime)
				CalResult(info.TestPoker)
				endTime := time.Now().Nanosecond()
				costTime := endTime - startTime
				fmt.Printf("结束时间=%d,AI算法耗时【%d】微秒\n\n", endTime, costTime/1000)
			} else if len(info.TestPoker) == 5 || len(info.TestPoker) == 3 {
				startTime := time.Now().Nanosecond()
				fmt.Printf("测试组合=%s,牌型={%s},开始时间=%d\n", info.Desc, info.Pokers, startTime)
				d := NewDun(info.TestPoker)
				endTime := time.Now().Nanosecond()
				costTime := endTime - startTime
				fmt.Printf("结束时间=%d,计算牌型=【%s】，AI算法耗时【%d】微秒\n\n", endTime, d.Type, costTime/1000)
			}
		}
	}
}


