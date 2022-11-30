package debitsam

//go:generate swagger generate server --exclude-main --skip-validation -A debitsam-service -t gen -f ./api/spec.yml  --principal models.Principal
//go:generate swagger -q generate client -A debitsam-service -f api/spec.yml -c pkg/client -m gen/models --principal models.Principal
