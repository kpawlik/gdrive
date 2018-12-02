module github.com/kpawlik/gdrive/downloadGD

require (
	cloud.google.com/go v0.33.1 // indirect
	github.com/kpawlik/gdrive/gdrive v0.0.0
	golang.org/x/net v0.0.0-20181201002055-351d144fa1fc
	golang.org/x/oauth2 v0.0.0-20181128211412-28207608b838
	google.golang.org/api v0.0.0-20181129220737-af4fc4062c26
)

replace github.com/kpawlik/gdrive/gdrive => ../gdrive
