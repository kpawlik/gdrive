module github.com/kpawlik/gdrive/uploadGD

require (
	cloud.google.com/go v0.33.1
	github.com/kpawlik/gdrive/gdrive v0.0.0
	golang.org/x/net v0.0.0-20181129055619-fae4c4e3ad76
	golang.org/x/oauth2 v0.0.0-20181128211412-28207608b838
	google.golang.org/api v0.0.0-20181129220737-af4fc4062c26
)

replace github.com/kpawlik/gdrive/gdrive => ../gdrive
