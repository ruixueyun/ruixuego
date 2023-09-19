// Copyright (c) 2021. Homeland Interactive Technology Ltd. All rights reserved.

package main

import (
	"fmt"
	"github.com/ruixueyun/ruixuego"
)

const (
	testAppID     = "wltestapp"
	testChannelID = "wltestchannel"
	testAppKey    = "00000000000000000000000"
	rankID        = "0_200_202_dayly"
)

func main() {
	err := ruixuego.Init(&ruixuego.Config{
		APIDomain:   "https://domain.com",
		AppKeys:     map[string]map[string]string{testAppID: {testChannelID: testAppKey}},
		CPKey:       "00000000000000000000000",
		CPID:        1000049,
		ServiceMark: "aabbcc",
	})
	if err != nil {
		panic(err)
	}
	// err = ruixuego.GetDefaultClient().SetCustom(testAppID, "rxuR4bwM27Y1JQwtAQn6H39y_9VrkEgR", "123")
	// if err != nil {
	//	panic(err)
	// }
	//
	// err = ruixuego.GetDefaultClient().AddFriend("rxuR4bwM27Y1JQwtAQn6H39y_9VrkEgR", "rxufPJeZoyrX3eWuNxLMSNK5N6x04jLl", "aaaa1", "bbb13")
	// if err != nil {
	//	panic(err)
	// }
	//
	// ret, err := ruixuego.GetDefaultClient().FriendList("rxuR4bwM27Y1JQwtAQn6H39y_9VrkEgR")
	// if err != nil {
	//	panic(err)
	// }
	// for _, ru := range ret {
	//	fmt.Printf("%+v\n", ru)
	// }
	//
	// yes, err := ruixuego.GetDefaultClient().IsFriend("rxuR4bwM27Y1JQwtAQn6H39y_9VrkEgR", "rxufPJeZoyrX3eWuNxLMSNK5N6x04jLl")
	// if err != nil {
	//	panic(err)
	// }
	// fmt.Println(yes)

	// err = ruixuego.GetDefaultClient().DelFriend("rxuR4bwM27Y1JQwtAQn6H39y_9VrkEgR", "rxufPJeZoyrX3eWuNxLMSNK5N6x04jLl")
	// if err != nil {
	// 	panic(err)
	// }

	// err = ruixuego.GetDefaultClient().CreateRank(rankID, time.Now(), time.Now().Add(10*time.Hour*24))
	// if err != nil {
	//	panic(err)
	// }
	//
	//ret, err := ruixuego.GetDefaultClient().IMSChannelUsersCount([]string{"$4$worldChannel"})
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Printf("IMSChannelUsersCount --> ret:%+v\n", ret)

	//
	// err = ruixuego.GetDefaultClient().RankSetScore(rankID, "rxuSl4QZoNk0G1HY2-Za6GlO7wO-p_ej", 5000)
	// if err != nil {
	//	panic(err)
	// }
	//
	// rankMember, err := ruixuego.GetDefaultClient().QueryUserRank(rankID, "rxuSl4QZoNk0G1HY2-Za6GlO7wO-p_ej")
	// if err != nil {
	//	panic(err)
	// }
	// fmt.Printf("QueryUserRank --> rank:%+v\n", rankMember)
	//
	// rankList, err := ruixuego.GetDefaultClient().GetRankList(rankID)
	// if err != nil {
	//	panic(err)
	// }
	// for _, rank := range rankList {
	//	fmt.Printf("rankList --> rank:%+v\n", rank)
	// }
	//
	// err = ruixuego.GetDefaultClient().CloseRank(rankID)
	// if err != nil {
	//	panic(err)
	// }

	//var loginOpenId = `HTFqREbcznvBQiUyANGmQ56T5KuGxoYw9EW/3kJiOhHeijQuQfePua6B5wT4V4NElFDtzDBiey/JWmdUWPAjgdvZizKXPN5UUW1cJP/mZqBfhdMxE0r1p6wyNOWB7960h5v4y2XHDJZhwihc3peES9FPSbGgpxvP90C7ASO6QmQPXBHQeP4xN33CLsM8BCplzWWNiiBCNHp8P3R+F/cxIrmPTQvuG9YkzZffjyMf9NwwdSz5i5lR35i7Lt4fwg1zwXPnNAWY4PtC2r5LwvyS3lUfm/wjjmUlMYwDfRrPP1VTLUdGQw5Sd7vx4lu73BrPdwOExkYlWpRVsfFjZP6Dyg==`
	////var openId = "rxuCAMkg_9QojLHlybEClG74woeB6EK7"
	//ruixuego.AddAESKey("100", "100", []byte("a0ac930adc99fed5b3467a47d757f907"))
	//openIdData, err := ruixuego.DecryptOpenIDData("100", "100", loginOpenId)
	//if err != nil {
	//	panic(err)
	//}

	tasks := make([]*ruixuego.GreenRequestTask, 0, 1)
	task := &ruixuego.GreenRequestTask{
		Tag: 123456,
		URL: "https://avatarfile.weile.com/temp/523872878-1687165437",
	}
	tasks = append(tasks, task)

	resp, err := ruixuego.GetDefaultClient().RiskGreenStrongScan([]string{"baselineCheck_pro"}, tasks, "")
	if err != nil {
		panic(err)
	}
	fmt.Println("resp +", resp.Results[0])
	//
	//fmt.Println("done")

	//resp, err := ruixuego.GetDefaultClient().RiskSensitive("你好我是毛泽东")
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println("resp +", resp)

}
