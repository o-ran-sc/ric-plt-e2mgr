module e2mgr

require (
	gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common v1.2.1
	gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities v1.2.1
	gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/reader v1.2.1
	gerrit.o-ran-sc.org/r/ric-plt/sdlgo v0.8.0
	github.com/golang/protobuf v1.4.2
	github.com/gorilla/mux v1.7.0
	github.com/kr/pretty v0.2.0 // indirect
	github.com/magiconair/properties v1.8.1
	github.com/mitchellh/mapstructure v1.3.2 // indirect
	github.com/pelletier/go-toml v1.8.0 // indirect
	github.com/pkg/errors v0.8.1
	github.com/spf13/afero v1.2.2 // indirect
	github.com/spf13/cast v1.3.1 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/spf13/viper v1.4.0
	github.com/stretchr/objx v0.2.0 // indirect
	github.com/stretchr/testify v1.5.1
	go.uber.org/multierr v1.2.0 // indirect
	go.uber.org/zap v1.11.0
	golang.org/x/sys v0.0.0-20200602225109-6fdc65e7d980 // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/yaml.v2 v2.3.0
)

replace gerrit.o-ran-sc.org/r/ric-plt/sdlgo => gerrit.o-ran-sc.org/r/ric-plt/sdlgo.git v0.8.0

go 1.13
