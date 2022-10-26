// Copyright (c) 2021. Homeland Interactive Technology Ltd. All rights reserved.

package main

import (
	"fmt"
	"git.jiaxianghudong.com/ruixuesdk/ruixuego"
	"time"
)

const (
	testAppID     = "wltestapp"
	testChannelID = "wltestchannel"
	testAppKey    = "a463deade4b15d5ac5398f97cdaeab65"
	rankID        = "0_200_202_dayly"
)

func main() {
	err := ruixuego.Init(&ruixuego.Config{
		APIDomain: "http://ruixue.weiletest.com",
		AppKeys:   map[string]map[string]string{testAppID: {testChannelID: testAppKey}},
		CPKey:     "0984cde09ebe42fd167510c727f57f71",
		CPID:      1000049,
	})
	if err != nil {
		panic(err)
	}
	err = ruixuego.GetDefaultClient().SetCustom(testAppID, "rxuR4bwM27Y1JQwtAQn6H39y_9VrkEgR", "123")
	if err != nil {
		panic(err)
	}

	err = ruixuego.GetDefaultClient().AddFriend("rxuR4bwM27Y1JQwtAQn6H39y_9VrkEgR", "rxufPJeZoyrX3eWuNxLMSNK5N6x04jLl", "aaaa1", "bbb13")
	if err != nil {
		panic(err)
	}

	ret, err := ruixuego.GetDefaultClient().FriendList("rxuR4bwM27Y1JQwtAQn6H39y_9VrkEgR")
	if err != nil {
		panic(err)
	}
	for _, ru := range ret {
		fmt.Printf("%+v\n", ru)
	}

	yes, err := ruixuego.GetDefaultClient().IsFriend("rxuR4bwM27Y1JQwtAQn6H39y_9VrkEgR", "rxufPJeZoyrX3eWuNxLMSNK5N6x04jLl")
	if err != nil {
		panic(err)
	}
	fmt.Println(yes)

	// err = ruixuego.GetDefaultClient().DelFriend("rxuR4bwM27Y1JQwtAQn6H39y_9VrkEgR", "rxufPJeZoyrX3eWuNxLMSNK5N6x04jLl")
	// if err != nil {
	// 	panic(err)
	// }

	err = ruixuego.GetDefaultClient().CreateRank(rankID, time.Now(), time.Now().Add(10*time.Hour*24))
	if err != nil {
		panic(err)
	}

	err = ruixuego.GetDefaultClient().RankAddScore(rankID, "rxufb2GKDBm5n4u7gYkcOLbOv7J4RrVP", 1000)
	if err != nil {
		panic(err)
	}

	err = ruixuego.GetDefaultClient().RankSetScore(rankID, "rxuSl4QZoNk0G1HY2-Za6GlO7wO-p_ej", 5000)
	if err != nil {
		panic(err)
	}

	rankMember, err := ruixuego.GetDefaultClient().QueryUserRank(rankID, "rxuSl4QZoNk0G1HY2-Za6GlO7wO-p_ej")
	if err != nil {
		panic(err)
	}
	fmt.Printf("QueryUserRank --> rank:%+v\n", rankMember)

	rankList, err := ruixuego.GetDefaultClient().GetRankList(rankID)
	if err != nil {
		panic(err)
	}
	for _, rank := range rankList {
		fmt.Printf("rankList --> rank:%+v\n", rank)
	}

	err = ruixuego.GetDefaultClient().CloseRank(rankID)
	if err != nil {
		panic(err)
	}

	fmt.Println("done")
}
