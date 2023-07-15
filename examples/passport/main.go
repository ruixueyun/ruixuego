package passport

import "github.com/ruixueyun/ruixuego"

func main() {
	err := ruixuego.Init(&ruixuego.Config{
		APIDomain: "https://domain.com",
		CPKey:     "00000000000000000000000",
		CPID:      1000005,
		ProductID: "425",
	})
	if err != nil {
		panic(err)
	}
	ruixuego.GetDefaultClient().UpdateCPuserID("rxOpenID", "cpUserID")
}
