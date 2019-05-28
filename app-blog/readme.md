1. First install dep
    curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

2. After that, go into project dir like /home/deepak/go/src/github.com/deepaksinghkushwah/blog

3. To add package, issue following command
    dep ensure -add "package name"
                like
    dep ensure -add "golang.org/x/crypto/bcrypt"
