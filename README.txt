# mockgen -destination=mocks/mock_doer.go -package=mocks github.com/sgreben/testing-with-gomock/doer Doer

https://blog.codecentric.de/en/2017/08/gomock-tutorial/

//go:generate 

mockgen -destination=mocks/repo/mock_repo.go -package=mocks code.mine/dating_server/repo Repo