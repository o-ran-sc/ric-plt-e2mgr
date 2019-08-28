module e2mgr

require (
	gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common v1.0.19
	gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities v1.0.19
	gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/reader v1.0.19
	gerrit.o-ran-sc.org/r/ric-plt/sdlgo v0.2.0
	github.com/asaskevich/govalidator v0.0.0-20190424111038-f61b66f89f4a // indirect
	github.com/go-ozzo/ozzo-validation v3.5.0+incompatible
	github.com/golang/protobuf v1.3.1
	github.com/gorilla/context v1.1.1 // indirect
	github.com/gorilla/mux v1.6.2
	github.com/pkg/errors v0.8.1
	github.com/spf13/viper v1.4.0
	github.com/stretchr/testify v1.3.0
	go.uber.org/zap v1.10.0
	gopkg.in/yaml.v2 v2.2.2
)

replace gerrit.o-ran-sc.org/r/ric-plt/sdlgo => gerrit.o-ran-sc.org/r/ric-plt/sdlgo.git v0.2.0
