imageName="registry.cn-qingdao.aliyuncs.com/wisati/paddle"
containerName="${imageName}"

docker run -it  -v "/Users/wyn/Library/Mobile Documents/com~apple~CloudDocs/Daocloud/GOProject/src/github.com/IsolationWyn/paddle":/go/src/github.com/IsolationWyn/paddle --privileged -p 2345:2345 $imageName

