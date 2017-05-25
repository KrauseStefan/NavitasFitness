
SF="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
DIR="$SF/App"

cd $DIR

GO_TEST_PATH="$HOME/TEST-GOPATH"

rm $GO_TEST_PATH -rf

mkdir -p $GO_TEST_PATH
mkdir -p "$GO_TEST_PATH/src"

for dir in $(ls -d */); do
    if [[ $dir =~ ^[A-Z].* ]]; then
        ln -sf "$DIR/${dir::-1}" "$GO_TEST_PATH/src/"
    fi;
done

export GOPATH="$GOPATH:$GO_TEST_PATH"

gcloud app deploy NavitasFitness/app.yaml
