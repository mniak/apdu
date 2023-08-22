module github.com/mniak/apdu

go 1.20

require (
	github.com/brianvoe/gofakeit/v6 v6.23.0
	github.com/ebfe/scard v0.0.0-20230420082256-7db3f9b7c8a7
	github.com/mniak/krypton v0.0.0-20230721155408-50f12342b13f
	github.com/pkg/errors v0.9.1
	github.com/samber/lo v1.38.1
	github.com/stretchr/testify v1.8.4
)

require (
	github.com/asaskevich/govalidator v0.0.0-20230301143203-a9d515a09cc2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-ozzo/ozzo-validation/v4 v4.3.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/exp v0.0.0-20220303212507-bbda1eaf7a17 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/mniak/krypton => ../krypton
