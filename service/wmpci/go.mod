module github.com/suitcase/butler/wmpci

go 1.24.1

require (
	github.com/google/uuid v1.6.0
	github.com/sirupsen/logrus v1.9.3
	github.com/suitcase/butler v0.0.0-00010101000000-000000000000
	go.mongodb.org/mongo-driver v1.17.1
	gopkg.in/yaml.v2 v2.4.0
)

require github.com/stretchr/testify v1.9.0 // indirect

require (
	github.com/golang/snappy v0.0.4 // indirect
	github.com/klauspost/compress v1.17.4 // indirect
	github.com/montanaflynn/stats v0.7.1 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20240726163527-a2c0da244d78 // indirect; indirect·
	golang.org/x/crypto v0.36.0 // indirect
	golang.org/x/sync v0.12.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	golang.org/x/text v0.23.0 // indirect
)

replace github.com/suitcase/butler => ../

replace github.com/suitcase/butler/built_wmp => ../built_wmp
