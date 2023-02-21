package server

func Init(apiToken string) {
	r := NewRouter(apiToken)
	r.Run(":5000")
}
