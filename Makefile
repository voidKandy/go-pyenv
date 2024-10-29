test-integration:
	go test ./pyenv -run TestIntegration  -count=1 -v
test-dependencies:
	go test ./pyenv -run TestDependencies  -count=1 -v
test-compression:
	go test ./pyenv -run TestCompression  -count=1 -v
test-all:
	go test ./pyenv -run "(TestIntegration|TestDependencies)"  -count=1 -v
test-remove:
	go test ./pyenv -run TestRemove  -count=1 -v