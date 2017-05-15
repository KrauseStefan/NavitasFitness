
cd websrc
& npm install
& npm -new_console start
cd ..

cd ipn-simulator
& npm -new_console start
cd ..

& go get -v ./...
& dev_appserver.py -new_console --dev_appserver_log_level=warning .

